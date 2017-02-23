package hoverfly

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/SpectoLabs/goproxy"
	authBackend "github.com/SpectoLabs/hoverfly/core/authentication/backends"
	"github.com/SpectoLabs/hoverfly/core/cache"
	"github.com/SpectoLabs/hoverfly/core/matching"
	"github.com/SpectoLabs/hoverfly/core/metrics"
	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/SpectoLabs/hoverfly/core/modes"
)

// orPanic - wrapper for logging errors
func orPanic(err error) {
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Panic("Got error.")
	}
}

// Hoverfly provides access to hoverfly - updating/starting/stopping proxy, http client and configuration, cache access
type Hoverfly struct {
	RequestCache   cache.Cache
	RequestMatcher matching.RequestMatcher
	MetadataCache  cache.Cache
	Authentication authBackend.Authentication
	HTTP           *http.Client
	Cfg            *Configuration
	Counter        *metrics.CounterByMode

	ResponseDelays models.ResponseDelays

	Proxy   *goproxy.ProxyHttpServer
	SL      *StoppableListener
	mu      sync.Mutex
	version string

	modeMap map[string]modes.Mode

	Simulation *models.Simulation
}

// GetNewHoverfly returns a configured ProxyHttpServer and DBClient
func GetNewHoverfly(cfg *Configuration, requestCache, metadataCache cache.Cache, authentication authBackend.Authentication) *Hoverfly {
	simulation := models.NewSimulation()

	requestMatcher := matching.RequestMatcher{
		RequestCache: requestCache,
		Webserver:    &cfg.Webserver,
		Simulation:   simulation,
	}

	h := &Hoverfly{
		RequestCache:   requestCache,
		MetadataCache:  metadataCache,
		Authentication: authentication,
		HTTP:           GetDefaultHoverflyHTTPClient(cfg.TLSVerification, cfg.UpstreamProxy),
		Cfg:            cfg,
		Counter:        metrics.NewModeCounter([]string{modes.Simulate, modes.Synthesize, modes.Modify, modes.Capture}),
		ResponseDelays: &models.ResponseDelayList{},
		RequestMatcher: requestMatcher,
		Simulation:     simulation,
	}

	modeMap := make(map[string]modes.Mode)

	modeMap[modes.Capture] = modes.CaptureMode{Hoverfly: h}
	modeMap[modes.Simulate] = modes.SimulateMode{Hoverfly: h}
	modeMap[modes.Modify] = modes.ModifyMode{Hoverfly: h}
	modeMap[modes.Synthesize] = modes.SynthesizeMode{Hoverfly: h}

	h.modeMap = modeMap

	h.version = "v0.10.2"

	return h
}

func GetDefaultHoverflyHTTPClient(tlsVerification bool, upstreamProxy string) *http.Client {

	var proxyURL func(*http.Request) (*url.URL, error)
	if upstreamProxy == "" {
		proxyURL = http.ProxyURL(nil)
	} else {
		u, err := url.Parse(upstreamProxy)
		if err != nil {
			log.Fatalf("Could not parse upstream proxy: ", err.Error())
		}
		proxyURL = http.ProxyURL(u)
	}

	return &http.Client{CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}, Transport: &http.Transport{
		Proxy:           proxyURL,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: tlsVerification},
	}}
}

// StartProxy - starts proxy with current configuration, this method is non blocking.
func (hf *Hoverfly) StartProxy() error {

	rebuildHashes(hf.RequestCache, hf.Cfg.Webserver)

	if hf.Cfg.ProxyPort == "" {
		return fmt.Errorf("Proxy port is not set!")
	}

	if hf.Cfg.Webserver {
		hf.Proxy = NewWebserverProxy(hf)
	} else {
		hf.Proxy = NewProxy(hf)
	}

	log.WithFields(log.Fields{
		"destination": hf.Cfg.Destination,
		"port":        hf.Cfg.ProxyPort,
		"mode":        hf.Cfg.GetMode(),
	}).Info("current proxy configuration")

	// creating TCP listener
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", hf.Cfg.ProxyPort))
	if err != nil {
		return err
	}

	sl, err := NewStoppableListener(listener)
	if err != nil {
		return err
	}
	hf.SL = sl
	server := http.Server{}

	hf.Cfg.ProxyControlWG.Add(1)

	go func() {
		defer func() {
			log.Info("sending done signal")
			hf.Cfg.ProxyControlWG.Done()
		}()
		log.Info("serving proxy")
		server.Handler = hf.Proxy
		log.Warn(server.Serve(sl))
	}()

	return nil
}

// StopProxy - stops proxy
func (hf *Hoverfly) StopProxy() {
	hf.SL.Stop()
	hf.Cfg.ProxyControlWG.Wait()
}

// processRequest - processes incoming requests and based on proxy state (record/playback)
// returns HTTP response.
func (hf *Hoverfly) processRequest(req *http.Request) *http.Response {
	requestDetails, err := models.NewRequestDetailsFromHttpRequest(req)
	if err != nil {
		return modes.ErrorResponse(req, err, "Could not interpret HTTP request")
	}

	mode := hf.Cfg.GetMode()

	response, err := hf.modeMap[mode].Process(req, requestDetails)

	// Don't delete the error
	// and definitely don't delay people in capture mode
	if err != nil || mode == modes.Capture {
		return response
	}

	respDelay := hf.ResponseDelays.GetDelay(requestDetails)
	if respDelay != nil {
		respDelay.Execute()
	}

	return response
}

// DoRequest - performs request and returns response that should be returned to client and error
func (hf *Hoverfly) DoRequest(request *http.Request) (*http.Response, error) {

	// We can't have this set. And it only contains "/pkg/net/http/" anyway
	request.RequestURI = ""

	requestBody, _ := ioutil.ReadAll(request.Body)

	request.Body = ioutil.NopCloser(bytes.NewReader(requestBody))

	resp, err := hf.HTTP.Do(request)

	request.Body = ioutil.NopCloser(bytes.NewReader(requestBody))
	if err != nil {
		return nil, err
	}

	resp.Header.Set("hoverfly", "Was-Here")

	return resp, nil

}

// GetResponse returns stored response from cache
func (hf *Hoverfly) GetResponse(requestDetails models.RequestDetails) (*models.ResponseDetails, *matching.MatchingError) {
	return hf.RequestMatcher.GetResponse(&requestDetails)
}

// save gets request fingerprint, extracts request body, status code and headers, then saves it to cache
func (hf *Hoverfly) Save(request *models.RequestDetails, response *models.ResponseDetails) error {

	pair := models.RequestResponsePair{
		Request:  *request,
		Response: *response,
	}

	err := hf.RequestMatcher.SaveRequestResponsePair(&pair)
	if err != nil {
		return err
	}

	return nil
}

func (this Hoverfly) ApplyMiddleware(pair models.RequestResponsePair) (models.RequestResponsePair, error) {
	if this.Cfg.Middleware.IsSet() {
		return this.Cfg.Middleware.Execute(pair)
	}

	return pair, nil
}

func (this Hoverfly) IsMiddlewareSet() bool {
	return this.Cfg.Middleware.IsSet()
}

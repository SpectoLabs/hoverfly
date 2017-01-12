package hoverfly

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	authBackend "github.com/SpectoLabs/hoverfly/core/authentication/backends"
	"github.com/SpectoLabs/hoverfly/core/cache"
	"github.com/SpectoLabs/hoverfly/core/matching"
	"github.com/SpectoLabs/hoverfly/core/metrics"
	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/SpectoLabs/hoverfly/core/modes"
	"github.com/rusenask/goproxy"
)

// SimulateMode - default mode when Hoverfly looks for captured requests to respond
const SimulateMode = "simulate"

// SynthesizeMode - all requests are sent to middleware to create response
const SynthesizeMode = "synthesize"

// ModifyMode - middleware is applied to outgoing and incoming traffic
const ModifyMode = "modify"

// CaptureMode - requests are captured and stored in cache
const CaptureMode = "capture"

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
	Hooks          ActionTypeHooks

	ResponseDelays models.ResponseDelays

	Proxy *goproxy.ProxyHttpServer
	SL    *StoppableListener
	mu    sync.Mutex

	modeMap map[string]modes.Mode
}

// GetNewHoverfly returns a configured ProxyHttpServer and DBClient
func GetNewHoverfly(cfg *Configuration, requestCache, metadataCache cache.Cache, authentication authBackend.Authentication) *Hoverfly {
	requestMatcher := matching.RequestMatcher{
		RequestCache:  requestCache,
		TemplateStore: matching.RequestTemplateStore{},
		Webserver:     &cfg.Webserver,
	}

	h := &Hoverfly{
		RequestCache:   requestCache,
		MetadataCache:  metadataCache,
		Authentication: authentication,
		HTTP:           GetDefaultHoverflyHTTPClient(cfg.TLSVerification),
		Cfg:            cfg,
		Counter:        metrics.NewModeCounter([]string{SimulateMode, SynthesizeMode, ModifyMode, CaptureMode}),
		Hooks:          make(ActionTypeHooks),
		ResponseDelays: &models.ResponseDelayList{},
		RequestMatcher: requestMatcher,
	}

	modeMap := make(map[string]modes.Mode)

	modeMap["capture"] = modes.CaptureMode{Hoverfly: h}
	modeMap["simulate"] = modes.SimulateMode{Hoverfly: h}
	modeMap["modify"] = modes.ModifyMode{Hoverfly: h}
	modeMap["synthesize"] = modes.SynthesizeMode{Hoverfly: h}

	h.modeMap = modeMap

	return h
}

func GetDefaultHoverflyHTTPClient(tlsVerification bool) *http.Client {
	return &http.Client{CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}, Transport: &http.Transport{
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
	if err != nil || mode == CaptureMode {
		return response
	}

	respDelay := hf.ResponseDelays.GetDelay(requestDetails)
	if respDelay != nil {
		respDelay.Execute()
	}

	return response
}

// AddHook - adds a hook to DBClient
func (hf *Hoverfly) AddHook(hook Hook) {
	hf.Hooks.Add(hook)
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
		log.WithFields(log.Fields{
			"mode":   hf.Cfg.Mode,
			"host":   request.Host,
			"method": request.Method,
			"path":   request.URL.Path,
		}).Error("HTTP request failed: " + err.Error())
		return nil, err
	}

	log.WithFields(log.Fields{
		"mode":   hf.Cfg.Mode,
		"host":   request.Host,
		"method": request.Method,
		"path":   request.URL.Path,
	}).Debug("response from external service got successfuly!")

	resp.Header.Set("hoverfly", "Was-Here")

	return resp, nil

}

// GetResponse returns stored response from cache
func (hf *Hoverfly) GetResponse(requestDetails models.RequestDetails) (*models.ResponseDetails, *matching.MatchingError) {
	return hf.RequestMatcher.GetResponse(&requestDetails)
}

// save gets request fingerprint, extracts request body, status code and headers, then saves it to cache
func (hf *Hoverfly) Save(request *models.RequestDetails, response *models.ResponseDetails) {

	pair := models.RequestResponsePair{
		Request:  *request,
		Response: *response,
	}

	err := hf.RequestMatcher.SaveRequestResponsePair(&pair)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Failed to save payload")
	}

	pairBytes, err := pair.Encode()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Failed to serialize payload")
	} else {
		// hook
		var en Entry
		en.ActionType = ActionTypeRequestCaptured
		en.Message = "captured"
		en.Time = time.Now()
		en.Data = pairBytes

		if err := hf.Hooks.Fire(ActionTypeRequestCaptured, &en); err != nil {
			log.WithFields(log.Fields{
				"error":      err.Error(),
				"message":    en.Message,
				"actionType": ActionTypeRequestCaptured,
			}).Error("failed to fire hook")
		}
	}

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

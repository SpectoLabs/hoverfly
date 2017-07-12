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

	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/SpectoLabs/goproxy"
	"github.com/SpectoLabs/hoverfly/core/authentication/backends"
	"github.com/SpectoLabs/hoverfly/core/cache"
	"github.com/SpectoLabs/hoverfly/core/journal"
	"github.com/SpectoLabs/hoverfly/core/matching"
	"github.com/SpectoLabs/hoverfly/core/metrics"
	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/SpectoLabs/hoverfly/core/modes"
	"github.com/SpectoLabs/hoverfly/core/templating"
	"github.com/SpectoLabs/hoverfly/core/util"
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
	CacheMatcher   matching.CacheMatcher
	MetadataCache  cache.Cache
	Authentication backends.Authentication
	HTTP           *http.Client
	Cfg            *Configuration
	Counter        *metrics.CounterByMode

	Proxy   *goproxy.ProxyHttpServer
	SL      *StoppableListener
	mu      sync.Mutex
	version string

	modeMap map[string]modes.Mode

	Simulation    *models.Simulation
	StoreLogsHook *StoreLogsHook
	Journal       *journal.Journal
}

func NewHoverfly() *Hoverfly {
	authBackend := backends.NewCacheBasedAuthBackend(cache.NewInMemoryCache(), cache.NewInMemoryCache())

	hoverfly := &Hoverfly{
		Simulation:     models.NewSimulation(),
		Authentication: authBackend,
		Counter:        metrics.NewModeCounter([]string{modes.Simulate, modes.Synthesize, modes.Modify, modes.Capture}),
		StoreLogsHook:  NewStoreLogsHook(),
		Journal:        journal.NewJournal(),
		Cfg:            InitSettings(),
	}

	hoverfly.version = "v0.13.0"

	log.AddHook(hoverfly.StoreLogsHook)

	modeMap := make(map[string]modes.Mode)

	modeMap[modes.Capture] = &modes.CaptureMode{Hoverfly: hoverfly}
	modeMap[modes.Simulate] = &modes.SimulateMode{Hoverfly: hoverfly, MatchingStrategy: "strongest"}
	modeMap[modes.Modify] = &modes.ModifyMode{Hoverfly: hoverfly}
	modeMap[modes.Synthesize] = &modes.SynthesizeMode{Hoverfly: hoverfly}

	hoverfly.modeMap = modeMap

	hoverfly.HTTP = GetDefaultHoverflyHTTPClient(hoverfly.Cfg.TLSVerification, hoverfly.Cfg.UpstreamProxy)

	return hoverfly
}

func NewHoverflyWithConfiguration(cfg *Configuration) *Hoverfly {
	hoverfly := NewHoverfly()

	var requestCache cache.Cache
	if !cfg.DisableCache {
		requestCache = cache.NewInMemoryCache()
	}

	hoverfly.MetadataCache = cache.NewInMemoryCache()

	hoverfly.CacheMatcher = matching.CacheMatcher{
		RequestCache: requestCache,
		Webserver:    cfg.Webserver,
	}

	hoverfly.Cfg = cfg
	hoverfly.HTTP = GetDefaultHoverflyHTTPClient(cfg.TLSVerification, cfg.UpstreamProxy)

	return hoverfly
}

// GetNewHoverfly returns a configured ProxyHttpServer and DBClient
func GetNewHoverfly(cfg *Configuration, requestCache, metadataCache cache.Cache, authentication backends.Authentication) *Hoverfly {
	hoverfly := NewHoverfly()

	if cfg.DisableCache {
		requestCache = nil
	}

	hoverfly.CacheMatcher = matching.CacheMatcher{
		RequestCache: requestCache,
		Webserver:    cfg.Webserver,
	}

	hoverfly.MetadataCache = metadataCache
	hoverfly.Authentication = authentication
	hoverfly.HTTP = GetDefaultHoverflyHTTPClient(cfg.TLSVerification, cfg.UpstreamProxy)
	hoverfly.Cfg = cfg

	return hoverfly
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
		TLSClientConfig: &tls.Config{InsecureSkipVerify: !tlsVerification},
	}}
}

// StartProxy - starts proxy with current configuration, this method is non blocking.
func (hf *Hoverfly) StartProxy() error {

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

	respDelay := hf.Simulation.ResponseDelays.GetDelay(requestDetails)
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

	cachedResponse, cacheErr := hf.CacheMatcher.GetCachedResponse(&requestDetails)
	if cacheErr == nil && cachedResponse.MatchingPair == nil {
		return nil, matching.MissedError(cachedResponse.ClosestMiss)
	} else if cacheErr == nil {
		return &cachedResponse.MatchingPair.Response, nil
	}

	var pair *models.RequestMatcherResponsePair
	var err *models.MatchError

	mode := (hf.modeMap[modes.Simulate]).(*modes.SimulateMode)

	strongestMatch := strings.ToLower(mode.MatchingStrategy) == "strongest"

	// Matching
	if strongestMatch {
		pair, err = matching.StrongestMatchRequestMatcher(requestDetails, hf.Cfg.Webserver, hf.Simulation)
	} else {
		pair, err = matching.FirstMatchRequestMatcher(requestDetails, hf.Cfg.Webserver, hf.Simulation)
	}

	// Templating
	if err == nil && pair.Response.Templated == true {
		responseBody, err := templating.ApplyTemplate(&requestDetails, pair.Response.Body)
		if err == nil {
			pair.Response.Body = responseBody
		}
	}

	hf.CacheMatcher.SaveRequestMatcherResponsePair(requestDetails, pair, err)

	if err != nil {
		log.WithFields(log.Fields{
			"error":       err.Error(),
			"query":       requestDetails.Query,
			"path":        requestDetails.Path,
			"destination": requestDetails.Destination,
			"method":      requestDetails.Method,
		}).Warn("Failed to find matching request from simulation")

		return nil, matching.MissedError(err.ClosestMiss)
	}

	return &pair.Response, nil
}

// save gets request fingerprint, extracts request body, status code and headers, then saves it to cache
func (hf *Hoverfly) Save(request *models.RequestDetails, response *models.ResponseDetails, headersWhitelist []string) error {
	body := &models.RequestFieldMatchers{
		ExactMatch: util.StringToPointer(request.Body),
	}
	contentType := util.GetContentTypeFromHeaders(request.Headers)
	if contentType == "json" {
		body = &models.RequestFieldMatchers{
			JsonMatch: util.StringToPointer(request.Body),
		}
	} else if contentType == "xml" {
		body = &models.RequestFieldMatchers{
			XmlMatch: util.StringToPointer(request.Body),
		}
	}

	var headers map[string][]string
	if headersWhitelist == nil {
		headersWhitelist = []string{}
	}

	if len(headersWhitelist) >= 1 && headersWhitelist[0] == "*" {
		headers = request.Headers
	} else {
		headers = map[string][]string{}
		for _, header := range headersWhitelist {
			headerValues := request.Headers[header]
			if len(headerValues) > 0 {
				headers[header] = headerValues
			}
		}
	}

	pair := models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Path: &models.RequestFieldMatchers{
				ExactMatch: util.StringToPointer(request.Path),
			},
			Method: &models.RequestFieldMatchers{
				ExactMatch: util.StringToPointer(request.Method),
			},
			Destination: &models.RequestFieldMatchers{
				ExactMatch: util.StringToPointer(request.Destination),
			},
			Scheme: &models.RequestFieldMatchers{
				ExactMatch: util.StringToPointer(request.Scheme),
			},
			Query: &models.RequestFieldMatchers{
				ExactMatch: util.StringToPointer(request.QueryString()),
			},
			Body:    body,
			Headers: headers,
		},
		Response: *response,
	}

	hf.Simulation.AddRequestMatcherResponsePair(&pair)

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

func (this Hoverfly) GetSimulationPairsCount() int {
	return len(this.Simulation.MatchingPairs)
}

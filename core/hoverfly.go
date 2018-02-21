package hoverfly

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"sync"

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
	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
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
	Authentication backends.Authentication

	HTTP    *http.Client
	Cfg     *Configuration
	Counter *metrics.CounterByMode

	Proxy   *goproxy.ProxyHttpServer
	SL      *StoppableListener
	mu      sync.Mutex
	version string

	modeMap map[string]modes.Mode

	state map[string]string

	Simulation    *models.Simulation
	StoreLogsHook *StoreLogsHook
	Journal       *journal.Journal
	templator     *templating.Templator

	responsesDiff map[v2.SimpleRequestDefinitionView][]string
}

func NewHoverfly() *Hoverfly {

	authBackend := backends.NewCacheBasedAuthBackend(cache.NewInMemoryCache(), cache.NewInMemoryCache())

	hoverfly := &Hoverfly{
		Simulation:     models.NewSimulation(),
		Authentication: authBackend,
		Counter:        metrics.NewModeCounter([]string{modes.Simulate, modes.Synthesize, modes.Modify, modes.Capture, modes.Spy, modes.Diff}),
		StoreLogsHook:  NewStoreLogsHook(),
		Journal:        journal.NewJournal(),
		Cfg:            InitSettings(),
		state:          make(map[string]string),
		templator:      templating.NewTemplator(),
		responsesDiff:  make(map[v2.SimpleRequestDefinitionView][]string),
	}

	hoverfly.version = "v0.15.1"

	log.AddHook(hoverfly.StoreLogsHook)

	modeMap := make(map[string]modes.Mode)

	modeMap[modes.Capture] = &modes.CaptureMode{Hoverfly: hoverfly}
	modeMap[modes.Simulate] = &modes.SimulateMode{Hoverfly: hoverfly, MatchingStrategy: "strongest"}
	modeMap[modes.Modify] = &modes.ModifyMode{Hoverfly: hoverfly}
	modeMap[modes.Synthesize] = &modes.SynthesizeMode{Hoverfly: hoverfly}
	modeMap[modes.Spy] = &modes.SpyMode{Hoverfly: hoverfly}
	modeMap[modes.Diff] = &modes.DiffMode{Hoverfly: hoverfly}

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

	hoverfly.CacheMatcher = matching.CacheMatcher{
		RequestCache: requestCache,
		Webserver:    cfg.Webserver,
	}

	hoverfly.Cfg = cfg
	hoverfly.HTTP = GetDefaultHoverflyHTTPClient(cfg.TLSVerification, cfg.UpstreamProxy)

	return hoverfly
}

// GetNewHoverfly returns a configured ProxyHttpServer and DBClient
func GetNewHoverfly(cfg *Configuration, requestCache cache.Cache, authentication backends.Authentication) *Hoverfly {
	hoverfly := NewHoverfly()

	if cfg.DisableCache {
		requestCache = nil
	}

	hoverfly.CacheMatcher = matching.CacheMatcher{
		RequestCache: requestCache,
		Webserver:    cfg.Webserver,
	}

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
		Proxy: proxyURL,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: !tlsVerification,
			Renegotiation:      tls.RenegotiateFreelyAsClient,
		},
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
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", hf.Cfg.ListenOnHost, hf.Cfg.ProxyPort))
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

	modeName := hf.Cfg.GetMode()
	mode := hf.modeMap[modeName]
	response, err := mode.Process(req, requestDetails)

	if modeName == modes.Diff {
		hf.handleResponseDiff(req, mode)
	}

	// Don't delete the error
	// and definitely don't delay people in capture mode
	if err != nil || modeName == modes.Capture {
		return response
	}

	respDelay := hf.Simulation.ResponseDelays.GetDelay(requestDetails)
	if respDelay != nil {
		respDelay.Execute()
	}

	return response
}

func (hf *Hoverfly) handleResponseDiff(request *http.Request, mode modes.Mode) {
	switch mode.(type) {
	case *modes.DiffMode:
		diffMode := mode.(*modes.DiffMode)
		errorMessage := diffMode.GetMessage()

		if errorMessage.GetErrorMessage() != "" {
			requestView := v2.SimpleRequestDefinitionView{
				Method: request.Method,
				Host:   request.URL.Host,
				Path:   request.URL.Path,
				Query:  request.URL.RawQuery,
			}

			diffs := hf.responsesDiff[requestView]
			hf.responsesDiff[requestView] = append(diffs, errorMessage.GetErrorMessage())
		}
	}
}

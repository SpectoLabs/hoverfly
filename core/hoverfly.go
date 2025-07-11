package hoverfly

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/SpectoLabs/goproxy"
	"github.com/SpectoLabs/hoverfly/core/action"
	"github.com/SpectoLabs/hoverfly/core/authentication/backends"
	"github.com/SpectoLabs/hoverfly/core/cache"
	"github.com/SpectoLabs/hoverfly/core/delay"
	v2 "github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/journal"
	"github.com/SpectoLabs/hoverfly/core/matching"
	"github.com/SpectoLabs/hoverfly/core/metrics"
	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/SpectoLabs/hoverfly/core/modes"
	"github.com/SpectoLabs/hoverfly/core/state"
	"github.com/SpectoLabs/hoverfly/core/templating"
	log "github.com/sirupsen/logrus"
)

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

	state *state.State

	Simulation             *models.Simulation
	StoreLogsHook          *StoreLogsHook
	Journal                *journal.Journal
	templator              *templating.Templator
	PostServeActionDetails *action.PostServeActionDetails
	responsesDiff          map[v2.SimpleRequestDefinitionView][]v2.DiffReport
}

func NewHoverfly() *Hoverfly {

	authBackend := backends.NewCacheBasedAuthBackend(cache.NewInMemoryCache(), cache.NewInMemoryCache())

	newJournal := journal.NewJournal()
	hoverfly := &Hoverfly{
		Simulation:             models.NewSimulation(),
		Authentication:         authBackend,
		Counter:                metrics.NewModeCounter([]string{modes.Simulate, modes.Synthesize, modes.Modify, modes.Capture, modes.Spy, modes.Diff}),
		StoreLogsHook:          NewStoreLogsHook(),
		Journal:                newJournal,
		Cfg:                    InitSettings(),
		state:                  state.NewState(),
		templator:              templating.NewEnrichedTemplator(newJournal),
		responsesDiff:          make(map[v2.SimpleRequestDefinitionView][]v2.DiffReport),
		PostServeActionDetails: action.NewPostServeActionDetails(),
	}

	hoverfly.version = "v1.11.2"

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

	var requestCache cache.FastCache
	if !cfg.DisableCache {
		if cfg.CacheSize > 0 {
			requestCache, _ = cache.NewLRUCache(cfg.CacheSize)
		} else {
			// Backward compatibility, always set default cache if cache size is not configured
			requestCache = cache.NewDefaultLRUCache()
		}
	}

	hoverfly.CacheMatcher = matching.CacheMatcher{
		Webserver:    cfg.Webserver,
		RequestCache: requestCache,
	}

	hoverfly.Cfg = cfg
	hoverfly.HTTP = GetDefaultHoverflyHTTPClient(cfg.TLSVerification, cfg.UpstreamProxy)

	return hoverfly
}

// GetNewHoverfly returns a configured ProxyHttpServer and DBClient
func GetNewHoverfly(cfg *Configuration, requestCache cache.FastCache, authentication backends.Authentication) *Hoverfly {
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
func (hf *Hoverfly) processRequest(req *http.Request) (*http.Response, chan string) {
	if hf.Cfg.CORS.Enabled {
		response := hf.Cfg.CORS.InterceptPreflightRequest(req)
		if response != nil {
			return response, nil
		}
	}
	requestDetails, err := models.NewRequestDetailsFromHttpRequest(req)
	if err != nil {
		return modes.ErrorResponse(req, err, "Could not interpret HTTP request").Response, nil
	}

	modeName := hf.Cfg.GetMode()
	mode := hf.modeMap[modeName]
	result, err := mode.Process(req, requestDetails)

	if err == nil && hf.Cfg.CORS.Enabled {
		hf.Cfg.CORS.AddCORSHeaders(req, result.Response)
	}

	// and definitely don't delay people in capture mode
	// Don't delete the error
	if err != nil || modeName == modes.Capture {
		return result.Response, nil
	}

	if result.IsResponseDelayable() {
		log.Debug("Applying response delay")
		hf.applyResponseDelay(result)
	} else {
		log.Debug("Applying global delay")
		hf.applyGlobalDelay(requestDetails)
	}

	if result.PostServeActionInputDetails != nil && result.PostServeActionInputDetails.PostServeAction != "" {
		if postServeAction, ok := hf.PostServeActionDetails.Actions[result.PostServeActionInputDetails.PostServeAction]; ok {
			journalIDChannel := make(chan string, 1)
			go postServeAction.Execute(result.PostServeActionInputDetails.Pair, journalIDChannel, hf.Journal)
			return result.Response, journalIDChannel
		} else if hf.PostServeActionDetails.FallbackAction != nil {
			journalIDChannel := make(chan string, 1)
			go hf.PostServeActionDetails.FallbackAction.Execute(result.PostServeActionInputDetails.Pair, journalIDChannel, hf.Journal)
			return result.Response, journalIDChannel
		}
	}

	return result.Response, nil
}

func (hf *Hoverfly) applyResponseDelay(result modes.ProcessResult) {
	if result.FixedDelay > 0 {
		time.Sleep(time.Duration(result.FixedDelay) * time.Millisecond)
	}

	if result.LogNormalDelay != nil {
		logNormalDelay := delay.NewLogNormalGenerator(
			result.LogNormalDelay.Min, result.LogNormalDelay.Max,
			result.LogNormalDelay.Mean, result.LogNormalDelay.Median,
		).GenerateDelay()

		time.Sleep(time.Duration(logNormalDelay) * time.Millisecond)
	}
}

func (hf *Hoverfly) applyGlobalDelay(requestDetails models.RequestDetails) {
	respDelay := hf.Simulation.ResponseDelays.GetDelay(requestDetails)
	if respDelay != nil {
		respDelay.Execute()
	}

	respDelayLogNormal := hf.Simulation.ResponseDelaysLogNormal.GetDelay(requestDetails)
	if respDelayLogNormal != nil {
		respDelayLogNormal.Execute()
	}
}

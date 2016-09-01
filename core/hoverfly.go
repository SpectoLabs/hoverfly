package hoverfly

import (
	"bytes"
	"crypto/tls"
	"fmt"
	log "github.com/Sirupsen/logrus"
	authBackend "github.com/SpectoLabs/hoverfly/core/authentication/backends"
	"github.com/SpectoLabs/hoverfly/core/cache"
	"github.com/SpectoLabs/hoverfly/core/matching"
	"github.com/SpectoLabs/hoverfly/core/metrics"
	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/rusenask/goproxy"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
	"sync"
	"time"
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
		HTTP: &http.Client{Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: cfg.TLSVerification},
		}},
		Cfg:            cfg,
		Counter:        metrics.NewModeCounter([]string{SimulateMode, SynthesizeMode, ModifyMode, CaptureMode}),
		Hooks:          make(ActionTypeHooks),
		ResponseDelays: &models.ResponseDelayList{},
		RequestMatcher: requestMatcher,
	}
	return h
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

// UpdateDestination - updates proxy with new destination regexp
func (hf *Hoverfly) UpdateDestination(destination string) (err error) {
	_, err = regexp.Compile(destination)
	if err != nil {
		return fmt.Errorf("destination is not a valid regular expression string")
	}

	hf.mu.Lock()
	hf.StopProxy()
	hf.Cfg.Destination = destination
	err = hf.StartProxy()
	hf.mu.Unlock()
	return
}

func (hf *Hoverfly) SetMiddleware(middleware string) error {
	if middleware == "" {
		hf.Cfg.Middleware = middleware
		return nil
	}
	originalPair := models.RequestResponsePair{
		Request: models.RequestDetails{
			Path:        "/",
			Method:      "GET",
			Destination: "www.test.com",
			Scheme:      "",
			Query:       "",
			Body:        "",
			Headers:     map[string][]string{"test_header": []string{"true"}},
		},
		Response: models.ResponseDetails{
			Status:  200,
			Body:    "ok",
			Headers: map[string][]string{"test_header": []string{"true"}},
		},
	}
	c := NewConstructor(nil, originalPair)
	err := c.ApplyMiddleware(middleware)
	if err != nil {
		return err
	}

	hf.Cfg.Middleware = middleware
	return nil
}

func (hf *Hoverfly) UpdateResponseDelays(responseDelays models.ResponseDelayList) {
	hf.ResponseDelays = &responseDelays
	log.Info("Response delay config updated on hoverfly")
}

func hoverflyError(req *http.Request, err error, msg string, statusCode int) *http.Response {
	return goproxy.NewResponse(req,
		goproxy.ContentTypeText, statusCode,
		fmt.Sprintf("Hoverfly Error! %s. Got error: %s \n", msg, err.Error()))
}

// processRequest - processes incoming requests and based on proxy state (record/playback)
// returns HTTP response.
func (hf *Hoverfly) processRequest(req *http.Request) (*http.Response) {
	var response *http.Response

	mode := hf.Cfg.GetMode()

	requestDetails, err := models.NewRequestDetailsFromHttpRequest(req)
	if err != nil {
		return hoverflyError(req, err, "Could not interpret HTTP request", http.StatusServiceUnavailable)
	}

	if mode == CaptureMode {
		var err error
		response, err = hf.captureRequest(req)

		if err != nil {
			return hoverflyError(req, err, "Could not capture request", http.StatusServiceUnavailable)
		}
		log.WithFields(log.Fields{
			"mode":        mode,
			"middleware":  hf.Cfg.Middleware,
			"path":        req.URL.Path,
			"rawQuery":    req.URL.RawQuery,
			"method":      req.Method,
			"destination": req.Host,
		}).Info("request and response captured")

		return response

	} else if mode == SynthesizeMode {
		var err error
		response, err = SynthesizeResponse(req, requestDetails, hf.Cfg.Middleware)

		if err != nil {
			return hoverflyError(req, err, "Could not create synthetic response!", http.StatusServiceUnavailable)
		}

		log.WithFields(log.Fields{
			"mode":        mode,
			"middleware":  hf.Cfg.Middleware,
			"path":        req.URL.Path,
			"rawQuery":    req.URL.RawQuery,
			"method":      req.Method,
			"destination": req.Host,
		}).Info("synthetic response created successfuly")

	} else if mode == ModifyMode {
		var err error
		response, err = hf.modifyRequestResponse(req, requestDetails, hf.Cfg.Middleware)

		if err != nil {
			log.WithFields(log.Fields{
				"error":      err.Error(),
				"middleware": hf.Cfg.Middleware,
			}).Error("Got error when performing request modification")
			return hoverflyError(req, err, fmt.Sprintf("Middleware (%s) failed or something else happened!", hf.Cfg.Middleware), http.StatusServiceUnavailable)
		}

	} else {
		var err *matching.MatchingError
		response, err = hf.getResponse(req, requestDetails)
		if err != nil {
			return hoverflyError(req, err, err.Error(), err.StatusCode)
		}
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

// captureRequest saves request for later playback
func (hf *Hoverfly) captureRequest(req *http.Request) (*http.Response, error) {

	// this is mainly for testing, since when you create
	if req.Body == nil {
		req.Body = ioutil.NopCloser(bytes.NewBuffer([]byte("")))
	}

	reqBody, err := ioutil.ReadAll(req.Body)

	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
			"mode":  "capture",
		}).Error("Got error when reading request body")
	}

	// outputting request body if verbose logging is set
	log.WithFields(log.Fields{
		"body": string(reqBody),
		"mode": "capture",
	}).Debug("got request body")

	// forwarding request
	req.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))

	req, resp, err := hf.doRequest(req)

	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
			"mode":  "capture",
		}).Error("Got error when reading body after being modified by middleware")
	}

	reqBody, err = ioutil.ReadAll(req.Body)
	req.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))

	if err == nil {
		respBody, err := extractBody(resp)

		if err != nil {

			log.WithFields(log.Fields{
				"error": err.Error(),
				"mode":  "capture",
			}).Error("Failed to copy response body.")

			return resp, err
		}

		// saving response body with request/response meta to cache
		hf.save(req, reqBody, resp, respBody)
	}

	// return new response or error here
	return resp, err
}

// doRequest performs original request and returns response that should be returned to client and error (if there is one)
func (hf *Hoverfly) doRequest(request *http.Request) (*http.Request, *http.Response, error) {

	// We can't have this set. And it only contains "/pkg/net/http/" anyway
	request.RequestURI = ""

	if hf.Cfg.Middleware != "" {
		// middleware is provided, modifying request
		var requestResponsePair models.RequestResponsePair

		rd, err := models.NewRequestDetailsFromHttpRequest(request)
		if err != nil {
			return nil, nil, err
		}
		requestResponsePair.Request = rd

		c := NewConstructor(request, requestResponsePair)
		err = c.ApplyMiddleware(hf.Cfg.Middleware)

		if err != nil {
			log.WithFields(log.Fields{
				"mode":   hf.Cfg.Mode,
				"error":  err.Error(),
				"host":   request.Host,
				"method": request.Method,
				"path":   request.URL.Path,
			}).Error("could not forward request, middleware failed to modify request.")
			return nil, nil, err
		}

		request, err = c.ReconstructRequest()

		if err != nil {
			return nil, nil, err
		}
	}

	requestBody, _ := ioutil.ReadAll(request.Body)

	request.Body = ioutil.NopCloser(bytes.NewReader(requestBody))

	resp, err := hf.HTTP.Do(request)

	request.Body = ioutil.NopCloser(bytes.NewReader(requestBody))

	if err != nil {
		log.WithFields(log.Fields{
			"mode":   hf.Cfg.Mode,
			"error":  err.Error(),
			"host":   request.Host,
			"method": request.Method,
			"path":   request.URL.Path,
		}).Error("could not forward request, failed to do an HTTP request.")
		return nil, nil, err
	}

	log.WithFields(log.Fields{
		"mode":   hf.Cfg.Mode,
		"host":   request.Host,
		"method": request.Method,
		"path":   request.URL.Path,
	}).Debug("response from external service got successfuly!")

	resp.Header.Set("hoverfly", "Was-Here")

	return request, resp, nil

}

// getResponse returns stored response from cache
func (hf *Hoverfly) getResponse(req *http.Request, requestDetails models.RequestDetails) (*http.Response, *matching.MatchingError) {

	responseDetails, matchErr := hf.RequestMatcher.GetResponse(&requestDetails)
	if matchErr != nil {
		return nil, matchErr
	}

	pair := &models.RequestResponsePair{
		Request:  requestDetails,
		Response: *responseDetails,
	}

	c := NewConstructor(req, *pair)
	if hf.Cfg.Middleware != "" {
		_ = c.ApplyMiddleware(hf.Cfg.Middleware)
	}

	return c.ReconstructResponse(), nil
}

// modifyRequestResponse modifies outgoing request and then modifies incoming response, neither request nor response
// is saved to cache.
func (hf *Hoverfly) modifyRequestResponse(req *http.Request, requestDetails models.RequestDetails, middleware string) (*http.Response, error) {

	// modifying request
	req, resp, err := hf.doRequest(req)

	if err != nil {
		return nil, err
	}

	// preparing payload
	bodyBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.WithFields(log.Fields{
			"error":      err.Error(),
			"middleware": middleware,
		}).Error("Failed to read response body after sending modified request")
		return nil, err
	}

	r := models.ResponseDetails{
		Status:  resp.StatusCode,
		Body:    string(bodyBytes),
		Headers: resp.Header,
	}

	requestResponsePair := models.RequestResponsePair{Response: r, Request: requestDetails}

	c := NewConstructor(req, requestResponsePair)
	// applying middleware to modify response
	err = c.ApplyMiddleware(middleware)

	if err != nil {
		return nil, err
	}

	newResponse := c.ReconstructResponse()

	log.WithFields(log.Fields{
		"status":      newResponse.StatusCode,
		"middleware":  middleware,
		"mode":        ModifyMode,
		"path":        c.requestResponsePair.Request.Path,
		"rawQuery":    c.requestResponsePair.Request.Query,
		"method":      c.requestResponsePair.Request.Method,
		"destination": c.requestResponsePair.Request.Destination,
		// original here
		"originalPath":        req.URL.Path,
		"originalRawQuery":    req.URL.RawQuery,
		"originalMethod":      req.Method,
		"originalDestination": req.Host,
	}).Info("request and response modified, returning")

	return newResponse, nil
}

// save gets request fingerprint, extracts request body, status code and headers, then saves it to cache
func (hf *Hoverfly) save(req *http.Request, reqBody []byte, resp *http.Response, respBody []byte) {

	if resp == nil {
		resp = emptyResp
	} else {
		responseObj := models.ResponseDetails{
			Status:  resp.StatusCode,
			Body:    string(respBody),
			Headers: resp.Header,
		}

		requestObj := models.RequestDetails{
			Path:        req.URL.Path,
			Method:      req.Method,
			Destination: req.Host,
			Scheme:      req.URL.Scheme,
			Query:       req.URL.RawQuery,
			Body:        string(reqBody),
			Headers:     req.Header,
		}

		pair := models.RequestResponsePair{
			Response: responseObj,
			Request:  requestObj,
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
}

package hoverfly

import (
	"crypto/tls"
	"fmt"
	log "github.com/Sirupsen/logrus"
	authBackend "github.com/SpectoLabs/hoverfly/core/authentication/backends"
	"github.com/SpectoLabs/hoverfly/core/cache"
	"github.com/SpectoLabs/hoverfly/core/metrics"
	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/rusenask/goproxy"
	"net"
	"net/http"
	"regexp"
	"sync"
	"io/ioutil"
	"bytes"
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
	MetadataCache  cache.Cache
	Authentication authBackend.Authentication
	HTTP           *http.Client
	Cfg            *Configuration
	Counter        *metrics.CounterByMode
	Hooks          ActionTypeHooks

	Proxy *goproxy.ProxyHttpServer
	SL    *StoppableListener
	mu    sync.Mutex
}

// GetNewHoverfly returns a configured ProxyHttpServer and DBClient
func GetNewHoverfly(cfg *Configuration, requestCache, metadataCache cache.Cache, authentication authBackend.Authentication) *Hoverfly {
	h := &Hoverfly{
		RequestCache:   requestCache,
		MetadataCache:  metadataCache,
		Authentication: authentication,
		HTTP: &http.Client{Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: cfg.TLSVerification},
		}},
		Cfg:     cfg,
		Counter: metrics.NewModeCounter([]string{SimulateMode, SynthesizeMode, ModifyMode, CaptureMode}),
		Hooks:   make(ActionTypeHooks),
	}
	return h
}

// StartProxy - starts proxy with current configuration, this method is non blocking.
func (d *Hoverfly) StartProxy() error {

	if d.Cfg.ProxyPort == "" {
		return fmt.Errorf("Proxy port is not set!")
	}

	if d.Cfg.Webserver {
		d.Proxy = NewWebserverProxy(d)
	} else {
		d.Proxy = NewProxy(d)
	}
	
	log.WithFields(log.Fields{
		"destination": d.Cfg.Destination,
		"port":        d.Cfg.ProxyPort,
		"mode":        d.Cfg.GetMode(),
	}).Info("current proxy configuration")

	// creating TCP listener
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", d.Cfg.ProxyPort))
	if err != nil {
		return err
	}

	sl, err := NewStoppableListener(listener)
	if err != nil {
		return err
	}
	d.SL = sl
	server := http.Server{}

	d.Cfg.ProxyControlWG.Add(1)

	go func() {
		defer func() {
			log.Info("sending done signal")
			d.Cfg.ProxyControlWG.Done()
		}()
		log.Info("serving proxy")
		server.Handler = d.Proxy
		log.Warn(server.Serve(sl))
	}()

	return nil
}

// StopProxy - stops proxy
func (d *Hoverfly) StopProxy() {
	d.SL.Stop()
	d.Cfg.ProxyControlWG.Wait()
}

// UpdateDestination - updates proxy with new destination regexp
func (d *Hoverfly) UpdateDestination(destination string) (err error) {
	_, err = regexp.Compile(destination)
	if err != nil {
		return fmt.Errorf("destination is not a valid regular expression string")
	}

	d.mu.Lock()
	d.StopProxy()
	d.Cfg.Destination = destination
	err = d.StartProxy()
	d.mu.Unlock()
	return
}

func (d *Hoverfly) UpdateResponseDelays(responseDelays []models.ResponseDelay) {
	d.Cfg.ResponseDelays = responseDelays
	log.Info("Response delay config updated on hoverfly")
}

func hoverflyError(req *http.Request, err error, msg string, statusCode int) *http.Response {
	return goproxy.NewResponse(req,
		goproxy.ContentTypeText, statusCode,
		fmt.Sprintf("Hoverfly Error! %s. Got error: %s \n", msg, err.Error()))
}

// processRequest - processes incoming requests and based on proxy state (record/playback)
// returns HTTP response.
func (d *Hoverfly) processRequest(req *http.Request) (*http.Request, *http.Response) {

	mode := d.Cfg.GetMode()

	if mode == CaptureMode {
		newResponse, err := d.captureRequest(req)

		if err != nil {
			return req, hoverflyError(req, err, "Could not capture request", http.StatusServiceUnavailable)
		}
		log.WithFields(log.Fields{
			"mode":        mode,
			"middleware":  d.Cfg.Middleware,
			"path":        req.URL.Path,
			"rawQuery":    req.URL.RawQuery,
			"method":      req.Method,
			"destination": req.Host,
		}).Info("request and response captured")

		return req, newResponse

	} else if mode == SynthesizeMode {
		response, err := SynthesizeResponse(req, d.Cfg.Middleware)

		if err != nil {
			return req, hoverflyError(req, err, "Could not create synthetic response!", http.StatusServiceUnavailable)
		}

		log.WithFields(log.Fields{
			"mode":        mode,
			"middleware":  d.Cfg.Middleware,
			"path":        req.URL.Path,
			"rawQuery":    req.URL.RawQuery,
			"method":      req.Method,
			"destination": req.Host,
		}).Info("synthetic response created successfuly")

		respDelay := d.Cfg.GetDelay(req.Host)
		if (respDelay != nil) {
			respDelay.Execute()
		}

		return req, response

	} else if mode == ModifyMode {

		response, err := d.modifyRequestResponse(req, d.Cfg.Middleware)

		if err != nil {
			log.WithFields(log.Fields{
				"error":      err.Error(),
				"middleware": d.Cfg.Middleware,
			}).Error("Got error when performing request modification")
			return req, hoverflyError(
				req,
				err,
				fmt.Sprintf("Middleware (%s) failed or something else happened!", d.Cfg.Middleware),
				http.StatusServiceUnavailable)
		}

		respDelay := d.Cfg.GetDelay(req.Host)
		if (respDelay != nil) {
			respDelay.Execute()
		}

		// returning modified response
		return req, response
	}

	newResponse := d.getResponse(req)

	respDelay := d.Cfg.GetDelay(req.Host)
	if (respDelay != nil) {
		respDelay.Execute()
	}

	return req, newResponse

}

// AddHook - adds a hook to DBClient
func (d *Hoverfly) AddHook(hook Hook) {
	d.Hooks.Add(hook)
}

// captureRequest saves request for later playback
func (d *Hoverfly) captureRequest(req *http.Request) (*http.Response, error) {

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

	req, resp, err := d.doRequest(req)

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
		d.save(req, reqBody, resp, respBody)
	}

	// return new response or error here
	return resp, err
}

// doRequest performs original request and returns response that should be returned to client and error (if there is one)
func (d *Hoverfly) doRequest(request *http.Request) (*http.Request, *http.Response, error) {

	// We can't have this set. And it only contains "/pkg/net/http/" anyway
	request.RequestURI = ""

	if d.Cfg.Middleware != "" {
		// middleware is provided, modifying request
		var payload models.Payload

		rd, err := getRequestDetails(request)
		if err != nil {
			return nil, nil, err
		}
		payload.Request = rd

		c := NewConstructor(request, payload)
		err = c.ApplyMiddleware(d.Cfg.Middleware)

		if err != nil {
			log.WithFields(log.Fields{
				"mode":   d.Cfg.Mode,
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

	resp, err := d.HTTP.Do(request)

	request.Body = ioutil.NopCloser(bytes.NewReader(requestBody))

	if err != nil {
		log.WithFields(log.Fields{
			"mode":   d.Cfg.Mode,
			"error":  err.Error(),
			"host":   request.Host,
			"method": request.Method,
			"path":   request.URL.Path,
		}).Error("could not forward request, failed to do an HTTP request.")
		return nil, nil, err
	}

	log.WithFields(log.Fields{
		"mode":   d.Cfg.Mode,
		"host":   request.Host,
		"method": request.Method,
		"path":   request.URL.Path,
	}).Debug("response from external service got successfuly!")

	resp.Header.Set("hoverfly", "Was-Here")

	return request, resp, nil

}

// getResponse returns stored response from cache
func (d *Hoverfly) getResponse(req *http.Request) *http.Response {

	if req.Body == nil {
		req.Body = ioutil.NopCloser(bytes.NewBuffer([]byte("")))
	}

	reqBody, err := ioutil.ReadAll(req.Body)

	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Got error when reading request body")
	}

	key := d.getRequestFingerprint(req, reqBody)

	payloadBts, err := d.RequestCache.Get([]byte(key))

	if err == nil {
		// getting cache response
		payload, err := models.NewPayloadFromBytes(payloadBts)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
				"value": string(payloadBts),
				"key":   key,
			}).Error("Failed to decode payload")
			return hoverflyError(req, err, "Failed to simulate", http.StatusInternalServerError)
		}

		c := NewConstructor(req, *payload)

		if d.Cfg.Middleware != "" {
			_ = c.ApplyMiddleware(d.Cfg.Middleware)
		}

		response := c.ReconstructResponse()

		log.WithFields(log.Fields{
			"key":         key,
			"mode":        SimulateMode,
			"middleware":  d.Cfg.Middleware,
			"path":        req.URL.Path,
			"rawQuery":    req.URL.RawQuery,
			"method":      req.Method,
			"destination": req.Host,
			"status":      payload.Response.Status,
			"bodyLength":  response.ContentLength,
		}).Info("Response found, returning")

		return response

	}

	log.WithFields(log.Fields{
		"key":         key,
		"error":       err.Error(),
		"query":       req.URL.RawQuery,
		"path":        req.URL.RawPath,
		"destination": req.Host,
		"method":      req.Method,
	}).Warn("Failed to retrieve response from cache")
	// return error? if we return nil - proxy forwards request to original destination
	return hoverflyError(req, err, "Could not find recorded request, please record it first!", http.StatusPreconditionFailed)
}

// modifyRequestResponse modifies outgoing request and then modifies incoming response, neither request nor response
// is saved to cache.
func (d *Hoverfly) modifyRequestResponse(req *http.Request, middleware string) (*http.Response, error) {

	// getting request details
	rd, err := getRequestDetails(req)
	if err != nil {
		log.WithFields(log.Fields{
			"error":      err.Error(),
			"middleware": middleware,
		}).Error("Failed to get request details")
		return nil, err
	}

	// modifying request
	req, resp, err := d.doRequest(req)

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

	payload := models.Payload{Response: r, Request: rd}

	c := NewConstructor(req, payload)
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
		"path":        c.payload.Request.Path,
		"rawQuery":    c.payload.Request.Query,
		"method":      c.payload.Request.Method,
		"destination": c.payload.Request.Destination,
		// original here
		"originalPath":        req.URL.Path,
		"originalRawQuery":    req.URL.RawQuery,
		"originalMethod":      req.Method,
		"originalDestination": req.Host,
	}).Info("request and response modified, returning")

	return newResponse, nil
}

// getRequestFingerprint returns request hash
func (d *Hoverfly) getRequestFingerprint(req *http.Request, requestBody []byte) string {
	r := models.RequestDetails{
		Path:        req.URL.Path,
		Method:      req.Method,
		Destination: req.Host,
		Query:       req.URL.RawQuery,
		Body:        string(requestBody),
		Headers:     req.Header,
	}

	return r.Hash()
}

// save gets request fingerprint, extracts request body, status code and headers, then saves it to cache
func (d *Hoverfly) save(req *http.Request, reqBody []byte, resp *http.Response, respBody []byte) {
	// record request here
	key := d.getRequestFingerprint(req, reqBody)

	if resp == nil {
		resp = emptyResp
	} else {
		responseObj := models.ResponseDetails{
			Status:  resp.StatusCode,
			Body:    string(respBody),
			Headers: resp.Header,
		}

		log.WithFields(log.Fields{
			"path":          req.URL.Path,
			"rawQuery":      req.URL.RawQuery,
			"requestMethod": req.Method,
			"bodyLen":       len(reqBody),
			"destination":   req.Host,
			"hashKey":       key,
		}).Debug("Capturing")

		requestObj := models.RequestDetails{
			Path:        req.URL.Path,
			Method:      req.Method,
			Destination: req.Host,
			Scheme:      req.URL.Scheme,
			Query:       req.URL.RawQuery,
			Body:        string(reqBody),
			Headers:     req.Header,
		}

		payload := models.Payload{
			Response: responseObj,
			Request:  requestObj,
		}

		bts, err := payload.Encode()

		// hook
		var en Entry
		en.ActionType = ActionTypeRequestCaptured
		en.Message = "captured"
		en.Time = time.Now()
		en.Data = bts

		if err := d.Hooks.Fire(ActionTypeRequestCaptured, &en); err != nil {
			log.WithFields(log.Fields{
				"error":      err.Error(),
				"message":    en.Message,
				"actionType": ActionTypeRequestCaptured,
			}).Error("failed to fire hook")
		}

		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Error("Failed to serialize payload")
		} else {
			d.RequestCache.Set([]byte(key), bts)
		}
	}
}

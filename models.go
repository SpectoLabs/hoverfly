package hoverfly

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
	"sync"
	"time"
	log "github.com/Sirupsen/logrus"
	authBackend "github.com/SpectoLabs/hoverfly/authentication/backends"
	"github.com/SpectoLabs/hoverfly/cache"
	"github.com/SpectoLabs/hoverfly/metrics"
	"github.com/rusenask/goproxy"
	"github.com/SpectoLabs/hoverfly/models"
)


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
	mu sync.Mutex
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
	d.UpdateProxy()
	err = d.StartProxy()
	d.mu.Unlock()
	return
}

// StartProxy - starts proxy with current configuration, this method is non blocking.
func (d *Hoverfly) StartProxy() error {

	if d.Cfg.ProxyPort == "" {
		return fmt.Errorf("Proxy port is not set!")
	}

	if d.Proxy == nil {
		d.UpdateProxy()
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

// AddHook - adds a hook to DBClient
func (d *Hoverfly) AddHook(hook Hook) {
	d.Hooks.Add(hook)
}


var emptyResp = &http.Response{}





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

func copyBody(body io.ReadCloser) (resp1, resp2 io.ReadCloser, err error) {
	var buf bytes.Buffer
	if _, err = buf.ReadFrom(body); err != nil {
		return nil, nil, err
	}
	if err = body.Close(); err != nil {
		return nil, nil, err
	}
	return ioutil.NopCloser(&buf), ioutil.NopCloser(bytes.NewReader(buf.Bytes())), nil
}

func extractBody(resp *http.Response) (extract []byte, err error) {
	save := resp.Body
	savecl := resp.ContentLength

	save, resp.Body, err = copyBody(resp.Body)

	if err != nil {
		return
	}
	defer resp.Body.Close()
	extract, err = ioutil.ReadAll(resp.Body)

	resp.Body = save
	resp.ContentLength = savecl
	if err != nil {
		return nil, err
	}
	return extract, nil
}

func extractRequestBody(req *http.Request) (extract []byte, err error) {
	save := req.Body
	savecl := req.ContentLength

	save, req.Body, err = copyBody(req.Body)

	if err != nil {
		return
	}
	defer req.Body.Close()
	extract, err = ioutil.ReadAll(req.Body)

	req.Body = save
	req.ContentLength = savecl
	if err != nil {
		return nil, err
	}
	return extract, nil
}

// getRequestDetails - extracts request details
func getRequestDetails(req *http.Request) (requestObj models.RequestDetails, err error) {
	if req.Body == nil {
		req.Body = ioutil.NopCloser(bytes.NewBuffer([]byte("")))
	}

	reqBody, err := extractRequestBody(req)

	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
			"mode":  "capture",
		}).Error("Got error while reading request body")
		return
	}

	requestObj = models.RequestDetails{
		Path:        req.URL.Path,
		Method:      req.Method,
		Destination: req.Host,
		Scheme:      req.URL.Scheme,
		Query:       req.URL.RawQuery,
		Body:        string(reqBody),
		Headers:     req.Header,
	}
	return
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

// ActionType - action type can be things such as "RequestCaptured", "GotResponse" - anything
type ActionType string

// ActionTypeRequestCaptured - default action type name for identifying
const ActionTypeRequestCaptured = "requestCaptured"

// ActionTypeWipeDB - default action type for wiping database
const ActionTypeWipeDB = "wipeDatabase"

// ActionTypeConfigurationChanged - default action name for identifying configuration changes
const ActionTypeConfigurationChanged = "configurationChanged"

// Entry - holds information about action, based on action type - other clients will be able to decode
// the data field.
type Entry struct {
	// Contains encoded data
	Data []byte

	// Time at which the action entry was fired
	Time time.Time

	ActionType ActionType

	// Message, can carry additional information
	Message string
}

// Hook - an interface to add dynamic hooks to extend functionality
type Hook interface {
	ActionTypes() []ActionType
	Fire(*Entry) error
}

// ActionTypeHooks type for storing the hooks
type ActionTypeHooks map[ActionType][]Hook

// Add a hook
func (hooks ActionTypeHooks) Add(hook Hook) {
	for _, ac := range hook.ActionTypes() {
		hooks[ac] = append(hooks[ac], hook)
	}
}

// Fire all the hooks for the passed ActionType
func (hooks ActionTypeHooks) Fire(ac ActionType, entry *Entry) error {
	for _, hook := range hooks[ac] {
		if err := hook.Fire(entry); err != nil {
			return err
		}
	}

	return nil
}

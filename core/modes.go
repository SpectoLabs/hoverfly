package hoverfly

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/SpectoLabs/hoverfly/core/modes"
	"github.com/SpectoLabs/hoverfly/core/util"
)

type Modify struct {
	hoverfly *Hoverfly
}

func (this Modify) Process(request *http.Request, details models.RequestDetails) (*http.Response, error) {
	response, err := this.hoverfly.modifyRequestResponse(request, details)

	if err != nil {
		log.WithFields(log.Fields{
			"error":      err.Error(),
			"middleware": this.hoverfly.Cfg.Middleware,
		}).Error("Got error when performing request modification")
		return hoverflyError(request, err, fmt.Sprintf("Middleware (%s) failed or something else happened!", this.hoverfly.Cfg.Middleware), http.StatusServiceUnavailable), err
	}

	return response, nil
}

type Synthesize struct {
	hoverfly *Hoverfly
}

func (this Synthesize) Process(request *http.Request, details models.RequestDetails) (*http.Response, error) {
	pair := models.RequestResponsePair{Request: details}

	log.WithFields(log.Fields{
		"middleware":  this.hoverfly.Cfg.Middleware.toString(),
		"body":        details.Body,
		"destination": details.Destination,
	}).Debug("Synthesizing new response")

	if !this.hoverfly.Cfg.Middleware.IsSet() {
		err := fmt.Errorf("Middleware not set")
		return hoverflyError(request, err, "Synthesize failed, middleware not provided", http.StatusServiceUnavailable), err
	}

	pair, err := this.hoverfly.ApplyMiddlewareIfSet(pair)

	if err != nil {
		return hoverflyError(request, err, "Could not create synthetic response!", http.StatusServiceUnavailable), err
	}

	log.WithFields(log.Fields{
		"mode":        this.hoverfly.Cfg.Mode,
		"middleware":  this.hoverfly.Cfg.Middleware,
		"path":        request.URL.Path,
		"rawQuery":    request.URL.RawQuery,
		"method":      request.Method,
		"destination": request.Host,
	}).Info("synthetic response created successfuly")

	return modes.ReconstructResponse(request, pair), nil
}

type Capture struct {
	hoverfly *Hoverfly
}

func (this Capture) Process(request *http.Request, details models.RequestDetails) (*http.Response, error) {
	// response, err := this.hoverfly.captureRequest(request)

	// this is mainly for testing, since when you create
	if request.Body == nil {
		request.Body = ioutil.NopCloser(bytes.NewBuffer([]byte("")))
	}

	// outputting request body if verbose logging is set
	log.WithFields(log.Fields{
		"body": details.Body,
		"mode": "capture",
	}).Debug("got request body")

	modifiedReq, response, err := this.hoverfly.DoRequest(request)

	if err != nil {
		log.WithFields(log.Fields{
			"error":       err.Error(),
			"mode":        "capture",
			"Path":        request.URL.Path,
			"Method":      request.Method,
			"Destination": request.Host,
			"Scheme":      request.URL.Scheme,
			"Query":       request.URL.RawQuery,
			"Body":        string(details.Body),
			"Headers":     request.Header,
		}).Error("Got error when executing request")

		return nil, err
	}

	requestObj, _ := models.NewRequestDetailsFromHttpRequest(modifiedReq)

	respBody, _ := util.GetResponseBody(response)

	responseObj := &models.ResponseDetails{
		Status:  response.StatusCode,
		Body:    string(respBody),
		Headers: response.Header,
	}

	// saving response body with request/response meta to cache
	this.hoverfly.Save(&requestObj, responseObj)

	if err != nil {
		return hoverflyError(request, err, "Could not capture request", http.StatusServiceUnavailable), err
	}
	log.WithFields(log.Fields{
		"mode":        "capture",
		"middleware":  this.hoverfly.Cfg.Middleware,
		"path":        request.URL.Path,
		"rawQuery":    request.URL.RawQuery,
		"method":      request.Method,
		"destination": request.Host,
	}).Info("request and response captured")

	return response, nil
}

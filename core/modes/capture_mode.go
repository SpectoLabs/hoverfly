package modes

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/SpectoLabs/hoverfly/core/util"

	log "github.com/Sirupsen/logrus"
)

type HoverflyCapture interface {
	ApplyMiddleware(models.RequestResponsePair) (models.RequestResponsePair, error)
	DoRequest(*http.Request) (*http.Response, error)
	Save(*models.RequestDetails, *models.ResponseDetails)
}

type CaptureMode struct {
	Hoverfly HoverflyCapture
}

func (this CaptureMode) Process(request *http.Request, details models.RequestDetails) (*http.Response, error) {
	// this is mainly for testing, since when you create
	if request.Body == nil {
		request.Body = ioutil.NopCloser(bytes.NewBuffer([]byte("")))
	}

	// outputting request body if verbose logging is set
	log.WithFields(log.Fields{
		"body": details.Body,
		"mode": "capture",
	}).Debug("got request body")

	pair, err := this.Hoverfly.ApplyMiddleware(models.RequestResponsePair{Request: details})
	if err != nil {
		return ErrorResponse(request, err, "There was an error when applying middleware to http request"), err
	}

	modifiedRequest, err := ReconstructRequest(pair)
	if err != nil {
		return ErrorResponse(request, err, "There was an error when rebuilding the modified http request"), err
	}

	response, err := this.Hoverfly.DoRequest(modifiedRequest)
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

		return ErrorResponse(request, err, "There was an error when forwarding the request to the intended desintation"), err
	}

	requestObj, err := models.NewRequestDetailsFromHttpRequest(modifiedRequest)
	if err != nil {
		return ErrorResponse(modifiedRequest, err, "There was an error reading the request body"), err
	}

	respBody, _ := util.GetResponseBody(response)

	responseObj := &models.ResponseDetails{
		Status:  response.StatusCode,
		Body:    string(respBody),
		Headers: response.Header,
	}

	// saving response body with request/response meta to cache
	this.Hoverfly.Save(&requestObj, responseObj)

	log.WithFields(log.Fields{
		"mode": "capture",
		// "middleware":  this.Hoverfly.Cfg.Middleware,
		"path":        request.URL.Path,
		"rawQuery":    request.URL.RawQuery,
		"method":      request.Method,
		"destination": request.Host,
	}).Info("request and response captured")

	return response, nil
}

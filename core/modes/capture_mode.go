package modes

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/SpectoLabs/hoverfly/core/util"

	log "github.com/Sirupsen/logrus"
)

type CaptureMode struct {
	Hoverfly Hoverfly
}

func (this CaptureMode) Process(request *http.Request, details models.RequestDetails) (*http.Response, error) {
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

	modifiedReq, response, err := this.Hoverfly.DoRequest(request)

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
	this.Hoverfly.Save(&requestObj, responseObj)

	if err != nil {
		return errorResponse(request, err, "Could not capture request"), err
	}
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

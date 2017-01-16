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
	Save(*models.RequestDetails, *models.ResponseDetails) error
}

type CaptureMode struct {
	Hoverfly HoverflyCapture
}

func (this CaptureMode) Process(request *http.Request, details models.RequestDetails) (*http.Response, error) {
	// this is mainly for testing, since when you create
	if request.Body == nil {
		request.Body = ioutil.NopCloser(bytes.NewBuffer([]byte("")))
	}

	pair, err := this.Hoverfly.ApplyMiddleware(models.RequestResponsePair{Request: details})
	if err != nil {
		return ReturnErrorAndLog(request, err, &pair, "There was an error when applying middleware to http request", "capture")
	}

	modifiedRequest, err := ReconstructRequest(pair)
	if err != nil {
		return ReturnErrorAndLog(request, err, &pair, "There was an error when applying middleware to http request", "capture")
	}

	response, err := this.Hoverfly.DoRequest(modifiedRequest)
	if err != nil {
		return ReturnErrorAndLog(request, err, &pair, "There was an error when forwarding the request to the intended desintation", "capture")
	}

	respBody, _ := util.GetResponseBody(response)

	responseObj := &models.ResponseDetails{
		Status:  response.StatusCode,
		Body:    string(respBody),
		Headers: response.Header,
	}

	// saving response body with request/response meta to cache
	err = this.Hoverfly.Save(&pair.Request, responseObj)
	if err != nil {
		return ReturnErrorAndLog(request, err, &pair, "There was an error when saving request and response", "capture")
	}

	log.WithFields(log.Fields{
		"mode":     "capture",
		"request":  GetRequestLogFields(&pair.Request),
		"response": GetResponseLogFields(&pair.Response),
	}).Info("request and response captured")

	return response, nil
}

package modes

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/SpectoLabs/hoverfly/core/util"

	log "github.com/sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
)

type HoverflyCapture interface {
	ApplyMiddleware(models.RequestResponsePair) (models.RequestResponsePair, error)
	DoRequest(*http.Request) (*http.Response, error)
	Save(*models.RequestDetails, *models.ResponseDetails, []string, bool) error
}

type CaptureMode struct {
	Hoverfly  HoverflyCapture
	Arguments ModeArguments
}

func (this *CaptureMode) View() v2.ModeView {
	return v2.ModeView{
		Mode: Capture,
		Arguments: v2.ModeArgumentsView{
			Headers:          this.Arguments.Headers,
			MatchingStrategy: this.Arguments.MatchingStrategy,
			Stateful:         this.Arguments.Stateful,
		},
	}
}

func (this *CaptureMode) SetArguments(arguments ModeArguments) {
	this.Arguments = arguments
}

func (this *CaptureMode) GetArguments(arguments ModeArguments) {
	this.Arguments = arguments
}

func (this CaptureMode) Process(request *http.Request, details models.RequestDetails) (*http.Response, error) {
	// this is mainly for testing, since when you create
	if request.Body == nil {
		request.Body = ioutil.NopCloser(bytes.NewBuffer([]byte("")))
	}

	pair, err := this.Hoverfly.ApplyMiddleware(models.RequestResponsePair{Request: details})
	if err != nil {
		return ReturnErrorAndLog(request, err, &pair, "There was an error when applying middleware to http request", Capture)
	}

	modifiedRequest, err := ReconstructRequest(pair)
	if err != nil {
		return ReturnErrorAndLog(request, err, &pair, "There was an error when preparing request for pass through", Capture)
	}

	response, err := this.Hoverfly.DoRequest(modifiedRequest)
	if err != nil {
		return ReturnErrorAndLog(request, err, &pair, "There was an error when forwarding the request to the intended destination", Capture)
	}

	respBody, _ := util.GetResponseBody(response)
	respHeaders := util.GetResponseHeaders(response)

	responseObj := &models.ResponseDetails{
		Status:  response.StatusCode,
		Body:    string(respBody),
		Headers: respHeaders,
	}

	if this.Arguments.Headers == nil {
		this.Arguments.Headers = []string{}
	}

	// saving response body with request/response meta to cache
	err = this.Hoverfly.Save(&pair.Request, responseObj, this.Arguments.Headers, this.Arguments.Stateful)
	if err != nil {
		return ReturnErrorAndLog(request, err, &pair, "There was an error when saving request and response", Capture)
	}

	log.WithFields(log.Fields{
		"mode":     Capture,
		"request":  GetRequestLogFields(&pair.Request),
		"response": GetResponseLogFields(&pair.Response),
	}).Info("request and response captured")

	return response, nil
}

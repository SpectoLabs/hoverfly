package modes

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/SpectoLabs/hoverfly/core/util"

	v2 "github.com/SpectoLabs/hoverfly/core/handlers/v2"
	log "github.com/sirupsen/logrus"
)

type HoverflyCapture interface {
	ApplyMiddleware(models.RequestResponsePair) (models.RequestResponsePair, error)
	DoRequest(*http.Request) (*http.Response, *time.Duration, error)
	Save(*models.RequestDetails, *models.ResponseDetails, *ModeArguments) error
}

type CaptureMode struct {
	Hoverfly  HoverflyCapture
	Arguments ModeArguments
}

func (this *CaptureMode) View() v2.ModeView {
	return v2.ModeView{
		Mode: Capture,
		Arguments: v2.ModeArgumentsView{
			Headers:            this.Arguments.Headers,
			Stateful:           this.Arguments.Stateful,
			OverwriteDuplicate: this.Arguments.OverwriteDuplicate,
			CaptureDelay:       this.Arguments.CaptureDelay,
		},
	}
}

func (this *CaptureMode) SetArguments(arguments ModeArguments) {
	this.Arguments = arguments
}

func (this *CaptureMode) GetArguments(arguments ModeArguments) {
	this.Arguments = arguments
}

func (this CaptureMode) Process(request *http.Request, details models.RequestDetails) (ProcessResult, error) {
	request.ParseForm()
	// this is mainly for testing, since when you create
	if request.Body == nil {
		request.Body = io.NopCloser(bytes.NewBuffer([]byte("")))
	}

	pair, err := this.Hoverfly.ApplyMiddleware(models.RequestResponsePair{Request: details})
	if err != nil {
		return ReturnErrorAndLog(request, err, &pair, "There was an error when applying middleware to http request", Capture)
	}

	modifiedRequest, err := ReconstructRequest(pair)
	if err != nil {
		return ReturnErrorAndLog(request, err, &pair, "There was an error when preparing request for pass through", Capture)
	}

	response, duration, err := this.Hoverfly.DoRequest(modifiedRequest)
	if err != nil {
		return ReturnErrorAndLog(request, err, &pair, "There was an error when forwarding the request to the intended destination", Capture)
	}

	respBody, _ := util.GetResponseBody(response)
	respHeaders := util.GetResponseHeaders(response)

	delayInMs := 0
	if this.Arguments.CaptureDelay {
		delayInMs = int(duration.Milliseconds())
	}

	responseObj := &models.ResponseDetails{
		Status:     response.StatusCode,
		Body:       respBody,
		Headers:    respHeaders,
		FixedDelay: delayInMs,
	}

	if this.Arguments.Headers == nil {
		this.Arguments.Headers = []string{}
	}

	// saving response body with request/response meta to cache
	err = this.Hoverfly.Save(&pair.Request, responseObj, &this.Arguments)
	if err != nil {
		return ReturnErrorAndLog(request, err, &pair, "There was an error when saving request and response", Capture)
	}

	log.WithFields(log.Fields{
		"mode":     Capture,
		"request":  GetRequestLogFields(&pair.Request),
		"response": GetResponseLogFields(&pair.Response),
	}).Info("request and response captured")

	return newProcessResult(response, pair.Response.FixedDelay, pair.Response.LogNormalDelay), nil
}

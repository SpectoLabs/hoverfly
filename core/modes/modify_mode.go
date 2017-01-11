package modes

import (
	"io/ioutil"
	"net/http"

	"github.com/SpectoLabs/hoverfly/core/models"

	log "github.com/Sirupsen/logrus"
)

type HoverflyModify interface {
	DoRequest(*http.Request) (*http.Request, *http.Response, error)
	ApplyMiddleware(models.RequestResponsePair) (models.RequestResponsePair, error)
}

type ModifyMode struct {
	Hoverfly HoverflyModify
}

func (this ModifyMode) Process(request *http.Request, details models.RequestDetails) (*http.Response, error) {
	modifiedRequest, resp, err := this.Hoverfly.DoRequest(request)
	if err != nil {
		return errorResponse(request, err, "There was an error when forwarding the request to the intended desintation"), err
	}

	// preparing payload
	bodyBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
			// "middleware": this.hoverfly.Cfg.Middleware,
		}).Error("Failed to read response body after sending modified request")
		return errorResponse(request, err, "Middleware failed or something else happened!"), err
	}

	r := models.ResponseDetails{
		Status:  resp.StatusCode,
		Body:    string(bodyBytes),
		Headers: resp.Header,
	}

	modifiedRequestDetails, err := models.NewRequestDetailsFromHttpRequest(modifiedRequest)
	if err != nil {
		return errorResponse(request, err, "There was an error when reading modified request body"), err
	}

	requestResponsePair := models.RequestResponsePair{Response: r, Request: modifiedRequestDetails}

	newPairs, err := this.Hoverfly.ApplyMiddleware(requestResponsePair)
	if err != nil {
		return errorResponse(request, err, "There was an error when executing middleware"), err
	}

	log.WithFields(log.Fields{
		"status": newPairs.Response.Status,
		// "middleware":  hf.Cfg.Middleware.toString(),
		"mode":        "modify",
		"path":        newPairs.Request.Path,
		"rawQuery":    newPairs.Request.Query,
		"method":      newPairs.Request.Method,
		"destination": newPairs.Request.Destination,
		// original here
		"originalPath":        modifiedRequest.URL.Path,
		"originalRawQuery":    modifiedRequest.URL.RawQuery,
		"originalMethod":      modifiedRequest.Method,
		"originalDestination": modifiedRequest.Host,
	}).Info("request and response modified, returning")

	return ReconstructResponse(modifiedRequest, newPairs), nil
}

package modes

import (
	"io/ioutil"
	"net/http"

	"github.com/SpectoLabs/hoverfly/core/models"

	log "github.com/Sirupsen/logrus"
)

type ModifyMode struct {
	Hoverfly Hoverfly
}

func (this ModifyMode) Process(request *http.Request, details models.RequestDetails) (*http.Response, error) {
	req, resp, err := this.Hoverfly.DoRequest(request)
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

	requestResponsePair := models.RequestResponsePair{Response: r, Request: details}

	newPairs, err := this.Hoverfly.ApplyMiddleware(requestResponsePair)
	if err != nil {
		return errorResponse(request, err, "Middleware failed or something else happened!"), err
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
		"originalPath":        req.URL.Path,
		"originalRawQuery":    req.URL.RawQuery,
		"originalMethod":      req.Method,
		"originalDestination": req.Host,
	}).Info("request and response modified, returning")

	return ReconstructResponse(req, newPairs), nil
}

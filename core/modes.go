package hoverfly

import (
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/SpectoLabs/hoverfly/core/modes"
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

	pair, err := this.hoverfly.ApplyMiddleware(pair)

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

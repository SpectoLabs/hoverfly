package hoverfly

import (
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/core/models"
)

type Mode interface {
	Process(*http.Request, models.RequestDetails) (*http.Response, error)
}

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
	response, err := SynthesizeResponse(request, details, &this.hoverfly.Cfg.Middleware)

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

	return response, nil
}

type Capture struct {
	hoverfly *Hoverfly
}

func (this Capture) Process(request *http.Request, details models.RequestDetails) (*http.Response, error) {
	response, err := this.hoverfly.captureRequest(request)

	if err != nil {
		return hoverflyError(request, err, "Could not capture request", http.StatusServiceUnavailable), err
	}
	log.WithFields(log.Fields{
		"mode":        this.hoverfly.Cfg.Mode,
		"middleware":  this.hoverfly.Cfg.Middleware,
		"path":        request.URL.Path,
		"rawQuery":    request.URL.RawQuery,
		"method":      request.Method,
		"destination": request.Host,
	}).Info("request and response captured")

	return response, nil
}

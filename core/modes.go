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

type Simulate struct {
	hoverfly *Hoverfly
}

func (this Simulate) Process(request *http.Request, details models.RequestDetails) (*http.Response, error) {
	response, err := this.hoverfly.getResponse(request, details)
	if err != nil {
		return hoverflyError(request, err, err.Error(), err.StatusCode), err
	}

	return response, nil
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

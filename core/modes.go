package hoverfly

import (
	"net/http"

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

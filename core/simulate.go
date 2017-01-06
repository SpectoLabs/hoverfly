package hoverfly

import (
	"net/http"

	"github.com/SpectoLabs/hoverfly/core/models"
)

type Simulate struct {
	hoverfly *Hoverfly
}

func (this Simulate) Process(request *http.Request, details models.RequestDetails) (*http.Response, error) {
	response, err := this.hoverfly.GetResponse(details)
	if err != nil {
		return hoverflyError(request, err, err.Error(), err.StatusCode), err
	}

	pair := models.RequestResponsePair{
		Request:  details,
		Response: *response,
	}

	if this.hoverfly.Cfg.Middleware.IsSet() {
		pair, _ = this.hoverfly.Cfg.Middleware.Execute(pair)
	}

	c := NewConstructor(request, pair)

	return c.ReconstructResponse(), nil
}

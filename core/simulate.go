package hoverfly

import (
	"net/http"

	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/SpectoLabs/hoverfly/core/modes"
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

	pair, _ = this.hoverfly.ApplyMiddlewareIfSet(pair)
	// TODO: If there is an error, should Hoverfly return an error via http.Response
	// or should it just log.Error the message and return the original pair?

	return modes.ReconstructResponse(request, pair), nil
}

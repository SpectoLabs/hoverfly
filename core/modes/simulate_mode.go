package modes

import (
	"net/http"

	"github.com/SpectoLabs/hoverfly/core/models"
)

type SimulateMode struct {
	Hoverfly Hoverfly
}

func (this SimulateMode) Process(request *http.Request, details models.RequestDetails) (*http.Response, error) {
	response, err := this.Hoverfly.GetResponse(details)
	if err != nil {
		return errorResponse(request, err, err.Error(), err.StatusCode), err
	}

	pair := models.RequestResponsePair{
		Request:  details,
		Response: *response,
	}

	pair, _ = this.Hoverfly.ApplyMiddlewareIfSet(pair)
	// TODO: If there is an error, should Hoverfly return an error via http.Response
	// or should it just log.Error the message and return the original pair?

	return ReconstructResponse(request, pair), nil
}

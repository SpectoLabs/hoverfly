package modes

import (
	"net/http"

	"github.com/SpectoLabs/hoverfly/core/models"
)

type SimulateMode struct {
	Hoverfly Hoverfly
}

func (this SimulateMode) Process(request *http.Request, details models.RequestDetails) (*http.Response, error) {
	response, matchingErr := this.Hoverfly.GetResponse(details)
	if matchingErr != nil {
		return errorResponse(request, matchingErr, matchingErr.Error(), matchingErr.StatusCode), matchingErr
	}

	pair := models.RequestResponsePair{
		Request:  details,
		Response: *response,
	}

	pair, err := this.Hoverfly.ApplyMiddlewareIfSet(pair)
	if err != nil {
		return errorResponse(request, err, "Error when executing middleware", http.StatusServiceUnavailable), err
	}

	return ReconstructResponse(request, pair), nil
}

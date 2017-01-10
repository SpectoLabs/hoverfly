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
		return errorResponse(request, matchingErr, "There was an error when matching"), matchingErr
	}

	pair := models.RequestResponsePair{
		Request:  details,
		Response: *response,
	}

	pair, err := this.Hoverfly.ApplyMiddleware(pair)
	if err != nil {
		return errorResponse(request, err, "There was an error when executing middleware"), err
	}

	return ReconstructResponse(request, pair), nil
}

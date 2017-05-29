package modes

import (
	"net/http"

	"github.com/SpectoLabs/hoverfly/core/matching"
	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
)

type HoverflySimulate interface {
	GetResponse(models.RequestDetails, bool) (*models.ResponseDetails, *matching.MatchingError)
	ApplyMiddleware(models.RequestResponsePair) (models.RequestResponsePair, error)
}

type SimulateMode struct {
	Hoverfly HoverflySimulate
	MatchingStrategy string
}

func (this *SimulateMode) View() v2.ModeView {
	return v2.ModeView{
		Mode: Simulate,
		Arguments: v2.ModeArgumentsView{
			MatchingStrategy: &this.MatchingStrategy,
		},
	}
}

func (this *SimulateMode) SetArguments(arguments ModeArguments) {
	if arguments.MatchingStrategy == nil {
		this.MatchingStrategy = "STRONGEST"
	} else {
		this.MatchingStrategy = *arguments.MatchingStrategy
	}
}

func (this SimulateMode) Process(request *http.Request, details models.RequestDetails) (*http.Response, error) {
	pair := models.RequestResponsePair{
		Request: details,
	}

	response, matchingErr := this.Hoverfly.GetResponse(details, this.MatchingStrategy != "FIRST")
	if matchingErr != nil {
		return ReturnErrorAndLog(request, matchingErr, &pair, "There was an error when matching", Simulate)
	}

	pair.Response = *response

	pair, err := this.Hoverfly.ApplyMiddleware(pair)
	if err != nil {
		return ReturnErrorAndLog(request, err, &pair, "There was an error when executing middleware", Simulate)
	}

	return ReconstructResponse(request, pair), nil
}

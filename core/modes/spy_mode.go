package modes

import (
	"net/http"

	"github.com/SpectoLabs/hoverfly/core/errors"

	log "github.com/sirupsen/logrus"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/models"
)

type HoverflySpy interface {
	GetResponse(models.RequestDetails) (*models.ResponseDetails, *errors.HoverflyError)
	ApplyMiddleware(models.RequestResponsePair) (models.RequestResponsePair, error)
	DoRequest(*http.Request) (*http.Response, error)
}

type SpyMode struct {
	Hoverfly         HoverflySpy
	MatchingStrategy string
}

func (this *SpyMode) View() v2.ModeView {
	return v2.ModeView{
		Mode: Spy,
		Arguments: v2.ModeArgumentsView{
			MatchingStrategy: &this.MatchingStrategy,
		},
	}
}

func (this *SpyMode) SetArguments(arguments ModeArguments) {
	if arguments.MatchingStrategy == nil {
		this.MatchingStrategy = "strongest"
	} else {
		this.MatchingStrategy = *arguments.MatchingStrategy
	}
}

//TODO: We should only need one of these two parameters
func (this SpyMode) Process(request *http.Request, details models.RequestDetails) (*http.Response, error) {
	pair := models.RequestResponsePair{
		Request: details,
	}

	response, matchingErr := this.Hoverfly.GetResponse(details)

	if matchingErr != nil {
		log.Info("Going to call real server")
		modifiedRequest, err := ReconstructRequest(pair)
		if err == nil {
			response, err := this.Hoverfly.DoRequest(modifiedRequest)
			if err == nil {
				log.Info("Going to return response from real server")
				return response, nil
			}
		}
	}

	if matchingErr != nil {
		return ReturnErrorAndLog(request, matchingErr, &pair, "There was an error when matching", Simulate)
	}

	pair.Response = *response

	if pair, err := this.Hoverfly.ApplyMiddleware(pair); err == nil {
		return ReconstructResponse(request, pair), nil
	} else {
		return ReturnErrorAndLog(request, err, &pair, "There was an error when executing middleware", Simulate)
	}
}

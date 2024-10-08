package modes

import (
	"net/http"
	"time"

	"github.com/SpectoLabs/hoverfly/core/errors"
	v2 "github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/util"

	log "github.com/sirupsen/logrus"

	"github.com/SpectoLabs/hoverfly/core/models"
)

type HoverflySpy interface {
	GetResponse(models.RequestDetails) (*models.ResponseDetails, *errors.HoverflyError)
	ApplyMiddleware(models.RequestResponsePair) (models.RequestResponsePair, error)
	DoRequest(*http.Request) (*http.Response, *time.Duration, error)
	Save(*models.RequestDetails, *models.ResponseDetails, *ModeArguments) error
}

type SpyMode struct {
	Hoverfly  HoverflySpy
	Arguments ModeArguments
}

func (this *SpyMode) View() v2.ModeView {
	return v2.ModeView{
		Mode: Spy,
		Arguments: v2.ModeArgumentsView{
			MatchingStrategy:   this.Arguments.MatchingStrategy,
			CaptureOnMiss:      this.Arguments.CaptureOnMiss,
			Stateful:           this.Arguments.Stateful,
			Headers:            this.Arguments.Headers,
			OverwriteDuplicate: this.Arguments.OverwriteDuplicate,
		},
	}
}

func (this *SpyMode) SetArguments(arguments ModeArguments) {
	var matchingStrategy string
	if arguments.MatchingStrategy == nil || *arguments.MatchingStrategy == "" {
		matchingStrategy = "strongest"
	} else {
		matchingStrategy = *arguments.MatchingStrategy
	}
	this.Arguments = ModeArguments{
		MatchingStrategy:   &matchingStrategy,
		Headers:            arguments.Headers,
		Stateful:           arguments.Stateful,
		OverwriteDuplicate: arguments.OverwriteDuplicate,
		CaptureOnMiss:      arguments.CaptureOnMiss,
	}
}

// TODO: We should only need one of these two parameters
func (this SpyMode) Process(request *http.Request, details models.RequestDetails) (ProcessResult, error) {
	pair := models.RequestResponsePair{
		Request: details,
	}

	response, matchingErr := this.Hoverfly.GetResponse(details)

	if matchingErr != nil {
		log.Info("Going to call real server")
		modifiedRequest, err := ReconstructRequest(pair)
		if err != nil {
			return ReturnErrorAndLog(request, err, &pair, "There was an error when reconstructing the request.", Spy)
		}
		response, duration, err := this.Hoverfly.DoRequest(modifiedRequest)
		if err == nil {

			if this.Arguments.CaptureOnMiss {
				respBody, _ := util.GetResponseBody(response)
				respHeaders := util.GetResponseHeaders(response)

				responseObj := &models.ResponseDetails{
					Status:     response.StatusCode,
					Body:       respBody,
					Headers:    respHeaders,
					FixedDelay: int(duration.Milliseconds()),
				}
				if this.Arguments.Headers == nil {
					this.Arguments.Headers = []string{}
				}
				err = this.Hoverfly.Save(&pair.Request, responseObj, &this.Arguments)
				if err != nil {
					return ReturnErrorAndLog(request, err, &pair, "There was an error when saving request and response", Spy)
				}
			}
			log.Info("Going to return response from real server")
			return newProcessResult(response, 0, nil), nil
		} else {
			return ReturnErrorAndLog(request, err, &pair, "There was an error when forwarding the request to the intended destination", Spy)
		}
	}

	pair.Response = *response

	pair, err := this.Hoverfly.ApplyMiddleware(pair)
	if err != nil {
		return ReturnErrorAndLog(request, err, &pair, "There was an error when executing middleware", Spy)
	}

	return newProcessResult(
		ReconstructResponse(request, pair),
		pair.Response.FixedDelay,
		pair.Response.LogNormalDelay,
	), nil
}

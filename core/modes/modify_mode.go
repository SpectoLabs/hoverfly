package modes

import (
	"io/ioutil"
	"math"
	"net/http"
	"time"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/models"
)

type HoverflyModify interface {
	ApplyMiddleware(models.RequestResponsePair) (models.RequestResponsePair, error)
	DoRequest(*http.Request) (*http.Response, *time.Duration, error)
}

type ModifyMode struct {
	Hoverfly HoverflyModify
}

func (this *ModifyMode) View() v2.ModeView {
	return v2.ModeView{
		Mode: Modify,
	}
}

func (this *ModifyMode) SetArguments(arguments ModeArguments) {}

func (this ModifyMode) Process(request *http.Request, details models.RequestDetails) (ProcessResult, error) {
	pair, err := this.Hoverfly.ApplyMiddleware(models.RequestResponsePair{Request: details})
	if err != nil {
		return ReturnErrorAndLog(request, err, &pair, "There was an error when executing middleware", Modify)
	}

	modifiedRequest, err := ReconstructRequest(pair)
	if err != nil {
		return ReturnErrorAndLog(request, err, &pair, "There was an error when rebuilding the modified http request", Modify)
	}

	resp, duration, err := this.Hoverfly.DoRequest(modifiedRequest)
	if err != nil {
		return ReturnErrorAndLog(request, err, &pair, "There was an error when forwarding the request to the intended destination", Modify)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ReturnErrorAndLog(request, err, &pair, "There was an error when reading the http response body", Modify)
	}

	pair.Response = models.ResponseDetails{
		Status:     resp.StatusCode,
		Body:       string(bodyBytes),
		Headers:    resp.Header,
		FixedDelay: int(math.Ceil(duration.Seconds())),
	}

	pair, err = this.Hoverfly.ApplyMiddleware(pair)
	if err != nil {
		return ReturnErrorAndLog(request, err, &pair, "There was an error when executing middleware", Modify)
	}

	response := ReconstructResponse(modifiedRequest, pair)
	return newProcessResult(response, pair.Response.FixedDelay, pair.Response.LogNormalDelay), nil
}

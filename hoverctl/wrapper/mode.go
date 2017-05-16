package wrapper

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/hoverctl/configuration"
)

// GetMode will go the state endpoint in Hoverfly, parse the JSON response and return the mode of Hoverfly
func GetMode(target configuration.Target) (string, error) {
	response, err := doRequest(target, "GET", v2ApiMode, "", nil)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	apiResponse := createAPIStateResponse(response)

	return apiResponse.Mode, nil
}

// Set will go the state endpoint in Hoverfly, sending JSON that will set the mode of Hoverfly
func SetModeWithArguments(target configuration.Target, modeView v2.ModeView) (string, error) {
	if modeView.Mode != "simulate" && modeView.Mode != "capture" &&
		modeView.Mode != "modify" && modeView.Mode != "synthesize" {
		return "", errors.New(modeView.Mode + " is not a valid mode")
	}
	bytes, err := json.Marshal(modeView)
	if err != nil {
		return "", err
	}

	response, err := doRequest(target, "PUT", v2ApiMode, string(bytes), nil)
	if err != nil {
		return "", err
	}

	if response.StatusCode == http.StatusBadRequest {
		return "", handlerError(response)
	}

	apiResponse := createAPIStateResponse(response)

	return apiResponse.Mode, nil
}

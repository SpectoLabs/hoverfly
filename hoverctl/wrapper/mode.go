package wrapper

import (
	"encoding/json"
	"errors"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/hoverctl/configuration"
)

// GetMode will go the state endpoint in Hoverfly, parse the JSON response and return the mode of Hoverfly
func GetMode(target configuration.Target) (*v2.ModeView, error) {
	response, err := doRequest(target, "GET", v2ApiMode, "", nil)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	err = handleResponseError(response, "Could not retrieve mode")
	if err != nil {
		return nil, err
	}

	var modeView v2.ModeView

	err = UnmarshalToInterface(response, &modeView)
	if err != nil {
		return nil, err
	}

	return &modeView, nil
}

// Set will go the state endpoint in Hoverfly, sending JSON that will set the mode of Hoverfly
func SetModeWithArguments(target configuration.Target, modeView *v2.ModeView) (string, error) {
	if modeView.Mode != "simulate" && modeView.Mode != "capture" &&
		modeView.Mode != "modify" && modeView.Mode != "synthesize" &&
		modeView.Mode != "spy" && modeView.Mode != "diff" {
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

	err = handleResponseError(response, "Could not set mode")
	if err != nil {
		return "", err
	}

	var modeViewResponse v2.ModeView

	err = UnmarshalToInterface(response, &modeViewResponse)
	if err != nil {
		return "", err
	}

	return modeViewResponse.Mode, nil
}

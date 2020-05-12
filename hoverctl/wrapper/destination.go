package wrapper

import (
	"encoding/json"
	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/hoverctl/configuration"
)

// GetDestination will go the destination endpoint in Hoverfly, parse the JSON response and return the destination of Hoverfly
func GetDestination(target configuration.Target) (string, error) {
	response, err := doRequest(target, "GET", v2ApiDestination, "", nil)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	err = handleResponseError(response, "Could not retrieve destination")
	if err != nil {
		return "", err
	}

	var destinationView v2.DestinationView

	err = UnmarshalToInterface(response, &destinationView)
	if err != nil {
		return "", err
	}

	return destinationView.Destination, nil
}

// SetDestination will go the destination endpoint in Hoverfly, sending JSON that will set the destination of Hoverfly
func SetDestination(target configuration.Target, destination string) (string, error) {

	destinationReq := map[string]string{"destination": destination}
	bytes, _ := json.Marshal(destinationReq) // JSON encode in case there are special chars
	reqBody := string(bytes)

	response, err := doRequest(target, "PUT", v2ApiDestination, reqBody, nil)
	if err != nil {
		return "", err
	}

	err = handleResponseError(response, "Could not set destination")
	if err != nil {
		return "", err
	}

	var destinationView v2.DestinationView

	err = UnmarshalToInterface(response, &destinationView)
	if err != nil {
		return "", err
	}

	return destinationView.Destination, nil
}

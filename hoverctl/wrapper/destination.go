package wrapper

import "github.com/SpectoLabs/hoverfly/hoverctl/configuration"

// GetDestination will go the destination endpoint in Hoverfly, parse the JSON response and return the destination of Hoverfly
func GetDestination(target configuration.Target) (string, error) {
	response, err := doRequest(target, "GET", v2ApiDestination, "", nil)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	apiResponse := createAPIStateResponse(response)

	return apiResponse.Destination, nil
}

// SetDestination will go the destination endpoint in Hoverfly, sending JSON that will set the destination of Hoverfly
func SetDestination(target configuration.Target, destination string) (string, error) {
	response, err := doRequest(target, "PUT", v2ApiDestination, `{"destination":"`+destination+`"}`, nil)
	if err != nil {
		return "", err
	}

	apiResponse := createAPIStateResponse(response)

	return apiResponse.Destination, nil
}

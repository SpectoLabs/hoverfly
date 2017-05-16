package wrapper

import (
	"encoding/json"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/hoverctl/configuration"
)

// GetMiddle will go the middleware endpoint in Hoverfly, parse the JSON response and return the middleware of Hoverfly
func GetMiddleware(target configuration.Target) (v2.MiddlewareView, error) {
	response, err := doRequest(target, "GET", v2ApiMiddleware, "", nil)
	if err != nil {
		return v2.MiddlewareView{}, err
	}

	defer response.Body.Close()

	middlewareResponse := createMiddlewareSchema(response)

	return middlewareResponse, nil
}

func SetMiddleware(target configuration.Target, binary, script, remote string) (v2.MiddlewareView, error) {
	middlewareRequest := &v2.MiddlewareView{
		Binary: binary,
		Script: script,
		Remote: remote,
	}

	marshalledMiddleware, err := json.Marshal(middlewareRequest)
	if err != nil {
		return v2.MiddlewareView{}, err
	}

	response, err := doRequest(target, "PUT", v2ApiMiddleware, string(marshalledMiddleware), nil)
	if err != nil {
		return v2.MiddlewareView{}, err
	}

	err = handleResponseError(response, "Hoverfly could not execute this middleware")
	if err != nil {
		return v2.MiddlewareView{}, err
	}

	apiResponse := createMiddlewareSchema(response)

	return apiResponse, nil
}

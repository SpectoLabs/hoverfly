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

	err = handleResponseError(response, "Could not retrieve middleware")
	if err != nil {
		return v2.MiddlewareView{}, err
	}

	var middlewareView v2.MiddlewareView

	err = UnmarshalToInterface(response, &middlewareView)
	if err != nil {
		return v2.MiddlewareView{}, err
	}

	return middlewareView, nil
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

	defer response.Body.Close()

	err = handleResponseError(response, "Could not set middleware, it may have failed the test")
	if err != nil {
		return v2.MiddlewareView{}, err
	}

	var middlewareView v2.MiddlewareView

	err = UnmarshalToInterface(response, &middlewareView)
	if err != nil {
		return v2.MiddlewareView{}, err
	}

	return middlewareView, nil
}

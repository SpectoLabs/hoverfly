package wrapper

import (
	"encoding/json"

	v2 "github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/hoverctl/configuration"
)

func GetAllPostServeActions(target configuration.Target) (v2.PostServeActionDetailsView, error) {

	response, err := doRequest(target, "GET", v2ApiPostServeAction, "", nil)
	if err != nil {
		return v2.PostServeActionDetailsView{}, err
	}

	defer response.Body.Close()

	err = handleResponseError(response, "Could not retrieve all post serve actions")
	if err != nil {
		return v2.PostServeActionDetailsView{}, err
	}

	var postServeActionDetailsView v2.PostServeActionDetailsView

	err = UnmarshalToInterface(response, &postServeActionDetailsView)
	if err != nil {
		return v2.PostServeActionDetailsView{}, err
	}

	return postServeActionDetailsView, nil
}

func SetPostServeAction(actionName, binary, scriptContent string, delayInMs int, target configuration.Target) error {

	actionRequest := v2.ActionView{
		ActionName:    actionName,
		Binary:        binary,
		ScriptContent: scriptContent,
		DelayInMs:     delayInMs,
	}
	marshalledAction, err := json.Marshal(actionRequest)
	if err != nil {
		return err
	}

	response, err := doRequest(target, "PUT", v2ApiPostServeAction, string(marshalledAction), nil)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	err = handleResponseError(response, "Could not set post serve action")
	if err != nil {
		return err
	}
	return nil
}

func DeletePostServeAction(actionName string, target configuration.Target) error {

	response, err := doRequest(target, "DELETE", v2ApiPostServeAction+"/"+actionName, "", nil)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	err = handleResponseError(response, "Could not delete post serve action")
	if err != nil {
		return err
	}

	return nil
}

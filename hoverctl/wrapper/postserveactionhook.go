package wrapper

import (
	"encoding/json"

	v2 "github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/hoverctl/configuration"
)

func SetPostServeActionHook(hookName, binary, scriptContent string, delayInMilliSeconds int, target configuration.Target) error {

	hookRequest := v2.HookView{
		HookName:            hookName,
		Binary:              binary,
		ScriptContent:       scriptContent,
		DelayInMilliSeconds: delayInMilliSeconds,
	}
	marshalledHook, err := json.Marshal(hookRequest)
	if err != nil {
		return err
	}

	_, err = doRequest(target, "POST", v2ApiPostServeActionHook, string(marshalledHook), nil)
	if err != nil {
		return err
	}
	return nil
}

func DeletePostServeActionHook(hookName string, target configuration.Target) error {

	_, err := doRequest(target, "DELETE", v2ApiPostServeActionHook+"?name="+hookName, "", nil)
	if err != nil {
		return err
	}
	return nil
}

package wrapper

import (
	"encoding/json"
	"io/ioutil"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/hoverctl/configuration"
)

func GetCurrentState(target configuration.Target) (map[string]string, error) {

	res, err := doRequest(target, "GET", v2ApiState, "", nil)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	currentState := &v2.StateView{}

	err = json.Unmarshal(responseBody, currentState)

	if err != nil {
		return nil, err
	}

	return currentState.State, nil
}

func PatchCurrentState(target configuration.Target, key, value string) error {

	marshal, err := json.Marshal(&v2.StateView{
		State: map[string]string{
			key: value,
		},
	})

	if err != nil {
		return err
	}

	_, err = doRequest(target, "PATCH", v2ApiState, string(marshal), nil)

	return err
}

func DeleteCurrentState(target configuration.Target) error {

	_, err := doRequest(target, "DELETE", v2ApiState, "", nil)

	return err
}

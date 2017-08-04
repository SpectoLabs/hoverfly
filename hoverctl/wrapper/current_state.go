package wrapper

import (
	"encoding/json"
	ioutil "io/ioutil"

	"github.com/SpectoLabs/hoverfly/hoverctl/configuration"
)

func GetCurrentState(target configuration.Target) (map[string]string, error) {

	res, err := doRequest(target, "GET", v2ApiState, "", nil)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	bytes, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	var currentState = make(map[string]string)

	err = json.Unmarshal(bytes, &currentState)

	if err != nil {
		return nil, err
	}

	return currentState, nil
}

func PatchCurrentState(target configuration.Target, key, value string) error {

	toPatch := make(map[string]string)
	toPatch[key] = value

	marshal, err := json.Marshal(toPatch)

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

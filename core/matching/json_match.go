package matching

import (
	"encoding/json"
	"reflect"
)

func JsonMatch(matchingString string, toMatch string) bool {
	var matchingObject map[string]interface{}
	err := json.Unmarshal([]byte(matchingString), &matchingObject)
	if err != nil {
		return false
	}

	var toMatchObject map[string]interface{}
	err = json.Unmarshal([]byte(toMatch), &toMatchObject)
	if err != nil {
		return false
	}

	return reflect.DeepEqual(matchingObject, toMatchObject)
}

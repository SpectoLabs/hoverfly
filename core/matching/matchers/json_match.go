package matchers

import (
	"encoding/json"
	"reflect"
)

var Json = "json"

func JsonMatch(match interface{}, toMatch string) bool {
	matchString, ok := match.(string)
	if !ok {
		return false
	}

	if matchString == toMatch {
		return true
	}
	var matchingObject map[string]interface{}
	err := json.Unmarshal([]byte(matchString), &matchingObject)
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

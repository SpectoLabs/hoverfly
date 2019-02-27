package matchers

import (
	"encoding/json"
	"reflect"
)

var Json = "json"

func JsonMatch(match interface{}, toMatch string) (matched bool, result string) {
	matchString, ok := match.(string)
	if !ok {
		return
	}

	if matchString == toMatch {
		return true, toMatch
	}
	var matchingObject map[string]interface{}
	err := json.Unmarshal([]byte(matchString), &matchingObject)
	if err != nil {
		return
	}

	var toMatchObject map[string]interface{}
	err = json.Unmarshal([]byte(toMatch), &toMatchObject)
	if err != nil {
		return
	}

	matched = reflect.DeepEqual(matchingObject, toMatchObject)
	if matched {
		result = toMatch
	}

	return
}

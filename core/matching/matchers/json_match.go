package matchers

import (
	"encoding/json"
	"reflect"
)

var Json = "json"

func JsonMatch(match interface{}, toMatch string, config map[string]interface{}) (string, bool) {
	matchString, ok := match.(string)
	if !ok {
		return "", false
	}

	if matchString == toMatch {
		return toMatch, true
	}
	var matchingObject interface{}
	err := json.Unmarshal([]byte(matchString), &matchingObject)
	if err != nil {
		return "", false
	}

	var toMatchObject interface{}
	err = json.Unmarshal([]byte(toMatch), &toMatchObject)
	if err != nil {
		return "", false
	}

	if isMatched := reflect.DeepEqual(matchingObject, toMatchObject); isMatched {
		return toMatch, true
	}
	return "", false
}

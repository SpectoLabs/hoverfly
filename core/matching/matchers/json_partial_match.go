package matchers

import "encoding/json"

var JsonPartial = "jsonpartial"

func JsonPartialMatch(match interface{}, toMatch string) bool {
	var expected, actual map[string]interface{}
	matchString, ok := match.(string)
	if !ok {
		return false
	}
	err0 := json.Unmarshal([]byte(matchString), &expected)
	err1 := json.Unmarshal([]byte(toMatch), &actual)
	if err0 != nil || err1 != nil {
		return false
	}
	return isMapContainsMap(expected, actual)
}

func isMapContainsMap(match, toMatch interface{}) bool {
	var expected, actual map[string]interface{}
	expected, ok0 := match.(map[string]interface{})
	actual, ok1 := toMatch.(map[string]interface{})
	if !ok0 || !ok1 {
		return false
	}

	for key, val := range expected {
		if innerMap, ok := val.(map[string]interface{}); ok {
			result := isMapContainsMap(innerMap, actual[key])
			if !result {
				return false
			}
		} else {
			_, exist := actual[key]
			if !exist {
				return false
			}
		}
	}
	return true
}

package matching

import (
	"encoding/json"
	"reflect"

	"github.com/SpectoLabs/hoverfly/core/util"
)

func JsonMatch(matchingString string, toMatch string) bool {
	minifiedMatchingString, err := util.MinifyJson(matchingString)
	if err != nil {
		return false
	}

	minifiedToMatch, err := util.MinifyJson(toMatch)
	if err != nil {
		return false
	}

	var matchingJson, toMatchJson map[string]interface{}

	err = json.Unmarshal([]byte(minifiedMatchingString), &matchingJson)
	if err != nil {
		return false
	}

	err = json.Unmarshal([]byte(minifiedToMatch), &toMatchJson)
	if err != nil {
		return false
	}

	return reflect.DeepEqual(matchingJson, toMatchJson)
}

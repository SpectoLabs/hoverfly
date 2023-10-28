package matchers

import (
	"github.com/SpectoLabs/hoverfly/core/util"
)

var JsonPath = "jsonpath"

func JsonPathMatch(match interface{}, toMatch string) bool {
	matchString, ok := match.(string)
	if !ok {
		return false
	}

	matchString = util.PrepareJsonPathQuery(matchString)
	returnedString, err := util.JsonPathExecution(matchString, toMatch)
	if err != nil || returnedString == matchString {
		return false
	}

	return true
}

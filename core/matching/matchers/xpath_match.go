package matchers

import (
	"github.com/SpectoLabs/hoverfly/core/util"
)

var Xpath = "xpath"

func XpathMatch(match interface{}, toMatch string) bool {
	matchString, ok := match.(string)
	if !ok {
		return false
	}

	results, err := util.XpathExecution(matchString, toMatch)
	if err != nil {
		return false
	}

	return results.Bool()
}

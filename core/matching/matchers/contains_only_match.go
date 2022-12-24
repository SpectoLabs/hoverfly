package matchers

import (
	"strings"

	"github.com/SpectoLabs/hoverfly/core/util"
)

var ContainsOnly = "containsonly"

func ContainsOnlyMatch(match interface{}, toMatch string) bool {
	matchStringArr, ok := match.([]string)
	if !ok {
		return false
	}
	toMatchArr := strings.Split(toMatch, ";")
	return util.ContainsOnly(matchStringArr, toMatchArr)
}

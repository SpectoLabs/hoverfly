package matchers

import (
	"strings"

	"github.com/SpectoLabs/hoverfly/core/util"
)

var Contains = "contains"

func ContainsMatch(match interface{}, toMatch string) bool {
	matchStringArr, ok := match.([]string)
	if !ok {
		return false
	}
	toMatchArr := strings.Split(toMatch, ";")
	return util.Contains(matchStringArr, toMatchArr)
}

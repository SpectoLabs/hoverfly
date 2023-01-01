package matchers

import (
	"strings"

	"github.com/SpectoLabs/hoverfly/core/util"
)

var ContainsExactly = "containsexactly"

func ContainsExactlyMatch(match interface{}, toMatch string) bool {
	matchStringArr, ok := util.GetStringArray(match)
	if !ok {
		return false
	}
	toMatchArr := strings.Split(toMatch, ";")
	return util.Identical(matchStringArr, toMatchArr)
}

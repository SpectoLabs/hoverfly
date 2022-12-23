package matchers

import (
	"sort"
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
	sort.Strings(matchStringArr)
	sort.Strings(toMatchArr)
	return util.Identical(matchStringArr, toMatchArr)
}

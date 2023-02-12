package matchers

import (
	"strings"

	"github.com/SpectoLabs/hoverfly/core/util"
)

const (
	IgnoreUnknown     = "ignoreUnknown"
	IgnoreOrder       = "ignoreOrder"
	IgnoreOccurrences = "ignoreOccurrences"
)

var Array = "array"

func ArrayMatchWithoutConfig(match interface{}, toMatch string) bool {

	return ArrayMatch(match, toMatch, nil)
}

func ArrayMatch(match interface{}, toMatch string, config map[string]interface{}) bool {
	matchStringArr, ok := util.GetStringArray(match)
	if !ok {
		return false
	}
	toMatchArr := strings.Split(toMatch, ";")
	ignoreUnknown := util.GetBoolOrDefault(config, IgnoreUnknown, false)
	ignoreOrder := util.GetBoolOrDefault(config, IgnoreOrder, false)
	ignoreOccurrences := util.GetBoolOrDefault(config, IgnoreOccurrences, false)

	return (ignoreUnknown || hasAllKnown(matchStringArr, toMatchArr)) &&
		(ignoreOccurrences || hasSameNoOfOccurrences(matchStringArr, toMatchArr)) &&
		(ignoreOrder || isInSameOrder(matchStringArr, toMatchArr))
}

func hasSameNoOfOccurrences(matchGroup, toMatch []string) bool {

	matchGroupSet := make(map[string]int)
	for _, value := range matchGroup {
		matchGroupSet[value] = matchGroupSet[value] + 1
	}
	toMatchSet := make(map[string]int)
	for _, value := range toMatch {
		toMatchSet[value] = toMatchSet[value] + 1
	}

	for key, value := range matchGroupSet {
		if toMatchSet[key] != value {
			return false
		}
	}
	return true
}

func hasAllKnown(matchGroup, toMatch []string) bool {

	matchGroupSet := make(map[string]bool)
	for _, value := range matchGroup {
		matchGroupSet[value] = true
	}

	for _, value := range toMatch {
		if _, found := matchGroupSet[value]; !found {
			return false
		}
	}
	return true
}

func isInSameOrder(matchGroup, toMatch []string) bool {

	index := 0
	for _, value := range toMatch {
		if index < len(matchGroup) && value == matchGroup[index] {
			index++
		} else if index == len(matchGroup) {
			break
		}
	}
	return index == len(matchGroup)
}

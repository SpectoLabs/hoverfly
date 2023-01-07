package matchers

import (
	"strings"

	"github.com/SpectoLabs/hoverfly/core/util"
)

const (
	IGNORE_UNKNOWN     = "ignoreunknown"
	IGNORE_ORDER       = "ignoreorder"
	IGNORE_OCCURRENCES = "ignoreoccurrences"
)

var Array = "array"

func ArrayMatch(match interface{}, toMatch string, config map[string]interface{}) (string, bool) {
	matchStringArr, ok := util.GetStringArray(match)
	if !ok {
		return "", false
	}
	toMatchArr := strings.Split(toMatch, ";")
	ignoreUnknown := util.GetBoolOrDefault(config, IGNORE_UNKNOWN, false)
	ignoreOrder := util.GetBoolOrDefault(config, IGNORE_ORDER, false)
	ignoreOccurrences := util.GetBoolOrDefault(config, IGNORE_OCCURRENCES, false)

	isMatched := (ignoreUnknown || hasAllKnown(matchStringArr, toMatchArr)) &&
		(ignoreOccurrences || hasSameNoOfOccurrences(matchStringArr, toMatchArr)) &&
		(ignoreOrder || isInSameOrder(matchStringArr, toMatchArr))
	if isMatched {
		return toMatch, isMatched
	}
	return "", false
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

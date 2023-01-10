package matchers

import (
	"encoding/json"
	"fmt"

	"github.com/SpectoLabs/hoverfly/core/util"
	log "github.com/sirupsen/logrus"
)

const (
	IGNORE_UNKNOWN     = "ignoreunknown"
	IGNORE_ORDER       = "ignoreorder"
	IGNORE_OCCURRENCES = "ignoreoccurrences"
)

var Array = "array"

func ArrayMatch(match interface{}, toMatch string, config map[string]interface{}) bool {
	matchStringArr, ok := util.GetStringArray(match)
	if !ok {
		return false
	}
	toMatchArr, err := unMarshalArray(toMatch)
	if err != nil {
		return false
	}
	ignoreUnknown := util.GetBoolOrDefault(config, IGNORE_UNKNOWN, false)
	ignoreOrder := util.GetBoolOrDefault(config, IGNORE_ORDER, false)
	ignoreOccurrences := util.GetBoolOrDefault(config, IGNORE_OCCURRENCES, false)

	return (ignoreUnknown || hasAllKnown(matchStringArr, toMatchArr)) &&
		(ignoreOccurrences || hasSameNoOfOccurrences(matchStringArr, toMatchArr)) &&
		(ignoreOrder || isInSameOrder(matchStringArr, toMatchArr))
}

func unMarshalArray(jsonArrStr string) ([]string, error) {

	var arr []interface{}
	if err := json.Unmarshal([]byte(jsonArrStr), &arr); err != nil {
		log.Errorf("Cannot unmarshal to array %s", err.Error())
		return []string{}, err
	}
	var strArr []string
	for _, value := range arr {
		strArr = append(strArr, fmt.Sprint(value))
	}
	return strArr, nil
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

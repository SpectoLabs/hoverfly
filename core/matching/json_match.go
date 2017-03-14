package matching

import (
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

	return reflect.DeepEqual(minifiedMatchingString, minifiedToMatch)
}

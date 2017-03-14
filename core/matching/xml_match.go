package matching

import (
	"reflect"

	"github.com/SpectoLabs/hoverfly/core/util"
)

func XmlMatch(matchingString string, toMatch string) bool {
	minifiedMatchingString, err := util.MinifyXml(matchingString)
	if err != nil {
		return false
	}

	minifiedToMatch, err := util.MinifyXml(toMatch)
	if err != nil {
		return false
	}

	return reflect.DeepEqual(minifiedMatchingString, minifiedToMatch)
}

package matchers

import (
	"reflect"

	"github.com/SpectoLabs/hoverfly/core/util"
)

var Xml = "xml"

func XmlMatch(match interface{}, toMatch string) bool {
	matchString, ok := match.(string)
	if !ok {
		return false
	}

	minifiedMatch, err := util.MinifyXml(matchString)
	if err != nil {
		return false
	}

	minifiedToMatch, err := util.MinifyXml(toMatch)
	if err != nil {
		return false
	}

	return reflect.DeepEqual(minifiedMatch, minifiedToMatch)
}

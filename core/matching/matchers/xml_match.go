package matchers

import (
	"reflect"

	"github.com/SpectoLabs/hoverfly/core/util"
)

var Xml = "xml"

func XmlMatch(match interface{}, toMatch string) (matched bool, result string) {
	matchString, ok := match.(string)
	if !ok {
		return
	}

	minifiedMatch, err := util.MinifyXml(matchString)
	if err != nil {
		return
	}

	minifiedToMatch, err := util.MinifyXml(toMatch)
	if err != nil {
		return
	}

	matched = reflect.DeepEqual(minifiedMatch, minifiedToMatch)

	if matched {
		result = toMatch
	}
	return
}

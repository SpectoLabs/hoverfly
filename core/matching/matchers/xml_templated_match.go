package matchers

import (
	"regexp"

	"github.com/beevik/etree"
)

var XmlTemplated = "xmltemplated"

var ignoreExpr = regexp.MustCompile("^\\s*{{\\s*ignore\\s*}}\\s*$")
var regExpr = regexp.MustCompile("^\\s*{{\\s*regex:(.*)}}\\s*$")

func XmlTemplatedMatch(match interface{}, toMatch string) bool {
	matchString, ok := match.(string)
	if !ok {
		return false
	}

	// parse xml in mock data into dom tree
	expected := etree.NewDocument()
	if err := expected.ReadFromString(matchString); err != nil {
		return false
	}

	// parse xml in actual request body into dom tree
	actual := etree.NewDocument()
	if err := actual.ReadFromString(toMatch); err != nil {
		return false
	}

	// tree matching
	return compareTree(expected.Root(), actual.Root())
}

func compareTree(expected *etree.Element, actual *etree.Element) bool {
	// compare constructure
	// step 1. compare tag name
	if expected.Tag != actual.Tag {
		return false
	}
	// step 2. compare node content
	// case 1: leaf
	if isLeaf(expected) {
		// compare text content
		return compareValue(expected.Text(), actual.Text())
	}
	// case 2: children element matching
	actualChildren := actual.ChildElements()
	// for each expected tag
	for _, match := range expected.ChildElements() {
		// find one in actual
		matched := false
		for i, ele := range actualChildren {
			if compareTree(match, ele) {
				matched = true
				// remove matched
				actualChildren = append(actualChildren[:i], actualChildren[i+1:]...)
				break
			}
		}
		// all elements in actual data is not matched
		// or, too many elements in expected data
		if matched == false {
			return false
		}
	}
	// too many elements in actual data
	if len(actualChildren) > 0 {
		return false
	}
	return true
}

// check element text content
func compareValue(expected string, actual string) bool {
	// pattern 1: ignore value  => always be true
	if ignoreExpr.MatchString(expected) {
		return true
	}
	// pattern 2: regex
	// parse node content
	group := regExpr.FindStringSubmatch(expected)
	// if it matches the template like {{regex: ... }}  ==> otherwise, take it as plain text
	if len(group) > 1 {
		matcher, err := regexp.Compile(group[1])
		// can not compile regular expression --> invalid regex --> false
		if err != nil {
			return false
		}
		// use regular expression to match actual value
		return matcher.MatchString(actual)
	}
	// pattern 3: exact equal
	return expected == actual
}

// check if an element is leaf
func isLeaf(expected *etree.Element) bool {
	// it is a leaf node if it does not have any children
	return len(expected.ChildElements()) == 0
}

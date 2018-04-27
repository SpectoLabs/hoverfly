package matchers

import "regexp"

var Regex = "regex"

func RegexMatch(match interface{}, toMatch string) bool {
	matchString, ok := match.(string)
	if !ok {
		return false
	}

	result, err := regexp.MatchString(matchString, toMatch)
	if err != nil {
		return false
	}

	return result
}

package matchers

import "regexp"

var Regex = "regex"

func RegexMatch(match interface{}, toMatch string, config map[string]interface{}) (string, bool) {
	matchString, ok := match.(string)
	if !ok {
		return "", false
	}

	result, err := regexp.MatchString(matchString, toMatch)
	if err != nil {
		return "", false
	}

	if result {
		return toMatch, result
	}
	return "", false
}

package matchers

import "regexp"

var Regex = "regex"

func RegexMatch(match interface{}, toMatch string) (matched bool, result string) {
	matchString, ok := match.(string)
	if !ok {
		return
	}

	matched, err := regexp.MatchString(matchString, toMatch)
	if err != nil {
		return
	}

	if matched {
		result = toMatch
	}

	return
}

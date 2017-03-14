package matching

import "regexp"

func RegexMatch(matchingString string, toMatch string) bool {
	match, err := regexp.MatchString(matchingString, toMatch)
	if err != nil {
		return false
	}

	return match
}

package matchers

var Exact = "exact"

func ExactMatch(match interface{}, toMatch string, config map[string]interface{}) (string, bool) {
	matchString, ok := match.(string)
	if !ok {
		return "", false
	}

	if matchString == toMatch {
		return toMatch, true
	}
	return "", false
}

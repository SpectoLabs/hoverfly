package matchers

var Exact = "exact"

func ExactMatch(match interface{}, toMatch string) (matched bool, result string) {
	matchString, ok := match.(string)
	if !ok {
		return
	}

	matched = matchString == toMatch

	if matched {
		result = toMatch
	}
	return
}

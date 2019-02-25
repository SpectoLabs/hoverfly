package matchers

var Exact = "exact"

func ExactMatch(match interface{}, toMatch string) (bool, string) {
	matchString, ok := match.(string)
	if !ok {
		return false, ""
	}

	// TODO only return toMatch string if the matching is true
	return matchString == toMatch, toMatch
}

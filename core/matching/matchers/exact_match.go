package matchers

func ExactMatch(match interface{}, toMatch string) bool {
	matchString, ok := match.(string)
	if !ok {
		return false
	}

	return matchString == toMatch
}

package matchers

func ExactMatch(matchingString string, toMatch string) bool {
	return matchingString == toMatch
}

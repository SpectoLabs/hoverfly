package matchers

var ContainsExactly = "containsexactly"

func ContainsExactlyMatch(match interface{}, toMatch string) bool {
	return ArrayMatch(match, toMatch, nil)
}

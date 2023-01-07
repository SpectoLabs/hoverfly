package matchers

var ContainsExactly = "containsexactly"

func ContainsExactlyMatch(match interface{}, toMatch string, config map[string]interface{}) (string, bool) {
	return ArrayMatch(match, toMatch, nil)
}

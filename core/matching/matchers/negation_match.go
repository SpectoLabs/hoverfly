package matchers

var Negation = "negation"

func NegationMatch(match interface{}, toMatch string) bool {
	matchString, ok := match.(string)
	if ok {
		return matchString != toMatch
	}
	return true
}

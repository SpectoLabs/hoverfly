package matchers

import "github.com/ryanuber/go-glob"

var Glob = "glob"

func GlobMatch(match interface{}, toMatch string) (matched bool, result string) {
	matchString, ok := match.(string)
	if !ok {
		return
	}

	matched = glob.Glob(matchString, toMatch)
	if matched {
		result = toMatch
	}
	return
}

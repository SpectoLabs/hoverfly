package matchers

import "github.com/ryanuber/go-glob"

var Glob = "glob"

func GlobMatch(match interface{}, toMatch string, config map[string]interface{}) (string, bool) {
	matchString, ok := match.(string)
	if !ok {
		return "", false
	}

	if matched := glob.Glob(matchString, toMatch); matched {
		return toMatch, true
	}
	return "", false
}

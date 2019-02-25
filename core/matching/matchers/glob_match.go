package matchers

import "github.com/ryanuber/go-glob"

var Glob = "glob"

func GlobMatch(match interface{}, toMatch string) (bool, string) {
	matchString, ok := match.(string)
	if !ok {
		return false, ""
	}

	return glob.Glob(matchString, toMatch), toMatch
}

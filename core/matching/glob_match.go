package matching

import glob "github.com/ryanuber/go-glob"

func GlobMatch(matchingString string, toMatch string) bool {
	return glob.Glob(matchingString, toMatch)
}

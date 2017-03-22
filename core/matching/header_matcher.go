package matching

import (
	"strings"

	glob "github.com/ryanuber/go-glob"
)

func HeaderMatcher(matchingHeaders, toMatch map[string][]string) bool {

	for matcherHeaderKey, matcherHeaderValues := range matchingHeaders {
		for requestHeaderKey, requestHeaderValues := range toMatch {
			delete(toMatch, requestHeaderKey)
			toMatch[strings.ToLower(requestHeaderKey)] = requestHeaderValues

		}

		toMatchHeaderValues, toMatchHeaderValuesFound := toMatch[strings.ToLower(matcherHeaderKey)]
		if !toMatchHeaderValuesFound {
			return false
		}

		for _, matcherHeaderValue := range matcherHeaderValues {
			matcherHeaderValueMatched := false
			for _, toMatchHeaderValue := range toMatchHeaderValues {
				if glob.Glob(strings.ToLower(matcherHeaderValue), strings.ToLower(toMatchHeaderValue)) {
					matcherHeaderValueMatched = true
				}
			}

			if !matcherHeaderValueMatched {
				return false
			}
		}
	}
	return true
}

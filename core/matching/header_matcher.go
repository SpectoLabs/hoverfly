package matching

import (
	"strings"

	glob "github.com/ryanuber/go-glob"
	"fmt"
)

func CountlessHeaderMatcher(matchingHeaders, toMatch map[string][]string) * FieldMatch {

	for matcherHeaderKey, matcherHeaderValues := range matchingHeaders {

		// Make everything lowercase, as headers are case insensitive
		for requestHeaderKey, requestHeaderValues := range toMatch {
			delete(toMatch, requestHeaderKey)
			toMatch[strings.ToLower(requestHeaderKey)] = requestHeaderValues
		}

		toMatchHeaderValues, toMatchHeaderValuesFound := toMatch[strings.ToLower(matcherHeaderKey)]
		if !toMatchHeaderValuesFound {
			return countlessFieldMatch(false)
		}

		for _, matcherHeaderValue := range matcherHeaderValues {
			fmt.Println(toMatchHeaderValues)

			matcherHeaderValueMatched := false
			for _, toMatchHeaderValue := range toMatchHeaderValues {
				if glob.Glob(strings.ToLower(matcherHeaderValue), strings.ToLower(toMatchHeaderValue)) {
					matcherHeaderValueMatched = true
				}
			}

			if !matcherHeaderValueMatched {
				return countlessFieldMatch(false)
			}
		}
	}
	return countlessFieldMatch(true)
}

func CountingHeaderMatcher(matchingHeaders, toMatch map[string][]string) * FieldMatch {

	matched := true
	var matchCount int

	for matcherHeaderKey, matcherHeaderValues := range matchingHeaders {

		// Make everything lowercase, as headers are case insensitive
		for requestHeaderKey, requestHeaderValues := range toMatch {
			delete(toMatch, requestHeaderKey)
			toMatch[strings.ToLower(requestHeaderKey)] = requestHeaderValues
		}

		toMatchHeaderValues, found := toMatch[strings.ToLower(matcherHeaderKey)]
		if !found {
			matched = false
		}

		for _, matcherHeaderValue := range matcherHeaderValues {
			matcherHeaderValueMatched := false
			for _, toMatchHeaderValue := range toMatchHeaderValues {
				if glob.Glob(strings.ToLower(matcherHeaderValue), strings.ToLower(toMatchHeaderValue)) {
					matcherHeaderValueMatched = true
					matchCount++
				}
			}

			if !matcherHeaderValueMatched {
				matched = false
			}
		}
	}
	return &FieldMatch{
		Matched: matched,
		TotalMatches: matchCount,
	}
}
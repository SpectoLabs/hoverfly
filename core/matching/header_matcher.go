package matching

import (
	"strings"

	"github.com/SpectoLabs/hoverfly/core/models"
)

func HeaderMatching(requestMatcher models.RequestMatcher, toMatch map[string][]string) *FieldMatch {

	// // Make everything lowercase, as headers are case insensitive
	// for requestHeaderKey, requestHeaderValues := range toMatch {
	// 	delete(toMatch, requestHeaderKey)
	// 	toMatch[strings.ToLower(requestHeaderKey)] = requestHeaderValues
	// }

	matched := true
	var score int

	requestMatcherHeadersWithMatchers := requestMatcher.Headers

	for matcherHeaderKey, matcherHeaderValue := range requestMatcherHeadersWithMatchers {
		// Make everything lowercase, as headers are case insensitive
		for requestHeaderKey, requestHeaderValues := range toMatch {
			delete(toMatch, requestHeaderKey)
			toMatch[strings.ToLower(requestHeaderKey)] = requestHeaderValues
		}
		matcherHeaderValueMatched := false

		toMatchHeaderValues, found := toMatch[strings.ToLower(matcherHeaderKey)]
		if !found {
			matched = false
			continue
		}

		fieldMatch := FieldMatcher(matcherHeaderValue, strings.Join(toMatchHeaderValues, ";"))
		matcherHeaderValueMatched = fieldMatch.Matched
		score += fieldMatch.Score

		if !matcherHeaderValueMatched {
			matched = false
		}
	}

	return &FieldMatch{
		Matched: matched,
		Score:   score,
	}
}

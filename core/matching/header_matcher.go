package matching

import (
	"strings"

	"github.com/SpectoLabs/hoverfly/core/models"
	glob "github.com/ryanuber/go-glob"
)

func HeaderMatching(requestMatcher models.RequestMatcher, toMatch map[string][]string) *FieldMatch {

	// // Make everything lowercase, as headers are case insensitive
	// for requestHeaderKey, requestHeaderValues := range toMatch {
	// 	delete(toMatch, requestHeaderKey)
	// 	toMatch[strings.ToLower(requestHeaderKey)] = requestHeaderValues
	// }

	matched := true
	var score int

	requestMatcherHeadersWithMatchers := requestMatcher.HeadersWithMatchers

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

	requestMatcherHeaders := requestMatcher.Headers

	for matcherHeaderKey, matcherHeaderValues := range requestMatcherHeaders {
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
					score++
				}
			}

			if !matcherHeaderValueMatched {
				matched = false
			}
		}
	}
	return &FieldMatch{
		Matched: matched,
		Score:   score,
	}
}

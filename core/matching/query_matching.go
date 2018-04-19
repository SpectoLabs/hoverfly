package matching

import (
	"strings"

	"github.com/SpectoLabs/hoverfly/core/models"
)

func QueryMatching(requestMatcher models.RequestMatcher, toMatch map[string][]string) *FieldMatch {

	matched := true
	var matchScore int

	requestMatcherQueriesWithMatchers := requestMatcher.QueriesWithMatchers

	for matcherQueryKey, matcherQueryValue := range requestMatcherQueriesWithMatchers {
		matcherHeaderValueMatched := false

		toMatchQueryValues, found := toMatch[strings.ToLower(matcherQueryKey)]
		if !found {
			matched = false
			continue
		}

		fieldMatch := ScoredFieldMatcher(matcherQueryValue, strings.Join(toMatchQueryValues, ";"))
		matcherHeaderValueMatched = fieldMatch.Matched
		matchScore += fieldMatch.MatchScore

		if !matcherHeaderValueMatched {
			matched = false
		}
	}

	return &FieldMatch{
		Matched:    matched,
		MatchScore: matchScore,
	}
}

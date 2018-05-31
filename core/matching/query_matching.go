package matching

import (
	"strings"

	"github.com/SpectoLabs/hoverfly/core/models"
)

func QueryMatching(requestMatcher models.RequestMatcher, toMatch map[string][]string) *FieldMatch {

	matched := true
	var score int

	for key, value := range toMatch {
		toMatch[strings.ToLower(key)] = value
	}

	for matcherQueryKey, matcherQueryValue := range requestMatcher.Query {
		matcherHeaderValueMatched := false

		toMatchQueryValues, found := toMatch[strings.ToLower(matcherQueryKey)]
		if !found {
			matched = false
			continue
		}

		fieldMatch := FieldMatcher(matcherQueryValue, strings.Join(toMatchQueryValues, ";"))
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

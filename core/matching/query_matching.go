package matching

import (
	"strings"

	"github.com/SpectoLabs/hoverfly/core/models"
)

func QueryMatching(requestMatcher models.RequestMatcher, toMatch map[string][]string) *FieldMatch {

	matched := true
	var score int

	if requestMatcher.Query == nil {
		return &FieldMatch{
			Matched: true,
			Score:   1,
		}
	}

	if len(*requestMatcher.Query) == 0 {
		if len(toMatch) == 0 {
			return &FieldMatch{
				Matched: true,
				Score:   1,
			}
		}
		return &FieldMatch{
			Matched: false,
			Score:   0,
		}
	}

	lowercaseKeyMap := make(map[string][]string)
	for key, value := range toMatch {
		lowercaseKeyMap[strings.ToLower(key)] = value
	}

	for matcherQueryKey, matcherQueryValue := range *requestMatcher.Query {
		matcherHeaderValueMatched := false

		toMatchQueryValues, found := lowercaseKeyMap[strings.ToLower(matcherQueryKey)]
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

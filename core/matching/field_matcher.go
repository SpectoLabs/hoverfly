package matching

import (
	"strings"

	"github.com/SpectoLabs/hoverfly/core/matching/matchers"
	"github.com/SpectoLabs/hoverfly/core/models"
)

func FieldMatcher(fields []models.RequestFieldMatchers, toMatch string) *FieldMatch {

	fieldMatch := &FieldMatch{Matched: true}

	if fields == nil || len(fields) == 0 {
		return fieldMatch
	}

	for _, field := range fields {
		if isMatching(field, toMatch) {
			if field.Matcher == matchers.Exact || (field.Matcher == matchers.Array && field.Config == nil) {
				fieldMatch.Score = fieldMatch.Score + 2
			} else {
				fieldMatch.Score = fieldMatch.Score + 1
			}
		} else {
			fieldMatch.Matched = false
		}
	}

	return fieldMatch
}

func isMatching(field models.RequestFieldMatchers, toMatch string) bool {
	currentMatcher := field
	actual := toMatch
	result := false
	for {

		var matcherDetails matchers.MatcherDetails
		isMatched := false
		if currentMatcher.Config == nil {
			matcherDetails = matchers.Matchers[strings.ToLower(currentMatcher.Matcher)]
			isMatched = matcherDetails.MatcherFunction.(func(interface{}, string) bool)(currentMatcher.Value, actual)

		} else {
			matcherDetails = matchers.MatchersWithConfig[strings.ToLower(currentMatcher.Matcher)]
			isMatched = matcherDetails.MatcherFunction.(func(interface{}, string, map[string]interface{}) bool)(currentMatcher.Value, actual, currentMatcher.Config)

		}
		if !isMatched {
			return false
		}
		/* it ll break if match value generator is nil.. incase where we are matching complete details(exact match, containsexactlymatch, jsonmatch or xmlmatch)
		no need of matcher chaining in such scenarios and if it is there then it will be ignored
		*/
		if currentMatcher.DoMatch == nil || matcherDetails.MatchValueGenerator == nil {
			result = isMatched
			break
		}
		actual = matcherDetails.MatchValueGenerator(currentMatcher.Value, actual)
		currentMatcher = *currentMatcher.DoMatch

	}
	return result
}

type FieldMatch struct {
	Matched bool
	Score   int
}

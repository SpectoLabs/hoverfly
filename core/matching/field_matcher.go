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
			if field.Matcher == matchers.Exact || field.Matcher == matchers.ContainsExactly {
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
	isMatched := false
	if len(field.Config) > 0 {
		isMatched = matchers.MatchersWithConfig[strings.ToLower(field.Matcher)](field.Value, toMatch, field.Config)
	} else {
		isMatched = matchers.Matchers[strings.ToLower(field.Matcher)](field.Value, toMatch)
	}
	return isMatched
}

type FieldMatch struct {
	Matched bool
	Score   int
}

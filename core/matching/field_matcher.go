package matching

import (
	"github.com/SpectoLabs/hoverfly/core/matching/matchers"
	"github.com/SpectoLabs/hoverfly/core/models"
)

func FieldMatcher(fields []models.RequestFieldMatchers, toMatch string) *FieldMatch {

	fieldMatch := &FieldMatch{Matched: true}

	if fields == nil || len(fields) == 0 {
		return fieldMatch
	}

	for _, field := range fields {
		if matchers.Matchers[field.Matcher](field.Value, toMatch) {
			if field.Matcher == matchers.Exact {
				fieldMatch.MatchScore = fieldMatch.MatchScore + 2
			} else {
				fieldMatch.MatchScore = fieldMatch.MatchScore + 1
			}
		} else {
			fieldMatch.Matched = false
		}
	}

	return fieldMatch
}

func FieldMatchWithNoScore(matched bool) *FieldMatch {
	return &FieldMatch{
		Matched:    matched,
		MatchScore: 0,
	}
}

type FieldMatch struct {
	Matched    bool
	MatchScore int
}

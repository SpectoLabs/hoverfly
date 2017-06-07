package matching

import (
	"github.com/SpectoLabs/hoverfly/core/models"
)

func UnscoredFieldMatcher(field *models.RequestFieldMatchers, toMatch string) *FieldMatch {
	if field == nil {
		return FieldMatchWithNoScore(true)
	}

	if field.ExactMatch != nil && !ExactMatch(*field.ExactMatch, toMatch) {
		return FieldMatchWithNoScore(false)
	}

	if field.XmlMatch != nil && !XmlMatch(*field.XmlMatch, toMatch) {
		return FieldMatchWithNoScore(false)
	}

	if field.XpathMatch != nil && !XpathMatch(*field.XpathMatch, toMatch) {
		return FieldMatchWithNoScore(false)
	}

	if field.JsonMatch != nil && !JsonMatch(*field.JsonMatch, toMatch) {
		return FieldMatchWithNoScore(false)
	}

	if field.JsonPathMatch != nil && !JsonPathMatch(*field.JsonPathMatch, toMatch) {
		return FieldMatchWithNoScore(false)
	}

	if field.RegexMatch != nil && !RegexMatch(*field.RegexMatch, toMatch) {
		return FieldMatchWithNoScore(false)
	}

	if field.GlobMatch != nil && !GlobMatch(*field.GlobMatch, toMatch) {
		return FieldMatchWithNoScore(false)
	}

	return FieldMatchWithNoScore(true)
}

func ScoredFieldMatcher(field *models.RequestFieldMatchers, toMatch string) *FieldMatch {

	fieldMatch := &FieldMatch{Matched: true}

	if field == nil {
		return fieldMatch
	}

	if field.ExactMatch != nil {
		if ExactMatch(*field.ExactMatch, toMatch) {
			fieldMatch.MatchScore++
		} else {
			fieldMatch.Matched = false
		}
	}

	if field.XmlMatch != nil {
		if XmlMatch(*field.XmlMatch, toMatch) {
			fieldMatch.MatchScore++
		} else {
			fieldMatch.Matched = false
		}
	}

	if field.XpathMatch != nil {
		if XpathMatch(*field.XpathMatch, toMatch) {
			fieldMatch.MatchScore++
		} else {
			fieldMatch.Matched = false
		}
	}

	if field.JsonMatch != nil {
		if JsonMatch(*field.JsonMatch, toMatch) {
			fieldMatch.MatchScore++
		} else {
			fieldMatch.Matched = false
		}
	}

	if field.JsonPathMatch != nil {
		if JsonPathMatch(*field.JsonPathMatch, toMatch) {
			fieldMatch.MatchScore++
		} else {
			fieldMatch.Matched = false
		}
	}

	if field.RegexMatch != nil {
		if RegexMatch(*field.RegexMatch, toMatch) {
			fieldMatch.MatchScore++
		} else {
			fieldMatch.Matched = false
		}
	}

	if field.GlobMatch != nil {
		if GlobMatch(*field.GlobMatch, toMatch) {
			fieldMatch.MatchScore++
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

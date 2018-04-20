package matching

import (
	"github.com/SpectoLabs/hoverfly/core/matching/matchers"
	"github.com/SpectoLabs/hoverfly/core/models"
)

func UnscoredFieldMatcher(field *models.RequestFieldMatchers, toMatch string) *FieldMatch {
	if field == nil {
		return FieldMatchWithNoScore(true)
	}

	switch field.Matcher {
	case "exact":
		if !matchers.ExactMatch(field.Value.(string), toMatch) {
			return FieldMatchWithNoScore(false)
		}
	case "xml":
		if !matchers.XmlMatch(field.Value.(string), toMatch) {
			return FieldMatchWithNoScore(false)
		}
	case "xpath":
		if !matchers.XpathMatch(field.Value.(string), toMatch) {
			return FieldMatchWithNoScore(false)
		}
	case "json":
		if !matchers.JsonMatch(field.Value.(string), toMatch) {
			return FieldMatchWithNoScore(false)
		}
	case "jsonpath":
		if !matchers.JsonPathMatch(field.Value.(string), toMatch) {
			return FieldMatchWithNoScore(false)
		}
	case "regex":
		if !matchers.RegexMatch(field.Value.(string), toMatch) {
			return FieldMatchWithNoScore(false)
		}
	case "glob":
		if !matchers.GlobMatch(field.Value.(string), toMatch) {
			return FieldMatchWithNoScore(false)
		}
	}

	return FieldMatchWithNoScore(true)
}

func ScoredFieldMatcher(field *models.RequestFieldMatchers, toMatch string) *FieldMatch {

	fieldMatch := &FieldMatch{Matched: true}

	if field == nil {
		return fieldMatch
	}

	if field.ExactMatch != nil {
		if matchers.ExactMatch(*field.ExactMatch, toMatch) {
			fieldMatch.MatchScore = fieldMatch.MatchScore + 2
		} else {
			fieldMatch.Matched = false
		}
	}

	if field.XmlMatch != nil {
		if matchers.XmlMatch(*field.XmlMatch, toMatch) {
			fieldMatch.MatchScore++
		} else {
			fieldMatch.Matched = false
		}
	}

	if field.XpathMatch != nil {
		if matchers.XpathMatch(*field.XpathMatch, toMatch) {
			fieldMatch.MatchScore++
		} else {
			fieldMatch.Matched = false
		}
	}

	if field.JsonMatch != nil {
		if matchers.JsonMatch(*field.JsonMatch, toMatch) {
			fieldMatch.MatchScore++
		} else {
			fieldMatch.Matched = false
		}
	}

	if field.JsonPathMatch != nil {
		if matchers.JsonPathMatch(*field.JsonPathMatch, toMatch) {
			fieldMatch.MatchScore++
		} else {
			fieldMatch.Matched = false
		}
	}

	if field.RegexMatch != nil {
		if matchers.RegexMatch(*field.RegexMatch, toMatch) {
			fieldMatch.MatchScore++
		} else {
			fieldMatch.Matched = false
		}
	}

	if field.GlobMatch != nil {
		if matchers.GlobMatch(*field.GlobMatch, toMatch) {
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

package matching

import (
	"github.com/SpectoLabs/hoverfly/core/matching/matchers"
	"github.com/SpectoLabs/hoverfly/core/models"
)

func UnscoredFieldMatcher(fields []models.RequestFieldMatchers, toMatch string) *FieldMatch {
	if fields == nil || len(fields) == 0 {
		return FieldMatchWithNoScore(true)
	}

	field := fields[0]

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

func ScoredFieldMatcher(fields []models.RequestFieldMatchers, toMatch string) *FieldMatch {

	fieldMatch := &FieldMatch{Matched: true}

	if fields == nil || len(fields) == 0 {
		return fieldMatch
	}

	field := fields[0]

	switch field.Matcher {
	case "exact":
		if matchers.ExactMatch(field.Value.(string), toMatch) {
			fieldMatch.MatchScore = fieldMatch.MatchScore + 2
		} else {
			fieldMatch.Matched = false
		}
	case "xml":
		if matchers.XmlMatch(field.Value.(string), toMatch) {
			fieldMatch.MatchScore = fieldMatch.MatchScore + 1
		} else {
			fieldMatch.Matched = false
		}
	case "xpath":
		if matchers.XpathMatch(field.Value.(string), toMatch) {
			fieldMatch.MatchScore = fieldMatch.MatchScore + 1
		} else {
			fieldMatch.Matched = false
		}
	case "json":
		if matchers.JsonMatch(field.Value.(string), toMatch) {
			fieldMatch.MatchScore = fieldMatch.MatchScore + 1
		} else {
			fieldMatch.Matched = false
		}
	case "jsonpath":
		if matchers.JsonPathMatch(field.Value.(string), toMatch) {
			fieldMatch.MatchScore = fieldMatch.MatchScore + 1
		} else {
			fieldMatch.Matched = false
		}
	case "regex":
		if matchers.RegexMatch(field.Value.(string), toMatch) {
			fieldMatch.MatchScore = fieldMatch.MatchScore + 1
		} else {
			fieldMatch.Matched = false
		}
	case "glob":
		if matchers.GlobMatch(field.Value.(string), toMatch) {
			fieldMatch.MatchScore = fieldMatch.MatchScore + 1
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

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

	field := fields[0]

	switch field.Matcher {
	case matchers.Exact:
		if matchers.ExactMatch(field.Value, toMatch) {
			fieldMatch.MatchScore = fieldMatch.MatchScore + 2
		} else {
			fieldMatch.Matched = false
		}
	case matchers.Xml:
		if matchers.XmlMatch(field.Value, toMatch) {
			fieldMatch.MatchScore = fieldMatch.MatchScore + 1
		} else {
			fieldMatch.Matched = false
		}
	case matchers.Xpath:
		if matchers.XpathMatch(field.Value, toMatch) {
			fieldMatch.MatchScore = fieldMatch.MatchScore + 1
		} else {
			fieldMatch.Matched = false
		}
	case matchers.Json:
		if matchers.JsonMatch(field.Value, toMatch) {
			fieldMatch.MatchScore = fieldMatch.MatchScore + 1
		} else {
			fieldMatch.Matched = false
		}
	case matchers.JsonPath:
		if matchers.JsonPathMatch(field.Value, toMatch) {
			fieldMatch.MatchScore = fieldMatch.MatchScore + 1
		} else {
			fieldMatch.Matched = false
		}
	case matchers.Regex:
		if matchers.RegexMatch(field.Value, toMatch) {
			fieldMatch.MatchScore = fieldMatch.MatchScore + 1
		} else {
			fieldMatch.Matched = false
		}
	case matchers.Glob:
		if matchers.GlobMatch(field.Value, toMatch) {
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

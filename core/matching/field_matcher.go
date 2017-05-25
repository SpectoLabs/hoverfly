package matching

import (
	"github.com/SpectoLabs/hoverfly/core/models"
)

func CountlessFieldMatcher(field *models.RequestFieldMatchers, toMatch string) *FieldMatch {
	if field == nil {
		return countlessFieldMatch(true)
	}

	if field.ExactMatch != nil && !ExactMatch(*field.ExactMatch, toMatch) {
		return countlessFieldMatch(false)
	}

	if field.XmlMatch != nil && !XmlMatch(*field.XmlMatch, toMatch) {
		return countlessFieldMatch(false)
	}

	if field.XpathMatch != nil && !XpathMatch(*field.XpathMatch, toMatch) {
		return countlessFieldMatch(false)
	}

	if field.JsonMatch != nil && !JsonMatch(*field.JsonMatch, toMatch) {
		return countlessFieldMatch(false)
	}

	if field.JsonPathMatch != nil && !JsonPathMatch(*field.JsonPathMatch, toMatch) {
		return countlessFieldMatch(false)
	}

	if field.RegexMatch != nil && !RegexMatch(*field.RegexMatch, toMatch) {
		return countlessFieldMatch(false)
	}

	if field.GlobMatch != nil && !GlobMatch(*field.GlobMatch, toMatch) {
		return countlessFieldMatch(false)
	}

	return countlessFieldMatch(true)
}

func CountingFieldMatcher(field *models.RequestFieldMatchers, toMatch string) *FieldMatch {

	fieldMatch := &FieldMatch{Matched: true}

	if field == nil {
		return fieldMatch
	}

	if field.ExactMatch != nil {
		if ExactMatch(*field.ExactMatch, toMatch) {
			fieldMatch.TotalMatches++
		} else {
			fieldMatch.Matched = false
		}
	}

	if field.XmlMatch != nil {
		if XmlMatch(*field.XmlMatch, toMatch) {
			fieldMatch.TotalMatches++
		} else {
			fieldMatch.Matched = false
		}
	}

	if field.XpathMatch != nil {
		if XpathMatch(*field.XpathMatch, toMatch) {
			fieldMatch.TotalMatches++
		} else {
			fieldMatch.Matched = false
		}
	}

	if field.JsonMatch != nil {
		if JsonMatch(*field.JsonMatch, toMatch) {
			fieldMatch.TotalMatches++
		} else {
			fieldMatch.Matched = false
		}
	}

	if field.JsonPathMatch != nil {
		if JsonPathMatch(*field.JsonPathMatch, toMatch) {
			fieldMatch.TotalMatches++
		} else {
			fieldMatch.Matched = false
		}
	}

	if field.RegexMatch != nil {
		if RegexMatch(*field.RegexMatch, toMatch) {
			fieldMatch.TotalMatches++
		} else {
			fieldMatch.Matched = false
		}
	}

	if field.GlobMatch != nil {
		if GlobMatch(*field.GlobMatch, toMatch) {
			fieldMatch.TotalMatches++
		} else {
			fieldMatch.Matched = false
		}
	}

	return fieldMatch
}

func countlessFieldMatch(matched bool) *FieldMatch {
	return &FieldMatch{
		Matched:      matched,
		TotalMatches: 0,
	}
}

type FieldMatch struct {
	Matched      bool
	TotalMatches int
}

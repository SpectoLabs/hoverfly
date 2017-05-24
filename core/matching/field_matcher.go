package matching

import (
	"github.com/SpectoLabs/hoverfly/core/models"
)

func FieldMatcher(field *models.RequestFieldMatchers, toMatch string) * FieldMatch {
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

func countlessFieldMatch(matched bool)  * FieldMatch  {
	return &FieldMatch{
		Matched: matched,
		TotalMatches: 0,
	}
}

type FieldMatch struct {
	Matched bool
	TotalMatches int
}
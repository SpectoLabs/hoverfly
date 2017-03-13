package matching

import (
	"github.com/SpectoLabs/hoverfly/core/models"
)

func FieldMatcher(field *models.RequestFieldMatchers, toMatch string) bool {
	if field == nil {
		return true
	}

	if field.ExactMatch != nil {
		return ExactMatch(field.ExactMatch, toMatch)
	}

	if field.XpathMatch != nil {
		return XpathMatch(field.XpathMatch, toMatch)
	}

	if field.JsonPathMatch != nil {
		return JsonMatch(field.JsonPathMatch, toMatch)
	}

	if field.RegexMatch != nil {
		return RegexMatch(field.RegexMatch, toMatch)
	}

	if field.GlobMatch != nil {
		return GlobMatch(field.GlobMatch, toMatch)
	}

	return false
}

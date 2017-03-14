package matching

import (
	"github.com/SpectoLabs/hoverfly/core/models"
)

func FieldMatcher(field *models.RequestFieldMatchers, toMatch string) bool {
	if field == nil {
		return true
	}

	if field.ExactMatch != nil && !ExactMatch(*field.ExactMatch, toMatch) {
		return false
	}

	if field.XmlMatch != nil && !XmlMatch(*field.XmlMatch, toMatch) {
		return false
	}

	if field.XpathMatch != nil && !XpathMatch(*field.XpathMatch, toMatch) {
		return false
	}

	if field.JsonMatch != nil && !JsonMatch(*field.JsonMatch, toMatch) {
		return false
	}

	if field.JsonPathMatch != nil && !JsonPathMatch(*field.JsonPathMatch, toMatch) {
		return false
	}

	if field.RegexMatch != nil && !RegexMatch(*field.RegexMatch, toMatch) {
		return false
	}

	if field.GlobMatch != nil && !GlobMatch(*field.GlobMatch, toMatch) {
		return false
	}

	return true
}

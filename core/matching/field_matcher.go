package matching

import (
	"github.com/SpectoLabs/hoverfly/core/models"
)

func FieldMatcher(field *models.RequestFieldMatchers, toMatch string) bool {
	matches := []bool{}
	if field == nil {
		return true
	}

	if field.ExactMatch != nil {
		matches = append(matches, ExactMatch(*field.ExactMatch, toMatch))
	}

	if field.XpathMatch != nil {
		matches = append(matches, XpathMatch(field.XpathMatch, toMatch))
	}

	if field.JsonPathMatch != nil {
		matches = append(matches, JsonMatch(field.JsonPathMatch, toMatch))
	}

	if field.RegexMatch != nil {
		matches = append(matches, RegexMatch(field.RegexMatch, toMatch))
	}

	if field.GlobMatch != nil {
		matches = append(matches, GlobMatch(field.GlobMatch, toMatch))
	}

	for _, match := range matches {
		if !match {
			return false
		}

	}

	return true
}

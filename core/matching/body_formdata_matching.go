package matching

import (
	"github.com/SpectoLabs/hoverfly/core/models"
)

func BodyMatching(fields []models.RequestFieldMatchers, req models.RequestDetails) *FieldMatch {

	matched := true
	hasForm := false
	var score int

	if len(fields) == 0 {
		return &FieldMatch{
			Matched: matched,
			Score:   1,
		}
	}

	for _, field := range fields {
		if field.Matcher == "form" {
			hasForm = true
			formMatchers := field.Value.(map[string][]models.RequestFieldMatchers)
			formMatched := processFormMatcher(formMatchers, req.FormData)
			if !formMatched.Matched {
				matched = false
			}
			score += formMatched.Score
		}
	}
	if !hasForm {
		bodyMatched := FieldMatcher(fields, req.Body)
		if !bodyMatched.Matched {
			matched = false
		}
		score += bodyMatched.Score
	}

	return &FieldMatch{
		Matched: matched,
		Score:   score,
	}
}

func processFormMatcher(formFields map[string][]models.RequestFieldMatchers, formData map[string][]string) *FieldMatch {
	matched := true
	var score int

	for formField, formMatchers := range formFields {
		formValue, ok := formData[formField]
		if !ok {
			matched = false
			continue
		}
		formMatched := FieldMatcher(formMatchers, formValue[0])
		if !formMatched.Matched {
			matched = false
		}
		score += formMatched.Score

	}
	return &FieldMatch{
		Matched: matched,
		Score:   score,
	}
}

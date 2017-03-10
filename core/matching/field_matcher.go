package matching

import (
	"github.com/SpectoLabs/hoverfly/core/models"
	glob "github.com/ryanuber/go-glob"
)

func FieldMatcher(field *models.RequestFieldMatchers, toMatch string) bool {
	return !(field != nil && field.ExactMatch != nil && !glob.Glob(*field.ExactMatch, toMatch))
}

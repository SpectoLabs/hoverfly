package matching

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/SpectoLabs/hoverfly/core/util"
	. "github.com/onsi/gomega"
)

func Test_FieldMatcher_MatchesTrueWithExactMatch(t *testing.T) {
	RegisterTestingT(t)

	Expect(FieldMatcher(&models.RequestFieldMatchers{
		ExactMatch: util.StringToPointer("yes"),
	}, "yes")).To(BeTrue())
}

func Test_FieldMatcher_MatchesFalseWithIncorrectExactMatch(t *testing.T) {
	RegisterTestingT(t)

	Expect(FieldMatcher(&models.RequestFieldMatchers{
		ExactMatch: util.StringToPointer("yes"),
	}, "no")).To(BeFalse())
}

func Test_FieldMatcher_MatchesTrueWithNilMatchers(t *testing.T) {
	RegisterTestingT(t)

	Expect(FieldMatcher(nil, "no")).To(BeTrue())
}

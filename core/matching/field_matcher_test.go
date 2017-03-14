package matching_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/matching"
	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/SpectoLabs/hoverfly/core/util"
	. "github.com/onsi/gomega"
)

func Test_FieldMatcher_MatchesTrueWithNilMatchers(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.FieldMatcher(nil, "no")).To(BeTrue())
}

func Test_FieldMatcher_MatchesTrueWithACombinationOfMatchers(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.FieldMatcher(&models.RequestFieldMatchers{
		ExactMatch: util.StringToPointer("testtesttest"),
		RegexMatch: util.StringToPointer("test"),
	}, `testtesttest`)).To(BeTrue())
}

func Test_FieldMatcher_MatchesFalseWithACombinationOfMatchers(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.FieldMatcher(&models.RequestFieldMatchers{
		ExactMatch: util.StringToPointer("testtesttest"),
		RegexMatch: util.StringToPointer("tst"),
	}, `testtesttest`)).To(BeFalse())
}

func Test_FieldMatcher_MatchesFalseWithADifferentCombinationOfMatchers(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.FieldMatcher(&models.RequestFieldMatchers{
		GlobMatch:  util.StringToPointer("*test"),
		RegexMatch: util.StringToPointer("tst"),
	}, `testtesttest`)).To(BeFalse())
}

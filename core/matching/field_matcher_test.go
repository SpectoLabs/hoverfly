package matching_test

import (
	"encoding/xml"
	"testing"

	"github.com/SpectoLabs/hoverfly/core/matching"
	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/SpectoLabs/hoverfly/core/util"
	. "github.com/onsi/gomega"
)

func Test_FieldMatcher_MatchesTrue_WithNilMatchers(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.FieldMatcher(nil, "test")).To(BeTrue())
}

func Test_FieldMatcher_MatchesTrueWithJsonMatch(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.FieldMatcher(&models.RequestFieldMatchers{
		JsonMatch: util.StringToPointer(`{"test":true}`),
	}, `{"test": true}`)).To(BeTrue())
}

func Test_FieldMatcher_MatchesFalseWithJsonMatch(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.FieldMatcher(&models.RequestFieldMatchers{
		JsonMatch: util.StringToPointer(`{"test":true}`),
	}, `{"test": [ ] }`)).To(BeFalse())
}

func Test_FieldMatcher_MatchesTrue_WithMatchersNotDefined(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.FieldMatcher(&models.RequestFieldMatchers{}, "test")).To(BeTrue())
}

func Test_FieldMatcher_WithMultipleMatchers_MatchesTrue(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.FieldMatcher(&models.RequestFieldMatchers{
		ExactMatch: util.StringToPointer("testtesttest"),
		RegexMatch: util.StringToPointer("test"),
	}, `testtesttest`)).To(BeTrue())
}

func Test_FieldMatcher_WithMultipleMatchers_AlsoMatchesTrue(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.FieldMatcher(&models.RequestFieldMatchers{
		XpathMatch: util.StringToPointer("/list/item[1]/field"),
		RegexMatch: util.StringToPointer("test"),
	}, xml.Header+"<list><item><field>test</field></item></list>")).To(BeTrue())
}

func Test_FieldMatcher_WithMultipleMatchers_MatchesFalse(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.FieldMatcher(&models.RequestFieldMatchers{
		ExactMatch: util.StringToPointer("testtesttest"),
		RegexMatch: util.StringToPointer("tst"),
	}, `testtesttest`)).To(BeFalse())
}

func Test_FieldMatcher__WithMultipleMatchers_AlsoMatchesFalse(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.FieldMatcher(&models.RequestFieldMatchers{
		GlobMatch:     util.StringToPointer("*test"),
		JsonPathMatch: util.StringToPointer("$.test[1]"),
	}, `testtesttest`)).To(BeFalse())
}

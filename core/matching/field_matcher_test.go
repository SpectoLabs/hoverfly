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

	Expect(matching.FieldMatcher(nil, "test").Matched).To(BeTrue())
}

func Test_FieldMatcher_MatchesTrueWithJsonMatch(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.FieldMatcher(&models.RequestFieldMatchers{
		JsonMatch: util.StringToPointer(`{"test":true}`),
	}, `{"test": true}`).Matched).To(BeTrue())
}

func Test_FieldMatcher_MatchesFalseWithJsonMatch(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.FieldMatcher(&models.RequestFieldMatchers{
		JsonMatch: util.StringToPointer(`{"test":true}`),
	}, `{"test": [ ] }`).Matched).To(BeFalse())
}

func Test_FieldMatcher_MatchesTrueWithXmlMatch(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.FieldMatcher(&models.RequestFieldMatchers{
		XmlMatch: util.StringToPointer(`<document></document>`),
	}, `<document></document>`).Matched).To(BeTrue())
}

func Test_FieldMatcher_MatchesFalseWithXmlMatch(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.FieldMatcher(&models.RequestFieldMatchers{
		XmlMatch: util.StringToPointer(`<document></document>`),
	}, `<document>
		<test>data</test>
	</document>`).Matched).To(BeFalse())
}

func Test_FieldMatcher_MatchesTrue_WithMatchersNotDefined(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.FieldMatcher(&models.RequestFieldMatchers{}, "test").Matched).To(BeTrue())
}

func Test_FieldMatcher_WithMultipleMatchers_MatchesTrue(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.FieldMatcher(&models.RequestFieldMatchers{
		ExactMatch: util.StringToPointer("testtesttest"),
		RegexMatch: util.StringToPointer("test"),
	}, `testtesttest`).Matched).To(BeTrue())
}

func Test_FieldMatcher_WithMultipleMatchers_AlsoMatchesTrue(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.FieldMatcher(&models.RequestFieldMatchers{
		XpathMatch: util.StringToPointer("/list/item[1]/field"),
		RegexMatch: util.StringToPointer("test"),
	}, xml.Header+"<list><item><field>test</field></item></list>").Matched).To(BeTrue())
}

func Test_FieldMatcher_WithMultipleMatchers_MatchesFalse(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.FieldMatcher(&models.RequestFieldMatchers{
		ExactMatch: util.StringToPointer("testtesttest"),
		RegexMatch: util.StringToPointer("tst"),
	}, `testtesttest`).Matched).To(BeFalse())
}

func Test_FieldMatcher__WithMultipleMatchers_AlsoMatchesFalse(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.FieldMatcher(&models.RequestFieldMatchers{
		GlobMatch:     util.StringToPointer("*test"),
		JsonPathMatch: util.StringToPointer("$.test[1]"),
	}, `testtesttest`).Matched).To(BeFalse())
}

//func Test_FieldMatcher_CountsMatchesOnMatch(t *testing.T)  {
//	RegisterTestingT(t)
//
//	_, count := matching.FieldMatcher(&models.RequestFieldMatchers{
//		GlobMatch:  util.StringToPointer("*test"),
//		ExactMatch: util.StringToPointer("testtesttest"),
//		RegexMatch: util.StringToPointer(".*"),
//	}, `testtesttest`)
//
//	Expect(count).To(Equal(3))
//}

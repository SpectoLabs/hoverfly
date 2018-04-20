package matching_test

import (
	"encoding/xml"
	"testing"

	"github.com/SpectoLabs/hoverfly/core/matching"
	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/SpectoLabs/hoverfly/core/util"
	. "github.com/onsi/gomega"
)

func Test_CountlessFieldMatcher_MatchesTrue_WithNilMatchers(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.UnscoredFieldMatcher(nil, "test").Matched).To(BeTrue())
}

func Test_CountlessFieldMatcher_MatchesTrueWithJsonMatch(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.UnscoredFieldMatcher(&models.RequestFieldMatchers{
		Matcher:   "json",
		Value:     `{"test":true}`,
		JsonMatch: util.StringToPointer(`{"test":true}`),
	}, `{"test": true}`).Matched).To(BeTrue())
}

func Test_CountlessFieldMatcher_MatchesFalseWithJsonMatch(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.UnscoredFieldMatcher(&models.RequestFieldMatchers{
		Matcher:   "json",
		Value:     `{"test":true}`,
		JsonMatch: util.StringToPointer(`{"test":true}`),
	}, `{"test": [ ] }`).Matched).To(BeFalse())
}

func Test_CountlessFieldMatcher_MatchesTrueWithXmlMatch(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.UnscoredFieldMatcher(&models.RequestFieldMatchers{
		XmlMatch: util.StringToPointer(`<document></document>`),
	}, `<document></document>`).Matched).To(BeTrue())
}

func Test_CountlessFieldMatcher_MatchesFalseWithXmlMatch(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.UnscoredFieldMatcher(&models.RequestFieldMatchers{
		Matcher:  "xml",
		Value:    "`<document></document>`",
		XmlMatch: util.StringToPointer(`<document></document>`),
	}, `<document>
		<test>data</test>
	</document>`).Matched).To(BeFalse())
}

func Test_CountlessFieldMatcher_MatchesTrue_WithMatchersNotDefined(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.UnscoredFieldMatcher(&models.RequestFieldMatchers{}, "test").Matched).To(BeTrue())
}

func Test_CountlessFieldMatcher_WithMultipleMatchers_MatchesTrue(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.UnscoredFieldMatcher(&models.RequestFieldMatchers{
		ExactMatch: util.StringToPointer("testtesttest"),
		RegexMatch: util.StringToPointer("test"),
	}, `testtesttest`).Matched).To(BeTrue())
}

func Test_CountlessFieldMatcher_WithMultipleMatchers_AlsoMatchesTrue(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.UnscoredFieldMatcher(&models.RequestFieldMatchers{
		XpathMatch: util.StringToPointer("/list/item[1]/field"),
		RegexMatch: util.StringToPointer("test"),
	}, xml.Header+"<list><item><field>test</field></item></list>").Matched).To(BeTrue())
}

func Test_ScoredFieldMatcher_MatchesTrue_WithNilMatchers(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.ScoredFieldMatcher(nil, "test").Matched).To(BeTrue())
}

func Test_ScoredFieldMatcher_MatchesTrueWithJsonMatch(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.ScoredFieldMatcher(&models.RequestFieldMatchers{
		Matcher:   "json",
		Value:     `{"test":true}`,
		JsonMatch: util.StringToPointer(`{"test":true}`),
	}, `{"test": true}`).Matched).To(BeTrue())
}

func Test_ScoredFieldMatcher_MatchesFalseWithJsonMatch(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.ScoredFieldMatcher(&models.RequestFieldMatchers{
		Matcher:   "json",
		Value:     `{"test":true}`,
		JsonMatch: util.StringToPointer(`{"test":true}`),
	}, `{"test": [ ] }`).Matched).To(BeFalse())
}

func Test_ScoredFieldMatcher_MatchesTrueWithXmlMatch(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.ScoredFieldMatcher(&models.RequestFieldMatchers{
		XmlMatch: util.StringToPointer(`<document></document>`),
	}, `<document></document>`).Matched).To(BeTrue())
}

func Test_ScoredFieldMatcher_MatchesFalseWithXmlMatch(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.ScoredFieldMatcher(&models.RequestFieldMatchers{
		Matcher:  "xml",
		Value:    "<document></document>",
		XmlMatch: util.StringToPointer(`<document></document>`),
	}, `<document>
		<test>data</test>
	</document>`).Matched).To(BeFalse())
}

func Test_ScoredFieldMatcher_MatchesTrue_WithMatchersNotDefined(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.ScoredFieldMatcher(&models.RequestFieldMatchers{}, "test").Matched).To(BeTrue())
}

func Test_ScoredFieldMatcher_WithMultipleMatchers_MatchesTrue(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.ScoredFieldMatcher(&models.RequestFieldMatchers{
		ExactMatch: util.StringToPointer("testtesttest"),
		RegexMatch: util.StringToPointer("test"),
	}, `testtesttest`).Matched).To(BeTrue())
}

func Test_ScoredFieldMatcher_WithExactMatch_ScoresDouble(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.ScoredFieldMatcher(&models.RequestFieldMatchers{
		Matcher:    "exact",
		Value:      "testtesttest",
		ExactMatch: util.StringToPointer("testtesttest"),
	}, `testtesttest`).MatchScore).To(Equal(2))
}

func Test_ScoredFieldMatcher_WithMultipleMatchers_AlsoMatchesTrue(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.ScoredFieldMatcher(&models.RequestFieldMatchers{
		XpathMatch: util.StringToPointer("/list/item[1]/field"),
		RegexMatch: util.StringToPointer("test"),
	}, xml.Header+"<list><item><field>test</field></item></list>").Matched).To(BeTrue())
}

func Test_ScoredFieldMatcher_CountZero_WhenFieldIsNil(t *testing.T) {
	RegisterTestingT(t)

	// Glob, regex, and exact
	matcher := matching.ScoredFieldMatcher(nil, `testtesttest`)

	Expect(matcher.Matched).To(BeTrue())
	Expect(matcher.MatchScore).To(Equal(0))
}

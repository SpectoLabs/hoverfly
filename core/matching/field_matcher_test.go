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
		JsonMatch: util.StringToPointer(`{"test":true}`),
	}, `{"test": true}`).Matched).To(BeTrue())
}

func Test_CountlessFieldMatcher_MatchesFalseWithJsonMatch(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.UnscoredFieldMatcher(&models.RequestFieldMatchers{
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

func Test_CountlessFieldMatcher_WithMultipleMatchers_MatchesFalse(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.UnscoredFieldMatcher(&models.RequestFieldMatchers{
		ExactMatch: util.StringToPointer("testtesttest"),
		RegexMatch: util.StringToPointer("tst"),
	}, `testtesttest`).Matched).To(BeFalse())
}

func Test_CountlessFieldMatcher__WithMultipleMatchers_AlsoMatchesFalse(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.UnscoredFieldMatcher(&models.RequestFieldMatchers{
		GlobMatch:     util.StringToPointer("*test"),
		JsonPathMatch: util.StringToPointer("$.test[1]"),
	}, `testtesttest`).Matched).To(BeFalse())
}

func Test_CountingFieldMatcher_MatchesTrue_WithNilMatchers(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.ScoredFieldMatcher(nil, "test").Matched).To(BeTrue())
}

func Test_CountingFieldMatcher_MatchesTrueWithJsonMatch(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.ScoredFieldMatcher(&models.RequestFieldMatchers{
		JsonMatch: util.StringToPointer(`{"test":true}`),
	}, `{"test": true}`).Matched).To(BeTrue())
}

func Test_CountingFieldMatcher_MatchesFalseWithJsonMatch(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.ScoredFieldMatcher(&models.RequestFieldMatchers{
		JsonMatch: util.StringToPointer(`{"test":true}`),
	}, `{"test": [ ] }`).Matched).To(BeFalse())
}

func Test_CountingFieldMatcher_MatchesTrueWithXmlMatch(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.ScoredFieldMatcher(&models.RequestFieldMatchers{
		XmlMatch: util.StringToPointer(`<document></document>`),
	}, `<document></document>`).Matched).To(BeTrue())
}

func Test_CountingFieldMatcher_MatchesFalseWithXmlMatch(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.ScoredFieldMatcher(&models.RequestFieldMatchers{
		XmlMatch: util.StringToPointer(`<document></document>`),
	}, `<document>
		<test>data</test>
	</document>`).Matched).To(BeFalse())
}

func Test_CountingFieldMatcher_MatchesTrue_WithMatchersNotDefined(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.ScoredFieldMatcher(&models.RequestFieldMatchers{}, "test").Matched).To(BeTrue())
}

func Test_CountingFieldMatcher_WithMultipleMatchers_MatchesTrue(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.ScoredFieldMatcher(&models.RequestFieldMatchers{
		ExactMatch: util.StringToPointer("testtesttest"),
		RegexMatch: util.StringToPointer("test"),
	}, `testtesttest`).Matched).To(BeTrue())
}

func Test_CountingFieldMatcher_WithMultipleMatchers_AlsoMatchesTrue(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.ScoredFieldMatcher(&models.RequestFieldMatchers{
		XpathMatch: util.StringToPointer("/list/item[1]/field"),
		RegexMatch: util.StringToPointer("test"),
	}, xml.Header+"<list><item><field>test</field></item></list>").Matched).To(BeTrue())
}

func Test_CountingFieldMatcher_WithMultipleMatchers_MatchesFalse(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.ScoredFieldMatcher(&models.RequestFieldMatchers{
		ExactMatch: util.StringToPointer("testtesttest"),
		RegexMatch: util.StringToPointer("tst"),
	}, `testtesttest`).Matched).To(BeFalse())
}

func Test_CountingFieldMatcher__WithMultipleMatchers_AlsoMatchesFalse(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.ScoredFieldMatcher(&models.RequestFieldMatchers{
		GlobMatch:     util.StringToPointer("*test"),
		JsonPathMatch: util.StringToPointer("$.test[1]"),
	}, `testtesttest`).Matched).To(BeFalse())
}

func Test_CountingFieldMatcher_CountsMatches_WhenThereIsAMatch(t *testing.T) {
	RegisterTestingT(t)

	// Glob, regex, and exact
	matcher := matching.ScoredFieldMatcher(&models.RequestFieldMatchers{
		GlobMatch:  util.StringToPointer("*test"),
		RegexMatch: util.StringToPointer(".*"),
		ExactMatch: util.StringToPointer("testtesttest"),
	}, `testtesttest`)

	Expect(matcher.Matched).To(BeTrue())
	Expect(matcher.MatchScore).To(Equal(3))

	// JSON and JSONPath
	matcher = matching.ScoredFieldMatcher(&models.RequestFieldMatchers{
		JsonMatch:     util.StringToPointer(`{"test":true}`),
		JsonPathMatch: util.StringToPointer(`$.test`),
	}, `{"test":true}`)

	Expect(matcher.Matched).To(BeTrue())
	Expect(matcher.MatchScore).To(Equal(2))

	// XML and XMLPath
	matcher = matching.ScoredFieldMatcher(&models.RequestFieldMatchers{
		XmlMatch:   util.StringToPointer(xml.Header + "<list><item><field>test</field></item></list>"),
		XpathMatch: util.StringToPointer(`/list/item[1]/field`),
	}, xml.Header+"<list><item><field>test</field></item></list>")

	Expect(matcher.Matched).To(BeTrue())
	Expect(matcher.MatchScore).To(Equal(2))
}

func Test_CountingFieldMatcher_CountsMatches_WhenThereIsNoMatch(t *testing.T) {
	RegisterTestingT(t)

	// Glob, regex, and exact
	matcher := matching.ScoredFieldMatcher(&models.RequestFieldMatchers{
		GlobMatch:     util.StringToPointer("*test"),
		RegexMatch:    util.StringToPointer(".*"),
		ExactMatch:    util.StringToPointer("testtesttest"),
		JsonPathMatch: util.StringToPointer(`$.notmatch`),
	}, `testtesttest`)

	Expect(matcher.Matched).To(BeFalse())
	Expect(matcher.MatchScore).To(Equal(3))
}

func Test_CountingFieldMatcher_CountZero_WhenFieldIsNil(t *testing.T) {
	RegisterTestingT(t)

	// Glob, regex, and exact
	matcher := matching.ScoredFieldMatcher(nil, `testtesttest`)

	Expect(matcher.Matched).To(BeTrue())
	Expect(matcher.MatchScore).To(Equal(0))
}

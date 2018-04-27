package matching_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/matching"
	"github.com/SpectoLabs/hoverfly/core/matching/matchers"
	"github.com/SpectoLabs/hoverfly/core/models"
	. "github.com/onsi/gomega"
)

func Test_FieldMatcher_MatchesTrue_WithNilMatchers(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.FieldMatcher(nil, "test").Matched).To(BeTrue())
}

func Test_FieldMatcher_MatchesTrueWithJsonMatch(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.FieldMatcher([]models.RequestFieldMatchers{
		{
			Matcher: matchers.Json,
			Value:   `{"test":true}`,
		},
	}, `{"test": true}`).Matched).To(BeTrue())
}

func Test_FieldMatcher_MatchesFalseWithJsonMatch(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.FieldMatcher([]models.RequestFieldMatchers{
		{
			Matcher: matchers.Json,
			Value:   `{"test":true}`,
		},
	}, `{"test": [ ] }`).Matched).To(BeFalse())
}

func Test_FieldMatcher_MatchesTrueWithXmlMatch(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.FieldMatcher([]models.RequestFieldMatchers{
		{
			Matcher: matchers.Xml,
			Value:   `<document></document>`,
		},
	}, `<document></document>`).Matched).To(BeTrue())
}

func Test_FieldMatcher_MatchesFalseWithXmlMatch(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.FieldMatcher([]models.RequestFieldMatchers{
		{
			Matcher: matchers.Xml,
			Value:   "<document></document>",
		},
	}, `<document>
		<test>data</test>
	</document>`).Matched).To(BeFalse())
}

func Test_FieldMatcher_MatchesTrue_WithMatchersNotDefined(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.FieldMatcher([]models.RequestFieldMatchers{}, "test").Matched).To(BeTrue())
}

func Test_FieldMatcher_WithExactMatch_ScoresDouble(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.FieldMatcher([]models.RequestFieldMatchers{
		{
			Matcher: matchers.Exact,
			Value:   "testtesttest",
		},
	}, `testtesttest`).MatchScore).To(Equal(2))
}

func Test_FieldMatcher_CountZero_WhenFieldIsNil(t *testing.T) {
	RegisterTestingT(t)

	// Glob, regex, and exact
	matcher := matching.FieldMatcher(nil, `testtesttest`)

	Expect(matcher.Matched).To(BeTrue())
	Expect(matcher.MatchScore).To(Equal(0))
}

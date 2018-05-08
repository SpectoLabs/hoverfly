package matching_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/matching"
	"github.com/SpectoLabs/hoverfly/core/matching/matchers"
	"github.com/SpectoLabs/hoverfly/core/models"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

type fieldMatcherTest struct {
	name        string
	matchers    []models.RequestFieldMatchers
	toMatch     string
	equals      types.GomegaMatcher
	scoreEquals types.GomegaMatcher
}

var fieldMatcherTests = []fieldMatcherTest{
	{
		name:        "MatchesTrue_WithNilMatchers",
		matchers:    nil,
		toMatch:     "test",
		equals:      BeTrue(),
		scoreEquals: Equal(0),
	},
	{
		name: "MatchesTrueWithDefaultMatcherWhichIsExactMatch",
		matchers: []models.RequestFieldMatchers{
			{
				Value: `test`,
			},
		},
		toMatch: "test",
		equals:  BeTrue(),
	},
	{
		name: "MatchesTrueWithJsonMatch",
		matchers: []models.RequestFieldMatchers{
			{
				Matcher: matchers.Json,
				Value:   `{"test":true}`,
			},
		},
		toMatch: `{"test":true}`,
		equals:  BeTrue(),
	},
	{
		name: "MatchesFalseWithJsonMatch",
		matchers: []models.RequestFieldMatchers{
			{
				Matcher: matchers.Json,
				Value:   `{"test":true}`,
			},
		},
		toMatch: "test",
		equals:  BeFalse(),
	},
	{
		name: "MatchesTrueWithXmlMatch",
		matchers: []models.RequestFieldMatchers{
			{
				Matcher: matchers.Xml,
				Value:   `<document></document>`,
			},
		},
		toMatch: `<document></document>`,
		equals:  BeTrue(),
	},
	{
		name: "MatchesFalseWithXmlMatch",
		matchers: []models.RequestFieldMatchers{
			{
				Matcher: matchers.Xml,
				Value:   "<document></document>",
			},
		},
		toMatch: `<document>
		<test>data</test>
	</document>`,
		equals: BeFalse(),
	},
	{
		name:     "MatchesTrue_WithMatchersNotDefined",
		matchers: []models.RequestFieldMatchers{},
		toMatch:  "test",
		equals:   BeTrue(),
	},
	{
		name: "WithExactMatch_ScoresDouble(",
		matchers: []models.RequestFieldMatchers{
			{
				Matcher: matchers.Exact,
				Value:   "test",
			},
		},
		toMatch:     "test",
		scoreEquals: Equal(2),
	},
	{
		name: "WithMultipleMatchers_MatchesOnBoth",
		matchers: []models.RequestFieldMatchers{
			{
				Matcher: matchers.Exact,
				Value:   "test",
			},
			{
				Matcher: matchers.Exact,
				Value:   "test",
			},
		},
		toMatch: "test",
		equals:  BeTrue(),
	},
	{
		name: "WithMultipleMatchers_MatchesOne",
		matchers: []models.RequestFieldMatchers{
			{
				Matcher: matchers.Exact,
				Value:   "test",
			},
			{
				Matcher: matchers.Exact,
				Value:   "nottest",
			},
		},
		toMatch: "test",
		equals:  BeFalse(),
	},
	{
		name: "FieldMatcher_WithMultipleMatchers_ScoresDouble",
		matchers: []models.RequestFieldMatchers{
			{
				Matcher: matchers.Exact,
				Value:   "test",
			},
			{
				Matcher: matchers.Exact,
				Value:   "test",
			},
		},
		toMatch:     "test",
		scoreEquals: Equal(4),
	},
}

func Test_FieldMatcher(t *testing.T) {
	RegisterTestingT(t)

	for _, test := range fieldMatcherTests {
		result := matching.FieldMatcher(test.matchers, test.toMatch)
		if test.equals != nil {
			Expect(result.Matched).To(test.equals, test.name)
		}
		if test.scoreEquals != nil {
			Expect(result.Score).To(test.scoreEquals, test.name)
		}
	}

}

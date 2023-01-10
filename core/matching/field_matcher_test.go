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
	{
		name: "MatcherNameShouldBeCaseInsensitive",
		matchers: []models.RequestFieldMatchers{
			{
				Matcher: "XML",
				Value:   `<document></document>`,
			},
		},
		toMatch: `<document></document>`,
		equals:  BeTrue(),
	},
	{
		name: "MatcherChaining1",
		matchers: []models.RequestFieldMatchers{
			{
				Matcher: "xpath",
				Value:   "/document/id",
				DoMatch: &models.RequestFieldMatchers{
					Matcher: "exact",
					Value:   "12345",
				},
			},
		},
		toMatch: "<document><id>12345</id><name>Test</name></document>",
		equals:  BeTrue(),
	},
	{
		name: "MatcherChaining3",
		matchers: []models.RequestFieldMatchers{
			{
				Matcher: "jsonpath",
				Value:   "$.testArr",
				DoMatch: &models.RequestFieldMatchers{
					Matcher: "array",
					Value:   []string{"q1", "q2", "q3"},
					Config: map[string]interface{}{
						matchers.IGNORE_OCCURRENCES: false,
						matchers.IGNORE_ORDER:       false,
						matchers.IGNORE_UNKNOWN:     false,
					},
				},
			},
		},
		toMatch: `{"testArr":["q1", "q2", "q3"]}`,
		equals:  BeTrue(),
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

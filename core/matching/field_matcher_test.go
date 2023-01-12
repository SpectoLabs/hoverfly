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
		name: "MatcherChaining2",
		matchers: []models.RequestFieldMatchers{
			{
				Matcher: "xpath",
				Value:   "/document/details",
				DoMatch: &models.RequestFieldMatchers{
					Matcher: "jsonpath",
					Value:   "$.name",
					DoMatch: &models.RequestFieldMatchers{
						Matcher: "glob",
						Value:   "*es*",
						DoMatch: &models.RequestFieldMatchers{
							Matcher: "exact",
							Value:   "Test",
						},
					},
				},
			},
		},
		toMatch: `<document><details>{"name":"Test", "id":"12345"}</details></document>`,
		equals:  BeTrue(),
	},
	{
		name: "TestJwtMatcher",
		matchers: []models.RequestFieldMatchers{
			{
				Matcher: "jwt",
				Value:   `{"header":{"alg":"HS256"},"payload":{"sub":"1234567890","name":"John Doe"}}`,
			},
		},
		toMatch: `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c`,
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

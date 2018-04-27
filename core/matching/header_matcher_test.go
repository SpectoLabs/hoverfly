package matching_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/matching"
	"github.com/SpectoLabs/hoverfly/core/matching/matchers"
	"github.com/SpectoLabs/hoverfly/core/models"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

type headerMatchingTest struct {
	name                string
	headers             map[string][]string
	headersWithMatchers map[string][]models.RequestFieldMatchers
	toMatchHeaders      map[string][]string
	equals              types.GomegaMatcher
	matchEquals         types.GomegaMatcher
}

var tests = []headerMatchingTest{
	headerMatchingTest{
		name: "basic",
		headers: map[string][]string{
			"header1": {"val1"},
			"header2": {"val2"},
		},
		toMatchHeaders: map[string][]string{
			"header1": {"val1"},
			"header2": {"val2"},
		},
		equals: BeTrue(),
	},
	{
		name: "IgnoreTestCaseInsensitive",
		headers: map[string][]string{
			"header1": {"val1"},
			"header2": {"val2"},
		},
		toMatchHeaders: map[string][]string{
			"HEADER1": {"val1"},
			"Header2": {"VAL2"},
		},
		equals: BeTrue(),
	},
	{
		name: "MatchingHeadersHasMoreKeysThanRequestMatchesFalse",
		headers: map[string][]string{
			"header1": {"val1"},
			"header2": {"val2"},
		},
		toMatchHeaders: map[string][]string{
			"header1": {"val1"},
		},
		equals: BeFalse(),
	},
	{
		name: "MatchingHeadersHasMoreValuesThanRequestMatchesFalse",
		headers: map[string][]string{
			"header2": {"val1", "val2"},
		},
		toMatchHeaders: map[string][]string{
			"header2": {"val1"},
		},
		equals: BeFalse(),
	},
	{
		name: "MatchingHeadersHasLessKeysThanRequestMatchesTrue",
		headers: map[string][]string{
			"header2": {"val2"},
		},
		toMatchHeaders: map[string][]string{
			"HEADER1": {"val1"},
			"header2": {"val2"},
		},
		equals: BeTrue(),
	},
	{
		name: "RequestHeadersContainsAllOFMatcherHeaders",
		headers: map[string][]string{
			"header2": {"val2"},
		},
		toMatchHeaders: map[string][]string{
			"header2": {"val1", "val2"},
		},
		equals: BeTrue(),
	},
	{
		name: "ShouldNotMatchUnlessAllHeadersValuesAreFound",
		headers: map[string][]string{
			"header2": {"val1", "val2"},
		},
		toMatchHeaders: map[string][]string{
			"header2": {"val1", "nomatch"},
		},
		equals: BeFalse(),
	},
	{
		name: "CountsMatches_WhenThereIsAMatch",
		headers: map[string][]string{
			"header1": {"val1", "val2"},
		},
		toMatchHeaders: map[string][]string{
			"header1":     {"val1", "val2"},
			"extraHeader": {"extraHeader1", "extraHeader2"},
		},
		equals:      BeTrue(),
		matchEquals: Equal(2),
	},
	{
		name: "CountsMatches_WhenThereIsAMatch_2",
		headers: map[string][]string{
			"header1": {"val1", "val2"},
			"header2": {"val3"},
		},
		toMatchHeaders: map[string][]string{
			"header1":     {"val1", "val2"},
			"header2":     {"val3", "extra"},
			"extraHeader": {"extraHeader1", "extraHeader2"},
		},
		equals:      BeTrue(),
		matchEquals: Equal(3),
	},
	{
		name: "basic",
		headers: map[string][]string{
			"header1": {"val1", "val2"},
			"header2": {"val3", "nomatch"},
		},
		toMatchHeaders: map[string][]string{
			"header1":     {"val1", "val2"},
			"header2":     {"val3", "extra"},
			"extraHeader": {"extraHeader1", "extraHeader2"},
		},
		equals:      BeFalse(),
		matchEquals: Equal(3),
	},
	{
		name: "headersWithMatchers 1 header 1 value",
		headersWithMatchers: map[string][]models.RequestFieldMatchers{
			"header1": {
				{
					Matcher: matchers.Exact,
					Value:   "val1",
				},
			},
		},
		toMatchHeaders: map[string][]string{
			"header1": {"val1"},
		},
		equals:      BeTrue(),
		matchEquals: Equal(2),
	},
	{
		name: "headersWithMatchers 1 header 2 values",
		headersWithMatchers: map[string][]models.RequestFieldMatchers{
			"header1": {
				{
					Matcher: matchers.Exact,
					Value:   "val1;val2",
				},
			},
		},
		toMatchHeaders: map[string][]string{
			"header1": {"val1", "val2"},
		},
		equals:      BeTrue(),
		matchEquals: Equal(2),
	},
	{
		name: "headersWithMatchers fail",
		headersWithMatchers: map[string][]models.RequestFieldMatchers{
			"header1": {
				{
					Matcher: matchers.Exact,
					Value:   "val1",
				},
			},
		},
		toMatchHeaders: map[string][]string{
			"header1": {"val2"},
		},
		equals:      BeFalse(),
		matchEquals: Equal(0),
	},
	{
		name: "headersWithMatchers 2 headers",
		headersWithMatchers: map[string][]models.RequestFieldMatchers{
			"header1": {
				{
					Matcher: matchers.Exact,
					Value:   "val1",
				},
			},
			"header2": {
				{
					Matcher: matchers.Glob,
					Value:   "*a*",
				},
			},
		},
		toMatchHeaders: map[string][]string{
			"header1": {"val1"},
			"header2": {"val1"},
		},
		equals:      BeTrue(),
		matchEquals: Equal(3),
	},
	{
		name: "headersWithMatchers 2 headers fail",
		headersWithMatchers: map[string][]models.RequestFieldMatchers{
			"header1": {
				{
					Matcher: matchers.Exact,
					Value:   "val1",
				},
			},
			"header2": {
				{
					Matcher: matchers.Glob,
					Value:   "*a*",
				},
			},
		},
		toMatchHeaders: map[string][]string{
			"header1": {"val1"},
		},
		equals:      BeFalse(),
		matchEquals: Equal(2),
	},
	{
		name: "headersWithMatchers case insensitive",
		headersWithMatchers: map[string][]models.RequestFieldMatchers{
			"HEADER1": {
				{
					Matcher: matchers.Exact,
					Value:   "val1",
				},
			},
		},
		toMatchHeaders: map[string][]string{
			"header1": {"val1"},
		},
		equals:      BeTrue(),
		matchEquals: Equal(2),
	},
	{
		name: "headersWithMatchers headers case insensitive",
		headersWithMatchers: map[string][]models.RequestFieldMatchers{
			"HEADER1": {
				{
					Matcher: matchers.Exact,
					Value:   "val1",
				},
			},
		},
		toMatchHeaders: map[string][]string{
			"Header1": {"val1"},
		},
		equals:      BeTrue(),
		matchEquals: Equal(2),
	},
	{
		name: "headersWithMatchers case insensitive fail",
		headersWithMatchers: map[string][]models.RequestFieldMatchers{
			"HEADER1": {
				{
					Matcher: matchers.Exact,
					Value:   "val1",
				},
			},
		},
		toMatchHeaders: map[string][]string{
			"soemthing-else": {"val1"},
		},
		equals:      BeFalse(),
		matchEquals: Equal(0),
	},
	{
		name: "headersWithMatchers defaults to original headers",
		headersWithMatchers: map[string][]models.RequestFieldMatchers{
			"soemthing-else": {
				{
					Matcher: matchers.Exact,
					Value:   "val1",
				},
			},
		},
		headers: map[string][]string{
			"soemthing-else": {"val1"},
		},
		toMatchHeaders: map[string][]string{
			"soemthing-else": {"val1"},
		},
		equals:      BeTrue(),
		matchEquals: Equal(3),
	},
	{
		name: "headersWithMatchers defaults to original headers fail",
		headersWithMatchers: map[string][]models.RequestFieldMatchers{
			"soemthing-else": {
				{
					Matcher: matchers.Exact,
					Value:   "val1",
				},
			},
		},
		headers: map[string][]string{
			"HEADER1": {"val1"},
		},
		toMatchHeaders: map[string][]string{
			"soemthing-else": {"val1"},
		},
		equals:      BeFalse(),
		matchEquals: Equal(2),
	},
}

func Test_HeaderMatching(t *testing.T) {
	RegisterTestingT(t)

	for _, test := range tests {
		result := matching.HeaderMatching(models.RequestMatcher{
			Headers:             test.headers,
			HeadersWithMatchers: test.headersWithMatchers,
		}, test.toMatchHeaders)

		Expect(result.Matched).To(test.equals, test.name)
		if test.matchEquals != nil {
			Expect(result.MatchScore).To(test.matchEquals, test.name)
		}
	}

}

package matching_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/matching"
	"github.com/SpectoLabs/hoverfly/core/matching/matchers"
	"github.com/SpectoLabs/hoverfly/core/models"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

type queryMatchingTest struct {
	name                string
	queriesWithMatchers map[string][]models.RequestFieldMatchers
	toMatchQueries      map[string][]string
	equals              types.GomegaMatcher
	matchEquals         types.GomegaMatcher
}

var queryMatchingTests = []queryMatchingTest{
	{
		name: "basic",
		queriesWithMatchers: map[string][]models.RequestFieldMatchers{
			"query1": {
				{
					Matcher: matchers.Exact,
					Value:   "val1",
				},
			},
		},
		toMatchQueries: map[string][]string{
			"query1": {"val1"},
		},
		equals: BeTrue(),
	},
	{
		name: "basic fail",
		queriesWithMatchers: map[string][]models.RequestFieldMatchers{
			"query1": {
				{
					Matcher: matchers.Exact,
					Value:   "val1",
				},
			},
		},
		toMatchQueries: map[string][]string{
			"query1": {"val2"},
		},
		equals: BeFalse(),
	},
	{
		name: "2 query parameters",
		queriesWithMatchers: map[string][]models.RequestFieldMatchers{
			"query1": {
				{
					Matcher: matchers.Exact,
					Value:   "val1",
				},
			},
			"query2": {
				{
					Matcher: matchers.Glob,
					Value:   "*a*",
				},
			},
		},
		toMatchQueries: map[string][]string{
			"query1": {"val1"},
			"query2": {"val1"},
		},
		equals:      BeTrue(),
		matchEquals: Equal(3),
	},
	{
		name: "2 query parameters fail missing query",
		queriesWithMatchers: map[string][]models.RequestFieldMatchers{
			"query1": {
				{
					Matcher: matchers.Exact,
					Value:   "val1",
				},
			},
			"query2": {
				{
					Matcher: matchers.Glob,
					Value:   "*a*",
				},
			},
		},
		toMatchQueries: map[string][]string{
			"query1": {"val1"},
		},
		equals:      BeFalse(),
		matchEquals: Equal(2),
	},
	{
		name: "2 query parameters fail bad match",
		queriesWithMatchers: map[string][]models.RequestFieldMatchers{
			"query1": {
				{
					Matcher: matchers.Exact,
					Value:   "val1",
				},
			},
			"query2": {
				{
					Matcher: matchers.Glob,
					Value:   "*a*",
				},
			},
		},
		toMatchQueries: map[string][]string{
			"query1": {"val1"},
			"query2": {"vol1"},
		},
		equals:      BeFalse(),
		matchEquals: Equal(2),
	},
	{
		name: "Can handle different cases 1",
		queriesWithMatchers: map[string][]models.RequestFieldMatchers{
			"urlPattern": {
				{
					Matcher: matchers.Glob,
					Value:   "test-(.+).com",
				},
			},
		},
		toMatchQueries: map[string][]string{
			"urlPattern": {"test-(.+).com"},
		},
		equals:      BeTrue(),
		matchEquals: Equal(1),
	},
	{
		name: "Can handle different cases 2",
		queriesWithMatchers: map[string][]models.RequestFieldMatchers{
			"urlPattern": {
				{
					Matcher: matchers.Glob,
					Value:   "test-(.+).com",
				},
			},
		},
		toMatchQueries: map[string][]string{
			"URLPATTERN": {"test-(.+).com"},
		},
		equals:      BeTrue(),
		matchEquals: Equal(1),
	},
}

func Test_QueryMatching(t *testing.T) {
	RegisterTestingT(t)

	for _, test := range queryMatchingTests {
		result := matching.QueryMatching(models.RequestMatcher{
			QueriesWithMatchers: test.queriesWithMatchers,
		}, test.toMatchQueries)

		Expect(result.Matched).To(test.equals, test.name)
		if test.matchEquals != nil {
			Expect(result.Score).To(test.matchEquals, test.name)
		}
	}

}

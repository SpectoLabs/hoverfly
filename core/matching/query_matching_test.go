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
	queriesWithMatchers *models.QueryRequestFieldMatchers
	toMatchQueries      map[string][]string
	equals              types.GomegaMatcher
	matchEquals         types.GomegaMatcher
}

var queryMatchingTests = []queryMatchingTest{
	{
		name:                "nil",
		queriesWithMatchers: nil,
		toMatchQueries: map[string][]string{
			"query1": {"val1"},
		},
		equals: BeTrue(),
	},
	{
		name:                "empty",
		queriesWithMatchers: &models.QueryRequestFieldMatchers{},
		toMatchQueries: map[string][]string{
			"query1": {"val1"},
		},
		equals: BeFalse(),
	},
	{
		name: "basic",
		queriesWithMatchers: &models.QueryRequestFieldMatchers{
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
		queriesWithMatchers: &models.QueryRequestFieldMatchers{
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
		queriesWithMatchers: &models.QueryRequestFieldMatchers{
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
		queriesWithMatchers: &models.QueryRequestFieldMatchers{
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
		queriesWithMatchers: &models.QueryRequestFieldMatchers{
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
		queriesWithMatchers: &models.QueryRequestFieldMatchers{
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
		queriesWithMatchers: &models.QueryRequestFieldMatchers{
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
			Query: test.queriesWithMatchers,
		}, test.toMatchQueries)

		Expect(result.Matched).To(test.equals, test.name)
		if test.matchEquals != nil {
			Expect(result.Score).To(test.matchEquals, test.name)
		}
	}

}

func Test_QueryMatching_ShouldNotModifySourceQueries(t *testing.T) {
	RegisterTestingT(t)

	toMatch := map[string][]string{
		"urlPattern": {"test-(.+).com"},
	}

	result := matching.QueryMatching(models.RequestMatcher{
		Query: &models.QueryRequestFieldMatchers{
			"urlPattern": {
				{
					Matcher: matchers.Glob,
					Value:   "test-(.+).com",
				},
			},
		},
	}, toMatch)

	Expect(result.Matched).To(BeTrue())
	Expect(len(toMatch)).To(Equal(1))
	Expect(toMatch["urlPattern"]).To(Equal([]string{"test-(.+).com"}))
}

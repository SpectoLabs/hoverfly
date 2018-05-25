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
	name           string
	headers        map[string][]models.RequestFieldMatchers
	toMatchHeaders map[string][]string
	equals         types.GomegaMatcher
	matchEquals    types.GomegaMatcher
}

var tests = []headerMatchingTest{
	{
		name: "headersWithMatchers 1 header 1 value",
		headers: map[string][]models.RequestFieldMatchers{
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
		headers: map[string][]models.RequestFieldMatchers{
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
		headers: map[string][]models.RequestFieldMatchers{
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
		headers: map[string][]models.RequestFieldMatchers{
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
		headers: map[string][]models.RequestFieldMatchers{
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
		headers: map[string][]models.RequestFieldMatchers{
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
		headers: map[string][]models.RequestFieldMatchers{
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
		headers: map[string][]models.RequestFieldMatchers{
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
}

func Test_HeaderMatching(t *testing.T) {
	RegisterTestingT(t)

	for _, test := range tests {
		result := matching.HeaderMatching(models.RequestMatcher{
			Headers: test.headers,
		}, test.toMatchHeaders)

		Expect(result.Matched).To(test.equals, test.name)
		if test.matchEquals != nil {
			Expect(result.Score).To(test.matchEquals, test.name)
		}
	}

}

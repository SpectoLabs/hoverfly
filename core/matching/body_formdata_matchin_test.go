package matching_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/matching"
	"github.com/SpectoLabs/hoverfly/core/matching/matchers"
	"github.com/SpectoLabs/hoverfly/core/models"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

type bodyMatchingTest struct {
	name        string
	matchers    []models.RequestFieldMatchers
	toMatch     models.RequestDetails
	equals      types.GomegaMatcher
	matchEquals types.GomegaMatcher
}

var bodyMatchingTests = []bodyMatchingTest{
	{
		name:     "nil",
		matchers: nil,
		toMatch: models.RequestDetails{
			Body: "foo",
		},
		equals: BeTrue(),
	},
	{
		name: "MatchesTrueWithJsonMatch",
		matchers: []models.RequestFieldMatchers{
			{
				Matcher: matchers.Json,
				Value:   `{"test":true}`,
			},
		},
		toMatch: models.RequestDetails{
			Body: `{"test":true}`,
		},
		equals: BeTrue(),
	},
	{
		name: "MatchesFalseWithJsonMatch",
		matchers: []models.RequestFieldMatchers{
			{
				Matcher: matchers.Json,
				Value:   `{"test":true}`,
			},
		},
		toMatch: models.RequestDetails{
			Body: `{"test": ""}`,
		},
		equals: BeFalse(),
	},
	{
		name: "MatchesTrueWithFormMatch",
		matchers: []models.RequestFieldMatchers{
			{
				Matcher: "form",
				Value: map[string][]models.RequestFieldMatchers{
					"name": {
						{
							Matcher: matchers.Exact,
							Value:   "foo",
						},
					},
				},
			},
		},
		toMatch: models.RequestDetails{
			FormData: map[string][]string{"name": {"foo"}},
		},
		equals: BeTrue(),
	},
	{
		name: "MatchesFalseWithFormMatch",
		matchers: []models.RequestFieldMatchers{
			{
				Matcher: "form",
				Value: map[string][]models.RequestFieldMatchers{
					"name": {
						{
							Matcher: matchers.Exact,
							Value:   "foo",
						},
					},
				},
			},
		},
		toMatch: models.RequestDetails{
			FormData: map[string][]string{"name": {"Bar"}},
		},
		equals: BeFalse(),
	},
}

func Test_BodyMatching(t *testing.T) {
	RegisterTestingT(t)

	for _, test := range bodyMatchingTests {
		result := matching.BodyMatching(test.matchers, test.toMatch)

		Expect(result.Matched).To(test.equals, test.name)
		if test.matchEquals != nil {
			Expect(result.Score).To(test.matchEquals, test.name)
		}
	}

}

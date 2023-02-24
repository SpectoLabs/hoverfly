package matching_test

import (
	"fmt"
	"testing"

	"github.com/SpectoLabs/hoverfly/v2/core/matching"
	"github.com/SpectoLabs/hoverfly/v2/core/matching/matchers"
	"github.com/SpectoLabs/hoverfly/v2/core/models"
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
				Value: map[string]interface{}{
					"name": []map[string]interface{}{
						{
							"matcher": matchers.Exact,
							"value":   "foo",
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
				Value: map[string]interface{}{
					"name": []map[string]interface{}{
						{
							"matcher": matchers.Exact,
							"value":   "foo",
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
		fmt.Println("Test matchers is ", test.matchers)
		result := matching.BodyMatching(test.matchers, test.toMatch)

		Expect(result.Matched).To(test.equals, test.name)
		if test.matchEquals != nil {
			Expect(result.Score).To(test.matchEquals, test.name)
		}
	}

}

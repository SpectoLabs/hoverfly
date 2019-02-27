package matchers_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/matching/matchers"
	. "github.com/onsi/gomega"
)


var jsonPathMatchTests = []matchTest{
	{
		name:    "MatchesFalseWithIncorrectDataType",
		match:   1,
		toMatch: "yes",
		matched: false,
		result:  "",
	},
	{
		name:    "MatchesFalseWithInvalidJsonPath",
		match:   "test",
		toMatch: `{"test": "field"}`,
		matched: false,
		result:  "",
	},
	{
		name:    "MatchesTrueWithJsonMatch_GetSingleElement",
		match:   "$.test",
		toMatch: `{"test": "field"}`,
		matched: true,
		result:  "field",
	},
	{
		name:    "MatchesTrueWithJsonMatch_GetSingleArrayElement",
		match:   "$.test",
		toMatch: `{"test": [{"field": 123}]}`,
		matched: true,
		result:  `[{"field":123}]`,
	},
	//{
	//	name:    "MatchesTrueWithJsonMatch_GetSingleElement",
	//	match:   "$.test",
	//	toMatch: `{"test": 123}`,
	//	matched: true,
	//	result:  123,
	//},
	{
		name:    "MatchesFalse_IfJsonElementNotFound",
		match:   "$.test",
		toMatch: `{"not-test": "field"}`,
		matched: false,
		result:  "",
	},
	{
		name:    "MatchesFalseWithIncorrectJsonMatch_GetSingleElement",
		match:   "$.notAField",
		toMatch: `{"test": "field"}`,
		matched: false,
		result:  "",
	},
	{
		name:    "MatchesTrueWithJsonMatch_GetElementFromArray",
		match:   "$.test[1]",
		toMatch: `{"test": [{}, {}]}`,
		matched: true,
		result:  "{}",
	},
	{
		name:    "MatchesFalseWithIncorrectJsonMatch_GetElementFromArray",
		match:   "$.test[2]",
		toMatch: `{"test": [{}, {}]}`,
		matched: false,
		result:  "",
	},
	{
		name:    "MatchesTrueWithJsonMatch_WithExpression",
		match:   "$.test[?(@.field == \"test\")]",
		toMatch: `{"test": [{"field": "test"}]}`,
		matched: true,
		result:  `{"field":"test"}`,
	},
	{
		name:    "MatchesFalseWithIncorrectJsonMatch_WithExpression",
		match:   "$.test[*]?(@.field == \"test\")",
		toMatch: `{"test": [{"field": "not-test"}]}`,
		matched: false,
		result:  "",
	},

	// TODO more jsonpath tests required

	// TODO the following JSONPath expressions are not supported at the moment
	//{
	//	name:      "MatchesTrueWithJsonMatch_WithRootLevelObjectFilter",
	//	match:     "$[?(@.field == \"test\")]",
	//	toMatch:   `{"field": "test"}`,
	//	matched: true,
	//	result:    `{"field":"test"}`,
	//},
	//{
	//	name:      "MatchesTrueWithJsonMatch_WithObjectFilter",
	//	match:     "$.test[?(@.field == \"test\")]",
	//	toMatch:   `{"test": {"field": "test"}}`,
	//	matched: true,
	//	result:    `{"field":"test"}`,
	//},
}

func Test_JsonPathMatch(t *testing.T) {
	RegisterTestingT(t)

	for _, test := range jsonPathMatchTests {
		t.Run(test.name, func(t *testing.T) {

			isMatched, result := matchers.JsonPathMatch(test.match, test.toMatch)

			Expect(isMatched).To(Equal(test.matched))
			Expect(result).To(Equal(test.result))
		})
	}
}
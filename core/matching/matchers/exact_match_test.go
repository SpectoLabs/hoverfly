package matchers_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/matching/matchers"
	. "github.com/onsi/gomega"
)

type matchTest struct {
	name    string
	match   interface{}
	toMatch string
	matched bool
	result  interface{}
}

var exactMatchTests = []matchTest{
	{
		name:    "MatchesFalseWithIncorrectDataType",
		match:   1,
		toMatch: "yes",
		matched: false,
		result:  "",
	},
	{
		name:    "MatchesTrueWithExactMatch",
		match:   "yes",
		toMatch: "yes",
		matched: true,
		result:  "yes",
	},
	{
		name:    "MatchesFalseWithIncorrectExactMatch",
		match:   "yes",
		toMatch: "no",
		matched: false,
		result:  "",
	},
	{
		name:    "MatchesTrueWithJSON",
		match:   `{"test":{"json":true,"minified":true}}`,
		toMatch: `{"test":{"json":true,"minified":true}}`,
		matched: true,
		result:  `{"test":{"json":true,"minified":true}}`,
	},
	{
		name: 		"MatchesFalseWithUnminifiedJSON",
		match: 		`{"test":{"json":true,"minified":true}}`,
		toMatch: 	`{
		"test": {
			"json": true,
			"minified": true
		}
	}`,
		matched: false,
		result:  "",
	},
}

func Test_ExactMatch(t *testing.T) {

	RegisterTestingT(t)

	for _, test := range exactMatchTests {
		t.Run(test.name, func(t *testing.T) {

			isMatched, result := matchers.ExactMatch(test.match, test.toMatch)

			Expect(isMatched).To(Equal(test.matched))
			Expect(result).To(Equal(test.result))
		})
	}
}

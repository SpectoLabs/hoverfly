package matchers_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/matching/matchers"
	. "github.com/onsi/gomega"
)

var regexMatchTests = []matchTest{
	{
		name:    "MatchesFalseWithIncorrectDataType",
		match:   1,
		toMatch: "yes",
		matched: false,
		result:  "",
	},
	{
		name:    "MatchesTrueWithRegexMatch",
		match:   "t[o|a|e]st",
		toMatch: "test",
		matched: true,
		result:  "test",
	},
	{
		name:    "MatchesFalseWithIncorrectRegexMatch",
		match:   "t[o|a]st",
		toMatch: "test",
		matched: false,
		result:  "",
	},
}

func Test_RegexMatch(t *testing.T) {
	RegisterTestingT(t)

	for _, test := range regexMatchTests {
		t.Run(test.name, func(t *testing.T) {

			isMatched, result := matchers.RegexMatch(test.match, test.toMatch)

			Expect(isMatched).To(Equal(test.matched))
			Expect(result).To(Equal(test.result))
		})
	}
}
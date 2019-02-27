package matchers_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/matching/matchers"
	. "github.com/onsi/gomega"
)

var globMatchTests = []matchTest {
	{
		name:    "MatchesFalseWithIncorrectDataType",
		match:   1,
		toMatch: "yes",
		matched: false,
		result:  "",
	},
	{
		name:    "MatchesTrueWithGlobMatch",
		match:   "t*st",
		toMatch: "test",
		matched: true,
		result:  "test",
	},
	{
		name:    "MatchesZeroExtraCharactersAtEnd",
		match:   "test*",
		toMatch: "test",
		matched: true,
		result:  "test",
	},
	{
		name:    "MatchesZeroExtraCharactersAtStart",
		match:   "*test",
		toMatch: "test",
		matched: true,
		result:  "test",
	},
	{
		name:    "MatchesZeroExtraCharactersAtStartAndEnd",
		match:   "*test*",
		toMatch: "test",
		matched: true,
		result:  "test",
	},
	{
		name:    "MatchesUpperCase",
		match:   "*est",
		toMatch: "Test",
		matched: true,
		result:  "Test",
	},
	{
		name:    "MatchesLowerCase",
		match:   "*est",
		toMatch: "test",
		matched: true,
		result:  "test",
	},
	{
		name:    "MatchesAstrik",
		match:   "*est",
		toMatch: "*est",
		matched: true,
		result:  "*est",
	},
	{
		name:    "MatchesFalseWithGlobMatch_UpperCase",
		match:   "*esT",
		toMatch: "test",
		matched: false,
		result:  "",
	},
	{
		name:    "MatchesFalseWithIncorrectGlobMatch",
		match:   "t*st",
		toMatch: "tset",
		matched: false,
		result:  "",
	},
}

func Test_GlobMatch(t *testing.T) {
	RegisterTestingT(t)

	for _, test := range globMatchTests {
		t.Run(test.name, func(t *testing.T) {

			isMatched, result := matchers.GlobMatch(test.match, test.toMatch)

			Expect(isMatched).To(Equal(test.matched))
			Expect(result).To(Equal(test.result))
		})
	}
}

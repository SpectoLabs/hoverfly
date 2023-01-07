package matchers_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/matching/matchers"
	. "github.com/onsi/gomega"
)

func Test_GlobMatch_MatchesFalseWithIncorrectDataType(t *testing.T) {
	RegisterTestingT(t)

	_, isMatched := matchers.GlobMatch(1, "yes", nil)
	Expect(isMatched).To(BeFalse())
}

func Test_GlobMatch_MatchesTrueWithGlobMatch(t *testing.T) {
	RegisterTestingT(t)

	matchedValue, isMatched := matchers.GlobMatch("t*st", `test`, nil)
	expectTrueWithTestString(isMatched, matchedValue)
}

func Test_GlobMatch_MatchesTrueWithGlobMatch_MatchesZeroExtraCharactersAtEnd(t *testing.T) {
	RegisterTestingT(t)

	matchedValue, isMatched := matchers.GlobMatch("test*", `test`, nil)
	expectTrueWithTestString(isMatched, matchedValue)
}

func Test_GlobMatch_MatchesTrueWithGlobMatch_MatchesZeroExtraCharactersAtStart(t *testing.T) {
	RegisterTestingT(t)

	matchedValue, isMatched := matchers.GlobMatch("*test", `test`, nil)
	expectTrueWithTestString(isMatched, matchedValue)
}

func Test_GlobMatch_MatchesTrueWithGlobMatch_MatchesZeroExtraCharactersAtStartAndEnd(t *testing.T) {
	RegisterTestingT(t)

	matchedValue, isMatched := matchers.GlobMatch("*test*", `test`, nil)
	expectTrueWithTestString(isMatched, matchedValue)
}

func Test_GlobMatch_MatchesTrueWithGlobMatch_MatchesUpperCase(t *testing.T) {
	RegisterTestingT(t)

	matchedValue, isMatched := matchers.GlobMatch("*est", `Test`, nil)
	Expect(isMatched).To(BeTrue())
	Expect(matchedValue).Should(Equal("Test"))
}

func Test_GlobMatch_MatchesTrueWithGlobMatch_MatchesLowerCase(t *testing.T) {
	RegisterTestingT(t)

	matchedValue, isMatched := matchers.GlobMatch("*est", `test`, nil)
	expectTrueWithTestString(isMatched, matchedValue)
}

func Test_GlobMatch_MatchesTrueWithGlobMatch_MatchesAstrik(t *testing.T) {
	RegisterTestingT(t)

	_, isMatched1 := matchers.GlobMatch("*est", `*est`, nil)
	Expect(isMatched1).To(BeTrue())

	_, isMatched2 := matchers.GlobMatch("t*est", `t*est`, nil)
	Expect(isMatched2).To(BeTrue())

	_, isMatched3 := matchers.GlobMatch("test*", `test*`, nil)
	Expect(isMatched3).To(BeTrue())
}

func Test_GlobMatch_MatchesFalseWithGlobMatch_UpperCase(t *testing.T) {
	RegisterTestingT(t)

	_, isMatched := matchers.GlobMatch("*esT", `test`, nil)
	Expect(isMatched).To(BeFalse())
}

func Test_GlobMatch_MatchesFalseWithIncorrectGlobMatch(t *testing.T) {
	RegisterTestingT(t)

	_, isMatched := matchers.GlobMatch("t*st", `tset`, nil)
	Expect(isMatched).To(BeFalse())
}

func expectTrueWithTestString(isMatched bool, matchedValue string) {
	Expect(isMatched).To(BeTrue())
	Expect(matchedValue).Should(Equal("test"))
}

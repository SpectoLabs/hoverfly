package matchers_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/matching/matchers"
	. "github.com/onsi/gomega"
)

func Test_GlobMatch_MatchesFalseWithIncorrectDataType(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.GlobMatch(1, "yes")).To(BeFalse())
}

func Test_GlobMatch_MatchesTrueWithGlobMatch(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.GlobMatch("t*st", `test`)).To(BeTrue())
}

func Test_GlobMatch_MatchesTrueWithGlobMatch_MatchesZeroExtraCharactersAtEnd(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.GlobMatch("test*", `test`)).To(BeTrue())
}

func Test_GlobMatch_MatchesTrueWithGlobMatch_MatchesZeroExtraCharactersAtStart(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.GlobMatch("*test", `test`)).To(BeTrue())
}

func Test_GlobMatch_MatchesTrueWithGlobMatch_MatchesZeroExtraCharactersAtStartAndEnd(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.GlobMatch("*test*", `test`)).To(BeTrue())
}

func Test_GlobMatch_MatchesTrueWithGlobMatch_MatchesUpperCase(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.GlobMatch("*est", `Test`)).To(BeTrue())
}

func Test_GlobMatch_MatchesTrueWithGlobMatch_MatchesLowerCase(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.GlobMatch("*est", `test`)).To(BeTrue())
}

func Test_GlobMatch_MatchesTrueWithGlobMatch_MatchesAstrik(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.GlobMatch("*est", `*est`)).To(BeTrue())
	Expect(matchers.GlobMatch("t*est", `t*est`)).To(BeTrue())
	Expect(matchers.GlobMatch("test*", `test*`)).To(BeTrue())
}

func Test_GlobMatch_MatchesFalseWithGlobMatch_UpperCase(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.GlobMatch("*esT", `test`)).To(BeFalse())
}

func Test_GlobMatch_MatchesFalseWithIncorrectGlobMatch(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.GlobMatch("t*st", `tset`)).To(BeFalse())
}

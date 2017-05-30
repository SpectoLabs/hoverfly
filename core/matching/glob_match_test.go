package matching_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/matching"
	. "github.com/onsi/gomega"
)

func Test_GlobMatch_MatchesTrueWithGlobMatch(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.GlobMatch("t*st", `test`)).To(BeTrue())
}

func Test_GlobMatch_MatchesTrueWithGlobMatch_MatchesZeroExtraCharactersAtEnd(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.GlobMatch("test*", `test`)).To(BeTrue())
}

func Test_GlobMatch_MatchesTrueWithGlobMatch_MatchesZeroExtraCharactersAtStart(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.GlobMatch("*test", `test`)).To(BeTrue())
}

func Test_GlobMatch_MatchesTrueWithGlobMatch_MatchesZeroExtraCharactersAtStartAndEnd(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.GlobMatch("*test*", `test`)).To(BeTrue())
}

func Test_GlobMatch_MatchesTrueWithGlobMatch_MatchesUpperCase(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.GlobMatch("*est", `Test`)).To(BeTrue())
}

func Test_GlobMatch_MatchesTrueWithGlobMatch_MatchesLowerCase(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.GlobMatch("*est", `test`)).To(BeTrue())
}

func Test_GlobMatch_MatchesFalseWithGlobMatch_UpperCase(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.GlobMatch("*esT", `test`)).To(BeFalse())
}

func Test_GlobMatch_MatchesFalseWithIncorrectGlobMatch(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.GlobMatch("t*st", `tset`)).To(BeFalse())
}

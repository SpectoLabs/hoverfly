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

func Test_GlobMatch_MatchesFalseWithIncorrectGlobMatch(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.GlobMatch("t*st", `tset`)).To(BeFalse())
}

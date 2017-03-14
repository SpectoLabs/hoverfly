package matching_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/matching"
	. "github.com/onsi/gomega"
)

func Test_RegexMatch_MatchesTrueWithRegexMatch(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.RegexMatch("t[o|a|e]st", `test`)).To(BeTrue())
}

func Test_RegexMatch_MatchesFalseWithIncorrectRegexMatch(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.RegexMatch("t[o|a]st", `test`)).To(BeFalse())
}

package matchers_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/matching/matchers"
	. "github.com/onsi/gomega"
)

func Test_RegexMatch_MatchesFalseWithIncorrectDataType(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.RegexMatch(1, "yes")).To(BeFalse())
}
func Test_RegexMatch_MatchesTrueWithRegexMatch(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.RegexMatch("t[o|a|e]st", `test`)).To(BeTrue())
}

func Test_RegexMatch_MatchesFalseWithIncorrectRegexMatch(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.RegexMatch("t[o|a]st", `test`)).To(BeFalse())
}

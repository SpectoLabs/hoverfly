package matchers_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/matching/matchers"
	. "github.com/onsi/gomega"
)

func Test_NegationMatch_MatchesTrueWithIncorrectDataType(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.NegationMatch(1, "yes")).To(BeTrue())
}

func Test_NegationMatch_MatchesTrueWithNil(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.NegationMatch(nil, "yes")).To(BeTrue())
}

func Test_NegationMatch_MatchesFalseWithExactMatch(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.NegationMatch("yes", "yes")).To(BeFalse())
}

func Test_NegationMatch_MatchesTrueWithIncorrectExactMatch(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.NegationMatch("yes", "no")).To(BeTrue())
}

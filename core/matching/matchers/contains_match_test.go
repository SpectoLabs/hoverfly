package matchers_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/matching/matchers"
	. "github.com/onsi/gomega"
)

func Test_ContainsMatch_MatchesFalseWithIncorrectDataType(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.ContainsMatch("hello", "yes")).To(BeFalse())
}

func Test_ContainsMatch_MatchesTrueWithArrayContainingValues(t *testing.T) {
	RegisterTestingT(t)

	arr := [3]string{"q1", "q2", "q3"}
	Expect(matchers.ContainsMatch(arr[:], "q1;q2")).To(BeTrue())
}

func Test_ContainsMatch_MatchesFalseArrayNotContainingSomeValues(t *testing.T) {
	RegisterTestingT(t)

	arr := [3]string{"q1", "q2", "q3"}
	Expect(matchers.ExactMatch(arr[:], "q1;q4")).To(BeFalse())
}

func Test_ContainsMatch_MatchesFalseArrayNotContainingAllValues(t *testing.T) {
	RegisterTestingT(t)

	arr := [4]string{"q1", "q2", "q3", "q4"}
	Expect(matchers.ExactMatch(arr[:], "q5;q6")).To(BeFalse())
}

func Test_ContainsMatch_MatchesFalseArrayIsEmpty(t *testing.T) {
	RegisterTestingT(t)

	arr := [0]string{}
	Expect(matchers.ExactMatch(arr, "q5;q5")).To(BeFalse())
}

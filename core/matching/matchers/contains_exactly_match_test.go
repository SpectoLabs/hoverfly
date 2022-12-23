package matchers_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/matching/matchers"
	. "github.com/onsi/gomega"
)

func Test_ContainsExactlyMatch_MatchesFalseWithIncorrectDataType(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.ContainsExactlyMatch("hello", "yes")).To(BeFalse())
}

func Test_ContainsExactlyMatch_MatchesTrueWithIdenticalArray(t *testing.T) {
	RegisterTestingT(t)

	arr := [3]string{"q1", "q2", "q3"}
	Expect(matchers.ContainsExactlyMatch(arr[:], "q1;q2;q3")).To(BeTrue())
}

func Test_ContainsExactlyMatch_MatchesFalseSameArrayInDifferentOrder(t *testing.T) {
	RegisterTestingT(t)

	arr := [3]string{"q1", "q2", "q3"}
	Expect(matchers.ContainsExactlyMatch(arr[:], "q1;q3;q2")).To(BeFalse())
}

func Test_ContainsExactlyMatch_MatchesFalseDifferentArray(t *testing.T) {
	RegisterTestingT(t)

	arr := [4]string{"q1", "q2", "q3", "q4"}
	Expect(matchers.ContainsExactlyMatch(arr[:], "q5;q6")).To(BeFalse())
}

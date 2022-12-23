package matchers_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/matching/matchers"
	. "github.com/onsi/gomega"
)

func Test_ContainsOnlyMatch_MatchesFalseWithIncorrectDataType(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.ContainsOnlyMatch("hello", "yes")).To(BeFalse())
}

func Test_ContainsOnlyMatch_MatchesTrueWithIdenticalArray(t *testing.T) {
	RegisterTestingT(t)

	arr := [3]string{"q1", "q2", "q3"}
	Expect(matchers.ContainsOnlyMatch(arr[:], "q1;q2;q3")).To(BeTrue())
}

func Test_ContainsOnlyMatch_MatchesTrueSameArrayInDifferentOrder(t *testing.T) {
	RegisterTestingT(t)

	arr := [3]string{"q1", "q2", "q3"}
	Expect(matchers.ContainsOnlyMatch(arr[:], "q1;q3;q2")).To(BeTrue())
}

func Test_ContainsOnlyMatch_MatchesFalseDifferentArray(t *testing.T) {
	RegisterTestingT(t)

	arr := [4]string{"q1", "q2", "q3", "q4"}
	Expect(matchers.ContainsOnlyMatch(arr[:], "q5;q6")).To(BeFalse())
}

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

func Test_ContainsMatch_MatchesTrueWithArrayContainingAllValues(t *testing.T) {
	RegisterTestingT(t)

	arr := [3]string{"q1", "q2", "q3"}
	Expect(matchers.ContainsMatch(arr[:], "q1;q2")).To(BeTrue())
}

func Test_ContainsMatch_MatchesTrueWithArrayContainingAllValuesWithMatcherAsArrayOfInterface(t *testing.T) {
	RegisterTestingT(t)

	arr := [3]interface{}{"q1", "q2", "q3"}
	Expect(matchers.ContainsMatch(arr[:], "q1;q2")).To(BeTrue())
}

func Test_ContainsMatch_MatchesTrueWithArrayContainingSomeValues(t *testing.T) {
	RegisterTestingT(t)

	arr := [3]string{"q1", "q2", "q3"}
	Expect(matchers.ContainsMatch(arr[:], "q1;q5;q6")).To(BeTrue())
}

func Test_ContainsMatch_MatchesFalseWithArrayNotContainingAllValues(t *testing.T) {
	RegisterTestingT(t)

	arr := [4]string{"q1", "q2", "q3", "q4"}
	Expect(matchers.ContainsMatch(arr[:], "q5;q6")).To(BeFalse())
}

func Test_ContainsMatch_MatchesFalseWithArrayIsEmpty(t *testing.T) {
	RegisterTestingT(t)

	arr := [0]string{}
	Expect(matchers.ContainsMatch(arr, "q5;q6")).To(BeFalse())
}

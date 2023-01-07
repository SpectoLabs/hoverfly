package matchers_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/matching/matchers"
	. "github.com/onsi/gomega"
)

func Test_ContainsExactlyMatch_MatchesFalseWithIncorrectDataType(t *testing.T) {
	RegisterTestingT(t)

	_, isMatched := matchers.ContainsExactlyMatch("hello", "yes", nil)
	Expect(isMatched).To(BeFalse())
}

func Test_ContainsExactlyMatch_MatchesTrueWithIdenticalArray(t *testing.T) {
	RegisterTestingT(t)

	arr := [3]string{"q1", "q2", "q3"}
	matchedValue, isMatched := matchers.ContainsExactlyMatch(arr[:], "q1;q2;q3", nil)
	Expect(isMatched).To(BeTrue())
	Expect(matchedValue).Should(Equal("q1;q2;q3"))
}

func Test_ContainsExactlyMatch_MatchesTrueWithIdenticalArrayWithMatcherAsArrayOfInterface(t *testing.T) {
	RegisterTestingT(t)

	arr := [3]interface{}{"q1", "q2", "q3"}
	matchedValue, isMatched := matchers.ContainsExactlyMatch(arr[:], "q1;q2;q3", nil)
	Expect(isMatched).To(BeTrue())
	Expect(matchedValue).Should(Equal("q1;q2;q3"))
}

func Test_ContainsExactlyMatch_MatchesFalseWithSameArrayInDifferentOrder(t *testing.T) {
	RegisterTestingT(t)

	arr := [3]string{"q1", "q2", "q3"}
	_, isMatched := matchers.ContainsExactlyMatch(arr[:], "q1;q3;q2", nil)
	Expect(isMatched).To(BeFalse())
}

func Test_ContainsExactlyMatch_MatchesFalseWithDifferentArray(t *testing.T) {
	RegisterTestingT(t)

	arr := [4]string{"q1", "q2", "q3", "q4"}
	_, isMatched := matchers.ContainsExactlyMatch(arr[:], "q5;q6", nil)
	Expect(isMatched).To(BeFalse())
}

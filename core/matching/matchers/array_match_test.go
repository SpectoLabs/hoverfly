package matchers_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/matching/matchers"
	. "github.com/onsi/gomega"
)

func Test_ArrayMatch_ReturnsFalseWithIncorrectDataType(t *testing.T) {
	RegisterTestingT(t)

	configMap := make(map[string]interface{})
	_, isMatched := matchers.ArrayMatch("hello", "yes", configMap)
	Expect(isMatched).To(BeFalse())
}

func Test_ArrayMatch_ReturnsTrueWithIdenticalArray(t *testing.T) {
	RegisterTestingT(t)

	configMap := getConfiguration(false, false, false)
	arr := [3]string{"q1", "q2", "q3"}
	matchedValue, isMatched := matchers.ArrayMatch(arr[:], "q1;q2;q3", configMap)
	Expect(isMatched).To(BeTrue())
	Expect(matchedValue).Should(Equal(matchedValue))
}

func Test_ArrayMatch_ReturnsTrueWithAllKnownsInArrayAndNotIgnoringUnkowns(t *testing.T) {
	RegisterTestingT(t)

	configMap := getConfiguration(false, true, true)
	arr := [3]string{"q1", "q2", "q3"}
	matchedValue, isMatched := matchers.ArrayMatch(arr[:], "q1;q3;q2;q1;q3", configMap)
	Expect(isMatched).To(BeTrue())
	Expect(matchedValue).Should(Equal("q1;q3;q2;q1;q3"))
}
func Test_ArrayMatch_ReturnsFalseWithUnkownsInArrayAndNotIgnoringUnkowns(t *testing.T) {
	RegisterTestingT(t)

	configMap := getConfiguration(false, true, true)
	arr := [3]string{"q1", "q2", "q3"}
	_, isMatched := matchers.ArrayMatch(arr[:], "q1;q4;q3;q2", configMap)
	Expect(isMatched).To(BeFalse())
}

func Test_ArrayMatch_ReturnsTrueWithInSameOrderAndNotIgnoringOrder(t *testing.T) {
	RegisterTestingT(t)

	configMap := getConfiguration(true, true, false)
	arr := [3]string{"q1", "q2", "q3"}
	matchedValue, isMatched := matchers.ArrayMatch(arr[:], "q1;q2;q3;q2;q4", configMap)
	Expect(isMatched).To(BeTrue())
	Expect(matchedValue).Should(Equal("q1;q2;q3;q2;q4"))
}

func Test_ArrayMatch_ReturnsFalseWithOutOfOrderAndNotIgnoringOrder(t *testing.T) {
	RegisterTestingT(t)

	configMap := getConfiguration(true, true, false)
	arr := [3]string{"q1", "q2", "q3"}
	_, isMatched := matchers.ArrayMatch(arr[:], "q1;q3;q3;q2;q4", configMap)
	Expect(isMatched).To(BeFalse())
}

func Test_ArrayMatch_ReturnsTrueWithSameOccurrencesAndNotIgnoringOccurrences(t *testing.T) {
	RegisterTestingT(t)

	configMap := getConfiguration(true, false, true)
	arr := [3]string{"q1", "q2", "q3"}
	matchedValue, isMatched := matchers.ArrayMatch(arr[:], "q1;q3;q0;q2;q4", configMap)
	Expect(isMatched).To(BeTrue())
	Expect(matchedValue).Should(Equal("q1;q3;q0;q2;q4"))
}

func Test_ArrayMatch_ReturnsFalseWithDifferentNoOfOccurrencesAndNotIgnoringOccurrences(t *testing.T) {
	RegisterTestingT(t)

	configMap := getConfiguration(true, false, true)
	arr := [3]string{"q1", "q2", "q3"}
	_, isMatched := matchers.ArrayMatch(arr[:], "q1;q3;q3;q2;q4", configMap)
	Expect(isMatched).To(BeFalse())
}

func getConfiguration(ignoreUnknown, ignoreOccurrences, ignoreOrder bool) map[string]interface{} {

	configMap := make(map[string]interface{})
	configMap[matchers.IGNORE_UNKNOWN] = ignoreUnknown
	configMap[matchers.IGNORE_ORDER] = ignoreOrder
	configMap[matchers.IGNORE_OCCURRENCES] = ignoreOccurrences
	return configMap
}

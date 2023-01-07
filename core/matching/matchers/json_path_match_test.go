package matchers_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/matching/matchers"
	. "github.com/onsi/gomega"
)

func Test_JsonPathMatch_MatchesFalseWithIncorrectDataType(t *testing.T) {
	RegisterTestingT(t)

	_, isMatched := matchers.JsonPathMatch(1, "yes", nil)
	Expect(isMatched).To(BeFalse())
}
func Test_JsonPathMatch_MatchesFalseWithInvalidJsonPath(t *testing.T) {
	RegisterTestingT(t)

	_, isMatched := matchers.JsonPathMatch("test", `{"test": "field"}`, nil)
	Expect(isMatched).To(BeFalse())
}

func Test_JsonPathMatch_MatchesTrueWithJsonMatch_GetSingleElement(t *testing.T) {
	RegisterTestingT(t)

	matchedContext, isMatched := matchers.JsonPathMatch("$.test", `{"test": "field"}`, nil)
	Expect(isMatched).To(BeTrue())
	Expect(matchedContext).Should(Equal("field"))
}

func Test_JsonPathMatch_MatchesFalseWithIncorrectJsonMatch_GetSingleElement(t *testing.T) {
	RegisterTestingT(t)

	_, isMatched := matchers.JsonPathMatch("$.notAField", `{"test": "field"}`, nil)
	Expect(isMatched).To(BeFalse())
}

func Test_JsonPathMatch_MatchesTrueWithJsonMatch_GetElementFromArray(t *testing.T) {
	RegisterTestingT(t)

	matchedValue, isMatched := matchers.JsonPathMatch("$.test[1]", `{"test": [{}, {}]}`, nil)
	Expect(matchedValue).Should(Equal("{}"))
	Expect(isMatched).To(BeTrue())
}

func Test_JsonPathMatch_MatchesFalseWithIncorrectJsonMatch_GetElementFromArray(t *testing.T) {
	RegisterTestingT(t)

	_, isMatched := matchers.JsonPathMatch("$.test[2]", `{"test": [{}, {}]}`, nil)
	Expect(isMatched).To(BeFalse())
}

func Test_JsonPathMatch_MatchesTrueWithJsonMatch_GetArrayElement(t *testing.T) {
	RegisterTestingT(t)

	matchedValue, isMatched := matchers.JsonPathMatch("$.test", `{"test": ["hello", "world"]}`, nil)
	Expect(isMatched).To(BeTrue())
	Expect(matchedValue).Should(Equal("[\"hello\",\"world\"]"))
}

func Test_JsonPathMatch_MatchesTrueWithJsonMatch_WithExpression(t *testing.T) {
	RegisterTestingT(t)

	matchedValue, isMatched := matchers.JsonPathMatch("$.test[?(@.field == \"test\")]", `{"test": [{"field": "test"}]}`, nil)
	Expect(isMatched).To(BeTrue())
	Expect(matchedValue).Should(Equal("{\"field\":\"test\"}"))
}

func Test_JsonPathMatch_MatchesFalseWithIncorrectJsonMatch_WithExpression(t *testing.T) {
	RegisterTestingT(t)

	_, isMatched := matchers.JsonPathMatch("$.test[*]?(@.field == \"test\")", `{"test": [{"field": "not-test"}]}`, nil)
	Expect(isMatched).To(BeFalse())
}

func Test_JsonPathMatch_MatchesTrueWithJsonMatch_GetSingleElement_WhereRootIsArray(t *testing.T) {
	RegisterTestingT(t)

	matchedValue, isMatched := matchers.JsonPathMatch("$[0].test", `[{"test": "field"}]`, nil)
	Expect(isMatched).To(BeTrue())
	Expect(matchedValue).Should(Equal("field"))
}

func Test_JsonPathMatch_MatchesTrueWithJsonMatch_GetCompleteObject_WhereElementIsObject(t *testing.T) {
	RegisterTestingT(t)

	matchedValue, isMatched := matchers.JsonPathMatch("$.test", `{"test": {"field1":"value1", "field2":"value2"}}`, nil)
	Expect(isMatched).To(BeTrue())
	Expect(matchedValue).Should(Equal("{\"field1\":\"value1\",\"field2\":\"value2\"}"))
}

func Test_JsonPathMatch_MatchesTrueWithJsonMatch_GetArray_WhereElementIsArrayObject(t *testing.T) {
	RegisterTestingT(t)

	matchedValue, isMatched := matchers.JsonPathMatch("$.test[1:3]", `{"test": [{"field1":"value1"}, {"field2":"value2"}, {"field3":"value3"}, {"field4":"value4"}]}`, nil)
	Expect(isMatched).To(BeTrue())
	Expect(matchedValue).Should(Equal("[{\"field2\":\"value2\"},{\"field3\":\"value3\"}]"))
}

// TODO the following JSONPath expressions are not supported at the moment
//func Test_JsonPathMatch_MatchesTrueWithJsonMatch_WithRootLevelObjectFilter(t *testing.T) {
//	RegisterTestingT(t)
//
//	Expect(matchers.JsonPathMatch("$[?(@.field == \"test\")]", `{"field": "test"}`)).To(BeTrue())
//}
//
//func Test_JsonPathMatch_MatchesTrueWithJsonMatch_WithObjectFilter(t *testing.T) {
//	RegisterTestingT(t)
//
//	Expect(matchers.JsonPathMatch("$.test[?(@.field == \"test\")]", `{"test": {"field": "test"}}`)).To(BeTrue())
//}

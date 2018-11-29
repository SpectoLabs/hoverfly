package matchers_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/matching/matchers"
	. "github.com/onsi/gomega"
)

func Test_JsonPathMatch_MatchesFalseWithIncorrectDataType(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.JsonPathMatch(1, "yes")).To(BeFalse())
}
func Test_JsonPathMatch_MatchesFalseWithInvalidJsonPath(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.JsonPathMatch("test", `{"test": "field"}`)).To(BeFalse())
}

func Test_JsonPathMatch_MatchesTrueWithJsonMatch_GetSingleElement(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.JsonPathMatch("$.test", `{"test": "field"}`)).To(BeTrue())
}

func Test_JsonPathMatch_MatchesFalseWithIncorrectJsonMatch_GetSingleElement(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.JsonPathMatch("$.notAField", `{"test": "field"}`)).To(BeFalse())
}

func Test_JsonPathMatch_MatchesTrueWithJsonMatch_GetElementFromArray(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.JsonPathMatch("$.test[1]", `{"test": [{}, {}]}`)).To(BeTrue())
}

func Test_JsonPathMatch_MatchesFalseWithIncorrectJsonMatch_GetElementFromArray(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.JsonPathMatch("$.test[2]", `{"test": [{}, {}]}`)).To(BeFalse())
}

func Test_JsonPathMatch_MatchesTrueWithJsonMatch_WithExpression(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.JsonPathMatch("$.test[?(@.field == \"test\")]", `{"test": [{"field": "test"}]}`)).To(BeTrue())
}

func Test_JsonPathMatch_MatchesFalseWithIncorrectJsonMatch_WithExpression(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.JsonPathMatch("$.test[*]?(@.field == \"test\")", `{"test": [{"field": "not-test"}]}`)).To(BeFalse())
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

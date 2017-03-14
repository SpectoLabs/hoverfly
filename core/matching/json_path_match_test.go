package matching_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/matching"
	. "github.com/onsi/gomega"
)

func Test_JsonPathMatch_MatchesFalseWithInvalidJsonPath(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.JsonPathMatch("test", `{"test": "field"}`)).To(BeFalse())
}

func Test_JsonPathMatch_MatchesTrueWithJsonMatch_GetSingleElement(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.JsonPathMatch("$.test", `{"test": "field"}`)).To(BeTrue())
}

func Test_JsonPathMatch_MatchesFalseWithIncorrectJsonMatch_GetSingleElement(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.JsonPathMatch("$.notAField", `{"test": "field"}`)).To(BeFalse())
}

func Test_JsonPathMatch_MatchesTrueWithJsonMatch_GetElementFromArray(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.JsonPathMatch("$.test[1]", `{"test": [{}, {}]}`)).To(BeTrue())
}

func Test_JsonPathMatch_MatchesFalseWithIncorrectJsonMatch_GetElementFromArray(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.JsonPathMatch("$.test[2]", `{"test": [{}, {}]}`)).To(BeFalse())
}

func Test_JsonPathMatch_MatchesTrueWithJsonMatch_WithExpression(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.JsonPathMatch("$.test[*]?(@.field == \"test\")", `{"test": [{"field": "test"}]}`)).To(BeTrue())
}

func Test_JsonPathMatch_MatchesFalseWithIncorrectJsonMatch_WithExpression(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.JsonPathMatch("$.test[*]?(@.field == \"test\")", `{"test": [{"field": "not-test"}]}`)).To(BeFalse())
}

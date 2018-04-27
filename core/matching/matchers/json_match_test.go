package matchers_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/matching/matchers"
	. "github.com/onsi/gomega"
)

func Test_JsonMatch_MatchesFalseWithIncorrectDataType(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.JsonMatch(1, "yes")).To(BeFalse())
}
func Test_JsonMatch_MatchesTrueWithJSON(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.JsonMatch(`{"test":{"json":true,"minified":true}}`, `{"test":{"json":true,"minified":true}}`)).To(BeTrue())
}

func Test_JsonMatch_MatchesTrueWithJSON_InADifferentOrder(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.JsonMatch(`{"test":{"minified":true, "json":true}}`, `{"test":{"json":true,"minified":true}}`)).To(BeTrue())
}

func Test_JsonMatch_MatchesTrueWithUnminifiedJSON(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.JsonMatch(`{"test":{"json":true,"minified":true}}`, `{
		"test": {
			"json": true,
			"minified": true
		}
	}`)).To(BeTrue())
}

func Test_JsonMatch_MatchesFalseWithInvalidJSONAsMatcher(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.JsonMatch(`"test":"json":true,"minified"`, `{
		"test": {
			"json": true,
			"minified": true
		}
	}`)).To(BeFalse())
}

func Test_JsonMatch_MatchesFalseWithInvalidJSON(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.JsonMatch(`{"test":{"json":true,"minified":true}}`, `{
		"test": {
			"json": true,
			"minified": 
		}
	}`)).To(BeFalse())
}

func Test_JsonMatch_MatchesTrueWithTwoEmptyString(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.JsonMatch(``, ``)).To(BeTrue())
}

func Test_JsonMatch_MatchesFalseAgainstEmptyString(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.JsonMatch(`{
		"test": {
			"json": true,
			"minified": 
		}
	}`, ``)).To(BeFalse())
}

func Test_JsonMatch_MatchesFalseWithEmptyString(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.JsonMatch(``, `{
		"test": {
			"json": true,
			"minified": 
		}
	}`)).To(BeFalse())
}

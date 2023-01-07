package matchers_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/matching/matchers"
	. "github.com/onsi/gomega"
)

func Test_JsonMatch_MatchesFalseWithIncorrectDataType(t *testing.T) {
	RegisterTestingT(t)

	_, isMatched := matchers.JsonMatch(1, "yes", nil)
	Expect(isMatched).To(BeFalse())
}
func Test_JsonMatch_MatchesTrueWithJSON(t *testing.T) {
	RegisterTestingT(t)

	matchedValue, isMatch := matchers.JsonMatch(`{"test":{"json":true,"minified":true}}`, `{"test":{"json":true,"minified":true}}`, nil)
	Expect(isMatch).To(BeTrue())
	Expect(matchedValue).Should(Equal(`{"test":{"json":true,"minified":true}}`))
}

func Test_JsonMatch_MatchesTrueWithJSON_InADifferentOrder(t *testing.T) {
	RegisterTestingT(t)

	matchedValue, isMatch := matchers.JsonMatch(`{"test":{"minified":true, "json":true}}`, `{"test":{"json":true,"minified":true}}`, nil)
	Expect(isMatch).To(BeTrue())
	Expect(matchedValue).Should(Equal(`{"test":{"json":true,"minified":true}}`))
}

func Test_JsonMatch_MatchesTrueWithUnminifiedJSON(t *testing.T) {
	RegisterTestingT(t)

	matchedValue, isMatched := matchers.JsonMatch(`{"test":{"json":true,"minified":true}}`, `{
		"test": {
			"json": true,
			"minified": true
		}
	}`, nil)
	Expect(isMatched).To(BeTrue())
	Expect(matchedValue).Should(Equal(`{
		"test": {
			"json": true,
			"minified": true
		}
	}`))
}

func Test_JsonMatch_MatchesFalseWithInvalidJSONAsMatcher(t *testing.T) {
	RegisterTestingT(t)

	_, isMatched := matchers.JsonMatch(`"test":"json":true,"minified"`, `{
		"test": {
			"json": true,
			"minified": true
		}
	}`, nil)
	Expect(isMatched).To(BeFalse())
}

func Test_JsonMatch_MatchesFalseWithInvalidJSON(t *testing.T) {
	RegisterTestingT(t)

	_, isMatched := matchers.JsonMatch(`{"test":{"json":true,"minified":true}}`, `{
		"test": {
			"json": true,
			"minified": 
		}
	}`, nil)
	Expect(isMatched).To(BeFalse())
}

func Test_JsonMatch_MatchesTrueWithTwoEmptyString(t *testing.T) {
	RegisterTestingT(t)

	matchedValue, isMatched := matchers.JsonMatch(``, ``, nil)
	Expect(isMatched).To(BeTrue())
	Expect(matchedValue).Should(Equal(``))
}

func Test_JsonMatch_MatchesFalseAgainstEmptyString(t *testing.T) {
	RegisterTestingT(t)

	_, isMatched := matchers.JsonMatch(`{
		"test": {
			"json": true,
			"minified": 
		}
	}`, ``, nil)
	Expect(isMatched).To(BeFalse())
}

func Test_JsonMatch_MatchesFalseWithEmptyString(t *testing.T) {
	RegisterTestingT(t)

	_, isMatched := matchers.JsonMatch(``, `{
		"test": {
			"json": true,
			"minified": 
		}
	}`, nil)
	Expect(isMatched).To(BeFalse())
}

// Should not ignore JSON array order by default
func Test_JsonMatch_MatchesFalseWithJSONRootAsArray_InADifferentOrder(t *testing.T) {
	RegisterTestingT(t)

	_, isMatched := matchers.JsonMatch(`[{"minified": true}, {"json": true}]`, `[{"json":true},{"minified":true}]`, nil)
	Expect(isMatched).To(BeFalse())
}

func Test_JsonMatch_MatchesTrueWithUnminifiedJSONRootAsArray(t *testing.T) {
	RegisterTestingT(t)

	matchedValue, isMatched := matchers.JsonMatch(`[{"minified": true}, {"json": true}]`, `[
		{
			"minified": true
		}, {
			"json": true
		}
	]`, nil)
	Expect(isMatched).To(BeTrue())
	Expect(matchedValue).Should(Equal(`[
		{
			"minified": true
		}, {
			"json": true
		}
	]`))
}

func Test_JsonMatch_MatchesTrueWithJSONRootAsArray_WithDataInDifferentOrder(t *testing.T) {
	RegisterTestingT(t)

	matchedValue, isMatched := matchers.JsonMatch(`[{"minified":true, "json":true}]`, `[{"json":true,"minified":true}]`, nil)
	Expect(isMatched).To(BeTrue())
	Expect(matchedValue).Should(Equal(`[{"json":true,"minified":true}]`))
}

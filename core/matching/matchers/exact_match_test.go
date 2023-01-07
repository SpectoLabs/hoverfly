package matchers_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/matching/matchers"
	. "github.com/onsi/gomega"
)

func Test_ExactMatch_MatchesFalseWithIncorrectDataType(t *testing.T) {
	RegisterTestingT(t)

	_, isMatched := matchers.ExactMatch(1, "yes", nil)
	Expect(isMatched).To(BeFalse())
}

func Test_ExactMatch_MatchesTrueWithExactMatch(t *testing.T) {
	RegisterTestingT(t)

	matchedValue, isMatched := matchers.ExactMatch("yes", "yes", nil)
	Expect(isMatched).To(BeTrue())
	Expect(matchedValue).Should(Equal("yes"))
}

func Test_ExactMatch_MatchesFalseWithIncorrectExactMatch(t *testing.T) {
	RegisterTestingT(t)

	_, isMatched := matchers.ExactMatch("yes", "no", nil)
	Expect(isMatched).To(BeFalse())
}

func Test_ExactMatch_MatchesTrueWithJSON(t *testing.T) {
	RegisterTestingT(t)

	_, isMatched := matchers.ExactMatch(`{"test":{"json":true,"minified":true}}`, `{"test":{"json":true,"minified":true}}`, nil)
	Expect(isMatched).To(BeTrue())
}

func Test_ExactMatch_MatchesTrueWithUnminifiedJSON(t *testing.T) {
	RegisterTestingT(t)

	_, isMatchedValue := matchers.ExactMatch(`{"test":{"json":true,"minified":true}}`, `{
		"test": {
			"json": true,
			"minified": true
		}
	}`, nil)
	Expect(isMatchedValue).To(BeFalse())
}

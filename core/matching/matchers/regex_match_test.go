package matchers_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/matching/matchers"
	. "github.com/onsi/gomega"
)

func Test_RegexMatch_MatchesFalseWithIncorrectDataType(t *testing.T) {
	RegisterTestingT(t)

	_, isMatched := matchers.RegexMatch(1, "yes", nil)
	Expect(isMatched).To(BeFalse())
}
func Test_RegexMatch_MatchesTrueWithRegexMatch(t *testing.T) {
	RegisterTestingT(t)

	matchedValue, isMatched := matchers.RegexMatch("t[o|a|e]st", `test`, nil)
	expectTrueWithTestString(isMatched, matchedValue)
}

func Test_RegexMatch_MatchesFalseWithIncorrectRegexMatch(t *testing.T) {
	RegisterTestingT(t)

	_, isMatched := matchers.RegexMatch("t[o|a]st", `test`, nil)
	Expect(isMatched).To(BeFalse())
}

package matchers_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/v2/core/matching/matchers"
	. "github.com/onsi/gomega"
)

func Test_ExactMatch_MatchesFalseWithIncorrectDataType(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.ExactMatch(1, "yes")).To(BeFalse())
}

func Test_ExactMatch_MatchesTrueWithExactMatch(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.ExactMatch("yes", "yes")).To(BeTrue())
}

func Test_ExactMatch_MatchesFalseWithIncorrectExactMatch(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.ExactMatch("yes", "no")).To(BeFalse())
}

func Test_ExactMatch_MatchesTrueWithJSON(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.ExactMatch(`{"test":{"json":true,"minified":true}}`, `{"test":{"json":true,"minified":true}}`)).To(BeTrue())
}

func Test_ExactMatch_MatchesTrueWithUnminifiedJSON(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.ExactMatch(`{"test":{"json":true,"minified":true}}`, `{
		"test": {
			"json": true,
			"minified": true
		}
	}`)).To(BeFalse())
}

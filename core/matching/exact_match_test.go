package matching_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/matching"
	. "github.com/onsi/gomega"
)

func Test_ExactMatch_MatchesTrueWithExactMatch(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.ExactMatch("yes", "yes")).To(BeTrue())
}

func Test_ExactMatch_MatchesFalseWithIncorrectExactMatch(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.ExactMatch("yes", "no")).To(BeFalse())
}

func Test_ExactMatch_MatchesTrueWithJSON(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.ExactMatch(`{"test":{"json":true,"minified":true}}`, `{"test":{"json":true,"minified":true}}`)).To(BeTrue())
}

func Test_ExactMatch_MatchesTrueWithUnminifiedJSON(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.ExactMatch(`{"test":{"json":true,"minified":true}}`, `{
		"test": {
			"json": true,
			"minified": true
		}
	}`)).To(BeFalse())
}

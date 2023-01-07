package matchers_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/matching/matchers"
	. "github.com/onsi/gomega"
)

func Test_XmlMatch_MatchesFalseWithIncorrectDataType(t *testing.T) {
	RegisterTestingT(t)

	_, isMatched := matchers.XmlMatch(1, "yes", nil)
	Expect(isMatched).To(BeFalse())
}
func Test_XmlMatch_MatchesTrueWithXML(t *testing.T) {
	RegisterTestingT(t)

	matchedValue, isMatched := matchers.XmlMatch(`<xml><document><test></document>`, `<xml><document><test></document>`, nil)
	Expect(isMatched).To(BeTrue())
	Expect(matchedValue).Should(Equal(`<xml><document><test></document>`))
}

func Test_XmlMatch_MatchesTrueWithUnminifiedXml(t *testing.T) {
	RegisterTestingT(t)

	matchedValue, isMatched := matchers.XmlMatch(`<xml>
	<document>
		<test key="value">cat</test>
	</document>`, `<xml><document><test key="value">cat</test></document>`, nil)
	Expect(isMatched).To(BeTrue())
	Expect(matchedValue).Should(Equal(`<xml><document><test key="value">cat</test></document>`))
}

func Test_XmlMatch_MatchesFalseWithNotMatchingXml(t *testing.T) {
	RegisterTestingT(t)

	_, isMatched := matchers.XmlMatch(`<xml>
	<document>
		<test key="value">cat</test>
	</document>`, `<xml><document><test key="different">cat</test></document>`, nil)
	Expect(isMatched).To(BeFalse())
}

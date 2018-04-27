package matchers_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/matching/matchers"
	. "github.com/onsi/gomega"
)

func Test_XmlMatch_MatchesFalseWithIncorrectDataType(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.XmlMatch(1, "yes")).To(BeFalse())
}
func Test_XmlMatch_MatchesTrueWithXML(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.XmlMatch(`<xml><document><test></document>`, `<xml><document><test></document>`)).To(BeTrue())
}

func Test_XmlMatch_MatchesTrueWithUnminifiedXml(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.XmlMatch(`<xml>
		<document>
			<test key="value">cat</test>
		</document>`, `<xml><document><test key="value">cat</test></document>`)).To(BeTrue())
}

func Test_XmlMatch_MatchesFalseWithNotMatchingXml(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.XmlMatch(`<xml>
		<document>
			<test key="value">cat</test>
		</document>`, `<xml><document><test key="different">cat</test></document>`)).To(BeFalse())
}

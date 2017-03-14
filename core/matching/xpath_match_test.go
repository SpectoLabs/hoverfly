package matching_test

import (
	"encoding/xml"
	"testing"

	"github.com/SpectoLabs/hoverfly/core/matching"
	. "github.com/onsi/gomega"
)

func Test_XpathMatch_MatchesTrue(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.XpathMatch("/root/text", xml.Header+"<root><text>test</text></root>")).To(BeTrue())
}

func Test_XpathMatch_MatchesFalseWithIncorectXpathMatch(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.XpathMatch("/pop", xml.Header+"<root><text>test</text></root>")).To(BeFalse())
}

func Test_XpathMatch_MatchesTrue_GetAnElementFromAnArray(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.XpathMatch("/list/item[1]/field", xml.Header+"<list><item><field>test</field></item></list>")).To(BeTrue())
}

func Test_XpathMatch_MatchesFalse_GetAnElementFromAnArray(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.XpathMatch("/list/item[1]/pop", xml.Header+"<list><item><field>test</field></item></list>")).To(BeFalse())
}

func Test_XpathMatch_MatchesTrue_GetAttributeFromElement(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.XpathMatch("/list/item/field[@test]", xml.Header+"<list><item><field test=\"value\">test</field></item></list>")).To(BeTrue())
}

func Test_XpathMatch_MatchesFalse_GetAttributeFromElement(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.XpathMatch("/list/item/field[@pop]", xml.Header+"<list><item><field test=\"value\">test</field></item></list>")).To(BeFalse())
}

func Test_XpathMatch_MatchesTrue_GetElementWithNoValue(t *testing.T) {
	RegisterTestingT(t)

	Expect(matching.XpathMatch("/list/item/field", xml.Header+"<list><item><field></field></item></list>")).To(BeTrue())
}

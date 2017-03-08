package matching

import (
	"encoding/xml"
	"testing"

	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/SpectoLabs/hoverfly/core/util"
	. "github.com/onsi/gomega"
)

func Test_FieldMatcher_MatchesTrueWithExactMatch(t *testing.T) {
	RegisterTestingT(t)

	Expect(FieldMatcher(&models.RequestFieldMatchers{
		ExactMatch: util.StringToPointer("yes"),
	}, "yes")).To(BeTrue())
}

func Test_FieldMatcher_MatchesFalseWithIncorrectExactMatch(t *testing.T) {
	RegisterTestingT(t)

	Expect(FieldMatcher(&models.RequestFieldMatchers{
		ExactMatch: util.StringToPointer("yes"),
	}, "no")).To(BeFalse())
}

func Test_FieldMatcher_MatchesTrueWithNilMatchers(t *testing.T) {
	RegisterTestingT(t)

	Expect(FieldMatcher(nil, "no")).To(BeTrue())
}

func Test_FieldMatcher_MatchesTrueWithXpathMatch(t *testing.T) {
	RegisterTestingT(t)

	Expect(FieldMatcher(&models.RequestFieldMatchers{
		XpathMatch: util.StringToPointer("/root/text"),
	}, xml.Header+"<root><text>test</text></root>")).To(BeTrue())
}

func Test_FieldMatcher_MatchesFalseWithIncorectXpathMatch(t *testing.T) {
	RegisterTestingT(t)

	Expect(FieldMatcher(&models.RequestFieldMatchers{
		XpathMatch: util.StringToPointer("/pop"),
	}, xml.Header+"<root><text>test</text></root>")).To(BeFalse())
}

func Test_FieldMatcher_MatchesTrueWithXpathMatch_GetAnElementFromAnArray(t *testing.T) {
	RegisterTestingT(t)

	Expect(FieldMatcher(&models.RequestFieldMatchers{
		XpathMatch: util.StringToPointer("/list/item[1]/field"),
	}, xml.Header+"<list><item><field>test</field></item></list>")).To(BeTrue())
}

func Test_FieldMatcher_MatchesFalseWithInvalidXpathMatch_GetAnElementFromAnArray(t *testing.T) {
	RegisterTestingT(t)

	Expect(FieldMatcher(&models.RequestFieldMatchers{
		XpathMatch: util.StringToPointer("/list/item[1]/pop"),
	}, xml.Header+"<list><item><field>test</field></item></list>")).To(BeFalse())
}

func Test_FieldMatcher_MatchesTrueWithXpathMatch_GetAttributeFromElement(t *testing.T) {
	RegisterTestingT(t)

	Expect(FieldMatcher(&models.RequestFieldMatchers{
		XpathMatch: util.StringToPointer("/list/item/field[@test]"),
	}, xml.Header+"<list><item><field test=\"value\">test</field></item></list>")).To(BeTrue())
}

func Test_FieldMatcher_MatchesFalseWithInvalidXpathMatch_GetAttributeFromElement(t *testing.T) {
	RegisterTestingT(t)

	Expect(FieldMatcher(&models.RequestFieldMatchers{
		XpathMatch: util.StringToPointer("/list/item/field[@pop]"),
	}, xml.Header+"<list><item><field test=\"value\">test</field></item></list>")).To(BeFalse())
}

func Test_FieldMatcher_MatchesTrueWithXpathMatch_GetElementWithNoValue(t *testing.T) {
	RegisterTestingT(t)

	Expect(FieldMatcher(&models.RequestFieldMatchers{
		XpathMatch: util.StringToPointer("/list/item/field"),
	}, xml.Header+"<list><item><field></field></item></list>")).To(BeTrue())
}

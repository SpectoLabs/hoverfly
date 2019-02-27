package matchers_test

import (
	"encoding/xml"
	"testing"

	"github.com/SpectoLabs/hoverfly/core/matching/matchers"
	. "github.com/onsi/gomega"
)

var xpathMatchTests = []matchTest{
	{
		name:    "MatchesFalseWithIncorrectDataType",
		match:   1,
		toMatch: "yes",
		matched: false,
		result:  "",
	},
	{
		name:    "MatchesTrue",
		match:   "/root/text",
		toMatch: xml.Header+"<root><text>test</text></root>",
		matched: true,
		result:  "test",
	},
	{
		name:    "MatchesFalseWithIncorrectXpathMatch",
		match:   "/pop",
		toMatch: xml.Header+"<root><text>test</text></root>",
		matched: false,
		result:  "",
	},
	{
		name:    "MatchesTrue_GetAnElementFromAnArray",
		match:   "/list/item[1]/field",
		toMatch: xml.Header+"<list><item><field>test</field></item></list>",
		matched: true,
		result:  "test",
	},
	{
		name:    "MatchesFalse_GetAnElementFromAnArray",
		match:   "/list/item[1]/pop",
		toMatch: xml.Header+"<list><item><field>test</field></item></list>",
		matched: false,
		result:  "",
	},
	{
		name:    "MatchesTrue_GetAttributeFromElement",
		match:   "/list/item/field[@test]",
		toMatch: xml.Header+"<list><item><field test=\"value\">test</field></item></list>",
		matched: true,
		result:  "test",
	},
	{
		name:    "MatchesFalse_GetAttributeFromElement",
		match:   "/list/item/field[@pop]",
		toMatch: xml.Header+"<list><item><field test=\"value\">test</field></item></list>",
		matched: false,
		result:  "",
	},
	{
		name:    "MatchesTrue_GetElementWithNoValue",
		match:   "/list/item/field",
		toMatch: xml.Header+"<list><item><field></field></item></list>",
		matched: true,
		result:  "",
	},
	{
		name:    "MatchesTrue_WithoutHeader",
		match:   "/list/item/field",
		toMatch: "<list><item><field></field></item></list>",
		matched: true,
		result:  "",
	},
}

func Test_XpathMatch(t *testing.T) {
	RegisterTestingT(t)

	for _, test := range xpathMatchTests {
		t.Run(test.name, func(t *testing.T) {

			isMatched, result := matchers.XpathMatch(test.match, test.toMatch)

			Expect(isMatched).To(Equal(test.matched))
			Expect(result).To(Equal(test.result))
		})
	}
}
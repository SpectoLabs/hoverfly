package matchers_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/matching/matchers"
	. "github.com/onsi/gomega"
)

var xmlMatchTests = []matchTest{
	{
		name:    "MatchesFalseWithIncorrectDataType",
		match:   1,
		toMatch: "yes",
		matched: false,
		result:  "",
	},
	{
		name:    "MatchesTrueWithXML",
		match:   `<xml><document><test></document>`,
		toMatch: `<xml><document><test></document>`,
		matched: true,
		result:  `<xml><document><test></document>`,
	},
	{
		name:      "MatchesTrueWithUnminifiedXml",
		match:     `<xml>
		<document>
			<test key="value">cat</test>
		</document>`,
		toMatch: `<xml><document><test key="value">cat</test></document>`,
		matched: true,
		result:  `<xml><document><test key="value">cat</test></document>`,
	},
	{
		name:      "MatchesFalseWithNotMatchingXml",
		match:     `<xml>
		<document>
			<test key="value">cat</test>
		</document>`,
		toMatch: `<xml><document><test key="different">cat</test></document>`,
		matched: false,
		result:  "",
	},
}

func Test_XmlMatch(t *testing.T) {
	RegisterTestingT(t)

	for _, test := range xmlMatchTests {
		t.Run(test.name, func(t *testing.T) {

			isMatched, result := matchers.XmlMatch(test.match, test.toMatch)

			Expect(isMatched).To(Equal(test.matched))
			Expect(result).To(Equal(test.result))
		})
	}
}

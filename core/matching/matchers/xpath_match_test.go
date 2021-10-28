package matchers_test

import (
	"encoding/xml"
	"testing"

	"github.com/SpectoLabs/hoverfly/core/matching/matchers"
	. "github.com/onsi/gomega"
)

func Test_XpathMatch_MatchesFalseWithIncorrectDataType(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.XpathMatch(1, "yes")).To(BeFalse())
}

func Test_XpathMatch_MatchesTrue(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.XpathMatch("/root/text", xml.Header+"<root><text>test</text></root>")).To(BeTrue())
}

func Test_XpathMatch_MatchesFalseWithIncorectXpathMatch(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.XpathMatch("/pop", xml.Header+"<root><text>test</text></root>")).To(BeFalse())
}

func Test_XpathMatch_MatchesTrue_GetAnElementFromAnArray(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.XpathMatch("/list/item[1]/field", xml.Header+"<list><item><field>test</field></item></list>")).To(BeTrue())
}

func Test_XpathMatch_MatchesFalse_GetAnElementFromAnArray(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.XpathMatch("/list/item[1]/pop", xml.Header+"<list><item><field>test</field></item></list>")).To(BeFalse())
}

func Test_XpathMatch_MatchesTrue_GetAttributeFromElement(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.XpathMatch("/list/item/field[@test]", xml.Header+"<list><item><field test=\"value\">test</field></item></list>")).To(BeTrue())
}

func Test_XpathMatch_MatchesFalse_GetAttributeFromElement(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.XpathMatch("/list/item/field[@pop]", xml.Header+"<list><item><field test=\"value\">test</field></item></list>")).To(BeFalse())
}

func Test_XpathMatch_MatchesTrue_GetElementWithNoValue(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.XpathMatch("/list/item/field", xml.Header+"<list><item><field></field></item></list>")).To(BeTrue())
}

func Test_XpathMatch_MatchesTrue_WithoutHeader(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.XpathMatch("/list/item/field", "<list><item><field></field></item></list>")).To(BeTrue())
}

func Test_XpathMatch_MatchesTrue_WithNameSpacePrefix(t *testing.T) {
	RegisterTestingT(t)

	Expect(matchers.XpathMatch("/soapenv:Envelope/soapenv:Header/head:MessageHeader/head:From/head:Id",
		`<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:head="http://www.test.com/Header_01" xmlns:v1="http://www.test.com/GetStatement/v1">
			   <soapenv:Header>
				  <head:MessageHeader>
					 <head:From>
						<head:Id>Test</head:Id>
					 </head:From>
				  </head:MessageHeader>
			   </soapenv:Header>
			   <soapenv:Body>
				  <v1:GetCMSAccountStatementReq>         
				  </v1:GetCMSAccountStatementReq>
			   </soapenv:Body>
			</soapenv:Envelope>`)).To(BeTrue())
}

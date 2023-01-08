package matchers_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/matching/matchers"
	. "github.com/onsi/gomega"
)

func Test_IdentityValueGenerator_ReturnsActualValue(t *testing.T) {
	RegisterTestingT(t)

	value := matchers.IdentityValueGenerator("*es*", "Test")
	Expect(value).Should(Equal(value))
}

func Test_JsonPathValueGenerator_ReturnsStringValue(t *testing.T) {
	RegisterTestingT(t)
	value := matchers.JsonPathMatcherValueGenerator("$.field", `{"field": 1234}`)
	Expect(value).Should(Equal("1234"))
}

func Test_JsonPathValueGenerator_ReturnsArrayValue(t *testing.T) {
	RegisterTestingT(t)
	value := matchers.JsonPathMatcherValueGenerator("$.field", `{"field":["test1","test2","test3","test4"]}`)
	Expect(value).Should(Equal(`["test1","test2","test3","test4"]`))
}

func Test_JsonPathValueGenerator_ReturnsArrayObject(t *testing.T) {
	RegisterTestingT(t)

	value := matchers.JsonPathMatcherValueGenerator("$.field[1:3]", `{"field":[{"field1":"value1"}, {"field2":"value2"}, {"field3":"value3"}, {"field4":"value4"}]}`)
	Expect(value).Should(Equal(`[{"field2":"value2"},{"field3":"value3"}]`))
}

func Test_JsonPathValueGenerator_ReturnsObject(t *testing.T) {
	RegisterTestingT(t)

	value := matchers.JsonPathMatcherValueGenerator("$.field", `{"field":{"key1":"value1"}}`)
	Expect(value).Should(Equal(`{"key1":"value1"}`))
}

func Test_XPathValueGenerator_ReturnsValue(t *testing.T) {
	RegisterTestingT(t)

	value := matchers.XPathMatchValueGenerator("/document/name", "<document><id>1234</id><name>Test</name></document>")
	Expect(value).Should(Equal(value))
}

func Test_XPathValueGenerator_ReturnsEmbeddedJson(t *testing.T) {
	RegisterTestingT(t)

	value := matchers.XPathMatchValueGenerator("/document/details", `<document><details>{"id":1234,"name":"test"}</details></document>`)
	Expect(value).Should(Equal(`{"id":1234,"name":"test"}`))
}

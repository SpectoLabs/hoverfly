package templating_test

import (
	"testing"
	. "github.com/onsi/gomega"
	"net/http"
	"github.com/SpectoLabs/hoverfly/core/templating"
)

func Test_ShouldCreateTemplatingDataPathParamsFromRequest(t *testing.T) {
	RegisterTestingT(t)

	actual := NewTemplatingDataFomRequest("http://www.test.com/foo/bar?cheese=1&ham=2&cheese=3")

	Expect(actual.Request.PathParam).To(ConsistOf("foo", "bar"))
}

func Test_ShouldCreateTemplatingDataPathParamsFromRequestWithNoPathParams(t *testing.T) {
	RegisterTestingT(t)

	actual := NewTemplatingDataFomRequest("http://www.test.com?cheese=1&ham=2&cheese=3")

	Expect(actual.Request.PathParam).To(BeEmpty())
}

func Test_ShouldCreateTemplatingDataQueryParamsFromRequest(t *testing.T) {
	RegisterTestingT(t)

	actual := NewTemplatingDataFomRequest("http://www.test.com/foo/bar?cheese=1&ham=2&cheese=3")

	Expect(actual.Request.QueryParam).To(HaveKeyWithValue("cheese", []string{"1", "3"}))
	Expect(actual.Request.QueryParam).To(HaveKeyWithValue("ham", []string{"2"}))
	Expect(actual.Request.QueryParam).To(HaveLen(2))
}

func Test_ShouldCreateTemplatingDataQueryParamsFromRequestWithNoQueryParams(t *testing.T) {
	RegisterTestingT(t)

	actual := NewTemplatingDataFomRequest("http://www.test.com/foo/bar")

	Expect(actual.Request.QueryParam).To(BeEmpty())
}

func Test_ShouldCreateTemplatingDataHttpScheme(t *testing.T) {
	RegisterTestingT(t)

	actual := NewTemplatingDataFomRequest("http://www.test.com/foo/bar")

	Expect(actual.Request.QueryParam).To(BeEmpty())
}

func Test_ShouldCreateTemplatingDataQueryScheme(t *testing.T) {
	RegisterTestingT(t)

	actual := NewTemplatingDataFomRequest("http://www.test.com/foo/bar")

	Expect(actual.Request.Scheme).To(Equal("http"))

	actual = NewTemplatingDataFomRequest("https://www.test.com/foo/bar")

	Expect(actual.Request.Scheme).To(Equal("https"))
}

func TestApplyTemplateWithQueryParams(t *testing.T) {
	RegisterTestingT(t)

	request, err := http.NewRequest("GET", "http://www.foo.com/foo/bar?singular=1&multiple=2&multiple=3", nil)
	Expect(err).To(BeNil())

	template, err := templating.ApplyTemplate(request, `
Scheme: {{ Request.Scheme }}

Query param value: {{ Request.QueryParam.singular }}

Query param value by index: {{ Request.QueryParam.multiple.[0] }}
Query param value by index: {{ Request.QueryParam.multiple.[1] }}

List of query param values: {{ Request.QueryParam.multiple}}
Looping through query params: {{#each Request.QueryParam.multiple}}{{ this }}-{{/each}}

Path param value: {{ Request.PathParam.[0] }}
All path param values: {{ Request.PathParam }}
Looping through path params: {{#each Request.PathParam}}{{ this }}-{{/each}}`)

	Expect(err).To(BeNil())

	Expect(template).To(Equal(`
Scheme: http

Query param value: 1

Query param value by index: 2
Query param value by index: 3

List of query param values: 23
Looping through query params: 2-3-

Path param value: foo
All path param values: foobar
Looping through path params: foo-bar-`))
}

func TestTemplatingWithParametersWhichDoNotExistDoNotErrorAndAreEmpty(t *testing.T) {
	RegisterTestingT(t)

	request, err := http.NewRequest("GET", "http://www.foo.com", nil)
	Expect(err).To(BeNil())

	template, err := templating.ApplyTemplate(request, `
Scheme:{{ Request.Scheme }}

Query param value:{{ Request.QueryParam.singular }}

Query param value by index:{{ Request.QueryParam.multiple.[0] }}
Query param value by index:{{ Request.QueryParam.multiple.[1] }}

List of query param values:{{ Request.QueryParam.multiple}}
Looping through query params:{{#each Request.QueryParam.multiple}}{{ this }}{{/each}}

Path param value:{{ Request.PathParam.[0] }}
All path param values:{{ Request.PathParam }}
Looping through path params:{{#each Request.PathParam}}{{ this }}{{/each}}`)

	Expect(err).To(BeNil())

	Expect(template).To(Equal(`
Scheme:http

Query param value:

Query param value by index:
Query param value by index:

List of query param values:
Looping through query params:

Path param value:
All path param values:
Looping through path params:`))
}

func NewTemplatingDataFomRequest(url string) *templating.TemplatingData {
	r, err := http.NewRequest("GET", url, nil)

	Expect(err).To(BeNil())

	return templating.NewTemplatingDataFromRequest(r)
}

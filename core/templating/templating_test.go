package templating_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/SpectoLabs/hoverfly/core/templating"
	. "github.com/onsi/gomega"
)

func Test_ShouldCreateTemplatingDataPathsFromRequest(t *testing.T) {
	RegisterTestingT(t)

	actual := templating.NewTemplator().NewTemplatingData(
		&models.RequestDetails{
			Scheme:      "http",
			Destination: "test.com",
			Path:        "/foo/bar",
		},
		&models.Literals{},
		&models.Variables{},
		make(map[string]string),
	)

	Expect(actual.Request.Path).To(ConsistOf("foo", "bar"))
}

func Test_ShouldCreateTemplatingDataPathsFromRequestWithNoPaths(t *testing.T) {
	RegisterTestingT(t)

	actual := templating.NewTemplator().NewTemplatingData(
		&models.RequestDetails{
			Scheme:      "http",
			Destination: "test.com",
		},
		&models.Literals{},
		&models.Variables{},
		make(map[string]string),
	)

	Expect(actual.Request.Path).To(BeEmpty())
}

func Test_ShouldCreateTemplatingDataQueryParamsFromRequest(t *testing.T) {
	RegisterTestingT(t)

	actual := templating.NewTemplator().NewTemplatingData(
		&models.RequestDetails{
			Scheme:      "http",
			Destination: "test.com",
			Query: map[string][]string{
				"cheese": {"1", "3"},
				"ham":    {"2"},
			},
		},
		&models.Literals{},
		&models.Variables{},
		make(map[string]string),
	)

	Expect(actual.Request.QueryParam).To(HaveKeyWithValue("cheese", []string{"1", "3"}))
	Expect(actual.Request.QueryParam).To(HaveKeyWithValue("ham", []string{"2"}))
	Expect(actual.Request.QueryParam).To(HaveLen(2))
}

func Test_ShouldCreateTemplatingDataQueryParamsFromRequestWithNoQueryParams(t *testing.T) {
	RegisterTestingT(t)

	actual := templating.NewTemplator().NewTemplatingData(
		&models.RequestDetails{
			Scheme:      "http",
			Destination: "test.com",
		},
		&models.Literals{},
		&models.Variables{},
		make(map[string]string),
	)

	Expect(actual.Request.QueryParam).To(BeEmpty())
}

func Test_ShouldCreateTemplatingDataHttpScheme(t *testing.T) {
	RegisterTestingT(t)

	actual := templating.NewTemplator().NewTemplatingData(
		&models.RequestDetails{
			Scheme:      "http",
			Destination: "test.com",
		},
		&models.Literals{},
		&models.Variables{},
		make(map[string]string),
	)

	Expect(actual.Request.Scheme).To(Equal("http"))
}

func Test_ShouldCreateTemplatingDataHeaderFromRequest(t *testing.T) {
	RegisterTestingT(t)

	actual := templating.NewTemplator().NewTemplatingData(
		&models.RequestDetails{
			Scheme:      "http",
			Destination: "test.com",
			Headers: map[string][]string{
				"cheese": {"1", "3"},
				"ham":    {"2"},
			},
		},
		&models.Literals{},
		&models.Variables{},
		make(map[string]string),
	)

	Expect(actual.Request.Header).To(HaveKeyWithValue("cheese", []string{"1", "3"}))
	Expect(actual.Request.Header).To(HaveKeyWithValue("ham", []string{"2"}))
	Expect(actual.Request.Header).To(HaveLen(2))
}

func Test_ShouldCreateTemplatingDataHeaderFromRequestWithNoHeader(t *testing.T) {
	RegisterTestingT(t)

	actual := templating.NewTemplator().NewTemplatingData(
		&models.RequestDetails{
			Scheme:      "http",
			Destination: "test.com",
		},
		&models.Literals{},
		&models.Variables{},
		make(map[string]string),
	)

	Expect(actual.Request.Header).To(BeEmpty())
}

func TestApplyTemplateWithRequestDetails(t *testing.T) {
	RegisterTestingT(t)

	requestDetails := &models.RequestDetails{
		Method:      "GET",
		Scheme:      "http",
		Destination: "foo.com",
		Path:        "/foo/bar",
		Query: map[string][]string{
			"singular": {"1"},
			"multiple": {"2", "3"},
		},
		Headers: map[string][]string{
			"X-Singular": {"a"},
			"X-Multiple": {"b", "c"},
		},
	}

	template, err := ApplyTemplate(requestDetails,
		make(map[string]string),
		`
Scheme: {{ Request.Scheme }}

Query param value: {{ Request.QueryParam.singular }}

Query param value by index: {{ Request.QueryParam.multiple.[0] }}
Query param value by index: {{ Request.QueryParam.multiple.[1] }}

List of query param values: {{ Request.QueryParam.multiple}}
Looping through query params: {{#each Request.QueryParam.multiple}}{{ this }}-{{/each}}

Header value: {{ Request.Header.X-Singular }}
Header value by index: {{ Request.Header.X-Multiple.[0] }}
Header value by index: {{ Request.Header.X-Multiple.[1] }}
List of header values: {{ Request.Header.X-Multiple}}

Path param value: {{ Request.Path.[0] }}
All path param values: {{ Request.Path }}
Looping through path params: {{#each Request.Path}}{{ this }}-{{/each}}`)

	Expect(err).To(BeNil())

	Expect(template).To(Equal(`
Scheme: http

Query param value: 1

Query param value by index: 2
Query param value by index: 3

List of query param values: 23
Looping through query params: 2-3-

Header value: a
Header value by index: b
Header value by index: c
List of header values: bc

Path param value: foo
All path param values: foobar
Looping through path params: foo-bar-`))
}

func TestTemplatingWithParametersWhichDoNotExistDoNotErrorAndAreEmpty(t *testing.T) {
	RegisterTestingT(t)

	requestDetails := &models.RequestDetails{
		Method:      "GET",
		Scheme:      "http",
		Destination: "foo.com",
	}

	template, err := ApplyTemplate(requestDetails,
		map[string]string{
			"one": "A",
			"two": "B",
		},
		`
Scheme:{{ Request.Scheme }}

Query param value:{{ Request.QueryParam.singular }}

Query param value by index:{{ Request.QueryParam.multiple.[0] }}
Query param value by index:{{ Request.QueryParam.multiple.[1] }}

List of query param values:{{ Request.QueryParam.multiple}}
Looping through query params:{{#each Request.QueryParam.multiple}}{{ this }}{{/each}}

Header value: {{ Request.Header.X-Singular }}
Header value by index: {{ Request.Header.X-Multiple.[0] }}
Header value by index: {{ Request.Header.X-Multiple.[1] }}
List of header values: {{ Request.Header.X-Multiple}}

Path param value:{{ Request.Path.[0] }}
All path param values:{{ Request.Path }}
Looping through path params:{{#each Request.Path}}{{ this }}{{/each}}

State One: {{ State.one }}
State Two: {{ State.two }}`)

	Expect(err).To(BeNil())

	Expect(template).To(Equal(`
Scheme:http

Query param value:

Query param value by index:
Query param value by index:

List of query param values:
Looping through query params:

Header value: 
Header value by index: 
Header value by index: 
List of header values: 

Path param value:
All path param values:
Looping through path params:

State One: A
State Two: B`))
}

func Test_ApplyTemplate_now(t *testing.T) {
	RegisterTestingT(t)

	template, err := ApplyTemplate(&models.RequestDetails{}, make(map[string]string), `{{now "" "unix"}}`)

	Expect(err).To(BeNil())

	Expect(template).To(Not(Equal(ContainSubstring(`{{now "" "unix"}}`))))
}

func Test_ApplyTemplate_randomString(t *testing.T) {
	RegisterTestingT(t)

	template, err := ApplyTemplate(&models.RequestDetails{}, make(map[string]string), `{{randomString}}`)

	Expect(err).To(BeNil())

	Expect(template).To(Not(Equal(ContainSubstring(`{{randomString}}`))))
}

func Test_ApplyTemplate_randomStringLength(t *testing.T) {
	RegisterTestingT(t)

	template, err := ApplyTemplate(&models.RequestDetails{}, make(map[string]string), `{{randomStringLength 2}}`)

	Expect(err).To(BeNil())

	Expect(template).To(Not(Equal(ContainSubstring(`{{randomStringLength 2}}`))))
	Expect(template).To(HaveLen(2))
}

func Test_ApplyTemplate_randomBoolean(t *testing.T) {
	RegisterTestingT(t)

	template, err := ApplyTemplate(&models.RequestDetails{}, make(map[string]string), `{{randomBoolean}}`)

	Expect(err).To(BeNil())

	Expect(template).To(Not(Equal(ContainSubstring(`{{randomBoolean}}`))))
}

func Test_ApplyTemplate_randomInteger(t *testing.T) {
	RegisterTestingT(t)

	template, err := ApplyTemplate(&models.RequestDetails{}, make(map[string]string), `{{randomInteger}}`)

	Expect(err).To(BeNil())

	Expect(template).To(Not(Equal(ContainSubstring(`{{randomInteger}}`))))
}

func Test_ApplyTemplate_randomIntegerRange(t *testing.T) {
	RegisterTestingT(t)

	template, err := ApplyTemplate(&models.RequestDetails{}, make(map[string]string), `{{randomIntegerRange 7 8}}`)

	Expect(err).To(BeNil())

	Expect(template).To(Not(Equal(ContainSubstring(`{{randomIntegerRange 7 8}}`))))
}

func Test_ApplyTemplate_randomFloat(t *testing.T) {
	RegisterTestingT(t)

	template, err := ApplyTemplate(&models.RequestDetails{}, make(map[string]string), `{{randomFloat}}`)

	Expect(err).To(BeNil())

	Expect(template).To(Not(Equal(ContainSubstring(`{{randomFloat}}`))))
}

func Test_ApplyTemplate_randomFloatRange(t *testing.T) {
	RegisterTestingT(t)

	template, err := ApplyTemplate(&models.RequestDetails{}, make(map[string]string), `{{randomFloatRange 7.0 8.0}}`)

	Expect(err).To(BeNil())

	Expect(template).To(Not(Equal(ContainSubstring(`{{randomFloatRange 7.0 8.0}}`))))
}

func Test_ApplyTemplate_randomEmail(t *testing.T) {
	RegisterTestingT(t)

	template, err := ApplyTemplate(&models.RequestDetails{}, make(map[string]string), `{{randomEmai}}`)

	Expect(err).To(BeNil())
	Expect(template).To(Not(Equal(ContainSubstring(`{{randomEmail}}`))))
}

func Test_ApplyTemplate_randomIPv4(t *testing.T) {
	RegisterTestingT(t)

	template, err := ApplyTemplate(&models.RequestDetails{}, make(map[string]string), `{{randomIPv4}}`)

	Expect(err).To(BeNil())

	Expect(template).To(Not(Equal(ContainSubstring(`{{randomIPv4}}`))))
}

func Test_ApplyTemplate_randomIPv6(t *testing.T) {
	RegisterTestingT(t)

	template, err := ApplyTemplate(&models.RequestDetails{}, make(map[string]string), `{{randomIPv6}}`)

	Expect(err).To(BeNil())

	Expect(template).To(Not(Equal(ContainSubstring(`{{randomIPv6}}`))))
}

func Test_ApplyTemplate_randomUuid(t *testing.T) {
	RegisterTestingT(t)

	template, err := ApplyTemplate(&models.RequestDetails{}, make(map[string]string), `{{randomUuid}}`)

	Expect(err).To(BeNil())

	Expect(template).To(Not(Equal(ContainSubstring(`{{randomUuid}}`))))
}

func Test_ApplyTemplate_Request_Body_Jsonpath(t *testing.T) {
	RegisterTestingT(t)

	template, err := ApplyTemplate(&models.RequestDetails{
		Body: `{ "name": "Ben" }`,
	}, make(map[string]string), `{{ Request.Body 'jsonpath' '$.name' }}`)

	Expect(err).To(BeNil())

	Expect(template).To(Equal("Ben"))
}

func Test_ApplyTemplate_Request_Body_JsonPath_Unescaped(t *testing.T) {
	RegisterTestingT(t)

	template, err := ApplyTemplate(&models.RequestDetails{
		Body: `{ "name": "O'Reilly" }`,
	}, make(map[string]string), `{{{ Request.Body 'jsonpath' '$.name' }}}`)

	Expect(err).To(BeNil())

	Expect(template).To(Equal("O'Reilly"))
}

func Test_ApplyTemplate_ReplaceStringInQueryParams(t *testing.T) {
	RegisterTestingT(t)

	template, err := ApplyTemplate(&models.RequestDetails{
		Query: map[string][]string{
			"sound": {"oink,oink,oink"},
		},
	}, make(map[string]string), `{{ replace Request.QueryParam.sound 'oink' 'moo' }}`)

	Expect(err).To(BeNil())

	Expect(template).To(Equal(`moo,moo,moo`))
}

func Test_VarSetToNilInCaseOfInvalidArgsPassed(t *testing.T) {
	RegisterTestingT(t)
	templator := templating.NewTemplator()

	vars := &models.Variables{
		models.Variable{
			Name:     "varOne",
			Function: "faker",
		},
	}

	actual := templator.NewTemplatingData(
		&models.RequestDetails{
			Scheme:      "http",
			Destination: "test.com",
		},
		&models.Literals{},
		vars,
		make(map[string]string),
	)

	Expect(actual.Vars["varOne"]).To(BeNil())

}

func ApplyTemplate(requestDetails *models.RequestDetails, state map[string]string, responseBody string) (string, error) {
	templator := templating.NewTemplator()
	template, _ := templator.ParseTemplate(responseBody)

	return templator.RenderTemplate(template, requestDetails, &models.Literals{}, &models.Variables{}, state)
}

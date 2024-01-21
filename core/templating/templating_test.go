package templating_test

import (
	"github.com/SpectoLabs/hoverfly/core/journal"
	"testing"

	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/SpectoLabs/hoverfly/core/templating"
	. "github.com/onsi/gomega"
)

// We need to run the csv templating tests first because the datasource is part of the templatingHelper which can be registered only once
// for each runtime, raymond unfortunately doesn't provide a way to unregister helpers for us to isolate the test data.
func Test_ApplyTemplate_ParseCsvAndReturnMatchedString(t *testing.T) {
	RegisterTestingT(t)

	template, err := ApplyTemplate(&models.RequestDetails{}, make(map[string]string), `{{csv 'test-csv1' 'Id' '2' 'Marks'}}`)

	Expect(err).To(BeNil())
	Expect(template).To(Equal(`56`))
}

func Test_ApplyTemplate_ParseCsvAndReturnFallbackStringIfNoMatchFound(t *testing.T) {
	RegisterTestingT(t)

	template, err := ApplyTemplate(&models.RequestDetails{}, make(map[string]string), `{{csv 'test-csv1' 'Id' '51' 'Marks'}}`)

	Expect(err).To(BeNil())
	Expect(template).To(Equal(`ABSENT`))
}

func Test_ApplyTemplate_ParseCsvAndReturnQueryStringIfNoMatchFound(t *testing.T) {
	RegisterTestingT(t)

	template, err := ApplyTemplate(&models.RequestDetails{}, make(map[string]string), `{{csv 'test-csv2' 'Id' '51' 'Marks'}}`)

	Expect(err).To(BeNil())
	Expect(template).To(Equal(`{{ csv test-csv2 Id 51 Marks }}`))
}

func Test_ApplyTemplate_ParseCsvByPassingRequestParamAndReturnMatchValue(t *testing.T) {
	RegisterTestingT(t)

	template, err := ApplyTemplate(&models.RequestDetails{
		Query: map[string][]string{"Id": {"1"}},
	}, make(map[string]string), `{{csv 'test-csv2' 'Id' 'Request.QueryParam.Id.[0]' 'Marks'}}`)

	Expect(err).To(BeNil())
	Expect(template).To(Equal(`55`))
}

func Test_ApplyTemplate_EachBlockWithCsvTemplatingFunction(t *testing.T) {
	RegisterTestingT(t)

	template, err := ApplyTemplate(&models.RequestDetails{
		Query: map[string][]string{"Ids": {"1", "2"}},
	}, make(map[string]string), `{{#each (Request.QueryParam.Ids) }}{{csv 'test-csv2' 'Id' this 'Marks'}} | {{/each}}`)

	Expect(err).To(BeNil())
	Expect(template).To(Equal(`55 | 56 | `))
}

func Test_ApplyTemplate_EachBlockWithCsvTemplatingFunctionAndLargeInteger(t *testing.T) {
	RegisterTestingT(t)

	template, err := ApplyTemplate(&models.RequestDetails{
		Body: `{"ids": [1, 5553686208582]}`,
	}, make(map[string]string), `{{#each (Request.Body 'jsonpath' '$.ids') }}{{csv 'test-csv2' 'Id' this 'Marks'}} | {{/each}}`)

	Expect(err).To(BeNil())
	Expect(template).To(Equal(`55 | 66 | `))
}

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
		&journal.Journal{},
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
		&journal.Journal{},
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
		&journal.Journal{},
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
		&journal.Journal{},
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
		&journal.Journal{},
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
		&journal.Journal{},
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
		&journal.Journal{},
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

func Test_ApplyTemplate_Request_Body_Jsonpath_LargeInt(t *testing.T) {
	RegisterTestingT(t)

	template, err := ApplyTemplate(&models.RequestDetails{
		Body: `{ "id": 5553686208582 }`,
	}, make(map[string]string), `{{ Request.Body 'jsonpath' '$.id' }}`)

	Expect(err).To(BeNil())

	Expect(template).To(Equal("5553686208582"))
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
		&journal.Journal{},
	)

	Expect(actual.Vars["varOne"]).To(BeNil())

}

func Test_VarSetToProperValueInCaseOfRequestDetailsPassedAsArgument(t *testing.T) {
	RegisterTestingT(t)
	templator := templating.NewTemplator()
	argumentsArray := toInterfaceSlice([]string{"Request.Path.[1]", ","})
	vars := &models.Variables{
		models.Variable{
			Name:      "splitRequestPath",
			Function:  "split",
			Arguments: argumentsArray,
		},
	}

	actual := templator.NewTemplatingData(
		&models.RequestDetails{
			Path: "/part1/foo,bar",
		},
		&models.Literals{},
		vars,
		make(map[string]string),
		&journal.Journal{},
	)

	Expect(actual.Vars["splitRequestPath"]).ToNot(BeNil())
	Expect(len(actual.Vars["splitRequestPath"].([]string))).To(Equal(2))
	Expect(actual.Vars["splitRequestPath"].([]string)[0]).To(Equal("foo"))
	Expect(actual.Vars["splitRequestPath"].([]string)[1]).To(Equal("bar"))

}

func Test_ApplyTemplate_add_integers(t *testing.T) {
	RegisterTestingT(t)

	template, err := ApplyTemplate(&models.RequestDetails{}, make(map[string]string), `{{ add 1 2 '0'}}`)

	Expect(err).To(BeNil())

	Expect(template).To(Equal("3"))
}

func Test_ApplyTemplate_add_floats(t *testing.T) {
	RegisterTestingT(t)

	template, err := ApplyTemplate(&models.RequestDetails{}, make(map[string]string), `{{ add 0.1 1.34 '0.00'}}`)

	Expect(err).To(BeNil())

	Expect(template).To(Equal("1.44"))
}

func Test_ApplyTemplate_add_floats_withRoundUp(t *testing.T) {
	RegisterTestingT(t)

	template, err := ApplyTemplate(&models.RequestDetails{}, make(map[string]string), `{{ add 0.1 1.34 '0.0'}} and {{ add 0.1 1.56 '0.0'}}`)

	Expect(err).To(BeNil())

	Expect(template).To(Equal("1.4 and 1.7"))
}

func Test_ApplyTemplate_add_number_without_format(t *testing.T) {
	RegisterTestingT(t)

	template, err := ApplyTemplate(&models.RequestDetails{}, make(map[string]string), `{{ add 0.1 1.34 ''}} and {{ add 1 2 ''}} and {{ add 0 0 ''}}`)

	Expect(err).To(BeNil())

	Expect(template).To(Equal("1.44 and 3 and 0"))
}

func Test_ApplyTemplate_add_NotNumber(t *testing.T) {
	RegisterTestingT(t)

	template, err := ApplyTemplate(&models.RequestDetails{}, make(map[string]string), `{{ add 'a' 'b' '0.00'}}`)

	Expect(err).To(BeNil())

	Expect(template).To(Equal("NaN"))
}

func Test_ApplyTemplate_subtract_numbers(t *testing.T) {
	RegisterTestingT(t)

	template, err := ApplyTemplate(&models.RequestDetails{}, make(map[string]string), `{{ subtract 10 0.99 ''}}`)

	Expect(err).To(BeNil())

	Expect(template).To(Equal("9.01"))
}

func Test_ApplyTemplate_mutiply_numbers(t *testing.T) {
	RegisterTestingT(t)

	template, err := ApplyTemplate(&models.RequestDetails{}, make(map[string]string), `{{ multiply 10 0.99 ''}}`)

	Expect(err).To(BeNil())

	Expect(template).To(Equal("9.9"))
}

func Test_ApplyTemplate_divide_numbers(t *testing.T) {
	RegisterTestingT(t)

	template, err := ApplyTemplate(&models.RequestDetails{}, make(map[string]string), `{{ divide 10 2.5 ''}}`)

	Expect(err).To(BeNil())

	Expect(template).To(Equal("4"))
}

func Test_ApplyTemplate_Arithmetic_Ops_With_Each_Block(t *testing.T) {
	RegisterTestingT(t)

	templator := templating.NewTemplator()

	requestDetails := &models.RequestDetails{
		Body: `{"lineitems":{"lineitem":[{"upc":"1001","quantity":"1","price":"3.50"},{"upc":"1002","quantity":"2","price":"4.50"}]}}`,
	}
	responseBody := `{{#each (Request.Body 'jsonpath' '$.lineitems.lineitem') }} {{ addToArray 'subtotal' (multiply (this.price) (this.quantity) '') }} {{/each}} total: {{ sum (getArray 'subtotal') '0.00' }}`

	template, _ := templator.ParseTemplate(responseBody)
	state := make(map[string]string)
	result, err := templator.RenderTemplate(template, requestDetails, &models.Literals{}, &models.Variables{}, state, &journal.Journal{})

	Expect(err).To(BeNil())
	Expect(result).To(Equal(` 3.5  9  total: 12.50`))

	// Running the second time should produce the same result because each execution has its own context data.
	result, err = templator.RenderTemplate(template, requestDetails, &models.Literals{}, &models.Variables{}, state, &journal.Journal{})
	Expect(err).To(BeNil())
	Expect(result).To(Equal(` 3.5  9  total: 12.50`))
}

func Test_ApplyTemplate_PutAndGetValue(t *testing.T) {
	RegisterTestingT(t)

	template, err := ApplyTemplate(&models.RequestDetails{
		Body: `{ "id": 5553686208582 }`,
	}, make(map[string]string), `{{ putValue 'id' (Request.Body 'jsonpath' '$.id') }} The ID was {{ getValue 'id' }}`)

	Expect(err).To(BeNil())

	Expect(template).To(Equal("5553686208582 The ID was 5553686208582"))
}

func toInterfaceSlice(arguments []string) []interface{} {
	argumentsArray := make([]interface{}, len(arguments))

	for i, s := range arguments {
		argumentsArray[i] = s
	}
	return argumentsArray
}

func ApplyTemplate(requestDetails *models.RequestDetails, state map[string]string, responseBody string) (string, error) {
	templator := templating.NewTemplator()
	dataSource1, _ := templating.NewCsvDataSource("test-csv1", "id,name,marks\n1,Test1,55\n2,Test2,56\n*,Dummy,ABSENT")
	dataSource2, _ := templating.NewCsvDataSource("test-csv2", "id,name,marks\n1,Test1,55\n2,Test2,56\n5553686208582,Test3,66\n")
	templator.TemplateHelper.TemplateDataSource.SetDataSource("test-csv1", dataSource1)
	templator.TemplateHelper.TemplateDataSource.SetDataSource("test-csv2", dataSource2)

	template, _ := templator.ParseTemplate(responseBody)
	return templator.RenderTemplate(template, requestDetails, &models.Literals{}, &models.Variables{}, state, &journal.Journal{})
}

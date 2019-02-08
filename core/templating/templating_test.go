package templating_test

import (
	"testing"

	"time"

	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/SpectoLabs/hoverfly/core/templating"
	. "github.com/onsi/gomega"
)

func Test_ShouldCreateTemplatingDataPathsFromRequest(t *testing.T) {
	RegisterTestingT(t)

	actual := templating.NewTemplatingDataFromRequest(
		&models.RequestDetails{
			Scheme:      "http",
			Destination: "test.com",
			Path:        "/foo/bar",
		},
		make(map[string]string),
	)

	Expect(actual.Request.Path).To(ConsistOf("foo", "bar"))
}

func Test_ShouldCreateTemplatingDataPathsFromRequestWithNoPaths(t *testing.T) {
	RegisterTestingT(t)

	actual := templating.NewTemplatingDataFromRequest(
		&models.RequestDetails{
			Scheme:      "http",
			Destination: "test.com",
		},
		make(map[string]string),
	)

	Expect(actual.Request.Path).To(BeEmpty())
}

func Test_ShouldCreateTemplatingDataQueryParamsFromRequest(t *testing.T) {
	RegisterTestingT(t)

	actual := templating.NewTemplatingDataFromRequest(
		&models.RequestDetails{
			Scheme:      "http",
			Destination: "test.com",
			Query: map[string][]string{
				"cheese": {"1", "3"},
				"ham":    {"2"},
			},
		},
		make(map[string]string),
	)

	Expect(actual.Request.QueryParam).To(HaveKeyWithValue("cheese", []string{"1", "3"}))
	Expect(actual.Request.QueryParam).To(HaveKeyWithValue("ham", []string{"2"}))
	Expect(actual.Request.QueryParam).To(HaveLen(2))
}

func Test_ShouldCreateTemplatingDataQueryParamsFromRequestWithNoQueryParams(t *testing.T) {
	RegisterTestingT(t)

	actual := templating.NewTemplatingDataFromRequest(&models.RequestDetails{
		Scheme:      "http",
		Destination: "test.com",
	},
		make(map[string]string),
	)

	Expect(actual.Request.QueryParam).To(BeEmpty())
}

func Test_ShouldCreateTemplatingDataHttpScheme(t *testing.T) {
	RegisterTestingT(t)

	actual := templating.NewTemplatingDataFromRequest(&models.RequestDetails{
		Scheme:      "http",
		Destination: "test.com",
	},
		make(map[string]string),
	)

	Expect(actual.Request.Scheme).To(Equal("http"))
}

func TestApplyTemplateWithQueryParams(t *testing.T) {
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

Path param value:
All path param values:
Looping through path params:

State One: A
State Two: B`))
}

func TestTemplatingWithHelperMethodsForDates(t *testing.T) {
	RegisterTestingT(t)

	template, err := ApplyTemplate(&models.RequestDetails{}, make(map[string]string), `{{iso8601DateTime}}`)

	Expect(err).To(BeNil())

	Expect(template).To(Equal(time.Now().UTC().Format("2006-01-02T15:04:05Z07:00")))

	template, err = ApplyTemplate(&models.RequestDetails{
		Query: map[string][]string{
			"plusDays": {"2"},
		},
	}, make(map[string]string), `{{iso8601DateTimePlusDays Request.QueryParam.plusDays}}`)

	Expect(err).To(BeNil())

	Expect(template).To(Equal(time.Now().AddDate(0, 0, 2).UTC().Format("2006-01-02T15:04:05Z07:00")))
}

func Test_ApplyTemplate_currentDateTime(t *testing.T) {
	RegisterTestingT(t)

	template, err := ApplyTemplate(&models.RequestDetails{}, make(map[string]string), `{{currentDateTime "2006-01-02T15:04:05Z07:00"}}`)

	Expect(err).To(BeNil())

	Expect(template).To(Not(Equal(ContainSubstring(`{{currentDateTime "2006-01-02T15:04:05Z07:00"}}`))))
}

func Test_ApplyTemplate_currentDateTimeAdd(t *testing.T) {
	RegisterTestingT(t)

	template, err := ApplyTemplate(&models.RequestDetails{}, make(map[string]string), `{{currentDateTimeAdd "5m" "2006-01-02T15:04:05Z07:00"}}`)

	Expect(err).To(BeNil())

	Expect(template).To(Not(Equal(ContainSubstring(`{{currentDateTimeAdd "5m" "2006-01-02T15:04:05Z07:00"}}`))))
}

func Test_ApplyTemplate_currentDateTimeSubtract(t *testing.T) {
	RegisterTestingT(t)

	template, err := ApplyTemplate(&models.RequestDetails{}, make(map[string]string), `{{currentDateTimeSubtract "5m" "2006-01-02T15:04:05Z07:00"}}`)

	Expect(err).To(BeNil())

	Expect(template).To(Not(Equal(ContainSubstring(`{{currentDateTimeSubtract "5m" "2006-01-02T15:04:05Z07:00"}}`))))
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

func Test_ApplyTemplate_Request_Body(t *testing.T) {
	RegisterTestingT(t)

	template, err := ApplyTemplate(&models.RequestDetails{}, make(map[string]string), `{{ Request.Body jsonPath "$.test" }}`)

	Expect(err).To(BeNil())

	Expect(template).To(Not(Equal(ContainSubstring(`{{Request.Body jsonPath \"$.test\"}}`))))
}

func ApplyTemplate(requestDetails *models.RequestDetails, state map[string]string, responseBody string) (string, error) {
	templator := templating.NewTemplator()
	template, _ := templator.ParseTemplate(responseBody)

	return templator.RenderTemplate(template, requestDetails, state)
}
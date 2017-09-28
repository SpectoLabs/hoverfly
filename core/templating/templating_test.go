package templating_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/SpectoLabs/hoverfly/core/templating"
	. "github.com/onsi/gomega"
	"time"
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

	template, err := templating.NewTemplator().ApplyTemplate(requestDetails,
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

	template, err := templating.NewTemplator().ApplyTemplate(requestDetails,
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

	template, err := templating.NewTemplator().ApplyTemplate(&models.RequestDetails{}, make(map[string]string), `{{iso8601DateTime}}`)

	Expect(err).To(BeNil())

	Expect(template).To(Equal(time.Now().UTC().Format("2006-01-02T15:04:05Z07:00")))

	template, err = templating.NewTemplator().ApplyTemplate(&models.RequestDetails{
		Query: map[string][]string{
			"plusDays": {"2"},
		},
	}, make(map[string]string), `{{iso8601DateTimePlusDays Request.QueryParam.plusDays}}`)

	Expect(err).To(BeNil())

	Expect(template).To(Equal(time.Now().AddDate(0, 0, 2).UTC().Format("2006-01-02T15:04:05Z07:00")))
}

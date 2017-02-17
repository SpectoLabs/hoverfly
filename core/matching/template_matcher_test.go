package matching

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/models"
	. "github.com/SpectoLabs/hoverfly/core/util"
	. "github.com/onsi/gomega"
)

var (
	unit = TemplateMatcher{}
)

func Test_Match_EmptyTemplateShouldMatchOnAnyRequest(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.Templates = append(simulation.Templates, models.RequestTemplateResponsePair{
		RequestTemplate: models.RequestTemplate{},
		Response: models.ResponseDetails{
			Body: "test-body",
		},
	})

	r := models.RequestDetails{
		Method:      "GET",
		Destination: "somehost.com",
		Headers: map[string][]string{
			"sdv": []string{"ascd"},
		},
	}
	result, _ := unit.Match(r, false, simulation)

	Expect(result.Body).To(Equal("test-body"))
}

func Test_Match_TemplateShouldMatchOnBody(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.Templates = append(simulation.Templates, models.RequestTemplateResponsePair{
		RequestTemplate: models.RequestTemplate{
			Body: StringToPointer("body"),
		},
		Response: models.ResponseDetails{
			Body: "body",
		},
	})

	r := models.RequestDetails{
		Body: "body",
	}
	result, err := unit.Match(r, false, simulation)
	Expect(err).To(BeNil())

	Expect(result.Body).To(Equal("body"))
}

func Test_Match_ReturnResponseWhenAllHeadersMatch(t *testing.T) {
	RegisterTestingT(t)

	headers := map[string][]string{
		"header1": []string{"val1"},
		"header2": []string{"val2"},
	}

	simulation := models.NewSimulation()

	simulation.Templates = append(simulation.Templates, models.RequestTemplateResponsePair{
		RequestTemplate: models.RequestTemplate{
			Headers: headers,
		},
		Response: models.ResponseDetails{
			Body: "test-body",
		},
	})

	r := models.RequestDetails{
		Method:      "GET",
		Destination: "http://somehost.com",
		Headers: map[string][]string{
			"header1": []string{"val1"},
			"header2": []string{"val2"},
		},
	}

	result, _ := unit.Match(r, false, simulation)

	Expect(result.Body).To(Equal("test-body"))
}

func Test_Match_ReturnNilWhenOneHeaderNotPresentInRequest(t *testing.T) {
	RegisterTestingT(t)

	headers := map[string][]string{
		"header1": []string{"val1"},
		"header2": []string{"val2"},
	}

	simulation := models.NewSimulation()

	simulation.Templates = append(simulation.Templates, models.RequestTemplateResponsePair{
		RequestTemplate: models.RequestTemplate{
			Headers: headers,
		},
		Response: models.ResponseDetails{
			Body: "test-body",
		},
	})

	r := models.RequestDetails{
		Method:      "GET",
		Destination: "http://somehost.com",
		Headers: map[string][]string{
			"header1": []string{"val1"},
		},
	}

	result, _ := unit.Match(r, false, simulation)

	Expect(result).To(BeNil())
}

func Test_Match_ReturnNilWhenOneHeaderValueDifferent(t *testing.T) {
	RegisterTestingT(t)

	headers := map[string][]string{
		"header1": []string{"val1"},
		"header2": []string{"val2"},
	}

	simulation := models.NewSimulation()

	simulation.Templates = append(simulation.Templates, models.RequestTemplateResponsePair{
		RequestTemplate: models.RequestTemplate{
			Headers: headers,
		},
		Response: models.ResponseDetails{
			Body: "test-body",
		},
	})

	r := models.RequestDetails{
		Method:      "GET",
		Destination: "somehost.com",
		Headers: map[string][]string{
			"header1": []string{"val1"},
			"header2": []string{"different"},
		},
	}
	result, _ := unit.Match(r, false, simulation)

	Expect(result).To(BeNil())
}

func Test_Match_ReturnResponseWithMultiValuedHeaderMatch(t *testing.T) {
	RegisterTestingT(t)

	headers := map[string][]string{
		"header1": []string{"val1-a", "val1-b"},
		"header2": []string{"val2"},
	}

	simulation := models.NewSimulation()

	simulation.Templates = append(simulation.Templates, models.RequestTemplateResponsePair{
		RequestTemplate: models.RequestTemplate{
			Headers: headers,
		},
		Response: models.ResponseDetails{
			Body: "test-body",
		},
	})

	r := models.RequestDetails{
		Method:      "GET",
		Destination: "http://somehost.com",
		Body:        "test-body",
		Headers: map[string][]string{
			"header1": []string{"val1-a", "val1-b"},
			"header2": []string{"val2"},
		},
	}
	result, _ := unit.Match(r, false, simulation)

	Expect(result.Body).To(Equal("test-body"))
}

func Test_Match_ReturnNilWithDifferentMultiValuedHeaders(t *testing.T) {
	RegisterTestingT(t)

	headers := map[string][]string{
		"header1": []string{"val1-a", "val1-b"},
		"header2": []string{"val2"},
	}

	simulation := models.NewSimulation()

	simulation.Templates = append(simulation.Templates, models.RequestTemplateResponsePair{
		RequestTemplate: models.RequestTemplate{
			Headers: headers,
		},
		Response: models.ResponseDetails{
			Body: "test-body",
		},
	})

	r := models.RequestDetails{
		Method:      "GET",
		Destination: "http://somehost.com",
		Headers: map[string][]string{
			"header1": []string{"val1-a", "val1-differnet"},
			"header2": []string{"val2"},
		},
	}

	result, _ := unit.Match(r, false, simulation)

	Expect(result).To(BeNil())
}

func Test_Match_EndpointMatchWithHeaders(t *testing.T) {
	RegisterTestingT(t)

	headers := map[string][]string{
		"header1": []string{"val1-a", "val1-b"},
		"header2": []string{"val2"},
	}

	destination := "testhost.com"
	method := "GET"
	path := "/a/1"
	query := "q=test"

	simulation := models.NewSimulation()

	simulation.Templates = append(simulation.Templates, models.RequestTemplateResponsePair{
		RequestTemplate: models.RequestTemplate{
			Headers:     headers,
			Destination: &destination,
			Path:        &path,
			Method:      &method,
			Query:       &query,
		},
		Response: models.ResponseDetails{
			Body: "test-body",
		},
	})

	r := models.RequestDetails{
		Method:      "GET",
		Destination: "testhost.com",
		Path:        "/a/1",
		Query:       "q=test",
		Headers: map[string][]string{
			"header1": []string{"val1-a", "val1-b"},
			"header2": []string{"val2"},
		},
	}
	result, _ := unit.Match(r, false, simulation)

	Expect(result.Body).To(Equal("test-body"))
}

func Test_Match_EndpointMismatchWithHeadersReturnsNil(t *testing.T) {
	RegisterTestingT(t)

	headers := map[string][]string{
		"header1": []string{"val1-a", "val1-b"},
		"header2": []string{"val2"},
	}

	destination := "testhost.com"
	method := "GET"
	path := "/a/1"
	query := "q=test"

	simulation := models.NewSimulation()

	simulation.Templates = append(simulation.Templates, models.RequestTemplateResponsePair{
		RequestTemplate: models.RequestTemplate{
			Headers:     headers,
			Destination: &destination,
			Path:        &path,
			Method:      &method,
			Query:       &query,
		},
		Response: models.ResponseDetails{
			Body: "test-body",
		},
	})

	r := models.RequestDetails{
		Method:      "GET",
		Destination: "http://testhost.com",
		Path:        "/a/1",
		Query:       "q=different",
		Headers: map[string][]string{
			"header1": []string{"val1-a", "val1-b"},
			"header2": []string{"val2"},
		},
	}

	result, _ := unit.Match(r, false, simulation)

	Expect(result).To(BeNil())
}

func Test_Match_AbleToMatchAnEmptyPathInAReasonableWay(t *testing.T) {
	RegisterTestingT(t)

	response := models.ResponseDetails{
		Body: "test-body",
	}
	destination := "testhost.com"
	method := "GET"
	path := ""
	query := "q=test"
	simulation := models.NewSimulation()

	simulation.Templates = append(simulation.Templates, models.RequestTemplateResponsePair{
		RequestTemplate: models.RequestTemplate{
			Destination: &destination,
			Path:        &path,
			Method:      &method,
			Query:       &query,
		},
		Response: response,
	})

	r := models.RequestDetails{
		Method:      "GET",
		Destination: "testhost.com",
		Query:       "q=test",
	}
	result, _ := unit.Match(r, false, simulation)

	Expect(result.Body).To(Equal("test-body"))

	r = models.RequestDetails{
		Method:      "GET",
		Destination: "testhost.com",
		Path:        "/a/1",
		Query:       "q=test",
	}

	result, _ = unit.Match(r, false, simulation)

	Expect(result).To(BeNil())
}

func Test_Match_RequestTemplateResponsePairCanBeConvertedToARequestResponsePairView_WhileIncomplete(t *testing.T) {
	RegisterTestingT(t)

	method := "POST"

	requestTemplateResponsePair := models.RequestTemplateResponsePair{
		RequestTemplate: models.RequestTemplate{
			Method: &method,
		},
		Response: models.ResponseDetails{
			Body: "template matched",
		},
	}

	pairView := requestTemplateResponsePair.ConvertToV1RequestResponsePairView()

	Expect(pairView.Request.RequestType).To(Equal(StringToPointer("template")))
	Expect(pairView.Request.Method).To(Equal(StringToPointer("POST")))
	Expect(pairView.Request.Destination).To(BeNil())
	Expect(pairView.Request.Path).To(BeNil())
	Expect(pairView.Request.Scheme).To(BeNil())
	Expect(pairView.Request.Query).To(BeNil())

	Expect(pairView.Response.Body).To(Equal("template matched"))
}

func Test_Match_TemplatesCanUseGlobsOnDestinationAndBeMatched(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.Templates = append(simulation.Templates, models.RequestTemplateResponsePair{
		RequestTemplate: models.RequestTemplate{
			Destination: StringToPointer("*.com"),
		},
		Response: models.ResponseDetails{
			Body: "template matched",
		},
	})

	request := models.RequestDetails{
		Method:      "GET",
		Destination: "testhost.com",
		Path:        "/api/1",
	}

	response, err := unit.Match(request, false, simulation)
	Expect(err).To(BeNil())

	Expect(response.Body).To(Equal("template matched"))
}

func Test_Match_TemplatesCanUseGlobsOnPathAndBeMatched(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.Templates = append(simulation.Templates, models.RequestTemplateResponsePair{
		RequestTemplate: models.RequestTemplate{
			Path: StringToPointer("/api/*"),
		},
		Response: models.ResponseDetails{
			Body: "template matched",
		},
	})

	request := models.RequestDetails{
		Method:      "GET",
		Destination: "testhost.com",
		Path:        "/api/1",
	}

	response, err := unit.Match(request, false, simulation)
	Expect(err).To(BeNil())

	Expect(response.Body).To(Equal("template matched"))
}

func Test_Match_TemplatesCanUseGlobsOnMethodAndBeMatched(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.Templates = append(simulation.Templates, models.RequestTemplateResponsePair{
		RequestTemplate: models.RequestTemplate{
			Method: StringToPointer("*T"),
		},
		Response: models.ResponseDetails{
			Body: "template matched",
		},
	})

	request := models.RequestDetails{
		Method:      "GET",
		Destination: "testhost.com",
		Path:        "/api/1",
	}

	response, err := unit.Match(request, false, simulation)
	Expect(err).To(BeNil())

	Expect(response.Body).To(Equal("template matched"))
}

func Test_Match_TemplatesCanUseGlobsOnSchemeAndBeMatched(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.Templates = append(simulation.Templates, models.RequestTemplateResponsePair{
		RequestTemplate: models.RequestTemplate{
			Scheme: StringToPointer("H*"),
		},
		Response: models.ResponseDetails{
			Body: "template matched",
		},
	})

	request := models.RequestDetails{
		Method:      "GET",
		Destination: "testhost.com",
		Scheme:      "http",
		Path:        "/api/1",
	}

	response, err := unit.Match(request, false, simulation)
	Expect(err).To(BeNil())

	Expect(response.Body).To(Equal("template matched"))
}

func Test_Match_TemplatesCanUseGlobsOnQueryAndBeMatched(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.Templates = append(simulation.Templates, models.RequestTemplateResponsePair{
		RequestTemplate: models.RequestTemplate{
			Query: StringToPointer("q=*"),
		},
		Response: models.ResponseDetails{
			Body: "template matched",
		},
	})

	request := models.RequestDetails{
		Method:      "GET",
		Destination: "testhost.com",
		Path:        "/api/1",
		Query:       "q=anything-i-want",
	}

	response, err := unit.Match(request, false, simulation)
	Expect(err).To(BeNil())

	Expect(response.Body).To(Equal("template matched"))
}

func Test_Match_TemplatesCanUseGlobsOnBodyndBeMatched(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.Templates = append(simulation.Templates, models.RequestTemplateResponsePair{
		RequestTemplate: models.RequestTemplate{
			Body: StringToPointer(`{"json": "object", "key": *}`),
		},
		Response: models.ResponseDetails{
			Body: "template matched",
		},
	})

	request := models.RequestDetails{
		Method:      "GET",
		Destination: "testhost.com",
		Path:        "/api/1",
		Body:        `{"json": "object", "key": "value"}`,
	}

	response, err := unit.Match(request, false, simulation)
	Expect(err).To(BeNil())

	Expect(response.Body).To(Equal("template matched"))
}

func Test_Match_TemplatesCanUseGlobsOnBodyAndNotMatchWhenTheBodyIsWrong(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.Templates = append(simulation.Templates, models.RequestTemplateResponsePair{
		RequestTemplate: models.RequestTemplate{
			Body: StringToPointer(`{"json": "object", "key": *}`),
		},
		Response: models.ResponseDetails{
			Body: "template matched",
		},
	})

	request := models.RequestDetails{
		Method:      "GET",
		Destination: "testhost.com",
		Path:        "/api/1",
		Body:        `[{"json": "objects", "key": "value"}]`,
	}

	_, err := unit.Match(request, false, simulation)
	Expect(err).ToNot(BeNil())
}

func Test_Match_TemplatesCanUseGlobsOnHeadersAndBeMatched(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.Templates = append(simulation.Templates, models.RequestTemplateResponsePair{
		RequestTemplate: models.RequestTemplate{
			Headers: map[string][]string{
				"unique-header": []string{"*"},
			},
		},
		Response: models.ResponseDetails{
			Body: "template matched",
		},
	})

	request := models.RequestDetails{
		Method:      "GET",
		Destination: "testhost.com",
		Path:        "/api/1",
		Headers: map[string][]string{
			"unique-header": []string{"totally-unique"},
		},
	}

	response, err := unit.Match(request, false, simulation)
	Expect(err).To(BeNil())

	Expect(response.Body).To(Equal("template matched"))
}

func Test_Match_RequestTemplateResponsePair_ConvertToRequestResponsePairView_CanBeConvertedToARequestResponsePairView_WhileIncomplete(t *testing.T) {
	RegisterTestingT(t)

	method := "POST"

	requestTemplateResponsePair := models.RequestTemplateResponsePair{
		RequestTemplate: models.RequestTemplate{
			Method: &method,
		},
		Response: models.ResponseDetails{
			Body: "template matched",
		},
	}

	pairView := requestTemplateResponsePair.ConvertToRequestResponsePairView()

	Expect(pairView.Request.RequestType).To(Equal(StringToPointer("template")))
	Expect(pairView.Request.Method).To(Equal(StringToPointer("POST")))
	Expect(pairView.Request.Destination).To(BeNil())
	Expect(pairView.Request.Path).To(BeNil())
	Expect(pairView.Request.Scheme).To(BeNil())
	Expect(pairView.Request.Query).To(BeNil())

	Expect(pairView.Response.Body).To(Equal("template matched"))
}

func Test_headerMatch(t *testing.T) {
	RegisterTestingT(t)

	tmplHeaders := map[string][]string{
		"header1": []string{"val1"},
		"header2": []string{"val2"},
	}

	res := headerMatch(tmplHeaders, tmplHeaders)
	Expect(res).To(BeTrue())
}

func Test_headerMatch_IgnoreTestCaseInsensitive(t *testing.T) {
	RegisterTestingT(t)

	tmplHeaders := map[string][]string{
		"header1": []string{"val1"},
		"header2": []string{"val2"},
	}
	reqHeaders := map[string][]string{
		"HEADER1": []string{"val1"},
		"Header2": []string{"VAL2"},
	}
	res := headerMatch(tmplHeaders, reqHeaders)
	Expect(res).To(BeTrue())
}

func Test_HeaderMatch_MatchingTemplateHasMoreHeaderKeysThanRequestMatchesFalse(t *testing.T) {
	RegisterTestingT(t)

	tmplHeaders := map[string][]string{
		"header1": []string{"val1"},
		"header2": []string{"val2"},
	}
	reqHeaders := map[string][]string{
		"header1": []string{"val1"},
	}
	res := headerMatch(tmplHeaders, reqHeaders)
	Expect(res).To(BeFalse())
}

func Test_headerMatch_MatchingTemplateHasMoreHeaderValuesThanRequestMatchesFalse(t *testing.T) {
	RegisterTestingT(t)

	tmplHeaders := map[string][]string{
		"header2": []string{"val1", "val2"},
	}
	reqHeaders := map[string][]string{
		"header2": []string{"val1"},
	}
	res := headerMatch(tmplHeaders, reqHeaders)
	Expect(res).To(BeFalse())
}

func Test_headerMatch_MatchingRequestHasMoreHeaderKeysThanTemplateMatchesFalse(t *testing.T) {
	RegisterTestingT(t)

	tmplHeaders := map[string][]string{
		"header2": []string{"val2"},
	}
	reqHeaders := map[string][]string{
		"HEADER1": []string{"val1"},
		"header2": []string{"val2"},
	}
	res := headerMatch(tmplHeaders, reqHeaders)
	Expect(res).To(BeTrue())
}

func Test_headerMatch_MatchingRequestHasMoreHeaderValuesThanTemplateMatchesFalse(t *testing.T) {
	RegisterTestingT(t)

	tmplHeaders := map[string][]string{
		"header2": []string{"val2"},
	}
	reqHeaders := map[string][]string{
		"header2": []string{"val1", "val2"},
	}
	res := headerMatch(tmplHeaders, reqHeaders)
	Expect(res).To(BeTrue())
}

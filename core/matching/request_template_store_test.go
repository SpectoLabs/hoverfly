package matching

import (
	"github.com/SpectoLabs/hoverfly/core/models"
	. "github.com/onsi/gomega"
	"testing"
	. "github.com/SpectoLabs/hoverfly/core/util"
)

func TestEmptyTemplateShouldMatchOnAnyRequest(t *testing.T) {
	RegisterTestingT(t)

	response := models.ResponseDetails{
		Body: "test-body",
	}
	templateEntry := RequestTemplateResponsePair{
		RequestTemplate: RequestTemplate{},
		Response:        response,
	}
	store := RequestTemplateStore{templateEntry}

	r := models.RequestDetails{
		Method:      "GET",
		Destination: "somehost.com",
		Headers: map[string][]string{
			"sdv": []string{"ascd"},
		},
	}
	result, _ := store.GetResponse(r, false)

	Expect(result.Body).To(Equal("test-body"))
}

func TestTemplateShouldMatchOnBody(t *testing.T) {
	RegisterTestingT(t)

	response := models.ResponseDetails{
		Body: "body",
	}
	templateEntry := RequestTemplateResponsePair{
		RequestTemplate: RequestTemplate{Body: &response.Body},
		Response:        response,
	}
	store := RequestTemplateStore{templateEntry}

	r := models.RequestDetails{
		Body: "body",
	}
	result, err := store.GetResponse(r, false)
	Expect(err).To(BeNil())

	Expect(result.Body).To(Equal("body"))
}

func TestReturnResponseWhenAllHeadersMatch(t *testing.T) {
	RegisterTestingT(t)

	response := models.ResponseDetails{
		Body: "test-body",
	}
	headers := map[string][]string{
		"header1": []string{"val1"},
		"header2": []string{"val2"},
	}
	templateEntry := RequestTemplateResponsePair{
		RequestTemplate: RequestTemplate{
			Headers: headers,
		},
		Response: response,
	}
	store := RequestTemplateStore{templateEntry}

	r := models.RequestDetails{
		Method:      "GET",
		Destination: "http://somehost.com",
		Headers: map[string][]string{
			"header1": []string{"val1"},
			"header2": []string{"val2"},
		},
	}

	result, _ := store.GetResponse(r, false)

	Expect(result.Body).To(Equal("test-body"))
}

func TestReturnNilWhenOneHeaderNotPresentInRequest(t *testing.T) {
	RegisterTestingT(t)

	response := models.ResponseDetails{
		Body: "test-body",
	}
	headers := map[string][]string{
		"header1": []string{"val1"},
		"header2": []string{"val2"},
	}
	templateEntry := RequestTemplateResponsePair{
		RequestTemplate: RequestTemplate{
			Headers: headers,
		},
		Response: response,
	}
	store := RequestTemplateStore{templateEntry}

	r := models.RequestDetails{
		Method:      "GET",
		Destination: "http://somehost.com",
		Headers: map[string][]string{
			"header1": []string{"val1"},
		},
	}

	result, _ := store.GetResponse(r, false)

	Expect(result).To(BeNil())
}

func TestReturnNilWhenOneHeaderValueDifferent(t *testing.T) {
	RegisterTestingT(t)

	response := models.ResponseDetails{
		Body: "test-body",
	}
	headers := map[string][]string{
		"header1": []string{"val1"},
		"header2": []string{"val2"},
	}
	templateEntry := RequestTemplateResponsePair{
		RequestTemplate: RequestTemplate{
			Headers: headers,
		},
		Response: response,
	}
	store := RequestTemplateStore{templateEntry}

	r := models.RequestDetails{
		Method:      "GET",
		Destination: "somehost.com",
		Headers: map[string][]string{
			"header1": []string{"val1"},
			"header2": []string{"different"},
		},
	}
	result, _ := store.GetResponse(r, false)

	Expect(result).To(BeNil())
}

func TestReturnResponseWithMultiValuedHeaderMatch(t *testing.T) {
	RegisterTestingT(t)

	response := models.ResponseDetails{
		Body: "test-body",
	}
	headers := map[string][]string{
		"header1": []string{"val1-a", "val1-b"},
		"header2": []string{"val2"},
	}
	templateEntry := RequestTemplateResponsePair{
		RequestTemplate: RequestTemplate{
			Headers: headers,
		},
		Response: response,
	}
	store := RequestTemplateStore{templateEntry}

	r := models.RequestDetails{
		Method:      "GET",
		Destination: "http://somehost.com",
		Body:        "test-body",
		Headers: map[string][]string{
			"header1": []string{"val1-a", "val1-b"},
			"header2": []string{"val2"},
		},
	}
	result, _ := store.GetResponse(r, false)

	Expect(result.Body).To(Equal("test-body"))
}

func TestReturnNilWithDifferentMultiValuedHeaders(t *testing.T) {
	RegisterTestingT(t)

	response := models.ResponseDetails{
		Body: "test-body",
	}
	headers := map[string][]string{
		"header1": []string{"val1-a", "val1-b"},
		"header2": []string{"val2"},
	}
	templateEntry := RequestTemplateResponsePair{
		RequestTemplate: RequestTemplate{
			Headers: headers,
		},
		Response: response,
	}

	store := RequestTemplateStore{templateEntry}

	r := models.RequestDetails{
		Method:      "GET",
		Destination: "http://somehost.com",
		Headers: map[string][]string{
			"header1": []string{"val1-a", "val1-differnet"},
			"header2": []string{"val2"},
		},
	}

	result, _ := store.GetResponse(r, false)

	Expect(result).To(BeNil())
}

func TestHeaderMatch(t *testing.T) {
	RegisterTestingT(t)

	tmplHeaders := map[string][]string{
		"header1": []string{"val1"},
		"header2": []string{"val2"},
	}

	res := headerMatch(tmplHeaders, tmplHeaders)
	Expect(res).To(BeTrue())
}

func IgnoreTestCaseInsensitiveHeaderMatch(t *testing.T) {
	RegisterTestingT(t)

	tmplHeaders := map[string][]string{
		"header1": []string{"val1"},
		"header2": []string{"val2"},
	}
	reqHeaders := map[string][]string{
		"HEADER1": []string{"val1"},
		"Header2": []string{"val2"},
	}
	res := headerMatch(tmplHeaders, reqHeaders)
	Expect(res).To(BeTrue())
}

func TestEndpointMatchWithHeaders(t *testing.T) {
	RegisterTestingT(t)

	response := models.ResponseDetails{
		Body: "test-body",
	}
	headers := map[string][]string{
		"header1": []string{"val1-a", "val1-b"},
		"header2": []string{"val2"},
	}
	destination := "testhost.com"
	method := "GET"
	path := "/a/1"
	query := "q=test"
	templateEntry := RequestTemplateResponsePair{
		RequestTemplate: RequestTemplate{
			Headers:     headers,
			Destination: &destination,
			Path:        &path,
			Method:      &method,
			Query:       &query,
		},
		Response: response,
	}
	store := RequestTemplateStore{templateEntry}

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
	result, _ := store.GetResponse(r, false)

	Expect(result.Body).To(Equal("test-body"))
}

func TestEndpointMismatchWithHeadersReturnsNil(t *testing.T) {
	RegisterTestingT(t)

	response := models.ResponseDetails{
		Body: "test-body",
	}
	headers := map[string][]string{
		"header1": []string{"val1-a", "val1-b"},
		"header2": []string{"val2"},
	}
	destination := "testhost.com"
	method := "GET"
	path := "/a/1"
	query := "q=test"
	templateEntry := RequestTemplateResponsePair{
		RequestTemplate: RequestTemplate{
			Headers:     headers,
			Destination: &destination,
			Path:        &path,
			Method:      &method,
			Query:       &query,
		},
		Response: response,
	}
	store := RequestTemplateStore{templateEntry}

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

	result, _ := store.GetResponse(r, false)

	Expect(result).To(BeNil())
}

func TestAbleToMatchAnEmptyPathInAReasonableWay(t *testing.T) {
	RegisterTestingT(t)

	response := models.ResponseDetails{
		Body: "test-body",
	}
	destination := "testhost.com"
	method := "GET"
	path := ""
	query := "q=test"
	templateEntry := RequestTemplateResponsePair{
		RequestTemplate: RequestTemplate{
			Destination: &destination,
			Path:        &path,
			Method:      &method,
			Query:       &query,
		},
		Response: response,
	}
	store := RequestTemplateStore{templateEntry}

	r := models.RequestDetails{
		Method:      "GET",
		Destination: "testhost.com",
		Query:       "q=test",
	}
	result, _ := store.GetResponse(r, false)

	Expect(result.Body).To(Equal("test-body"))

	r = models.RequestDetails{
		Method:      "GET",
		Destination: "testhost.com",
		Path:        "/a/1",
		Query:       "q=test",
	}

	result, _ = store.GetResponse(r, false)

	Expect(result).To(BeNil())
}

func TestRequestTemplateResponsePairCanBeConvertedToARequestResponsePairView_WhileIncomplete(t *testing.T) {
	RegisterTestingT(t)

	method := "POST"

	requestTemplateResponsePair := RequestTemplateResponsePair{
		RequestTemplate: RequestTemplate{
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

func TestTemplatesCanUseGlobsOnDestinationAndBeMatched(t *testing.T) {
	RegisterTestingT(t)

	requestTemplateResponsePair := RequestTemplateResponsePair{
		RequestTemplate: RequestTemplate{
			Destination: StringToPointer("*.com"),
		},
		Response: models.ResponseDetails{
			Body: "template matched",
		},
	}

	store := RequestTemplateStore{requestTemplateResponsePair}

	request := models.RequestDetails{
		Method:      "GET",
		Destination: "testhost.com",
		Path:        "/api/1",
	}

	response, err :=store.GetResponse(request, false)
	Expect(err).To(BeNil())

	Expect(response.Body).To(Equal("template matched"))
}

func TestTemplatesCanUseGlobsOnPathAndBeMatched(t *testing.T) {
	RegisterTestingT(t)

	requestTemplateResponsePair := RequestTemplateResponsePair{
		RequestTemplate: RequestTemplate{
			Path: StringToPointer("/api/*"),
		},
		Response: models.ResponseDetails{
			Body: "template matched",
		},
	}

	store := RequestTemplateStore{requestTemplateResponsePair}

	request := models.RequestDetails{
		Method:      "GET",
		Destination: "testhost.com",
		Path:        "/api/1",
	}

	response, err :=store.GetResponse(request, false)
	Expect(err).To(BeNil())

	Expect(response.Body).To(Equal("template matched"))
}

func TestTemplatesCanUseGlobsOnMethodAndBeMatched(t *testing.T) {
	RegisterTestingT(t)

	requestTemplateResponsePair := RequestTemplateResponsePair{
		RequestTemplate: RequestTemplate{
			Method: StringToPointer("*T"),
		},
		Response: models.ResponseDetails{
			Body: "template matched",
		},
	}

	store := RequestTemplateStore{requestTemplateResponsePair}

	request := models.RequestDetails{
		Method:      "GET",
		Destination: "testhost.com",
		Path:        "/api/1",
	}

	response, err :=store.GetResponse(request, false)
	Expect(err).To(BeNil())

	Expect(response.Body).To(Equal("template matched"))
}

func TestTemplatesCanUseGlobsOnSchemeAndBeMatched(t *testing.T) {
	RegisterTestingT(t)

	requestTemplateResponsePair := RequestTemplateResponsePair{
		RequestTemplate: RequestTemplate{
			Scheme: StringToPointer("H*"),
		},
		Response: models.ResponseDetails{
			Body: "template matched",
		},
	}

	store := RequestTemplateStore{requestTemplateResponsePair}

	request := models.RequestDetails{
		Method:      "GET",
		Destination: "testhost.com",
		Scheme:      "http",
		Path:        "/api/1",
	}

	response, err :=store.GetResponse(request, false)
	Expect(err).To(BeNil())

	Expect(response.Body).To(Equal("template matched"))
}

func TestTemplatesCanUseGlobsOnQueryAndBeMatched(t *testing.T) {
	RegisterTestingT(t)

	requestTemplateResponsePair := RequestTemplateResponsePair{
		RequestTemplate: RequestTemplate{
			Query: StringToPointer("q=*"),
		},
		Response: models.ResponseDetails{
			Body: "template matched",
		},
	}

	store := RequestTemplateStore{requestTemplateResponsePair}

	request := models.RequestDetails{
		Method:      "GET",
		Destination: "testhost.com",
		Path:        "/api/1",
		Query: 	     "q=anything-i-want",
	}

	response, err :=store.GetResponse(request, false)
	Expect(err).To(BeNil())

	Expect(response.Body).To(Equal("template matched"))
}

func TestTemplatesCanUseGlobsOnBodyndBeMatched(t *testing.T) {
	RegisterTestingT(t)

	requestTemplateResponsePair := RequestTemplateResponsePair{
		RequestTemplate: RequestTemplate{
			Body: StringToPointer(`{"json": "object", "key": *}`),
		},
		Response: models.ResponseDetails{
			Body: "template matched",
		},
	}

	store := RequestTemplateStore{requestTemplateResponsePair}

	request := models.RequestDetails{
		Method:      "GET",
		Destination: "testhost.com",
		Path:        "/api/1",
		Body: 	     `{"json": "object", "key": "value"}`,
	}

	response, err :=store.GetResponse(request, false)
	Expect(err).To(BeNil())

	Expect(response.Body).To(Equal("template matched"))
}

func TestTemplatesCanUseGlobsOnBodyAndNotMatchWhenTheBodyIsWrong(t *testing.T) {
	RegisterTestingT(t)

	requestTemplateResponsePair := RequestTemplateResponsePair{
		RequestTemplate: RequestTemplate{
			Body: StringToPointer(`{"json": "object", "key": *}`),
		},
		Response: models.ResponseDetails{
			Body: "template matched",
		},
	}

	store := RequestTemplateStore{requestTemplateResponsePair}

	request := models.RequestDetails{
		Method:      "GET",
		Destination: "testhost.com",
		Path:        "/api/1",
		Body: 	     `[{"json": "objects", "key": "value"}]`,
	}

	_, err :=store.GetResponse(request, false)
	Expect(err).ToNot(BeNil())
}
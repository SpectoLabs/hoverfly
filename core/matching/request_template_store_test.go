package matching

import (
	"testing"
	"github.com/SpectoLabs/hoverfly/core/models"
	. "github.com/onsi/gomega"
)

func TestEmptyTemplateShouldMatchOnAnyRequest(t *testing.T) {
	RegisterTestingT(t)

	response := models.ResponseDetails{
		Body: "test-body",
	}
	templateEntry := RequestTemplatePayload{
		RequestTemplate: RequestTemplate{},
		Response: response,
	}
	store := RequestTemplateStore{templateEntry}

	r := models.RequestDetails{
		Method:"GET",
		Destination: "somehost.com",
		Headers: map[string][]string {
			"sdv": []string{"ascd"},
		},
	}
	result, _ := store.GetPayload(r, false)

	Expect(result.Response.Body).To(Equal("test-body"))
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
	templateEntry := RequestTemplatePayload{
		RequestTemplate: RequestTemplate{
			Headers: headers,
		},
		Response: response,
	}
	store := RequestTemplateStore{templateEntry}

	r := models.RequestDetails{
		Method:"GET",
		Destination: "http://somehost.com",
		Headers: map[string][]string {
			"header1": []string{"val1"},
			"header2": []string{"val2"},
		},
	}

	result, _ := store.GetPayload(r, false)

	Expect(result.Response.Body).To(Equal("test-body"))
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
	templateEntry := RequestTemplatePayload{
		RequestTemplate: RequestTemplate{
			Headers: headers,
		},
		Response: response,
	}
	store := RequestTemplateStore{templateEntry}

	r := models.RequestDetails{
		Method:"GET",
		Destination: "http://somehost.com",
		Headers: map[string][]string {
			"header1": []string{"val1"},
		},
	}

	result, _ := store.GetPayload(r, false)

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
	templateEntry := RequestTemplatePayload{
		RequestTemplate: RequestTemplate{
			Headers: headers,
		},
		Response: response,
	}
	store := RequestTemplateStore{templateEntry}

	r := models.RequestDetails{
		Method:"GET",
		Destination: "somehost.com",
		Headers: map[string][]string {
			"header1": []string{"val1"},
			"header2": []string{"different"},
		},
	}
	result, _ := store.GetPayload(r, false)

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
	templateEntry := RequestTemplatePayload{
		RequestTemplate: RequestTemplate{
			Headers: headers,
		},
		Response: response,
	}
	store := RequestTemplateStore{templateEntry}

	r := models.RequestDetails{
		Method:"GET",
		Destination: "http://somehost.com",
		Body: "test-body",
		Headers: map[string][]string {
			"header1": []string{"val1-a", "val1-b"},
			"header2": []string{"val2"},
		},
	}
	result, _ := store.GetPayload(r, false)

	Expect(result.Response.Body).To(Equal("test-body"))
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
	templateEntry := RequestTemplatePayload{
		RequestTemplate: RequestTemplate{
			Headers: headers,
		},
		Response: response,
	}

	store := RequestTemplateStore{templateEntry}

	r := models.RequestDetails{
		Method:"GET",
		Destination: "http://somehost.com",
		Headers: map[string][]string {
			"header1": []string{"val1-a", "val1-differnet"},
			"header2": []string{"val2"},
		},
	}

	result, _ := store.GetPayload(r, false)

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
	templateEntry := RequestTemplatePayload{
		RequestTemplate: RequestTemplate{
			Headers: headers,
			Destination: &destination,
			Path: &path,
			Method: &method,
			Query: &query,
		},
		Response: response,
	}
	store := RequestTemplateStore{templateEntry}

	r := models.RequestDetails{
		Method:"GET",
		Destination: "testhost.com",
		Path: "/a/1",
		Query: "q=test",
		Headers: map[string][]string {
			"header1": []string{"val1-a", "val1-b"},
			"header2": []string{"val2"},
		},
	}
	result, _ := store.GetPayload(r, false)

	Expect(result.Response.Body).To(Equal("test-body"))
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
	templateEntry := RequestTemplatePayload{
		RequestTemplate: RequestTemplate{
			Headers: headers,
			Destination: &destination,
			Path: &path,
			Method: &method,
			Query: &query,
		},
		Response: response,
	}
	store := RequestTemplateStore{templateEntry}

	r := models.RequestDetails{
		Method:"GET",
		Destination: "http://testhost.com",
		Path: "/a/1",
		Query: "q=different",
		Headers: map[string][]string {
			"header1": []string{"val1-a", "val1-b"},
			"header2": []string{"val2"},
		},
	}

	result, _ := store.GetPayload(r, false)

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
	templateEntry := RequestTemplatePayload{
		RequestTemplate: RequestTemplate{
			Destination: &destination,
			Path: &path,
			Method: &method,
			Query: &query,
		},
		Response: response,
	}
	store := RequestTemplateStore{templateEntry}

	r := models.RequestDetails{
		Method:"GET",
		Destination: "testhost.com",
		Query: "q=test",
	}
	result, _ := store.GetPayload(r, false)

	Expect(result.Response.Body).To(Equal("test-body"))

	r = models.RequestDetails{
		Method:"GET",
		Destination: "testhost.com",
		Path: "/a/1",
		Query: "q=test",
	}

	result, _ = store.GetPayload(r, false)

	Expect(result).To(BeNil())
}
package matching

import (
	"testing"
	"github.com/SpectoLabs/hoverfly/core/models"
	"net/http"
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

	r, _ := http.NewRequest("GET", "http://somehost.com", nil)
	r.Header = http.Header{
		"sdv": []string{"ascd"},
	}
	result, _ := store.GetPayload(r, []byte("test-body"))

	//var rd models.RequestDetails
	Expect(result.Response.Body).To(Equal("test-body"))
	//Expect(result.Request).To(BeFalse())
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

	r, _ := http.NewRequest("GET", "http://somehost.com", nil)
	r.Header = http.Header{
		"header1": []string{"val1"},
		"header2": []string{"val2"},
	}
	result, _ := store.GetPayload(r, []byte("test-body"))

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

	r, _ := http.NewRequest("GET", "http://somehost.com", nil)
	r.Header = http.Header{
		"header1": []string{"val1"},
	}
	result, _ := store.GetPayload(r, []byte("test-body"))

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

	r, _ := http.NewRequest("GET", "http://somehost.com", nil)
	r.Header = http.Header{
		"header1": []string{"val1"},
		"header2": []string{"different"},
	}
	result, _ := store.GetPayload(r, []byte("test-body"))

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

	r, _ := http.NewRequest("GET", "http://somehost.com", nil)
	r.Header = http.Header{
		"header1": []string{"val1-a", "val1-b"},
		"header2": []string{"val2"},
	}
	result, _ := store.GetPayload(r, []byte("test-body"))

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

	r, _ := http.NewRequest("GET", "http://somehost.com", nil)
	r.Header = http.Header{
		"header1": []string{"val1-a", "val1-differnet"},
		"header2": []string{"val2"},
	}
	result, _ := store.GetPayload(r, []byte("test-body"))

	Expect(result).To(BeNil())
}

func TestHeaderMatch(t *testing.T) {
	RegisterTestingT(t)

	tmplHeaders := map[string][]string{
		"header1": []string{"val1"},
		"header2": []string{"val2"},
	}
	reqHeaders := http.Header{
		"header1": []string{"val1"},
		"header2": []string{"val2"},
	}
	res := headerMatch(tmplHeaders, reqHeaders)
	Expect(res).To(BeTrue())
}

func IgnoreTestCaseInsensitiveHeaderMatch(t *testing.T) {
	RegisterTestingT(t)

	tmplHeaders := map[string][]string{
		"header1": []string{"val1"},
		"header2": []string{"val2"},
	}
	reqHeaders := http.Header{
		"HEADER1": []string{"val1"},
		"Header2": []string{"val2"},
	}
	res := headerMatch(tmplHeaders, reqHeaders)
	Expect(res).To(BeTrue())
}
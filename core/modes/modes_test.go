package modes_test

import (
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/SpectoLabs/hoverfly/core/modes"
	. "github.com/onsi/gomega"
)

func Test_ReconstructResponse_ReturnsAResponseWithCorrectStatus(t *testing.T) {
	RegisterTestingT(t)

	req, _ := http.NewRequest("GET", "http://example.com", nil)

	pair := models.RequestResponsePair{
		Response: models.ResponseDetails{
			Status: 404,
		},
	}

	response := modes.ReconstructResponse(req, pair)

	Expect(response.StatusCode).To(Equal(404))
}

func Test_ReconstructResponse_ReturnsAResponseWithBody(t *testing.T) {
	RegisterTestingT(t)

	req, _ := http.NewRequest("GET", "http://example.com", nil)

	pair := models.RequestResponsePair{
		Response: models.ResponseDetails{
			Body: "test body",
		},
	}

	response := modes.ReconstructResponse(req, pair)

	responseBody, err := ioutil.ReadAll(response.Body)
	Expect(err).To(BeNil())

	Expect(string(responseBody)).To(Equal("test body"))
}

func Test_ReconstructResponse_AddsHeadersToResponse(t *testing.T) {
	RegisterTestingT(t)

	req, _ := http.NewRequest("GET", "http://example.com", nil)

	pair := models.RequestResponsePair{}

	headers := make(map[string][]string)
	headers["Header"] = []string{"one"}

	pair.Response.Headers = headers

	response := modes.ReconstructResponse(req, pair)

	Expect(response.Header.Get("Header")).To(Equal(headers["Header"][0]))
}

func Test_ReconstructResponse_AddsMultipleHeaderValuesToResponse(t *testing.T) {
	RegisterTestingT(t)

	req, _ := http.NewRequest("GET", "http://example.com", nil)

	pair := models.RequestResponsePair{}

	headers := make(map[string][]string)
	headers["Header"] = []string{"one", "two", "three"}

	pair.Response.Headers = headers

	response := modes.ReconstructResponse(req, pair)
	values, ok := response.Header["Header"]
	Expect(ok).To(BeTrue())

	Expect(len(values)).To(Equal(3))
	Expect(values[0]).To(Equal("one"))
	Expect(values[1]).To(Equal("two"))
	Expect(values[2]).To(Equal("three"))
}

func Test_ReconstructResponse_CanReturnACompleteHttpResponseWithAllFieldsFilled(t *testing.T) {
	RegisterTestingT(t)

	req, _ := http.NewRequest("GET", "http://example.com", nil)

	pair := models.RequestResponsePair{
		Response: models.ResponseDetails{
			Status: 201,
			Body:   "test body",
		},
	}

	headers := make(map[string][]string)
	headers["Header"] = []string{"header test"}
	headers["Other"] = []string{"header"}
	pair.Response.Headers = headers

	response := modes.ReconstructResponse(req, pair)

	Expect(response.StatusCode).To(Equal(201))

	responseBody, err := ioutil.ReadAll(response.Body)
	Expect(err).To(BeNil())

	Expect(string(responseBody)).To(Equal("test body"))

	Expect(response.Header.Get("Header")).To(Equal(headers["Header"][0]))
	Expect(response.Header.Get("Other")).To(Equal(headers["Other"][0]))
}

func Test_errorResponse_ShouldAlwaysBeABadGatway(t *testing.T) {
	RegisterTestingT(t)

	response := modes.ErrorResponse(&http.Request{}, errors.New(""), "An error was got")

	Expect(response.StatusCode).To(Equal(http.StatusBadGateway))
}

func Test_errorResponse_ShouldAlwaysIncludeBothMessageAndErrorInResponseBody(t *testing.T) {
	RegisterTestingT(t)

	response := modes.ErrorResponse(&http.Request{}, errors.New("error doing something"), "This is a test error")

	responseBody, err := ioutil.ReadAll(response.Body)
	Expect(err).To(BeNil())

	Expect(string(responseBody)).To(ContainSubstring("Hoverfly Error!"))
	Expect(string(responseBody)).To(ContainSubstring("This is a test error"))
	Expect(string(responseBody)).To(ContainSubstring("error doing something"))
}

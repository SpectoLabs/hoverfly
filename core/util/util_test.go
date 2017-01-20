package util

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	. "github.com/onsi/gomega"
)

func Test_GetRequestBody_GettingTheRequestBodyGetsTheCorrectData(t *testing.T) {
	RegisterTestingT(t)

	request := &http.Request{}
	request.Body = ioutil.NopCloser(bytes.NewBuffer([]byte("test")))

	requestBody, err := GetRequestBody(request)
	Expect(err).To(BeNil())

	Expect(requestBody).To(Equal("test"))
}

func Test_GetRequestBody_GettingTheRequestBodySetsTheSameBodyAgain(t *testing.T) {
	RegisterTestingT(t)

	request := &http.Request{}
	request.Body = ioutil.NopCloser(bytes.NewBuffer([]byte("test-preserve")))

	_, err := GetRequestBody(request)
	Expect(err).To(BeNil())

	newRequestBody, err := ioutil.ReadAll(request.Body)
	Expect(err).To(BeNil())

	Expect(string(newRequestBody)).To(Equal("test-preserve"))
}

func Test_GetResponseBody_GettingTheResponseBodyGetsTheCorrectData(t *testing.T) {
	RegisterTestingT(t)

	response := &http.Response{}
	response.Body = ioutil.NopCloser(bytes.NewBuffer([]byte("test")))

	responseBody, err := GetResponseBody(response)
	Expect(err).To(BeNil())

	Expect(responseBody).To(Equal("test"))

}

func Test_GetResponseBody_GettingTheResponseBodySetsTheSameBodyAgain(t *testing.T) {
	RegisterTestingT(t)

	response := &http.Response{}
	response.Body = ioutil.NopCloser(bytes.NewBuffer([]byte("test-preserve")))

	_, err := GetResponseBody(response)
	Expect(err).To(BeNil())

	newResponseBody, err := ioutil.ReadAll(response.Body)
	Expect(err).To(BeNil())

	Expect(string(newResponseBody)).To(Equal("test-preserve"))
}

func Test_SortQueryString_ReordersQueryStringAlphabetically(t *testing.T) {
	RegisterTestingT(t)

	Expect(SortQueryString("e=e&d=d")).To(Equal("d=d&e=e"))
}

func Test_SortQueryString_ReordersQueryValuesAlphabetically(t *testing.T) {
	RegisterTestingT(t)

	Expect(SortQueryString("e=e&e=d")).To(Equal("e=d&e=e"))
}

func Test_SortQueryString_ReordersQueryValuesNumerically(t *testing.T) {
	RegisterTestingT(t)

	Expect(SortQueryString("e=2&e=1")).To(Equal("e=1&e=2"))
}

func Test_SortQueryString_ReordersQueryValuesAlphanumerically(t *testing.T) {
	RegisterTestingT(t)

	Expect(SortQueryString("e=2&e=d&e=1&e=e")).To(Equal("e=1&e=2&e=d&e=e"))
}

func Test_SortQueryString_KeepsAsteriskInTact(t *testing.T) {
	RegisterTestingT(t)

	Expect(SortQueryString("&e=*")).To(Equal("e=*"))
}

func Test_SortQueryString_PreservesEqualsAndEmptyValueQuery(t *testing.T) {
	RegisterTestingT(t)

	Expect(SortQueryString("e=")).To(Equal("e="))
}

func Test_SortQueryString_PreservesNoEqualsAndEmptyValueQuery(t *testing.T) {
	RegisterTestingT(t)

	Expect(SortQueryString("e")).To(Equal("e"))
}

func Test_SortQueryString_PreservesBothEqualsAndNoEqualsWithEmptyValue(t *testing.T) {
	RegisterTestingT(t)

	Expect(SortQueryString("a&b&c=&d&e=&f=")).To(Equal("a&b&c=&d&e=&f="))
}

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

	request, err := http.NewRequest("GET", "http://hoverfly.io", bytes.NewBuffer([]byte("test")))
	Expect(err).To(BeNil())

	requestBody, err := GetRequestBody(request)
	Expect(err).To(BeNil())

	Expect(requestBody).To(Equal("test"))

}

func Test_GetRequestBody_GettingTheRequestBodySetsTheSameBodyAgain(t *testing.T) {
	RegisterTestingT(t)

	request, err := http.NewRequest("GET", "http://hoverfly.io", bytes.NewBuffer([]byte("test-preserve")))
	Expect(err).To(BeNil())

	_, err = GetRequestBody(request)
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

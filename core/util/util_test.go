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

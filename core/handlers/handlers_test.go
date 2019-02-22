package handlers_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/SpectoLabs/hoverfly/core/handlers"
	. "github.com/onsi/gomega"
)

func Test_WriteResponse_WritesResponseToBody(t *testing.T) {
	RegisterTestingT(t)

	response := httptest.NewRecorder()
	handlers.WriteResponse(response, []byte("Test body"))

	Expect(response.Code).To(Equal(200))
	Expect(response.Header()["Content-Type"]).To(ContainElement("text/plain; charset=utf-8"))

	bodyBytes, err := ioutil.ReadAll(response.Body)
	Expect(err).To(BeNil())
	Expect(string(bodyBytes)).Should(Equal("Test body"))
}

func Test_WriteResponse_CorrectlySetsContentType(t *testing.T) {
	RegisterTestingT(t)

	jsonResponse := httptest.NewRecorder()
	handlers.WriteResponse(jsonResponse, []byte("{\"body\": \"json\"}"))

	Expect(jsonResponse.Header()["Content-Type"]).To(ContainElement("application/json; charset=utf-8"))

	xmlResponse := httptest.NewRecorder()
	handlers.WriteResponse(xmlResponse, []byte("<?xml>"))
	Expect(xmlResponse.Header()["Content-Type"]).To(ContainElement("text/xml; charset=utf-8"))

}

func Test_WriteResponseError_WritesErrorMessage(t *testing.T) {
	RegisterTestingT(t)

	response := httptest.NewRecorder()

	handlers.WriteErrorResponse(response, "This is an error", 5555)

	Expect(response.Code).To(Equal(5555))
	Expect(response.Header()["Content-Type"]).To(ContainElement("application/json; charset=utf-8"))

	errorView, err := unmarshalErrorView(response.Body)
	Expect(err).To(BeNil())

	Expect(errorView.Error).To(Equal("This is an error"))
}

func unmarshalErrorView(buffer *bytes.Buffer) (handlers.ErrorView, error) {
	body, err := ioutil.ReadAll(buffer)
	if err != nil {
		return handlers.ErrorView{}, err
	}

	var errorView handlers.ErrorView

	err = json.Unmarshal(body, &errorView)
	if err != nil {
		return handlers.ErrorView{}, err
	}

	return errorView, nil
}

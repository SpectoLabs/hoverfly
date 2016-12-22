package v2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	. "github.com/onsi/gomega"
)

type HoverflyMiddlewareStub struct {
	Binary     string
	Script     string
	Middleware string
	Remote     string
}

func (this HoverflyMiddlewareStub) GetMiddlewareV2() (string, string, string) {
	return this.Binary, this.Script, this.Remote
}

func (this *HoverflyMiddlewareStub) SetMiddlewareV2(binary, script, remote string) error {
	this.Binary = binary
	this.Script = script
	this.Remote = remote
	if script == "error" {
		return fmt.Errorf("error")
	}

	return nil
}

func TestHoverflyMiddlewareHandlerGetReturnsTheCorrectMiddleware(t *testing.T) {
	RegisterTestingT(t)

	stubHoverfly := &HoverflyMiddlewareStub{
		Binary: "test",
		Script: "middleware",
		Remote: "remote",
	}

	unit := HoverflyMiddlewareHandler{Hoverfly: stubHoverfly}

	request, err := http.NewRequest("GET", "", nil)
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Get, request)

	Expect(response.Code).To(Equal(http.StatusOK))

	middlewareView, err := unmarshalMiddlewareView(response.Body)
	Expect(err).To(BeNil())
	Expect(middlewareView.Binary).To(Equal("test"))
	Expect(middlewareView.Script).To(Equal("middleware"))
	Expect(middlewareView.Remote).To(Equal("remote"))
}

func TestHoverflyMiddlewareHandlerPutSetsTheNewMiddlewarendReplacesTheTestMiddleware(t *testing.T) {
	RegisterTestingT(t)

	stubHoverfly := &HoverflyMiddlewareStub{Binary: "test"}
	unit := HoverflyMiddlewareHandler{Hoverfly: stubHoverfly}

	middlewareView := &MiddlewareView{Binary: "python", Script: "new-middleware"}

	bodyBytes, err := json.Marshal(middlewareView)
	Expect(err).To(BeNil())

	request, err := http.NewRequest("PUT", "", ioutil.NopCloser(bytes.NewBuffer(bodyBytes)))
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Put, request)
	Expect(response.Code).To(Equal(http.StatusOK))
	Expect(stubHoverfly.Binary).To(Equal("python"))
	Expect(stubHoverfly.Script).To(Equal("new-middleware"))

	middlewareViewResponse, err := unmarshalMiddlewareView(response.Body)
	Expect(err).To(BeNil())

	Expect(middlewareViewResponse.Binary).To(Equal("python"))
	Expect(middlewareViewResponse.Script).To(Equal("new-middleware"))
}

func TestHoverflyMiddlewareHandlerPutWill422ErrorIfHoverflyErrors(t *testing.T) {
	RegisterTestingT(t)

	var stubHoverfly HoverflyMiddlewareStub
	unit := HoverflyMiddlewareHandler{Hoverfly: &stubHoverfly}

	middlewareView := &MiddlewareView{Script: "error"}

	bodyBytes, err := json.Marshal(middlewareView)
	Expect(err).To(BeNil())

	request, err := http.NewRequest("PUT", "", ioutil.NopCloser(bytes.NewBuffer(bodyBytes)))
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Put, request)
	Expect(response.Code).To(Equal(http.StatusUnprocessableEntity))

	errorViewResponse, err := unmarshalErrorView(response.Body)
	Expect(err).To(BeNil())

	Expect(errorViewResponse.Error).To(Equal("Invalid middleware: error"))
}

func TestHoverflyMiddlewareHandlerPutWill400ErrorIfJsonIsBad(t *testing.T) {
	RegisterTestingT(t)

	var stubHoverfly HoverflyMiddlewareStub
	unit := HoverflyMiddlewareHandler{Hoverfly: &stubHoverfly}

	bodyBytes := []byte("{{}{}}")

	request, err := http.NewRequest("PUT", "/api/v2/hoverfly/mode", ioutil.NopCloser(bytes.NewBuffer(bodyBytes)))
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Put, request)
	Expect(response.Code).To(Equal(http.StatusBadRequest))

	errorViewResponse, err := unmarshalErrorView(response.Body)
	Expect(err).To(BeNil())

	Expect(errorViewResponse.Error).To(Equal("Malformed JSON"))
}

func unmarshalMiddlewareView(buffer *bytes.Buffer) (MiddlewareView, error) {
	body, err := ioutil.ReadAll(buffer)
	if err != nil {
		return MiddlewareView{}, err
	}

	var middlewareView MiddlewareView

	err = json.Unmarshal(body, &middlewareView)
	if err != nil {
		return MiddlewareView{}, err
	}

	return middlewareView, nil
}

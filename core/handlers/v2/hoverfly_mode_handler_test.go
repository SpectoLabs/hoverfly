package v2

import (
	"testing"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"github.com/codegangsta/negroni"
	"encoding/json"
	"io/ioutil"
	"bytes"
	"fmt"
)

type HoverflyStub struct {
	Mode string
}

func (this HoverflyStub) GetMode() string {
	return this.Mode
}

func (this *HoverflyStub) SetMode(mode string) error {
	this.Mode = mode
	if mode == "error" {
		return fmt.Errorf("This is an error")
	}

	return nil
}

func TestGetReturnsTheCorrectMode(t *testing.T) {
	RegisterTestingT(t)

	stubHoverfly := &HoverflyStub{Mode: "test-mode"}
	unit := HoverflyModeHandler{Hoverfly: stubHoverfly}

	request, err := http.NewRequest("GET", "/api/v2/hoverfly/mode", nil)
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Get, request)

	Expect(response.Code).To(Equal(http.StatusOK))

	modeView, err := unmarshalModeView(response.Body)
	Expect(err).To(BeNil())
	Expect(modeView.Mode).To(Equal("test-mode"))
}

func TestPutSetsTheNewModeAndReplacesTheTestMode(t *testing.T) {
	RegisterTestingT(t)

	stubHoverfly := &HoverflyStub{Mode: "test-mode"}
	unit := HoverflyModeHandler{Hoverfly: stubHoverfly}

	modeView := &ModeView{Mode: "new-mode"}

	bodyBytes, err := json.Marshal(modeView)
	Expect(err).To(BeNil())

	request, err := http.NewRequest("PUT", "/api/v2/hoverfly/mode", ioutil.NopCloser(bytes.NewBuffer(bodyBytes)))
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Put, request)
	Expect(response.Code).To(Equal(http.StatusOK))
	Expect(stubHoverfly.Mode).To(Equal("new-mode"))

	modeViewResponse, err := unmarshalModeView(response.Body)
	Expect(err).To(BeNil())

	Expect(modeViewResponse.Mode).To(Equal("new-mode"))
}

func TestPutWill422ErrorIfHoverflyErrors(t *testing.T) {
	RegisterTestingT(t)

	var stubHoverfly HoverflyStub
	unit := HoverflyModeHandler{Hoverfly: &stubHoverfly}

	modeView := &ModeView{Mode: "error"}

	bodyBytes, err := json.Marshal(modeView)
	Expect(err).To(BeNil())

	request, err := http.NewRequest("PUT", "/api/v2/hoverfly/mode", ioutil.NopCloser(bytes.NewBuffer(bodyBytes)))
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Put, request)
	Expect(response.Code).To(Equal(http.StatusUnprocessableEntity))

	errorViewResponse, err := unmarshalErrorView(response.Body)
	Expect(err).To(BeNil())

	Expect(errorViewResponse.Error).To(Equal("This is an error"))
}

func TestPutWill400ErrorIfJsonIsBad(t *testing.T) {
	RegisterTestingT(t)

	var stubHoverfly HoverflyStub
	unit := HoverflyModeHandler{Hoverfly: &stubHoverfly}

	bodyBytes := []byte("{{}{}}")

	request, err := http.NewRequest("PUT", "/api/v2/hoverfly/mode", ioutil.NopCloser(bytes.NewBuffer(bodyBytes)))
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Put, request)
	Expect(response.Code).To(Equal(http.StatusBadRequest))

	errorViewResponse, err := unmarshalErrorView(response.Body)
	Expect(err).To(BeNil())

	Expect(errorViewResponse.Error).To(Equal("Malformed JSON"))
}

func makeRequestOnHandler(handlerFunc negroni.HandlerFunc, request *http.Request) *httptest.ResponseRecorder {
	responseRecorder := httptest.NewRecorder()
	negroni := negroni.New(handlerFunc)
	negroni.ServeHTTP(responseRecorder, request)
	return responseRecorder
}

func unmarshalModeView(buffer *bytes.Buffer) (ModeView, error) {
	body, err := ioutil.ReadAll(buffer)
	if err != nil {
		return ModeView{}, err
	}

	var modeView ModeView

	err = json.Unmarshal(body, &modeView)
	if err != nil {
		return ModeView{}, err
	}

	return modeView, nil
}

func unmarshalErrorView(buffer *bytes.Buffer) (ErrorView, error) {
	body, err := ioutil.ReadAll(buffer)
	if err != nil {
		return ErrorView{}, err
	}

	var errorView ErrorView

	err = json.Unmarshal(body, &errorView)
	if err != nil {
		return ErrorView{}, err
	}

	return errorView, nil
}
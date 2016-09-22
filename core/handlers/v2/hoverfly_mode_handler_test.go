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
	fmt.Println(mode)
	return nil
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

func TestGetReturnsTheCorrectMode(t *testing.T) {
	RegisterTestingT(t)

	var stubHoverfly HoverflyStub
	stubHoverfly.Mode = "test-mode"
	unit := HoverflyModeHandler{Hoverfly: &stubHoverfly}

	request, err := http.NewRequest("GET", "/api/v2/hoverfly/mode", nil)
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Get, request)

	Expect(response.Code).To(Equal(http.StatusOK))

	modeView, err := unmarshalModeView(response.Body)
	Expect(err).To(BeNil())
	Expect(modeView.Mode).To(Equal("test-mode"))
}

func TestPutSetsTheNewMode(t *testing.T) {
	RegisterTestingT(t)

	var stubHoverfly HoverflyStub
	unit := HoverflyModeHandler{Hoverfly: &stubHoverfly}

	modeView := &ModeView{Mode: "simulate"}

	bodyBytes, err := json.Marshal(modeView)
	Expect(err).To(BeNil())

	request, err := http.NewRequest("GET", "/api/v2/hoverfly/mode", ioutil.NopCloser(bytes.NewBuffer(bodyBytes)))
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Put, request)
	Expect(response.Code).To(Equal(http.StatusOK))
	Expect(stubHoverfly.Mode).To(Equal("simulate"))

	modeViewResponse, err := unmarshalModeView(response.Body)
	Expect(err).To(BeNil())

	Expect(modeViewResponse.Mode).To(Equal("simulate"))
}
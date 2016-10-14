package v2

import (
	"bytes"
	"encoding/json"
	"fmt"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
	"testing"
)

type HoverflyModeStub struct {
	Mode string
}

func (this HoverflyModeStub) GetMode() string {
	return this.Mode
}

func (this *HoverflyModeStub) SetMode(mode string) error {
	this.Mode = mode
	if mode == "error" {
		return fmt.Errorf("This is an error")
	}

	return nil
}

func TestGetReturnsTheCorrectMode(t *testing.T) {
	RegisterTestingT(t)

	stubHoverfly := &HoverflyModeStub{Mode: "test-mode"}
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

	stubHoverfly := &HoverflyModeStub{Mode: "test-mode"}
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

	var stubHoverfly HoverflyModeStub
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

	var stubHoverfly HoverflyModeStub
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

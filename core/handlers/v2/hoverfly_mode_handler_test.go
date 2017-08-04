package v2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/SpectoLabs/hoverfly/core/util"
	. "github.com/onsi/gomega"
)

type HoverflyModeStub struct {
	ModeView ModeView
}

func (this HoverflyModeStub) GetMode() ModeView {
	return this.ModeView
}

func (this *HoverflyModeStub) SetModeWithArguments(modeView ModeView) error {
	this.ModeView = modeView
	if modeView.Mode == "error" {
		return fmt.Errorf("This is an error")
	}

	return nil
}

func TestGetReturnsTheCorrectModeAndArguments(t *testing.T) {
	RegisterTestingT(t)

	stubHoverfly := &HoverflyModeStub{ModeView{
		Mode: "test-mode",
	}}

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

	stubHoverfly := &HoverflyModeStub{ModeView{
		Mode: "test-mode",
	}}
	unit := HoverflyModeHandler{Hoverfly: stubHoverfly}

	modeView := &ModeView{Mode: "new-mode"}

	bodyBytes, err := json.Marshal(modeView)
	Expect(err).To(BeNil())

	request, err := http.NewRequest("PUT", "/api/v2/hoverfly/mode", ioutil.NopCloser(bytes.NewBuffer(bodyBytes)))
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Put, request)
	Expect(response.Code).To(Equal(http.StatusOK))
	Expect(stubHoverfly.ModeView.Mode).To(Equal("new-mode"))

	modeViewResponse, err := unmarshalModeView(response.Body)
	Expect(err).To(BeNil())

	Expect(modeViewResponse.Mode).To(Equal("new-mode"))
}

func TestPutSetsTheArguments(t *testing.T) {
	RegisterTestingT(t)

	stubHoverfly := &HoverflyModeStub{ModeView{
		Mode: "test-mode",
	}}

	unit := HoverflyModeHandler{Hoverfly: stubHoverfly}

	modeView := &ModeView{
		Mode: "mode",
		Arguments: ModeArgumentsView{
			Headers:          []string{"argument"},
			MatchingStrategy: util.StringToPointer("strategy"),
		}}

	bodyBytes, err := json.Marshal(modeView)
	Expect(err).To(BeNil())

	request, err := http.NewRequest("PUT", "/api/v2/hoverfly/mode", ioutil.NopCloser(bytes.NewBuffer(bodyBytes)))
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Put, request)
	Expect(response.Code).To(Equal(http.StatusOK))

	_, err = unmarshalModeView(response.Body)
	Expect(err).To(BeNil())

	Expect(*stubHoverfly.ModeView.Arguments.MatchingStrategy).To(Equal("strategy"))
	Expect(stubHoverfly.ModeView.Arguments.Headers).To(ContainElement("argument"))
	Expect(stubHoverfly.ModeView.Mode).To(Equal("mode"))
}

func TestPutResetsTheArgumentsWhenNotSet(t *testing.T) {
	RegisterTestingT(t)

	stubHoverfly := &HoverflyModeStub{ModeView{
		Mode: "test-mode",
	}}

	unit := HoverflyModeHandler{Hoverfly: stubHoverfly}

	modeView := &ModeView{Arguments: ModeArgumentsView{
		Headers: []string{"argument"},
	}}

	bodyBytes, err := json.Marshal(modeView)
	Expect(err).To(BeNil())

	request, err := http.NewRequest("PUT", "/api/v2/hoverfly/mode", ioutil.NopCloser(bytes.NewBuffer(bodyBytes)))
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Put, request)
	Expect(response.Code).To(Equal(http.StatusOK))

	modeView.Arguments = ModeArgumentsView{}
	bodyBytes, err = json.Marshal(modeView)
	Expect(err).To(BeNil())

	request, err = http.NewRequest("PUT", "/api/v2/hoverfly/mode", ioutil.NopCloser(bytes.NewBuffer(bodyBytes)))
	Expect(err).To(BeNil())

	response = makeRequestOnHandler(unit.Put, request)
	Expect(response.Code).To(Equal(http.StatusOK))

	_, err = unmarshalModeView(response.Body)
	Expect(err).To(BeNil())

	Expect(stubHoverfly.ModeView.Arguments).To(Equal(ModeArgumentsView{}))
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
	Expect(response.Code).To(Equal(http.StatusBadRequest))

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

func Test_HoverflyModeHandler_Options_GetsOptions(t *testing.T) {
	RegisterTestingT(t)

	var stubHoverfly HoverflyModeStub
	unit := HoverflyModeHandler{Hoverfly: &stubHoverfly}

	request, err := http.NewRequest("OPTIONS", "/api/v2/hoverfly/mode", nil)
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Options, request)

	Expect(response.Code).To(Equal(http.StatusOK))
	Expect(response.Header().Get("Allow")).To(Equal("OPTIONS, GET, PUT"))
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

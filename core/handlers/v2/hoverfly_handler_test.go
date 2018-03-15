package v2

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/SpectoLabs/hoverfly/core/metrics"
	. "github.com/onsi/gomega"
)

type HoverflyStub struct{}

func (this HoverflyStub) GetDestination() string {
	return "test-destination.com"
}

func (this HoverflyStub) GetMode() ModeView {
	return ModeView{
		Mode: "test-mode",
		Arguments: ModeArgumentsView{
			Headers: []string{"test-header"},
		}}
}

func (this HoverflyStub) GetMiddleware() (string, string, string) {
	return "test-binary", "test-script", "test-remote"
}

func (this HoverflyStub) GetStats() metrics.Stats {
	metrics := metrics.Stats{
		Counters: make(map[string]int64),
	}

	metrics.Counters["countOne"] = int64(1)
	metrics.Counters["countTwo"] = int64(2)

	return metrics
}

func (this HoverflyStub) GetVersion() string {
	return "test-version"
}

func (this HoverflyStub) GetUpstreamProxy() string {
	return "test-proxy.com:8080"
}

func (this *HoverflyStub) GetState() map[string]string {
	return nil
}

func (this *HoverflyStub) SetState(state map[string]string) {
}

func (this *HoverflyStub) PatchState(state map[string]string) {
}

func (this *HoverflyStub) ClearState() {
}

func (this *HoverflyStub) IsWebServer() bool {
	return false
}

func (this *HoverflyStub) GetDiff() map[SimpleRequestDefinitionView][]DiffReport {
	return nil
}

func (this *HoverflyStub) ClearDiff() {
}

func TestHoverflyHandlerGetReturnsTheCorrectMode(t *testing.T) {
	RegisterTestingT(t)

	stubHoverfly := &HoverflyStub{}
	unit := HoverflyHandler{Hoverfly: stubHoverfly}

	request, err := http.NewRequest("GET", "", nil)
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Get, request)

	Expect(response.Code).To(Equal(http.StatusOK))

	hoverflyView, err := unmarshalHoverflyView(response.Body)
	Expect(err).To(BeNil())
	Expect(hoverflyView.Destination).To(Equal("test-destination.com"))
	Expect(hoverflyView.Mode).To(Equal("test-mode"))
	Expect(hoverflyView.Arguments.Headers).To(ContainElement("test-header"))
	Expect(hoverflyView.Binary).To(Equal("test-binary"))
	Expect(hoverflyView.Script).To(Equal("test-script"))
	Expect(hoverflyView.Remote).To(Equal("test-remote"))
	Expect(hoverflyView.Version).To(Equal("test-version"))
	Expect(hoverflyView.UpstreamProxy).To(Equal("test-proxy.com:8080"))
	Expect(hoverflyView.IsWebServer).To(BeFalse())
}

func Test_HoverflyHandler_Options_GetsOptions(t *testing.T) {
	RegisterTestingT(t)

	var stubHoverfly HoverflyStub
	unit := HoverflyHandler{Hoverfly: &stubHoverfly}

	request, err := http.NewRequest("OPTIONS", "/api/v2/hoverfly", nil)
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Options, request)

	Expect(response.Code).To(Equal(http.StatusOK))
	Expect(response.Header().Get("Allow")).To(Equal("OPTIONS, GET"))
}

func unmarshalHoverflyView(buffer *bytes.Buffer) (HoverflyView, error) {
	body, err := ioutil.ReadAll(buffer)
	if err != nil {
		return HoverflyView{}, err
	}

	var hoverflyView HoverflyView

	err = json.Unmarshal(body, &hoverflyView)
	if err != nil {
		return HoverflyView{}, err
	}

	return hoverflyView, nil
}

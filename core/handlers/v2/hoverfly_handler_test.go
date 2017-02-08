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

func (this HoverflyStub) GetMode() string {
	return "test-mode"
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
	Expect(hoverflyView.Binary).To(Equal("test-binary"))
	Expect(hoverflyView.Script).To(Equal("test-script"))
	Expect(hoverflyView.Remote).To(Equal("test-remote"))
	Expect(hoverflyView.Version).To(Equal("test-version"))
	Expect(hoverflyView.UpstreamProxy).To(Equal("test-proxy.com:8080"))
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

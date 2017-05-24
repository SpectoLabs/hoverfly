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

type HoverflyUsageStub struct{}

func (this HoverflyUsageStub) GetStats() metrics.Stats {
	metrics := metrics.Stats{
		Counters: make(map[string]int64),
	}

	metrics.Counters["countOne"] = int64(1)
	metrics.Counters["countTwo"] = int64(2)

	return metrics
}

func TestHoverflyUsageHandlerGetReturnsMetrics(t *testing.T) {
	RegisterTestingT(t)

	stubHoverfly := &HoverflyUsageStub{}
	unit := HoverflyUsageHandler{Hoverfly: stubHoverfly}

	request, err := http.NewRequest("GET", "", nil)
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Get, request)

	Expect(response.Code).To(Equal(http.StatusOK))

	usageView, err := unmarshalUsageView(response.Body)
	Expect(err).To(BeNil())

	Expect(usageView.Usage.Counters).To(HaveLen(2))
	Expect(usageView.Usage.Counters).To(HaveKeyWithValue("countOne", int64(1)))
	Expect(usageView.Usage.Counters).To(HaveKeyWithValue("countTwo", int64(2)))
}

func Test_HoverflyUsageHandler_Options_GetsOptions(t *testing.T) {
	RegisterTestingT(t)

	var stubHoverfly HoverflyUsageStub
	unit := HoverflyUsageHandler{Hoverfly: &stubHoverfly}

	request, err := http.NewRequest("OPTIONS", "/api/v2/hoverfly/usage", nil)
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Options, request)

	Expect(response.Code).To(Equal(http.StatusOK))
	Expect(response.Header().Get("Allow")).To(Equal("OPTIONS, GET"))
}

func unmarshalUsageView(buffer *bytes.Buffer) (UsageView, error) {
	body, err := ioutil.ReadAll(buffer)
	if err != nil {
		return UsageView{}, err
	}

	var metricsView UsageView

	err = json.Unmarshal(body, &metricsView)
	if err != nil {
		return UsageView{}, err
	}

	return metricsView, nil
}

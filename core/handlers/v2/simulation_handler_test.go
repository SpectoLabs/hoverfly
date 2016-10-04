package v2

import (
	"bytes"
	"encoding/json"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
	"testing"
	"github.com/SpectoLabs/hoverfly/core/util"
	"github.com/SpectoLabs/hoverfly/core/handlers/v1"
)

type HoverflySimulationStub struct{}

func (this HoverflySimulationStub) GetSimulation() (SimulationView, error) {
	pairOne := RequestResponsePairView{
		Request: RequestDetailsView{
			Destination: util.StringToPointer("test.com"),
			Path: util.StringToPointer("/testing"),
		},
		Response: ResponseDetailsView{
			Body: "test-body",
		},

	}

	return SimulationView {
		DataView {
			RequestResponsePairs: []RequestResponsePairView{pairOne},
			GlobalActions: GlobalActionsView{
				Delays: []v1.ResponseDelayView{
					{
						HttpMethod: "GET",
						Delay: 100,
					},
				},
			},
		},
		MetaView {
			SchemaVersion: "v1",
			HoverflyVersion: "test",
			TimeExported: "now",
		},
	}, nil
}

func TestSimulationHandlerGetReturnsSimulation(t *testing.T) {
	RegisterTestingT(t)

	stubHoverfly := &HoverflySimulationStub{}
	unit := SimulationHandler{Hoverfly: stubHoverfly}

	request, err := http.NewRequest("GET", "", nil)
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Get, request)

	Expect(response.Code).To(Equal(http.StatusOK))

	simulationView, err := unmarshalSimulationView(response.Body)
	Expect(err).To(BeNil())

	Expect(simulationView.DataView.RequestResponsePairs).To(HaveLen(1))

	Expect(simulationView.DataView.RequestResponsePairs[0].Request.Destination).To(Equal(util.StringToPointer("test.com")))
	Expect(simulationView.DataView.RequestResponsePairs[0].Request.Path).To(Equal(util.StringToPointer("/testing")))

	Expect(simulationView.DataView.RequestResponsePairs[0].Response.Body).To(Equal("test-body"))

	Expect(simulationView.DataView.GlobalActions.Delays).To(HaveLen(1))
	Expect(simulationView.DataView.GlobalActions.Delays[0].HttpMethod).To(Equal("GET"))
	Expect(simulationView.DataView.GlobalActions.Delays[0].Delay).To(Equal(100))

	Expect(simulationView.MetaView.SchemaVersion).To(Equal("v1"))
	Expect(simulationView.MetaView.HoverflyVersion).To(Equal("test"))
	Expect(simulationView.MetaView.TimeExported).To(Equal("now"))
}

func unmarshalSimulationView(buffer *bytes.Buffer) (SimulationView, error) {
	body, err := ioutil.ReadAll(buffer)
	if err != nil {
		return SimulationView{}, err
	}

	var simulationView SimulationView

	err = json.Unmarshal(body, &simulationView)
	if err != nil {
		return SimulationView{}, err
	}

	return simulationView, nil
}

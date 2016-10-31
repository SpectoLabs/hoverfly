package v2

import (
	"bytes"
	"encoding/json"
	"github.com/SpectoLabs/hoverfly/core/handlers/v1"
	"github.com/SpectoLabs/hoverfly/core/util"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
	"testing"
)

type HoverflySimulationStub struct {
	Deleted bool
	Simulation SimulationView
}

func (this HoverflySimulationStub) GetSimulation() (SimulationView, error) {
	pairOne := RequestResponsePairView{
		Request: RequestDetailsView{
			Destination: util.StringToPointer("test.com"),
			Path:        util.StringToPointer("/testing"),
		},
		Response: ResponseDetailsView{
			Body: "test-body",
		},
	}

	return SimulationView{
		DataView{
			RequestResponsePairs: []RequestResponsePairView{pairOne},
			GlobalActions: GlobalActionsView{
				Delays: []v1.ResponseDelayView{
					{
						HttpMethod: "GET",
						Delay:      100,
					},
				},
			},
		},
		MetaView{
			SchemaVersion:   "v1",
			HoverflyVersion: "test",
			TimeExported:    "now",
		},
	}, nil
}

func (this *HoverflySimulationStub) DeleteSimulation() error {
	this.Deleted = true
	return nil
}

func (this *HoverflySimulationStub) PutSimulation(simulation SimulationView) (error) {
	this.Simulation = simulation
	return nil
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

func TestSimulationHandlerDeleteCallsDelete(t *testing.T) {
	RegisterTestingT(t)

	stubHoverfly := &HoverflySimulationStub{}
	Expect(stubHoverfly.Deleted).To(BeFalse())

	unit := SimulationHandler{Hoverfly: stubHoverfly}

	request, err := http.NewRequest("DELETE", "", nil)
	Expect(err).To(BeNil())

	makeRequestOnHandler(unit.Delete, request)

	Expect(stubHoverfly.Deleted).To(BeTrue())
}

func TestSimulationHandler_Delete_CallsGetAfterDelete(t *testing.T) {
	RegisterTestingT(t)

	stubHoverfly := &HoverflySimulationStub{}

	unit := SimulationHandler{Hoverfly: stubHoverfly}

	request, err := http.NewRequest("DELETE", "", nil)
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Delete, request)

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

func TestSimulationHandler_Put_PassesDataIntoHoverfly(t *testing.T) {
	RegisterTestingT(t)

	stubHoverfly := &HoverflySimulationStub{}

	unit := SimulationHandler{Hoverfly: stubHoverfly}

	request, err := http.NewRequest("PUT", "", ioutil.NopCloser(bytes.NewBuffer([]byte(`
	{
		"data": {
			"pairs": [
				{
					"request": {
						"destination": "test.org"
					},
					"response": {
						"status": 200
					}
				}
			],

			"globalActions": {
				"delays": [
					{
						"urlPattern": "test.org",
						"httpMethod": "GET",
						"delay": 200
					}
				]
			}

		}
	}
	`))))
	Expect(err).To(BeNil())

	makeRequestOnHandler(unit.Put, request)

	Expect(stubHoverfly.Simulation).ToNot(BeNil())
	Expect(stubHoverfly.Simulation.RequestResponsePairs).ToNot(BeNil())

	Expect(stubHoverfly.Simulation.RequestResponsePairs[0].Request.Destination).To(Equal(util.StringToPointer("test.org")))
	Expect(stubHoverfly.Simulation.RequestResponsePairs[0].Response.Status).To(Equal(200))

	Expect(stubHoverfly.Simulation.GlobalActions.Delays[0].UrlPattern).To(Equal("test.org"))
	Expect(stubHoverfly.Simulation.GlobalActions.Delays[0].HttpMethod).To(Equal("GET"))
	Expect(stubHoverfly.Simulation.GlobalActions.Delays[0].Delay).To(Equal(200))
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

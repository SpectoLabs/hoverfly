package v2

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"fmt"

	"github.com/SpectoLabs/hoverfly/core/handlers/v1"
	"github.com/SpectoLabs/hoverfly/core/util"
	. "github.com/onsi/gomega"
)

type HoverflySimulationStub struct {
	Deleted    bool
	Simulation SimulationViewV2
}

func (this HoverflySimulationStub) GetSimulation() (SimulationViewV2, error) {
	pairOne := RequestMatcherResponsePairViewV2{
		RequestMatcher: RequestMatcherViewV2{
			Destination: &RequestFieldMatchersView{
				ExactMatch: util.StringToPointer("test.com"),
			},
			Path: &RequestFieldMatchersView{
				ExactMatch: util.StringToPointer("/testing"),
			},
		},
		Response: ResponseDetailsView{
			Body: "test-body",
		},
	}

	return SimulationViewV2{
		DataViewV2{
			RequestResponsePairs: []RequestMatcherResponsePairViewV2{pairOne},
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
			SchemaVersion:   "v2",
			HoverflyVersion: "test",
			TimeExported:    "now",
		},
	}, nil
}

func (this *HoverflySimulationStub) DeleteSimulation() {
	this.Deleted = true
}

func (this *HoverflySimulationStub) PutSimulation(simulation SimulationViewV2) error {
	this.Simulation = simulation
	return nil
}

type HoverflySimulationErrorStub struct{}

func (this HoverflySimulationErrorStub) GetSimulation() (SimulationViewV2, error) {
	return SimulationViewV2{}, fmt.Errorf("error")
}

func (this *HoverflySimulationErrorStub) DeleteSimulation() {}

func (this *HoverflySimulationErrorStub) PutSimulation(simulation SimulationViewV2) error {
	return fmt.Errorf("error")
}

func TestSimulationHandler_Get_ReturnsSimulation(t *testing.T) {
	RegisterTestingT(t)

	stubHoverfly := &HoverflySimulationStub{}
	unit := SimulationHandler{Hoverfly: stubHoverfly}

	request, err := http.NewRequest("GET", "", nil)
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Get, request)

	Expect(response.Code).To(Equal(http.StatusOK))

	simulationView, err := unmarshalSimulationViewV2(response.Body)
	Expect(err).To(BeNil())

	Expect(simulationView.DataViewV2.RequestResponsePairs).To(HaveLen(1))

	Expect(simulationView.DataViewV2.RequestResponsePairs[0].RequestMatcher.Destination.ExactMatch).To(Equal(util.StringToPointer("test.com")))
	Expect(simulationView.DataViewV2.RequestResponsePairs[0].RequestMatcher.Path.ExactMatch).To(Equal(util.StringToPointer("/testing")))

	Expect(simulationView.DataViewV2.RequestResponsePairs[0].Response.Body).To(Equal("test-body"))

	Expect(simulationView.DataViewV2.GlobalActions.Delays).To(HaveLen(1))
	Expect(simulationView.DataViewV2.GlobalActions.Delays[0].HttpMethod).To(Equal("GET"))
	Expect(simulationView.DataViewV2.GlobalActions.Delays[0].Delay).To(Equal(100))

	Expect(simulationView.MetaView.SchemaVersion).To(Equal("v2"))
	Expect(simulationView.MetaView.HoverflyVersion).To(Equal("test"))
	Expect(simulationView.MetaView.TimeExported).To(Equal("now"))
}

func TestSimulationHandler_Get_ReturnsErrorIfHoverflyErrors(t *testing.T) {
	RegisterTestingT(t)

	stubHoverfly := &HoverflySimulationErrorStub{}
	unit := SimulationHandler{Hoverfly: stubHoverfly}

	request, err := http.NewRequest("GET", "", nil)
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Get, request)

	Expect(response.Code).To(Equal(http.StatusInternalServerError))

	errorView, err := unmarshalErrorView(response.Body)
	Expect(err).To(BeNil())

	Expect(errorView.Error).To(Equal("error"))
}

func TestSimulationHandler_Delete_CallsDelete(t *testing.T) {
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

	simulationView, err := unmarshalSimulationViewV2(response.Body)
	Expect(err).To(BeNil())

	Expect(simulationView.DataViewV2.RequestResponsePairs).To(HaveLen(1))

	Expect(simulationView.DataViewV2.RequestResponsePairs[0].RequestMatcher.Destination.ExactMatch).To(Equal(util.StringToPointer("test.com")))
	Expect(simulationView.DataViewV2.RequestResponsePairs[0].RequestMatcher.Path.ExactMatch).To(Equal(util.StringToPointer("/testing")))

	Expect(simulationView.DataViewV2.RequestResponsePairs[0].Response.Body).To(Equal("test-body"))

	Expect(simulationView.DataViewV2.GlobalActions.Delays).To(HaveLen(1))
	Expect(simulationView.DataViewV2.GlobalActions.Delays[0].HttpMethod).To(Equal("GET"))
	Expect(simulationView.DataViewV2.GlobalActions.Delays[0].Delay).To(Equal(100))

	Expect(simulationView.MetaView.SchemaVersion).To(Equal("v2"))
	Expect(simulationView.MetaView.HoverflyVersion).To(Equal("test"))
	Expect(simulationView.MetaView.TimeExported).To(Equal("now"))
}

func TestSimulationHandler_Delete_ErrorReturnsWithoutGet(t *testing.T) {
	RegisterTestingT(t)

	stubHoverfly := &HoverflySimulationErrorStub{}

	unit := SimulationHandler{Hoverfly: stubHoverfly}

	request, err := http.NewRequest("DELETE", "", nil)
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Delete, request)

	errorView, err := unmarshalErrorView(response.Body)
	Expect(err).To(BeNil())

	Expect(errorView.Error).To(Equal("error"))
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
						"destination": {
							"exactMatch": "test.org"
						}
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
		},
		"meta": {
			"schemaVersion": "v2"
		}
	}
	`))))
	Expect(err).To(BeNil())

	makeRequestOnHandler(unit.Put, request)

	Expect(stubHoverfly.Simulation).ToNot(BeNil())
	Expect(stubHoverfly.Simulation.RequestResponsePairs).ToNot(BeNil())

	Expect(stubHoverfly.Simulation.RequestResponsePairs[0].RequestMatcher.Destination.ExactMatch).To(Equal(util.StringToPointer("test.org")))
	Expect(stubHoverfly.Simulation.RequestResponsePairs[0].Response.Status).To(Equal(200))

	Expect(stubHoverfly.Simulation.GlobalActions.Delays[0].UrlPattern).To(Equal("test.org"))
	Expect(stubHoverfly.Simulation.GlobalActions.Delays[0].HttpMethod).To(Equal("GET"))
	Expect(stubHoverfly.Simulation.GlobalActions.Delays[0].Delay).To(Equal(200))
}

func TestSimulationHandler_Put_CallsDelete(t *testing.T) {
	RegisterTestingT(t)

	stubHoverfly := &HoverflySimulationStub{}

	unit := SimulationHandler{Hoverfly: stubHoverfly}

	request, err := http.NewRequest("PUT", "", ioutil.NopCloser(bytes.NewBuffer([]byte(`
	{
		"data": {
			"pairs": [
				{
					"request": {
						"destination": {
							"exactMatch": "test.org"
						}
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
		},
		"meta": {
			"schemaVersion": "v2"
		}
	}
	`))))
	Expect(err).To(BeNil())

	makeRequestOnHandler(unit.Put, request)

	Expect(stubHoverfly.Deleted).To(BeTrue())
}

func TestSimulationHandler_Put_ReturnsErrorIfJsonDoesntMatchSchema_MissingDataKey(t *testing.T) {
	RegisterTestingT(t)

	stubHoverfly := &HoverflySimulationErrorStub{}

	unit := SimulationHandler{Hoverfly: stubHoverfly}

	request, err := http.NewRequest("PUT", "", ioutil.NopCloser(bytes.NewBuffer([]byte(`{"meta": {"schemaVersion": "v2"}}`))))
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Put, request)

	errorView, err := unmarshalErrorView(response.Body)
	Expect(err).To(BeNil())

	Expect(response.Result().StatusCode).To(Equal(400))
	Expect(errorView.Error).To(Equal("Invalid v2 simulation: data is required"))
}

func TestSimulationHandler_Put_ReturnsErrorIfJsonDoesntMatchSchema_EmptyObject(t *testing.T) {
	RegisterTestingT(t)

	stubHoverfly := &HoverflySimulationErrorStub{}

	unit := SimulationHandler{Hoverfly: stubHoverfly}

	request, err := http.NewRequest("PUT", "", ioutil.NopCloser(bytes.NewBuffer([]byte(`{}`))))
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Put, request)

	errorView, err := unmarshalErrorView(response.Body)
	Expect(err).To(BeNil())

	Expect(response.Result().StatusCode).To(Equal(400))
	Expect(errorView.Error).To(Equal(`Invalid JSON, missing "meta" object`))
}

func TestSimulationHandler_Put_ReturnsErrorIfJsonIsNotValid(t *testing.T) {
	RegisterTestingT(t)

	stubHoverfly := &HoverflySimulationErrorStub{}

	unit := SimulationHandler{Hoverfly: stubHoverfly}

	request, err := http.NewRequest("PUT", "", ioutil.NopCloser(bytes.NewBuffer([]byte(`{notdata: {{]]}[}]}""}`))))
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Put, request)

	errorView, err := unmarshalErrorView(response.Body)
	Expect(err).To(BeNil())

	Expect(response.Result().StatusCode).To(Equal(400))
	Expect(errorView.Error).To(Equal("Invalid JSON"))
}

func unmarshalSimulationViewV2(buffer *bytes.Buffer) (SimulationViewV2, error) {
	body, err := ioutil.ReadAll(buffer)
	if err != nil {
		return SimulationViewV2{}, err
	}

	var simulationView SimulationViewV2

	err = json.Unmarshal(body, &simulationView)
	if err != nil {
		return SimulationViewV2{}, err
	}

	return simulationView, nil
}

package v2

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"fmt"

	"github.com/SpectoLabs/hoverfly/core/handlers/v1"
	"github.com/SpectoLabs/hoverfly/core/matching/matchers"
	. "github.com/onsi/gomega"
)

type HoverflySimulationStub struct {
	Deleted    bool
	Simulation SimulationViewV6
	UrlPattern string
	Filtered   bool
}

func (this HoverflySimulationStub) GetSimulation() (SimulationViewV6, error) {
	pairOne := RequestMatcherResponsePairViewV6{
		RequestMatcher: RequestMatcherViewV6{
			Destination: []MatcherViewV6{
				NewMatcherViewV6(matchers.Exact, "test.com"),
			},
			Path: []MatcherViewV6{
				NewMatcherViewV6(matchers.Exact, "/testing"),
			},
		},
		Response: ResponseDetailsViewV6{
			Body: "test-body",
		},
	}

	return SimulationViewV6{
		DataViewV6{
			RequestResponsePairs: []RequestMatcherResponsePairViewV6{pairOne},
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
			SchemaVersion:   "v3",
			HoverflyVersion: "test",
			TimeExported:    "now",
		},
	}, nil
}

func (this *HoverflySimulationStub) GetFilteredSimulation(urlPattern string) (SimulationViewV6, error) {
	this.Filtered = true
	this.UrlPattern = urlPattern
	return this.GetSimulation()
}

func (this *HoverflySimulationStub) DeleteSimulation() {
	this.Deleted = true
}

func (this *HoverflySimulationStub) PutSimulation(simulation SimulationViewV6, overrideExisting bool) SimulationImportResult {
	this.Simulation = simulation
	this.Deleted = overrideExisting
	return SimulationImportResult{}
}

type HoverflySimulationErrorStub struct{}

func (this HoverflySimulationErrorStub) GetSimulation() (SimulationViewV6, error) {
	return SimulationViewV6{}, fmt.Errorf("error")
}

func (this HoverflySimulationErrorStub) GetFilteredSimulation(urlPattern string) (SimulationViewV6, error) {
	return SimulationViewV6{}, fmt.Errorf("error")
}

func (this *HoverflySimulationErrorStub) DeleteSimulation() {}

func (this *HoverflySimulationErrorStub) PutSimulation(simulation SimulationViewV6, overrideExisting bool) SimulationImportResult {
	return SimulationImportResult{
		err: fmt.Errorf("error"),
	}
}

type HoverflySimulationWarningStub struct{}

func (this HoverflySimulationWarningStub) GetSimulation() (SimulationViewV6, error) {
	return SimulationViewV6{}, fmt.Errorf("error")
}

func (this HoverflySimulationWarningStub) GetFilteredSimulation(urlPattern string) (SimulationViewV6, error) {
	return SimulationViewV6{}, fmt.Errorf("error")
}

func (this *HoverflySimulationWarningStub) DeleteSimulation() {}

func (this *HoverflySimulationWarningStub) PutSimulation(simulation SimulationViewV6, overrideExisting bool) SimulationImportResult {
	return SimulationImportResult{
		WarningMessages: []SimulationImportWarning{{"This is a warning", "url"}},
	}
}

func TestSimulationHandler_Get_ReturnsSimulation(t *testing.T) {
	RegisterTestingT(t)

	stubHoverfly := &HoverflySimulationStub{}
	unit := SimulationHandler{Hoverfly: stubHoverfly}

	request, err := http.NewRequest("GET", "", nil)
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Get, request)

	Expect(response.Code).To(Equal(http.StatusOK))

	simulationView, err := unmarshalSimulationViewV6(response.Body)
	Expect(err).To(BeNil())

	Expect(simulationView.DataViewV6.RequestResponsePairs).To(HaveLen(1))

	Expect(simulationView.DataViewV6.RequestResponsePairs[0].RequestMatcher.Destination[0].Matcher).To(Equal("exact"))
	Expect(simulationView.DataViewV6.RequestResponsePairs[0].RequestMatcher.Destination[0].Value).To(Equal("test.com"))

	Expect(simulationView.DataViewV6.RequestResponsePairs[0].RequestMatcher.Path[0].Matcher).To(Equal("exact"))
	Expect(simulationView.DataViewV6.RequestResponsePairs[0].RequestMatcher.Path[0].Value).To(Equal("/testing"))

	Expect(simulationView.DataViewV6.RequestResponsePairs[0].Response.Body).To(Equal("test-body"))

	Expect(simulationView.DataViewV6.GlobalActions.Delays).To(HaveLen(1))
	Expect(simulationView.DataViewV6.GlobalActions.Delays[0].HttpMethod).To(Equal("GET"))
	Expect(simulationView.DataViewV6.GlobalActions.Delays[0].Delay).To(Equal(100))

	Expect(simulationView.MetaView.SchemaVersion).To(Equal("v3"))
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

func TestSimulationHandler_Get_WithEmptyUrlPatternShouldNotFilterSimulation(t *testing.T) {
	RegisterTestingT(t)

	stubHoverfly := &HoverflySimulationStub{}
	unit := SimulationHandler{Hoverfly: stubHoverfly}

	request, err := http.NewRequest("GET", "?urlPattern=", nil)
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Get, request)

	Expect(response.Code).To(Equal(http.StatusOK))

	simulationView, err := unmarshalSimulationViewV6(response.Body)
	Expect(err).To(BeNil())

	Expect(simulationView.DataViewV6.RequestResponsePairs).To(HaveLen(1))
	Expect(stubHoverfly.Filtered).To(BeFalse())
}

func TestSimulationHandler_Get_WithUrlPatternShouldFilterSimulation(t *testing.T) {
	RegisterTestingT(t)

	stubHoverfly := &HoverflySimulationStub{}
	unit := SimulationHandler{Hoverfly: stubHoverfly}

	request, err := http.NewRequest("GET", "?urlPattern=foo.com", nil)
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Get, request)

	Expect(response.Code).To(Equal(http.StatusOK))

	simulationView, err := unmarshalSimulationViewV6(response.Body)
	Expect(err).To(BeNil())

	Expect(simulationView.DataViewV6.RequestResponsePairs).To(HaveLen(1))
	Expect(stubHoverfly.Filtered).To(BeTrue())
	Expect(stubHoverfly.UrlPattern).To(Equal("foo.com"))
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

	simulationView, err := unmarshalSimulationViewV6(response.Body)
	Expect(err).To(BeNil())

	Expect(simulationView.DataViewV6.RequestResponsePairs).To(HaveLen(1))

	Expect(simulationView.DataViewV6.RequestResponsePairs[0].RequestMatcher.Destination[0].Matcher).To(Equal("exact"))
	Expect(simulationView.DataViewV6.RequestResponsePairs[0].RequestMatcher.Destination[0].Value).To(Equal("test.com"))

	Expect(simulationView.DataViewV6.RequestResponsePairs[0].RequestMatcher.Path[0].Matcher).To(Equal("exact"))
	Expect(simulationView.DataViewV6.RequestResponsePairs[0].RequestMatcher.Path[0].Value).To(Equal("/testing"))

	Expect(simulationView.DataViewV6.RequestResponsePairs[0].Response.Body).To(Equal("test-body"))

	Expect(simulationView.DataViewV6.GlobalActions.Delays).To(HaveLen(1))
	Expect(simulationView.DataViewV6.GlobalActions.Delays[0].HttpMethod).To(Equal("GET"))
	Expect(simulationView.DataViewV6.GlobalActions.Delays[0].Delay).To(Equal(100))

	Expect(simulationView.MetaView.SchemaVersion).To(Equal("v3"))
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
			"schemaVersion": "v3"
		}
	}
	`))))
	Expect(err).To(BeNil())

	makeRequestOnHandler(unit.Put, request)

	Expect(stubHoverfly.Simulation).ToNot(BeNil())
	Expect(stubHoverfly.Simulation.RequestResponsePairs).ToNot(BeNil())

	Expect(stubHoverfly.Simulation.RequestResponsePairs[0].RequestMatcher.Destination[0].Matcher).To(Equal("exact"))
	Expect(stubHoverfly.Simulation.RequestResponsePairs[0].RequestMatcher.Destination[0].Value).To(Equal("test.org"))
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
			"schemaVersion": "v3"
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

	request, err := http.NewRequest("PUT", "", ioutil.NopCloser(bytes.NewBuffer([]byte(`{"meta": {"schemaVersion": "v3"}}`))))
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Put, request)

	errorView, err := unmarshalErrorView(response.Body)
	Expect(err).To(BeNil())

	Expect(response.Result().StatusCode).To(Equal(400))
	Expect(errorView.Error).To(Equal("Invalid v3 simulation: [Error for <data>: data is required]"))
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

func TestSimulationHandler_Put_ReturnsWarnings(t *testing.T) {
	RegisterTestingT(t)

	stubHoverfly := &HoverflySimulationWarningStub{}

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
				"schemaVersion": "v3"
			}
		}
		`))))
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Put, request)

	resultView, err := unmarshalResultView(response.Body)
	Expect(err).To(BeNil())

	Expect(response.Result().StatusCode).To(Equal(http.StatusOK))
	Expect(resultView.WarningMessages[0].Message).To(Equal("This is a warning"))
	Expect(resultView.WarningMessages[0].DocsLink).To(Equal("url"))

}

func TestSimulationHandler_Post_PassesDataIntoHoverfly(t *testing.T) {
	RegisterTestingT(t)

	stubHoverfly := &HoverflySimulationStub{}

	unit := SimulationHandler{Hoverfly: stubHoverfly}

	request, err := http.NewRequest("POST", "", ioutil.NopCloser(bytes.NewBuffer([]byte(`
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
			"schemaVersion": "v3"
		}
	}
	`))))
	Expect(err).To(BeNil())

	makeRequestOnHandler(unit.Post, request)

	Expect(stubHoverfly.Simulation).ToNot(BeNil())
	Expect(stubHoverfly.Simulation.RequestResponsePairs).ToNot(BeNil())

	Expect(stubHoverfly.Simulation.RequestResponsePairs[0].RequestMatcher.Destination[0].Matcher).To(Equal("exact"))
	Expect(stubHoverfly.Simulation.RequestResponsePairs[0].RequestMatcher.Destination[0].Value).To(Equal("test.org"))
	Expect(stubHoverfly.Simulation.RequestResponsePairs[0].Response.Status).To(Equal(200))

	Expect(stubHoverfly.Simulation.GlobalActions.Delays[0].UrlPattern).To(Equal("test.org"))
	Expect(stubHoverfly.Simulation.GlobalActions.Delays[0].HttpMethod).To(Equal("GET"))
	Expect(stubHoverfly.Simulation.GlobalActions.Delays[0].Delay).To(Equal(200))
}

func TestSimulationHandler_Post_NotCallsDelete(t *testing.T) {
	RegisterTestingT(t)

	stubHoverfly := &HoverflySimulationStub{}

	unit := SimulationHandler{Hoverfly: stubHoverfly}

	request, err := http.NewRequest("POST", "", ioutil.NopCloser(bytes.NewBuffer([]byte(`
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
			"schemaVersion": "v3"
		}
	}
	`))))
	Expect(err).To(BeNil())

	makeRequestOnHandler(unit.Post, request)

	Expect(stubHoverfly.Deleted).To(BeFalse())
}

func Test_SimulationHandler_Options_GetsOptions(t *testing.T) {
	RegisterTestingT(t)

	var stubHoverfly HoverflySimulationStub
	unit := SimulationHandler{Hoverfly: &stubHoverfly}

	request, err := http.NewRequest("OPTIONS", "/api/v2/simulation", nil)
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Options, request)

	Expect(response.Code).To(Equal(http.StatusOK))
	Expect(response.Header().Get("Allow")).To(Equal("OPTIONS, GET, PUT, DELETE"))
}

func Test_SimulationHandler_OptionsSchema_GetsOptions(t *testing.T) {
	RegisterTestingT(t)

	var stubHoverfly HoverflySimulationStub
	unit := SimulationHandler{Hoverfly: &stubHoverfly}

	request, err := http.NewRequest("OPTIONS", "/api/v2/simulation/schema", nil)
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.OptionsSchema, request)

	Expect(response.Code).To(Equal(http.StatusOK))
	Expect(response.Header().Get("Allow")).To(Equal("OPTIONS, GET"))
}

func unmarshalSimulationViewV6(buffer *bytes.Buffer) (SimulationViewV6, error) {
	body, err := ioutil.ReadAll(buffer)
	if err != nil {
		return SimulationViewV6{}, err
	}

	var simulationView SimulationViewV6

	err = json.Unmarshal(body, &simulationView)
	if err != nil {
		return SimulationViewV6{}, err
	}

	return simulationView, nil
}

func unmarshalResultView(buffer *bytes.Buffer) (SimulationImportResult, error) {
	body, err := ioutil.ReadAll(buffer)
	if err != nil {
		return SimulationImportResult{}, err
	}

	var result SimulationImportResult

	err = json.Unmarshal(body, &result)
	if err != nil {
		return SimulationImportResult{}, err
	}

	return result, nil
}

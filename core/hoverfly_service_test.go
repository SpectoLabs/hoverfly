package hoverfly

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SpectoLabs/hoverfly/core/handlers/v1"
	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/matching/matchers"
	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/SpectoLabs/hoverfly/core/modes"
	"github.com/gorilla/mux"
	. "github.com/onsi/gomega"
)

var (
	pairOne = v2.RequestMatcherResponsePairViewV5{
		RequestMatcher: v2.RequestMatcherViewV5{
			Destination: []v2.MatcherViewV5{
				v2.NewMatcherView(matchers.Exact, "test.com"),
			},
			Path: []v2.MatcherViewV5{
				v2.NewMatcherView(matchers.Exact, "/testing"),
			},
		},
		Response: v2.ResponseDetailsViewV5{
			Body: "test-body",
		},
	}

	pairTwo = v2.RequestMatcherResponsePairViewV5{
		RequestMatcher: v2.RequestMatcherViewV5{
			Path: []v2.MatcherViewV5{
				{
					Matcher: matchers.Exact,
					Value:   "/path",
				},
			},
		},
		Response: v2.ResponseDetailsViewV5{
			Body: "pair2-body",
		},
	}

	delayOne = v1.ResponseDelayView{
		UrlPattern: ".",
		HttpMethod: "GET",
		Delay:      200,
	}

	delayTwo = v1.ResponseDelayView{
		UrlPattern: "test.com",
		Delay:      201,
	}

	delayLogNormalOne = v1.ResponseDelayLogNormalView{
		UrlPattern: ".",
		HttpMethod: "GET",
		Min:        100,
		Max:        400,
		Mean:       300,
		Median:     200,
	}

	delayLogNormalTwo = v1.ResponseDelayLogNormalView{
		UrlPattern: "test.com",
		Min:        101,
		Max:        401,
		Mean:       301,
		Median:     201,
	}
)

func processHandlerOkay(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)

	var newPairView v2.RequestResponsePairViewV1

	json.Unmarshal(body, &newPairView)

	newPairView.Response.Body = "You got straight up messed with"

	pairViewBytes, _ := json.Marshal(newPairView)
	w.Write(pairViewBytes)
}

func Test_Hoverfly_SetDestination_SetDestination(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	unit.Cfg.ProxyPort = "5556"
	err := unit.StartProxy()
	Expect(err).To(BeNil())
	unit.SetDestination("newdest")

	Expect(unit.Cfg.Destination).To(Equal("newdest"))
}

func Test_Hoverfly_SetDestination_UpdateDestinationEmpty(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	unit.Cfg.ProxyPort = "5557"
	unit.StartProxy()
	err := unit.SetDestination("e^^**#")
	Expect(err).ToNot(BeNil())
}

func Test_Hoverfly_GetSimulation_ReturnsBlankSimulation_ifThereIsNoData(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	simulation, err := unit.GetSimulation()
	Expect(err).To(BeNil())

	Expect(simulation.RequestResponsePairs).To(HaveLen(0))
	Expect(simulation.GlobalActions.Delays).To(HaveLen(0))

	Expect(simulation.MetaView.SchemaVersion).To(Equal("v5"))
	Expect(simulation.MetaView.HoverflyVersion).To(MatchRegexp(`v\d+.\d+.\d+(-rc.\d)*`))
	Expect(simulation.MetaView.TimeExported).ToNot(BeNil())
}

func Test_Hoverfly_GetSimulation_ReturnsASingleRequestResponsePair(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	unit.Simulation.AddPair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "test.com",
				},
			},
		},
		Response: models.ResponseDetails{
			Status: 200,
			Body:   "test-body",
		},
	})

	simulation, err := unit.GetSimulation()
	Expect(err).To(BeNil())

	Expect(simulation.RequestResponsePairs).To(HaveLen(1))

	Expect(simulation.RequestResponsePairs[0].RequestMatcher.Destination).To(HaveLen(1))
	Expect(simulation.RequestResponsePairs[0].RequestMatcher.Destination[0].Matcher).To(Equal("exact"))
	Expect(simulation.RequestResponsePairs[0].RequestMatcher.Destination[0].Value).To(Equal("test.com"))
	Expect(simulation.RequestResponsePairs[0].RequestMatcher.Path).To(BeNil())
	Expect(simulation.RequestResponsePairs[0].RequestMatcher.Method).To(BeNil())
	Expect(simulation.RequestResponsePairs[0].RequestMatcher.DeprecatedQuery).To(BeNil())
	Expect(simulation.RequestResponsePairs[0].RequestMatcher.Scheme).To(BeNil())
	Expect(simulation.RequestResponsePairs[0].RequestMatcher.Headers).To(HaveLen(0))

	Expect(simulation.RequestResponsePairs[0].Response.Status).To(Equal(200))
	Expect(simulation.RequestResponsePairs[0].Response.EncodedBody).To(BeFalse())
	Expect(simulation.RequestResponsePairs[0].Response.Body).To(Equal("test-body"))
	Expect(simulation.RequestResponsePairs[0].Response.Headers).To(HaveLen(0))

	Expect(nil).To(BeNil())
}

func Test_Hoverfly_GetSimulation_ReturnsMultipleRequestResponsePairs(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	unit.Simulation.AddPair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "testhost-0.com",
				},
			},
			Path: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "/test",
				},
			},
		},
		Response: models.ResponseDetails{
			Status: 200,
			Body:   "test",
		},
	})

	unit.Simulation.AddPair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "testhost-1.com",
				},
			},
			Path: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "/test",
				},
			},
		},
		Response: models.ResponseDetails{
			Status: 200,
			Body:   "test",
		},
	})

	simulation, err := unit.GetSimulation()
	Expect(err).To(BeNil())

	Expect(simulation.DataViewV5.RequestResponsePairs).To(HaveLen(2))

	Expect(simulation.RequestResponsePairs[0].RequestMatcher.Destination[0].Matcher).To(Equal("exact"))
	Expect(simulation.RequestResponsePairs[0].RequestMatcher.Destination[0].Value).To(Equal("testhost-0.com"))
	Expect(simulation.RequestResponsePairs[0].RequestMatcher.Path[0].Matcher).To(Equal("exact"))
	Expect(simulation.RequestResponsePairs[0].RequestMatcher.Path[0].Value).To(Equal("/test"))

	Expect(simulation.DataViewV5.RequestResponsePairs[0].Response.Status).To(Equal(200))
	Expect(simulation.DataViewV5.RequestResponsePairs[0].Response.Body).To(Equal("test"))

	Expect(simulation.RequestResponsePairs[1].RequestMatcher.Destination[0].Matcher).To(Equal("exact"))
	Expect(simulation.RequestResponsePairs[1].RequestMatcher.Destination[0].Value).To(Equal("testhost-1.com"))
	Expect(simulation.RequestResponsePairs[1].RequestMatcher.Path[0].Matcher).To(Equal("exact"))
	Expect(simulation.RequestResponsePairs[1].RequestMatcher.Path[0].Value).To(Equal("/test"))

	Expect(simulation.DataViewV5.RequestResponsePairs[1].Response.Status).To(Equal(200))
	Expect(simulation.DataViewV5.RequestResponsePairs[1].Response.Body).To(Equal("test"))
}

func Test_Hoverfly_GetSimulation_ReturnsMultipleDelays(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	delay1 := models.ResponseDelay{
		UrlPattern: "test-pattern",
		Delay:      100,
	}

	delay2 := models.ResponseDelay{
		HttpMethod: "test",
		Delay:      200,
	}

	responseDelays := models.ResponseDelayList{delay1, delay2}

	unit.Simulation.ResponseDelays = &responseDelays

	simulation, err := unit.GetSimulation()
	Expect(err).To(BeNil())

	Expect(simulation.DataViewV5.GlobalActions.Delays).To(HaveLen(2))

	Expect(simulation.DataViewV5.GlobalActions.Delays[0].UrlPattern).To(Equal("test-pattern"))
	Expect(simulation.DataViewV5.GlobalActions.Delays[0].HttpMethod).To(Equal(""))
	Expect(simulation.DataViewV5.GlobalActions.Delays[0].Delay).To(Equal(100))

	Expect(simulation.DataViewV5.GlobalActions.Delays[1].UrlPattern).To(Equal(""))
	Expect(simulation.DataViewV5.GlobalActions.Delays[1].HttpMethod).To(Equal("test"))
	Expect(simulation.DataViewV5.GlobalActions.Delays[1].Delay).To(Equal(200))
}

func Test_Hoverfly_GetSimulation_ReturnsMultipleDelaysLogNormal(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	delay1 := models.ResponseDelayLogNormal{
		UrlPattern: "test-pattern",
		Min:        100,
		Max:        400,
		Mean:       300,
		Median:     200,
	}

	delay2 := models.ResponseDelayLogNormal{
		HttpMethod: "test",
		Min:        101,
		Max:        401,
		Mean:       301,
		Median:     201,
	}

	responseDelays := models.ResponseDelayLogNormalList{delay1, delay2}

	unit.Simulation.ResponseDelaysLogNormal = &responseDelays

	simulation, err := unit.GetSimulation()
	Expect(err).To(BeNil())

	Expect(simulation.DataViewV5.GlobalActions.DelaysLogNormal).To(HaveLen(2))

	Expect(simulation.DataViewV5.GlobalActions.DelaysLogNormal[0].Min).To(Equal(delay1.Min))
	Expect(simulation.DataViewV5.GlobalActions.DelaysLogNormal[0].Max).To(Equal(delay1.Max))
	Expect(simulation.DataViewV5.GlobalActions.DelaysLogNormal[0].Mean).To(Equal(delay1.Mean))
	Expect(simulation.DataViewV5.GlobalActions.DelaysLogNormal[0].Median).To(Equal(delay1.Median))
	Expect(simulation.DataViewV5.GlobalActions.DelaysLogNormal[1].Min).To(Equal(delay2.Min))
	Expect(simulation.DataViewV5.GlobalActions.DelaysLogNormal[1].Max).To(Equal(delay2.Max))
	Expect(simulation.DataViewV5.GlobalActions.DelaysLogNormal[1].Mean).To(Equal(delay2.Mean))
	Expect(simulation.DataViewV5.GlobalActions.DelaysLogNormal[1].Median).To(Equal(delay2.Median))
}

func Test_Hoverfly_GetFilteredSimulation_WithPlainTextUrlQuery(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	unit.Simulation.AddPair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "foo.com",
				},
			},
		},
	})

	unit.Simulation.AddPair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "bar.com",
				},
			},
		},
	})

	simulation, err := unit.GetFilteredSimulation("bar.com")
	Expect(err).To(BeNil())

	Expect(simulation.RequestResponsePairs).To(HaveLen(1))

	Expect(simulation.RequestResponsePairs[0].RequestMatcher.Destination[0].Matcher).To(Equal("exact"))
	Expect(simulation.RequestResponsePairs[0].RequestMatcher.Destination[0].Value).To(Equal("bar.com"))
}

func Test_Hoverfly_GetFilteredSimulation_WithRegexUrlQuery(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	unit.Simulation.AddPair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "foo.com",
				},
			},
		},
	})

	unit.Simulation.AddPair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "test-1.com",
				},
			},
		},
	})

	unit.Simulation.AddPair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "test-2.com",
				},
			},
		},
	})

	simulation, err := unit.GetFilteredSimulation("test-(.+).com")
	Expect(err).To(BeNil())

	Expect(simulation.RequestResponsePairs).To(HaveLen(2))

	Expect(simulation.RequestResponsePairs[0].RequestMatcher.Destination[0].Matcher).To(Equal("exact"))
	Expect(simulation.RequestResponsePairs[0].RequestMatcher.Destination[0].Value).To(Equal("test-1.com"))
	Expect(simulation.RequestResponsePairs[1].RequestMatcher.Destination[0].Matcher).To(Equal("exact"))
	Expect(simulation.RequestResponsePairs[1].RequestMatcher.Destination[0].Value).To(Equal("test-2.com"))
}

func Test_Hoverfly_GetFilteredSimulation_ReturnBlankSimulation_IfThereIsNoMatch(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	unit.Simulation.AddPair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "foo.com",
				},
			},
		},
	})

	simulation, err := unit.GetFilteredSimulation("test-(.+).com")
	Expect(err).To(BeNil())

	Expect(simulation.RequestResponsePairs).To(HaveLen(0))
	Expect(simulation.GlobalActions.Delays).To(HaveLen(0))
	Expect(simulation.GlobalActions.DelaysLogNormal).To(HaveLen(0))

	Expect(simulation.MetaView.SchemaVersion).To(Equal("v5"))
	Expect(simulation.MetaView.HoverflyVersion).To(MatchRegexp(`v\d+.\d+.\d+(-rc.\d)*`))
	Expect(simulation.MetaView.TimeExported).ToNot(BeNil())
}

func Test_Hoverfly_GetFilteredSimulationReturnError_OnInvalidRegexQuery(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	unit.Simulation.AddPair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "foo.com",
				},
			},
		},
	})

	_, err := unit.GetFilteredSimulation("test-(.+.com")
	Expect(err).NotTo(BeNil())
}

func Test_Hoverfly_GetFilteredSimulation_WithUrlQueryContainingPath(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	unit.Simulation.AddPair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "foo.com",
				},
			},
			Path: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "/api/v1",
				},
			},
		},
	})

	unit.Simulation.AddPair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "foo.com",
				},
			},
			Path: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "/api/v2",
				},
			},
		},
	})

	unit.Simulation.AddPair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "bar.com",
				},
			},
			Path: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "/api/v1",
				},
			},
		},
	})

	simulation, err := unit.GetFilteredSimulation("foo.com/api/v1")
	Expect(err).To(BeNil())

	Expect(simulation.RequestResponsePairs).To(HaveLen(1))

	Expect(simulation.RequestResponsePairs[0].RequestMatcher.Destination[0].Matcher).To(Equal("exact"))
	Expect(simulation.RequestResponsePairs[0].RequestMatcher.Destination[0].Value).To(Equal("foo.com"))
	Expect(simulation.RequestResponsePairs[0].RequestMatcher.Path[0].Matcher).To(Equal("exact"))
	Expect(simulation.RequestResponsePairs[0].RequestMatcher.Path[0].Value).To(Equal("/api/v1"))
}

func Test_Hoverfly_PutSimulation_ImportsRecordings(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	simulationToImport := v2.SimulationViewV5{
		v2.DataViewV5{
			RequestResponsePairs: []v2.RequestMatcherResponsePairViewV5{pairOne},
			GlobalActions: v2.GlobalActionsView{
				Delays: []v1.ResponseDelayView{},
			},
		},
		v2.MetaView{},
	}

	unit.PutSimulation(simulationToImport)

	importedSimulation, err := unit.GetSimulation()
	Expect(err).To(BeNil())

	Expect(importedSimulation).ToNot(BeNil())

	Expect(importedSimulation.RequestResponsePairs).ToNot(BeNil())
	Expect(importedSimulation.RequestResponsePairs).To(HaveLen(1))

	Expect(importedSimulation.RequestResponsePairs[0].RequestMatcher.Destination[0].Matcher).To(Equal("exact"))
	Expect(importedSimulation.RequestResponsePairs[0].RequestMatcher.Destination[0].Value).To(Equal("test.com"))
	Expect(importedSimulation.RequestResponsePairs[0].RequestMatcher.Path[0].Matcher).To(Equal("exact"))
	Expect(importedSimulation.RequestResponsePairs[0].RequestMatcher.Path[0].Value).To(Equal("/testing"))

	Expect(importedSimulation.RequestResponsePairs[0].Response.Body).To(Equal("test-body"))
}

func Test_Hoverfly_PutSimulation_ImportsSimulationViews(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	simulationToImport := v2.SimulationViewV5{
		v2.DataViewV5{
			RequestResponsePairs: []v2.RequestMatcherResponsePairViewV5{pairTwo},
			GlobalActions: v2.GlobalActionsView{
				Delays: []v1.ResponseDelayView{},
			},
		},
		v2.MetaView{},
	}

	unit.PutSimulation(simulationToImport)

	importedSimulation, err := unit.GetSimulation()
	Expect(err).To(BeNil())

	Expect(importedSimulation).ToNot(BeNil())

	Expect(importedSimulation.RequestResponsePairs).ToNot(BeNil())
	Expect(importedSimulation.RequestResponsePairs).To(HaveLen(1))

	Expect(importedSimulation.RequestResponsePairs[0].RequestMatcher.Destination).To(BeNil())
	Expect(importedSimulation.RequestResponsePairs[0].RequestMatcher.Path[0].Matcher).To(Equal("exact"))
	Expect(importedSimulation.RequestResponsePairs[0].RequestMatcher.Path[0].Value).To(Equal("/path"))

	Expect(importedSimulation.RequestResponsePairs[0].Response.Body).To(Equal("pair2-body"))
}

func Test_Hoverfly_PutSimulation_ImportsDelays(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	simulationToImport := v2.SimulationViewV5{
		v2.DataViewV5{
			RequestResponsePairs: []v2.RequestMatcherResponsePairViewV5{},
			GlobalActions: v2.GlobalActionsView{
				Delays: []v1.ResponseDelayView{delayOne, delayTwo},
			},
		},
		v2.MetaView{},
	}

	err := unit.PutSimulation(simulationToImport)
	Expect(err.GetError()).To(BeNil())

	delays := unit.Simulation.ResponseDelays.ConvertToResponseDelayPayloadView()
	Expect(delays).ToNot(BeNil())

	Expect(delays.Data).To(HaveLen(2))

	Expect(delays.Data[0].UrlPattern).To(Equal("."))
	Expect(delays.Data[0].HttpMethod).To(Equal("GET"))
	Expect(delays.Data[0].Delay).To(Equal(200))

	Expect(delays.Data[1].UrlPattern).To(Equal("test.com"))
	Expect(delays.Data[1].HttpMethod).To(Equal(""))
	Expect(delays.Data[1].Delay).To(Equal(201))
}

func Test_Hoverfly_PutSimulation_ImportsDelaysLogNormal(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	simulationToImport := v2.SimulationViewV5{
		v2.DataViewV5{
			RequestResponsePairs: []v2.RequestMatcherResponsePairViewV5{},
			GlobalActions: v2.GlobalActionsView{
				DelaysLogNormal: []v1.ResponseDelayLogNormalView{delayLogNormalOne, delayLogNormalTwo},
			},
		},
		v2.MetaView{},
	}

	err := unit.PutSimulation(simulationToImport)
	Expect(err.GetError()).To(BeNil())

	delays := unit.Simulation.ResponseDelaysLogNormal.ConvertToResponseDelayLogNormalPayloadView()
	Expect(delays).ToNot(BeNil())

	Expect(delays.Data).To(HaveLen(2))

	Expect(delays.Data[0]).To(Equal(delayLogNormalOne))
	Expect(delays.Data[1]).To(Equal(delayLogNormalTwo))
}

func Test_Hoverfly_GetMiddleware_ReturnsCorrectValuesFromMiddleware(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})
	unit.Cfg.Middleware.SetBinary("python")
	unit.Cfg.Middleware.SetScript(pythonMiddlewareBasic)

	binary, script, remote := unit.GetMiddleware()
	Expect(binary).To(Equal("python"))
	Expect(script).To(Equal(pythonMiddlewareBasic))
	Expect(remote).To(Equal(""))
}

func Test_Hoverfly_GetMiddleware_ReturnsEmptyStringsWhenNeitherIsSet(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	binary, script, remote := unit.GetMiddleware()
	Expect(binary).To(Equal(""))
	Expect(script).To(Equal(""))
	Expect(remote).To(Equal(""))
}

func Test_Hoverfly_GetMiddleware_ReturnsBinaryIfJustBinarySet(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})
	unit.Cfg.Middleware.SetBinary("python")

	binary, script, remote := unit.GetMiddleware()
	Expect(binary).To(Equal("python"))
	Expect(script).To(Equal(""))
	Expect(remote).To(Equal(""))
}

func Test_Hoverfly_GetMiddleware_ReturnsRemotefJustRemoteSet(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})
	unit.Cfg.Middleware.Remote = "test.com"

	binary, script, remote := unit.GetMiddleware()
	Expect(binary).To(Equal(""))
	Expect(script).To(Equal(""))
	Expect(remote).To(Equal("test.com"))
}

func Test_Hoverfly_SetMiddleware_CanSetBinaryAndScript(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	err := unit.SetMiddleware("python", pythonMiddlewareBasic, "")
	Expect(err).To(BeNil())

	Expect(unit.Cfg.Middleware.Binary).To(Equal("python"))

	script, err := unit.Cfg.Middleware.GetScript()
	Expect(script).To(Equal(pythonMiddlewareBasic))
	Expect(err).To(BeNil())
}

func Test_Hoverfly_SetMiddleware_CanSetRemote(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/process", processHandlerOkay).Methods("POST")
	server := httptest.NewServer(muxRouter)
	defer server.Close()

	err := unit.SetMiddleware("", "", server.URL+"/process")
	Expect(err).To(BeNil())

	Expect(unit.Cfg.Middleware.Binary).To(Equal(""))

	script, _ := unit.Cfg.Middleware.GetScript()
	Expect(script).To(Equal(""))

	Expect(unit.Cfg.Middleware.Remote).To(Equal(server.URL + "/process"))
}

func Test_Hoverfly_SetMiddleware_WillErrorIfGivenBadRemote(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	err := unit.SetMiddleware("", "", "[]somemadeupwebsite*&*^&$%^")
	Expect(err).ToNot(BeNil())

	Expect(unit.Cfg.Middleware.Binary).To(Equal(""))

	script, _ := unit.Cfg.Middleware.GetScript()
	Expect(script).To(Equal(""))

	Expect(unit.Cfg.Middleware.Remote).To(Equal(""))
}

func Test_Hoverfly_SetMiddleware_WillErrorIfGivenScriptAndNoBinaryAndWillNotChangeMiddleware(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})
	unit.Cfg.Middleware.SetBinary("python")
	unit.Cfg.Middleware.SetScript("test-script")

	err := unit.SetMiddleware("", pythonMiddlewareBasic, "")
	Expect(err).ToNot(BeNil())

	Expect(unit.Cfg.Middleware.Binary).To(Equal("python"))

	script, _ := unit.Cfg.Middleware.GetScript()
	Expect(script).To(Equal("test-script"))
}

func Test_Hoverfly_SetMiddleware_WillDeleteMiddlewareSettingsIfEmptyBinaryAndScriptAndRemote(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})
	unit.Cfg.Middleware.SetBinary("python")
	unit.Cfg.Middleware.SetScript("test-script")

	err := unit.SetMiddleware("", "", "")
	Expect(err).To(BeNil())

	Expect(unit.Cfg.Middleware.Binary).To(Equal(""))

	script, _ := unit.Cfg.Middleware.GetScript()
	Expect(script).To(Equal(""))
}

func Test_Hoverfly_SetMiddleware_WontSetMiddlewareIfCannotRunScript(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	err := unit.SetMiddleware("python", "ewfaet4rafgre", "")
	Expect(err).ToNot(BeNil())

	Expect(unit.Cfg.Middleware.Binary).To(Equal(""))

	script, _ := unit.Cfg.Middleware.GetScript()
	Expect(script).To(Equal(""))
}

func Test_Hoverfly_SetMiddleware_WillSetBinaryWithNoScript(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	err := unit.SetMiddleware("cat", "", "")
	Expect(err).To(BeNil())

	Expect(unit.Cfg.Middleware.Binary).To(Equal("cat"))

	script, _ := unit.Cfg.Middleware.GetScript()
	Expect(script).To(Equal(""))
}

func Test_Hoverfly_GetVersion_GetsVersion(t *testing.T) {
	RegisterTestingT(t)

	unit := Hoverfly{
		version: "test-version",
	}

	Expect(unit.GetVersion()).To(Equal("test-version"))
}

func Test_Hoverfly_GetUpstreamProxy_GetsUpstreamProxy(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{
		UpstreamProxy: "upstream-proxy.org",
	})

	Expect(unit.GetUpstreamProxy()).To(Equal("upstream-proxy.org"))
}

func Test_Hoverfly_IsWebServer_GetsIsWebServer(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{
		Webserver: true,
	})

	Expect(unit.IsWebServer()).To(BeTrue())
}

func Test_Hoverfly_SetModeWithArguments_CanSetModeToCapture(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	Expect(unit.SetModeWithArguments(
		v2.ModeView{
			Mode: "capture",
		})).To(BeNil())
	Expect(unit.Cfg.Mode).To(Equal("capture"))
}

func Test_Hoverfly_SetModeWithArguments_CanSetModeToSimulate(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	Expect(unit.SetModeWithArguments(
		v2.ModeView{
			Mode: "simulate",
		})).To(BeNil())
	Expect(unit.Cfg.Mode).To(Equal("simulate"))
}

func Test_Hoverfly_SetModeWithArguments_CanSetModeToModify(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	Expect(unit.SetModeWithArguments(
		v2.ModeView{
			Mode: "modify",
		})).To(BeNil())
	Expect(unit.Cfg.Mode).To(Equal("modify"))
}

func Test_Hoverfly_SetModeWithArguments_CanSetModeToSynthesize(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	Expect(unit.SetModeWithArguments(
		v2.ModeView{
			Mode: "synthesize",
		})).To(BeNil())
	Expect(unit.Cfg.Mode).To(Equal("synthesize"))
}

func Test_Hoverfly_SetModeWithArguments_CannotSetModeToSomethingInvalid(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	Expect(unit.SetModeWithArguments(
		v2.ModeView{
			Mode: "mode",
		})).ToNot(BeNil())
	Expect(unit.Cfg.Mode).To(Equal(""))

	Expect(unit.SetModeWithArguments(
		v2.ModeView{
			Mode: "hoverfly",
		})).ToNot(BeNil())
	Expect(unit.Cfg.Mode).To(Equal(""))
}

func Test_Hoverfly_SetModeWithArguments_SettingModeToCaptureWipesCache(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	unit.CacheMatcher.RequestCache.Set("test", "test_bytes")

	Expect(unit.SetModeWithArguments(
		v2.ModeView{
			Mode: "capture",
		})).To(BeNil())
	Expect(unit.Cfg.Mode).To(Equal("capture"))

	Expect(unit.CacheMatcher.RequestCache.RecordsCount()).To(Equal(0))
}

func Test_Hoverfly_SetModeWithArguments_Stateful(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	Expect(unit.SetModeWithArguments(v2.ModeView{
		Mode: "capture",
		Arguments: v2.ModeArgumentsView{
			Stateful: true,
		},
	})).To(Succeed())

	storedMode := unit.modeMap[modes.Capture].View()
	Expect(storedMode.Arguments.Stateful).To(BeTrue())
}

func Test_Hoverfly_SetModeWithArguments_AsteriskCanOnlyBeValidAsTheOnlyHeader(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	Expect(unit.SetModeWithArguments(
		v2.ModeView{
			Mode: "capture",
		})).To(BeNil())
	Expect(unit.Cfg.Mode).To(Equal("capture"))

	Expect(unit.SetModeWithArguments(v2.ModeView{
		Arguments: v2.ModeArgumentsView{
			Headers: []string{"Content-Type", "*"},
		},
	})).ToNot(Succeed())

	Expect(unit.SetModeWithArguments(v2.ModeView{
		Arguments: v2.ModeArgumentsView{
			Headers: []string{"*"},
		},
	})).ToNot(Succeed())
}

func Test_Hoverfly_AddDiff_AddEntry(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})
	Expect(unit.responsesDiff).To(HaveLen(0))

	key := v2.SimpleRequestDefinitionView{
		Host: "test.com",
	}

	unit.AddDiff(key, v2.DiffReport{Timestamp: "now", DiffEntries: []v2.DiffReportEntry{{}}})

	Expect(unit.responsesDiff).To(HaveLen(1))

	diffReports := unit.responsesDiff[key]
	Expect(diffReports).To(HaveLen(1))
	Expect(diffReports[0].Timestamp).To(Equal("now"))
}

func Test_Hoverfly_AddDiff_AppendsEntry(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})
	Expect(unit.responsesDiff).To(HaveLen(0))

	key := v2.SimpleRequestDefinitionView{
		Host: "test.com",
	}

	unit.AddDiff(key, v2.DiffReport{Timestamp: "now", DiffEntries: []v2.DiffReportEntry{{Actual: "1"}}})
	unit.AddDiff(key, v2.DiffReport{Timestamp: "now", DiffEntries: []v2.DiffReportEntry{{Actual: "2"}}})

	Expect(unit.responsesDiff).To(HaveLen(1))

	diffReports := unit.responsesDiff[key]
	Expect(diffReports).To(HaveLen(2))
	Expect(diffReports[0].DiffEntries[0].Actual).To(Equal("1"))
	Expect(diffReports[1].DiffEntries[0].Actual).To(Equal("2"))
}

func Test_Hoverfly_AddDiff_AddEntry_DiffrentKey(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})
	Expect(unit.responsesDiff).To(HaveLen(0))

	key := v2.SimpleRequestDefinitionView{
		Host: "test.com",
	}

	keyTwo := v2.SimpleRequestDefinitionView{
		Method: "POST",
		Host:   "test.com",
	}

	unit.AddDiff(key, v2.DiffReport{Timestamp: "now", DiffEntries: []v2.DiffReportEntry{{Actual: "1"}}})
	unit.AddDiff(keyTwo, v2.DiffReport{Timestamp: "now", DiffEntries: []v2.DiffReportEntry{{Actual: "2"}}})

	Expect(unit.responsesDiff).To(HaveLen(2))

	diffReports := unit.responsesDiff[key]
	Expect(diffReports).To(HaveLen(1))
	Expect(diffReports[0].DiffEntries[0].Actual).To(Equal("1"))

	diffReports = unit.responsesDiff[keyTwo]
	Expect(diffReports).To(HaveLen(1))
	Expect(diffReports[0].DiffEntries[0].Actual).To(Equal("2"))
}

func Test_Hoverfly_AddDiff_DoesntAddDiffReport_NoEntries(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})
	Expect(unit.responsesDiff).To(HaveLen(0))

	key := v2.SimpleRequestDefinitionView{
		Host: "test.com",
	}

	unit.AddDiff(key, v2.DiffReport{Timestamp: "now"})

	Expect(unit.responsesDiff).To(HaveLen(0))
}

func Test_Hoverfly_GetPACFile_GetsPACFile(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{
		PACFile: []byte("PACFILE"),
	})

	Expect(string(unit.GetPACFile())).To(Equal("PACFILE"))
}

func Test_Hoverfly_SetPACFile_SetsPACFile(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	unit.SetPACFile([]byte("PACFILE"))

	Expect(string(unit.Cfg.PACFile)).To(Equal("PACFILE"))
}

func Test_Hoverfly_SetPACFile_SetsPACFileToNilIfEmpty(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{
		PACFile: []byte("PACFILE"),
	})

	unit.SetPACFile([]byte(""))

	Expect(unit.Cfg.PACFile).To(BeNil())
}

func Test_Hoverfly_DeletePACFile(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{
		PACFile: []byte("PACFILE"),
	})

	unit.DeletePACFile()

	Expect(unit.Cfg.PACFile).To(BeNil())
}

package hoverfly

import (
	"net/http/httptest"
	"testing"

	"github.com/SpectoLabs/hoverfly/core/handlers/v1"
	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/SpectoLabs/hoverfly/core/util"
	"github.com/gorilla/mux"
	. "github.com/onsi/gomega"
)

var (
	pairOneRecording = v2.RequestResponsePairViewV1{
		Request: v2.RequestDetailsViewV1{
			Destination: util.StringToPointer("test.com"),
			Path:        util.StringToPointer("/testing"),
		},
		Response: v2.ResponseDetailsView{
			Body: "test-body",
		},
	}

	pairOneTemplate = v2.RequestResponsePairViewV1{
		Request: v2.RequestDetailsViewV1{
			Path: util.StringToPointer("/template"),
		},
		Response: v2.ResponseDetailsView{
			Body: "template-body",
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
)

func TestHoverflyGetSimulationReturnsBlankSimulation_ifThereIsNoData(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	simulation, err := unit.GetSimulation()
	Expect(err).To(BeNil())

	Expect(simulation.DataViewV1.RequestResponsePairViewV1).To(HaveLen(0))
	Expect(simulation.DataViewV1.GlobalActions.Delays).To(HaveLen(0))

	Expect(simulation.MetaView.SchemaVersion).To(Equal("v1"))
	Expect(simulation.MetaView.HoverflyVersion).To(MatchRegexp(`v\d+.\d+.\d+`))
	Expect(simulation.MetaView.TimeExported).ToNot(BeNil())
}

func TestHoverfly_GetSimulation_ReturnsASingleRequestResponsePairTemplate(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	unit.Simulation.Templates = append(unit.Simulation.Templates, models.RequestTemplateResponsePair{
		RequestTemplate: models.RequestTemplate{
			Destination: util.StringToPointer("test.com"),
		},
		Response: models.ResponseDetails{
			Status: 200,
			Body:   "test-template",
		},
	})

	simulation, err := unit.GetSimulation()
	Expect(err).To(BeNil())

	Expect(simulation.DataViewV1.RequestResponsePairViewV1).To(HaveLen(1))

	Expect(*simulation.DataViewV1.RequestResponsePairViewV1[0].Request.Destination).To(Equal("test.com"))
	Expect(simulation.DataViewV1.RequestResponsePairViewV1[0].Request.Path).To(BeNil())
	Expect(simulation.DataViewV1.RequestResponsePairViewV1[0].Request.Method).To(BeNil())
	Expect(simulation.DataViewV1.RequestResponsePairViewV1[0].Request.Query).To(BeNil())
	Expect(simulation.DataViewV1.RequestResponsePairViewV1[0].Request.Scheme).To(BeNil())
	Expect(simulation.DataViewV1.RequestResponsePairViewV1[0].Request.Headers).To(HaveLen(0))

	Expect(simulation.DataViewV1.RequestResponsePairViewV1[0].Response.Status).To(Equal(200))
	Expect(simulation.DataViewV1.RequestResponsePairViewV1[0].Response.EncodedBody).To(BeFalse())
	Expect(simulation.DataViewV1.RequestResponsePairViewV1[0].Response.Body).To(Equal("test-template"))
	Expect(simulation.DataViewV1.RequestResponsePairViewV1[0].Response.Headers).To(HaveLen(0))

	Expect(nil).To(BeNil())
}

func Test_Hoverfly_GetSimulation_ReturnsMultipleRequestResponsePairs(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	unit.Simulation.Templates = append(unit.Simulation.Templates, models.RequestTemplateResponsePair{
		RequestTemplate: models.RequestTemplate{
			Destination: util.StringToPointer("testhost.com"),
			Path: &models.RequestFieldMatchers{
				ExactMatch: util.StringToPointer("/test"),
			},
		},
		Response: models.ResponseDetails{
			Status: 200,
			Body:   "test",
		},
	})

	unit.Simulation.Templates = append(unit.Simulation.Templates, models.RequestTemplateResponsePair{
		RequestTemplate: models.RequestTemplate{
			Destination: util.StringToPointer("testhost.com"),
			Path: &models.RequestFieldMatchers{
				ExactMatch: util.StringToPointer("/test"),
			},
		},
		Response: models.ResponseDetails{
			Status: 200,
			Body:   "test",
		},
	})

	simulation, err := unit.GetSimulation()
	Expect(err).To(BeNil())

	Expect(simulation.DataViewV1.RequestResponsePairViewV1).To(HaveLen(2))

	Expect(*simulation.DataViewV1.RequestResponsePairViewV1[0].Request.Destination).To(Equal("testhost.com"))
	Expect(*simulation.DataViewV1.RequestResponsePairViewV1[0].Request.Path).To(Equal("/test"))

	Expect(simulation.DataViewV1.RequestResponsePairViewV1[0].Response.Status).To(Equal(200))
	Expect(simulation.DataViewV1.RequestResponsePairViewV1[0].Response.Body).To(Equal("test"))

	Expect(*simulation.DataViewV1.RequestResponsePairViewV1[1].Request.Destination).To(Equal("testhost.com"))
	Expect(*simulation.DataViewV1.RequestResponsePairViewV1[1].Request.Path).To(Equal("/test"))

	Expect(simulation.DataViewV1.RequestResponsePairViewV1[1].Response.Status).To(Equal(200))
	Expect(simulation.DataViewV1.RequestResponsePairViewV1[1].Response.Body).To(Equal("test"))
}

func TestHoverflyGetSimulationReturnsMultipleDelays(t *testing.T) {
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

	Expect(simulation.DataViewV1.GlobalActions.Delays).To(HaveLen(2))

	Expect(simulation.DataViewV1.GlobalActions.Delays[0].UrlPattern).To(Equal("test-pattern"))
	Expect(simulation.DataViewV1.GlobalActions.Delays[0].HttpMethod).To(Equal(""))
	Expect(simulation.DataViewV1.GlobalActions.Delays[0].Delay).To(Equal(100))

	Expect(simulation.DataViewV1.GlobalActions.Delays[1].UrlPattern).To(Equal(""))
	Expect(simulation.DataViewV1.GlobalActions.Delays[1].HttpMethod).To(Equal("test"))
	Expect(simulation.DataViewV1.GlobalActions.Delays[1].Delay).To(Equal(200))
}

func TestHoverfly_PutSimulation_ImportsRecordings(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	simulationToImport := v2.SimulationViewV1{
		v2.DataViewV1{
			RequestResponsePairViewV1: []v2.RequestResponsePairViewV1{pairOneRecording},
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

	Expect(importedSimulation.RequestResponsePairViewV1).ToNot(BeNil())
	Expect(importedSimulation.RequestResponsePairViewV1).To(HaveLen(1))

	Expect(importedSimulation.RequestResponsePairViewV1[0].Request.Destination).To(Equal(util.StringToPointer("test.com")))
	Expect(importedSimulation.RequestResponsePairViewV1[0].Request.Path).To(Equal(util.StringToPointer("/testing")))

	Expect(importedSimulation.RequestResponsePairViewV1[0].Response.Body).To(Equal("test-body"))
}

func TestHoverfly_PutSimulation_ImportsTemplates(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	simulationToImport := v2.SimulationViewV1{
		v2.DataViewV1{
			RequestResponsePairViewV1: []v2.RequestResponsePairViewV1{pairOneTemplate},
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

	Expect(importedSimulation.RequestResponsePairViewV1).ToNot(BeNil())
	Expect(importedSimulation.RequestResponsePairViewV1).To(HaveLen(1))

	Expect(importedSimulation.RequestResponsePairViewV1[0].Request.Destination).To(BeNil())
	Expect(importedSimulation.RequestResponsePairViewV1[0].Request.Path).To(Equal(util.StringToPointer("/template")))

	Expect(importedSimulation.RequestResponsePairViewV1[0].Response.Body).To(Equal("template-body"))
}

func TestHoverfly_PutSimulation_ImportsRecordingsAndTemplates(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	simulationToImport := v2.SimulationViewV1{
		v2.DataViewV1{
			RequestResponsePairViewV1: []v2.RequestResponsePairViewV1{pairOneRecording, pairOneTemplate},
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

	Expect(importedSimulation.RequestResponsePairViewV1).ToNot(BeNil())
	Expect(importedSimulation.RequestResponsePairViewV1).To(HaveLen(2))

	Expect(importedSimulation.RequestResponsePairViewV1[0].Request.Destination).To(Equal(util.StringToPointer("test.com")))
	Expect(importedSimulation.RequestResponsePairViewV1[0].Request.Path).To(Equal(util.StringToPointer("/testing")))

	Expect(importedSimulation.RequestResponsePairViewV1[0].Response.Body).To(Equal("test-body"))

	Expect(importedSimulation.RequestResponsePairViewV1[1].Request.Destination).To(BeNil())
	Expect(importedSimulation.RequestResponsePairViewV1[1].Request.Path).To(Equal(util.StringToPointer("/template")))

	Expect(importedSimulation.RequestResponsePairViewV1[1].Response.Body).To(Equal("template-body"))
}

func TestHoverfly_PutSimulation_ImportsDelays(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	simulationToImport := v2.SimulationViewV1{
		v2.DataViewV1{
			RequestResponsePairViewV1: []v2.RequestResponsePairViewV1{},
			GlobalActions: v2.GlobalActionsView{
				Delays: []v1.ResponseDelayView{delayOne, delayTwo},
			},
		},
		v2.MetaView{},
	}

	err := unit.PutSimulation(simulationToImport)
	Expect(err).To(BeNil())

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

func Test_Hoverfly_SetMiddleware_WillErrorIfGivenBadBinaryAndWillNotChangeMiddleware(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})
	unit.Cfg.Middleware.SetBinary("python")
	unit.Cfg.Middleware.SetScript("test-script")

	err := unit.SetMiddleware("this-isnt-a-binary", pythonMiddlewareBasic, "")
	Expect(err).ToNot(BeNil())

	Expect(unit.Cfg.Middleware.Binary).To(Equal("python"))

	script, _ := unit.Cfg.Middleware.GetScript()
	Expect(script).To(Equal("test-script"))
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

func Test_Hoverfly_SetMode_CanSetModeToCapture(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	Expect(unit.SetMode("capture")).To(BeNil())
	Expect(unit.Cfg.Mode).To(Equal("capture"))
}

func Test_Hoverfly_SetMode_CanSetModeToSimulate(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	Expect(unit.SetMode("simulate")).To(BeNil())
	Expect(unit.Cfg.Mode).To(Equal("simulate"))
}

func Test_Hoverfly_SetMode_CanSetModeToModify(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	Expect(unit.SetMode("modify")).To(BeNil())
	Expect(unit.Cfg.Mode).To(Equal("modify"))
}

func Test_Hoverfly_SetMode_CanSetModeToSynthesize(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	Expect(unit.SetMode("synthesize")).To(BeNil())
	Expect(unit.Cfg.Mode).To(Equal("synthesize"))
}

func Test_Hoverfly_SetMode_CannotSetModeToSomethingInvalid(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	Expect(unit.SetMode("mode")).ToNot(BeNil())
	Expect(unit.Cfg.Mode).To(Equal(""))

	Expect(unit.SetMode("hoverfly")).ToNot(BeNil())
	Expect(unit.Cfg.Mode).To(Equal(""))
}

func Test_Hoverfly_SetMode_SettingModeToCaptureWipesCache(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	unit.CacheMatcher.RequestCache.Set([]byte("test"), []byte("test_bytes"))

	Expect(unit.SetMode("capture")).To(BeNil())
	Expect(unit.Cfg.Mode).To(Equal("capture"))

	values, _ := unit.CacheMatcher.RequestCache.GetAllValues()
	Expect(values).To(HaveLen(0))
}

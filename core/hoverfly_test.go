package hoverfly

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/SpectoLabs/hoverfly/core/authentication/backends"
	"github.com/SpectoLabs/hoverfly/core/cache"
	"github.com/SpectoLabs/hoverfly/core/handlers/v1"
	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/SpectoLabs/hoverfly/core/util"
	. "github.com/onsi/gomega"
)

func TestGetNewHoverflyCheckConfig(t *testing.T) {
	RegisterTestingT(t)

	cfg := InitSettings()

	db := cache.GetDB("testing2.db")
	requestCache := cache.NewBoltDBCache(db, []byte("requestBucket"))
	metaCache := cache.NewBoltDBCache(db, []byte("metaBucket"))
	tokenCache := cache.NewBoltDBCache(db, []byte("tokenBucket"))
	userCache := cache.NewBoltDBCache(db, []byte("userBucket"))
	backend := backends.NewCacheBasedAuthBackend(tokenCache, userCache)

	dbClient := GetNewHoverfly(cfg, requestCache, metaCache, backend)

	Expect(dbClient.Cfg).To(Equal(cfg))

	// deleting this database
	os.Remove("testing2.db")
}

func TestGetNewHoverfly(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	unit.Cfg.ProxyPort = "6666"

	err := unit.StartProxy()
	Expect(err).To(BeNil())

	newResponse, err := http.Get(fmt.Sprintf("http://localhost:%s/", unit.Cfg.ProxyPort))
	Expect(err).To(BeNil())
	Expect(newResponse.StatusCode).To(Equal(http.StatusInternalServerError))

}

func Test_Hoverfly_processRequest_CaptureModeReturnsResponseAndSavesIt(t *testing.T) {
	RegisterTestingT(t)

	server, unit := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	r, err := http.NewRequest("GET", "http://somehost.com", nil)
	Expect(err).To(BeNil())

	unit.Cfg.SetMode("capture")

	resp := unit.processRequest(r)

	Expect(resp).ToNot(BeNil())
	Expect(resp.StatusCode).To(Equal(http.StatusCreated))

	Expect(unit.Simulation.Templates).To(HaveLen(1))
}

func Test_Hoverfly_processRequest_CanSimulateRequest(t *testing.T) {
	RegisterTestingT(t)

	server, unit := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	r, err := http.NewRequest("GET", "http://somehost.com", nil)
	Expect(err).To(BeNil())

	// capturing
	unit.Cfg.SetMode("capture")
	resp := unit.processRequest(r)

	Expect(resp).ToNot(BeNil())
	Expect(resp.StatusCode).To(Equal(http.StatusCreated))

	// virtualizing
	unit.Cfg.SetMode("simulate")
	newResp := unit.processRequest(r)

	Expect(newResp).ToNot(BeNil())
	Expect(newResp.StatusCode).To(Equal(http.StatusCreated))
}

func Test_Hoverfly_processRequest_CanUseMiddlewareToSynthesizeRequest(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	// getting reflect middleware
	err := dbClient.Cfg.Middleware.SetBinary("python")
	Expect(err).To(BeNil())

	err = dbClient.Cfg.Middleware.SetScript(pythonReflectBody)
	Expect(err).To(BeNil())

	bodyBytes := []byte("request_body_here")

	r, err := http.NewRequest("GET", "http://somehost.com", ioutil.NopCloser(bytes.NewBuffer(bodyBytes)))
	Expect(err).To(BeNil())

	dbClient.Cfg.SetMode("synthesize")
	newResp := dbClient.processRequest(r)

	Expect(newResp).ToNot(BeNil())
	Expect(newResp.StatusCode).To(Equal(http.StatusCreated))
	b, err := ioutil.ReadAll(newResp.Body)
	Expect(err).To(BeNil())
	Expect(string(b)).To(Equal(string(bodyBytes)))
}

func Test_Hoverfly_processRequest_CanModifyRequest(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	err := dbClient.Cfg.Middleware.SetBinary("python")
	Expect(err).To(BeNil())

	err = dbClient.Cfg.Middleware.SetScript(pythonModifyResponse)
	Expect(err).To(BeNil())

	r, err := http.NewRequest("POST", "http://somehost.com", nil)
	Expect(err).To(BeNil())

	dbClient.Cfg.SetMode("modify")
	newResp := dbClient.processRequest(r)

	Expect(newResp).ToNot(BeNil())

	Expect(newResp.StatusCode).To(Equal(http.StatusCreated))
	Expect(newResp.Header).To(HaveKeyWithValue("Hoverfly", []string{"Was-Here"}))
}

func Test_Hoverfly_GetResponse_CanReturnResponseFromCache(t *testing.T) {
	RegisterTestingT(t)

	server, unit := testTools(201, `{'message': 'here'}`)
	server.Close()

	unit.CacheMatcher.SaveRequestResponsePair(&models.RequestResponsePair{
		Request: models.RequestDetails{
			Destination: "somehost.com",
			Method:      "POST",
			Scheme:      "http",
		},
		Response: models.ResponseDetails{
			Status: 200,
			Body:   "cached response",
		},
	})

	response, err := unit.GetResponse(models.RequestDetails{
		Destination: "somehost.com",
		Method:      "POST",
		Scheme:      "http",
	})

	Expect(err).To(BeNil())
	Expect(response).ToNot(BeNil())

	Expect(response.Status).To(Equal(http.StatusOK))
	Expect(response.Body).To(Equal("cached response"))
}

func Test_Hoverfly_GetResponse_CanReturnResponseFromSimulationAndNotCache(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	unit.Simulation.AddRequestTemplateResponsePair(&models.RequestTemplateResponsePair{
		RequestTemplate: models.RequestTemplate{
			Destination: util.StringToPointer("somehost.com"),
			Method:      util.StringToPointer("POST"),
			Scheme:      util.StringToPointer("http"),
		},
		Response: models.ResponseDetails{
			Status: 200,
			Body:   "template response",
		},
	})

	response, err := unit.GetResponse(models.RequestDetails{
		Destination: "somehost.com",
		Method:      "POST",
		Scheme:      "http",
	})

	Expect(err).To(BeNil())
	Expect(response).ToNot(BeNil())

	Expect(response.Status).To(Equal(http.StatusOK))
	Expect(response.Body).To(Equal("template response"))
}

func Test_Hoverfly_GetResponse_WillCacheResponseIfNotInCache(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	unit.Simulation.AddRequestTemplateResponsePair(&models.RequestTemplateResponsePair{
		RequestTemplate: models.RequestTemplate{
			Destination: util.StringToPointer("somehost.com"),
			Method:      util.StringToPointer("POST"),
			Scheme:      util.StringToPointer("http"),
		},
		Response: models.ResponseDetails{
			Status: 200,
			Body:   "template response",
		},
	})

	unit.GetResponse(models.RequestDetails{
		Destination: "somehost.com",
		Method:      "POST",
		Scheme:      "http",
	})

	Expect(unit.RequestCache.RecordsCount()).Should(Equal(1))

	pairBytes, err := unit.RequestCache.Get([]byte("75b4ae6efa2a3f6d3ee6b9fed4d8c8c5"))
	Expect(err).To(BeNil())

	cachedRequestResponsePair, err := models.NewRequestResponsePairFromBytes(pairBytes)
	Expect(err).To(BeNil())

	Expect(cachedRequestResponsePair.Response.Body).To(Equal("template response"))

	unit.Simulation = models.NewSimulation()
	response, err := unit.GetResponse(models.RequestDetails{
		Destination: "somehost.com",
		Method:      "POST",
		Scheme:      "http",
	})

	Expect(err).To(BeNil())
	Expect(response).ToNot(BeNil())

	Expect(response.Status).To(Equal(http.StatusOK))
	Expect(response.Body).To(Equal("template response"))
}

type ResponseDelayListStub struct {
	gotDelays int
}

func (this *ResponseDelayListStub) Json() []byte {
	return nil
}

func (this *ResponseDelayListStub) Len() int {
	return this.Len()
}

func (this *ResponseDelayListStub) GetDelay(request models.RequestDetails) *models.ResponseDelay {
	this.gotDelays++
	return nil
}

func (this ResponseDelayListStub) ConvertToResponseDelayPayloadView() v1.ResponseDelayPayloadView {
	return v1.ResponseDelayPayloadView{}
}

func TestDelayAppliedToSuccessfulSimulateRequest(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	r, err := http.NewRequest("GET", "http://somehost.com", nil)
	Expect(err).To(BeNil())

	// capturing
	dbClient.Cfg.SetMode("capture")
	resp := dbClient.processRequest(r)

	Expect(resp.StatusCode).To(Equal(http.StatusCreated))

	// virtualizing
	dbClient.Cfg.SetMode("simulate")

	stub := ResponseDelayListStub{}
	dbClient.ResponseDelays = &stub

	newResp := dbClient.processRequest(r)

	Expect(newResp.StatusCode).To(Equal(http.StatusCreated))

	Expect(stub.gotDelays, Equal(1))
}

func TestDelayNotAppliedToFailedSimulateRequest(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	r, err := http.NewRequest("GET", "http://somehost.com", nil)
	Expect(err).To(BeNil())

	// virtualizing
	dbClient.Cfg.SetMode("simulate")

	stub := ResponseDelayListStub{}
	dbClient.ResponseDelays = &stub

	newResp := dbClient.processRequest(r)

	Expect(newResp.StatusCode).To(Equal(http.StatusBadGateway))

	Expect(stub.gotDelays).To(Equal(0))
}

func TestDelayNotAppliedToCaptureRequest(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	r, err := http.NewRequest("GET", "http://somehost.com", nil)
	Expect(err).To(BeNil())

	dbClient.Cfg.SetMode("capture")

	stub := ResponseDelayListStub{}
	dbClient.ResponseDelays = &stub

	resp := dbClient.processRequest(r)

	Expect(resp.StatusCode).To(Equal(http.StatusCreated))

	Expect(stub.gotDelays).To(Equal(0))
}

func TestDelayAppliedToSynthesizeRequest(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	err := dbClient.Cfg.Middleware.SetBinary("python")
	Expect(err).To(BeNil())

	err = dbClient.Cfg.Middleware.SetScript(pythonReflectBody)
	Expect(err).To(BeNil())

	bodyBytes := []byte("request_body_here")

	r, err := http.NewRequest("GET", "http://somehost.com", ioutil.NopCloser(bytes.NewBuffer(bodyBytes)))
	Expect(err).To(BeNil())

	dbClient.Cfg.SetMode("synthesize")

	stub := ResponseDelayListStub{}
	dbClient.ResponseDelays = &stub
	newResp := dbClient.processRequest(r)

	Expect(newResp.StatusCode).To(Equal(http.StatusCreated))

	Expect(stub.gotDelays).To(Equal(1))
}

func TestDelayNotAppliedToFailedSynthesizeRequest(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	err := dbClient.Cfg.Middleware.SetBinary("python")
	Expect(err).To(BeNil())

	err = dbClient.Cfg.Middleware.SetScript(pythonMiddlewareBad)
	Expect(err).To(BeNil())

	bodyBytes := []byte("request_body_here")

	r, err := http.NewRequest("GET", "http://somehost.com", ioutil.NopCloser(bytes.NewBuffer(bodyBytes)))
	Expect(err).To(BeNil())

	dbClient.Cfg.SetMode("synthesize")

	stub := ResponseDelayListStub{}
	dbClient.ResponseDelays = &stub
	newResp := dbClient.processRequest(r)

	Expect(newResp.StatusCode).To(Equal(http.StatusBadGateway))

	Expect(stub.gotDelays).To(Equal(0))
}

func TestDelayAppliedToSuccessfulMiddleware(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	err := dbClient.Cfg.Middleware.SetBinary("python")
	Expect(err).To(BeNil())

	err = dbClient.Cfg.Middleware.SetScript(pythonModifyResponse)
	Expect(err).To(BeNil())

	r, err := http.NewRequest("POST", "http://somehost.com", nil)
	Expect(err).To(BeNil())

	dbClient.Cfg.SetMode("modify")

	stub := ResponseDelayListStub{}
	dbClient.ResponseDelays = &stub
	newResp := dbClient.processRequest(r)

	Expect(newResp.StatusCode).To(Equal(http.StatusCreated))

	Expect(stub.gotDelays).To(Equal(1))
}

func TestDelayNotAppliedToFailedModifyRequest(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	err := unit.Cfg.Middleware.SetBinary("python")
	Expect(err).To(BeNil())

	err = unit.Cfg.Middleware.SetScript(pythonMiddlewareBad)
	Expect(err).To(BeNil())

	r, err := http.NewRequest("POST", "http://somehost.com", nil)
	Expect(err).To(BeNil())

	unit.Cfg.SetMode("modify")

	stub := ResponseDelayListStub{}
	unit.ResponseDelays = &stub
	newResp := unit.processRequest(r)

	Expect(newResp.StatusCode).To(Equal(http.StatusBadGateway))

	Expect(stub.gotDelays).To(Equal(0))
}

func Test_Hoverfly_DoRequest_DoesNotPanicWhenCannotMakeRequest(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	ioutil.NopCloser(bytes.NewBuffer([]byte("")))
	request, err := http.NewRequest("GET", "w.specto.fake", ioutil.NopCloser(bytes.NewBuffer([]byte(""))))
	Expect(err).To(BeNil())

	response, err := unit.DoRequest(request)
	Expect(response).To(BeNil())
	Expect(err).ToNot(BeNil())
}

func Test_Hoverfly_DoRequest_FailedHTTP(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	requestBody := []byte("fizz=buzz")

	body := ioutil.NopCloser(bytes.NewBuffer(requestBody))

	req, err := http.NewRequest("POST", "http://capture_body.com", body)
	Expect(err).To(BeNil())

	_, err = unit.DoRequest(req)
	Expect(err).ToNot(BeNil())
}

// TestCaptureHeader tests whether request gets new header assigned
func Test_DoRequest_AddsHoverflyHeaderOnSuccessfulRequest(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	req, err := http.NewRequest("GET", "http://example.com", ioutil.NopCloser(bytes.NewBuffer([]byte(""))))
	Expect(err).To(BeNil())

	response, err := unit.DoRequest(req)

	Expect(response.Header.Get("hoverfly")).To(Equal("Was-Here"))
}

func Test_Hoverfly_Save_SavesRequestAndResponseToSimulation(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	unit.Save(&models.RequestDetails{
		Body:        "testbody",
		Destination: "testdestination",
		Headers:     map[string][]string{"testheader": []string{"testvalue"}},
		Method:      "testmethod",
		Path:        "/testpath",
		Query:       "?query=test",
		Scheme:      "http",
	}, &models.ResponseDetails{
		Body:    "testresponsebody",
		Headers: map[string][]string{"testheader": []string{"testvalue"}},
		Status:  200,
	})

	Expect(unit.Simulation.Templates).To(HaveLen(1))

	Expect(*unit.Simulation.Templates[0].RequestTemplate.Body).To(Equal("testbody"))
	Expect(*unit.Simulation.Templates[0].RequestTemplate.Destination).To(Equal("testdestination"))
	Expect(unit.Simulation.Templates[0].RequestTemplate.Headers).To(HaveKeyWithValue("testheader", []string{"testvalue"}))
	Expect(*unit.Simulation.Templates[0].RequestTemplate.Method).To(Equal("testmethod"))
	Expect(*unit.Simulation.Templates[0].RequestTemplate.Path).To(Equal("/testpath"))
	Expect(*unit.Simulation.Templates[0].RequestTemplate.Query).To(Equal("?query=test"))
	Expect(*unit.Simulation.Templates[0].RequestTemplate.Scheme).To(Equal("http"))

	Expect(unit.Simulation.Templates[0].Response.Body).To(Equal("testresponsebody"))
	Expect(unit.Simulation.Templates[0].Response.Headers).To(HaveKeyWithValue("testheader", []string{"testvalue"}))
	Expect(unit.Simulation.Templates[0].Response.Status).To(Equal(200))
}

func Test_Hoverfly_Save_SavesIncompleteRequestAndResponseToSimulation(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	unit.Save(&models.RequestDetails{
		Destination: "testdestination",
	}, &models.ResponseDetails{
		Body:    "testresponsebody",
		Headers: map[string][]string{"testheader": []string{"testvalue"}},
		Status:  200,
	})

	Expect(unit.Simulation.Templates).To(HaveLen(1))

	// Expect(unit.Simulation.Templates[0].RequestTemplate.Body).To(BeNil())
	Expect(*unit.Simulation.Templates[0].RequestTemplate.Destination).To(Equal("testdestination"))
	// Expect(unit.Simulation.Templates[0].RequestTemplate.Headers).To(BeNil())
	// Expect(*unit.Simulation.Templates[0].RequestTemplate.Method).To(BeNil())
	// Expect(*unit.Simulation.Templates[0].RequestTemplate.Path).To(BeNil())
	// Expect(*unit.Simulation.Templates[0].RequestTemplate.Query).To(BeNil())
	// Expect(*unit.Simulation.Templates[0].RequestTemplate.Scheme).To(BeNil())

	Expect(unit.Simulation.Templates[0].Response.Body).To(Equal("testresponsebody"))
	Expect(unit.Simulation.Templates[0].Response.Headers).To(HaveKeyWithValue("testheader", []string{"testvalue"}))
	Expect(unit.Simulation.Templates[0].Response.Status).To(Equal(200))
}

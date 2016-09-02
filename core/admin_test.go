package hoverfly

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/SpectoLabs/hoverfly/core/views"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAllRecords(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()
	m := getBoneRouter(dbClient)

	req, err := http.NewRequest("GET", "/api/records", nil)
	Expect(err).To(BeNil())

	//The response recorder used to record HTTP responses
	respRec := httptest.NewRecorder()

	m.ServeHTTP(respRec, req)

	Expect(respRec.Code, http.StatusOK)

	body, err := ioutil.ReadAll(respRec.Body)

	pair := views.RequestResponsePairPayload{}
	err = json.Unmarshal(body, &pair)

	Expect(len(pair.Data)).To(Equal(0))
}

func TestGetAllRecordsWRecords(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()

	// inserting some payloads
	for i := 0; i < 5; i++ {
		req, err := http.NewRequest("GET", fmt.Sprintf("http://example.com/q=%d", i), nil)
		Expect(err).To(BeNil())
		dbClient.captureRequest(req)
	}
	// performing query
	m := getBoneRouter(dbClient)

	req, err := http.NewRequest("GET", "/api/records", nil)
	Expect(err).To(BeNil())

	//The response recorder used to record HTTP responses
	respRec := httptest.NewRecorder()

	m.ServeHTTP(respRec, req)

	Expect(respRec.Code).To(Equal(http.StatusOK))

	body, err := ioutil.ReadAll(respRec.Body)

	pair := views.RequestResponsePairPayload{}
	err = json.Unmarshal(body, &pair)

	Expect(len(pair.Data)).To(Equal(5))
}

func TestGetRecordsCount(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()
	m := getBoneRouter(dbClient)

	req, err := http.NewRequest("GET", "/api/count", nil)
	Expect(err).To(BeNil())

	//The response recorder used to record HTTP responses
	respRec := httptest.NewRecorder()

	m.ServeHTTP(respRec, req)

	Expect(respRec.Code, http.StatusOK)

	body, err := ioutil.ReadAll(respRec.Body)

	rc := recordsCount{}
	err = json.Unmarshal(body, &rc)

	Expect(rc.Count).To(Equal(0))
}

func TestGetRecordsCountWRecords(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()

	// inserting some payloads
	for i := 0; i < 5; i++ {
		req, err := http.NewRequest("GET", fmt.Sprintf("http://example.com/q=%d", i), nil)
		Expect(err).To(BeNil())
		dbClient.captureRequest(req)
	}
	// performing query
	m := getBoneRouter(dbClient)

	req, err := http.NewRequest("GET", "/api/count", nil)
	Expect(err).To(BeNil())

	//The response recorder used to record HTTP responses
	respRec := httptest.NewRecorder()

	m.ServeHTTP(respRec, req)

	Expect(respRec.Code, http.StatusOK)

	body, err := ioutil.ReadAll(respRec.Body)

	rc := recordsCount{}
	err = json.Unmarshal(body, &rc)

	Expect(rc.Count).To(Equal(5))
}

func TestExportImportRecords(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()
	m := getBoneRouter(dbClient)

	// inserting some payloads
	for i := 0; i < 5; i++ {
		req, err := http.NewRequest("GET", fmt.Sprintf("http://example.com/q=%d", i), nil)
		Expect(err).To(BeNil())
		dbClient.captureRequest(req)
	}

	req, err := http.NewRequest("GET", "/api/records", nil)
	Expect(err).To(BeNil())

	//The response recorder used to record HTTP responses
	respRec := httptest.NewRecorder()

	m.ServeHTTP(respRec, req)

	Expect(respRec.Code, http.StatusOK)

	body, err := ioutil.ReadAll(respRec.Body)

	// deleting records
	err = dbClient.RequestCache.DeleteData()
	Expect(err).To(BeNil())

	// using body to import records again
	importReq, err := http.NewRequest("POST", "/api/records", ioutil.NopCloser(bytes.NewBuffer(body)))
	//The response recorder used to record HTTP responses
	importRec := httptest.NewRecorder()

	m.ServeHTTP(importRec, importReq)
	Expect(respRec.Code, http.StatusOK)

	// records should be there
	pairBytes, err := dbClient.RequestCache.GetAllValues()
	Expect(err).To(BeNil())
	Expect(len(pairBytes)).To(Equal(5))
}

func TestDeleteHandler(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()
	m := getBoneRouter(dbClient)

	// inserting some payloads
	for i := 0; i < 5; i++ {
		req, err := http.NewRequest("GET", fmt.Sprintf("http://example.com/q=%d", i), nil)
		Expect(err).To(BeNil())
		dbClient.captureRequest(req)
	}

	// checking whether we have records
	pairBytes, err := dbClient.RequestCache.GetAllValues()
	Expect(err).To(BeNil())
	Expect(len(pairBytes)).To(Equal(5))

	// deleting through handler
	deleteReq, err := http.NewRequest("DELETE", "/api/records", nil)
	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, deleteReq)
	Expect(rec.Code, http.StatusOK)
}

func TestDeleteHandlerNoBucket(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()
	m := getBoneRouter(dbClient)

	// deleting through handler
	importReq, err := http.NewRequest("DELETE", "/api/records", nil)
	Expect(err).To(BeNil())
	//The response recorder used to record HTTP responses
	importRec := httptest.NewRecorder()

	m.ServeHTTP(importRec, importReq)
	Expect(importRec.Code, http.StatusOK)
}

func TestGetState(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()
	m := getBoneRouter(dbClient)

	// setting initial mode
	dbClient.Cfg.SetMode(SimulateMode)

	req, err := http.NewRequest("GET", "/api/state", nil)
	Expect(err).To(BeNil())
	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	Expect(rec.Code, http.StatusOK)

	body, err := ioutil.ReadAll(rec.Body)

	sr := stateRequest{}
	err = json.Unmarshal(body, &sr)

	Expect(sr.Mode).To(Equal(SimulateMode))
}

func TestSetSimulateState(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()
	m := getBoneRouter(dbClient)

	// setting mode to capture
	dbClient.Cfg.SetMode("capture")

	// preparing to set mode through rest api
	var resp stateRequest
	resp.Mode = SimulateMode

	requestBytes, err := json.Marshal(&resp)
	Expect(err).To(BeNil())

	// deleting through handler
	req, err := http.NewRequest("POST", "/api/state", ioutil.NopCloser(bytes.NewBuffer(requestBytes)))
	Expect(err).To(BeNil())

	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	Expect(rec.Code).To(Equal(http.StatusOK))

	// checking mode
	Expect(dbClient.Cfg.GetMode()).To(Equal(SimulateMode))
}

func TestSetCaptureState(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()
	m := getBoneRouter(dbClient)

	// setting mode to simulate
	dbClient.Cfg.SetMode(SimulateMode)

	// preparing to set mode through rest api
	var resp stateRequest
	resp.Mode = "capture"

	requestBytes, err := json.Marshal(&resp)
	Expect(err).To(BeNil())

	// deleting through handler
	req, err := http.NewRequest("POST", "/api/state", ioutil.NopCloser(bytes.NewBuffer(requestBytes)))
	Expect(err).To(BeNil())

	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	Expect(rec.Code).To(Equal(http.StatusOK))

	// checking mode
	Expect(dbClient.Cfg.GetMode()).To(Equal("capture"))
}

func TestSetModifyState(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()
	m := getBoneRouter(dbClient)

	// setting mode to simulate
	dbClient.Cfg.SetMode(SimulateMode)

	// preparing to set mode through rest api
	var resp stateRequest
	resp.Mode = ModifyMode

	requestBytes, err := json.Marshal(&resp)
	Expect(err).To(BeNil())

	// deleting through handler
	req, err := http.NewRequest("POST", "/api/state", ioutil.NopCloser(bytes.NewBuffer(requestBytes)))
	Expect(err).To(BeNil())

	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	Expect(rec.Code).To(Equal(http.StatusOK))

	// checking mode
	Expect(dbClient.Cfg.GetMode()).To(Equal(ModifyMode))
}

func TestSetSynthesizeState(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()
	m := getBoneRouter(dbClient)

	// setting mode to simulate
	dbClient.Cfg.SetMode(SimulateMode)

	// preparing to set mode through rest api
	var resp stateRequest
	resp.Mode = SynthesizeMode

	requestBytes, err := json.Marshal(&resp)
	Expect(err).To(BeNil())

	// deleting through handler
	req, err := http.NewRequest("POST", "/api/state", ioutil.NopCloser(bytes.NewBuffer(requestBytes)))
	Expect(err).To(BeNil())

	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	Expect(rec.Code).To(Equal(http.StatusOK))

	// checking mode
	Expect(dbClient.Cfg.GetMode()).To(Equal(SynthesizeMode))
}

func TestSetRandomState(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()
	m := getBoneRouter(dbClient)

	// setting mode to simulate
	dbClient.Cfg.SetMode(SimulateMode)

	// preparing to set mode through rest api
	var resp stateRequest
	resp.Mode = "shouldnotwork"

	requestBytes, err := json.Marshal(&resp)
	Expect(err).To(BeNil())

	// deleting through handler
	req, err := http.NewRequest("POST", "/api/state", ioutil.NopCloser(bytes.NewBuffer(requestBytes)))
	Expect(err).To(BeNil())

	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	Expect(rec.Code).To(Equal(http.StatusBadRequest))

	// checking mode, should not have changed
	Expect(dbClient.Cfg.GetMode()).To(Equal(SimulateMode))
}

func TestSetNoBody(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()
	m := getBoneRouter(dbClient)

	// setting mode to simulate
	dbClient.Cfg.SetMode(SimulateMode)

	// setting state
	req, err := http.NewRequest("POST", "/api/state", nil)
	Expect(err).To(BeNil())

	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	Expect(rec.Code).To(Equal(http.StatusBadRequest))

	// checking mode, should not have changed
	Expect(dbClient.Cfg.GetMode()).To(Equal(SimulateMode))
}

func TestGetMiddleware(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()
	m := getBoneRouter(dbClient)

	dbClient.Cfg.Middleware = "python middleware_test.py"
	req, err := http.NewRequest("GET", "/api/middleware", nil)
	Expect(err).To(BeNil())

	rec := httptest.NewRecorder()
	m.ServeHTTP(rec, req)
	Expect(rec.Code).To(Equal(http.StatusOK))

	body, err := ioutil.ReadAll(rec.Body)
	Expect(err).To(BeNil())

	middlewareResponse := middlewareSchema{}
	err = json.Unmarshal(body, &middlewareResponse)
	Expect(err).To(BeNil())

	Expect(middlewareResponse.Middleware).To(Equal("python middleware_test.py"))
}

func TestSetMiddleware_WithValidMiddleware(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()
	m := getBoneRouter(dbClient)

	dbClient.Cfg.Middleware = "python examples/middleware/modify_request/modify_request.py"

	var middlewareReq middlewareSchema
	middlewareReq.Middleware = "python examples/middleware/delay_policy/add_random_delay.py"

	middlewareReqBytes, err := json.Marshal(&middlewareReq)
	Expect(err).To(BeNil())

	req, err := http.NewRequest("POST", "/api/middleware", ioutil.NopCloser(bytes.NewBuffer(middlewareReqBytes)))
	Expect(err).To(BeNil())

	rec := httptest.NewRecorder()
	m.ServeHTTP(rec, req)
	Expect(rec.Code).To(Equal(http.StatusOK))

	body, err := ioutil.ReadAll(rec.Body)
	Expect(err).To(BeNil())

	middlewareResp := middlewareSchema{}
	err = json.Unmarshal(body, &middlewareResp)
	Expect(err).To(BeNil())

	Expect(middlewareResp.Middleware).To(Equal("python examples/middleware/delay_policy/add_random_delay.py"))
	Expect(dbClient.Cfg.Middleware).To(Equal("python examples/middleware/delay_policy/add_random_delay.py"))
}

func TestSetMiddleware_WithInvalidMiddleware(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()
	m := getBoneRouter(dbClient)

	dbClient.Cfg.Middleware = "python examples/middleware/modify_request/modify_request.py"

	var middlewareReq middlewareSchema
	middlewareReq.Middleware = "definitely won't execute"

	middlewareReqBytes, err := json.Marshal(&middlewareReq)
	Expect(err).To(BeNil())

	req, err := http.NewRequest("POST", "/api/middleware", ioutil.NopCloser(bytes.NewBuffer(middlewareReqBytes)))
	Expect(err).To(BeNil())

	rec := httptest.NewRecorder()
	m.ServeHTTP(rec, req)
	Expect(rec.Code).To(Equal(http.StatusBadRequest))

	body, err := ioutil.ReadAll(rec.Body)
	Expect(err).To(BeNil())

	middlewareResp := middlewareSchema{}
	err = json.Unmarshal(body, &middlewareResp)

	Expect(err).ToNot(BeNil())
	Expect(string(body)).To(ContainSubstring("Invalid middleware"))

	Expect(dbClient.Cfg.Middleware).To(Equal("python examples/middleware/modify_request/modify_request.py"))
}

func TestSetMiddleware_WithEmptyMiddleware(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()
	m := getBoneRouter(dbClient)

	dbClient.Cfg.Middleware = "python examples/middleware/modify_request/modify_request.py"

	var middlewareReq middlewareSchema
	middlewareReq.Middleware = ""

	middlewareReqBytes, err := json.Marshal(&middlewareReq)
	Expect(err).To(BeNil())

	req, err := http.NewRequest("POST", "/api/middleware", ioutil.NopCloser(bytes.NewBuffer(middlewareReqBytes)))
	Expect(err).To(BeNil())

	rec := httptest.NewRecorder()
	m.ServeHTTP(rec, req)
	Expect(rec.Code).To(Equal(http.StatusOK))

	body, err := ioutil.ReadAll(rec.Body)
	Expect(err).To(BeNil())

	middlewareResp := middlewareSchema{}
	err = json.Unmarshal(body, &middlewareResp)

	Expect(err).To(BeNil())

	Expect(middlewareResp.Middleware).To(Equal(""))
	Expect(dbClient.Cfg.Middleware).To(Equal(""))
}

func TestStatsHandler(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()
	m := getBoneRouter(dbClient)

	// deleting through handler
	req, err := http.NewRequest("GET", "/api/stats", nil)
	Expect(err).To(BeNil())

	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	Expect(rec.Code).To(Equal(http.StatusOK))
}

func TestStatsHandlerSimulateMetrics(t *testing.T) {
	RegisterTestingT(t)

	// test metrics, increases simulate count by 1 and then checks through stats
	// handler whether it is visible through /stats handler
	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()
	m := getBoneRouter(dbClient)

	dbClient.Counter.Counters[SimulateMode].Inc(1)

	req, err := http.NewRequest("GET", "/api/stats", nil)
	Expect(err).To(BeNil())

	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	Expect(rec.Code).To(Equal(http.StatusOK))

	body, err := ioutil.ReadAll(rec.Body)

	sr := statsResponse{}
	err = json.Unmarshal(body, &sr)

	Expect(int(sr.Stats.Counters[SimulateMode])).To(Equal(1))
}

func TestStatsHandlerCaptureMetrics(t *testing.T) {
	RegisterTestingT(t)

	// test metrics, increases capture count by 1 and then checks through stats
	// handler whether it is visible through /stats handler
	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()
	m := getBoneRouter(dbClient)

	dbClient.Counter.Counters[CaptureMode].Inc(1)

	req, err := http.NewRequest("GET", "/api/stats", nil)
	Expect(err).To(BeNil())

	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	Expect(rec.Code).To(Equal(http.StatusOK))

	body, err := ioutil.ReadAll(rec.Body)

	sr := statsResponse{}
	err = json.Unmarshal(body, &sr)

	Expect(int(sr.Stats.Counters[CaptureMode])).To(Equal(1))
}

func TestStatsHandlerModifyMetrics(t *testing.T) {
	RegisterTestingT(t)

	// test metrics, increases modify count by 1 and then checks through stats
	// handler whether it is visible through /stats handler
	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()
	m := getBoneRouter(dbClient)

	dbClient.Counter.Counters[ModifyMode].Inc(1)

	req, err := http.NewRequest("GET", "/api/stats", nil)
	Expect(err).To(BeNil())

	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	Expect(rec.Code).To(Equal(http.StatusOK))

	body, err := ioutil.ReadAll(rec.Body)

	sr := statsResponse{}
	err = json.Unmarshal(body, &sr)

	Expect(int(sr.Stats.Counters[ModifyMode])).To(Equal(1))
}

func TestStatsHandlerSynthesizeMetrics(t *testing.T) {
	RegisterTestingT(t)

	// test metrics, increases synthesize count by 1 and then checks through stats
	// handler whether it is visible through /stats handler
	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()
	m := getBoneRouter(dbClient)

	dbClient.Counter.Counters[SynthesizeMode].Inc(1)

	req, err := http.NewRequest("GET", "/api/stats", nil)
	Expect(err).To(BeNil())

	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	Expect(rec.Code).To(Equal(http.StatusOK))

	body, err := ioutil.ReadAll(rec.Body)

	sr := statsResponse{}
	err = json.Unmarshal(body, &sr)

	Expect(int(sr.Stats.Counters[SynthesizeMode])).To(Equal(1))
}

func TestStatsHandlerRecordCountMetrics(t *testing.T) {
	RegisterTestingT(t)

	// test metrics, adds 5 new requests and then checks through stats
	// handler whether it is visible through /stats handler
	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()
	m := getBoneRouter(dbClient)

	// inserting some payloads
	for i := 0; i < 5; i++ {
		req, err := http.NewRequest("GET", fmt.Sprintf("http://example.com/q=%d", i), nil)
		Expect(err).To(BeNil())
		dbClient.captureRequest(req)
	}

	req, err := http.NewRequest("GET", "/api/stats", nil)
	Expect(err).To(BeNil())

	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	Expect(rec.Code).To(Equal(http.StatusOK))

	body, err := ioutil.ReadAll(rec.Body)

	sr := statsResponse{}
	err = json.Unmarshal(body, &sr)

	Expect(int(sr.RecordsCount)).To(Equal(5))
}

func TestSetMetadata(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()
	m := getBoneRouter(dbClient)

	// preparing to set mode through rest api
	var reqBody setMetadata
	reqBody.Key = "some_key"
	reqBody.Value = "some_val"

	reqBodyBytes, err := json.Marshal(&reqBody)
	Expect(err).To(BeNil())

	// deleting through handler
	req, err := http.NewRequest("PUT", "/api/metadata", ioutil.NopCloser(bytes.NewBuffer(reqBodyBytes)))
	Expect(err).To(BeNil())

	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	Expect(rec.Code).To(Equal(http.StatusCreated))

	// checking mode
	metaValue, err := dbClient.MetadataCache.Get([]byte("some_key"))
	Expect(err).To(BeNil())
	Expect(string(metaValue)).To(Equal("some_val"))
}

func TestSetMetadataBadBody(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()
	m := getBoneRouter(dbClient)

	// deleting through handler
	req, err := http.NewRequest("PUT", "/api/metadata", ioutil.NopCloser(bytes.NewBuffer([]byte("you shall not decode me!!"))))
	Expect(err).To(BeNil())

	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	Expect(rec.Code).To(Equal(http.StatusBadRequest))
}

func TestSetMetadataMissingKey(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()
	m := getBoneRouter(dbClient)

	// preparing to set mode through rest api
	var reqBody setMetadata
	// missing key
	reqBody.Value = "some_val"

	reqBodyBytes, err := json.Marshal(&reqBody)
	Expect(err).To(BeNil())

	// deleting through handler
	req, err := http.NewRequest("PUT", "/api/metadata", ioutil.NopCloser(bytes.NewBuffer(reqBodyBytes)))
	Expect(err).To(BeNil())

	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	Expect(rec.Code).To(Equal(http.StatusBadRequest))

	// checking response body
	body, err := ioutil.ReadAll(rec.Body)
	mr := messageResponse{}
	err = json.Unmarshal(body, &mr)

	Expect(mr.Message).To(Equal("Key not provided."))
}

func TestGetMetadata(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()
	m := getBoneRouter(dbClient)
	// adding some metadata
	for i := 0; i < 3; i++ {
		k := fmt.Sprintf("key_%d", i)
		v := fmt.Sprintf("val_%d", i)
		err := dbClient.MetadataCache.Set([]byte(k), []byte(v))
		Expect(err).To(BeNil())
	}

	req, err := http.NewRequest("GET", "/api/metadata", nil)
	Expect(err).To(BeNil())

	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	Expect(rec.Code).To(Equal(http.StatusOK))

	body, err := ioutil.ReadAll(rec.Body)

	sm := storedMetadata{}
	err = json.Unmarshal(body, &sm)

	Expect(len(sm.Data)).To(Equal(3))

	for i := 0; i < 3; i++ {
		k := fmt.Sprintf("key_%d", i)
		v := fmt.Sprintf("val_%d", i)
		Expect(sm.Data[k]).To(Equal(v))
	}
}

func TestDeleteMetadata(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()
	m := getBoneRouter(dbClient)
	// adding some metadata
	for i := 0; i < 3; i++ {
		k := fmt.Sprintf("key_%d", i)
		v := fmt.Sprintf("val_%d", i)
		err := dbClient.MetadataCache.Set([]byte(k), []byte(v))
		Expect(err).To(BeNil())
	}

	// checking that metadata is there
	allMeta, err := dbClient.MetadataCache.GetAllEntries()
	Expect(err).To(BeNil())
	Expect(len(allMeta)).To(Equal(3))

	// deleting it
	req, err := http.NewRequest("DELETE", "/api/metadata", nil)
	Expect(err).To(BeNil())
	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	Expect(rec.Code).To(Equal(http.StatusOK))

	// checking metadata again, should be zero
	allMeta, err = dbClient.MetadataCache.GetAllEntries()
	Expect(err).To(BeNil())
	Expect(len(allMeta)).To(Equal(0))
}

func TestDeleteMetadataEmpty(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()
	m := getBoneRouter(dbClient)

	// deleting it
	req, err := http.NewRequest("DELETE", "/api/metadata", nil)
	Expect(err).To(BeNil())
	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	Expect(rec.Code).To(Equal(http.StatusOK))

	// checking metadata again, should be zero
	allMeta, err := dbClient.MetadataCache.GetAllEntries()
	Expect(err).To(BeNil())
	Expect(len(allMeta)).To(Equal(0))
}

func TestGetResponseDelays(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()

	delay := models.ResponseDelay{
		UrlPattern: ".",
		Delay:      100,
	}
	delays := models.ResponseDelayList{delay}
	dbClient.UpdateResponseDelays(delays)

	m := getBoneRouter(dbClient)

	req, err := http.NewRequest("GET", "/api/delays", nil)
	Expect(err).To(BeNil())
	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	Expect(rec.Code).To(Equal(http.StatusOK))

	body, err := ioutil.ReadAll(rec.Body)

	sr := models.ResponseDelayPayload{}
	err = json.Unmarshal(body, &sr)

	// normal equality checking doesn't work on slices (!!)
	Expect(*sr.Data).To(Equal(delays))
}

func TestDeleteAllResponseDelaysHandler(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()
	delay := models.ResponseDelay{
		UrlPattern: ".",
		Delay:      100,
	}
	delays := models.ResponseDelayList{delay}
	dbClient.ResponseDelays = &delays
	m := getBoneRouter(dbClient)

	req, err := http.NewRequest("DELETE", "/api/delays", nil)
	Expect(err).To(BeNil())

	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	Expect(rec.Code).To(Equal(http.StatusOK))

	Expect(dbClient.ResponseDelays.Len()).To(Equal(0))
}

func TestUpdateResponseDelays(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()
	m := getBoneRouter(dbClient)

	delayOne := models.ResponseDelay{
		UrlPattern: ".",
		Delay:      100,
	}
	delayTwo := models.ResponseDelay{
		UrlPattern: "example",
		Delay:      100,
	}
	delays := models.ResponseDelayList{delayOne, delayTwo}
	delayJson := models.ResponseDelayPayload{Data: &delays}
	delayJsonBytes, err := json.Marshal(&delayJson)
	Expect(err).To(BeNil())

	req, err := http.NewRequest("PUT", "/api/delays", ioutil.NopCloser(bytes.NewBuffer(delayJsonBytes)))
	Expect(err).To(BeNil())

	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	Expect(rec.Code).To(Equal(http.StatusCreated))

	// normal equality checking doesn't work on slices (!!)
	Expect(dbClient.ResponseDelays).To(Equal(&delays))
}

func TestInvalidJSONSyntaxUpdateResponseDelays(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()
	m := getBoneRouter(dbClient)

	delayJson := "{aseuifhksejfc}"

	req, err := http.NewRequest("PUT", "/api/delays", ioutil.NopCloser(bytes.NewBuffer([]byte(delayJson))))
	Expect(err).To(BeNil())

	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	Expect(rec.Code).To(Equal(http.StatusBadRequest))

	// normal equality checking doesn't work on slices (!!)
	Expect(dbClient.ResponseDelays).To(Equal(&models.ResponseDelayList{}))
}

func TestInvalidJSONSemanticsUpdateResponseDelays(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()
	m := getBoneRouter(dbClient)

	delayJson := "{ \"madeupfield\" : \"somevalue\" }"

	req, err := http.NewRequest("PUT", "/api/delays", ioutil.NopCloser(bytes.NewBuffer([]byte(delayJson))))
	Expect(err).To(BeNil())

	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	Expect(rec.Code).To(Equal(422))

	// normal equality checking doesn't work on slices (!!)
	Expect(dbClient.ResponseDelays).To(Equal(&models.ResponseDelayList{}))
}

func TestJSONWithInvalidHostPatternUpdateResponseDelays(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()
	m := getBoneRouter(dbClient)

	delayJson := "{ \"data\": [{\"hostPattern\": \"*\", \"delay\": 100}] }"

	req, err := http.NewRequest("PUT", "/api/delays", ioutil.NopCloser(bytes.NewBuffer([]byte(delayJson))))
	Expect(err).To(BeNil())

	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	Expect(rec.Code).To(Equal(422))

	// normal equality checking doesn't work on slices (!!)
	Expect(dbClient.ResponseDelays).To(Equal(&models.ResponseDelayList{}))
}

func TestJSONWithMissingFieldUpdateResponseDelays(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()
	m := getBoneRouter(dbClient)

	delayJson := "{ \"data\" : [{\"hostPattern\": \".\"}] }"

	req, err := http.NewRequest("PUT", "/api/delays", ioutil.NopCloser(bytes.NewBuffer([]byte(delayJson))))
	Expect(err).To(BeNil())

	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	Expect(rec.Code).To(Equal(422))

	// normal equality checking doesn't work on slices (!!)
	Expect(dbClient.ResponseDelays).To(Equal(&models.ResponseDelayList{}))
}

func TestGetMissingURL(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()
	m := getBoneRouter(dbClient)

	req, err := http.NewRequest("GET", "/api/sdiughvksjv", nil)
	Expect(err).To(BeNil())

	//The response recorder used to record HTTP responses
	respRec := httptest.NewRecorder()

	m.ServeHTTP(respRec, req)

	Expect(respRec.Code, http.StatusNotFound)
}

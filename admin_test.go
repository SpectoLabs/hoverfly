package hoverfly

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAllRecords(t *testing.T) {
	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.Cache.DeleteData()
	m := getBoneRouter(*dbClient)

	req, err := http.NewRequest("GET", "/records", nil)
	expect(t, err, nil)

	//The response recorder used to record HTTP responses
	respRec := httptest.NewRecorder()

	m.ServeHTTP(respRec, req)

	expect(t, respRec.Code, http.StatusOK)

	body, err := ioutil.ReadAll(respRec.Body)

	rr := recordedRequests{}
	err = json.Unmarshal(body, &rr)

	expect(t, len(rr.Data), 0)
}

func TestGetAllRecordsWRecords(t *testing.T) {
	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.Cache.DeleteData()

	// inserting some payloads
	for i := 0; i < 5; i++ {
		req, err := http.NewRequest("GET", fmt.Sprintf("http://example.com/q=%d", i), nil)
		expect(t, err, nil)
		dbClient.captureRequest(req)
	}
	// performing query
	m := getBoneRouter(*dbClient)

	req, err := http.NewRequest("GET", "/records", nil)
	expect(t, err, nil)

	//The response recorder used to record HTTP responses
	respRec := httptest.NewRecorder()

	m.ServeHTTP(respRec, req)

	expect(t, respRec.Code, http.StatusOK)

	body, err := ioutil.ReadAll(respRec.Body)

	rr := recordedRequests{}
	err = json.Unmarshal(body, &rr)

	expect(t, len(rr.Data), 5)
}

func TestGetRecordsCount(t *testing.T) {
	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.Cache.DeleteData()
	m := getBoneRouter(*dbClient)

	req, err := http.NewRequest("GET", "/count", nil)
	expect(t, err, nil)

	//The response recorder used to record HTTP responses
	respRec := httptest.NewRecorder()

	m.ServeHTTP(respRec, req)

	expect(t, respRec.Code, http.StatusOK)

	body, err := ioutil.ReadAll(respRec.Body)

	rc := recordsCount{}
	err = json.Unmarshal(body, &rc)

	expect(t, rc.Count, 0)
}

func TestGetRecordsCountWRecords(t *testing.T) {
	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.Cache.DeleteData()

	// inserting some payloads
	for i := 0; i < 5; i++ {
		req, err := http.NewRequest("GET", fmt.Sprintf("http://example.com/q=%d", i), nil)
		expect(t, err, nil)
		dbClient.captureRequest(req)
	}
	// performing query
	m := getBoneRouter(*dbClient)

	req, err := http.NewRequest("GET", "/count", nil)
	expect(t, err, nil)

	//The response recorder used to record HTTP responses
	respRec := httptest.NewRecorder()

	m.ServeHTTP(respRec, req)

	expect(t, respRec.Code, http.StatusOK)

	body, err := ioutil.ReadAll(respRec.Body)

	rc := recordsCount{}
	err = json.Unmarshal(body, &rc)

	expect(t, rc.Count, 5)
}

func TestExportImportRecords(t *testing.T) {
	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.Cache.DeleteData()
	m := getBoneRouter(*dbClient)

	// inserting some payloads
	for i := 0; i < 5; i++ {
		req, err := http.NewRequest("GET", fmt.Sprintf("http://example.com/q=%d", i), nil)
		expect(t, err, nil)
		dbClient.captureRequest(req)
	}

	req, err := http.NewRequest("GET", "/records", nil)
	expect(t, err, nil)

	//The response recorder used to record HTTP responses
	respRec := httptest.NewRecorder()

	m.ServeHTTP(respRec, req)

	expect(t, respRec.Code, http.StatusOK)

	body, err := ioutil.ReadAll(respRec.Body)

	// deleting records
	err = dbClient.Cache.DeleteData()
	expect(t, err, nil)

	// using body to import records again
	importReq, err := http.NewRequest("POST", "/records", ioutil.NopCloser(bytes.NewBuffer(body)))
	//The response recorder used to record HTTP responses
	importRec := httptest.NewRecorder()

	m.ServeHTTP(importRec, importReq)
	expect(t, importRec.Code, http.StatusOK)

	// records should be there
	payloads, err := dbClient.Cache.GetAllRequests()
	expect(t, err, nil)
	expect(t, len(payloads), 5)

}

func TestDeleteHandler(t *testing.T) {
	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.Cache.DeleteData()
	m := getBoneRouter(*dbClient)

	// inserting some payloads
	for i := 0; i < 5; i++ {
		req, err := http.NewRequest("GET", fmt.Sprintf("http://example.com/q=%d", i), nil)
		expect(t, err, nil)
		dbClient.captureRequest(req)
	}

	// checking whether we have records
	payloads, err := dbClient.Cache.GetAllRequests()
	expect(t, err, nil)
	expect(t, len(payloads), 5)

	// deleting through handler
	deleteReq, err := http.NewRequest("DELETE", "/records", nil)
	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, deleteReq)
	expect(t, rec.Code, http.StatusOK)
}

func TestDeleteHandlerNoBucket(t *testing.T) {
	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.Cache.DeleteData()
	m := getBoneRouter(*dbClient)

	// deleting through handler
	importReq, err := http.NewRequest("DELETE", "/records", nil)
	expect(t, err, nil)
	//The response recorder used to record HTTP responses
	importRec := httptest.NewRecorder()

	m.ServeHTTP(importRec, importReq)
	expect(t, importRec.Code, http.StatusOK)
}

func TestGetState(t *testing.T) {
	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.Cache.DeleteData()
	m := getBoneRouter(*dbClient)

	// setting initial mode
	dbClient.Cfg.SetMode("virtualize")

	req, err := http.NewRequest("GET", "/state", nil)
	expect(t, err, nil)
	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	expect(t, rec.Code, http.StatusOK)

	body, err := ioutil.ReadAll(rec.Body)

	sr := stateRequest{}
	err = json.Unmarshal(body, &sr)

	expect(t, sr.Mode, "virtualize")
}

func TestSetVirtualizeState(t *testing.T) {
	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.Cache.DeleteData()
	m := getBoneRouter(*dbClient)

	// setting mode to capture
	dbClient.Cfg.SetMode("capture")

	// preparing to set mode through rest api
	var resp stateRequest
	resp.Mode = "virtualize"

	bts, err := json.Marshal(&resp)
	expect(t, err, nil)

	// deleting through handler
	req, err := http.NewRequest("POST", "/state", ioutil.NopCloser(bytes.NewBuffer(bts)))
	expect(t, err, nil)
	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	expect(t, rec.Code, http.StatusOK)

	// checking mode
	expect(t, dbClient.Cfg.GetMode(), "virtualize")
}

func TestSetCaptureState(t *testing.T) {
	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.Cache.DeleteData()
	m := getBoneRouter(*dbClient)

	// setting mode to virtualize
	dbClient.Cfg.SetMode("virtualize")

	// preparing to set mode through rest api
	var resp stateRequest
	resp.Mode = "capture"

	bts, err := json.Marshal(&resp)
	expect(t, err, nil)

	// deleting through handler
	req, err := http.NewRequest("POST", "/state", ioutil.NopCloser(bytes.NewBuffer(bts)))
	expect(t, err, nil)
	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	expect(t, rec.Code, http.StatusOK)

	// checking mode
	expect(t, dbClient.Cfg.GetMode(), "capture")
}

func TestSetModifyState(t *testing.T) {
	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.Cache.DeleteData()
	m := getBoneRouter(*dbClient)

	// setting mode to virtualize
	dbClient.Cfg.SetMode("virtualize")

	// preparing to set mode through rest api
	var resp stateRequest
	resp.Mode = "modify"

	bts, err := json.Marshal(&resp)
	expect(t, err, nil)

	// deleting through handler
	req, err := http.NewRequest("POST", "/state", ioutil.NopCloser(bytes.NewBuffer(bts)))
	expect(t, err, nil)
	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	expect(t, rec.Code, http.StatusOK)

	// checking mode
	expect(t, dbClient.Cfg.GetMode(), "modify")
}

func TestSetSynthesizeState(t *testing.T) {
	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.Cache.DeleteData()
	m := getBoneRouter(*dbClient)

	// setting mode to virtualize
	dbClient.Cfg.SetMode("virtualize")

	// preparing to set mode through rest api
	var resp stateRequest
	resp.Mode = "synthesize"

	bts, err := json.Marshal(&resp)
	expect(t, err, nil)

	// deleting through handler
	req, err := http.NewRequest("POST", "/state", ioutil.NopCloser(bytes.NewBuffer(bts)))
	expect(t, err, nil)
	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	expect(t, rec.Code, http.StatusOK)

	// checking mode
	expect(t, dbClient.Cfg.GetMode(), "synthesize")
}

func TestSetRandomState(t *testing.T) {
	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.Cache.DeleteData()
	m := getBoneRouter(*dbClient)

	// setting mode to virtualize
	dbClient.Cfg.SetMode("virtualize")

	// preparing to set mode through rest api
	var resp stateRequest
	resp.Mode = "shouldnotwork"

	bts, err := json.Marshal(&resp)
	expect(t, err, nil)

	// deleting through handler
	req, err := http.NewRequest("POST", "/state", ioutil.NopCloser(bytes.NewBuffer(bts)))
	expect(t, err, nil)
	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	expect(t, rec.Code, http.StatusBadRequest)

	// checking mode, should not have changed
	expect(t, dbClient.Cfg.GetMode(), "virtualize")
}

func TestSetNoBody(t *testing.T) {
	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.Cache.DeleteData()
	m := getBoneRouter(*dbClient)

	// setting mode to virtualize
	dbClient.Cfg.SetMode("virtualize")

	// setting state
	req, err := http.NewRequest("POST", "/state", nil)
	expect(t, err, nil)
	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	expect(t, rec.Code, http.StatusBadRequest)

	// checking mode, should not have changed
	expect(t, dbClient.Cfg.GetMode(), "virtualize")
}

func TestStatsHandler(t *testing.T) {
	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.Cache.DeleteData()
	m := getBoneRouter(*dbClient)

	// deleting through handler
	req, err := http.NewRequest("GET", "/stats", nil)

	expect(t, err, nil)
	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	expect(t, rec.Code, http.StatusOK)
}

func TestStatsHandlerVirtualizeMetrics(t *testing.T) {
	// test metrics, increases virtualize count by 1 and then checks through stats
	// handler whether it is visible through /stats handler
	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.Cache.DeleteData()
	m := getBoneRouter(*dbClient)

	dbClient.Counter.counterVirtualize.Inc(1)

	req, err := http.NewRequest("GET", "/stats", nil)

	expect(t, err, nil)
	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	expect(t, rec.Code, http.StatusOK)

	body, err := ioutil.ReadAll(rec.Body)

	sr := statsResponse{}
	err = json.Unmarshal(body, &sr)

	expect(t, int(sr.Stats.Counters[VirtualizeMode]), 1)
}

func TestStatsHandlerCaptureMetrics(t *testing.T) {
	// test metrics, increases capture count by 1 and then checks through stats
	// handler whether it is visible through /stats handler
	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.Cache.DeleteData()
	m := getBoneRouter(*dbClient)

	dbClient.Counter.counterCapture.Inc(1)

	req, err := http.NewRequest("GET", "/stats", nil)

	expect(t, err, nil)
	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	expect(t, rec.Code, http.StatusOK)

	body, err := ioutil.ReadAll(rec.Body)

	sr := statsResponse{}
	err = json.Unmarshal(body, &sr)

	expect(t, int(sr.Stats.Counters[CaptureMode]), 1)
}

func TestStatsHandlerModifyMetrics(t *testing.T) {
	// test metrics, increases modify count by 1 and then checks through stats
	// handler whether it is visible through /stats handler
	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.Cache.DeleteData()
	m := getBoneRouter(*dbClient)

	dbClient.Counter.counterModify.Inc(1)

	req, err := http.NewRequest("GET", "/stats", nil)

	expect(t, err, nil)
	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	expect(t, rec.Code, http.StatusOK)

	body, err := ioutil.ReadAll(rec.Body)

	sr := statsResponse{}
	err = json.Unmarshal(body, &sr)

	expect(t, int(sr.Stats.Counters[ModifyMode]), 1)
}

func TestStatsHandlerSynthesizeMetrics(t *testing.T) {
	// test metrics, increases synthesize count by 1 and then checks through stats
	// handler whether it is visible through /stats handler
	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.Cache.DeleteData()
	m := getBoneRouter(*dbClient)

	dbClient.Counter.counterSynthesize.Inc(1)

	req, err := http.NewRequest("GET", "/stats", nil)

	expect(t, err, nil)
	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	expect(t, rec.Code, http.StatusOK)

	body, err := ioutil.ReadAll(rec.Body)

	sr := statsResponse{}
	err = json.Unmarshal(body, &sr)

	expect(t, int(sr.Stats.Counters[SynthesizeMode]), 1)
}

func TestStatsHandlerRecordCountMetrics(t *testing.T) {
	// test metrics, adds 5 new requests and then checks through stats
	// handler whether it is visible through /stats handler
	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.Cache.DeleteData()
	m := getBoneRouter(*dbClient)

	// inserting some payloads
	for i := 0; i < 5; i++ {
		req, err := http.NewRequest("GET", fmt.Sprintf("http://example.com/q=%d", i), nil)
		expect(t, err, nil)
		dbClient.captureRequest(req)
	}

	req, err := http.NewRequest("GET", "/stats", nil)

	expect(t, err, nil)
	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	expect(t, rec.Code, http.StatusOK)

	body, err := ioutil.ReadAll(rec.Body)

	sr := statsResponse{}
	err = json.Unmarshal(body, &sr)

	expect(t, int(sr.RecordsCount), 5)
}

func TestSetMetadata(t *testing.T) {
	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.Cache.DeleteData()
	m := getBoneRouter(*dbClient)

	// preparing to set mode through rest api
	var reqBody setMetadata
	reqBody.Key = "some_key"
	reqBody.Value = "some_val"

	bts, err := json.Marshal(&reqBody)
	expect(t, err, nil)

	// deleting through handler
	req, err := http.NewRequest("PUT", "/metadata", ioutil.NopCloser(bytes.NewBuffer(bts)))
	expect(t, err, nil)
	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	expect(t, rec.Code, http.StatusCreated)

	// checking mode
	metaValue, err := dbClient.MD.Get([]byte("some_key"))
	expect(t, err, nil)
	expect(t, string(metaValue), "some_val")
}

func TestSetMetadataBadBody(t *testing.T) {
	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.Cache.DeleteData()
	m := getBoneRouter(*dbClient)

	// deleting through handler
	req, err := http.NewRequest("PUT", "/metadata", ioutil.NopCloser(bytes.NewBuffer([]byte("you shall not decode me!!"))))
	expect(t, err, nil)
	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	expect(t, rec.Code, http.StatusBadRequest)
}

func TestSetMetadataMissingKey(t *testing.T) {
	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.Cache.DeleteData()
	m := getBoneRouter(*dbClient)

	// preparing to set mode through rest api
	var reqBody setMetadata
	// missing key
	reqBody.Value = "some_val"

	bts, err := json.Marshal(&reqBody)
	expect(t, err, nil)

	// deleting through handler
	req, err := http.NewRequest("PUT", "/metadata", ioutil.NopCloser(bytes.NewBuffer(bts)))
	expect(t, err, nil)
	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	expect(t, rec.Code, http.StatusBadRequest)

	// checking response body
	body, err := ioutil.ReadAll(rec.Body)
	mr := messageResponse{}
	err = json.Unmarshal(body, &mr)

	expect(t, mr.Message, "Key not provided.")
}

func TestGetMetadata(t *testing.T) {
	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.Cache.DeleteData()
	m := getBoneRouter(*dbClient)
	// adding some metadata
	for i := 0; i < 3; i++ {
		k := fmt.Sprintf("key_%d", i)
		v := fmt.Sprintf("val_%d", i)
		err := dbClient.MD.Set(k, v)
		expect(t, err, nil)
	}

	req, err := http.NewRequest("GET", "/metadata", nil)
	expect(t, err, nil)
	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	expect(t, rec.Code, http.StatusOK)

	body, err := ioutil.ReadAll(rec.Body)

	sm := storedMetadata{}
	err = json.Unmarshal(body, &sm)

	expect(t, len(sm.Data), 3)

	for i := 0; i < 3; i++ {
		k := fmt.Sprintf("key_%d", i)
		v := fmt.Sprintf("val_%d", i)
		expect(t, sm.Data[k], v)
	}
}

func TestDeleteMetadata(t *testing.T) {
	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.Cache.DeleteData()
	m := getBoneRouter(*dbClient)
	// adding some metadata
	for i := 0; i < 3; i++ {
		k := fmt.Sprintf("key_%d", i)
		v := fmt.Sprintf("val_%d", i)
		err := dbClient.MD.Set(k, v)
		expect(t, err, nil)
	}

	// checking that metadata is there
	allMeta, err := dbClient.MD.GetAll()
	expect(t, err, nil)
	expect(t, len(allMeta), 3)

	// deleting it
	req, err := http.NewRequest("DELETE", "/metadata", nil)
	expect(t, err, nil)
	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	expect(t, rec.Code, http.StatusOK)

	// checking metadata again, should be zero
	allMeta, err = dbClient.MD.GetAll()
	expect(t, err, nil)
	expect(t, len(allMeta), 0)
}

func TestDeleteMetadataEmpty(t *testing.T) {
	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.Cache.DeleteData()
	m := getBoneRouter(*dbClient)

	// deleting it
	req, err := http.NewRequest("DELETE", "/metadata", nil)
	expect(t, err, nil)
	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	expect(t, rec.Code, http.StatusOK)

	// checking metadata again, should be zero
	allMeta, err := dbClient.MD.GetAll()
	expect(t, err, nil)
	expect(t, len(allMeta), 0)
}

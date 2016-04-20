package hoverfly

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/SpectoLabs/hoverfly/testutil"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAllRecords(t *testing.T) {
	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()
	m := getBoneRouter(*dbClient)

	req, err := http.NewRequest("GET", "/api/records", nil)
	testutil.Expect(t, err, nil)

	//The response recorder used to record HTTP responses
	respRec := httptest.NewRecorder()

	m.ServeHTTP(respRec, req)

	testutil.Expect(t, respRec.Code, http.StatusOK)

	body, err := ioutil.ReadAll(respRec.Body)

	rr := recordedRequests{}
	err = json.Unmarshal(body, &rr)

	testutil.Expect(t, len(rr.Data), 0)
}

func TestGetAllRecordsWRecords(t *testing.T) {
	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()

	// inserting some payloads
	for i := 0; i < 5; i++ {
		req, err := http.NewRequest("GET", fmt.Sprintf("http://example.com/q=%d", i), nil)
		testutil.Expect(t, err, nil)
		dbClient.captureRequest(req)
	}
	// performing query
	m := getBoneRouter(*dbClient)

	req, err := http.NewRequest("GET", "/api/records", nil)
	testutil.Expect(t, err, nil)

	//The response recorder used to record HTTP responses
	respRec := httptest.NewRecorder()

	m.ServeHTTP(respRec, req)

	testutil.Expect(t, respRec.Code, http.StatusOK)

	body, err := ioutil.ReadAll(respRec.Body)

	rr := recordedRequests{}
	err = json.Unmarshal(body, &rr)

	testutil.Expect(t, len(rr.Data), 5)
}

func TestGetRecordsCount(t *testing.T) {
	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()
	m := getBoneRouter(*dbClient)

	req, err := http.NewRequest("GET", "/api/count", nil)
	testutil.Expect(t, err, nil)

	//The response recorder used to record HTTP responses
	respRec := httptest.NewRecorder()

	m.ServeHTTP(respRec, req)

	testutil.Expect(t, respRec.Code, http.StatusOK)

	body, err := ioutil.ReadAll(respRec.Body)

	rc := recordsCount{}
	err = json.Unmarshal(body, &rc)

	testutil.Expect(t, rc.Count, 0)
}

func TestGetRecordsCountWRecords(t *testing.T) {
	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()

	// inserting some payloads
	for i := 0; i < 5; i++ {
		req, err := http.NewRequest("GET", fmt.Sprintf("http://example.com/q=%d", i), nil)
		testutil.Expect(t, err, nil)
		dbClient.captureRequest(req)
	}
	// performing query
	m := getBoneRouter(*dbClient)

	req, err := http.NewRequest("GET", "/api/count", nil)
	testutil.Expect(t, err, nil)

	//The response recorder used to record HTTP responses
	respRec := httptest.NewRecorder()

	m.ServeHTTP(respRec, req)

	testutil.Expect(t, respRec.Code, http.StatusOK)

	body, err := ioutil.ReadAll(respRec.Body)

	rc := recordsCount{}
	err = json.Unmarshal(body, &rc)

	testutil.Expect(t, rc.Count, 5)
}

func TestExportImportRecords(t *testing.T) {
	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()
	m := getBoneRouter(*dbClient)

	// inserting some payloads
	for i := 0; i < 5; i++ {
		req, err := http.NewRequest("GET", fmt.Sprintf("http://example.com/q=%d", i), nil)
		testutil.Expect(t, err, nil)
		dbClient.captureRequest(req)
	}

	req, err := http.NewRequest("GET", "/api/records", nil)
	testutil.Expect(t, err, nil)

	//The response recorder used to record HTTP responses
	respRec := httptest.NewRecorder()

	m.ServeHTTP(respRec, req)

	testutil.Expect(t, respRec.Code, http.StatusOK)

	body, err := ioutil.ReadAll(respRec.Body)

	// deleting records
	err = dbClient.RequestCache.DeleteData()
	testutil.Expect(t, err, nil)

	// using body to import records again
	importReq, err := http.NewRequest("POST", "/api/records", ioutil.NopCloser(bytes.NewBuffer(body)))
	//The response recorder used to record HTTP responses
	importRec := httptest.NewRecorder()

	m.ServeHTTP(importRec, importReq)
	testutil.Expect(t, importRec.Code, http.StatusOK)

	// records should be there
	payloads, err := dbClient.RequestCache.GetAllValues()
	testutil.Expect(t, err, nil)
	testutil.Expect(t, len(payloads), 5)

}

func TestDeleteHandler(t *testing.T) {
	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()
	m := getBoneRouter(*dbClient)

	// inserting some payloads
	for i := 0; i < 5; i++ {
		req, err := http.NewRequest("GET", fmt.Sprintf("http://example.com/q=%d", i), nil)
		testutil.Expect(t, err, nil)
		dbClient.captureRequest(req)
	}

	// checking whether we have records
	payloads, err := dbClient.RequestCache.GetAllValues()
	testutil.Expect(t, err, nil)
	testutil.Expect(t, len(payloads), 5)

	// deleting through handler
	deleteReq, err := http.NewRequest("DELETE", "/api/records", nil)
	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, deleteReq)
	testutil.Expect(t, rec.Code, http.StatusOK)
}

func TestDeleteHandlerNoBucket(t *testing.T) {
	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()
	m := getBoneRouter(*dbClient)

	// deleting through handler
	importReq, err := http.NewRequest("DELETE", "/api/records", nil)
	testutil.Expect(t, err, nil)
	//The response recorder used to record HTTP responses
	importRec := httptest.NewRecorder()

	m.ServeHTTP(importRec, importReq)
	testutil.Expect(t, importRec.Code, http.StatusOK)
}

func TestGetState(t *testing.T) {
	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()
	m := getBoneRouter(*dbClient)

	// setting initial mode
	dbClient.Cfg.SetMode(VirtualizeMode)

	req, err := http.NewRequest("GET", "/api/state", nil)
	testutil.Expect(t, err, nil)
	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	testutil.Expect(t, rec.Code, http.StatusOK)

	body, err := ioutil.ReadAll(rec.Body)

	sr := stateRequest{}
	err = json.Unmarshal(body, &sr)

	testutil.Expect(t, sr.Mode, VirtualizeMode)
}

func TestSetVirtualizeState(t *testing.T) {
	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()
	m := getBoneRouter(*dbClient)

	// setting mode to capture
	dbClient.Cfg.SetMode("capture")

	// preparing to set mode through rest api
	var resp stateRequest
	resp.Mode = VirtualizeMode

	bts, err := json.Marshal(&resp)
	testutil.Expect(t, err, nil)

	// deleting through handler
	req, err := http.NewRequest("POST", "/api/state", ioutil.NopCloser(bytes.NewBuffer(bts)))
	testutil.Expect(t, err, nil)
	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	testutil.Expect(t, rec.Code, http.StatusOK)

	// checking mode
	testutil.Expect(t, dbClient.Cfg.GetMode(), VirtualizeMode)
}

func TestSetCaptureState(t *testing.T) {
	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()
	m := getBoneRouter(*dbClient)

	// setting mode to virtualize
	dbClient.Cfg.SetMode(VirtualizeMode)

	// preparing to set mode through rest api
	var resp stateRequest
	resp.Mode = "capture"

	bts, err := json.Marshal(&resp)
	testutil.Expect(t, err, nil)

	// deleting through handler
	req, err := http.NewRequest("POST", "/api/state", ioutil.NopCloser(bytes.NewBuffer(bts)))
	testutil.Expect(t, err, nil)
	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	testutil.Expect(t, rec.Code, http.StatusOK)

	// checking mode
	testutil.Expect(t, dbClient.Cfg.GetMode(), "capture")
}

func TestSetModifyState(t *testing.T) {
	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()
	m := getBoneRouter(*dbClient)

	// setting mode to virtualize
	dbClient.Cfg.SetMode(VirtualizeMode)

	// preparing to set mode through rest api
	var resp stateRequest
	resp.Mode = ModifyMode

	bts, err := json.Marshal(&resp)
	testutil.Expect(t, err, nil)

	// deleting through handler
	req, err := http.NewRequest("POST", "/api/state", ioutil.NopCloser(bytes.NewBuffer(bts)))
	testutil.Expect(t, err, nil)
	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	testutil.Expect(t, rec.Code, http.StatusOK)

	// checking mode
	testutil.Expect(t, dbClient.Cfg.GetMode(), ModifyMode)
}

func TestSetSynthesizeState(t *testing.T) {
	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()
	m := getBoneRouter(*dbClient)

	// setting mode to virtualize
	dbClient.Cfg.SetMode(VirtualizeMode)

	// preparing to set mode through rest api
	var resp stateRequest
	resp.Mode = SynthesizeMode

	bts, err := json.Marshal(&resp)
	testutil.Expect(t, err, nil)

	// deleting through handler
	req, err := http.NewRequest("POST", "/api/state", ioutil.NopCloser(bytes.NewBuffer(bts)))
	testutil.Expect(t, err, nil)
	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	testutil.Expect(t, rec.Code, http.StatusOK)

	// checking mode
	testutil.Expect(t, dbClient.Cfg.GetMode(), SynthesizeMode)
}

func TestSetRandomState(t *testing.T) {
	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()
	m := getBoneRouter(*dbClient)

	// setting mode to virtualize
	dbClient.Cfg.SetMode(VirtualizeMode)

	// preparing to set mode through rest api
	var resp stateRequest
	resp.Mode = "shouldnotwork"

	bts, err := json.Marshal(&resp)
	testutil.Expect(t, err, nil)

	// deleting through handler
	req, err := http.NewRequest("POST", "/api/state", ioutil.NopCloser(bytes.NewBuffer(bts)))
	testutil.Expect(t, err, nil)
	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	testutil.Expect(t, rec.Code, http.StatusBadRequest)

	// checking mode, should not have changed
	testutil.Expect(t, dbClient.Cfg.GetMode(), VirtualizeMode)
}

func TestSetNoBody(t *testing.T) {
	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()
	m := getBoneRouter(*dbClient)

	// setting mode to virtualize
	dbClient.Cfg.SetMode(VirtualizeMode)

	// setting state
	req, err := http.NewRequest("POST", "/api/state", nil)
	testutil.Expect(t, err, nil)
	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	testutil.Expect(t, rec.Code, http.StatusBadRequest)

	// checking mode, should not have changed
	testutil.Expect(t, dbClient.Cfg.GetMode(), VirtualizeMode)
}

func TestStatsHandler(t *testing.T) {
	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()
	m := getBoneRouter(*dbClient)

	// deleting through handler
	req, err := http.NewRequest("GET", "/api/stats", nil)

	testutil.Expect(t, err, nil)
	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	testutil.Expect(t, rec.Code, http.StatusOK)
}

func TestStatsHandlerVirtualizeMetrics(t *testing.T) {
	// test metrics, increases virtualize count by 1 and then checks through stats
	// handler whether it is visible through /stats handler
	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()
	m := getBoneRouter(*dbClient)

	dbClient.Counter.Counters[VirtualizeMode].Inc(1)

	req, err := http.NewRequest("GET", "/api/stats", nil)

	testutil.Expect(t, err, nil)
	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	testutil.Expect(t, rec.Code, http.StatusOK)

	body, err := ioutil.ReadAll(rec.Body)

	sr := statsResponse{}
	err = json.Unmarshal(body, &sr)

	testutil.Expect(t, int(sr.Stats.Counters[VirtualizeMode]), 1)
}

func TestStatsHandlerCaptureMetrics(t *testing.T) {
	// test metrics, increases capture count by 1 and then checks through stats
	// handler whether it is visible through /stats handler
	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()
	m := getBoneRouter(*dbClient)

	dbClient.Counter.Counters[CaptureMode].Inc(1)

	req, err := http.NewRequest("GET", "/api/stats", nil)

	testutil.Expect(t, err, nil)
	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	testutil.Expect(t, rec.Code, http.StatusOK)

	body, err := ioutil.ReadAll(rec.Body)

	sr := statsResponse{}
	err = json.Unmarshal(body, &sr)

	testutil.Expect(t, int(sr.Stats.Counters[CaptureMode]), 1)
}

func TestStatsHandlerModifyMetrics(t *testing.T) {
	// test metrics, increases modify count by 1 and then checks through stats
	// handler whether it is visible through /stats handler
	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()
	m := getBoneRouter(*dbClient)

	dbClient.Counter.Counters[ModifyMode].Inc(1)

	req, err := http.NewRequest("GET", "/api/stats", nil)

	testutil.Expect(t, err, nil)
	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	testutil.Expect(t, rec.Code, http.StatusOK)

	body, err := ioutil.ReadAll(rec.Body)

	sr := statsResponse{}
	err = json.Unmarshal(body, &sr)

	testutil.Expect(t, int(sr.Stats.Counters[ModifyMode]), 1)
}

func TestStatsHandlerSynthesizeMetrics(t *testing.T) {
	// test metrics, increases synthesize count by 1 and then checks through stats
	// handler whether it is visible through /stats handler
	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()
	m := getBoneRouter(*dbClient)

	dbClient.Counter.Counters[SynthesizeMode].Inc(1)

	req, err := http.NewRequest("GET", "/api/stats", nil)

	testutil.Expect(t, err, nil)
	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	testutil.Expect(t, rec.Code, http.StatusOK)

	body, err := ioutil.ReadAll(rec.Body)

	sr := statsResponse{}
	err = json.Unmarshal(body, &sr)

	testutil.Expect(t, int(sr.Stats.Counters[SynthesizeMode]), 1)
}

func TestStatsHandlerRecordCountMetrics(t *testing.T) {
	// test metrics, adds 5 new requests and then checks through stats
	// handler whether it is visible through /stats handler
	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()
	m := getBoneRouter(*dbClient)

	// inserting some payloads
	for i := 0; i < 5; i++ {
		req, err := http.NewRequest("GET", fmt.Sprintf("http://example.com/q=%d", i), nil)
		testutil.Expect(t, err, nil)
		dbClient.captureRequest(req)
	}

	req, err := http.NewRequest("GET", "/api/stats", nil)

	testutil.Expect(t, err, nil)
	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	testutil.Expect(t, rec.Code, http.StatusOK)

	body, err := ioutil.ReadAll(rec.Body)

	sr := statsResponse{}
	err = json.Unmarshal(body, &sr)

	testutil.Expect(t, int(sr.RecordsCount), 5)
}

func TestSetMetadata(t *testing.T) {
	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()
	m := getBoneRouter(*dbClient)

	// preparing to set mode through rest api
	var reqBody setMetadata
	reqBody.Key = "some_key"
	reqBody.Value = "some_val"

	bts, err := json.Marshal(&reqBody)
	testutil.Expect(t, err, nil)

	// deleting through handler
	req, err := http.NewRequest("PUT", "/api/metadata", ioutil.NopCloser(bytes.NewBuffer(bts)))
	testutil.Expect(t, err, nil)
	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	testutil.Expect(t, rec.Code, http.StatusCreated)

	// checking mode
	metaValue, err := dbClient.MetadataCache.Get([]byte("some_key"))
	testutil.Expect(t, err, nil)
	testutil.Expect(t, string(metaValue), "some_val")
}

func TestSetMetadataBadBody(t *testing.T) {
	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()
	m := getBoneRouter(*dbClient)

	// deleting through handler
	req, err := http.NewRequest("PUT", "/api/metadata", ioutil.NopCloser(bytes.NewBuffer([]byte("you shall not decode me!!"))))
	testutil.Expect(t, err, nil)
	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	testutil.Expect(t, rec.Code, http.StatusBadRequest)
}

func TestSetMetadataMissingKey(t *testing.T) {
	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()
	m := getBoneRouter(*dbClient)

	// preparing to set mode through rest api
	var reqBody setMetadata
	// missing key
	reqBody.Value = "some_val"

	bts, err := json.Marshal(&reqBody)
	testutil.Expect(t, err, nil)

	// deleting through handler
	req, err := http.NewRequest("PUT", "/api/metadata", ioutil.NopCloser(bytes.NewBuffer(bts)))
	testutil.Expect(t, err, nil)
	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	testutil.Expect(t, rec.Code, http.StatusBadRequest)

	// checking response body
	body, err := ioutil.ReadAll(rec.Body)
	mr := messageResponse{}
	err = json.Unmarshal(body, &mr)

	testutil.Expect(t, mr.Message, "Key not provided.")
}

func TestGetMetadata(t *testing.T) {
	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()
	m := getBoneRouter(*dbClient)
	// adding some metadata
	for i := 0; i < 3; i++ {
		k := fmt.Sprintf("key_%d", i)
		v := fmt.Sprintf("val_%d", i)
		err := dbClient.MetadataCache.Set([]byte(k), []byte(v))
		testutil.Expect(t, err, nil)
	}

	req, err := http.NewRequest("GET", "/api/metadata", nil)
	testutil.Expect(t, err, nil)
	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	testutil.Expect(t, rec.Code, http.StatusOK)

	body, err := ioutil.ReadAll(rec.Body)

	sm := storedMetadata{}
	err = json.Unmarshal(body, &sm)

	testutil.Expect(t, len(sm.Data), 3)

	for i := 0; i < 3; i++ {
		k := fmt.Sprintf("key_%d", i)
		v := fmt.Sprintf("val_%d", i)
		testutil.Expect(t, sm.Data[k], v)
	}
}

func TestDeleteMetadata(t *testing.T) {
	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()
	m := getBoneRouter(*dbClient)
	// adding some metadata
	for i := 0; i < 3; i++ {
		k := fmt.Sprintf("key_%d", i)
		v := fmt.Sprintf("val_%d", i)
		err := dbClient.MetadataCache.Set([]byte(k), []byte(v))
		testutil.Expect(t, err, nil)
	}

	// checking that metadata is there
	allMeta, err := dbClient.MetadataCache.GetAllEntries()
	testutil.Expect(t, err, nil)
	testutil.Expect(t, len(allMeta), 3)

	// deleting it
	req, err := http.NewRequest("DELETE", "/api/metadata", nil)
	testutil.Expect(t, err, nil)
	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	testutil.Expect(t, rec.Code, http.StatusOK)

	// checking metadata again, should be zero
	allMeta, err = dbClient.MetadataCache.GetAllEntries()
	testutil.Expect(t, err, nil)
	testutil.Expect(t, len(allMeta), 0)
}

func TestDeleteMetadataEmpty(t *testing.T) {
	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()
	m := getBoneRouter(*dbClient)

	// deleting it
	req, err := http.NewRequest("DELETE", "/api/metadata", nil)
	testutil.Expect(t, err, nil)
	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	testutil.Expect(t, rec.Code, http.StatusOK)

	// checking metadata again, should be zero
	allMeta, err := dbClient.MetadataCache.GetAllEntries()
	testutil.Expect(t, err, nil)
	testutil.Expect(t, len(allMeta), 0)
}

package hoverfly

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SpectoLabs/hoverfly/core/handlers/v1"
	"github.com/SpectoLabs/hoverfly/core/models"
	. "github.com/onsi/gomega"
)

var adminApi = AdminApi{}

func TestGetRecordsCount(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})
	m := adminApi.getBoneRouter(unit)

	req, err := http.NewRequest("GET", "/api/count", nil)
	Expect(err).To(BeNil())

	//The response recorder used to record HTTP responses
	respRec := httptest.NewRecorder()

	m.ServeHTTP(respRec, req)

	Expect(respRec.Code, http.StatusOK)

	body, err := ioutil.ReadAll(respRec.Body)

	rc := v1.RecordsCount{}
	err = json.Unmarshal(body, &rc)

	Expect(rc.Count).To(Equal(0))
}

func TestGetRecordsCountWRecords(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	// inserting some payloads
	for i := 0; i < 5; i++ {
		req := &models.RequestDetails{
			Method:      "GET",
			Scheme:      "http",
			Destination: "example.com",
			Query:       fmt.Sprintf("q=%d", i),
		}

		unit.Save(req, &models.ResponseDetails{})
	}
	// performing query
	m := adminApi.getBoneRouter(unit)

	req, err := http.NewRequest("GET", "/api/count", nil)
	Expect(err).To(BeNil())

	//The response recorder used to record HTTP responses
	respRec := httptest.NewRecorder()

	m.ServeHTTP(respRec, req)

	Expect(respRec.Code, http.StatusOK)

	body, err := ioutil.ReadAll(respRec.Body)

	rc := v1.RecordsCount{}
	err = json.Unmarshal(body, &rc)

	Expect(rc.Count).To(Equal(5))
}

func TestGetState(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})
	m := adminApi.getBoneRouter(unit)

	// setting initial mode
	unit.Cfg.SetMode("simulate")

	req, err := http.NewRequest("GET", "/api/state", nil)
	Expect(err).To(BeNil())
	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	Expect(rec.Code, http.StatusOK)

	body, err := ioutil.ReadAll(rec.Body)

	sr := v1.StateRequest{}
	err = json.Unmarshal(body, &sr)

	Expect(sr.Mode).To(Equal("simulate"))
}

func TestSetSimulateState(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})
	m := adminApi.getBoneRouter(unit)

	// setting mode to capture
	unit.Cfg.SetMode("capture")

	// preparing to set mode through rest api
	var resp v1.StateRequest
	resp.Mode = "simulate"

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
	Expect(unit.Cfg.GetMode()).To(Equal("simulate"))
}

func TestSetCaptureState(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})
	m := adminApi.getBoneRouter(unit)

	// setting mode to simulate
	unit.Cfg.SetMode("simulate")

	// preparing to set mode through rest api
	var resp v1.StateRequest
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
	Expect(unit.Cfg.GetMode()).To(Equal("capture"))
}

func TestSetModifyState(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})
	m := adminApi.getBoneRouter(unit)

	// setting mode to simulate
	unit.Cfg.SetMode("simulate")

	// preparing to set mode through rest api
	var resp v1.StateRequest
	resp.Mode = "modify"

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
	Expect(unit.Cfg.GetMode()).To(Equal("modify"))
}

func TestSetSynthesizeState(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})
	m := adminApi.getBoneRouter(unit)

	// setting mode to simulate
	unit.Cfg.SetMode("simulate")

	// preparing to set mode through rest api
	var resp v1.StateRequest
	resp.Mode = "synthesize"

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
	Expect(unit.Cfg.GetMode()).To(Equal("synthesize"))
}

func TestSetRandomState(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})
	m := adminApi.getBoneRouter(unit)

	// setting mode to simulate
	unit.Cfg.SetMode("simulate")

	// preparing to set mode through rest api
	var resp v1.StateRequest
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
	Expect(unit.Cfg.GetMode()).To(Equal("simulate"))
}

func TestSetNoBody(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})
	m := adminApi.getBoneRouter(unit)

	// setting mode to simulate
	unit.Cfg.SetMode("simulate")

	// setting state
	req, err := http.NewRequest("POST", "/api/state", nil)
	Expect(err).To(BeNil())

	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	Expect(rec.Code).To(Equal(http.StatusBadRequest))

	// checking mode, should not have changed
	Expect(unit.Cfg.GetMode()).To(Equal("simulate"))
}

func TestStatsHandler(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})
	m := adminApi.getBoneRouter(unit)

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
	unit := NewHoverflyWithConfiguration(&Configuration{})
	m := adminApi.getBoneRouter(unit)

	unit.Counter.Counters["simulate"].Inc(1)

	req, err := http.NewRequest("GET", "/api/stats", nil)
	Expect(err).To(BeNil())

	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	Expect(rec.Code).To(Equal(http.StatusOK))

	body, err := ioutil.ReadAll(rec.Body)

	sr := v1.StatsResponse{}
	err = json.Unmarshal(body, &sr)

	Expect(int(sr.Stats.Counters["simulate"])).To(Equal(1))
}

func TestStatsHandlerCaptureMetrics(t *testing.T) {
	RegisterTestingT(t)

	// test metrics, increases capture count by 1 and then checks through stats
	// handler whether it is visible through /stats handler
	unit := NewHoverflyWithConfiguration(&Configuration{})
	m := adminApi.getBoneRouter(unit)

	unit.Counter.Counters["capture"].Inc(1)

	req, err := http.NewRequest("GET", "/api/stats", nil)
	Expect(err).To(BeNil())

	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	Expect(rec.Code).To(Equal(http.StatusOK))

	body, err := ioutil.ReadAll(rec.Body)

	sr := v1.StatsResponse{}
	err = json.Unmarshal(body, &sr)

	Expect(int(sr.Stats.Counters["capture"])).To(Equal(1))
}

func TestStatsHandlerModifyMetrics(t *testing.T) {
	RegisterTestingT(t)

	// test metrics, increases modify count by 1 and then checks through stats
	// handler whether it is visible through /stats handler
	unit := NewHoverflyWithConfiguration(&Configuration{})
	m := adminApi.getBoneRouter(unit)

	unit.Counter.Counters["modify"].Inc(1)

	req, err := http.NewRequest("GET", "/api/stats", nil)
	Expect(err).To(BeNil())

	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	Expect(rec.Code).To(Equal(http.StatusOK))

	body, err := ioutil.ReadAll(rec.Body)

	sr := v1.StatsResponse{}
	err = json.Unmarshal(body, &sr)

	Expect(int(sr.Stats.Counters["modify"])).To(Equal(1))
}

func TestStatsHandlerSynthesizeMetrics(t *testing.T) {
	RegisterTestingT(t)

	// test metrics, increases synthesize count by 1 and then checks through stats
	// handler whether it is visible through /stats handler
	unit := NewHoverflyWithConfiguration(&Configuration{})
	m := adminApi.getBoneRouter(unit)

	unit.Counter.Counters["synthesize"].Inc(1)

	req, err := http.NewRequest("GET", "/api/stats", nil)
	Expect(err).To(BeNil())

	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	Expect(rec.Code).To(Equal(http.StatusOK))

	body, err := ioutil.ReadAll(rec.Body)

	sr := v1.StatsResponse{}
	err = json.Unmarshal(body, &sr)

	Expect(int(sr.Stats.Counters["synthesize"])).To(Equal(1))
}

func TestStatsHandlerRecordCountMetrics(t *testing.T) {
	RegisterTestingT(t)

	// test metrics, adds 5 new requests and then checks through stats
	// handler whether it is visible through /stats handler
	unit := NewHoverflyWithConfiguration(&Configuration{})
	m := adminApi.getBoneRouter(unit)

	// inserting some payloads
	for i := 0; i < 5; i++ {
		req := &models.RequestDetails{
			Method:      "GET",
			Scheme:      "http",
			Destination: "example.com",
			Query:       fmt.Sprintf("q=%d", i),
		}

		resp := &http.Response{}
		resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte("")))

		unit.Save(req, &models.ResponseDetails{})
	}

	req, err := http.NewRequest("GET", "/api/stats", nil)
	Expect(err).To(BeNil())

	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	Expect(rec.Code).To(Equal(http.StatusOK))

	body, err := ioutil.ReadAll(rec.Body)

	sr := v1.StatsResponse{}
	err = json.Unmarshal(body, &sr)

	Expect(int(sr.RecordsCount)).To(Equal(5))
}

func TestSetMetadata(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})
	m := adminApi.getBoneRouter(unit)

	// preparing to set mode through rest api
	var reqBody v1.SetMetadata
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
	metaValue, err := unit.MetadataCache.Get([]byte("some_key"))
	Expect(err).To(BeNil())
	Expect(string(metaValue)).To(Equal("some_val"))
}

func TestSetMetadataBadBody(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})
	m := adminApi.getBoneRouter(unit)

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

	unit := NewHoverflyWithConfiguration(&Configuration{})
	m := adminApi.getBoneRouter(unit)

	// preparing to set mode through rest api
	var reqBody v1.SetMetadata
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
	mr := v1.MessageResponse{}
	err = json.Unmarshal(body, &mr)

	Expect(mr.Message).To(Equal("Key not provided."))
}

func TestGetMetadata(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})
	m := adminApi.getBoneRouter(unit)
	// adding some metadata
	for i := 0; i < 3; i++ {
		k := fmt.Sprintf("key_%d", i)
		v := fmt.Sprintf("val_%d", i)
		err := unit.MetadataCache.Set([]byte(k), []byte(v))
		Expect(err).To(BeNil())
	}

	req, err := http.NewRequest("GET", "/api/metadata", nil)
	Expect(err).To(BeNil())

	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	Expect(rec.Code).To(Equal(http.StatusOK))

	body, err := ioutil.ReadAll(rec.Body)

	sm := v1.StoredMetadata{}
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

	unit := NewHoverflyWithConfiguration(&Configuration{})
	m := adminApi.getBoneRouter(unit)
	// adding some metadata
	for i := 0; i < 3; i++ {
		k := fmt.Sprintf("key_%d", i)
		v := fmt.Sprintf("val_%d", i)
		err := unit.MetadataCache.Set([]byte(k), []byte(v))
		Expect(err).To(BeNil())
	}

	// checking that metadata is there
	allMeta, err := unit.MetadataCache.GetAllEntries()
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
	allMeta, err = unit.MetadataCache.GetAllEntries()
	Expect(err).To(BeNil())
	Expect(len(allMeta)).To(Equal(0))
}

func TestDeleteMetadataEmpty(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})
	m := adminApi.getBoneRouter(unit)

	// deleting it
	req, err := http.NewRequest("DELETE", "/api/metadata", nil)
	Expect(err).To(BeNil())
	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	Expect(rec.Code).To(Equal(http.StatusOK))

	// checking metadata again, should be zero
	allMeta, err := unit.MetadataCache.GetAllEntries()
	Expect(err).To(BeNil())
	Expect(len(allMeta)).To(Equal(0))
}

func TestGetResponseDelays(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	delay := v1.ResponseDelayView{
		UrlPattern: ".",
		HttpMethod: "GET",
		Delay:      100,
	}
	delays := []v1.ResponseDelayView{delay}

	delaysPayload := v1.ResponseDelayPayloadView{
		Data: delays,
	}

	unit.SetResponseDelays(delaysPayload)

	m := adminApi.getBoneRouter(unit)

	req, err := http.NewRequest("GET", "/api/delays", nil)
	Expect(err).To(BeNil())
	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	Expect(rec.Code).To(Equal(http.StatusOK))

	body, err := ioutil.ReadAll(rec.Body)

	sr := v1.ResponseDelayPayloadView{}
	err = json.Unmarshal(body, &sr)

	// normal equality checking doesn't work on slices (!!)
	delayList := []v1.ResponseDelayView{{UrlPattern: ".", HttpMethod: "GET", Delay: 100}}
	Expect(sr.Data).To(Equal(delayList))
}

func TestDeleteAllResponseDelaysHandler(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})
	delay := models.ResponseDelay{
		UrlPattern: ".",
		Delay:      100,
	}
	delays := models.ResponseDelayList{delay}
	unit.ResponseDelays = &delays
	m := adminApi.getBoneRouter(unit)

	req, err := http.NewRequest("DELETE", "/api/delays", nil)
	Expect(err).To(BeNil())

	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	Expect(rec.Code).To(Equal(http.StatusOK))

	Expect(unit.ResponseDelays.Len()).To(Equal(0))
}

func TestUpdateResponseDelays(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})
	m := adminApi.getBoneRouter(unit)

	delayOne := v1.ResponseDelayView{
		UrlPattern: ".",
		Delay:      100,
	}
	delayTwo := v1.ResponseDelayView{
		UrlPattern: "example",
		Delay:      100,
	}
	delays := []v1.ResponseDelayView{delayOne, delayTwo}
	delayJson := v1.ResponseDelayPayloadView{Data: delays}
	delayJsonBytes, err := json.Marshal(&delayJson)
	Expect(err).To(BeNil())

	req, err := http.NewRequest("PUT", "/api/delays", ioutil.NopCloser(bytes.NewBuffer(delayJsonBytes)))
	Expect(err).To(BeNil())

	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	Expect(rec.Code).To(Equal(http.StatusCreated))

	Expect(unit.ResponseDelays.ConvertToResponseDelayPayloadView()).To(Equal(delayJson))
}

func TestInvalidJSONSyntaxUpdateResponseDelays(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})
	m := adminApi.getBoneRouter(unit)

	delayJson := "{aseuifhksejfc}"

	req, err := http.NewRequest("PUT", "/api/delays", ioutil.NopCloser(bytes.NewBuffer([]byte(delayJson))))
	Expect(err).To(BeNil())

	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	Expect(rec.Code).To(Equal(http.StatusBadRequest))

	// normal equality checking doesn't work on slices (!!)
	Expect(unit.ResponseDelays).To(Equal(&models.ResponseDelayList{}))
}

func TestInvalidJSONSemanticsUpdateResponseDelays(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})
	m := adminApi.getBoneRouter(unit)

	delayJson := "{ \"madeupfield\" : \"somevalue\" }"

	req, err := http.NewRequest("PUT", "/api/delays", ioutil.NopCloser(bytes.NewBuffer([]byte(delayJson))))
	Expect(err).To(BeNil())

	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	Expect(rec.Code).To(Equal(422))

	// normal equality checking doesn't work on slices (!!)
	Expect(unit.ResponseDelays).To(Equal(&models.ResponseDelayList{}))
}

func TestJSONWithInvalidHostPatternUpdateResponseDelays(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})
	m := adminApi.getBoneRouter(unit)

	delayJson := "{ \"data\": [{\"hostPattern\": \"*\", \"delay\": 100}] }"

	req, err := http.NewRequest("PUT", "/api/delays", ioutil.NopCloser(bytes.NewBuffer([]byte(delayJson))))
	Expect(err).To(BeNil())

	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	Expect(rec.Code).To(Equal(422))

	// normal equality checking doesn't work on slices (!!)
	Expect(unit.ResponseDelays).To(Equal(&models.ResponseDelayList{}))
}

func TestJSONWithMissingFieldUpdateResponseDelays(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})
	m := adminApi.getBoneRouter(unit)

	delayJson := "{ \"data\" : [{\"hostPattern\": \".\"}] }"

	req, err := http.NewRequest("PUT", "/api/delays", ioutil.NopCloser(bytes.NewBuffer([]byte(delayJson))))
	Expect(err).To(BeNil())

	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	m.ServeHTTP(rec, req)
	Expect(rec.Code).To(Equal(422))

	// normal equality checking doesn't work on slices (!!)
	Expect(unit.ResponseDelays).To(Equal(&models.ResponseDelayList{}))
}

func TestGetMissingURL(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})
	m := adminApi.getBoneRouter(unit)

	req, err := http.NewRequest("GET", "/api/sdiughvksjv", nil)
	Expect(err).To(BeNil())

	//The response recorder used to record HTTP responses
	respRec := httptest.NewRecorder()

	m.ServeHTTP(respRec, req)

	Expect(respRec.Code, http.StatusNotFound)
}

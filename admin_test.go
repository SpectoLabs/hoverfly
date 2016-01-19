package main

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
	err = dbClient.cache.DeleteBucket(dbClient.cache.requestsBucket)
	expect(t, err, nil)

	// using body to import records again
	importReq, err := http.NewRequest("POST", "/records", ioutil.NopCloser(bytes.NewBuffer(body)))
	//The response recorder used to record HTTP responses
	importRec := httptest.NewRecorder()

	m.ServeHTTP(importRec, importReq)
	expect(t, importRec.Code, http.StatusOK)

	// records should be overwritten
	payloads, err := dbClient.cache.GetAllRequests()
	expect(t, err, nil)
	expect(t, len(payloads), 5)

}

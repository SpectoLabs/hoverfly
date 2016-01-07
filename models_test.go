package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {

	setup()

	retCode := m.Run()

	// your func
	teardown()

	// call with result of m.Run()
	os.Exit(retCode)
}

// TestCaptureHeader tests whether request gets new header assigned
func TestCaptureHeader(t *testing.T) {

	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()

	req, err := http.NewRequest("GET", "http://example.com", nil)
	expect(t, err, nil)

	response, err := dbClient.captureRequest(req)

	expect(t, response.Header.Get("hoverfly"), "Was-Here")
}

// TestRequestBodyCaptured tests whether request body is recorded
func TestRequestBodyCaptured(t *testing.T) {

	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()

	body := ioutil.NopCloser(bytes.NewBuffer([]byte("fizz=buzz")))

	req, err := http.NewRequest("POST", "http://capture_body.com", body)
	expect(t, err, nil)

	_, err = dbClient.captureRequest(req)
	expect(t, err, nil)

	// since capture
	time.Sleep(10 * time.Millisecond)

	fp := getRequestFingerprint(req)

	payloadBts, err := dbClient.cache.Get([]byte(fp))

	var payload Payload

	expect(t, err, nil)

	// getting cache response
	err = json.Unmarshal(payloadBts, &payload)

	expect(t, err, nil)

	expect(t, payload.Request.Body, "fizz=buzz")

}

// TestRequestFingerprint tests whether we get correct request ID
func TestRequestFingerprint(t *testing.T) {

	req, err := http.NewRequest("GET", "http://example.com", nil)
	expect(t, err, nil)

	fp := getRequestFingerprint(req)

	expect(t, fp, "92a65ed4ca2b7100037a4cba9afd15ea")

}

func TestDeleteAllRecords(t *testing.T) {

	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	// inserting some payloads
	for i := 0; i < 5; i++ {
		req, err := http.NewRequest("GET", fmt.Sprintf("http://delete_all_records.com/q=%d", i), nil)
		expect(t, err, nil)
		dbClient.captureRequest(req)
	}
	err := dbClient.cache.DeleteBucket(dbClient.cache.requestsBucket)
	expect(t, err, nil)
}

func TestPayloadEncodeDecode(t *testing.T) {
	resp := response{
		Status: 200,
		Body:   "body here",
	}

	payload := Payload{Response: resp}

	bts, err := payload.encode()
	expect(t, err, nil)

	pl, err := decodePayload(bts)
	expect(t, err, nil)
	expect(t, pl.Response.Body, resp.Body)
	expect(t, pl.Response.Status, resp.Status)

}

func TestPayloadEncodeEmpty(t *testing.T) {
	payload := Payload{}

	bts, err := payload.encode()
	expect(t, err, nil)

	_, err = decodePayload(bts)
	expect(t, err, nil)
}


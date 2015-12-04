package main

import (
	"fmt"
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	retCode := m.Run()

	// your func
	teardown()

	// call with result of m.Run()
	os.Exit(retCode)
}

// TestRecordHeader tests whether request gets new header assigned
func TestRecordHeader(t *testing.T) {

	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.cache.pool.Close()

	req, err := http.NewRequest("GET", "http://example.com", nil)
	expect(t, err, nil)

	response, err := dbClient.captureRequest(req)

	expect(t, response.Header.Get("Gen-proxy"), "Was-Here")
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
	defer dbClient.cache.pool.Close()

	// inserting some payloads
	for i := 0; i < 5; i++ {
		req, err := http.NewRequest("GET", fmt.Sprintf("http://example.com/q=%d", i), nil)
		expect(t, err, nil)
		dbClient.captureRequest(req)
	}
	// checking that keys are there
	keys, _ := dbClient.cache.getAllKeys()
	expect(t, len(keys) > 0, true)

	// deleting
	err := dbClient.deleteAllRecords()
	expect(t, err, nil)

	// checking whether all records were deleted
	keys, _ = dbClient.cache.getAllKeys()
	expect(t, len(keys), 0)
}

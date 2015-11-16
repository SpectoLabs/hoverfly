package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
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

	response, err := dbClient.recordRequest(req)

	expect(t, response.Header.Get("Gen-proxy"), "Was-Here")
}

// TestRecordingToCache tests cache wrapper get/set/delete operations
func TestRecordingToCache(t *testing.T) {

	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.cache.pool.Close()

	dbClient.cache.set("some_key", "value")

	value, err := redis.String(dbClient.cache.get("some_key"))

	expect(t, err, nil)

	expect(t, string(value), "value")

	err = dbClient.cache.delete("some_key")

	expect(t, err, nil)
}

// TestRequestFingerprint tests whether we get correct request ID
func TestRequestFingerprint(t *testing.T) {

	req, err := http.NewRequest("GET", "http://example.com", nil)
	expect(t, err, nil)

	fp := getRequestFingerprint(req)

	expect(t, fp, "92a65ed4ca2b7100037a4cba9afd15ea")

}

// TestGetAllRecords - tests recording and then getting responses
func TestGetAllRecords(t *testing.T) {

	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.cache.pool.Close()

	// inserting some payloads
	for i := 0; i < 5; i++ {
		req, err := http.NewRequest("GET", fmt.Sprintf("http://example.com/q=%d", i), nil)
		expect(t, err, nil)
		dbClient.recordRequest(req)
	}

	// getting all keys
	keys, _ := dbClient.cache.getAllKeys()
	expect(t, len(keys) > 0, true)
	// getting requests
	payloads, err := dbClient.getAllRecords()
	expect(t, err, nil)

	for _, payload := range payloads {
		expect(t, payload.Request.Method, "GET")
		expect(t, payload.Response.Status, 201)
	}

}

func TestDeleteAllRecords(t *testing.T) {

	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.cache.pool.Close()

	// inserting some payloads
	for i := 0; i < 5; i++ {
		req, err := http.NewRequest("GET", fmt.Sprintf("http://example.com/q=%d", i), nil)
		expect(t, err, nil)
		dbClient.recordRequest(req)
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

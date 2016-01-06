package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

// TestRecordingToCache tests cache wrapper get/set/delete operations
func TestRecordingToCache(t *testing.T) {

	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.cache.pool.Close()

	dbClient.cache.set("some_key", []byte("value"))

	value, err := dbClient.cache.get("some_key")

	expect(t, err, nil)

	expect(t, string(value), "value")

	err = dbClient.cache.delete("some_key")

	expect(t, err, nil)
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
		dbClient.captureRequest(req)
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

	dbClient.cache.DeleteBucket(dbClient.cache.requestsBucket)
}

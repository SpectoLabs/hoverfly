package main

import (
	"net/http"
	"testing"
)

func TestRecordHeader(t *testing.T) {
	server, client := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer client.cache.pool.Close()

	req, err := http.NewRequest("GET", "http://example.com", nil)
	expect(t, err, nil)

	response, err := client.recordRequest(req)

	expect(t, response.Header.Get("Gen-proxy"), "Was-Here")
}

// TestRecordingToCache tests cache get/set operations
func TestRecordingToCache(t *testing.T) {
	server, client := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer client.cache.pool.Close()

	client.cache.set("some_key", "value")

	value, err := client.cache.get("some_key")

	expect(t, err, nil)

	if err == nil {
		expect(t, value, "value")
	}
}

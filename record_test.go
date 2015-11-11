package main

import (
	"github.com/garyburd/redigo/redis"
	"net/http"
	"testing"
)

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

// TestRecordingToCache tests cache wrapper get/set operations
func TestRecordingToCache(t *testing.T) {

	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.cache.pool.Close()

	dbClient.cache.set("some_key", "value")

	value, err := redis.String(dbClient.cache.get("some_key"))

	expect(t, err, nil)

	if err == nil {
		expect(t, string(value), "value")
	}
}

// TestRequestFingerprint tests whether we get correct request ID
func TestRequestFingerprint(t *testing.T) {

	req, err := http.NewRequest("GET", "http://example.com", nil)
	expect(t, err, nil)

	fp := getRequestFingerprint(req)

	expect(t, fp, "92a65ed4ca2b7100037a4cba9afd15ea")

}

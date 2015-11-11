package main

import (
	"github.com/garyburd/redigo/redis"
	"net/http"
	"testing"
)

func TestRecordHeader(t *testing.T) {
	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.cache.pool.Close()

	req, err := http.NewRequest("GET", "http://example.com", nil)
	expect(t, err, nil)

	response, err := dbClient.recordRequest(req)

	expect(t, response.Header.Get("Gen-proxy"), "Was-Here")
}

// TestRecordingToCache tests cache get/set operations
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

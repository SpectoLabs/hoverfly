package main

import (
	"net/http"
	"testing"
)

func TestGetNewHoverflyCheckConfig(t *testing.T) {

	cfg := InitSettings()
	_, dbClient := getNewHoverfly(cfg)
	defer dbClient.cache.db.Close()

	expect(t, dbClient.cfg, cfg)
}

func TestProcessCaptureRequest(t *testing.T) {
	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.cache.DeleteBucket(dbClient.cache.requestsBucket)

	r, err := http.NewRequest("GET", "http://somehost.com", nil)
	expect(t, err, nil)

	dbClient.cfg.SetMode("capture")

	req, resp := dbClient.processRequest(r)

	refute(t, req, nil)
	refute(t, resp, nil)
	expect(t, resp.StatusCode, 201)
}

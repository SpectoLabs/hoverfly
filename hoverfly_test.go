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

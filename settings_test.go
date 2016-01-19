package hoverfly
package main

import (
	"os"
	"testing"
)

func TestSettingsAdminPortEnv(t *testing.T) {
	defer os.Setenv("AdminPort", "")

	os.Setenv("AdminPort", "5555")

	cfg := InitSettings()
	expect(t, cfg.adminPort, "5555")
}


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

func TestSettingsDefaultAdminPort(t *testing.T) {
	os.Setenv("AdminPort", "")
	cfg := InitSettings()
	expect(t, cfg.adminPort, DefaultAdminPort)
}

func TestSettingsProxyPortEnv(t *testing.T) {
	defer os.Setenv("ProxyPort", "")

	os.Setenv("ProxyPort", "6666")
	cfg := InitSettings()

	expect(t, cfg.proxyPort, "6666")
}

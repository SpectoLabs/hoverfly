package hoverfly

import (
	"os"
	"testing"
)

func TestSettingsAdminPortEnv(t *testing.T) {
	defer os.Setenv("AdminPort", "")

	os.Setenv("AdminPort", "5555")

	cfg := InitSettings()
	expect(t, cfg.AdminPort, "5555")
}

func TestSettingsDefaultAdminPort(t *testing.T) {
	os.Setenv("AdminPort", "")
	cfg := InitSettings()
	expect(t, cfg.AdminPort, DefaultAdminPort)
}

func TestSettingsProxyPortEnv(t *testing.T) {
	defer os.Setenv("ProxyPort", "")

	os.Setenv("ProxyPort", "6666")
	cfg := InitSettings()

	expect(t, cfg.ProxyPort, "6666")
}

func TestSettingsDefaultProxyPort(t *testing.T) {
	os.Setenv("ProxyPort", "")
	cfg := InitSettings()
	expect(t, cfg.ProxyPort, DefaultPort)
}

func TestSettingsDatabaseEnv(t *testing.T) {
	defer os.Setenv("HoverflyDB", "")

	os.Setenv("HoverflyDB", "testingX.db")
	cfg := InitSettings()

	expect(t, cfg.DatabaseName, "testingX.db")
}

func TestSettingsMiddlewareEnv(t *testing.T) {
	defer os.Setenv("HoverflyMiddleware", "")

	os.Setenv("HoverflyMiddleware", "./examples/middleware/x.go")
	cfg := InitSettings()

	expect(t, cfg.Middleware, "./examples/middleware/x.go")
}

// TestSetMode - tests SetMode function, however it doesn't test
// whether mutex works correctly or not
func TestSetMode(t *testing.T) {

	cfg := Configuration{}
	cfg.SetMode("virtualize")
	expect(t, cfg.Mode, "virtualize")
}

// TestGetMode - tests GetMode function, however it doesn't test
// whether mutex works correctly or not
func TestGetMode(t *testing.T) {
	cfg := Configuration{Mode: "capture"}

	expect(t, cfg.GetMode(), "capture")
}

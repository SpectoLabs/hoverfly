package hoverfly

import (
	"github.com/SpectoLabs/hoverfly/core/testutil"
	"os"
	"testing"
)

func TestSettingsAdminPortEnv(t *testing.T) {
	defer os.Setenv("AdminPort", "")

	os.Setenv("AdminPort", "5555")

	cfg := InitSettings()
	testutil.Expect(t, cfg.AdminPort, "5555")
}

func TestSettingsDefaultAdminPort(t *testing.T) {
	os.Setenv("AdminPort", "")
	cfg := InitSettings()
	testutil.Expect(t, cfg.AdminPort, DefaultAdminPort)
}

func TestSettingsProxyPortEnv(t *testing.T) {
	defer os.Setenv("ProxyPort", "")

	os.Setenv("ProxyPort", "6666")
	cfg := InitSettings()

	testutil.Expect(t, cfg.ProxyPort, "6666")
}

func TestSettingsDefaultProxyPort(t *testing.T) {
	os.Setenv("ProxyPort", "")
	cfg := InitSettings()
	testutil.Expect(t, cfg.ProxyPort, DefaultPort)
}

func TestSettingsMiddlewareEnv(t *testing.T) {
	defer os.Setenv("HoverflyMiddleware", "")

	os.Setenv("HoverflyMiddleware", "./examples/middleware/x.go")
	cfg := InitSettings()

	testutil.Expect(t, cfg.Middleware, "./examples/middleware/x.go")
}

// TestSetMode - tests SetMode function, however it doesn't test
// whether mutex works correctly or not
func TestSetMode(t *testing.T) {

	cfg := Configuration{}
	cfg.SetMode(SimulateMode)
	testutil.Expect(t, cfg.Mode, SimulateMode)
}

// TestGetMode - tests GetMode function, however it doesn't test
// whether mutex works correctly or not
func TestGetMode(t *testing.T) {
	cfg := Configuration{Mode: "capture"}

	testutil.Expect(t, cfg.GetMode(), "capture")
}

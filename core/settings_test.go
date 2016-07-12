package hoverfly

import (
	"github.com/SpectoLabs/hoverfly/core/models"
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

func TestGetDelayWithRegexMatch(t *testing.T) {
	delay := models.ResponseDelay{
		HostPattern: "example",
		Delay:       100,
	}
	delays := []models.ResponseDelay{delay}
	cfg := Configuration{ResponseDelays: delays}

	delayMatch := cfg.GetDelay("delayexample.com")
	testutil.Expect(t, *delayMatch, delay)

	delayMatch = cfg.GetDelay("nodelay.com")
	var nilDelay *models.ResponseDelay
	testutil.Expect(t, delayMatch, nilDelay)
}

func TestMultipleMatchingDelaysReturnsTheFirst(t *testing.T) {
	delayOne := models.ResponseDelay{
		HostPattern: "example.com",
		Delay:       100,
	}
	delayTwo := models.ResponseDelay{
		HostPattern: "example",
		Delay:       100,
	}
	delays := []models.ResponseDelay{delayOne, delayTwo}
	cfg := Configuration{ResponseDelays: delays}

	delayMatch := cfg.GetDelay("delayexample.com")
	testutil.Expect(t, *delayMatch, delayOne)
}

func Test_InitSettings_SetsTheWebserverFieldToFalse(t *testing.T) {
	unit := InitSettings()
	testutil.Expect(t, unit.Webserver, false)
}

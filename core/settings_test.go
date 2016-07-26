package hoverfly

import (
	. "github.com/onsi/gomega"
	"os"
	"testing"
)

func TestSettingsAdminPortEnv(t *testing.T) {
	RegisterTestingT(t)

	defer os.Setenv("AdminPort", "")

	os.Setenv("AdminPort", "5555")

	cfg := InitSettings()
	Expect(cfg.AdminPort).To(Equal("5555"))
}

func TestSettingsDefaultAdminPort(t *testing.T) {
	RegisterTestingT(t)

	os.Setenv("AdminPort", "")
	cfg := InitSettings()
	Expect(cfg.AdminPort).To(Equal(DefaultAdminPort))
}

func TestSettingsProxyPortEnv(t *testing.T) {
	RegisterTestingT(t)

	defer os.Setenv("ProxyPort", "")

	os.Setenv("ProxyPort", "6666")
	cfg := InitSettings()

	Expect(cfg.ProxyPort).To(Equal("6666"))
}

func TestSettingsDefaultProxyPort(t *testing.T) {
	RegisterTestingT(t)

	os.Setenv("ProxyPort", "")
	cfg := InitSettings()
	Expect(cfg.ProxyPort).To(Equal(DefaultPort))
}

func TestSettingsMiddlewareEnv(t *testing.T) {
	RegisterTestingT(t)

	defer os.Setenv("HoverflyMiddleware", "")

	os.Setenv("HoverflyMiddleware", "./examples/middleware/x.go")
	cfg := InitSettings()

	Expect(cfg.Middleware).To(Equal("./examples/middleware/x.go"))
}

// TestSetMode - tests SetMode function, however it doesn't test
// whether mutex works correctly or not
func TestSetMode(t *testing.T) {
	RegisterTestingT(t)

	cfg := Configuration{}
	cfg.SetMode(SimulateMode)
	Expect(cfg.Mode).To(Equal(SimulateMode))
}

// TestGetMode - tests GetMode function, however it doesn't test
// whether mutex works correctly or not
func TestGetMode(t *testing.T) {
	cfg := Configuration{Mode: "capture"}

	Expect(cfg.GetMode()).To(Equal("capture"))
}

func Test_InitSettings_SetsTheWebserverFieldToFalse(t *testing.T) {
	unit := InitSettings()
	Expect(unit.Webserver).To(BeFalse())
}

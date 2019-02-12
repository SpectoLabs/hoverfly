package hoverfly

import (
	"os"
	"testing"

	. "github.com/onsi/gomega"
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

func TestSettingsDefaultListenOnHost(t *testing.T) {
	RegisterTestingT(t)

	cfg := InitSettings()
	Expect(cfg.ListenOnHost).To(Equal("127.0.0.1"))
}

func TestSettingsMiddlewareEnv(t *testing.T) {
	RegisterTestingT(t)

	defer os.Setenv("HoverflyMiddleware", "")

	os.Setenv("HoverflyMiddleware", "ruby ../examples/middleware/modify_response/modify_response.rb")
	cfg := InitSettings()

	Expect(cfg.Middleware.Binary).To(Equal("ruby"))

	script, err := cfg.Middleware.GetScript()
	Expect(err).To(BeNil())

	Expect(script).To(Equal(rubyModifyResponse))
}

func Test_InitSettings_SetsModeToSimulate(t *testing.T) {
	RegisterTestingT(t)

	settings := InitSettings()

	Expect(settings.Mode).To(Equal("simulate"))
}

// TestSetMode - tests SetMode function, however it doesn't test
// whether mutex works correctly or not
func TestSetMode(t *testing.T) {
	RegisterTestingT(t)

	cfg := Configuration{}
	cfg.SetMode("simulate")
	Expect(cfg.Mode).To(Equal("simulate"))
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

func Test_SetUpstreamProxy_WillPrependHttpColonSlashSlashToProxyURLWithoutIt(t *testing.T) {
	RegisterTestingT(t)

	unit := InitSettings()

	unit.SetUpstreamProxy("localhost")

	Expect(unit.UpstreamProxy).To(Equal("http://localhost"))
}

func Test_SetUpstreamProxy_WillNotPrependHttpColonSlashSlashToProxyURLWithIt(t *testing.T) {
	RegisterTestingT(t)

	unit := InitSettings()

	unit.SetUpstreamProxy("http://localhost")

	Expect(unit.UpstreamProxy).To(Equal("http://localhost"))
}

func Test_SetUpstreamProxy_WillNotPrependHttpColonSlashSlashToProxyURLWithHTTPS(t *testing.T) {
	RegisterTestingT(t)

	unit := InitSettings()

	unit.SetUpstreamProxy("https://localhost")

	Expect(unit.UpstreamProxy).To(Equal("https://localhost"))
}

func Test_InitSettings_SetsPlainHttpTunnelingToFalse(t *testing.T) {
	RegisterTestingT(t)

	settings := InitSettings()

	Expect(settings.PlainHttpTunneling).To(Equal(false))
}


func Test_InitSettings_SetsDefaultCacheSize(t *testing.T) {
	RegisterTestingT(t)

	settings := InitSettings()

	Expect(settings.CacheSize).To(Equal(1000))
}
package wrapper

import (
	"io/ioutil"
	"os"
	"testing"

	. "github.com/onsi/gomega"
)

var (
	defaultConfig = Config{
		HoverflyHost:          "localhost",
		HoverflyAdminPort:     "8888",
		HoverflyProxyPort:     "8500",
		HoverflyDbType:        "memory",
		HoverflyWebserver:     false,
		HoverflyUsername:      "",
		HoverflyPassword:      "",
		HoverflyCertificate:   "",
		HoverflyKey:           "",
		HoverflyDisableTls:    false,
		HoverflyUpstreamProxy: "",
		HoverflyCacheDisable:  false,
		Targets:               map[string]Target{},
	}
)

func Test_GetConfigWillReturnTheDefaultValues(t *testing.T) {
	RegisterTestingT(t)

	SetConfigurationDefaults()
	result := GetConfig()

	Expect(*result).To(Equal(defaultConfig))
}

func Test_Config_SetHost_OverridesDefaultValueWithAHoverflyHost(t *testing.T) {
	RegisterTestingT(t)

	SetConfigurationDefaults()
	result := GetConfig().SetHost("testhost")

	expected := defaultConfig
	expected.HoverflyHost = "testhost"

	Expect(*result).To(Equal(expected))
}

func Test_Config_SetHost_DoesNotOverrideWhenEmpty(t *testing.T) {
	RegisterTestingT(t)

	SetConfigurationDefaults()
	result := GetConfig().SetHost("")

	expected := defaultConfig

	Expect(*result).To(Equal(expected))
}

func Test_Config_SetAdminPort_OverridesDefaultValueWithAHoverflyAdminPort(t *testing.T) {
	RegisterTestingT(t)

	SetConfigurationDefaults()
	result := GetConfig().SetAdminPort("5")

	expected := defaultConfig
	expected.HoverflyAdminPort = "5"

	Expect(*result).To(Equal(expected))
}

func Test_Config_SetAdminPort_DoesNotOverrideWhenEmpty(t *testing.T) {
	RegisterTestingT(t)

	SetConfigurationDefaults()
	result := GetConfig().SetAdminPort("")

	expected := defaultConfig

	Expect(*result).To(Equal(expected))
}

func Test_Config_SetProxyPort_OverridesDefaultValueWithAHoverflyProxyPort(t *testing.T) {
	RegisterTestingT(t)

	SetConfigurationDefaults()
	result := GetConfig().SetProxyPort("7")

	expected := defaultConfig
	expected.HoverflyProxyPort = "7"

	Expect(*result).To(Equal(expected))
}

func Test_Config_SetProxyPort_DoesNotOverrideWhenEmpty(t *testing.T) {
	RegisterTestingT(t)

	SetConfigurationDefaults()
	result := GetConfig().SetProxyPort("")

	expected := defaultConfig

	Expect(*result).To(Equal(expected))
}

func Test_Config_SetDbType_OverridesDefaultValueWithAHoverflyProxyPort(t *testing.T) {
	RegisterTestingT(t)

	SetConfigurationDefaults()
	result := GetConfig().SetDbType("boltdb")

	expected := defaultConfig
	expected.HoverflyDbType = "boltdb"

	Expect(*result).To(Equal(expected))
}

func Test_Config_SetDbType_DoesNotOverrideWhenEmpty(t *testing.T) {
	RegisterTestingT(t)

	SetConfigurationDefaults()
	result := GetConfig().SetProxyPort("")

	expected := defaultConfig

	Expect(*result).To(Equal(expected))
}

func Test_Config_SetUsername_OverridesDefaultValue(t *testing.T) {
	RegisterTestingT(t)

	SetConfigurationDefaults()
	result := GetConfig().SetUsername("benjih")

	expected := defaultConfig
	expected.HoverflyUsername = "benjih"

	Expect(*result).To(Equal(expected))
}

func Test_Config_SetUsername_DoesNotOverrideWhenEmpty(t *testing.T) {
	RegisterTestingT(t)

	SetConfigurationDefaults()
	result := GetConfig().SetUsername("")

	expected := defaultConfig

	Expect(*result).To(Equal(expected))
}

func Test_Config_SetPassword_OverridesDefaultValue(t *testing.T) {
	RegisterTestingT(t)

	SetConfigurationDefaults()
	result := GetConfig().SetPassword("burger-toucher")

	expected := defaultConfig
	expected.HoverflyPassword = "burger-toucher"

	Expect(*result).To(Equal(expected))
}

func Test_Config_SetPassword_DoesNotOverrideWhenEmpty(t *testing.T) {
	RegisterTestingT(t)

	SetConfigurationDefaults()
	result := GetConfig().SetPassword("")

	expected := defaultConfig

	Expect(*result).To(Equal(expected))
}

func Test_Config_SetCertificate_OverridesDefaultValue(t *testing.T) {
	RegisterTestingT(t)

	SetConfigurationDefaults()
	result := GetConfig().SetCertificate("/home/benjih/test/certificate.pem")

	expected := defaultConfig
	expected.HoverflyCertificate = "/home/benjih/test/certificate.pem"

	Expect(*result).To(Equal(expected))
}

func Test_Config_SetKey_DoesNotOverrideWhenEmpty(t *testing.T) {
	RegisterTestingT(t)

	SetConfigurationDefaults()
	result := GetConfig().SetKey("/home/benjih/test/key.pem")

	expected := defaultConfig
	expected.HoverflyKey = "/home/benjih/test/key.pem"

	Expect(*result).To(Equal(expected))
}

func Test_Config_SetWebserver_OverridesDefaultValue(t *testing.T) {
	RegisterTestingT(t)

	SetConfigurationDefaults()
	result := GetConfig().SetWebserver("webserver")

	expected := defaultConfig
	expected.HoverflyWebserver = true

	Expect(*result).To(Equal(expected))
}

func Test_Config_SetWebserver_DoesNotOverrideIfEmptyString(t *testing.T) {
	RegisterTestingT(t)

	SetConfigurationDefaults()
	result := GetConfig().SetWebserver("")

	expected := defaultConfig
	expected.HoverflyWebserver = false

	Expect(*result).To(Equal(expected))
}

func Test_Config_SetWebserver_DoesNotOverrideIfNotGivenProxyOrWebserver(t *testing.T) {
	RegisterTestingT(t)

	SetConfigurationDefaults()

	result := GetConfig()
	result.HoverflyWebserver = true
	result.SetWebserver("not-webserver")

	expected := defaultConfig
	expected.HoverflyWebserver = true

	Expect(*result).To(Equal(expected))
}

func Test_Config_SetWebserver_OverrideIfGivenProxy(t *testing.T) {
	RegisterTestingT(t)

	SetConfigurationDefaults()

	result := GetConfig()
	result.HoverflyWebserver = true
	result.SetWebserver("proxy")

	expected := defaultConfig
	expected.HoverflyWebserver = false

	Expect(*result).To(Equal(expected))
}

func Test_Config_DisableTls_OverridesDefaultValue(t *testing.T) {
	RegisterTestingT(t)

	SetConfigurationDefaults()
	result := GetConfig().DisableTls(true)

	expected := defaultConfig
	expected.HoverflyDisableTls = true

	Expect(*result).To(Equal(expected))
}

func Test_Config_DisableTls_DoesNotOverridesDefaultValueIfDefaultIsPositive(t *testing.T) {
	RegisterTestingT(t)

	SetConfigurationDefaults()

	result := GetConfig()
	result.HoverflyDisableTls = true
	result = result.DisableTls(false)

	expected := defaultConfig
	expected.HoverflyDisableTls = true

	Expect(*result).To(Equal(expected))
}

func Test_Config_SetUpstreamProxy_OverridesDefaultValue(t *testing.T) {
	RegisterTestingT(t)

	SetConfigurationDefaults()
	result := GetConfig().SetUpstreamProxy("hoverfly.io:8080")

	expected := defaultConfig
	expected.HoverflyUpstreamProxy = "hoverfly.io:8080"

	Expect(*result).To(Equal(expected))
}

func Test_Config_SetUpstreamProxy_DoesNotOverrideWhenEmpty(t *testing.T) {
	RegisterTestingT(t)

	SetConfigurationDefaults()
	result := GetConfig().SetUpstreamProxy("")

	expected := defaultConfig

	Expect(*result).To(Equal(expected))
}

func Test_Config_DisableCache_OverridesDefaultValue(t *testing.T) {
	RegisterTestingT(t)

	SetConfigurationDefaults()
	result := GetConfig().DisableCache(true)

	expected := defaultConfig
	expected.HoverflyCacheDisable = true

	Expect(*result).To(Equal(expected))
}

func Test_Config_DisableCache_DoesNotOverridesDefaultValueIfDefaultIsPositive(t *testing.T) {
	RegisterTestingT(t)

	SetConfigurationDefaults()

	result := GetConfig()
	result.HoverflyCacheDisable = true
	result = result.DisableCache(false)

	expected := defaultConfig
	expected.HoverflyCacheDisable = true

	Expect(*result).To(Equal(expected))
}

func Test_Config_WriteToFile_WritesTheConfigObjectToAFileInAYamlFormat(t *testing.T) {
	RegisterTestingT(t)

	config := Config{
		HoverflyHost:         "testhost",
		HoverflyAdminPort:    "1234",
		HoverflyProxyPort:    "4567",
		HoverflyDbType:       "boltdb",
		HoverflyUsername:     "username",
		HoverflyPassword:     "password",
		HoverflyWebserver:    true,
		HoverflyCertificate:  "/home/benjih/certificate.pem",
		HoverflyKey:          "/home/benjih/key.pem",
		HoverflyDisableTls:   true,
		HoverflyCacheDisable: true,
	}

	wd, _ := os.Getwd()
	hoverflyDirectory := HoverflyDirectory{
		Path: wd,
	}

	err := config.WriteToFile(hoverflyDirectory)

	Expect(err).To(BeNil())

	data, _ := ioutil.ReadFile(hoverflyDirectory.Path + "/config.yaml")
	os.Remove(hoverflyDirectory.Path + "/config.yaml")

	Expect(string(data)).To(ContainSubstring(`hoverfly.host: testhost`))
	Expect(string(data)).To(ContainSubstring("hoverfly.admin.port: \"1234\""))
	Expect(string(data)).To(ContainSubstring("hoverfly.proxy.port: \"4567\""))
	Expect(string(data)).To(ContainSubstring("hoverfly.db.type: boltdb"))
	Expect(string(data)).To(ContainSubstring("hoverfly.username: username"))
	Expect(string(data)).To(ContainSubstring("hoverfly.password: password"))
	Expect(string(data)).To(ContainSubstring("hoverfly.webserver: true"))
	Expect(string(data)).To(ContainSubstring("hoverfly.tls.certificate: /home/benjih/certificate.pem"))
	Expect(string(data)).To(ContainSubstring("hoverfly.tls.key: /home/benjih/key.pem"))
	Expect(string(data)).To(ContainSubstring("hoverfly.tls.disable: true"))
	Expect(string(data)).To(ContainSubstring("hoverfly.cache.disable: true"))
}

func Test_Config_WriteToFile_WritesTheDefaultConfigObjectToAFileInAYamlFormat(t *testing.T) {
	RegisterTestingT(t)

	SetConfigurationDefaults()
	config := GetConfig()

	wd, _ := os.Getwd()
	hoverflyDirectory := HoverflyDirectory{
		Path: wd,
	}

	err := config.WriteToFile(hoverflyDirectory)

	Expect(err).To(BeNil())

	data, _ := ioutil.ReadFile(hoverflyDirectory.Path + "/config.yaml")
	os.Remove(hoverflyDirectory.Path + "/config.yaml")

	Expect(string(data)).To(ContainSubstring(`hoverfly.host: localhost`))
	Expect(string(data)).To(ContainSubstring("hoverfly.admin.port: \"8888\""))
	Expect(string(data)).To(ContainSubstring("hoverfly.proxy.port: \"8500\""))
	Expect(string(data)).To(ContainSubstring("hoverfly.db.type: memory"))
	Expect(string(data)).To(ContainSubstring("hoverfly.username: \"\""))
	Expect(string(data)).To(ContainSubstring("hoverfly.password: \"\""))
	Expect(string(data)).To(ContainSubstring("hoverfly.webserver: false"))
	Expect(string(data)).To(ContainSubstring("hoverfly.tls.certificate: \"\""))
	Expect(string(data)).To(ContainSubstring("hoverfly.tls.key: \"\""))
	Expect(string(data)).To(ContainSubstring("hoverfly.tls.disable: false"))
	Expect(string(data)).To(ContainSubstring("hoverfly.cache.disable: false"))
}

func Test_Config_GetTarget_ReturnsTargetIfAlreadyExists(t *testing.T) {
	RegisterTestingT(t)

	unit := &Config{
		Targets: map[string]Target{
			"default": Target{
				AdminPort: 1234,
			},
		},
	}

	Expect(unit.GetTarget("default").AdminPort).To(Equal(1234))
}

func Test_Config_GetTarget_GetsDefaultIfTargetNameIsEmpty(t *testing.T) {
	RegisterTestingT(t)

	unit := &Config{
		Targets: map[string]Target{
			"default": Target{
				AdminPort: 1234,
			},
		},
	}

	Expect(unit.GetTarget("").AdminPort).To(Equal(1234))
}

func Test_Config_GetTarget_ReturnsNilIfTargetDoesntExist(t *testing.T) {
	RegisterTestingT(t)

	unit := defaultConfig

	Expect(unit.GetTarget("default")).To(BeNil())
}

func Test_Config_NewTarget_AddsTarget(t *testing.T) {
	RegisterTestingT(t)

	unit := defaultConfig

	unit.NewTarget(Target{
		Name:      "default",
		AdminPort: 1234,
	})

	Expect(unit.Targets).To(HaveLen(1))

	Expect(unit.Targets["default"].AdminPort).To(Equal(1234))
}

func Test_Config_DeleteTarget_DeletesTarget(t *testing.T) {
	RegisterTestingT(t)

	unit := defaultConfig

	unit.NewTarget(Target{
		Name:      "default",
		AdminPort: 1234,
	})

	Expect(unit.Targets).To(HaveLen(1))
	unit.DeleteTarget(*unit.GetTarget("default"))
	Expect(unit.Targets).To(HaveLen(0))
}

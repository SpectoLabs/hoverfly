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

func Test_Config_WriteToFile_WritesTheConfigObjectToAFileInAYamlFormat(t *testing.T) {
	RegisterTestingT(t)

	config := Config{
		HoverflyHost:        "testhost",
		HoverflyAdminPort:   "1234",
		HoverflyProxyPort:   "4567",
		HoverflyDbType:      "boltdb",
		HoverflyUsername:    "username",
		HoverflyPassword:    "password",
		HoverflyWebserver:   true,
		HoverflyCertificate: "/home/benjih/certificate.pem",
		HoverflyKey:         "/home/benjih/key.pem",
		HoverflyDisableTls:  true,
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
}

func Test_Config_BuildFlags_SettingWebserverToTrueAddsTheFlag(t *testing.T) {
	RegisterTestingT(t)

	unit := Config{
		HoverflyWebserver: true,
	}

	Expect(unit.BuildFlags()).To(HaveLen(1))
	Expect(unit.BuildFlags()[0]).To(Equal("-webserver"))
}

func Test_Config_BuildFlags_SettingWebserverToFalseDoesNotAddTheFlag(t *testing.T) {
	RegisterTestingT(t)

	unit := Config{
		HoverflyWebserver: false,
	}

	Expect(unit.BuildFlags()).To(HaveLen(0))
}

func Test_Config_BuildFlags_AdminPortSetsTheApFlag(t *testing.T) {
	RegisterTestingT(t)

	unit := Config{
		HoverflyAdminPort: "1234",
	}

	Expect(unit.BuildFlags()).To(HaveLen(1))
	Expect(unit.BuildFlags()[0]).To(Equal("-ap=1234"))
}

func Test_Config_BuildFlags_ProxyPortSetsThePpFlag(t *testing.T) {
	RegisterTestingT(t)

	unit := Config{
		HoverflyProxyPort: "3421",
	}

	Expect(unit.BuildFlags()).To(HaveLen(1))
	Expect(unit.BuildFlags()[0]).To(Equal("-pp=3421"))
}

func Test_Config_BuildFlags_DbTypeSetsTheDbFlag(t *testing.T) {
	RegisterTestingT(t)

	unit := Config{
		HoverflyDbType: "boltdb",
	}

	Expect(unit.BuildFlags()).To(HaveLen(1))
	Expect(unit.BuildFlags()[0]).To(Equal("-db=boltdb"))
}

func Test_Config_BuildFlags_CertificateSetsCertFlag(t *testing.T) {
	RegisterTestingT(t)

	unit := Config{
		HoverflyCertificate: "certificate.pem",
	}

	Expect(unit.BuildFlags()).To(HaveLen(1))
	Expect(unit.BuildFlags()[0]).To(Equal("-cert=certificate.pem"))
}

func Test_Config_BuildFlags_KeySetsKeyFlag(t *testing.T) {
	RegisterTestingT(t)

	unit := Config{
		HoverflyKey: "key.pem",
	}

	Expect(unit.BuildFlags()).To(HaveLen(1))
	Expect(unit.BuildFlags()[0]).To(Equal("-key=key.pem"))
}

func Test_Config_BuildFlags_DisableTlsSetsTlsVerificationFlagToFalse(t *testing.T) {
	RegisterTestingT(t)

	unit := Config{
		HoverflyDisableTls: true,
	}

	Expect(unit.BuildFlags()).To(HaveLen(1))
	Expect(unit.BuildFlags()[0]).To(Equal("-tls-verification=false"))
}

func Test_Config_BuildFlags_UpstreamProxySetsUpstreamProxyFlag(t *testing.T) {
	RegisterTestingT(t)

	unit := Config{
		HoverflyUpstreamProxy: "hoverfly.io:8080",
	}

	Expect(unit.BuildFlags()).To(HaveLen(1))
	Expect(unit.BuildFlags()[0]).To(Equal("-upstream-proxy=hoverfly.io:8080"))
}

func Test_Config_BuildFlags_CanBuildFlagsInCorrectOrderWithAllVariables(t *testing.T) {
	RegisterTestingT(t)

	unit := Config{
		HoverflyWebserver:   true,
		HoverflyCertificate: "certificate.pem",
		HoverflyKey:         "key.pem",
		HoverflyDisableTls:  true,
	}

	Expect(unit.BuildFlags()).To(HaveLen(4))
	Expect(unit.BuildFlags()[0]).To(Equal("-webserver"))
	Expect(unit.BuildFlags()[1]).To(Equal("-cert=certificate.pem"))
	Expect(unit.BuildFlags()[2]).To(Equal("-key=key.pem"))
	Expect(unit.BuildFlags()[3]).To(Equal("-tls-verification=false"))
}

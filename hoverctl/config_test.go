package main

import (
	"io/ioutil"
	"os"
	"testing"

	. "github.com/onsi/gomega"
)

var (
	defaultHoverflyHost      = "localhost"
	defaultHoverflyAdminPort = "8888"
	defaultHoverflyProxyPort = "8500"
	defaultHoverflyUsername  = ""
	defaultHoverflyPassword  = ""
	defaultHoverflyWebserver = false
)

func Test_GetConfigWillReturnTheDefaultValues(t *testing.T) {
	RegisterTestingT(t)

	SetConfigurationDefaults()
	result := GetConfig()

	Expect(result.HoverflyHost).To(Equal(defaultHoverflyHost))
	Expect(result.HoverflyAdminPort).To(Equal(defaultHoverflyAdminPort))
	Expect(result.HoverflyProxyPort).To(Equal(defaultHoverflyProxyPort))
	Expect(result.HoverflyUsername).To(Equal(defaultHoverflyUsername))
	Expect(result.HoverflyPassword).To(Equal(defaultHoverflyPassword))
	Expect(result.HoverflyWebserver).To(Equal(defaultHoverflyWebserver))
}

func Test_Config_SetHost_OverridesDefaultValueWithAHoverflyHost(t *testing.T) {
	RegisterTestingT(t)

	SetConfigurationDefaults()
	result := GetConfig().SetHost("testhost")

	Expect(result.HoverflyHost).To(Equal("testhost"))
	Expect(result.HoverflyAdminPort).To(Equal(defaultHoverflyAdminPort))
	Expect(result.HoverflyProxyPort).To(Equal(defaultHoverflyProxyPort))
	Expect(result.HoverflyUsername).To(Equal(defaultHoverflyUsername))
	Expect(result.HoverflyPassword).To(Equal(defaultHoverflyPassword))
}

func Test_Config_SetHost_DoesNotOverrideWhenEmpty(t *testing.T) {
	RegisterTestingT(t)

	SetConfigurationDefaults()
	result := GetConfig().SetHost("")

	Expect(result.HoverflyHost).To(Equal(defaultHoverflyHost))
	Expect(result.HoverflyAdminPort).To(Equal(defaultHoverflyAdminPort))
	Expect(result.HoverflyProxyPort).To(Equal(defaultHoverflyProxyPort))
	Expect(result.HoverflyUsername).To(Equal(defaultHoverflyUsername))
	Expect(result.HoverflyPassword).To(Equal(defaultHoverflyPassword))
}

func Test_Config_SetAdminPort_OverridesDefaultValueWithAHoverflyAdminPort(t *testing.T) {
	RegisterTestingT(t)

	SetConfigurationDefaults()
	result := GetConfig().SetAdminPort("5")

	Expect(result.HoverflyHost).To(Equal(defaultHoverflyHost))
	Expect(result.HoverflyAdminPort).To(Equal("5"))
	Expect(result.HoverflyProxyPort).To(Equal(defaultHoverflyProxyPort))
	Expect(result.HoverflyUsername).To(Equal(defaultHoverflyUsername))
	Expect(result.HoverflyPassword).To(Equal(defaultHoverflyPassword))
}

func Test_Config_SetAdminPort_DoesNotOverrideWhenEmpty(t *testing.T) {
	RegisterTestingT(t)

	SetConfigurationDefaults()
	result := GetConfig().SetAdminPort("")

	Expect(result.HoverflyHost).To(Equal(defaultHoverflyHost))
	Expect(result.HoverflyAdminPort).To(Equal(defaultHoverflyAdminPort))
	Expect(result.HoverflyProxyPort).To(Equal(defaultHoverflyProxyPort))
	Expect(result.HoverflyUsername).To(Equal(defaultHoverflyUsername))
	Expect(result.HoverflyPassword).To(Equal(defaultHoverflyPassword))
}

func Test_Config_SetProxyPort_OverridesDefaultValueWithAHoverflyProxyPort(t *testing.T) {
	RegisterTestingT(t)

	SetConfigurationDefaults()
	result := GetConfig().SetProxyPort("7")

	Expect(result.HoverflyHost).To(Equal(defaultHoverflyHost))
	Expect(result.HoverflyAdminPort).To(Equal(defaultHoverflyAdminPort))
	Expect(result.HoverflyProxyPort).To(Equal("7"))
	Expect(result.HoverflyUsername).To(Equal(defaultHoverflyUsername))
	Expect(result.HoverflyPassword).To(Equal(defaultHoverflyPassword))
}

func Test_Config_SetProxyPort_DoesNotOverrideWhenEmpty(t *testing.T) {
	RegisterTestingT(t)

	SetConfigurationDefaults()
	result := GetConfig().SetProxyPort("")

	Expect(result.HoverflyHost).To(Equal(defaultHoverflyHost))
	Expect(result.HoverflyAdminPort).To(Equal(defaultHoverflyAdminPort))
	Expect(result.HoverflyProxyPort).To(Equal(defaultHoverflyProxyPort))
	Expect(result.HoverflyUsername).To(Equal(defaultHoverflyUsername))
	Expect(result.HoverflyPassword).To(Equal(defaultHoverflyPassword))
}

func Test_Config_SetUsername_OverridesDefaultValueWithAHoverflyUsername(t *testing.T) {
	RegisterTestingT(t)

	SetConfigurationDefaults()
	result := GetConfig().SetUsername("benjih")

	Expect(result.HoverflyHost).To(Equal(defaultHoverflyHost))
	Expect(result.HoverflyAdminPort).To(Equal(defaultHoverflyAdminPort))
	Expect(result.HoverflyProxyPort).To(Equal(defaultHoverflyProxyPort))
	Expect(result.HoverflyUsername).To(Equal("benjih"))
	Expect(result.HoverflyPassword).To(Equal(defaultHoverflyPassword))
}

func Test_Config_SetUsername_DoesNotOverrideWhenEmpty(t *testing.T) {
	RegisterTestingT(t)

	SetConfigurationDefaults()
	result := GetConfig().SetUsername("")

	Expect(result.HoverflyHost).To(Equal(defaultHoverflyHost))
	Expect(result.HoverflyAdminPort).To(Equal(defaultHoverflyAdminPort))
	Expect(result.HoverflyProxyPort).To(Equal(defaultHoverflyProxyPort))
	Expect(result.HoverflyUsername).To(Equal(defaultHoverflyUsername))
	Expect(result.HoverflyPassword).To(Equal(defaultHoverflyPassword))
}

func Test_Config_SetPassword_OverridesDefaultValueWithAHoverflyPassword(t *testing.T) {
	RegisterTestingT(t)

	SetConfigurationDefaults()
	result := GetConfig().SetPassword("burger-toucher")

	Expect(result.HoverflyHost).To(Equal(defaultHoverflyHost))
	Expect(result.HoverflyAdminPort).To(Equal(defaultHoverflyAdminPort))
	Expect(result.HoverflyProxyPort).To(Equal(defaultHoverflyProxyPort))
	Expect(result.HoverflyUsername).To(Equal(defaultHoverflyUsername))
	Expect(result.HoverflyPassword).To(Equal("burger-toucher"))
}

func Test_Config_SetPassword_DoesNotOverrideWhenEmpty(t *testing.T) {
	RegisterTestingT(t)

	SetConfigurationDefaults()
	result := GetConfig().SetPassword("")

	Expect(result.HoverflyHost).To(Equal(defaultHoverflyHost))
	Expect(result.HoverflyAdminPort).To(Equal(defaultHoverflyAdminPort))
	Expect(result.HoverflyProxyPort).To(Equal(defaultHoverflyProxyPort))
	Expect(result.HoverflyUsername).To(Equal(defaultHoverflyUsername))
	Expect(result.HoverflyPassword).To(Equal(defaultHoverflyPassword))
}

func Test_Config_SetCertificate_OverridesDefaultValueWithAHoverflyPassword(t *testing.T) {
	RegisterTestingT(t)

	SetConfigurationDefaults()
	result := GetConfig().SetCertificate("/home/benjih/test/certificate.pem")

	Expect(result.HoverflyHost).To(Equal(defaultHoverflyHost))
	Expect(result.HoverflyAdminPort).To(Equal(defaultHoverflyAdminPort))
	Expect(result.HoverflyProxyPort).To(Equal(defaultHoverflyProxyPort))
	Expect(result.HoverflyUsername).To(Equal(defaultHoverflyUsername))
	Expect(result.HoverflyPassword).To(Equal(defaultHoverflyPassword))
	Expect(result.HoverflyCertificate).To(Equal("/home/benjih/test/certificate.pem"))
}

func Test_Config_SetKey_DoesNotOverrideWhenEmpty(t *testing.T) {
	RegisterTestingT(t)

	SetConfigurationDefaults()
	result := GetConfig().SetKey("/home/benjih/test/key.pem")

	Expect(result.HoverflyHost).To(Equal(defaultHoverflyHost))
	Expect(result.HoverflyAdminPort).To(Equal(defaultHoverflyAdminPort))
	Expect(result.HoverflyProxyPort).To(Equal(defaultHoverflyProxyPort))
	Expect(result.HoverflyUsername).To(Equal(defaultHoverflyUsername))
	Expect(result.HoverflyPassword).To(Equal(defaultHoverflyPassword))
	Expect(result.HoverflyKey).To(Equal("/home/benjih/test/key.pem"))
}

func Test_Config_SetWebserver_OverridesDefaultValueWithAHoverflyPassword(t *testing.T) {
	RegisterTestingT(t)

	SetConfigurationDefaults()
	result := GetConfig().SetWebserver(true)

	Expect(result.HoverflyHost).To(Equal(defaultHoverflyHost))
	Expect(result.HoverflyAdminPort).To(Equal(defaultHoverflyAdminPort))
	Expect(result.HoverflyProxyPort).To(Equal(defaultHoverflyProxyPort))
	Expect(result.HoverflyUsername).To(Equal(defaultHoverflyUsername))
	Expect(result.HoverflyPassword).To(Equal(defaultHoverflyPassword))
	Expect(result.HoverflyWebserver).To(Equal(true))
}

func Test_Config_WriteToFile_WritesTheConfigObjectToAFileInAYamlFormat(t *testing.T) {
	RegisterTestingT(t)

	SetConfigurationDefaults()
	config := GetConfig()
	config = config.SetHost("testhost")
	config = config.SetAdminPort("1234")
	config = config.SetProxyPort("4567")
	config = config.SetUsername("username")
	config = config.SetPassword("password")
	config = config.SetWebserver(true)
	config = config.SetCertificate("/home/benjih/certificate.pem")
	config = config.SetKey("/home/benjih/key.pem")

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
	Expect(string(data)).To(ContainSubstring("hoverfly.username: username"))
	Expect(string(data)).To(ContainSubstring("hoverfly.password: password"))
	Expect(string(data)).To(ContainSubstring("hoverfly.webserver: true"))
	Expect(string(data)).To(ContainSubstring("hoverfly.tls.certificate: /home/benjih/certificate.pem"))
	Expect(string(data)).To(ContainSubstring("hoverfly.tls.key: /home/benjih/key.pem"))
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
	Expect(string(data)).To(ContainSubstring("hoverfly.username: \"\""))
	Expect(string(data)).To(ContainSubstring("hoverfly.password: \"\""))
	Expect(string(data)).To(ContainSubstring("hoverfly.webserver: false"))
	Expect(string(data)).To(ContainSubstring("hoverfly.tls.certificate: \"\""))
	Expect(string(data)).To(ContainSubstring("hoverfly.tls.key: \"\""))
}

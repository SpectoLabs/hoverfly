package main

import (
	"io/ioutil"
	"os"
	"testing"

	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"
)

func Test_GetConfigWillReturnTheDefaultValues(t *testing.T) {
	RegisterTestingT(t)

	SetConfigurationDefaults()
	result := GetConfig("", "", "", "", "")

	Expect(result.HoverflyHost).To(Equal("localhost"))
	Expect(result.HoverflyAdminPort).To(Equal("8888"))
	Expect(result.HoverflyProxyPort).To(Equal("8500"))
	Expect(result.HoverflyUsername).To(Equal(""))
	Expect(result.HoverflyPassword).To(Equal(""))
}

func Test_GetConfigOverridesDefaultValueWithAHoverflyHost(t *testing.T) {
	RegisterTestingT(t)

	SetConfigurationDefaults()
	hoverflyHost := "testhost"
	result := GetConfig(hoverflyHost, "", "", "", "")

	Expect(result.HoverflyHost).To(Equal(hoverflyHost))
	Expect(result.HoverflyAdminPort).To(Equal("8888"))
	Expect(result.HoverflyProxyPort).To(Equal("8500"))
	Expect(result.HoverflyUsername).To(Equal(""))
	Expect(result.HoverflyPassword).To(Equal(""))
}

func Test_GetConfigOverridesDefaultValueWithAHoverflyAdminPort(t *testing.T) {
	RegisterTestingT(t)

	SetConfigurationDefaults()
	hoverflyAdminPort := "5"
	result := GetConfig("", hoverflyAdminPort, "", "", "")

	Expect(result.HoverflyHost).To(Equal("localhost"))
	Expect(result.HoverflyAdminPort).To(Equal(hoverflyAdminPort))
	Expect(result.HoverflyProxyPort).To(Equal("8500"))
	Expect(result.HoverflyUsername).To(Equal(""))
	Expect(result.HoverflyPassword).To(Equal(""))
}

func Test_GetConfigOverridesDefaultValueWithAHoverflyProxyPort(t *testing.T) {
	RegisterTestingT(t)

	SetConfigurationDefaults()
	hoverflyProxyPort := "7"
	result := GetConfig("", "", hoverflyProxyPort, "", "")

	Expect(result.HoverflyHost).To(Equal("localhost"))
	Expect(result.HoverflyAdminPort).To(Equal("8888"))
	Expect(result.HoverflyProxyPort).To(Equal(hoverflyProxyPort))
	Expect(result.HoverflyUsername).To(Equal(""))
	Expect(result.HoverflyPassword).To(Equal(""))
}

func Test_GetConfigOverridesDefaultValueWithAHoverflyUsername(t *testing.T) {
	RegisterTestingT(t)

	SetConfigurationDefaults()
	hoverflyUsername := "benjih"
	result := GetConfig("", "", "", hoverflyUsername, "")

	Expect(result.HoverflyHost).To(Equal("localhost"))
	Expect(result.HoverflyAdminPort).To(Equal("8888"))
	Expect(result.HoverflyProxyPort).To(Equal("8500"))
	Expect(result.HoverflyUsername).To(Equal(hoverflyUsername))
	Expect(result.HoverflyPassword).To(Equal(""))
}

func Test_GetConfigOverridesDefaultValueWithAHoverflyPassword(t *testing.T) {
	RegisterTestingT(t)

	SetConfigurationDefaults()
	hoverflyPassword := "mypassword123"
	result := GetConfig("", "", "", "", hoverflyPassword)

	Expect(result.HoverflyHost).To(Equal("localhost"))
	Expect(result.HoverflyAdminPort).To(Equal("8888"))
	Expect(result.HoverflyProxyPort).To(Equal("8500"))
	Expect(result.HoverflyUsername).To(Equal(""))
	Expect(result.HoverflyPassword).To(Equal(hoverflyPassword))
}

func Test_GetConfigOverridesDefaultValueWithAllOverrides(t *testing.T) {
	RegisterTestingT(t)

	SetConfigurationDefaults()
	hoverflyHost := "specto.io"
	hoverflyAdminPort := "7654"
	hoverflyProxyPort := "1523"
	hoverflyUsername := "hfuser"
	hoverflyPassword := "hfpassword"
	result := GetConfig(hoverflyHost, hoverflyAdminPort, hoverflyProxyPort, hoverflyUsername, hoverflyPassword)

	Expect(result.HoverflyHost).To(Equal(hoverflyHost))
	Expect(result.HoverflyAdminPort).To(Equal(hoverflyAdminPort))
	Expect(result.HoverflyProxyPort).To(Equal(hoverflyProxyPort))
	Expect(result.HoverflyUsername).To(Equal(hoverflyUsername))
	Expect(result.HoverflyPassword).To(Equal(hoverflyPassword))
}

func Test_ConfigWriteToFileWritesTheConfigObjectToAFileInAYamlFormat(t *testing.T) {
	RegisterTestingT(t)

	SetConfigurationDefaults()
	config := GetConfig("testhost", "1234", "4567", "username", "password")

	wd, _ := os.Getwd()
	hoverflyDirectory := HoverflyDirectory{
		Path: wd,
	}

	err := config.WriteToFile(hoverflyDirectory)

	Expect(err).To(BeNil())

	data, _ := ioutil.ReadFile(hoverflyDirectory.Path + "/config.yaml")
	os.Remove(hoverflyDirectory.Path + "/config.yaml")

	var result Config
	yaml.Unmarshal(data, &result)

	Expect(result.HoverflyHost).To(Equal("testhost"))
	Expect(result.HoverflyAdminPort).To(Equal("1234"))
	Expect(result.HoverflyProxyPort).To(Equal("4567"))
	Expect(result.HoverflyUsername).To(Equal("username"))
	Expect(result.HoverflyPassword).To(Equal("password"))
}

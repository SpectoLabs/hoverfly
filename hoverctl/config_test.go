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
	result := GetConfig()

	Expect(result.HoverflyHost).To(Equal("localhost"))
	Expect(result.HoverflyAdminPort).To(Equal("8888"))
	Expect(result.HoverflyProxyPort).To(Equal("8500"))
	Expect(result.HoverflyUsername).To(Equal(""))
	Expect(result.HoverflyPassword).To(Equal(""))
}

func Test_Config_SetHost_OverridesDefaultValueWithAHoverflyHost(t *testing.T) {
	RegisterTestingT(t)

	SetConfigurationDefaults()
	hoverflyHost := "testhost"
	result := GetConfig().SetHost(hoverflyHost)

	Expect(result.HoverflyHost).To(Equal(hoverflyHost))
	Expect(result.HoverflyAdminPort).To(Equal("8888"))
	Expect(result.HoverflyProxyPort).To(Equal("8500"))
	Expect(result.HoverflyUsername).To(Equal(""))
	Expect(result.HoverflyPassword).To(Equal(""))
}

func Test_Config_SetAdminPort_OverridesDefaultValueWithAHoverflyAdminPort(t *testing.T) {
	RegisterTestingT(t)

	SetConfigurationDefaults()
	hoverflyAdminPort := "5"
	result := GetConfig().SetAdminPort(hoverflyAdminPort)

	Expect(result.HoverflyHost).To(Equal("localhost"))
	Expect(result.HoverflyAdminPort).To(Equal(hoverflyAdminPort))
	Expect(result.HoverflyProxyPort).To(Equal("8500"))
	Expect(result.HoverflyUsername).To(Equal(""))
	Expect(result.HoverflyPassword).To(Equal(""))
}

func Test_Config_SetProxyPort_OverridesDefaultValueWithAHoverflyProxyPort(t *testing.T) {
	RegisterTestingT(t)

	SetConfigurationDefaults()
	hoverflyProxyPort := "7"
	result := GetConfig().SetProxyPort(hoverflyProxyPort)

	Expect(result.HoverflyHost).To(Equal("localhost"))
	Expect(result.HoverflyAdminPort).To(Equal("8888"))
	Expect(result.HoverflyProxyPort).To(Equal(hoverflyProxyPort))
	Expect(result.HoverflyUsername).To(Equal(""))
	Expect(result.HoverflyPassword).To(Equal(""))
}

func Test_Config_SetUsername_OverridesDefaultValueWithAHoverflyUsername(t *testing.T) {
	RegisterTestingT(t)

	SetConfigurationDefaults()
	hoverflyUsername := "benjih"
	result := GetConfig().SetUsername(hoverflyUsername)

	Expect(result.HoverflyHost).To(Equal("localhost"))
	Expect(result.HoverflyAdminPort).To(Equal("8888"))
	Expect(result.HoverflyProxyPort).To(Equal("8500"))
	Expect(result.HoverflyUsername).To(Equal(hoverflyUsername))
	Expect(result.HoverflyPassword).To(Equal(""))
}

func Test_Config_SetPassword_OverridesDefaultValueWithAHoverflyPassword(t *testing.T) {
	RegisterTestingT(t)

	SetConfigurationDefaults()
	hoverflyPassword := "mypassword123"
	result := GetConfig().SetPassword(hoverflyPassword)

	Expect(result.HoverflyHost).To(Equal("localhost"))
	Expect(result.HoverflyAdminPort).To(Equal("8888"))
	Expect(result.HoverflyProxyPort).To(Equal("8500"))
	Expect(result.HoverflyUsername).To(Equal(""))
	Expect(result.HoverflyPassword).To(Equal(hoverflyPassword))
}

func Test_Config_WriteToFile_WritesTheConfigObjectToAFileInAYamlFormat(t *testing.T) {
	RegisterTestingT(t)

	SetConfigurationDefaults()
	config := GetConfig()
	config = config.SetHost("testhost").SetAdminPort("1234").SetProxyPort("4567").SetUsername("username").SetPassword("password")

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

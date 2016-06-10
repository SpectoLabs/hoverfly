package main

import (
	"testing"
	. "github.com/onsi/gomega"
	"os"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)


func Test_GetConfigWillReturnTheDefaultValues(t *testing.T) {
	RegisterTestingT(t)

	SetConfigurationDefaults()
	result := GetConfig("", "", "")

	Expect(result.HoverflyHost).To(Equal("localhost"))
	Expect(result.HoverflyAdminPort).To(Equal("8888"))
	Expect(result.HoverflyProxyPort).To(Equal("8500"))
	Expect(result.SpectoLabHost).To(Equal("localhost"))
	Expect(result.SpectoLabPort).To(Equal("81"))
	Expect(result.SpectoLabApiKey).To(Equal(""))
}

func Test_GetConfigOverridesDefaultValueWithAHoverflyHost(t *testing.T) {
	RegisterTestingT(t)

	SetConfigurationDefaults()
	hoverflyHost := "testhost"
	result := GetConfig(hoverflyHost, "", "")


	Expect(result.HoverflyHost).To(Equal(hoverflyHost))
	Expect(result.HoverflyAdminPort).To(Equal("8888"))
	Expect(result.HoverflyProxyPort).To(Equal("8500"))
	Expect(result.SpectoLabHost).To(Equal("localhost"))
	Expect(result.SpectoLabPort).To(Equal("81"))
	Expect(result.SpectoLabApiKey).To(Equal(""))
}

func Test_GetConfigOverridesDefaultValueWithAHoverflyAdminPort(t *testing.T) {
	RegisterTestingT(t)

	SetConfigurationDefaults()
	hoverflyAdminPort := "5"
	result := GetConfig("", hoverflyAdminPort, "")

	Expect(result.HoverflyHost).To(Equal("localhost"))
	Expect(result.HoverflyAdminPort).To(Equal(hoverflyAdminPort))
	Expect(result.HoverflyProxyPort).To(Equal("8500"))
	Expect(result.SpectoLabHost).To(Equal("localhost"))
	Expect(result.SpectoLabPort).To(Equal("81"))
	Expect(result.SpectoLabApiKey).To(Equal(""))
}

func Test_GetConfigOverridesDefaultValueWithAHoverflyProxyPort(t *testing.T) {
	RegisterTestingT(t)

	SetConfigurationDefaults()
	hoverflyProxyPort := "7"
	result := GetConfig("", "", hoverflyProxyPort)

	Expect(result.HoverflyHost).To(Equal("localhost"))
	Expect(result.HoverflyAdminPort).To(Equal("8888"))
	Expect(result.HoverflyProxyPort).To(Equal(hoverflyProxyPort))
	Expect(result.SpectoLabHost).To(Equal("localhost"))
	Expect(result.SpectoLabPort).To(Equal("81"))
	Expect(result.SpectoLabApiKey).To(Equal(""))
}

func Test_GetConfigOverridesDefaultValueWithAllOverrides(t *testing.T) {
	RegisterTestingT(t)

	SetConfigurationDefaults()
	hoverflyHost := "specto.io"
	hoverflyAdminPort := "7654"
	hoverflyProxyPort := "1523"
	result := GetConfig(hoverflyHost, hoverflyAdminPort, hoverflyProxyPort)

	Expect(result.HoverflyHost).To(Equal(hoverflyHost))
	Expect(result.HoverflyAdminPort).To(Equal(hoverflyAdminPort))
	Expect(result.HoverflyProxyPort).To(Equal(hoverflyProxyPort))
	Expect(result.SpectoLabHost).To(Equal("localhost"))
	Expect(result.SpectoLabPort).To(Equal("81"))
	Expect(result.SpectoLabApiKey).To(Equal(""))
}

func Test_ConfigWriteToFileWritesTheConfigObjectToAFileInAYamlFormat(t *testing.T) {
	RegisterTestingT(t)

	SetConfigurationDefaults()
	config := GetConfig("", "", "")

	workingDir, _ := os.Getwd()
	config.WriteToFile(workingDir)

	data, _ := ioutil.ReadFile(workingDir + "/config.yaml")
	os.Remove(workingDir + "/config.yaml")

	var result Config
	yaml.Unmarshal(data, &result)

	Expect(result.HoverflyHost).To(Equal("localhost"))
	Expect(result.HoverflyAdminPort).To(Equal("8888"))
	Expect(result.HoverflyProxyPort).To(Equal("8500"))
	Expect(result.SpectoLabHost).To(Equal("localhost"))
	Expect(result.SpectoLabPort).To(Equal("81"))
	Expect(result.SpectoLabApiKey).To(Equal(""))
}
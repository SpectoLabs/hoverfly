package configuration

import (
	"bytes"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"testing"

	. "github.com/onsi/gomega"
)

var (
	defaultConfig = Config{
		DefaultTarget: "local",
		Targets: map[string]Target{
			"local": {
				Name:      "local",
				Host:      "localhost",
				AdminPort: 8888,
				ProxyPort: 8500,
			},
		},
	}
)

func Test_GetConfigWillReturnTheDefaultValues(t *testing.T) {
	RegisterTestingT(t)

	SetConfigurationDefaults()
	result := GetConfig()

	Expect(*result).To(Equal(defaultConfig))
}

func Test_GetConfigWillInitializeMissingDefaultValues(t *testing.T) {
	RegisterTestingT(t)

	viper.SetConfigType("yaml")
	var configSource = []byte(`
default: local
targets:
 local:
   name: local
   authenabled: false
   username: ""
   password: ""
 remote:
   name: remote
`)

	_ = viper.ReadConfig(bytes.NewBuffer(configSource))

	result := parseConfig()

	Expect(*result).To(Equal(Config{
		DefaultTarget: "local",
		Targets: map[string]Target{
			"local": {
				Name:      "local",
				Host:      "localhost",
				AdminPort: 8888,
				ProxyPort: 8500,
			},

			"remote": {
				Name:      "remote",
				Host:      "localhost",
				AdminPort: 8888,
				ProxyPort: 8500,
			},
		},
	}))
}

func Test_GetConfigWillReadConfigFromAYamlFile(t *testing.T) {
	RegisterTestingT(t)

	viper.SetConfigType("yaml")
	var configSource = []byte(`
default: local
targets:
  local:
    name: local
    host: localhost
    admin.port: 8888
    proxy.port: 8500
    authenabled: false
    username: ""
    password: ""
  remote:
    name: remote
    host: hoverfly.cloud
    admin.port: 2345
    proxy.port: 9875
    authenabled: true
    username: "admin"
    password: "123"
`)

	_ = viper.ReadConfig(bytes.NewBuffer(configSource))

	result := parseConfig()

	Expect(*result).To(Equal(Config{
		DefaultTarget: "local",
		Targets: map[string]Target{
			"local": {
				Name:      "local",
				Host:      "localhost",
				AdminPort: 8888,
				ProxyPort: 8500,
			},

			"remote": {
				Name:      "remote",
				Host:      "hoverfly.cloud",
				AdminPort: 2345,
				ProxyPort: 9875,
				AuthEnabled: true,
				Username: 	"admin",
				Password: 	"123",
			},
		},
	}))
}

func Test_Config_WriteToFile_WritesTheConfigObjectToAFileInAYamlFormat(t *testing.T) {
	RegisterTestingT(t)

	config := Config{
		Targets: map[string]Target{
			"test-target": {
				Name:      "test-target",
				AdminPort: 1234,
				ProxyPort: 8765,
			},
		},
	}

	wd, _ := os.Getwd()
	hoverflyDirectory := HoverflyDirectory{
		Path: wd,
	}

	err := config.WriteToFile(hoverflyDirectory)

	Expect(err).To(BeNil())

	data, _ := ioutil.ReadFile(hoverflyDirectory.Path + "/config.yaml")
	_ = os.Remove(hoverflyDirectory.Path + "/config.yaml")

	Expect(string(data)).To(ContainSubstring(`targets:`))
	Expect(string(data)).To(ContainSubstring(`test-target:`))
	Expect(string(data)).To(ContainSubstring("name: test-target"))
	Expect(string(data)).To(ContainSubstring("admin.port: 1234"))
	Expect(string(data)).To(ContainSubstring("proxy.port: 8765"))
}

func Test_Config_GetTarget_ReturnsTargetIfAlreadyExists(t *testing.T) {
	RegisterTestingT(t)

	unit := &Config{
		Targets: map[string]Target{
			"default": {
				AdminPort: 1234,
			},
		},
	}

	Expect(unit.GetTarget("default").AdminPort).To(Equal(1234))
}

func Test_Config_GetTarget_GetsCurrentTargetIfTargetNameIsEmpty(t *testing.T) {
	RegisterTestingT(t)

	unit := &Config{
		DefaultTarget: "default",
		Targets: map[string]Target{
			"default": {
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

	unit := Config{
		Targets: map[string]Target{},
	}

	unit.NewTarget(Target{
		Name:      "newtarget",
		AdminPort: 1234,
	})

	Expect(unit.Targets).To(HaveLen(1))

	Expect(unit.Targets["newtarget"].AdminPort).To(Equal(1234))
}

func Test_Config_DeleteTarget_DeletesTarget(t *testing.T) {
	RegisterTestingT(t)

	unit := Config{
		Targets: map[string]Target{
			"deleteme": {
				Name:      "deleteme",
				AdminPort: 1234,
			},
		},
	}

	Expect(unit.Targets).To(HaveLen(1))
	unit.DeleteTarget(*unit.GetTarget("deleteme"))
	Expect(unit.Targets).To(HaveLen(0))
}

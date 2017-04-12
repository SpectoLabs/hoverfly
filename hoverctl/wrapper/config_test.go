package wrapper

import (
	"io/ioutil"
	"os"
	"testing"

	. "github.com/onsi/gomega"
)

var (
	defaultConfig = Config{
		DefaultTarget: "default",
		Targets:       map[string]Target{},
	}
)

func Test_GetConfigWillReturnTheDefaultValues(t *testing.T) {
	RegisterTestingT(t)

	SetConfigurationDefaults()
	result := GetConfig()

	Expect(*result).To(Equal(defaultConfig))
}

func Test_Config_WriteToFile_WritesTheConfigObjectToAFileInAYamlFormat(t *testing.T) {
	RegisterTestingT(t)

	config := Config{
		Targets: map[string]Target{
			"test-target": Target{
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
	os.Remove(hoverflyDirectory.Path + "/config.yaml")

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
			"default": Target{
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

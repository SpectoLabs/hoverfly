package main

import (
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"
	"testing"
)

var localCache_testDirectory = "/tmp/hoverctl-tests"

func localCache_setup() {
	os.Mkdir(localCache_testDirectory, 0777)
}

func localCache_teardown() {
	os.RemoveAll(localCache_testDirectory)
}

func Test_LocalCache_WriteSimulation(t *testing.T) {
	RegisterTestingT(t)
	localCache_setup()

	localCache := LocalCache{URI: localCache_testDirectory}
	simulation := Simulation{Vendor: "vendor", Name: "name", Version: "v1"}

	err := localCache.WriteSimulation(simulation, []byte("hello"))

	Expect(err).To(BeNil())

	data, err := ioutil.ReadFile(localCache_testDirectory + "/vendor.name.v1.json")

	Expect(err).To(BeNil())
	Expect(string(data)).To(Equal("hello"))

	localCache_teardown()
}

func Test_LocalCache_WriteSimulation_WithJson(t *testing.T) {
	RegisterTestingT(t)
	localCache_setup()

	localCache := LocalCache{URI: localCache_testDirectory}
	simulation := Simulation{Vendor: "vendor", Name: "test", Version: "v1"}

	err := localCache.WriteSimulation(simulation, []byte(`{"key":"value"}`))

	Expect(err).To(BeNil())

	data, err := ioutil.ReadFile(localCache_testDirectory + "/vendor.test.v1.json")

	Expect(err).To(BeNil())
	Expect(string(data)).To(Equal(`{"key":"value"}`))

	localCache_teardown()
}

func Test_LocalCache_ReadSimulation(t *testing.T) {
	RegisterTestingT(t)
	localCache_setup()

	ioutil.WriteFile(localCache_testDirectory+"/vendor.name.v1.json", []byte("this is a test file"), 0644)

	localCache := LocalCache{URI: localCache_testDirectory}
	simulation := Simulation{Vendor: "vendor", Name: "name", Version: "v1"}

	data, err := localCache.ReadSimulation(simulation)

	Expect(err).To(BeNil())
	Expect(data).To(Equal([]byte("this is a test file")))

	localCache_teardown()
}

func Test_LocalCache_ReadSimulation_ErrorsWhenFileIsMissing(t *testing.T) {
	RegisterTestingT(t)
	localCache_setup()

	localCache := LocalCache{URI: localCache_testDirectory}
	simulation := Simulation{Vendor: "vendor", Name: "name", Version: "v1"}

	data, err := localCache.ReadSimulation(simulation)

	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("Simulation not found in local cache"))
	Expect(data).To(BeNil())

	localCache_teardown()
}

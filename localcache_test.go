package main

import (
	"testing"
	"io/ioutil"
	. "github.com/onsi/gomega"
	"os"
)

var testDirectory = "/tmp/hoverctl-tests"

func setup() {
	os.Mkdir(testDirectory, 0777)
}

func teardown() {
	os.RemoveAll(testDirectory)
}

func Test_LocalCache_WriteSimulation(t *testing.T) {
	RegisterTestingT(t)
	setup()

	localCache := LocalCache{Uri: testDirectory}
	simulation := Simulation{Vendor: "vendor", Name: "name", Version: "v1"}

	err := localCache.WriteSimulation(simulation, []byte("hello"))

	Expect(err).To(BeNil())

	data, err := ioutil.ReadFile(testDirectory + "/vendor.name.v1.hfile")

	Expect(err).To(BeNil())
	Expect(string(data)).To(Equal("hello"))

	teardown()
}

func Test_LocalCache_WriteSimulation_WithJson(t *testing.T) {
	RegisterTestingT(t)
	setup()

	localCache := LocalCache{Uri: testDirectory}
	simulation := Simulation{Vendor: "vendor", Name: "test", Version: "v1"}

	err := localCache.WriteSimulation(simulation, []byte(`{"key":"value"}`))

	Expect(err).To(BeNil())

	data, err := ioutil.ReadFile(testDirectory + "/vendor.test.v1.hfile")

	Expect(err).To(BeNil())
	Expect(string(data)).To(Equal(`{"key":"value"}`))

	teardown()
}

func Test_LocalCache_ReadSimulation(t *testing.T) {
	RegisterTestingT(t)
	setup()

	ioutil.WriteFile(testDirectory + "/vendor.name.v1.hfile", []byte("this is a test file"), 0644)

	localCache := LocalCache{Uri: testDirectory}
	simulation := Simulation{Vendor: "vendor", Name: "name", Version: "v1"}

	data, err := localCache.ReadSimulation(simulation)

	Expect(err).To(BeNil())
	Expect(data).To(Equal([]byte("this is a test file")))

	teardown()
}

func Test_LocalCache_ReadSimulation_ErrorsWhenFileIsMissing(t *testing.T) {
	RegisterTestingT(t)
	setup()

	localCache := LocalCache{Uri: testDirectory}
	simulation := Simulation{Vendor: "vendor", Name: "name", Version: "v1"}

	data, err := localCache.ReadSimulation(simulation)

	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("Simulation not found"))
	Expect(data).To(BeNil())

	teardown()
}
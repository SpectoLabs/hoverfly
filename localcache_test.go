package main

import (
	"testing"
	"io/ioutil"
	. "github.com/onsi/gomega"
	"os"
)

var testDirectory = "/tmp/hoverfly-cli-tests"

func setup() {
	os.Mkdir(testDirectory, 0777)
}

func teardown() {
	os.RemoveAll(testDirectory)
}

func Test_LocalCache_ReadSimulation(t *testing.T) {
	RegisterTestingT(t)
	setup()

	ioutil.WriteFile(testDirectory + "/vendor.name.v1.hfile", []byte("this is a test file"), 0644)

	localCache := LocalCache{Uri: testDirectory}

	data, err := localCache.ReadSimulation(Hoverfile{Vendor: "vendor", Name: "name", Version: "v1"})

	Expect(err).To(BeNil())
	Expect(data).To(Equal([]byte("this is a test file")))

	teardown()
}

func Test_LocalCache_ReadSimulation_ErrorsWhenFileIsMissing(t *testing.T) {
	RegisterTestingT(t)
	setup()

	localCache := LocalCache{Uri: testDirectory}

	data, err := localCache.ReadSimulation(Hoverfile{Vendor: "vendor", Name: "name", Version: "v1"})

	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("Simulation not found"))
	Expect(data).To(BeNil())

	teardown()
}
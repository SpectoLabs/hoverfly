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

func Test_LocalCache_WriteSimulation(t *testing.T) {
	RegisterTestingT(t)
	setup()


	ioutil.WriteFile("/tmp/vendor.name.v1.hfile", []byte("this is a test file"), 0644)

	localCache := LocalCache{Uri: "/tmp"}

	data, err := localCache.ReadSimulation(Hoverfile{Vendor: "vendor", Name: "name", Version: "v1"})

	Expect(err).To(BeNil())
	Expect(data).To(Equal([]byte("this is a test file")))

	teardown()
}
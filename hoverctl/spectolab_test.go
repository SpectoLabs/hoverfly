package main

import (
	"testing"
	. "github.com/onsi/gomega"
	"os"
)

func TestMain(m *testing.M) {
	returnCode := m.Run()
	os.Exit(returnCode)
}

func Test_SpectoLab_buildUrl_JoinsAPathToTheHost(t *testing.T) {
	RegisterTestingT(t)

	spectoLab := SpectoLab{Host: "test-host.com"}

	result := spectoLab.buildURL("/cats")
	Expect(result).To(Equal("test-host.com/cats"))
}

func Test_SpectoLab_buildAuthorizationHeaderValue_UsesApiKey(t *testing.T) {
	RegisterTestingT(t)

	spectoLab := SpectoLab{APIKey: "test-key"}

	result := spectoLab.buildAuthorizationHeaderValue()
	Expect(result).To(Equal("Bearer test-key"))
}
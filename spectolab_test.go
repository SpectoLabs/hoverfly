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

func Test_SpectoLab_buildBaseUrl_UsesHostAndPort(t *testing.T) {
	RegisterTestingT(t)

	spectoLab := SpectoLab{Host: "test-host", Port: "12432"}

	result := spectoLab.buildBaseURL()
	Expect(result).To(Equal("test-host:12432"))
}

func Test_SpectoLab_buildBaseUrl_JustHostDoesNotIncludeSemicolon(t *testing.T) {
	RegisterTestingT(t)

	spectoLab := SpectoLab{Host: "test-host.com", Port: ""}

	result := spectoLab.buildBaseURL()
	Expect(result).To(Equal("test-host.com"))
}

func Test_SpectoLab_buildAuthorizationHeaderValue_UsesApiKey(t *testing.T) {
	RegisterTestingT(t)

	spectoLab := SpectoLab{APIKey: "test-key"}

	result := spectoLab.buildAuthorizationHeaderValue()
	Expect(result).To(Equal("Bearer test-key"))
}
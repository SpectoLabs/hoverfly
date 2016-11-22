package main

import (
	"testing"

	. "github.com/onsi/gomega"
)

func Test_FlagsBuilder_BuildFlags_ReturnsAStringArray(t *testing.T) {
	RegisterTestingT(t)

	unit := FlagsBuilder{}

	Expect(unit.BuildFlags()).To(BeAssignableToTypeOf([]string{}))
	Expect(unit.BuildFlags()).To(HaveLen(0))
}

func Test_FlagsBuilder_BuildFlags_SettingWebserverToWebserverPutsTheCorrectFlag(t *testing.T) {
	RegisterTestingT(t)

	unit := FlagsBuilder{
		Webserver: "webserver",
	}

	Expect(unit.BuildFlags()).To(HaveLen(1))
	Expect(unit.BuildFlags()[0]).To(Equal("-webserver"))
}

func Test_FlagsBuilder_BuildFlags_CertificateSetsCertFlag(t *testing.T) {
	RegisterTestingT(t)

	unit := FlagsBuilder{
		Certificate: "certificate.pem",
	}

	Expect(unit.BuildFlags()).To(HaveLen(2))
	Expect(unit.BuildFlags()[0]).To(Equal("-cert"))
	Expect(unit.BuildFlags()[1]).To(Equal("certificate.pem"))
}

func Test_FlagsBuilder_BuildFlags_KeySetsKeyFlag(t *testing.T) {
	RegisterTestingT(t)

	unit := FlagsBuilder{
		Key: "key.pem",
	}

	Expect(unit.BuildFlags()).To(HaveLen(2))
	Expect(unit.BuildFlags()[0]).To(Equal("-key"))
	Expect(unit.BuildFlags()[1]).To(Equal("key.pem"))
}

func Test_FlagsBuilder_BuildFlags_CanBuildFlagsInCorrectOrderWithAllVariables(t *testing.T) {
	RegisterTestingT(t)

	unit := FlagsBuilder{
		Webserver:   "webserver",
		Certificate: "certificate.pem",
		Key:         "key.pem",
	}

	Expect(unit.BuildFlags()).To(HaveLen(5))
	Expect(unit.BuildFlags()[0]).To(Equal("-webserver"))
	Expect(unit.BuildFlags()[1]).To(Equal("-cert"))
	Expect(unit.BuildFlags()[2]).To(Equal("certificate.pem"))
	Expect(unit.BuildFlags()[3]).To(Equal("-key"))
	Expect(unit.BuildFlags()[4]).To(Equal("key.pem"))
}

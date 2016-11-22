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

	Expect(unit.BuildFlags()).To(HaveLen(1))
	Expect(unit.BuildFlags()[0]).To(Equal("-cert=certificate.pem"))
}

func Test_FlagsBuilder_BuildFlags_KeySetsKeyFlag(t *testing.T) {
	RegisterTestingT(t)

	unit := FlagsBuilder{
		Key: "key.pem",
	}

	Expect(unit.BuildFlags()).To(HaveLen(1))
	Expect(unit.BuildFlags()[0]).To(Equal("-key=key.pem"))
}

func Test_FlagsBuilder_BuildFlags_DisableTlsSetsTlsVerificationFlagToFalse(t *testing.T) {
	RegisterTestingT(t)

	unit := FlagsBuilder{
		DisableTls: true,
	}

	Expect(unit.BuildFlags()).To(HaveLen(1))
	Expect(unit.BuildFlags()[0]).To(Equal("-tls-verification=false"))
}

func Test_FlagsBuilder_BuildFlags_CanBuildFlagsInCorrectOrderWithAllVariables(t *testing.T) {
	RegisterTestingT(t)

	unit := FlagsBuilder{
		Webserver:   "webserver",
		Certificate: "certificate.pem",
		Key:         "key.pem",
		DisableTls:  true,
	}

	Expect(unit.BuildFlags()).To(HaveLen(4))
	Expect(unit.BuildFlags()[0]).To(Equal("-webserver"))
	Expect(unit.BuildFlags()[1]).To(Equal("-cert=certificate.pem"))
	Expect(unit.BuildFlags()[2]).To(Equal("-key=key.pem"))
	Expect(unit.BuildFlags()[3]).To(Equal("-tls-verification=false"))
}

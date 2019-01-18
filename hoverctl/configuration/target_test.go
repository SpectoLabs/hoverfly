package configuration

import (
	"testing"

	. "github.com/onsi/gomega"
)

func Test_NewTarget_ReturnsDefaultWithEmptyStrings(t *testing.T) {
	RegisterTestingT(t)

	Expect(NewTarget("", "", 0, 0)).To(Equal(&Target{
		Name:      "local",
		Host:      "localhost",
		AdminPort: 8888,
		ProxyPort: 8500,
	}))
}

func Test_NewTarget_OverridesNamefNotEmpty(t *testing.T) {
	RegisterTestingT(t)

	Expect(NewTarget("notdefault", "", 0, 0)).To(Equal(&Target{
		Name:      "notdefault",
		Host:      "localhost",
		AdminPort: 8888,
		ProxyPort: 8500,
	}))
}

func Test_NewTarget_OverridesHostfNotEmpty(t *testing.T) {
	RegisterTestingT(t)

	Expect(NewTarget("", "notlocalhost", 0, 0)).To(Equal(&Target{
		Name:      "local",
		Host:      "notlocalhost",
		AdminPort: 8888,
		ProxyPort: 8500,
	}))
}

func Test_NewTarget_OverridesAdminPortfNotEmpty(t *testing.T) {
	RegisterTestingT(t)

	Expect(NewTarget("", "", 1234, 0)).To(Equal(&Target{
		Name:      "local",
		Host:      "localhost",
		AdminPort: 1234,
		ProxyPort: 8500,
	}))
}

func Test_NewTarget_OverridesProxyPortfNotEmpty(t *testing.T) {
	RegisterTestingT(t)

	Expect(NewTarget("", "", 0, 8765)).To(Equal(&Target{
		Name:      "local",
		Host:      "localhost",
		AdminPort: 8888,
		ProxyPort: 8765,
	}))
}

func Test_Target_BuildFlags_AdminPortSetsTheApFlag(t *testing.T) {
	RegisterTestingT(t)

	unit := Target{
		AdminPort: 1234,
	}

	Expect(unit.BuildFlags()).To(HaveLen(1))
	Expect(unit.BuildFlags()[0]).To(Equal("-ap=1234"))
}

func Test_Target_BuildFlags_ProxyPortSetsThePpFlag(t *testing.T) {
	RegisterTestingT(t)

	unit := Target{
		ProxyPort: 3421,
	}

	Expect(unit.BuildFlags()).To(HaveLen(1))
	Expect(unit.BuildFlags()[0]).To(Equal("-pp=3421"))
}

func Test_Target_BuildFlags_SettingWebserverToTrueAddsTheFlag(t *testing.T) {
	RegisterTestingT(t)

	unit := Target{
		Webserver: true,
	}

	Expect(unit.BuildFlags()).To(HaveLen(1))
	Expect(unit.BuildFlags()[0]).To(Equal("-webserver"))
}

func Test_Target_BuildFlags_SettingWebserverToFalseDoesNotAddTheFlag(t *testing.T) {
	RegisterTestingT(t)

	unit := Target{
		Webserver: false,
	}

	Expect(unit.BuildFlags()).To(HaveLen(0))
}

func Test_Target_BuildFlags_DbTypeSetsTheDbFlag(t *testing.T) {
	RegisterTestingT(t)

	unit := Target{
		CachePath: "cache.db",
	}

	Expect(unit.BuildFlags()).To(HaveLen(2))
	Expect(unit.BuildFlags()[0]).To(Equal("-db=boltdb"))
	Expect(unit.BuildFlags()[1]).To(Equal("-db-path=cache.db"))
}

func Test_Target_BuildFlags_ListenOnSetsListenOnHostFlag(t *testing.T) {
	RegisterTestingT(t)

	unit := Target{
		ListenOnHost: "0.0.0.0",
	}

	Expect(unit.BuildFlags()).To(HaveLen(1))
	Expect(unit.BuildFlags()[0]).To(Equal("-listen-on-host=0.0.0.0"))
}

func Test_Target_BuildFlags_CertificateSetsCertFlag(t *testing.T) {
	RegisterTestingT(t)

	unit := Target{
		CertificatePath: "certificate.pem",
	}

	Expect(unit.BuildFlags()).To(HaveLen(1))
	Expect(unit.BuildFlags()[0]).To(Equal("-cert=certificate.pem"))
}

func Test_Target_BuildFlags_KeySetsKeyFlag(t *testing.T) {
	RegisterTestingT(t)

	unit := Target{
		KeyPath: "key.pem",
	}

	Expect(unit.BuildFlags()).To(HaveLen(1))
	Expect(unit.BuildFlags()[0]).To(Equal("-key=key.pem"))
}

func Test_Target_BuildFlags_DisableTlsSetsTlsVerificationFlagToFalse(t *testing.T) {
	RegisterTestingT(t)

	unit := Target{
		DisableTls: true,
	}

	Expect(unit.BuildFlags()).To(HaveLen(1))
	Expect(unit.BuildFlags()[0]).To(Equal("-tls-verification=false"))
}

func Test_Target_BuildFlags_UpstreamProxySetsUpstreamProxyFlag(t *testing.T) {
	RegisterTestingT(t)

	unit := Target{
		UpstreamProxyUrl: "hoverfly.io:8080",
	}

	Expect(unit.BuildFlags()).To(HaveLen(1))
	Expect(unit.BuildFlags()[0]).To(Equal("-upstream-proxy=hoverfly.io:8080"))
}

func Test_Target_BuildFlags_HttpsOnlySetsTheFlag(t *testing.T) {
	RegisterTestingT(t)

	unit := Target{
		HttpsOnly: true,
	}

	Expect(unit.BuildFlags()).To(HaveLen(1))
	Expect(unit.BuildFlags()[0]).To(Equal("-https-only"))
}

func Test_Target_BuildFlags_IfAuthEnabledThenIncludesUsernameAndPassword(t *testing.T) {
	RegisterTestingT(t)

	unit := Target{
		AuthEnabled: true,
		Username:    "benji",
		Password:    "password",
	}

	Expect(unit.BuildFlags()).To(HaveLen(5))
	Expect(unit.BuildFlags()[0]).To(Equal("-auth"))
	Expect(unit.BuildFlags()[1]).To(Equal("-username"))
	Expect(unit.BuildFlags()[2]).To(Equal("benji"))
	Expect(unit.BuildFlags()[3]).To(Equal("-password-hash"))
	Expect(unit.BuildFlags()[4]).To(ContainSubstring("$2a$10$"))
}

func Test_Target_BuildFlags_CacheDisableBuildsCorrectFlagWhenTrue(t *testing.T) {
	RegisterTestingT(t)

	unit := Target{
		DisableCache: true,
	}

	Expect(unit.BuildFlags()).To(HaveLen(1))
	Expect(unit.BuildFlags()[0]).To(Equal("-disable-cache"))
}

func Test_Target_BuildFlags_CacheDisableDoesNotBuildCorrectFlagWhenFalse(t *testing.T) {
	RegisterTestingT(t)

	unit := Target{
		DisableCache: false,
	}

	Expect(unit.BuildFlags()).To(HaveLen(0))
}

func Test_Target_BuildFlags_CanBuildFlagsInCorrectOrderWithAllVariables(t *testing.T) {
	RegisterTestingT(t)

	unit := Target{
		Webserver:       true,
		CertificatePath: "certificate.pem",
		KeyPath:         "key.pem",
		DisableTls:      true,
	}

	Expect(unit.BuildFlags()).To(HaveLen(4))
	Expect(unit.BuildFlags()[0]).To(Equal("-webserver"))
	Expect(unit.BuildFlags()[1]).To(Equal("-cert=certificate.pem"))
	Expect(unit.BuildFlags()[2]).To(Equal("-key=key.pem"))
	Expect(unit.BuildFlags()[3]).To(Equal("-tls-verification=false"))
}

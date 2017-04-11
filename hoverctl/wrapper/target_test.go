package wrapper

import (
	"testing"

	. "github.com/onsi/gomega"
)

func Test_NewTarget_ReturnsDefaultWithEmptyStrings(t *testing.T) {
	RegisterTestingT(t)

	Expect(NewTarget("", "", 0, 0)).To(Equal(&Target{
		Name:      "default",
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
		Name:      "default",
		Host:      "notlocalhost",
		AdminPort: 8888,
		ProxyPort: 8500,
	}))
}

func Test_NewTarget_OverridesAdminPortfNotEmpty(t *testing.T) {
	RegisterTestingT(t)

	Expect(NewTarget("", "", 1234, 0)).To(Equal(&Target{
		Name:      "default",
		Host:      "localhost",
		AdminPort: 1234,
		ProxyPort: 8500,
	}))
}

func Test_NewTarget_OverridesProxyPortfNotEmpty(t *testing.T) {
	RegisterTestingT(t)

	Expect(NewTarget("", "", 0, 8765)).To(Equal(&Target{
		Name:      "default",
		Host:      "localhost",
		AdminPort: 8888,
		ProxyPort: 8765,
	}))
}

func Test_getTargetsFromConfig_host(t *testing.T) {
	RegisterTestingT(t)

	targets := getTargetsFromConfig(map[string]interface{}{
		"default": map[interface{}]interface{}{
			"host": "test.org",
		},
	})

	Expect(targets).To(HaveLen(1))
	Expect(targets).To(HaveKeyWithValue("default", Target{
		Name: "default",
		Host: "test.org",
	}))
}

func Test_getTargetsFromConfig_adminport(t *testing.T) {
	RegisterTestingT(t)

	targets := getTargetsFromConfig(map[string]interface{}{
		"other": map[interface{}]interface{}{
			"admin.port": 1234,
		},
	})

	Expect(targets).To(HaveLen(1))
	Expect(targets).To(HaveKeyWithValue("other", Target{
		Name:      "other",
		AdminPort: 1234,
	}))
}

func Test_getTargetsFromConfig_proxyport(t *testing.T) {
	RegisterTestingT(t)

	targets := getTargetsFromConfig(map[string]interface{}{
		"otherother": map[interface{}]interface{}{
			"proxy.port": 8765,
		},
	})

	Expect(targets).To(HaveLen(1))
	Expect(targets).To(HaveKeyWithValue("otherother", Target{
		Name:      "otherother",
		ProxyPort: 8765,
	}))
}

func Test_getTargetsFromConfig_authtoken(t *testing.T) {
	RegisterTestingT(t)

	targets := getTargetsFromConfig(map[string]interface{}{
		"anotherother": map[interface{}]interface{}{
			"auth.token": "token123:456",
		},
	})

	Expect(targets).To(HaveLen(1))
	Expect(targets).To(HaveKeyWithValue("anotherother", Target{
		Name:      "anotherother",
		AuthToken: "token123:456",
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

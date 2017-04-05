package wrapper

import (
	"testing"

	. "github.com/onsi/gomega"
)

func Test_NewTarget_ReturnsDefaultWithEmptyStrings(t *testing.T) {
	RegisterTestingT(t)

	Expect(NewTarget("", "", "", "")).To(Equal(&Target{
		Name:      "default",
		Host:      "localhost",
		AdminPort: 8888,
		ProxyPort: 8500,
	}))
}

func Test_NewTarget_OverridesNamefNotEmpty(t *testing.T) {
	RegisterTestingT(t)

	Expect(NewTarget("notdefault", "", "", "")).To(Equal(&Target{
		Name:      "notdefault",
		Host:      "localhost",
		AdminPort: 8888,
		ProxyPort: 8500,
	}))
}

func Test_NewTarget_OverridesHostfNotEmpty(t *testing.T) {
	RegisterTestingT(t)

	Expect(NewTarget("", "notlocalhost", "", "")).To(Equal(&Target{
		Name:      "default",
		Host:      "notlocalhost",
		AdminPort: 8888,
		ProxyPort: 8500,
	}))
}

func Test_NewTarget_OverridesAdminPortfNotEmpty(t *testing.T) {
	RegisterTestingT(t)

	Expect(NewTarget("", "", "1234", "")).To(Equal(&Target{
		Name:      "default",
		Host:      "localhost",
		AdminPort: 1234,
		ProxyPort: 8500,
	}))
}

func Test_NewTarget_ErrorsIfAdminPortIsNotANumber(t *testing.T) {
	RegisterTestingT(t)

	_, err := NewTarget("", "", "notanumber", "")

	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("The admin port provided was not a number"))
}

func Test_NewTarget_OverridesProxyPortfNotEmpty(t *testing.T) {
	RegisterTestingT(t)

	Expect(NewTarget("", "", "", "8765")).To(Equal(&Target{
		Name:      "default",
		Host:      "localhost",
		AdminPort: 8888,
		ProxyPort: 8765,
	}))
}

func Test_NewTarget_ErrorsIfProxyPortIsNotANumber(t *testing.T) {
	RegisterTestingT(t)

	_, err := NewTarget("", "", "", "notanumber")

	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("The proxy port provided was not a number"))
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

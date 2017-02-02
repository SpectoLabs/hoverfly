package wrapper

import (
	"testing"

	. "github.com/onsi/gomega"
)

func Test_Hoverfly_isLocal_WhenLocalhost(t *testing.T) {
	RegisterTestingT(t)

	hoverfly := Hoverfly{Host: "localhost"}

	result := hoverfly.isLocal()

	Expect(result).To(BeTrue())
}

func Test_Hoverfly_isLocal_WhenLocalhostIP(t *testing.T) {
	RegisterTestingT(t)

	hoverfly := Hoverfly{Host: "127.0.0.1"}

	result := hoverfly.isLocal()

	Expect(result).To(BeTrue())
}

func Test_Hoverfly_isLocal_WhenAnotherDNS(t *testing.T) {
	RegisterTestingT(t)

	hoverfly := Hoverfly{Host: "specto.io"}

	result := hoverfly.isLocal()

	Expect(result).To(BeFalse())
}

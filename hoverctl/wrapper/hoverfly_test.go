package wrapper

import (
	"testing"

	. "github.com/onsi/gomega"
)

func Test_isLocal_WhenLocalhost(t *testing.T) {
	RegisterTestingT(t)

	Expect(isLocal("localhost")).To(BeTrue())
}

func Test_isLocal_WhenLocalhostIP(t *testing.T) {
	RegisterTestingT(t)

	Expect(isLocal("127.0.0.1")).To(BeTrue())
}

func Test_isLocal_WhenAnotherDNS(t *testing.T) {
	RegisterTestingT(t)

	Expect(isLocal("specto.io")).To(BeFalse())
}

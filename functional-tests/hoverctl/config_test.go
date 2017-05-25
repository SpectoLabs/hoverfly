package hoverctl_suite

import (
	"github.com/SpectoLabs/hoverfly/functional-tests"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("When I use hoverctl", func() {

	BeforeEach(func() {
		functional_tests.Run(hoverctlBinary, "targets", "update", "local")
	})

	AfterEach(func() {
		WipeConfig()
	})

	Context("I can get the configuration being used by hoverctl", func() {

		It("prints the location of the config.yaml", func() {
			output := functional_tests.Run(hoverctlBinary, "config")

			Expect(output).To(ContainSubstring("functional-tests/hoverctl/.hoverfly/config.yaml"))
		})

		It("prints the contents of the config.yaml", func() {
			functional_tests.Run(hoverctlBinary, "targets", "create", "config-test")

			output := functional_tests.Run(hoverctlBinary, "config")

			Expect(output).To(ContainSubstring(`targets:`))
			Expect(output).To(ContainSubstring(`config-test:`))
			Expect(output).To(ContainSubstring("name: config-test"))
			Expect(output).To(ContainSubstring("host: localhost"))
			Expect(output).To(ContainSubstring("admin.port: 8888"))
			Expect(output).To(ContainSubstring("proxy.port: 8500"))
		})
	})

	Context("Can get the target host", func() {

		It("returns the host config of a target", func() {
			functional_tests.Run(hoverctlBinary, "targets", "create", "test-host", "--host", "test.host.com")
			output := functional_tests.Run(hoverctlBinary, "-t", "test-host", "config", "host")

			Expect(output).To(ContainSubstring("test.host.com"))
		})

		It("returns the host config of the default target", func() {
			functional_tests.Run(hoverctlBinary, "targets", "create", "test-host-default", "--host", "test.host.default.com")
			functional_tests.Run(hoverctlBinary, "targets", "default", "test-host-default")
			output := functional_tests.Run(hoverctlBinary, "config", "host")

			Expect(output).To(ContainSubstring("test.host.default.com"))
		})
	})

	Context("Can get the target admin-port", func() {

		It("returns the admin-port config of a target", func() {
			functional_tests.Run(hoverctlBinary, "targets", "create", "test-admin-port", "--admin-port", "2345")
			output := functional_tests.Run(hoverctlBinary, "-t", "test-admin-port", "config", "admin-port")

			Expect(output).To(ContainSubstring("2345"))
		})

		It("returns the admin-port config of the default target", func() {
			functional_tests.Run(hoverctlBinary, "targets", "create", "test-admin-port-default", "--admin-port", "1234")
			functional_tests.Run(hoverctlBinary, "targets", "default", "test-admin-port-default")
			output := functional_tests.Run(hoverctlBinary, "config", "admin-port")

			Expect(output).To(ContainSubstring("1234"))
		})
	})

	Context("Can get the target proxy-port", func() {

		It("returns the proxy-port config of a target", func() {
			functional_tests.Run(hoverctlBinary, "targets", "create", "test-proxy-port", "--proxy-port", "2345")
			output := functional_tests.Run(hoverctlBinary, "-t", "test-proxy-port", "config", "proxy-port")

			Expect(output).To(ContainSubstring("2345"))
		})

		It("returns the proxy-port config of the default target", func() {
			functional_tests.Run(hoverctlBinary, "targets", "create", "test-proxy-port-default", "--proxy-port", "1234")
			functional_tests.Run(hoverctlBinary, "targets", "default", "test-proxy-port-default")
			output := functional_tests.Run(hoverctlBinary, "config", "proxy-port")

			Expect(output).To(ContainSubstring("1234"))
		})
	})

	Context("Can get the target auth-token", func() {

		var hoverfly *functional_tests.Hoverfly

		BeforeEach(func() {
			hoverfly = functional_tests.NewHoverfly()
			hoverfly.Start()
		})

		AfterEach(func() {
			hoverfly.Stop()
		})

		It("returns the auth-token config of a target", func() {
			functional_tests.Run(hoverctlBinary, "targets", "create", "test-auth-token", "--admin-port", hoverfly.GetAdminPort())

			functional_tests.Run(hoverctlBinary, "-t", "test-auth-token", "login", "--username", functional_tests.HoverflyUsername, "--password", functional_tests.HoverflyPassword)
			output := functional_tests.Run(hoverctlBinary, "-t", "test-auth-token", "config", "auth-token")

			Expect(output).ToNot(ContainSubstring("No auth token"))
		})

		It("returns the auth-token config of the default target", func() {
			functional_tests.Run(hoverctlBinary, "targets", "create", "test-auth-token-default", "--admin-port", hoverfly.GetAdminPort())
			functional_tests.Run(hoverctlBinary, "targets", "default", "test-auth-token-default")

			functional_tests.Run(hoverctlBinary, "login", "--username", functional_tests.HoverflyUsername, "--password", functional_tests.HoverflyPassword)
			output := functional_tests.Run(hoverctlBinary, "config", "auth-token")

			Expect(output).ToNot(ContainSubstring("No auth token"))
		})

		It("returns an error if there is no auth-token", func() {
			functional_tests.Run(hoverctlBinary, "targets", "create", "test-auth-token-default", "--admin-port", hoverfly.GetAdminPort())

			functional_tests.Run(hoverctlBinary, "login", "--username", functional_tests.HoverflyUsername, "--password", functional_tests.HoverflyPassword)
			output := functional_tests.Run(hoverctlBinary, "config", "auth-token")

			Expect(output).To(ContainSubstring("No auth token"))
		})
	})
})

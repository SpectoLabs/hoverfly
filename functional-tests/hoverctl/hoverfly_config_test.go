package hoverctl_end_to_end

import (
	"github.com/SpectoLabs/hoverfly/functional-tests"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("When I use hoverctl", func() {

	Context("I can get the configuration being used by hoverctl", func() {

		It("prints the location of the config.yaml", func() {
			WriteConfiguration("test", "7654", "5432")
			output := functional_tests.Run(hoverctlBinary, "config")

			Expect(output).To(ContainSubstring("functional-tests/hoverctl/.hoverfly/config.yaml"))
		})

		It("prints the contents of the config.yaml", func() {
			WriteConfigurationWithAuth("test", "7654", "5432", true, "benjih", "secretpassword")

			output := functional_tests.Run(hoverctlBinary, "config")

			Expect(output).To(ContainSubstring("hoverfly.host: test"))
			Expect(output).To(ContainSubstring(`hoverfly.admin.port: "7654"`))
			Expect(output).To(ContainSubstring(`hoverfly.proxy.port: "5432"`))
			Expect(output).To(ContainSubstring("hoverfly.webserver: true"))
			Expect(output).To(ContainSubstring("hoverfly.username: benjih"))
			Expect(output).To(ContainSubstring("hoverfly.password: secretpassword"))
		})

	})

})

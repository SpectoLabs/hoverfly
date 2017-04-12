package hoverctl_suite

import (
	"github.com/SpectoLabs/hoverfly/functional-tests"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("When I use hoverctl", func() {

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

})

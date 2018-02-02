package hoverctl_suite

import (
	"github.com/SpectoLabs/hoverfly/functional-tests"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("When I export with hoverctl", func() {

	Context("without providing a path to write to", func() {

		It("it should fail nicely", func() {
			output := functional_tests.Run(hoverctlBinary, "export")

			Expect(output).To(ContainSubstring("You have not provided a path to simulation"))
			Expect(output).To(ContainSubstring("Try hoverctl export --help for more information"))
		})
	})

	Context("with a target that doesn't exist", func() {
		It("should error", func() {
			output := functional_tests.Run(hoverctlBinary, "export", "--target", "test-target")

			Expect(output).To(ContainSubstring("test-target is not a target"))
			Expect(output).To(ContainSubstring("Run `hoverctl targets create test-target`"))
		})
	})

	Context("without providing an urlPattern", func() {
		It("it should fail nicely", func() {
			output := functional_tests.Run(hoverctlBinary, "export", "--url-pattern")

			Expect(output).To(ContainSubstring("Error: flag needs an argument: --url-pattern"))
		})
	})
})

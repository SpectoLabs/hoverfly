package hoverctl_suite

import (
	"github.com/SpectoLabs/hoverfly/functional-tests"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("When I use hoverctl", func() {

	var (
		hoverfly *functional_tests.Hoverfly
	)

	Describe("with a running hoverfly", func() {

		BeforeEach(func() {
			hoverfly = functional_tests.NewHoverfly()
			hoverfly.Start()

			functional_tests.Run(hoverctlBinary, "targets", "update", "local", "--admin-port", hoverfly.GetAdminPort())
		})

		AfterEach(func() {
			hoverfly.Stop()
		})

		Context("I can get the hoverfly's destination", func() {

			It("should return the destination", func() {
				output := functional_tests.Run(hoverctlBinary, "destination")
				Expect(output).To(ContainSubstring("Current Hoverfly destination is set to ."))
			})
		})

		Context("I can set hoverfly's destination", func() {

			It("sets the destination", func() {
				output := functional_tests.Run(hoverctlBinary, "destination", "example.org")
				Expect(output).To(ContainSubstring("Hoverfly destination has been set to example.org"))

				output = functional_tests.Run(hoverctlBinary, "destination")
				Expect(output).To(ContainSubstring("Current Hoverfly destination is set to example.org"))
			})
		})

		Context("I cannot set hoverfly's destination", func() {

			It("does not set the destination if regex is invalid", func() {
				output := functional_tests.Run(hoverctlBinary, "destination", "regex[[[[")
				Expect(output).To(ContainSubstring("Regex pattern does not compile"))

				output = functional_tests.Run(hoverctlBinary, "destination")
				Expect(output).To(ContainSubstring("Current Hoverfly destination is set to ."))
			})
		})
	})

	Describe("without a running hoverfly", func() {

		Context("we can test our regex with a --dry-run", func() {

			It("does not attempt the --dry-run the destination if regex is invalid", func() {
				output := functional_tests.Run(hoverctlBinary, "destination", "regex[[[[", "--dry-run", "doesntmatter.io")
				Expect(output).To(ContainSubstring("Regex pattern does not compile"))
			})

			It("does a dry run and tests if the regex matches the URL - which it does", func() {
				output := functional_tests.Run(hoverctlBinary, "destination", "hoverfly.io", "--dry-run", "hoverfly.io")
				Expect(output).To(ContainSubstring("The regex provided matches the dry-run URL"))
			})

			It("does a dry run and tests if the regex matches the URL - which it does not", func() {
				output := functional_tests.Run(hoverctlBinary, "destination", "specto.io", "--dry-run", "hoverfly.io")
				Expect(output).To(ContainSubstring("The regex provided does not match the dry-run URL"))
			})

		})
	})

	Context("with a target that doesn't exist", func() {
		It("should error", func() {
			output := functional_tests.Run(hoverctlBinary, "destination", "--target", "test-target")

			Expect(output).To(ContainSubstring("test-target is not a target"))
			Expect(output).To(ContainSubstring("Run `hoverctl targets create test-target`"))
		})
	})
})

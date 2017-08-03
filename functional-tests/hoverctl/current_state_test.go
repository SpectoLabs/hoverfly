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

		Describe("Managing Hoverflies data using the CLI", func() {

			Describe("when the current state is empty", func() {

				It("Returns empty when getting all", func() {
					output := functional_tests.Run(hoverctlBinary, "current-state", "get-all", "--admin-port="+hoverfly.GetAdminPort())
					Expect(output).To(ContainSubstring("The current state for Hoverfly is empty"))
				})

				It("Returns empty for missing key", func() {
					output := functional_tests.Run(hoverctlBinary, "current-state", "get", "foo", "--admin-port="+hoverfly.GetAdminPort())
					Expect(output).To(ContainSubstring("State is not set for the key: foo"))
				})

			})

			Describe("when mutating state", func() {

				It("Can set, get, and get-all and delete state", func() {
					output := functional_tests.Run(hoverctlBinary, "current-state", "set", "foo", "bar", "--admin-port="+hoverfly.GetAdminPort())
					Expect(output).To(ContainSubstring("Successfully set current-state key and value:\n\"foo\"=\"bar\""))

					output = functional_tests.Run(hoverctlBinary, "current-state", "get", "foo", "--admin-port="+hoverfly.GetAdminPort())
					Expect(output).To(ContainSubstring("Current state of \"foo\":\nbar"))

					output = functional_tests.Run(hoverctlBinary, "current-state", "set", "cheese", "ham", "--admin-port="+hoverfly.GetAdminPort())
					Expect(output).To(ContainSubstring("Successfully set current-state key and value:\n\"cheese\"=\"ham\""))

					output = functional_tests.Run(hoverctlBinary, "current-state", "get", "cheese", "--admin-port="+hoverfly.GetAdminPort())
					Expect(output).To(ContainSubstring("Current state of \"cheese\":\nham"))

					output = functional_tests.Run(hoverctlBinary, "current-state", "get-all", "--admin-port="+hoverfly.GetAdminPort())
					Expect(output).To(ContainSubstring("Current state of Hoverfly:\n\"cheese\"=\"ham\"\n\"foo\"=\"bar\""))

					output = functional_tests.Run(hoverctlBinary, "current-state", "delete-all", "--admin-port="+hoverfly.GetAdminPort())
					Expect(output).To(ContainSubstring("Current state has been deleted"))

					output = functional_tests.Run(hoverctlBinary, "current-state", "get-all", "--admin-port="+hoverfly.GetAdminPort())
					Expect(output).To(ContainSubstring("The current state for Hoverfly is empty"))

					output = functional_tests.Run(hoverctlBinary, "current-state", "get", "foo", "--admin-port="+hoverfly.GetAdminPort())
					Expect(output).To(ContainSubstring("State is not set for the key: foo"))
				})

			})

		})
	})
})

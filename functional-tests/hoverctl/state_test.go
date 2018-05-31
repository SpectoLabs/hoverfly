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

			Describe("when the state is empty", func() {

				It("Returns empty when getting all", func() {
					output := functional_tests.Run(hoverctlBinary, "state", "get-all")
					Expect(output).To(ContainSubstring("The state for Hoverfly is empty"))
				})

				It("Returns empty for missing key", func() {
					output := functional_tests.Run(hoverctlBinary, "state", "get", "foo")
					Expect(output).To(ContainSubstring("State is not set for the key: foo"))
				})

			})

			Describe("when mutating state", func() {

				It("Can set, get, and get-all and delete state", func() {
					output := functional_tests.Run(hoverctlBinary, "state", "set", "foo", "bar")
					Expect(output).To(ContainSubstring("Successfully set state key and value:\n\"foo\"=\"bar\""))

					output = functional_tests.Run(hoverctlBinary, "state", "get", "foo")
					Expect(output).To(ContainSubstring("State of \"foo\":\nbar"))

					output = functional_tests.Run(hoverctlBinary, "state", "set", "cheese", "ham")
					Expect(output).To(ContainSubstring("Successfully set state key and value:\n\"cheese\"=\"ham\""))

					output = functional_tests.Run(hoverctlBinary, "state", "get", "cheese")
					Expect(output).To(ContainSubstring("State of \"cheese\":\nham"))

					output = functional_tests.Run(hoverctlBinary, "state", "get-all")
					Expect(output).To(ContainSubstring("State of Hoverfly:\n"))
					Expect(output).To(ContainSubstring(`"cheese"="ham"`))
					Expect(output).To(ContainSubstring(`"foo"="bar"`))

					output = functional_tests.Run(hoverctlBinary, "state", "delete-all")
					Expect(output).To(ContainSubstring("State has been deleted"))

					output = functional_tests.Run(hoverctlBinary, "state", "get-all")
					Expect(output).To(ContainSubstring("The state for Hoverfly is empty"))

					output = functional_tests.Run(hoverctlBinary, "state", "get", "foo")
					Expect(output).To(ContainSubstring("State is not set for the key: foo"))
				})

			})

		})
	})
})

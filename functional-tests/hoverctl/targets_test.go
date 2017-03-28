package hoverctl_suite

import (
	"github.com/SpectoLabs/hoverfly/functional-tests"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("When using the `targets` command", func() {

	Context("viewing targets", func() {
		Context("with no targets", func() {

			It("should fail nicely", func() {
				output := functional_tests.Run(hoverctlBinary, "targets")

				Expect(output).To(ContainSubstring("No targets registered"))
			})
		})

		Context("with targets", func() {

			BeforeEach(func() {
				functional_tests.Run(hoverctlBinary, "targets", "create", "--target", "default", "--admin-port", "1234")
			})

			It("print targets", func() {
				output := functional_tests.Run(hoverctlBinary, "targets")

				Expect(output).To(ContainSubstring("TARGET NAME | ADMIN PORT "))
				Expect(output).To(ContainSubstring("default"))
				Expect(output).To(ContainSubstring("1234"))
			})
		})
	})

	Context("creating targets", func() {

		It("should create the target and print it", func() {
			output := functional_tests.Run(hoverctlBinary, "targets", "create", "--target", "default", "--admin-port", "1234")

			Expect(output).To(ContainSubstring("TARGET NAME | ADMIN PORT "))
			Expect(output).To(ContainSubstring("default"))
			Expect(output).To(ContainSubstring("1234"))
		})

		It("should fail nicely if no target name is provided", func() {
			output := functional_tests.Run(hoverctlBinary, "targets", "create")

			Expect(output).To(ContainSubstring("Cannot create a target without a name"))
		})

	})

	Context("deleting targets", func() {

		BeforeEach(func() {
			functional_tests.Run(hoverctlBinary, "targets", "create", "--target", "default", "--admin-port", "1234")
		})

		It("should delete targets and print nice empty message", func() {
			output := functional_tests.Run(hoverctlBinary, "targets", "delete", "--target", "default", "--force")

			Expect(output).To(ContainSubstring("No targets registered"))
		})

		It("should fail nicely if no target name is provided", func() {
			output := functional_tests.Run(hoverctlBinary, "targets", "delete")

			Expect(output).To(ContainSubstring("Cannot delete a target without a name"))
		})
	})

})

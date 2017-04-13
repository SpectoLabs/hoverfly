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
				functional_tests.Run(hoverctlBinary, "targets", "create", "default", "--admin-port", "1234", "--proxy-port", "8765", "--host", "localhost")
			})

			It("print targets", func() {
				output := functional_tests.Run(hoverctlBinary, "targets")

				Expect(output).To(ContainSubstring("TARGET NAME"))
				Expect(output).To(ContainSubstring("HOST"))
				Expect(output).To(ContainSubstring("ADMIN PORT"))
				Expect(output).To(ContainSubstring("PROXY PORT"))

				Expect(output).To(ContainSubstring("default"))
				Expect(output).To(ContainSubstring("localhost"))
				Expect(output).To(ContainSubstring("1234"))
				Expect(output).To(ContainSubstring("8765"))
			})
		})
	})

	Context("creating targets", func() {

		It("should create the target and print it", func() {

			output := functional_tests.Run(hoverctlBinary, "targets", "create", "new-target",
				"--pid", "12345",
				"--host", "localhost",
				"--admin-port", "1234",
				"--proxy-port", "8765",
			)

			Expect(output).To(ContainSubstring("TARGET NAME"))
			Expect(output).To(ContainSubstring("PID"))
			Expect(output).To(ContainSubstring("HOST"))
			Expect(output).To(ContainSubstring("ADMIN PORT"))
			Expect(output).To(ContainSubstring("PROXY PORT"))

			Expect(output).To(ContainSubstring("new-target"))
			Expect(output).To(ContainSubstring("12345"))
			Expect(output).To(ContainSubstring("localhost"))
			Expect(output).To(ContainSubstring("1234"))
			Expect(output).To(ContainSubstring("8765"))
		})

		It("should not create a target if no target name is provided", func() {
			output := functional_tests.Run(hoverctlBinary, "targets", "create")

			Expect(output).To(ContainSubstring("Cannot create a target without a name"))
		})

		It("should not create a target if target already exists", func() {
			functional_tests.Run(hoverctlBinary, "targets", "create", "exists")
			output := functional_tests.Run(hoverctlBinary, "targets", "create", "exists")

			Expect(output).To(ContainSubstring("Target exists already exists"))
			Expect(output).To(ContainSubstring("Use a different target name or run `hoverctl targets update exists`"))
		})
	})

	Context("updating targets", func() {

		It("should update the target and print it", func() {
			functional_tests.Run(hoverctlBinary, "targets", "create", "new-target")
			output := functional_tests.Run(hoverctlBinary, "targets", "update", "new-target",
				"--pid", "12345",
				"--host", "localhost",
				"--admin-port", "1234",
				"--proxy-port", "8765",
			)

			Expect(output).To(ContainSubstring("TARGET NAME"))
			Expect(output).To(ContainSubstring("PID"))
			Expect(output).To(ContainSubstring("HOST"))
			Expect(output).To(ContainSubstring("ADMIN PORT"))
			Expect(output).To(ContainSubstring("PROXY PORT"))

			Expect(output).To(ContainSubstring("new-target"))
			Expect(output).To(ContainSubstring("12345"))
			Expect(output).To(ContainSubstring("localhost"))
			Expect(output).To(ContainSubstring("1234"))
			Expect(output).To(ContainSubstring("8765"))
		})

		It("should not update a target if no target name is provided", func() {
			output := functional_tests.Run(hoverctlBinary, "targets", "update")

			Expect(output).To(ContainSubstring("Cannot update a target without a name"))
		})

		It("should not update a target if target does not exist exists", func() {
			output := functional_tests.Run(hoverctlBinary, "targets", "update", "not-exists")

			Expect(output).To(ContainSubstring("Target not-exists does not exist"))
			Expect(output).To(ContainSubstring("Use a different target name or run `hoverctl targets create not-exists`"))
		})
	})

	Context("deleting targets", func() {

		BeforeEach(func() {
			functional_tests.Run(hoverctlBinary, "targets", "create", "default", "--admin-port", "1234")
		})

		It("should delete targets and print nice empty message", func() {
			output := functional_tests.Run(hoverctlBinary, "targets", "delete", "default", "--force")

			Expect(output).To(ContainSubstring("No targets registered"))
		})

		It("should fail nicely if no target name is provided", func() {
			output := functional_tests.Run(hoverctlBinary, "targets", "delete")

			Expect(output).To(ContainSubstring("Cannot delete a target without a name"))
		})
	})

	Context("targets default", func() {

		BeforeEach(func() {
			functional_tests.Run(hoverctlBinary, "targets", "create", "default", "--admin-port", "1234")
		})

		It("should print the default target", func() {
			output := functional_tests.Run(hoverctlBinary, "targets", "default")

			Expect(output).To(ContainSubstring("TARGET NAME"))
			Expect(output).To(ContainSubstring("HOST"))
			Expect(output).To(ContainSubstring("ADMIN PORT"))
			Expect(output).To(ContainSubstring("PROXY PORT"))

			Expect(output).To(ContainSubstring("default"))
			Expect(output).To(ContainSubstring("localhost"))
			Expect(output).To(ContainSubstring("1234"))
			Expect(output).To(ContainSubstring("8500"))
		})

		It("should set the default target when given a target name", func() {
			functional_tests.Run(hoverctlBinary, "targets", "create", "alternative", "--admin-port", "1233")
			output := functional_tests.Run(hoverctlBinary, "targets", "default", "alternative")

			Expect(output).To(ContainSubstring("TARGET NAME"))
			Expect(output).To(ContainSubstring("HOST"))
			Expect(output).To(ContainSubstring("ADMIN PORT"))
			Expect(output).To(ContainSubstring("PROXY PORT"))

			Expect(output).To(ContainSubstring("alternative"))
			Expect(output).To(ContainSubstring("localhost"))
			Expect(output).To(ContainSubstring("1233"))
			Expect(output).To(ContainSubstring("8500"))
		})

		It("should error when given an invalid target name", func() {
			output := functional_tests.Run(hoverctlBinary, "targets", "default", "alternative")

			Expect(output).To(ContainSubstring("alternative is not a target\n\nRun `hoverctl targets new alternative`"))
		})

		It("should not set default when given an invalid target name ", func() {
			functional_tests.Run(hoverctlBinary, "targets", "default", "alternative")
			output := functional_tests.Run(hoverctlBinary, "targets", "default")

			Expect(output).To(ContainSubstring("TARGET NAME"))
			Expect(output).To(ContainSubstring("HOST"))
			Expect(output).To(ContainSubstring("ADMIN PORT"))
			Expect(output).To(ContainSubstring("PROXY PORT"))

			Expect(output).To(ContainSubstring("default"))
			Expect(output).To(ContainSubstring("localhost"))
			Expect(output).To(ContainSubstring("1234"))
			Expect(output).To(ContainSubstring("8500"))
		})
	})
})

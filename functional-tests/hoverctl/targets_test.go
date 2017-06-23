package hoverctl_suite

import (
	"github.com/SpectoLabs/hoverfly/functional-tests"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("When using the `targets` command", func() {

	Context("viewing targets", func() {

		Context("with targets", func() {

			BeforeEach(func() {
				functional_tests.Run(hoverctlBinary, "targets", "update", "local", "--admin-port", "1234", "--proxy-port", "8765", "--host", "localhost")
			})

			It("print default local target", func() {
				output := functional_tests.Run(hoverctlBinary, "targets")
				targets := functional_tests.TableToSliceMapStringString(output)

				Expect(targets).To(HaveLen(1))

				Expect(targets).To(HaveKey("local"))
				Expect(targets["local"]).To(Equal(map[string]string{
					"TARGET NAME": "local",
					"HOST":        "localhost",
					"ADMIN PORT":  "1234",
					"PROXY PORT":  "8765",
					"DEFAULT":     "X",
				}))
			})

			It("prints all targets and marks default", func() {
				functional_tests.Run(hoverctlBinary, "targets", "create", "one", "--admin-port", "1234", "--proxy-port", "8765", "--host", "localhost")
				functional_tests.Run(hoverctlBinary, "targets", "create", "two", "--admin-port", "1234", "--proxy-port", "8765", "--host", "localhost")
				output := functional_tests.Run(hoverctlBinary, "targets")
				targets := functional_tests.TableToSliceMapStringString(output)

				Expect(targets).To(HaveLen(3))

				Expect(targets).To(HaveKey("local"))
				Expect(targets["local"]).To(Equal(map[string]string{
					"TARGET NAME": "local",
					"HOST":        "localhost",
					"ADMIN PORT":  "1234",
					"PROXY PORT":  "8765",
					"DEFAULT":     "X",
				}))

				Expect(targets).To(HaveKey("one"))
				Expect(targets["one"]).To(Equal(map[string]string{
					"TARGET NAME": "one",
					"HOST":        "localhost",
					"ADMIN PORT":  "1234",
					"PROXY PORT":  "8765",
					"DEFAULT":     "",
				}))

				Expect(targets).To(HaveKey("two"))
				Expect(targets["two"]).To(Equal(map[string]string{
					"TARGET NAME": "two",
					"HOST":        "localhost",
					"ADMIN PORT":  "1234",
					"PROXY PORT":  "8765",
					"DEFAULT":     "",
				}))
			})
		})
	})

	Context("creating targets", func() {

		It("should create the target and print it", func() {

			output := functional_tests.Run(hoverctlBinary, "targets", "create", "new-target",
				"--host", "localhost",
				"--admin-port", "1234",
				"--proxy-port", "8765",
			)
			targets := functional_tests.TableToSliceMapStringString(output)

			Expect(targets).To(HaveKey("new-target"))
			Expect(targets["new-target"]).To(Equal(map[string]string{
				"TARGET NAME": "new-target",
				"HOST":        "localhost",
				"ADMIN PORT":  "1234",
				"PROXY PORT":  "8765",
				"DEFAULT":     "",
			}))
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
				"--host", "localhost",
				"--admin-port", "1234",
				"--proxy-port", "8765",
			)
			targets := functional_tests.TableToSliceMapStringString(output)

			Expect(targets).To(HaveKey("new-target"))
			Expect(targets["new-target"]).To(Equal(map[string]string{
				"TARGET NAME": "new-target",
				"HOST":        "localhost",
				"ADMIN PORT":  "1234",
				"PROXY PORT":  "8765",
				"DEFAULT":     "",
			}))
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
			functional_tests.Run(hoverctlBinary, "targets", "update", "local", "--admin-port", "1234")
		})

		It("should delete targets and print nice empty message", func() {
			output := functional_tests.Run(hoverctlBinary, "targets", "delete", "local", "--force")

			Expect(output).To(ContainSubstring("No targets registered"))
		})

		It("should fail nicely if no target name is provided", func() {
			output := functional_tests.Run(hoverctlBinary, "targets", "delete")

			Expect(output).To(ContainSubstring("Cannot delete a target without a name"))
		})
	})

	Context("targets default", func() {

		BeforeEach(func() {
			functional_tests.Run(hoverctlBinary, "targets", "update", "local", "--admin-port", "1234")
		})

		It("should print the default target", func() {
			output := functional_tests.Run(hoverctlBinary, "targets", "default")
			targets := functional_tests.TableToSliceMapStringString(output)

			Expect(targets).To(HaveLen(1))

			Expect(targets).To(HaveKey("local"))
			Expect(targets["local"]).To(Equal(map[string]string{
				"TARGET NAME": "local",
				"HOST":        "localhost",
				"ADMIN PORT":  "1234",
				"PROXY PORT":  "8500",
			}))
		})

		It("should set the default target when given a target name", func() {
			functional_tests.Run(hoverctlBinary, "targets", "create", "alternative", "--admin-port", "1233")
			output := functional_tests.Run(hoverctlBinary, "targets", "default", "alternative")
			targets := functional_tests.TableToSliceMapStringString(output)

			Expect(targets).To(HaveLen(1))

			Expect(targets).To(HaveKey("alternative"))
			Expect(targets["alternative"]).To(Equal(map[string]string{
				"TARGET NAME": "alternative",
				"HOST":        "localhost",
				"ADMIN PORT":  "1233",
				"PROXY PORT":  "8500",
			}))
		})

		It("the default should automatically get updated to local if the default does not exist", func() {
			functional_tests.Run(hoverctlBinary, "targets", "create", "newdefault")
			functional_tests.Run(hoverctlBinary, "targets", "default", "newdefault")
			functional_tests.Run(hoverctlBinary, "targets", "delete", "newdefault", "--force")

			output := functional_tests.Run(hoverctlBinary, "targets", "default")

			Expect(output).To(ContainSubstring("Default target `newdefault` not found, changing default target to `local`"))
		})

		It("should error when given an invalid target name", func() {
			output := functional_tests.Run(hoverctlBinary, "targets", "default", "alternative")

			Expect(output).To(ContainSubstring("alternative is not a target\n\nRun `hoverctl targets create alternative`"))
		})

		It("should not set default when given an invalid target name ", func() {
			functional_tests.Run(hoverctlBinary, "targets", "default", "alternative")
			output := functional_tests.Run(hoverctlBinary, "targets", "default")
			targets := functional_tests.TableToSliceMapStringString(output)

			Expect(targets).To(HaveLen(1))

			Expect(targets).To(HaveKey("local"))
			Expect(targets["local"]).To(Equal(map[string]string{
				"TARGET NAME": "local",
				"HOST":        "localhost",
				"ADMIN PORT":  "1234",
				"PROXY PORT":  "8500",
			}))
		})
	})
})

package hoverctl_suite

import (
	"github.com/SpectoLabs/hoverfly/functional-tests"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("hoverctl --set-default", func() {

	BeforeEach(func() {
		functional_tests.Run(hoverctlBinary, "targets", "create", "set-default-test1")
		functional_tests.Run(hoverctlBinary, "targets", "create", "set-default-test2")
	})

	Context("on a successful command using the --target flag", func() {
		It("should change the default target to the current target", func() {
			output := functional_tests.Run(hoverctlBinary, "config", "-v",
				"--target", "set-default-test2",
				"--set-default")

			Expect(output).To(ContainSubstring("default: local"))

			output = functional_tests.Run(hoverctlBinary, "config", "-v")
			Expect(output).To(ContainSubstring("default: set-default-test2"))
		})
	})

	Context("on a successful command without the --target flag", func() {
		It("should not change the default target", func() {
			output := functional_tests.Run(hoverctlBinary, "config", "-v",
				"--set-default")

			Expect(output).To(ContainSubstring("default: local"))

			output = functional_tests.Run(hoverctlBinary, "config", "-v")
			Expect(output).To(ContainSubstring("default: local"))
		})
	})

	Context("on an unsuccessful command using the --target flag", func() {
		It("should not change the default target", func() {
			functional_tests.Run(hoverctlBinary, "mode", "-v",
				"--target", "set-default-test2",
				"--set-default")

			output := functional_tests.Run(hoverctlBinary, "config", "-v")
			Expect(output).To(ContainSubstring("default: local"))
		})
	})
})

package hoverctl_suite

import (
	"strconv"

	"github.com/SpectoLabs/hoverfly/functional-tests"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("hoverctl `stop`", func() {

	var (
		hoverfly = functional_tests.NewHoverfly()
	)

	Context("without a running instance of Hoverfly", func() {
		BeforeEach(func() {
			functional_tests.Run(hoverctlBinary, "targets", "create", "-t", "default")
		})

		It("should return an error", func() {
			output := functional_tests.Run(hoverctlBinary, "stop")

			Expect(output).To(ContainSubstring("Hoverfly is not running"))
		})
	})

	Context("with an incorrect pid", func() {
		BeforeEach(func() {
			functional_tests.Run(hoverctlBinary, "targets", "create", "-t", "default", "--pid", "432111")
		})

		It("should return an error", func() {
			output := functional_tests.Run(hoverctlBinary, "stop")

			Expect(output).To(ContainSubstring("Could not kill Hoverfly [process 432111]"))
		})
	})

	Context("with a running instance of Hoverfly", func() {
		BeforeEach(func() {
			hoverfly.Start()
			functional_tests.Run(hoverctlBinary, "targets", "create", "-t", "default", "--pid", strconv.Itoa(hoverfly.GetPid()))
		})

		AfterEach(func() {
			hoverfly.Stop()
		})

		It("stops Hoverfly", func() {
			output := functional_tests.Run(hoverctlBinary, "stop")
			Expect(output).To(ContainSubstring("Hoverfly has been stopped"))
		})

		It("removes the pid from the target", func() {
			output := functional_tests.Run(hoverctlBinary, "stop")
			Expect(output).To(ContainSubstring("Hoverfly has been stopped"))

			output = functional_tests.Run(hoverctlBinary, "targets")

			targets := functional_tests.TableToSliceMapStringString(output)
			Expect(targets["default"]["PID"]).To(Equal("0"))
		})
	})

	Context("with a target that doesn't exist", func() {
		It("should error", func() {
			output := functional_tests.Run(hoverctlBinary, "stop", "--target", "test-target")

			Expect(output).To(ContainSubstring("test-target is not a target"))
			Expect(output).To(ContainSubstring("Run `hoverctl targets new test-target`"))
		})
	})
})

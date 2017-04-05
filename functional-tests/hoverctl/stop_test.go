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

	AfterEach(func() {
		functional_tests.Run(hoverctlBinary, "targets", "delete", "-t", "default")
	})

	Context("without a target", func() {
		It("should return an error", func() {
			output := functional_tests.Run(hoverctlBinary, "stop")

			Expect(output).To(ContainSubstring("Cannot stop an instance of Hoverfly without a target"))
		})
	})

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

		It("by stopping hoverfly", func() {
			output := functional_tests.Run(hoverctlBinary, "stop")

			Expect(output).To(ContainSubstring("Hoverfly has been stopped"))
		})
	})
})

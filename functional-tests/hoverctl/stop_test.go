package hoverctl_suite

import (
	"strconv"

	"github.com/SpectoLabs/hoverfly/functional-tests"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/phayes/freeport"
)

var _ = Describe("hoverctl `stop`", func() {

	var (
		hoverfly = functional_tests.NewHoverfly()
	)

	Context("without a running instance of Hoverfly", func() {
		BeforeEach(func() {
			functional_tests.Run(hoverctlBinary, "targets", "update", "local", "--admin-port", strconv.Itoa(freeport.GetPort()))
		})

		It("should return an error", func() {
			output := functional_tests.Run(hoverctlBinary, "stop")

			Expect(output).To(ContainSubstring("Target Hoverfly is not running"))
			Expect(output).To(ContainSubstring("Run `hoverctl start -t local` to start it"))
		})
	})

	Context("with a running instance of Hoverfly", func() {
		BeforeEach(func() {
			hoverfly.Start()
			functional_tests.Run(hoverctlBinary, "targets", "update", "local", "--admin-port", hoverfly.GetAdminPort())
		})

		AfterEach(func() {
			hoverfly.Stop()
		})

		It("stops Hoverfly", func() {
			output := functional_tests.Run(hoverctlBinary, "stop")
			Expect(output).To(ContainSubstring("Hoverfly has been stopped"))
		})
	})

	Context("with a target that doesn't exist", func() {
		It("should error", func() {
			output := functional_tests.Run(hoverctlBinary, "stop", "--target", "test-target")

			Expect(output).To(ContainSubstring("test-target is not a target"))
			Expect(output).To(ContainSubstring("Run `hoverctl targets create test-target`"))
		})
	})

	Context("with a target with a remote url", func() {
		BeforeEach(func() {
			functional_tests.Run(hoverctlBinary, "targets", "create", "remote", "--host", "hoverfly.io")
		})
		It("should error", func() {
			output := functional_tests.Run(hoverctlBinary, "stop", "--target", "remote")

			Expect(output).To(ContainSubstring("Unable to stop an instance of Hoverfly on a remote host (remote host: hoverfly.io)"))
			Expect(output).To(ContainSubstring("Run `hoverctl start --new-target <name>`"))
		})
	})
})

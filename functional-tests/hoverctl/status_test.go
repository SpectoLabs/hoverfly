package hoverctl_suite

import (
	"github.com/SpectoLabs/hoverfly/functional-tests"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("when I use hoverctl status", func() {

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

		Describe("should get the status of Hoverfly", func() {

			It("Print it", func() {
				output := functional_tests.Run(hoverctlBinary, "status")
				Expect(output).To(ContainSubstring("Hoverfly   | running"))
				Expect(output).To(ContainSubstring("Admin port |    " + hoverfly.GetAdminPort()))
				Expect(output).To(ContainSubstring("Proxy port |     8500"))
				Expect(output).To(ContainSubstring("Proxy type | forward"))
				Expect(output).To(ContainSubstring("Mode       | simulate"))
				Expect(output).To(ContainSubstring("Middleware | disabled"))
			})

			It("should get the mode from Hoverfly", func() {
				hoverfly.SetMode("capture")

				output := functional_tests.Run(hoverctlBinary, "status")
				Expect(output).To(ContainSubstring("Mode       | capture"))

				hoverfly.SetMode("synthesize")

				output = functional_tests.Run(hoverctlBinary, "status")
				Expect(output).To(ContainSubstring("Mode       | synthesize"))

				hoverfly.SetMode("modify")

				output = functional_tests.Run(hoverctlBinary, "status")
				Expect(output).To(ContainSubstring("Mode       | modify"))
			})

			It("should get the middleware from Hoverfly", func() {
				output := functional_tests.Run(hoverctlBinary, "status")
				Expect(output).To(ContainSubstring("Middleware | disabled"))

				hoverfly.SetMiddleware("python", functional_tests.Middleware)

				output = functional_tests.Run(hoverctlBinary, "status")
				Expect(output).To(ContainSubstring("Middleware | enabled"))

				Expect(output).To(ContainSubstring("Hoverfly is using local middleware with the command python and the script:"))
				Expect(output).To(ContainSubstring(functional_tests.Middleware))
			})

		})
	})

	Describe("with a running hoverfly as a webserver", func() {

		BeforeEach(func() {
			hoverfly = functional_tests.NewHoverfly()
			hoverfly.Start("-webserver")

			functional_tests.Run(hoverctlBinary, "targets", "update", "local", "--admin-port", hoverfly.GetAdminPort())
		})

		AfterEach(func() {
			hoverfly.Stop()
		})

		Describe("should get the status of Hoverfly", func() {

			It("should get proxy type from Hoverfly", func() {
				output := functional_tests.Run(hoverctlBinary, "status")
				Expect(output).To(ContainSubstring("Proxy type | reverse (webserver)"))
			})
		})
	})

	Describe("without a running hoverfly", func() {

		Describe("should not get the status of Hoverfly", func() {

			It("should print an error", func() {
				output := functional_tests.Run(hoverctlBinary, "status")
				Expect(output).To(ContainSubstring("Could not connect to Hoverfly at localhost:8888"))
			})
		})
	})
})

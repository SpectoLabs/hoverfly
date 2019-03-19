package hoverctl_suite

import (
	"github.com/SpectoLabs/hoverfly/functional-tests"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("When I use hoverfly-cli", func() {

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

		Context("I can get the hoverfly's mode", func() {

			It("when hoverfly is in simulate mode", func() {
				hoverfly.SetMode("simulate")

				output := functional_tests.Run(hoverctlBinary, "mode")

				Expect(output).To(ContainSubstring("Hoverfly is currently set to simulate mode"))
			})

			It("when hoverfly is in capture mode", func() {
				hoverfly.SetMode("capture")

				output := functional_tests.Run(hoverctlBinary, "mode")

				Expect(output).To(ContainSubstring("Hoverfly is currently set to capture mode"))
			})

			It("when hoverfly is in synthesize mode", func() {
				hoverfly.SetMode("synthesize")

				output := functional_tests.Run(hoverctlBinary, "mode")

				Expect(output).To(ContainSubstring("Hoverfly is currently set to synthesize mode"))
			})

			It("when hoverfly is in modify mode", func() {
				hoverfly.SetMode("modify")

				output := functional_tests.Run(hoverctlBinary, "mode")

				Expect(output).To(ContainSubstring("Hoverfly is currently set to modify mode"))
			})
		})

		Context("I can set hoverfly's mode", func() {

			It("to simulate mode", func() {
				output := functional_tests.Run(hoverctlBinary, "mode", "simulate")

				Expect(output).To(ContainSubstring("Hoverfly has been set to simulate mode"))

				output = functional_tests.Run(hoverctlBinary, "mode")

				Expect(output).To(ContainSubstring("Hoverfly is currently set to simulate mode"))
				Expect(hoverfly.GetMode().Mode).To(Equal(simulate))
			})

			It("to simulate mode with matching strategy of strongest", func() {
				output := functional_tests.Run(hoverctlBinary, "mode", "simulate", "--matching-strategy", "strongest")

				Expect(output).To(ContainSubstring("Hoverfly has been set to simulate mode with a matching strategy of 'strongest'"))

				output = functional_tests.Run(hoverctlBinary, "mode")

				Expect(output).To(ContainSubstring("Hoverfly is currently set to simulate mode with a matching strategy of 'strongest'"))
				Expect(hoverfly.GetMode().Mode).To(Equal(simulate))
			})

			It("to simulate mode with matching strategy of first", func() {
				output := functional_tests.Run(hoverctlBinary, "mode", "simulate", "--matching-strategy", "first")

				Expect(output).To(ContainSubstring("Hoverfly has been set to simulate mode with a matching strategy of 'first'"))

				output = functional_tests.Run(hoverctlBinary, "mode")

				Expect(output).To(ContainSubstring("Hoverfly is currently set to simulate mode with a matching strategy of 'first'"))
				Expect(hoverfly.GetMode().Mode).To(Equal(simulate))
			})

			It("to simulate mode with an invalid matching strategy", func() {
				output := functional_tests.Run(hoverctlBinary, "mode", "capture")
				output = functional_tests.Run(hoverctlBinary, "mode")

				output = functional_tests.Run(hoverctlBinary, "mode", "simulate", "--matching-strategy", "invalid")

				Expect(output).To(ContainSubstring("Only matching strategy of 'first' or 'strongest' is permitted"))

				output = functional_tests.Run(hoverctlBinary, "mode")

				Expect(output).To(ContainSubstring("Hoverfly is currently set to capture mode"))
				Expect(hoverfly.GetMode().Mode).To(Equal(capture))
			})

			It("to capture mode", func() {
				output := functional_tests.Run(hoverctlBinary, "mode", "capture")

				Expect(output).To(ContainSubstring("Hoverfly has been set to capture mode"))

				output = functional_tests.Run(hoverctlBinary, "mode")

				Expect(output).To(ContainSubstring("Hoverfly is currently set to capture mode"))
				Expect(hoverfly.GetMode().Mode).To(Equal(capture))
			})

			It("to capture mode and capture all request headers", func() {
				output := functional_tests.Run(hoverctlBinary, "mode", "capture", "--all-headers")

				Expect(output).To(ContainSubstring("Hoverfly has been set to capture mode and will capture all request headers"))

				output = functional_tests.Run(hoverctlBinary, "mode")

				Expect(output).To(ContainSubstring("Hoverfly is currently set to capture mode"))
				Expect(hoverfly.GetMode().Mode).To(Equal(capture))
			})

			It("to capture mode and capture one request header", func() {
				output := functional_tests.Run(hoverctlBinary, "mode", "capture", "--headers", "Content-Type")

				Expect(output).To(ContainSubstring("Hoverfly has been set to capture mode and will capture the following request headers: [Content-Type]"))

				output = functional_tests.Run(hoverctlBinary, "mode")

				Expect(output).To(ContainSubstring("Hoverfly is currently set to capture mode"))
				Expect(hoverfly.GetMode().Mode).To(Equal(capture))
			})

			It("to capture mode and capture two request headers", func() {
				output := functional_tests.Run(hoverctlBinary, "mode", "capture", "--headers", "Content-Type,User-Agent")

				Expect(output).To(ContainSubstring("Hoverfly has been set to capture mode and will capture the following request headers: [Content-Type User-Agent]"))

				output = functional_tests.Run(hoverctlBinary, "mode")

				Expect(output).To(ContainSubstring("Hoverfly is currently set to capture mode"))
				Expect(hoverfly.GetMode().Mode).To(Equal(capture))
			})

			It("to capture mode and error if one of the headers is an asterisk", func() {
				output := functional_tests.Run(hoverctlBinary, "mode", "capture", "--headers", "Content-Type,*")

				Expect(output).To(ContainSubstring("Must provide a list containing only an asterix, or a list containing only headers names"))
			})

			It("to synthesize mode", func() {
				output := functional_tests.Run(hoverctlBinary, "mode", "synthesize")

				Expect(output).To(ContainSubstring("Hoverfly has been set to synthesize mode"))

				output = functional_tests.Run(hoverctlBinary, "mode")

				Expect(output).To(ContainSubstring("Hoverfly is currently set to synthesize mode"))
				Expect(hoverfly.GetMode().Mode).To(Equal(synthesize))
			})

			It("to modify mode", func() {
				output := functional_tests.Run(hoverctlBinary, "mode", "modify")

				Expect(output).To(ContainSubstring("Hoverfly has been set to modify mode"))

				output = functional_tests.Run(hoverctlBinary, "mode")

				Expect(output).To(ContainSubstring("Hoverfly is currently set to modify mode"))
				Expect(hoverfly.GetMode().Mode).To(Equal(modify))
			})

			It("to spy mode", func() {
				output := functional_tests.Run(hoverctlBinary, "mode", "spy")

				Expect(output).To(ContainSubstring("Hoverfly has been set to spy mode"))

				output = functional_tests.Run(hoverctlBinary, "mode")

				Expect(output).To(ContainSubstring("Hoverfly is currently set to spy mode"))
				Expect(hoverfly.GetMode().Mode).To(Equal(spy))
			})

			It("to diff mode", func() {
				output := functional_tests.Run(hoverctlBinary, "mode", "diff")

				Expect(output).To(ContainSubstring("Hoverfly has been set to diff mode"))

				output = functional_tests.Run(hoverctlBinary, "mode")

				Expect(output).To(ContainSubstring("Hoverfly is currently set to diff mode"))
				Expect(hoverfly.GetMode().Mode).To(Equal(diff))
			})

			It("to diff mode and exclude all response headers", func() {
				output := functional_tests.Run(hoverctlBinary, "mode", "diff", "--all-headers")

				Expect(output).To(ContainSubstring("Hoverfly has been set to diff mode and will exclude all response headers from diffing"))

				output = functional_tests.Run(hoverctlBinary, "mode")

				Expect(output).To(ContainSubstring("Hoverfly is currently set to diff mode"))
				Expect(hoverfly.GetMode().Mode).To(Equal(diff))
			})

			It("to diff mode and exclude one response header", func() {
				output := functional_tests.Run(hoverctlBinary, "mode", "diff", "--headers", "Content-Type")

				Expect(output).To(ContainSubstring("Hoverfly has been set to diff mode and will exclude the following response headers from diffing: [Content-Type]"))

				output = functional_tests.Run(hoverctlBinary, "mode")

				Expect(output).To(ContainSubstring("Hoverfly is currently set to diff mode"))
				Expect(hoverfly.GetMode().Mode).To(Equal(diff))
			})
		})
	})

	Describe("with a running hoverfly set to run as a webserver", func() {

		BeforeEach(func() {
			hoverfly = functional_tests.NewHoverfly()
			hoverfly.Start("-webserver")

			functional_tests.Run(hoverctlBinary, "targets", "update", "local", "--admin-port", hoverfly.GetAdminPort())
		})

		AfterEach(func() {
			hoverfly.Stop()
		})

		Context("I can get the hoverfly's mode", func() {

			It("when hoverfly is in simulate mode", func() {
				output := functional_tests.Run(hoverctlBinary, "mode")

				Expect(output).To(ContainSubstring("Hoverfly is currently set to simulate mode"))
			})

			It("when hoverfly is in synthesize mode", func() {
				hoverfly.SetMode("synthesize")

				output := functional_tests.Run(hoverctlBinary, "mode")

				Expect(output).To(ContainSubstring("Hoverfly is currently set to synthesize mode"))
			})

			It("when hoverfly is in modify mode", func() {
				hoverfly.SetMode("modify")

				output := functional_tests.Run(hoverctlBinary, "mode")

				Expect(output).To(ContainSubstring("Hoverfly is currently set to modify mode"))
			})

		})

		Context("I can set hoverfly's mode", func() {

			It("to simulate mode", func() {
				output := functional_tests.Run(hoverctlBinary, "mode", "simulate")

				Expect(output).To(ContainSubstring("Hoverfly has been set to simulate mode"))

				output = functional_tests.Run(hoverctlBinary, "mode")

				Expect(output).To(ContainSubstring("Hoverfly is currently set to simulate mode"))
				Expect(hoverfly.GetMode().Mode).To(Equal(simulate))
			})

			It("to capture mode", func() {
				output := functional_tests.Run(hoverctlBinary, "mode", "capture")

				Expect(output).To(ContainSubstring("Cannot change the mode of Hoverfly to capture when running as a webserver"))

				output = functional_tests.Run(hoverctlBinary, "mode")

				Expect(output).To(ContainSubstring("Hoverfly is currently set to simulate mode"))
				Expect(hoverfly.GetMode().Mode).To(Equal(simulate))
			})

			It("to synthesize mode", func() {
				output := functional_tests.Run(hoverctlBinary, "mode", "synthesize")

				Expect(output).To(ContainSubstring("Hoverfly has been set to synthesize mode"))

				output = functional_tests.Run(hoverctlBinary, "mode")

				Expect(output).To(ContainSubstring("Hoverfly is currently set to synthesize mode"))
				Expect(hoverfly.GetMode().Mode).To(Equal(synthesize))
			})

			It("to modify mode", func() {
				output := functional_tests.Run(hoverctlBinary, "mode", "modify")

				Expect(output).To(ContainSubstring("Hoverfly has been set to modify mode"))

				output = functional_tests.Run(hoverctlBinary, "mode")

				Expect(output).To(ContainSubstring("Hoverfly is currently set to modify mode"))
				Expect(hoverfly.GetMode().Mode).To(Equal(modify))
			})
		})
	})

	Context("with a target that doesn't exist", func() {
		It("should error", func() {
			output := functional_tests.Run(hoverctlBinary, "mode", "--target", "test-target")

			Expect(output).To(ContainSubstring("test-target is not a target"))
			Expect(output).To(ContainSubstring("Run `hoverctl targets create test-target`"))
		})
	})
})

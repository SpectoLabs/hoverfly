package hoverctl_end_to_end

import (
	"os/exec"
	"strings"

	"github.com/SpectoLabs/hoverfly/functional-tests"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	hoverfly *functional_tests.Hoverfly
)

var _ = Describe("When I use hoverfly-cli", func() {

	Describe("with a running hoverfly", func() {

		BeforeEach(func() {
			hoverfly = functional_tests.NewHoverfly()
			hoverfly.Start()

			WriteConfiguration("localhost", hoverfly.GetAdminPort(), hoverfly.GetProxyPort())
		})

		AfterEach(func() {
			hoverfly.Stop()
		})

		Context("I can get the hoverfly's mode", func() {

			It("when hoverfly is in simulate mode", func() {
				hoverfly.SetMode("simulate")

				out, _ := exec.Command(hoverctlBinary, "mode").Output()

				output := strings.TrimSpace(string(out))
				Expect(output).To(ContainSubstring("Hoverfly is set to simulate mode"))
			})

			It("when hoverfly is in capture mode", func() {
				hoverfly.SetMode("capture")

				out, _ := exec.Command(hoverctlBinary, "mode").Output()

				output := strings.TrimSpace(string(out))
				Expect(output).To(ContainSubstring("Hoverfly is set to capture mode"))
			})

			It("when hoverfly is in synthesize mode", func() {
				hoverfly.SetMode("synthesize")

				out, _ := exec.Command(hoverctlBinary, "mode").Output()

				output := strings.TrimSpace(string(out))
				Expect(output).To(ContainSubstring("Hoverfly is set to synthesize mode"))
			})

			It("when hoverfly is in modify mode", func() {
				hoverfly.SetMode("modify")

				out, _ := exec.Command(hoverctlBinary, "mode").Output()

				output := strings.TrimSpace(string(out))
				Expect(output).To(ContainSubstring("Hoverfly is set to modify mode"))
			})
		})

		Context("I can set hoverfly's mode", func() {

			It("to simulate mode", func() {
				setOutput, _ := exec.Command(hoverctlBinary, "mode", "simulate").Output()

				output := strings.TrimSpace(string(setOutput))
				Expect(output).To(ContainSubstring("Hoverfly has been set to simulate mode"))

				getOutput, _ := exec.Command(hoverctlBinary, "mode").Output()

				output = strings.TrimSpace(string(getOutput))
				Expect(output).To(ContainSubstring("Hoverfly is set to simulate mode"))
				Expect(hoverfly.GetMode()).To(Equal(simulate))
			})

			It("to capture mode", func() {
				setOutput, _ := exec.Command(hoverctlBinary, "mode", "capture").Output()

				output := strings.TrimSpace(string(setOutput))
				Expect(output).To(ContainSubstring("Hoverfly has been set to capture mode"))

				getOutput, _ := exec.Command(hoverctlBinary, "mode").Output()

				output = strings.TrimSpace(string(getOutput))
				Expect(output).To(ContainSubstring("Hoverfly is set to capture mode"))
				Expect(hoverfly.GetMode()).To(Equal(capture))
			})

			It("to synthesize mode", func() {
				setOutput, _ := exec.Command(hoverctlBinary, "mode", "synthesize").Output()

				output := strings.TrimSpace(string(setOutput))
				Expect(output).To(ContainSubstring("Hoverfly has been set to synthesize mode"))

				getOutput, _ := exec.Command(hoverctlBinary, "mode").Output()

				output = strings.TrimSpace(string(getOutput))
				Expect(output).To(ContainSubstring("Hoverfly is set to synthesize mode"))
				Expect(hoverfly.GetMode()).To(Equal(synthesize))
			})

			It("to modify mode", func() {
				setOutput, _ := exec.Command(hoverctlBinary, "mode", "modify").Output()

				output := strings.TrimSpace(string(setOutput))
				Expect(output).To(ContainSubstring("Hoverfly has been set to modify mode"))

				getOutput, _ := exec.Command(hoverctlBinary, "mode").Output()

				output = strings.TrimSpace(string(getOutput))
				Expect(output).To(ContainSubstring("Hoverfly is set to modify mode"))
				Expect(hoverfly.GetMode()).To(Equal(modify))
			})
		})
	})

	Describe("with a running hoverfly set to run as a webserver", func() {

		BeforeEach(func() {
			hoverfly = functional_tests.NewHoverfly()
			hoverfly.Start("-webserver")

			WriteConfiguration("localhost", hoverfly.GetAdminPort(), hoverfly.GetProxyPort())
		})

		AfterEach(func() {
			hoverfly.Stop()
		})

		Context("I can get the hoverfly's mode", func() {

			It("when hoverfly is in simulate mode", func() {
				out, _ := exec.Command(hoverctlBinary, "mode").Output()

				output := strings.TrimSpace(string(out))
				Expect(output).To(ContainSubstring("Hoverfly is set to simulate mode"))
			})
		})

		Context("I cannot set hoverfly's mode", func() {

			It("to simulate mode", func() {
				setOutput, _ := exec.Command(hoverctlBinary, "mode", "simulate").Output()

				output := strings.TrimSpace(string(setOutput))
				Expect(output).To(ContainSubstring("Cannot change the mode of Hoverfly when running as a webserver"))

				getOutput, _ := exec.Command(hoverctlBinary, "mode").Output()

				output = strings.TrimSpace(string(getOutput))
				Expect(output).To(ContainSubstring("Hoverfly is set to simulate mode"))
				Expect(hoverfly.GetMode()).To(Equal(simulate))
			})

			It("to capture mode", func() {
				setOutput, _ := exec.Command(hoverctlBinary, "mode", "capture").Output()

				output := strings.TrimSpace(string(setOutput))
				Expect(output).To(ContainSubstring("Cannot change the mode of Hoverfly when running as a webserver"))

				getOutput, _ := exec.Command(hoverctlBinary, "mode").Output()

				output = strings.TrimSpace(string(getOutput))
				Expect(output).To(ContainSubstring("Hoverfly is set to simulate mode"))
				Expect(hoverfly.GetMode()).To(Equal(simulate))
			})

			It("to synthesize mode", func() {
				setOutput, _ := exec.Command(hoverctlBinary, "mode", "synthesize").Output()

				output := strings.TrimSpace(string(setOutput))
				Expect(output).To(ContainSubstring("Cannot change the mode of Hoverfly when running as a webserver"))

				getOutput, _ := exec.Command(hoverctlBinary, "mode").Output()

				output = strings.TrimSpace(string(getOutput))
				Expect(output).To(ContainSubstring("Hoverfly is set to simulate mode"))
				Expect(hoverfly.GetMode()).To(Equal(simulate))
			})

			It("to modify mode", func() {
				setOutput, _ := exec.Command(hoverctlBinary, "mode", "modify").Output()

				output := strings.TrimSpace(string(setOutput))
				Expect(output).To(ContainSubstring("Cannot change the mode of Hoverfly when running as a webserver"))

				getOutput, _ := exec.Command(hoverctlBinary, "mode").Output()

				output = strings.TrimSpace(string(getOutput))
				Expect(output).To(ContainSubstring("Hoverfly is set to simulate mode"))
				Expect(hoverfly.GetMode()).To(Equal(simulate))
			})
		})
	})
})

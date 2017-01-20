package hoverctl_end_to_end

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/phayes/freeport"
	"os/exec"
	"strconv"
	"strings"
)

var _ = Describe("When I use hoverfly-cli", func() {
	var (
		hoverflyCmd *exec.Cmd

		adminPort         = freeport.GetPort()
		adminPortAsString = strconv.Itoa(adminPort)

		proxyPort         = freeport.GetPort()
		proxyPortAsString = strconv.Itoa(proxyPort)
	)

	Describe("with a running hoverfly", func() {

		BeforeEach(func() {
			hoverflyCmd = startHoverfly(adminPort, proxyPort, workingDirectory)
			WriteConfiguration("localhost", adminPortAsString, proxyPortAsString)
		})

		AfterEach(func() {
			hoverflyCmd.Process.Kill()
		})

		Context("I can get the hoverfly's mode", func() {

			It("when hoverfly is in simulate mode", func() {
				SetHoverflyMode("simulate", adminPort)

				out, _ := exec.Command(hoverctlBinary, "mode").Output()

				output := strings.TrimSpace(string(out))
				Expect(output).To(ContainSubstring("Hoverfly is set to simulate mode"))
			})

			It("when hoverfly is in capture mode", func() {
				SetHoverflyMode("capture", adminPort)

				out, _ := exec.Command(hoverctlBinary, "mode").Output()

				output := strings.TrimSpace(string(out))
				Expect(output).To(ContainSubstring("Hoverfly is set to capture mode"))
			})

			It("when hoverfly is in synthesize mode", func() {
				SetHoverflyMode("synthesize", adminPort)

				out, _ := exec.Command(hoverctlBinary, "mode").Output()

				output := strings.TrimSpace(string(out))
				Expect(output).To(ContainSubstring("Hoverfly is set to synthesize mode"))
			})

			It("when hoverfly is in modify mode", func() {
				SetHoverflyMode("modify", adminPort)

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
				Expect(GetHoverflyMode(adminPort)).To(Equal(simulate))
			})

			It("to capture mode", func() {
				setOutput, _ := exec.Command(hoverctlBinary, "mode", "capture").Output()

				output := strings.TrimSpace(string(setOutput))
				Expect(output).To(ContainSubstring("Hoverfly has been set to capture mode"))

				getOutput, _ := exec.Command(hoverctlBinary, "mode").Output()

				output = strings.TrimSpace(string(getOutput))
				Expect(output).To(ContainSubstring("Hoverfly is set to capture mode"))
				Expect(GetHoverflyMode(adminPort)).To(Equal(capture))
			})

			It("to synthesize mode", func() {
				setOutput, _ := exec.Command(hoverctlBinary, "mode", "synthesize").Output()

				output := strings.TrimSpace(string(setOutput))
				Expect(output).To(ContainSubstring("Hoverfly has been set to synthesize mode"))

				getOutput, _ := exec.Command(hoverctlBinary, "mode").Output()

				output = strings.TrimSpace(string(getOutput))
				Expect(output).To(ContainSubstring("Hoverfly is set to synthesize mode"))
				Expect(GetHoverflyMode(adminPort)).To(Equal(synthesize))
			})

			It("to modify mode", func() {
				setOutput, _ := exec.Command(hoverctlBinary, "mode", "modify").Output()

				output := strings.TrimSpace(string(setOutput))
				Expect(output).To(ContainSubstring("Hoverfly has been set to modify mode"))

				getOutput, _ := exec.Command(hoverctlBinary, "mode").Output()

				output = strings.TrimSpace(string(getOutput))
				Expect(output).To(ContainSubstring("Hoverfly is set to modify mode"))
				Expect(GetHoverflyMode(adminPort)).To(Equal(modify))
			})
		})
	})

	Describe("with a running hoverfly set to run as a webserver", func() {

		BeforeEach(func() {
			hoverflyCmd = startHoverflyWebserver(adminPort, proxyPort, workingDirectory)
			WriteConfiguration("localhost", adminPortAsString, proxyPortAsString)
		})

		AfterEach(func() {
			hoverflyCmd.Process.Kill()
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
				Expect(GetHoverflyMode(adminPort)).To(Equal(simulate))
			})

			It("to capture mode", func() {
				setOutput, _ := exec.Command(hoverctlBinary, "mode", "capture").Output()

				output := strings.TrimSpace(string(setOutput))
				Expect(output).To(ContainSubstring("Cannot change the mode of Hoverfly when running as a webserver"))

				getOutput, _ := exec.Command(hoverctlBinary, "mode").Output()

				output = strings.TrimSpace(string(getOutput))
				Expect(output).To(ContainSubstring("Hoverfly is set to simulate mode"))
				Expect(GetHoverflyMode(adminPort)).To(Equal(simulate))
			})

			It("to synthesize mode", func() {
				setOutput, _ := exec.Command(hoverctlBinary, "mode", "synthesize").Output()

				output := strings.TrimSpace(string(setOutput))
				Expect(output).To(ContainSubstring("Cannot change the mode of Hoverfly when running as a webserver"))

				getOutput, _ := exec.Command(hoverctlBinary, "mode").Output()

				output = strings.TrimSpace(string(getOutput))
				Expect(output).To(ContainSubstring("Hoverfly is set to simulate mode"))
				Expect(GetHoverflyMode(adminPort)).To(Equal(simulate))
			})

			It("to modify mode", func() {
				setOutput, _ := exec.Command(hoverctlBinary, "mode", "modify").Output()

				output := strings.TrimSpace(string(setOutput))
				Expect(output).To(ContainSubstring("Cannot change the mode of Hoverfly when running as a webserver"))

				getOutput, _ := exec.Command(hoverctlBinary, "mode").Output()

				output = strings.TrimSpace(string(getOutput))
				Expect(output).To(ContainSubstring("Hoverfly is set to simulate mode"))
				Expect(GetHoverflyMode(adminPort)).To(Equal(simulate))
			})
		})
	})
})

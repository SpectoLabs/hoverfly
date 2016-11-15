package hoverctl_end_to_end

import (
	"os/exec"
	"strconv"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/phayes/freeport"
)

var _ = Describe("When I use hoverctl", func() {
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

		Context("I can get the hoverfly's destination", func() {

			It("should return the destination", func() {
				out, _ := exec.Command(hoverctlBinary, "destination").Output()

				output := strings.TrimSpace(string(out))
				Expect(output).To(ContainSubstring("The destination in Hoverfly is set to ."))
			})
		})

		Context("I can set hoverfly's destination", func() {

			It("sets the destination", func() {
				setOutput, _ := exec.Command(hoverctlBinary, "destination", "example.org").Output()

				output := strings.TrimSpace(string(setOutput))
				Expect(output).To(ContainSubstring("The destination in Hoverfly has been set to example.org"))

				getOutput, _ := exec.Command(hoverctlBinary, "destination").Output()

				output = strings.TrimSpace(string(getOutput))
				Expect(output).To(ContainSubstring("The destination in Hoverfly is set to example.org"))
			})
		})
	})
})

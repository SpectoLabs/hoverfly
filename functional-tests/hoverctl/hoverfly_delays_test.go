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

		Context("I can get the response delay config set on hoverfly", func() {

			It("when no delay is set", func() {
				SetHoverflyMode("simulate", adminPort)

				out, _ := exec.Command(hoverctlBinary, "delays").Output()

				output := strings.TrimSpace(string(out))
				Expect(output).To(ContainSubstring("Hoverfly has no delays configured"))
			})

		})

		Context("I can update the response delay config set on hoverfly", func() {

			It("when no delay is set", func() {
				SetHoverflyMode("simulate", adminPort)

				out, _ := exec.Command(hoverctlBinary, "delays", "testdata/delays.json").Output()

				output := strings.TrimSpace(string(out))
				Expect(output).To(ContainSubstring("Response delays set in Hoverfly"))
				Expect(output).To(ContainSubstring("host1 - 100ms"))
				Expect(output).To(ContainSubstring("POST | host2 - 110ms"))
			})

		})

	})
})

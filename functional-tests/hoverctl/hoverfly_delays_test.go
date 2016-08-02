package hoverfly_end_to_end_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os/exec"
	"strings"
	"strconv"
	"github.com/phayes/freeport"
)

var _ = Describe("When I use hoverfly-cli", func() {
	var (
		hoverflyCmd *exec.Cmd

		adminPort = freeport.GetPort()
		adminPortAsString = strconv.Itoa(adminPort)

		proxyPort = freeport.GetPort()
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
				// TODO: when no delays are set we should really have a nice message
				Expect(len(output)).To(Equal(0))
			})

		})

		Context("I can update the response delay config set on hoverfly", func() {

			It("when no delay is set", func() {
				SetHoverflyMode("simulate", adminPort)

				out, _ := exec.Command(hoverctlBinary, "delays", "testdata/delays.json").Output()

				output := strings.TrimSpace(string(out))
				Expect(output).To(ContainSubstring("Response delays set in Hoverfly"))
				Expect(output).To(ContainSubstring("UrlPattern:host1"))
				Expect(output).To(ContainSubstring("Delay:110"))
				Expect(output).To(ContainSubstring("UrlPattern:host2"))
				Expect(output).To(ContainSubstring("Delay:100"))
			})

		})



	})
})
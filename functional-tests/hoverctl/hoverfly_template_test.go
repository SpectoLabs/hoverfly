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

	Context("with a running hoverfly", func() {

		BeforeEach(func() {
			hoverflyCmd = startHoverfly(adminPort, proxyPort, workingDirectory)
			WriteConfiguration("localhost", adminPortAsString, proxyPortAsString)
		})

		AfterEach(func() {
			hoverflyCmd.Process.Kill()
		})

		Context("With no templates imported", func() {

			It("(running export) will print out an empty string", func() {

				out, _ := exec.Command(hoverctlBinary, "templates").Output()

				output := strings.TrimSpace(string(out))
				// TODO: when no delays are set we should really have a nice message
				Expect(len(output)).To(Equal(0))
			})

			It("is possible to set request templates by import", func() {

				out, _ := exec.Command(hoverctlBinary, "templates", "import", "testdata/request-template.json").Output()

				output := strings.TrimSpace(string(out))
				Expect(output).To(ContainSubstring("Response templates set in Hoverfly"))
				Expect(output).To(ContainSubstring("Path:/path1"))
				Expect(output).To(ContainSubstring("Path:/path2"))
			})
		})
	})
})
package hoverfly_end_to_end_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os/exec"
	"strings"
	"strconv"
	"github.com/phayes/freeport"
)

var _ = Describe("When I use hoverctl", func() {
	var (
		hoverflyCmd *exec.Cmd

		adminPort = freeport.GetPort()
		adminPortAsString = strconv.Itoa(adminPort)

		proxyPort = freeport.GetPort()
		proxyPortAsString = strconv.Itoa(proxyPort)
	)

	Describe("with a running hoverfly which has middleware configured", func() {

		BeforeEach(func() {
			hoverflyCmd = startHoverflyWithMiddleware(adminPort, proxyPort, workingDirectory, "python middleware.py")
			WriteConfiguration("localhost", adminPortAsString, proxyPortAsString)
		})

		AfterEach(func() {
			hoverflyCmd.Process.Kill()
		})

		It("I can get the hoverfly's middleware", func() {
			out, _ := exec.Command(hoverctlBinary, "middleware").Output()

			output := strings.TrimSpace(string(out))
			Expect(output).To(ContainSubstring("Hoverfly is currently set to run the following as middleware"))
			Expect(output).To(ContainSubstring("python middleware.py"))
		})

		It("I can set the hoverfly's middleware", func() {
			out, _ := exec.Command(hoverctlBinary, "middleware", `python testdata/add_random_delay.py`).Output()

			output := strings.TrimSpace(string(out))
			Expect(output).To(ContainSubstring("Hoverfly is now set to run the following as middleware"))
			Expect(output).To(ContainSubstring("python testdata/add_random_delay.py"))
		})
	})
})
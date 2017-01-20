package hoverctl_end_to_end

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/phayes/freeport"
	"os/exec"
	"strconv"
	"strings"
)

var _ = Describe("When I use hoverctl", func() {
	var (
		hoverflyCmd *exec.Cmd

		adminPort         = freeport.GetPort()
		adminPortAsString = strconv.Itoa(adminPort)

		proxyPort         = freeport.GetPort()
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

		It("I cannae set the hoverfly's middleware when specifying non-existing file", func() {
			out, _ := exec.Command(hoverctlBinary, "middleware", `python testdata/not_a_real_file.fake`).Output()

			output := strings.TrimSpace(string(out))
			Expect(output).To(ContainSubstring("Hoverfly could not execute this middleware"))
		})

		It("When I use the verbose flag, I see that python exited with status 2", func() {
			out, _ := exec.Command(hoverctlBinary, "-v", "middleware", `python testdata/not_a_real_file.fake`).Output()

			output := strings.TrimSpace(string(out))
			Expect(output).To(ContainSubstring("Hoverfly could not execute this middleware"))
			Expect(output).To(ContainSubstring("Invalid middleware: exit status 2"))
		})

		It("When I use the verbose flag, I see that notpython is not an executable", func() {
			out, _ := exec.Command(hoverctlBinary, "-v", "middleware", `notpython testdata/add_random_delay.py`).Output()

			output := strings.TrimSpace(string(out))
			Expect(output).To(ContainSubstring("Hoverfly could not execute this middleware"))
			Expect(output).To(ContainSubstring(`Invalid middleware: exec: \"notpython\": executable file not found in $PATH`))
		})
	})
})

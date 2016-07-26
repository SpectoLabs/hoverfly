package hoverfly_end_to_end_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os/exec"
	"github.com/phayes/freeport"
	"strconv"
	"os"
	"path/filepath"
	"io/ioutil"
	"github.com/dghubble/sling"
	"fmt"
	"strings"
)

var _ = Describe("When I use hoverctl", func() {
	var (
		adminPort = strconv.Itoa(freeport.GetPort())
		proxyPort = strconv.Itoa(freeport.GetPort())
	)

	BeforeEach(func() {
		WriteConfiguration("localhost", adminPort, proxyPort)
	})

	AfterEach(func() {
		exec.Command(hoverctlBinary, "stop", "--admin-port=" + adminPort, "--proxy-port=" + proxyPort).Output()
	})

	Context("I can get the logs using the log command", func() {

		It("should return the logs", func() {
			exec.Command(hoverctlBinary, "start", "--admin-port=" + adminPort, "--proxy-port=" + proxyPort).Output()

			out, _ := exec.Command(hoverctlBinary, "logs", "--admin-port=" + adminPort, "--proxy-port=" + proxyPort).Output()

			output := strings.TrimSpace(string(out))
			Expect(output).To(ContainSubstring("listening on :" + adminPort))
		})

		It("should return an error if the logs don't exist", func() {
			exec.Command(hoverctlBinary, "start", "--admin-port=" + adminPort, "--proxy-port=" + proxyPort).Output()
			
			out, _ := exec.Command(hoverctlBinary, "logs", "--admin-port=hotdogs", "--proxy-port=burgers").Output()

			output := strings.TrimSpace(string(out))
			Expect(output).To(ContainSubstring("Could not open Hoverfly log file"))
		})
	})

	Context("and start Hoverfly using hoverctl", func() {

		Context("the logs get captured in a .log file", func() {

			It("and I can see it has started", func() {
				exec.Command(hoverctlBinary, "start", "--admin-port=" + adminPort, "--proxy-port=" + proxyPort).Output()

				workingDir, _ := os.Getwd()
				filePath := filepath.Join(workingDir, ".hoverfly/", "hoverfly." + adminPort + "." + proxyPort +".log")

				file, err := ioutil.ReadFile(filePath)
				Expect(err).To(BeNil())

				Expect(string(file)).To(ContainSubstring("listening on :" + adminPort))
			})

			It("and they get updated when you use hoverfly", func() {
				exec.Command(hoverctlBinary, "start", "--admin-port=" + adminPort, "--proxy-port=" + proxyPort).Output()

				adminPortAsString, _ := strconv.Atoi(adminPort)

				SetHoverflyMode("capture", adminPortAsString)

				workingDir, _ := os.Getwd()
				filePath := filepath.Join(workingDir, ".hoverfly/", "hoverfly." + adminPort + "." + proxyPort +".log")

				file, err := ioutil.ReadFile(filePath)
				Expect(err).To(BeNil())

				Expect(string(file)).To(ContainSubstring("Handling state change request!"))
				Expect(string(file)).To(ContainSubstring(`{\"mode\":\"capture\"}`))
			})

			It("and the stderr is captured in the log file", func() {
				exec.Command(hoverctlBinary, "start", "--admin-port=" + adminPort, "--proxy-port=" + proxyPort).Output()

				req := sling.New().Post(fmt.Sprintf("http://localhost:%v/api/state", adminPort)).Body(strings.NewReader(`{"mode":"not-a-mode"}`))
				DoRequest(req)

				workingDir, _ := os.Getwd()
				filePath := filepath.Join(workingDir, ".hoverfly/", "hoverfly." + adminPort + "." + proxyPort +".log")

				file, err := ioutil.ReadFile(filePath)
				Expect(err).To(BeNil())

				Expect(string(file)).To(ContainSubstring("Wrong mode found, can't change state"))
				Expect(string(file)).To(ContainSubstring("not-a-mode"))
			})
		})
	})
})
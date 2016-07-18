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

	Context("I can get the logs using the log command", func() {

		BeforeEach(func() {
			_, err := exec.Command(hoverctlBinary, "start", "--admin-port=" + adminPort, "--proxy-port=" + proxyPort).Output()
			Expect(err).To(BeNil())
		})

		AfterEach(func() {
			_, err := exec.Command(hoverctlBinary, "stop", "--admin-port=" + adminPort, "--proxy-port=" + proxyPort).Output()
			Expect(err).To(BeNil())
		})

		It("should return the logs", func() {
			out, _ := exec.Command(hoverctlBinary, "logs", "--admin-port=" + adminPort, "--proxy-port=" + proxyPort).Output()

			output := strings.TrimSpace(string(out))
			Expect(output).To(ContainSubstring("listening on :" + adminPort))
		})

		It("should return an error if the logs don't exist", func() {
			out, _ := exec.Command(hoverctlBinary, "logs", "--admin-port=hotdogs", "--proxy-port=burgers").Output()

			output := strings.TrimSpace(string(out))
			Expect(output).To(ContainSubstring("Could not open Hoverfly log file"))
		})
	})

	Describe("and start Hoverfly using hoverctl", func() {

		Context("the logs get captured in a .log file", func() {
			BeforeEach(func() {
				_, err := exec.Command(hoverctlBinary, "start", "--admin-port=" + adminPort, "--proxy-port=" + proxyPort).Output()
				Expect(err).To(BeNil())
			})

			AfterEach(func() {
				_, err := exec.Command(hoverctlBinary, "stop", "--admin-port=" + adminPort, "--proxy-port=" + proxyPort).Output()
				Expect(err).To(BeNil())
			})

			It("and I can see it has started", func() {
				workingDir, _ := os.Getwd()
				filePath := filepath.Join(workingDir, ".hoverfly/", "hoverfly." + adminPort + "." + proxyPort +".log")

				file, err := ioutil.ReadFile(filePath)
				Expect(err).To(BeNil())

				Expect(string(file)).To(ContainSubstring("listening on :" + adminPort))
			})

			It("and they get updated when you use hoverfly", func() {

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
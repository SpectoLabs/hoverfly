package hoverctl_suite

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/phayes/freeport"
)

var _ = Describe("When I use hoverctl", func() {

	var (
		adminPort = strconv.Itoa(freeport.GetPort())
		proxyPort = strconv.Itoa(freeport.GetPort())
	)

	BeforeEach(func() {
		functional_tests.Run(hoverctlBinary, "targets", "create", "default", "--admin-port", adminPort, "--proxy-port", proxyPort)
	})

	AfterEach(func() {
		functional_tests.Run(hoverctlBinary, "stop")
	})

	Context("I can get the logs using the log command in JSON format", func() {

		It("should return the logs", func() {
			functional_tests.Run(hoverctlBinary, "start", "--admin-port="+adminPort, "--proxy-port="+proxyPort)

			output := functional_tests.Run(hoverctlBinary, "logs", "--json", "--admin-port="+adminPort, "--proxy-port="+proxyPort)

			Expect(output).To(ContainSubstring("listening on :" + adminPort))
		})

		It("should return an error if the logs don't exist", func() {
			functional_tests.Run(hoverctlBinary, "start", "--admin-port="+adminPort, "--proxy-port="+proxyPort)
			functional_tests.Run(hoverctlBinary, "targets", "create", "incorrect", "--admin-port", "12345", "--proxy-port", "65432")

			output := functional_tests.Run(hoverctlBinary, "logs", "--json", "-t", "incorrect")

			Expect(output).To(ContainSubstring("Could not open Hoverfly log file"))
		})
	})

	Context("I can get the logs using the log command in plaintext format", func() {

		It("should return the logs", func() {
			functional_tests.Run(hoverctlBinary, "start", "--admin-port="+adminPort, "--proxy-port="+proxyPort)

			output := functional_tests.Run(hoverctlBinary, "logs", "--admin-port="+adminPort, "--proxy-port="+proxyPort)

			Expect(output).To(ContainSubstring("level=info msg=\"Proxy prepared...\" Destination=. Mode=simulate ProxyPort=" + proxyPort))
		})

		It("should return an error if the logs don't exist", func() {
			functional_tests.Run(hoverctlBinary, "start", "--admin-port="+adminPort, "--proxy-port="+proxyPort)
			functional_tests.Run(hoverctlBinary, "targets", "create", "incorrect", "--admin-port", "12345", "--proxy-port", "65432")

			output := functional_tests.Run(hoverctlBinary, "logs", "-t", "incorrect")

			Expect(output).To(ContainSubstring("Could not connect to Hoverfly at localhost:12345"))
		})
	})

	Context("and start Hoverfly using hoverctl", func() {

		Context("the logs get captured in a .log file", func() {

			It("and I can see it has started", func() {
				functional_tests.Run(hoverctlBinary, "start", "--admin-port="+adminPort, "--proxy-port="+proxyPort)

				workingDir, _ := os.Getwd()
				filePath := filepath.Join(workingDir, ".hoverfly/", "hoverfly."+adminPort+"."+proxyPort+".log")

				file, err := ioutil.ReadFile(filePath)
				Expect(err).To(BeNil())

				Expect(string(file)).To(ContainSubstring("Admin interface is starting..."))
				Expect(string(file)).To(ContainSubstring("serving proxy"))
			})

			It("and they get updated when you use hoverfly", func() {
				functional_tests.Run(hoverctlBinary, "start", "--admin-port="+adminPort, "--proxy-port="+proxyPort)

				functional_tests.Run(hoverctlBinary, "mode", "capture")

				workingDir, _ := os.Getwd()
				filePath := filepath.Join(workingDir, ".hoverfly/", "hoverfly."+adminPort+"."+proxyPort+".log")

				file, err := ioutil.ReadFile(filePath)
				Expect(err).To(BeNil())

				Expect(string(file)).To(ContainSubstring("Started GET /api/health"))
			})

			It("and the stderr is captured in the log file", func() {
				functional_tests.Run(hoverctlBinary, "start", "--admin-port="+adminPort, "--proxy-port="+proxyPort)

				req := sling.New().Put(fmt.Sprintf("http://localhost:%v/api/v2/hoverfly/mode", adminPort)).Body(strings.NewReader(`{"mode":"not-a-mode"}`))
				functional_tests.DoRequest(req)

				workingDir, _ := os.Getwd()
				filePath := filepath.Join(workingDir, ".hoverfly/", "hoverfly."+adminPort+"."+proxyPort+".log")

				file, err := ioutil.ReadFile(filePath)
				Expect(err).To(BeNil())

				Expect(string(file)).To(ContainSubstring("Can't change mode to"))
				Expect(string(file)).To(ContainSubstring("not-a-mode"))
			})
		})
	})

	Context("with a target that doesn't exist", func() {
		It("should error", func() {
			output := functional_tests.Run(hoverctlBinary, "logs", "--target", "test-target")

			Expect(output).To(ContainSubstring("test-target is not a target"))
			Expect(output).To(ContainSubstring("Run `hoverctl targets new test-target`"))
		})
	})
})

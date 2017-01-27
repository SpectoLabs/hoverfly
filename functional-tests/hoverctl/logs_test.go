package hoverctl_end_to_end

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
		WriteConfiguration("localhost", adminPort, proxyPort)
	})

	AfterEach(func() {
		functional_tests.Run(hoverctlBinary, "stop", "--admin-port="+adminPort, "--proxy-port="+proxyPort)
	})

	Context("I can get the logs using the log command", func() {

		It("should return the logs", func() {
			functional_tests.Run(hoverctlBinary, "start", "--admin-port="+adminPort, "--proxy-port="+proxyPort)

			output := functional_tests.Run(hoverctlBinary, "logs", "--admin-port="+adminPort, "--proxy-port="+proxyPort)

			Expect(output).To(ContainSubstring("listening on :" + adminPort))
		})

		It("should return an error if the logs don't exist", func() {
			functional_tests.Run(hoverctlBinary, "start", "--admin-port="+adminPort, "--proxy-port="+proxyPort)

			output := functional_tests.Run(hoverctlBinary, "logs", "--admin-port=hotdogs", "--proxy-port=burgers")

			Expect(output).To(ContainSubstring("Could not open Hoverfly log file"))
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

				Expect(string(file)).To(ContainSubstring("listening on :" + adminPort))
			})

			It("and they get updated when you use hoverfly", func() {
				functional_tests.Run(hoverctlBinary, "start", "--admin-port="+adminPort, "--proxy-port="+proxyPort)

				functional_tests.Run(hoverctlBinary, "mode", "capture")

				workingDir, _ := os.Getwd()
				filePath := filepath.Join(workingDir, ".hoverfly/", "hoverfly."+adminPort+"."+proxyPort+".log")

				file, err := ioutil.ReadFile(filePath)
				Expect(err).To(BeNil())

				Expect(string(file)).To(ContainSubstring("Started GET /api/v2/hoverfly/mode"))
			})

			It("and the stderr is captured in the log file", func() {
				functional_tests.Run(hoverctlBinary, "start", "--admin-port="+adminPort, "--proxy-port="+proxyPort)

				req := sling.New().Post(fmt.Sprintf("http://localhost:%v/api/state", adminPort)).Body(strings.NewReader(`{"mode":"not-a-mode"}`))
				functional_tests.DoRequest(req)

				workingDir, _ := os.Getwd()
				filePath := filepath.Join(workingDir, ".hoverfly/", "hoverfly."+adminPort+"."+proxyPort+".log")

				file, err := ioutil.ReadFile(filePath)
				Expect(err).To(BeNil())

				Expect(string(file)).To(ContainSubstring("Wrong mode found, can't change state"))
				Expect(string(file)).To(ContainSubstring("not-a-mode"))
			})
		})
	})
})

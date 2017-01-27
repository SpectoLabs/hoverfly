package hoverctl_end_to_end

import (
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"

	"github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/phayes/freeport"
)

var _ = Describe("When I use hoverctl", func() {

	Describe("without a running hoverfly", func() {

		var (
			adminPort         = freeport.GetPort()
			adminPortAsString = strconv.Itoa(adminPort)

			proxyPort         = freeport.GetPort()
			proxyPortAsString = strconv.Itoa(proxyPort)
		)

		BeforeEach(func() {
			exec.Command(hoverctlBinary, "stop", "-v").Run()
			WriteConfiguration("localhost", adminPortAsString, proxyPortAsString)
		})

		AfterEach(func() {
			exec.Command(hoverctlBinary, "stop", "-v").Run()
		})

		Context("I can control a process of hoverfly", func() {

			It("by starting hoverfly", func() {
				setOutput, _ := exec.Command(hoverctlBinary, "start", "-v").Output()

				output := strings.TrimSpace(string(setOutput))
				Expect(output).To(ContainSubstring("Hoverfly is now running"))
				Expect(output).To(ContainSubstring("admin-port=" + adminPortAsString))
				Expect(output).To(ContainSubstring("proxy-port=" + proxyPortAsString))

				data, err := ioutil.ReadFile("./.hoverfly/hoverfly." + adminPortAsString + "." + proxyPortAsString + ".pid")

				if err != nil {
					Fail("Could not find pid file")
				}

				Expect(data).ToNot(BeEmpty())

				if _, err := strconv.Atoi(string(data)); err != nil {
					Fail("Pid file not have an integer in it")
				}
			})

			It("by stopping  hoverfly", func() {
				exec.Command(hoverctlBinary, "start").Run()

				setOutput, _ := exec.Command(hoverctlBinary, "stop").Output()

				output := strings.TrimSpace(string(setOutput))
				Expect(output).To(ContainSubstring("Hoverfly has been stopped"))

				_, err := ioutil.ReadFile("./.hoverfly/hoverfly." + adminPortAsString + "." + proxyPortAsString + ".pid")

				if err == nil {
					Fail("Found the pid file that should have been deleted")
				}
			})

			It("by starting and stopping hoverfly on a different admin port using a flag", func() {
				setOutput, _ := exec.Command(hoverctlBinary, "start", "--admin-port=11223").Output()

				output := strings.TrimSpace(string(setOutput))
				Expect(output).To(ContainSubstring("Hoverfly is now running"))
				Expect(output).To(ContainSubstring("admin-port=11223"))
				Expect(output).To(ContainSubstring("proxy-port=" + proxyPortAsString))

				data, err := ioutil.ReadFile("./.hoverfly/hoverfly.11223." + proxyPortAsString + ".pid")

				if err != nil {
					Fail("Could not find pid file")
				}

				Expect(data).ToNot(BeEmpty())

				if _, err := strconv.Atoi(string(data)); err != nil {
					Fail("Pid file not have an integer in it")
				}

				setOutput, _ = exec.Command(hoverctlBinary, "stop", "--admin-port=11223").Output()

				output = strings.TrimSpace(string(setOutput))
				Expect(output).To(ContainSubstring("Hoverfly has been stopped"))

				_, err = ioutil.ReadFile("./.hoverfly/hoverfly.11223." + proxyPortAsString + ".pid")

				if err == nil {
					Fail("Found the pid file that should have been deleted")
				}

			})

			It("by starting and stopping hoverfly on a different proxy port using a flag", func() {
				setOutput, _ := exec.Command(hoverctlBinary, "start", "--proxy-port=22113").Output()

				output := strings.TrimSpace(string(setOutput))
				Expect(output).To(ContainSubstring("Hoverfly is now running"))
				Expect(output).To(ContainSubstring("admin-port=" + adminPortAsString))
				Expect(output).To(ContainSubstring("proxy-port=22113"))

				data, err := ioutil.ReadFile("./.hoverfly/hoverfly." + adminPortAsString + ".22113.pid")

				if err != nil {
					Fail("Could not find pid file")
				}

				Expect(data).ToNot(BeEmpty())

				if _, err := strconv.Atoi(string(data)); err != nil {
					Fail("Pid file not have an integer in it")
				}

				setOutput, _ = exec.Command(hoverctlBinary, "stop", "--proxy-port=22113").Output()

				output = strings.TrimSpace(string(setOutput))
				Expect(output).To(ContainSubstring("Hoverfly has been stopped"))

				_, err = ioutil.ReadFile("./.hoverfly/hoverfly." + adminPortAsString + ".22113.pid")

				if err == nil {
					Fail("Found the pid file that should have been deleted")
				}
			})

			It("by starting and stopping hoverfly on a different admin and proxy port using both flag", func() {
				setOutput, _ := exec.Command(hoverctlBinary, "start", "--admin-port=11223", "--proxy-port=22113").Output()

				output := strings.TrimSpace(string(setOutput))
				Expect(output).To(ContainSubstring("Hoverfly is now running"))
				Expect(output).To(ContainSubstring("admin-port=11223"))
				Expect(output).To(ContainSubstring("proxy-port=22113"))

				data, err := ioutil.ReadFile("./.hoverfly/hoverfly.11223.22113.pid")

				if err != nil {
					Fail("Could not find pid file")
				}

				Expect(data).ToNot(BeEmpty())

				if _, err := strconv.Atoi(string(data)); err != nil {
					Fail("Pid file not have an integer in it")
				}

				setOutput, _ = exec.Command(hoverctlBinary, "stop", "--admin-port=11223", "--proxy-port=22113").Output()

				output = strings.TrimSpace(string(setOutput))
				Expect(output).To(ContainSubstring("Hoverfly has been stopped"))

				_, err = ioutil.ReadFile("./.hoverfly/hoverfly.11223.22113.pid")

				if err == nil {
					Fail("Found the pid file that should have been deleted")
				}
			})

			It("but you cannot start hoverfly if already running", func() {
				exec.Command(hoverctlBinary, "start").Run()

				setOutput, _ := exec.Command(hoverctlBinary, "start").Output()

				output := strings.TrimSpace(string(setOutput))
				Expect(output).To(ContainSubstring("Hoverfly is already running"))
			})

			It("but you cannot stop hoverfly if is not running", func() {
				setOutput, _ := exec.Command(hoverctlBinary, "stop").Output()

				output := strings.TrimSpace(string(setOutput))
				Expect(output).To(ContainSubstring("Hoverfly is not running"))
			})

		})

		Context("I can control a process of hoverfly running as a webserver", func() {

			It("by starting hoverfly as a webserver", func() {
				setOutput, _ := exec.Command(hoverctlBinary, "start", "webserver", "-v").Output()

				output := strings.TrimSpace(string(setOutput))
				Expect(output).To(ContainSubstring("Hoverfly is now running as a webserver"))
				Expect(output).To(ContainSubstring("admin-port=" + adminPortAsString))
				Expect(output).To(ContainSubstring("webserver-port=" + proxyPortAsString))

				data, err := ioutil.ReadFile("./.hoverfly/hoverfly." + adminPortAsString + "." + proxyPortAsString + ".pid")

				if err != nil {
					Fail("Could not find pid file")
				}

				Expect(data).ToNot(BeEmpty())

				if _, err := strconv.Atoi(string(data)); err != nil {
					Fail("Pid file not have an integer in it")
				}

				request := sling.New().Get("http://localhost:" + proxyPortAsString)
				response := functional_tests.DoRequest(request)

				responseBody, err := ioutil.ReadAll(response.Body)
				Expect(err).To(BeNil())

				Expect(string(responseBody)).ToNot(ContainSubstring("This is a proxy server"))
			})
		})

		Context("You can specify the certificate and key for hoverfly", func() {

			It("starts hoverfly with different certificate and key", func() {
				setOutput, _ := exec.Command(hoverctlBinary, "start", "--certificate", "testdata/cert.pem", "--key", "testdata/key.pem", "-v").Output()

				output := strings.TrimSpace(string(setOutput))
				Expect(output).To(ContainSubstring("Hoverfly is now running"))

				data, err := ioutil.ReadFile("./.hoverfly/hoverfly." + adminPortAsString + "." + proxyPortAsString + ".pid")

				if err != nil {
					Fail("Could not find pid file")
				}

				Expect(data).ToNot(BeEmpty())

				data, err = ioutil.ReadFile("./.hoverfly/hoverfly." + adminPortAsString + "." + proxyPortAsString + ".log")

				if err != nil {
					Fail("Could not find log file")
				}

				Expect(data).ToNot(BeEmpty())
				Expect(data).To(ContainSubstring("Default keys have been overwritten"))
			})
		})

		Context("You can disable tls for hoverfly", func() {

			It("starts hoverfly with tls verification turned off", func() {
				setOutput, _ := exec.Command(hoverctlBinary, "start", "--disable-tls", "-v").Output()

				output := strings.TrimSpace(string(setOutput))
				Expect(output).To(ContainSubstring("Hoverfly is now running"))

				data, err := ioutil.ReadFile("./.hoverfly/hoverfly." + adminPortAsString + "." + proxyPortAsString + ".pid")

				if err != nil {
					Fail("Could not find pid file")
				}

				Expect(data).ToNot(BeEmpty())

				data, err = ioutil.ReadFile("./.hoverfly/hoverfly." + adminPortAsString + "." + proxyPortAsString + ".log")

				if err != nil {
					Fail("Could not find log file")
				}

				Expect(data).ToNot(BeEmpty())
				Expect(data).To(ContainSubstring("tls certificate verification is now turned off!"))
			})
		})

		Context("You can start a hoverfly based on config from config.yml", func() {

			It("will start on the admin and proxy ports", func() {
				WriteConfiguration("localhost", "5543", "6478")
				setOutput, _ := exec.Command(hoverctlBinary, "start", "-v").Output()

				output := strings.TrimSpace(string(setOutput))
				Expect(output).To(ContainSubstring("hoverfly -ap=5543 -pp=6478"))
			})

			It("will start as a webserver", func() {
				WriteConfigurationWithAuth("localhost", "7654", "8765", true, "", "")
				setOutput, _ := exec.Command(hoverctlBinary, "start", "-v").Output()

				output := strings.TrimSpace(string(setOutput))
				Expect(output).To(ContainSubstring("hoverfly -ap=7654 -pp=8765 -db=memory -webserver"))
				Expect(output).To(ContainSubstring("Hoverfly is now running as a webserver"))
			})

		})

		Context("You can set db options for hoverfly", func() {

			It("starts hoverfly with boltdb for data persistence", func() {
				setOutput, _ := exec.Command(hoverctlBinary, "start", "--database", "boltdb", "-v").Output()

				output := strings.TrimSpace(string(setOutput))
				Expect(output).To(ContainSubstring("Hoverfly is now running"))

				data, err := ioutil.ReadFile("./.hoverfly/hoverfly." + adminPortAsString + "." + proxyPortAsString + ".pid")

				if err != nil {
					Fail("Could not find pid file")
				}

				Expect(data).ToNot(BeEmpty())

				data, err = ioutil.ReadFile("./.hoverfly/hoverfly." + adminPortAsString + "." + proxyPortAsString + ".log")

				if err != nil {
					Fail("Could not find log file")
				}

				Expect(data).ToNot(BeEmpty())
				Expect(data).To(ContainSubstring("Creating bolt db backend."))
			})
		})
	})
})

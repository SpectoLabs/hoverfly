package hoverfly_end_to_end_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os/exec"
	"strings"
	"io/ioutil"
	"strconv"
	"github.com/phayes/freeport"
)


var _ = Describe("When I use hoverctl", func() {

	Describe("without a running hoverfly", func() {

		var (
			adminPort = freeport.GetPort()
			adminPortAsString = strconv.Itoa(adminPort)

			proxyPort = freeport.GetPort()
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

				data, err := ioutil.ReadFile("./.hoverfly/hoverfly.11223." + proxyPortAsString + ".pid")

				if err != nil {
					Fail("Could not find pid file")
				}

				Expect(data).ToNot(BeEmpty())

				if _, err := strconv.Atoi(string(data)); err != nil {
					Fail("Pid file not have an integer in it")
				}

				GetHoverflyMode(11223)

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

				data, err := ioutil.ReadFile("./.hoverfly/hoverfly." + adminPortAsString + ".22113.pid")

				if err != nil {
					Fail("Could not find pid file")
				}

				Expect(data).ToNot(BeEmpty())

				if _, err := strconv.Atoi(string(data)); err != nil {
					Fail("Pid file not have an integer in it")
				}

				GetHoverflyMode(adminPort)

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

				data, err := ioutil.ReadFile("./.hoverfly/hoverfly.11223.22113.pid")

				if err != nil {
					Fail("Could not find pid file")
				}

				Expect(data).ToNot(BeEmpty())

				if _, err := strconv.Atoi(string(data)); err != nil {
					Fail("Pid file not have an integer in it")
				}

				GetHoverflyMode(11223)

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
	})
})
package hoverfly_end_to_end_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os/exec"
	"strings"
	"strconv"
	"github.com/phayes/freeport"
	"io/ioutil"
	"path/filepath"
	"os"
)

var _ = Describe("When I use hoverctl with a running an authenticated hoverfly", func() {
	var (
		hoverflyCmd *exec.Cmd

		adminPort = freeport.GetPort()
		adminPortAsString = strconv.Itoa(adminPort)

		proxyPort = freeport.GetPort()
		proxyPortAsString = strconv.Itoa(proxyPort)


		username = "ft_user"
		password = "ft_password"
	)

	Describe("and the credentials are in the hoverctl config", func() {

		BeforeEach(func() {
			hoverflyCmd = startHoverflyWithAuth(adminPort, proxyPort, workingDirectory, username, password)
			WriteConfigurationWithAuth("localhost", adminPortAsString, proxyPortAsString, username, password)
		})

		AfterEach(func() {
			hoverflyCmd.Process.Kill()
		})

		Context("you can get the mode", func() {

			It("and it returns the correct mode", func() {
				out, _ := exec.Command(hoverctlBinary, "mode").Output()

				output := strings.TrimSpace(string(out))
				Expect(output).To(ContainSubstring("Hoverfly is set to simulate mode"))
			})
		})

		Context("you can set the mode", func() {

			It("and it correctly sets it", func() {
				setOutput, _ := exec.Command(hoverctlBinary, "mode", "capture").Output()

				output := strings.TrimSpace(string(setOutput))
				Expect(output).To(ContainSubstring("Hoverfly has been set to capture mode"))

				getOutput, _ := exec.Command(hoverctlBinary, "mode").Output()

				output = strings.TrimSpace(string(getOutput))
				Expect(output).To(ContainSubstring("Hoverfly is set to capture mode"))
			})
		})

		Context("you can manage simulations", func() {
			workingDirectory, _ := os.Getwd()
			fileToWrite := filepath.Join(workingDirectory, "/.hoverfly/cache/benjih.test.latest.json")
			ioutil.WriteFile(fileToWrite,
				[]byte(`
					{
						"data": [{
							"request": {
								"path": "/api/bookings",
								"method": "POST",
								"destination": "www.my-test.com",
								"scheme": "http",
								"query": "",
								"body": "{\"flightId\": \"1\"}",
								"headers": {
									"Content-Type": [
										"application/json"
									]
								}
							},
							"response": {
								"status": 201,
								"body": "",
								"encodedBody": false,
								"headers": {
									"Location": [
										"http://localhost/api/bookings/1"
									]
								}
							}
						}]
					}`), 0644)

			It("by importing and exporting data", func() {
				setOutput, _ := exec.Command(hoverctlBinary, "import", "benjih/test:latest").Output()

				output := strings.TrimSpace(string(setOutput))
				Expect(output).To(ContainSubstring("benjih/test:latest imported successfully"))

				setOutput, _ = exec.Command(hoverctlBinary, "export", "benjih/test-copy:latest").Output()

				output = strings.TrimSpace(string(setOutput))
				Expect(output).To(ContainSubstring("benjih/test-copy:latest exported successfully"))


				exportFile, err := ioutil.ReadFile(workingDirectory + "/.hoverfly/cache/benjih.test-copy.latest.json")
				if err != nil {
					Fail("Failed reading test data")
				}

				Expect(string(exportFile)).To(ContainSubstring(`"path":"/api/bookings"`))
				Expect(string(exportFile)).To(ContainSubstring(`"body":"{\"flightId\": \"1\"}"`))
			})

			It("and then wiping hoverfly", func() {
				setOutput, _ := exec.Command(hoverctlBinary, "wipe").Output()

				output := strings.TrimSpace(string(setOutput))
				Expect(output).To(ContainSubstring("Hoverfly has been wiped"))
			})
		})
	})

	Describe("and the credentials are not the hoverctl config", func() {

		workingDirectory, _ := os.Getwd()
		fileToWrite := filepath.Join(workingDirectory, "/.hoverfly/cache/benjih.test.latest.json")
		ioutil.WriteFile(fileToWrite,
			[]byte(`
					{
						"data": [{
							"request": {
								"path": "/api/bookings",
								"method": "POST",
								"destination": "www.my-test.com",
								"scheme": "http",
								"query": "",
								"body": "{\"flightId\": \"1\"}",
								"headers": {
									"Content-Type": [
										"application/json"
									]
								}
							},
							"response": {
								"status": 201,
								"body": "",
								"encodedBody": false,
								"headers": {
									"Location": [
										"http://localhost/api/bookings/1"
									]
								}
							}
						}]
					}`), 0644)

		BeforeEach(func() {
			hoverflyCmd = startHoverflyWithAuth(adminPort, proxyPort, workingDirectory, username, password)
			WriteConfiguration("localhost", adminPortAsString, proxyPortAsString)
		})

		AfterEach(func() {
			hoverflyCmd.Process.Kill()
		})

		Context("you cannot get the mode", func() {

			It("and it returns an error", func() {
				out, _ := exec.Command(hoverctlBinary, "mode").Output()

				output := strings.TrimSpace(string(out))
				Expect(output).To(ContainSubstring("Hoverfly requires authentication"))
			})
		})

		Context("you cannot set the mode", func() {

			It("and it returns an error", func() {
				setOutput, _ := exec.Command(hoverctlBinary, "mode", "capture").Output()

				output := strings.TrimSpace(string(setOutput))
				Expect(output).To(ContainSubstring("Hoverfly requires authentication"))
			})
		})

		Context("you cannot manage simulations", func() {

			It("by importing data", func() {
				setOutput, _ := exec.Command(hoverctlBinary, "import", "benjih/test:latest").Output()

				output := strings.TrimSpace(string(setOutput))
				Expect(output).To(ContainSubstring("Hoverfly requires authentication"))
			})

			It("and then exporting the data", func() {
				setOutput, _ := exec.Command(hoverctlBinary, "export", "benjih/test:latest").Output()

				output := strings.TrimSpace(string(setOutput))
				Expect(output).To(ContainSubstring("Hoverfly requires authentication"))
			})

			It("and then wiping hoverfly", func() {
				setOutput, _ := exec.Command(hoverctlBinary, "wipe").Output()

				output := strings.TrimSpace(string(setOutput))
				Expect(output).To(ContainSubstring("Hoverfly requires authentication"))
			})
		})
	})
})
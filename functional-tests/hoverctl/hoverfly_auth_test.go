package hoverctl_end_to_end

import (
	"io/ioutil"
	"os/exec"
	"strings"

	"github.com/SpectoLabs/hoverfly/functional-tests"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("When I use hoverctl with a running an authenticated hoverfly", func() {
	var (
		username = "ft_user"
		password = "ft_password"
	)

	Describe("and the credentials are in the hoverctl config", func() {

		BeforeEach(func() {
			hoverfly = functional_tests.NewHoverfly()
			hoverfly.Start("-db", "boltdb", "-add", "-username", username, "-password")
			hoverfly.Stop()
			hoverfly.Start("-db", "boltdb", "-auth", "true")

			WriteConfigurationWithAuth("localhost", hoverfly.GetAdminPort(), hoverfly.GetProxyPort(), false, username, password)
		})

		AfterEach(func() {
			hoverfly.Stop()
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

			It("by importing and exporting data", func() {
				filePath := generateFileName()
				ioutil.WriteFile(filePath,
					[]byte(`
						{
							"data": {
								"pairs": [{
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
							},
							"meta": {
								"schemaVersion": "v1"
							}
						}`), 0644)
				setOutput, _ := exec.Command(hoverctlBinary, "import", filePath).Output()

				output := strings.TrimSpace(string(setOutput))
				Expect(output).To(ContainSubstring("Successfully imported from " + filePath))

				newFilePath := generateFileName()
				setOutput, _ = exec.Command(hoverctlBinary, "export", newFilePath).Output()

				output = strings.TrimSpace(string(setOutput))
				Expect(output).To(ContainSubstring("Successfully exported to " + newFilePath))

				exportFile, err := ioutil.ReadFile(newFilePath)
				if err != nil {
					Fail("Failed reading test data")
				}

				Expect(string(exportFile)).To(ContainSubstring(`"path": "/api/bookings"`))
				Expect(string(exportFile)).To(ContainSubstring(`"body": "{\"flightId\": \"1\"}"`))
			})

			It("and then delete simulations from hoverfly", func() {
				setOutput, _ := exec.Command(hoverctlBinary, "delete", "simulations").Output()

				output := strings.TrimSpace(string(setOutput))
				Expect(output).To(ContainSubstring("Simulations have been deleted from Hoverfly"))
			})
		})
	})

	Describe("and the credentials are not in the hoverctl config", func() {

		BeforeEach(func() {
			hoverfly = functional_tests.NewHoverfly()
			hoverfly.Start("-db", "boltdb", "-add", "-username", username, "-password", password)
			hoverfly.Stop()
			hoverfly.Start("-db", "boltdb", "-auth", "true")

			WriteConfiguration("localhost", hoverfly.GetAdminPort(), hoverfly.GetProxyPort())
		})

		AfterEach(func() {
			hoverfly.Stop()
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

			var filePath string

			BeforeEach(func() {
				filePath = generateFileName()
				ioutil.WriteFile(filePath,
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
			})

			It("by importing data", func() {
				setOutput, _ := exec.Command(hoverctlBinary, "import", filePath).Output()

				output := strings.TrimSpace(string(setOutput))
				Expect(output).To(ContainSubstring("Hoverfly requires authentication"))
			})

			It("and then exporting the data", func() {
				setOutput, _ := exec.Command(hoverctlBinary, "export", filePath).Output()

				output := strings.TrimSpace(string(setOutput))
				Expect(output).To(ContainSubstring("Hoverfly requires authentication"))
			})

			It("and then wiping hoverfly", func() {
				setOutput, _ := exec.Command(hoverctlBinary, "delete", "simulations").Output()

				output := strings.TrimSpace(string(setOutput))
				Expect(output).To(ContainSubstring("Hoverfly requires authentication"))
			})
		})
	})
})

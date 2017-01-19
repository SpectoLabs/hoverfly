package hoverctl_end_to_end

import (
	"fmt"
	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/phayes/freeport"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var _ = Describe("When I use hoverctl", func() {
	var (
		hoverflyCmd *exec.Cmd

		workingDir, _     = os.Getwd()
		adminPort         = freeport.GetPort()
		adminPortAsString = strconv.Itoa(adminPort)

		proxyPort = freeport.GetPort()

		hoverflyData = `
					{
						"data": [{
							"request": {
								"requestType": "snapshot",
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
					}`
	)

	Describe("with a running hoverfly", func() {

		BeforeEach(func() {
			hoverflyCmd = startHoverfly(adminPort, proxyPort, workingDir)
		})

		AfterEach(func() {
			hoverflyCmd.Process.Kill()
		})

		Describe("Managing Hoverflies data using the CLI", func() {

			BeforeEach(func() {
				DoRequest(sling.New().Post(fmt.Sprintf("http://localhost:%v/api/records", adminPort)).Body(strings.NewReader(hoverflyData)))

				resp := DoRequest(sling.New().Get(fmt.Sprintf("http://localhost:%v/api/records", adminPort)))
				bytes, _ := ioutil.ReadAll(resp.Body)
				Expect(string(bytes)).ToNot(Equal(`{"data":null}`))
			})

			It("can export", func() {

				// Export the data
				output, _ := exec.Command(hoverctlBinary, "export", "testuser1/simulation1", "--admin-port="+adminPortAsString).Output()

				Expect(output).To(ContainSubstring("testuser1/simulation1:latest exported successfully"))
				Expect(ioutil.ReadFile(hoverctlCacheDir + "/testuser1.simulation1.latest.json")).To(MatchJSON(hoverflyData))

			})

			It("can import", func() {

				err := ioutil.WriteFile(hoverctlCacheDir+"/testuser2.simulation2.latest.json", []byte(hoverflyData), 0644)
				Expect(err).To(BeNil())

				output, _ := exec.Command(hoverctlBinary, "import", "testuser2/simulation2", "--admin-port="+adminPortAsString).Output()

				Expect(output).To(ContainSubstring("testuser2/simulation2:latest imported successfully"))

				resp := DoRequest(sling.New().Get(fmt.Sprintf("http://localhost:%v/api/records", adminPort)))
				bytes, _ := ioutil.ReadAll(resp.Body)
				Expect(string(bytes)).To(MatchJSON(hoverflyData))
			})

		})
	})
})

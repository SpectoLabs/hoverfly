package hoverfly_end_to_end_test

import (
	. "github.com/onsi/ginkgo"
	"os/exec"
	"os"
	"strings"
	"github.com/phayes/freeport"
	"fmt"
	"github.com/dghubble/sling"
)

var _ = Describe("When I use hoverfly-cli", func() {
	var (
		hoverflyCmd *exec.Cmd

		workingDir, _ = os.Getwd()
		adminPort = freeport.GetPort()
		//adminPortAsString = strconv.Itoa(adminPort)

		proxyPort = freeport.GetPort()
	)

	Describe("with a running hoverfly", func() {

		BeforeEach(func() {
			hoverflyCmd = startHoverfly(adminPort, proxyPort, workingDir)
		})

		AfterEach(func() {
			hoverflyCmd.Process.Kill()
		})

		Describe("which contains some data", func() {

			BeforeEach(func() {
				sling.New().Post(fmt.Sprintf("http://localhost:%v/api/state", adminPort)).Body(strings.NewReader(`
					{
						"data": [{
							"request": {
								"path": "/api/bookings",
								"method": "POST",
								"destination": "www.my-test.com",
								"query": null,
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
					}`))
			})

			//Context("and the data is wiped", func() {
			//	output, _ := exec.Command(hoverctlBinary, "wipe").CombinedOutput()
			//	panic(string(output))
			//})
		})
	})
})
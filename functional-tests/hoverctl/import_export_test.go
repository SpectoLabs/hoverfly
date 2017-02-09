package hoverctl_end_to_end

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"net/http"
	"net/http/httptest"

	"github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("When I use hoverctl", func() {

	var (
		hoverfly *functional_tests.Hoverfly
	)

	var (
		v1HoverflyData = `
					{
						"data": [{
							"request": {
								"requestType": "recording",
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

		v2HoverflyData = `
			{
				"data": {
					"pairs": [{
						"response": {
							"status": 201,
							"body": "",
							"encodedBody": false,
							"headers": {
								"Location": ["http://localhost/api/bookings/1"]
							}
						},
						"request": {
							"requestType": "recording",
							"path": "/api/bookings",
							"method": "POST",
							"destination": "www.my-test.com",
							"scheme": "http",
							"query": "",
							"body": "{\"flightId\": \"1\"}",
							"headers": {
								"Content-Type": ["application/json"]
							}
						}
					}],
					"globalActions": {
						"delays": []
					}
				},
				"meta": {
					"schemaVersion": "v1",
					"hoverflyVersion": "v0.9.2",
					"timeExported": "2016-11-10T12:27:46Z"
				}
			}`

		v2HoverflySimulation = `"pairs":[{"response":{"status":201,"body":"","encodedBody":false,"headers":{"Location":["http://localhost/api/bookings/1"]}},"request":{"requestType":"recording","path":"/api/bookings","method":"POST","destination":"www.my-test.com","scheme":"http","query":"","body":"{\"flightId\": \"1\"}","headers":{"Content-Type":["application/json"]}}}],"globalActions":{"delays":[]}}`

		v2HoverflyMeta = `"meta":{"schemaVersion":"v1","hoverflyVersion":"v\d+.\d+.\d+","timeExported":`
	)

	Describe("with a running hoverfly", func() {

		BeforeEach(func() {
			hoverfly = functional_tests.NewHoverfly()
			hoverfly.Start()

			WriteConfiguration("localhost", hoverfly.GetAdminPort(), hoverfly.GetProxyPort())
		})

		AfterEach(func() {
			hoverfly.Stop()
		})

		Describe("Managing Hoverflies data using the CLI", func() {

			BeforeEach(func() {
				functional_tests.DoRequest(sling.New().Post(fmt.Sprintf("http://localhost:%v/api/records", hoverfly.GetAdminPort())).Body(strings.NewReader(v1HoverflyData)))

				resp := functional_tests.DoRequest(sling.New().Get(fmt.Sprintf("http://localhost:%v/api/records", hoverfly.GetAdminPort())))
				bytes, _ := ioutil.ReadAll(resp.Body)
				Expect(string(bytes)).ToNot(Equal(`{"data":null}`))
			})

			It("can export", func() {

				fileName := generateFileName()
				// Export the data
				output := functional_tests.Run(hoverctlBinary, "export", fileName, "--admin-port="+hoverfly.GetAdminPort())

				Expect(output).To(ContainSubstring("Successfully exported simulation to " + fileName))

				data, err := ioutil.ReadFile(fileName)
				Expect(err).To(BeNil())

				buffer := new(bytes.Buffer)
				json.Compact(buffer, data)

				Expect(buffer.String()).To(ContainSubstring(v2HoverflySimulation))
				Expect(buffer.String()).To(MatchRegexp(v2HoverflyMeta))
			})

			It("can import", func() {

				fileName := generateFileName()
				err := ioutil.WriteFile(fileName, []byte(v2HoverflyData), 0644)
				Expect(err).To(BeNil())

				output := functional_tests.Run(hoverctlBinary, "import", fileName, "--admin-port="+hoverfly.GetAdminPort())

				Expect(output).To(ContainSubstring("Successfully imported simulation from " + fileName))

				resp := functional_tests.DoRequest(sling.New().Get(fmt.Sprintf("http://localhost:%v/api/records", hoverfly.GetAdminPort())))
				bytes, _ := ioutil.ReadAll(resp.Body)
				Expect(string(bytes)).To(MatchJSON(v1HoverflyData))
			})

			It("can import over http", func() {
				ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					fmt.Fprintln(w, v2HoverflyData)
				}))
				defer ts.Close()

				output := functional_tests.Run(hoverctlBinary, "import", ts.URL, "--admin-port="+hoverfly.GetAdminPort())

				Expect(output).To(ContainSubstring("Successfully imported simulation from " + ts.URL))

				resp := functional_tests.DoRequest(sling.New().Get(fmt.Sprintf("http://localhost:%v/api/records", hoverfly.GetAdminPort())))
				bytes, _ := ioutil.ReadAll(resp.Body)
				Expect(string(bytes)).To(MatchJSON(v1HoverflyData))
			})

			It("can import v1 simulations", func() {

				fileName := generateFileName()
				err := ioutil.WriteFile(fileName, []byte(v1HoverflyData), 0644)
				Expect(err).To(BeNil())

				output := functional_tests.Run(hoverctlBinary, "import", "--v1", fileName, "--admin-port="+hoverfly.GetAdminPort())

				Expect(output).To(ContainSubstring("Successfully imported simulation from " + fileName))

				resp := functional_tests.DoRequest(sling.New().Get(fmt.Sprintf("http://localhost:%v/api/records", hoverfly.GetAdminPort())))
				bytes, _ := ioutil.ReadAll(resp.Body)
				Expect(string(bytes)).To(MatchJSON(v1HoverflyData))
			})

			It("cannot import incorrect json / missing meta", func() {

				fileName := generateFileName()
				err := ioutil.WriteFile(fileName, []byte(`
				{
					"data": {
						"pairs": [{
							"response": {
								"status": 201,
								"body": "",
								"encodedBody": false,
								"headers": {
									"Location": ["http://localhost/api/bookings/1"]
								}
							},
							"request": {
								"requestType": "recording",
								"path": "/api/bookings",
								"method": "POST",
								"destination": "www.my-test.com",
								"scheme": "http",
								"query": "",
								"body": "{\"flightId\": \"1\"}",
								"headers": {
									"Content-Type": ["application/json"]
								}
							}
						}],
						"globalActions": {
							"delays": []
						}
					}
				}`), 0644)
				Expect(err).To(BeNil())

				output := functional_tests.Run(hoverctlBinary, "import", fileName, "--admin-port="+hoverfly.GetAdminPort())

				Expect(output).To(ContainSubstring("Import to Hoverfly failed: Json did not match schema: Object->Key[meta].Value->Object"))

				resp := functional_tests.DoRequest(sling.New().Get(fmt.Sprintf("http://localhost:%v/api/records", hoverfly.GetAdminPort())))
				bytes, _ := ioutil.ReadAll(resp.Body)
				Expect(string(bytes)).To(MatchJSON(v1HoverflyData))
			})
		})
	})
})

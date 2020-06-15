package hoverctl_suite

import (
	"bytes"
	"encoding/json"
	"fmt"
	v2 "github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"io/ioutil"
	"os"

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
		hoverflyData = `
			{
				"data": {
					"pairs": [{
						"response": {
							"status": 201,
							"body": "",
							"encodedBody": false,
							"headers": {
								"Location": ["http://localhost/api/bookings/1"]
							},
							"templated": false,
                            "fixedDelay": 0
						},
						"request": {
							"path": [
								{
									"matcher": "exact",
									"value": "/api/bookings"
								}
							],
							"method": [
								{
									"matcher": "exact",
									"value": "POST"
								}
							],
							"destination": [
								{
									"matcher": "exact",
									"value": "www.my-test.com"
								}
							],
							"scheme": [
								{
									"matcher": "exact",
									"value": "http"
								}
							],
							"body": [
								{
									"matcher": "exact",
									"value": "{\"flightId\": \"1\"}"
								}
							],
							"headers": {
								"Content-Type": [
									{
										"matcher": "exact",
										"value": "application/json"
									}
								]
							}
						}
					}],
					"globalActions": {
						"delays": []
					}
				},
				"meta": {
					"schemaVersion": "v5.1",
					"hoverflyVersion": "v0.9.2",
					"timeExported": "2016-11-10T12:27:46Z"
				}
			}`

		hoverflyDataWithMultiplePairs = `
			{
				"data": {
					"pairs": [{
						"response": {
							"status": 201,
							"body": "",
							"encodedBody": false,
							"headers": {
								"Location": ["http://localhost/api/bookings/1"]
							},
							"templated": false
						},
						"request": {
							"path": [
								{
									"matcher": "exact",
									"value": "/api/bookings"
								}
							],
							"method": [
								{
									"matcher": "exact",
									"value": "POST"
								}
							],
							"destination": [
								{
									"matcher": "exact",
									"value": "www.my-test.com"
								}
							],
							"scheme": [
								{
									"matcher": "exact",
									"value": "http"
								}
							],
							"body": [
								{
									"matcher": "exact",
									"value": "{\"flightId\": \"1\"}"
								}
							],
							"headers": {
								"Content-Type": [
									{
										"matcher": "exact",
										"value": "application/json"
									}
								]
							}
						}
					}, {
						"response": {
							"status": 201,
							"body": "",
							"encodedBody": false,
							"headers": {
								"Location": ["http://localhost/api/bookings/1"]
							},
							"templated": false
						},
						"request": {
							"path": [
								{
									"matcher": "exact",
									"value": "/api/bookings"
								}
							],
							"method": [
								{
									"matcher": "exact",
									"value": "POST"
								}
							],
							"destination": [
								{
									"matcher": "exact",
									"value": "www.other-test.com"
								}
							]
						}
					}],
					"globalActions": {
						"delays": []
					}
				},
				"meta": {
					"schemaVersion": "v5.1",
					"hoverflyVersion": "v0.9.2",
					"timeExported": "2016-11-10T12:27:46Z"
				}
			}`

		hoverflySimulation = `"pairs":[{"request":{"path":[{"matcher":"exact","value":"/api/bookings"}],"method":[{"matcher":"exact","value":"POST"}],"destination":[{"matcher":"exact","value":"www.my-test.com"}],"scheme":[{"matcher":"exact","value":"http"}],"body":[{"matcher":"exact","value":"{\"flightId\": \"1\"}"}],"headers":{"Content-Type":[{"matcher":"exact","value":"application/json"}]}},"response":{"status":201,"body":"","bodyFile":"","encodedBody":false,"headers":{"Location":["http://localhost/api/bookings/1"]},"templated":false,"fixedDelay":0}}],"globalActions":{"delays":[],"delaysLogNormal":[]}}`

		hoverflyMeta = `"meta":{"schemaVersion":"v5.1","hoverflyVersion":"v\d+.\d+.\d+(-rc.\d)*","timeExported":`
	)

	Describe("with a running hoverfly", func() {

		BeforeEach(func() {
			hoverfly = functional_tests.NewHoverfly()
			hoverfly.Start()

			functional_tests.Run(hoverctlBinary, "targets", "update", "local", "--admin-port", hoverfly.GetAdminPort())
		})

		AfterEach(func() {
			hoverfly.Stop()
		})

		Describe("Exporting simulation with bodyFile", func() {
			It("can export bodyFile fields", func() {
				fileName := functional_tests.GenerateFileName()
				bodyFileName := functional_tests.GenerateFileName()
				bodyFileContent := `{"success": true}`

				err := ioutil.WriteFile(bodyFileName, []byte(bodyFileContent), 0644)
				Expect(err).To(BeNil())

				hoverfly.ImportSimulation(`{
	"data": {
		"pairs": [
			{
				"request": {
					"path": [
						{
							"matcher": "exact",
							"value": "/api/v1/booking"
						}
					]
				},
				"response": {
					"status": 200,
					"bodyFile": "`+bodyFileName+`"
				}
			}
		]
	},
	"meta": {
		"schemaVersion": "v6"
	}
}`)
				// remove bodyFile to be restored by export command later
				Expect(os.Remove(bodyFileName)).To(BeNil())

				output := functional_tests.Run(hoverctlBinary, "export", fileName)
				Expect(output).To(ContainSubstring("Successfully exported simulation to " + fileName))

				data, err := ioutil.ReadFile(fileName)
				Expect(err).To(BeNil())

				var view v2.SimulationViewV6
				functional_tests.Unmarshal(data, &view)

				Expect(view.DataViewV6.RequestResponsePairs[0].Response.BodyFile).To(Equal(bodyFileName))

				data, err = ioutil.ReadFile(bodyFileName)
				Expect(err).To(BeNil())
				Expect(string(data)).To(Equal(bodyFileContent))
			})

			It("resets body when bodyFile is provided", func() {
				fileName := functional_tests.GenerateFileName()
				bodyFileName := functional_tests.GenerateFileName()

				err := ioutil.WriteFile(bodyFileName, []byte(`{"success": true}`), 0644)
				Expect(err).To(BeNil())

				hoverfly.ImportSimulation(`{
	"data": {
		"pairs": [
			{
				"request": {
					"path": [
						{
							"matcher": "exact",
							"value": "/api/v1/booking"
						}
					]
				},
				"response": {
					"status": 200,
					"bodyFile": "`+bodyFileName+`",
					"body": "testing content"
				}
			}
		]
	},
	"meta": {
		"schemaVersion": "v6"
	}
}`)
				output := functional_tests.Run(hoverctlBinary, "export", fileName)
				Expect(output).To(ContainSubstring("Successfully exported simulation to " + fileName))

				data, err := ioutil.ReadFile(fileName)
				Expect(err).To(BeNil())

				var view v2.SimulationViewV6
				functional_tests.Unmarshal(data, &view)

				Expect(view.DataViewV6.RequestResponsePairs[0].Response.Body).To(BeEmpty())
				Expect(view.DataViewV6.RequestResponsePairs[0].Response.BodyFile).To(Equal(bodyFileName))
			})
		})

		Describe("Managing Hoverflies data using the CLI", func() {

			BeforeEach(func() {
				hoverfly.ImportSimulation(hoverflyData)
			})

			It("can export", func() {

				fileName := functional_tests.GenerateFileName()
				// Export the data
				output := functional_tests.Run(hoverctlBinary, "export", fileName)

				Expect(output).To(ContainSubstring("Successfully exported simulation to " + fileName))

				data, err := ioutil.ReadFile(fileName)
				Expect(err).To(BeNil())

				buffer := new(bytes.Buffer)
				json.Compact(buffer, data)

				Expect(buffer.String()).To(ContainSubstring(hoverflySimulation))
				Expect(buffer.String()).To(MatchRegexp(hoverflyMeta))
			})

			It("can export with url pattern", func() {

				hoverfly.ImportSimulation(hoverflyDataWithMultiplePairs)
				fileName := functional_tests.GenerateFileName()
				// Export the data
				output := functional_tests.Run(hoverctlBinary, "export", fileName, "--url-pattern=my-test.com")

				Expect(output).To(ContainSubstring("Successfully exported simulation to " + fileName))

				data, err := ioutil.ReadFile(fileName)
				Expect(err).To(BeNil())

				buffer := new(bytes.Buffer)
				json.Compact(buffer, data)

				Expect(buffer.String()).To(ContainSubstring(hoverflySimulation))
				Expect(buffer.String()).To(MatchRegexp(hoverflyMeta))
			})

			It("can import", func() {

				fileName := functional_tests.GenerateFileName()
				err := ioutil.WriteFile(fileName, []byte(hoverflyData), 0644)
				Expect(err).To(BeNil())

				output := functional_tests.Run(hoverctlBinary, "import", fileName)

				Expect(output).To(ContainSubstring("Successfully imported simulation from " + fileName))

				resp := functional_tests.DoRequest(sling.New().Get(fmt.Sprintf("http://localhost:%v/api/v2/simulation", hoverfly.GetAdminPort())))
				bytes, _ := ioutil.ReadAll(resp.Body)
				Expect(string(bytes)).To(ContainSubstring(hoverflySimulation))
				Expect(string(bytes)).To(MatchRegexp(hoverflyMeta))
			})

			It("can import over http", func() {
				ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					fmt.Fprintln(w, hoverflyData)
				}))
				defer ts.Close()

				output := functional_tests.Run(hoverctlBinary, "import", ts.URL)

				Expect(output).To(ContainSubstring("Successfully imported simulation from " + ts.URL))

				resp := functional_tests.DoRequest(sling.New().Get(fmt.Sprintf("http://localhost:%v/api/v2/simulation", hoverfly.GetAdminPort())))
				bytes, _ := ioutil.ReadAll(resp.Body)
				Expect(string(bytes)).To(ContainSubstring(hoverflySimulation))
				Expect(string(bytes)).To(MatchRegexp(hoverflyMeta))
			})

			// TODO: Fix this test
			// It("cannot import incorrect json / missing meta", func() {
			// 	hoverfly.ImportSimulation(v3HoverflyData)
			// 	fileName := generateFileName()
			// 	err := ioutil.WriteFile(fileName, []byte(`
			// 	{
			// 		"data": {
			// 			"pairs": [{
			// 				"response": {
			// 					"status": 201,
			// 					"body": "",
			// 					"encodedBody": false,
			// 					"headers": {
			// 						"Location": ["http://localhost/api/bookings/1"]
			// 					}
			// 				},
			// 				"request": {
			// 					"requestType": {
			// 						"exactMatch": recording"
			// 					},
			// 					"path": {
			// 						"exactMatch": "/api/bookings"
			// 					},
			// 					"method": {
			// 						"exactMatch": "POST"
			// 					},
			// 					"destination": {
			// 						"exactMatch": "www.my-test.com"
			// 					},
			// 					"scheme":  {
			// 						"exactMatch": "http"
			// 					},
			// 					"query": {
			// 						"exactMatch": ""
			// 					},
			// 					"body": {
			// 						"exactMatch": "{\"flightId\": \"1\"}"
			// 					},
			// 					"headers": {
			// 						"Content-Type": ["application/json"]
			// 					}
			// 				}
			// 			}],
			// 			"globalActions": {
			// 				"delays": []
			// 			}
			// 		}
			// 	}`), 0644)
			// 	Expect(err).To(BeNil())

			// 	output := functional_tests.Run(hoverctlBinary, "import", fileName, "--admin-port="+hoverfly.GetAdminPort())

			// 	Expect(output).To(ContainSubstring("Import to Hoverfly failed: Json did not match schema: Object->Key[meta].Value->Object"))

			// 	resp := functional_tests.DoRequest(sling.New().Get(fmt.Sprintf("http://localhost:%v/api/v2/simulation", hoverfly.GetAdminPort())))
			// 	bytes, _ := ioutil.ReadAll(resp.Body)
			// 	Expect(string(bytes)).To(ContainSubstring(v3HoverflySimulation))
			// 	Expect(string(bytes)).To(MatchRegexp(v3HoverflyMeta))
			// })
		})
	})
})

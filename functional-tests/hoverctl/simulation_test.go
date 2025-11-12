package hoverctl_suite

import (
	"io/ioutil"

	"github.com/SpectoLabs/hoverfly/functional-tests"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("When I add simulation with hoverctl", func() {

	var (
		hoverfly *functional_tests.Hoverfly
	)

	Context("without providing a path to write to", func() {

		It("it should fail nicely", func() {
			output := functional_tests.Run(hoverctlBinary, "simulation", "add")

			Expect(output).To(ContainSubstring("You have not provided a path to simulation"))
			Expect(output).To(ContainSubstring("Try hoverctl simulation add --help for more information"))
		})
	})

	Context("with a target that doesn't exist", func() {
		It("should error", func() {
			output := functional_tests.Run(hoverctlBinary, "simulation", "add", "--target", "test-target")

			Expect(output).To(ContainSubstring("test-target is not a target"))
			Expect(output).To(ContainSubstring("Run `hoverctl targets create test-target`"))
		})
	})

	Describe("with a running hoverfly", func() {

		BeforeEach(func() {
			hoverfly = functional_tests.NewHoverfly()
			hoverfly.Start()

			functional_tests.Run(hoverctlBinary, "targets", "update", "local", "--admin-port", hoverfly.GetAdminPort())
		})

		AfterEach(func() {
			hoverfly.Stop()
		})

		It("can import multiple simulations", func() {

			file1 := functional_tests.GenerateFileName()
			err := ioutil.WriteFile(file1, []byte(`{
				"data": {
					"pairs": [{
						"response": {
							"status": 200,
							"body": "5iNe8dxWH5Ca8pZqAfEHv3rgC0SsvKNLu6o3K",
							"encodedBody": false
						},
						"request": {
							"path": [{
								"matcher": "exact",
								"value": "/bar"
							}],
							"method": [{
								"matcher": "exact",
								"value": "GET"
							}],
							"body": [{
								"matcher": "exact",
								"value": ""
							}]
						}
					}],
					"globalActions": {
						"delays": []
					}
				},
				"meta": {
					"schemaVersion": "v5"
				}
			}`), 0644)
			Expect(err).To(BeNil())
			file2 := functional_tests.GenerateFileName()
			err = ioutil.WriteFile(file2, []byte(`{
				"data": {
					"pairs": [{
						"response": {
							"status": 200,
							"body": "5iNe8dxWH5Ca8pZqAfEHv3rgC0SsvKNLu6o3K",
							"encodedBody": false
						},
						"request": {
							"path": [{
								"matcher": "exact",
								"value": "/foo"
							}],
							"method": [{
								"matcher": "exact",
								"value": "GET"
							}],
							"body": [{
								"matcher": "exact",
								"value": ""
							}]
						}
					}],
					"globalActions": {
						"delays": []
					}
				},
				"meta": {
					"schemaVersion": "v5"
				}
			}`), 0644)
			Expect(err).To(BeNil())

			output := functional_tests.Run(hoverctlBinary, "simulation", "add", file1, file2)
			Expect(output).To(ContainSubstring("Successfully added simulation from " + file1))
			Expect(output).To(ContainSubstring("Successfully added simulation from " + file2))

		})
	})
})

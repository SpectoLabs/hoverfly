package hoverctl_suite

import (
	"io/ioutil"

	"github.com/SpectoLabs/hoverfly/functional-tests"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("When I import with hoverctl", func() {

	var (
		hoverfly *functional_tests.Hoverfly
	)

	Context("without providing a path to write to", func() {

		It("it should fail nicely", func() {
			output := functional_tests.Run(hoverctlBinary, "import")

			Expect(output).To(ContainSubstring("You have not provided a path to simulation"))
			Expect(output).To(ContainSubstring("Try hoverctl import --help for more information"))
		})
	})

	Context("with a target that doesn't exist", func() {
		It("should error", func() {
			output := functional_tests.Run(hoverctlBinary, "import", "--target", "test-target")

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

		It("can show warnings", func() {

			fileName := functional_tests.GenerateFileName()
			err := ioutil.WriteFile(fileName, []byte(`{
				"data": {
				  "pairs": [
					{
					  "response": {
						"status": 200,
						"body": "YmFzZTY0IGVuY29kZWQ=",
						"encodedBody": true,
						"headers": {
						  "Hoverfly": [
							"Was-Here"
						  ]
						},
						"templated": false
					  },
					  "request": {
						"path": {
						  "exactMatch": "/pages/keyconcepts/templates.html"
						},
						"method": {
						  "exactMatch": "GET"
						},
						"destination": {
						  "exactMatch": "docs.hoverfly.io"
						},
						"scheme": {
						  "exactMatch": "http"
						},
						"query": {
						  "exactMatch": "query=true"
						},
						"body": {
						  "exactMatch": ""
						}
					  }
					}
				  ],
				  "globalActions": {
					"delays": []
				  }
				},
				"meta": {
				  "schemaVersion": "v3",
				  "hoverflyVersion": "v0.13.0",
				  "timeExported": "2017-07-13T16:34:30+01:00"
				}
			  }`), 0644)
			Expect(err).To(BeNil())

			output := functional_tests.Run(hoverctlBinary, "import", fileName)
			Expect(output).To(ContainSubstring("WARNING: Usage of deprecated field `deprecatedQuery` on data.pairs[0].request.deprecatedQuery, please update your simulation to use `query` field"))
			Expect(output).To(ContainSubstring("https://hoverfly.readthedocs.io/en/latest/pages/troubleshooting/troubleshooting.html#why-does-my-simulation-have-a-deprecatedquery-field"))
			Expect(output).To(ContainSubstring("Successfully imported simulation "))

		})
	})
})

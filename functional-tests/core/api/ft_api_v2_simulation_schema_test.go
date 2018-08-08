package api_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"

	"github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("/api/v2/simulation/schema", func() {

	var (
		hoverfly *functional_tests.Hoverfly
	)

	BeforeEach(func() {
		hoverfly = functional_tests.NewHoverfly()
		hoverfly.Start()
	})

	AfterEach(func() {
		hoverfly.Stop()
	})

	Context("GET", func() {

		It("Should get the JSON schema", func() {
			req := sling.New().Get("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/simulation/schema")
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))

			fileBytes, err := ioutil.ReadFile("../../../schema.json")
			Expect(err).To(BeNil(), "schema.json not found")

			fileBuffer := new(bytes.Buffer)
			json.Compact(fileBuffer, fileBytes)

			responseJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())

			Expect(responseJson).To(Equal(fileBuffer.Bytes()))
		})
	})
})

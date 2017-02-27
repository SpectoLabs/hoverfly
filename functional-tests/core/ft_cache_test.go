package hoverfly_test

import (
	"io/ioutil"

	"github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/antonholmquist/jason"
	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Hoverfly cache", func() {

	var (
		hoverfly *functional_tests.Hoverfly
	)

	BeforeEach(func() {
		hoverfly = functional_tests.NewHoverfly()
		hoverfly.Start()
		hoverfly.ImportSimulation(functional_tests.JsonPayload)
		hoverfly.Proxy(sling.New().Get("http://template-server.com"))
	})

	AfterEach(func() {
		hoverfly.Stop()
	})

	It("should be flushed when changing to capture mode", func() {
		hoverfly.SetMode("capture")
		req := sling.New().Get("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/cache")
		res := functional_tests.DoRequest(req)
		Expect(res.StatusCode).To(Equal(200))
		responseJson, err := ioutil.ReadAll(res.Body)
		Expect(err).To(BeNil())

		jsonObject, err := jason.NewObjectFromBytes(responseJson)
		Expect(err).To(BeNil())

		cacheArray, err := jsonObject.GetObjectArray("cache")
		Expect(err).To(BeNil())

		Expect(cacheArray).To(HaveLen(0))
	})

	It("should be flushed when importing via the API", func() {
		hoverfly.ImportSimulation(functional_tests.JsonPayload)
		req := sling.New().Get("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/cache")
		res := functional_tests.DoRequest(req)
		Expect(res.StatusCode).To(Equal(200))
		responseJson, err := ioutil.ReadAll(res.Body)
		Expect(err).To(BeNil())

		jsonObject, err := jason.NewObjectFromBytes(responseJson)
		Expect(err).To(BeNil())

		cacheArray, err := jsonObject.GetObjectArray("cache")
		Expect(err).To(BeNil())

		Expect(cacheArray).To(HaveLen(0))
	})
})

package hoverfly_test

import (
	"io/ioutil"

	"github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/antonholmquist/jason"
	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("/api/v2/cache", func() {

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

	Context("GET", func() {

		It("should get request response pairs in cache", func() {
			req := sling.New().Get("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/cache")
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))
			responseJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())

			jsonObject, err := jason.NewObjectFromBytes(responseJson)
			Expect(err).To(BeNil())

			cacheArray, err := jsonObject.GetObjectArray("cache")
			Expect(err).To(BeNil())

			Expect(cacheArray).To(HaveLen(1))

			request, err := cacheArray[0].GetObject("request")

			Expect(request.GetString("body")).Should(Equal(""))
			Expect(request.GetString("destination")).Should(Equal("template-server.com"))
			Expect(request.GetString("method")).Should(Equal("GET"))
			Expect(request.GetString("path")).Should(Equal("/"))
			Expect(request.GetString("query")).Should(Equal(""))
			Expect(request.GetString("scheme")).Should(Equal("http"))

			requestHeaders, _ := request.GetObject("headers")
			Expect(requestHeaders.GetStringArray("Accept-Encoding")).Should(ContainElement("gzip"))
			Expect(requestHeaders.GetStringArray("User-Agent")).Should(ContainElement("Go-http-client/1.1"))

			response, err := cacheArray[0].GetObject("response")

			Expect(response.GetInt64("status")).Should(Equal(int64(200)))
			Expect(response.GetString("body")).Should(Equal("template match"))
			Expect(response.GetBoolean("encodedBody")).Should(BeFalse())
		})

		It("should get error when cache is disabled", func() {
			hoverfly.Stop()
			hoverfly.Start("-disable-cache")

			req := sling.New().Get("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/cache")
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(500))
			responseJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())

			jsonObject, err := jason.NewObjectFromBytes(responseJson)
			Expect(err).To(BeNil())

			Expect(jsonObject.GetString("error")).Should(Equal("No cache set"))
		})
	})

	Context("DELETE", func() {

		It("should flush cache", func() {
			req := sling.New().Delete("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/cache")
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

		It("should get error when cache is disabled", func() {
			hoverfly.Stop()
			hoverfly.Start("-disable-cache")

			req := sling.New().Get("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/cache")
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(500))
			responseJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())

			jsonObject, err := jason.NewObjectFromBytes(responseJson)
			Expect(err).To(BeNil())

			Expect(jsonObject.GetString("error")).Should(Equal("No cache set"))
		})
	})
})

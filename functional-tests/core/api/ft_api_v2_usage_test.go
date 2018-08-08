package api_test

import (
	"io/ioutil"

	"github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("/api/v2/hoverfly/usage", func() {

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

		It("Should get the usage counters", func() {
			req := sling.New().Get("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/hoverfly/usage")
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))
			modeJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(modeJson).To(Equal([]byte(`{"usage":{"counters":{"capture":0,"diff":0,"modify":0,"simulate":0,"spy":0,"synthesize":0}}}`)))
		})

		It("Should get the usage counters with 1 simulate request when a request has been made", func() {
			proxyReq := sling.New().Get("http://www.google.com")
			hoverfly.Proxy(proxyReq)
			req := sling.New().Get("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/hoverfly/usage")
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))
			modeJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(modeJson).To(Equal([]byte(`{"usage":{"counters":{"capture":0,"diff":0,"modify":0,"simulate":1,"spy":0,"synthesize":0}}}`)))
		})

		It("Should get the usage counters with 1 capture request when a request has been made", func() {
			hoverfly.SetMode("capture")

			proxyReq := sling.New().Get("http://www.google.com")
			hoverfly.Proxy(proxyReq)
			req := sling.New().Get("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/hoverfly/usage")
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))
			modeJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(modeJson).To(Equal([]byte(`{"usage":{"counters":{"capture":1,"diff":0,"modify":0,"simulate":0,"spy":0,"synthesize":0}}}`)))
		})

		It("Should get the usage counters with 1 modify request when a request has been made", func() {
			hoverfly.SetMode("modify")

			proxyReq := sling.New().Get("http://www.google.com")
			hoverfly.Proxy(proxyReq)
			req := sling.New().Get("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/hoverfly/usage")
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))
			modeJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(modeJson).To(Equal([]byte(`{"usage":{"counters":{"capture":0,"diff":0,"modify":1,"simulate":0,"spy":0,"synthesize":0}}}`)))
		})

		It("Should get the usage counters with 1 modify request when a request has been made", func() {
			hoverfly.SetMode("synthesize")

			proxyReq := sling.New().Get("http://www.google.com")
			hoverfly.Proxy(proxyReq)
			req := sling.New().Get("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/hoverfly/usage")
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))
			modeJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(modeJson).To(Equal([]byte(`{"usage":{"counters":{"capture":0,"diff":0,"modify":0,"simulate":0,"spy":0,"synthesize":1}}}`)))
		})

		It("Should get the usage counters with 1 spy request when a request has been made", func() {
			hoverfly.SetMode("spy")

			proxyReq := sling.New().Get("http://www.google.com")
			hoverfly.Proxy(proxyReq)
			req := sling.New().Get("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/hoverfly/usage")
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))
			modeJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(modeJson).To(Equal([]byte(`{"usage":{"counters":{"capture":0,"diff":0,"modify":0,"simulate":0,"spy":1,"synthesize":0}}}`)))
		})

		It("Should get the usage counters with 1 diff request when a request has been made", func() {
			hoverfly.SetMode("diff")

			proxyReq := sling.New().Get("http://www.google.com")
			hoverfly.Proxy(proxyReq)
			req := sling.New().Get("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/hoverfly/usage")
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))
			modeJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(modeJson).To(Equal([]byte(`{"usage":{"counters":{"capture":0,"diff":1,"modify":0,"simulate":0,"spy":0,"synthesize":0}}}`)))
		})
	})
})

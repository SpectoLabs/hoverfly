package hoverfly_test

import (
	"io/ioutil"
	"strings"

	"github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Interacting with the API", func() {

	BeforeEach(func() {
		hoverflyCmd = startHoverfly(adminPort, proxyPort)
	})

	AfterEach(func() {
		stopHoverfly()
	})

	Context("GET /api/v2/hoverfly", func() {

		It("Should get the hoverfly config", func() {
			req := sling.New().Get(hoverflyAdminUrl + "/api/v2/hoverfly")

			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))

			hoverflyJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(hoverflyJson).To(MatchRegexp(`"destination":"."`))
			Expect(hoverflyJson).To(MatchRegexp(`"middleware":{"binary":"","script":"","remote":""}`))
			Expect(hoverflyJson).To(MatchRegexp(`"usage":{"counters":{"capture":0,"modify":0,"simulate":0,"synthesize":0}}`))
			Expect(hoverflyJson).To(MatchRegexp(`"version":"v\d+.\d+.\d+"`))
			Expect(hoverflyJson).To(MatchRegexp(`"upstream-proxy":""`))
		})
	})

	Context("GET /api/v2/hoverfly/destination", func() {

		It("Should get the mode", func() {
			req := sling.New().Get(hoverflyAdminUrl + "/api/v2/hoverfly/destination")
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))
			modeJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(modeJson).To(Equal([]byte(`{"destination":"."}`)))
		})

	})

	Context("PUT /api/v2/hoverfly/destination", func() {

		It("Should put the mode", func() {
			req := sling.New().Put(hoverflyAdminUrl + "/api/v2/hoverfly/destination")
			req.Body(strings.NewReader(`{"destination":"test.com"}`))
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))
			modeJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(modeJson).To(Equal([]byte(`{"destination":"test.com"}`)))

			req = sling.New().Get(hoverflyAdminUrl + "/api/v2/hoverfly/destination")
			res = functional_tests.DoRequest(req)
			modeJson, err = ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(modeJson).To(Equal([]byte(`{"destination":"test.com"}`)))
		})

	})

	Context("GET /api/v2/hoverfly/mode", func() {

		It("Should get the mode", func() {
			req := sling.New().Get(hoverflyAdminUrl + "/api/v2/hoverfly/mode")
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))
			modeJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(modeJson).To(Equal([]byte(`{"mode":"simulate"}`)))
		})
	})

	Context("PUT /api/v2/hoverfly/mode", func() {

		It("Should put the mode", func() {
			req := sling.New().Put(hoverflyAdminUrl + "/api/v2/hoverfly/mode")
			req.Body(strings.NewReader(`{"mode":"capture"}`))
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))
			modeJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(modeJson).To(Equal([]byte(`{"mode":"capture"}`)))

			req = sling.New().Get(hoverflyAdminUrl + "/api/v2/hoverfly/mode")
			res = functional_tests.DoRequest(req)
			modeJson, err = ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(modeJson).To(Equal([]byte(`{"mode":"capture"}`)))
		})

	})

	Context("GET /api/v2/hoverfly/middleware", func() {

		It("Should get the middleware which should be blank", func() {
			req := sling.New().Get(hoverflyAdminUrl + "/api/v2/hoverfly/middleware")
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))
			modeJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(modeJson).To(Equal([]byte(`{"binary":"","script":"","remote":""}`)))
		})
	})

	Context("PUT /api/v2/hoverfly/middleware", func() {

		It("Should put the middleware", func() {
			req := sling.New().Put(hoverflyAdminUrl + "/api/v2/hoverfly/middleware")
			req.Body(strings.NewReader(`{"binary":"ruby", "script":"#!/usr/bin/env ruby\n# encoding: utf-8\nwhile payload = STDIN.gets\nnext unless payload\n\nSTDOUT.puts payload\nend"}`))
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))
			modeJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(modeJson).To(Equal([]byte(`{"binary":"ruby","script":"#!/usr/bin/env ruby\n# encoding: utf-8\nwhile payload = STDIN.gets\nnext unless payload\n\nSTDOUT.puts payload\nend","remote":""}`)))

			req = sling.New().Get(hoverflyAdminUrl + "/api/v2/hoverfly/middleware")
			res = functional_tests.DoRequest(req)
			modeJson, err = ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(modeJson).To(Equal([]byte(`{"binary":"ruby","script":"#!/usr/bin/env ruby\n# encoding: utf-8\nwhile payload = STDIN.gets\nnext unless payload\n\nSTDOUT.puts payload\nend","remote":""}`)))
		})

	})

	Context("GET /api/v2/hoverfly/usage", func() {

		It("Should get the usage counters", func() {
			req := sling.New().Get(hoverflyAdminUrl + "/api/v2/hoverfly/usage")
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))
			modeJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(modeJson).To(Equal([]byte(`{"usage":{"counters":{"capture":0,"modify":0,"simulate":0,"synthesize":0}}}`)))
		})

		It("Should get the usage counters with 1 simulate request when a request has been made", func() {
			proxyReq := sling.New().Get("http://www.google.com")
			DoRequestThroughProxy(proxyReq)
			req := sling.New().Get(hoverflyAdminUrl + "/api/v2/hoverfly/usage")
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))
			modeJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(modeJson).To(Equal([]byte(`{"usage":{"counters":{"capture":0,"modify":0,"simulate":1,"synthesize":0}}}`)))
		})

		It("Should get the usage counters with 1 capture request when a request has been made", func() {
			SetHoverflyMode("capture")

			proxyReq := sling.New().Get("http://www.google.com")
			DoRequestThroughProxy(proxyReq)
			req := sling.New().Get(hoverflyAdminUrl + "/api/v2/hoverfly/usage")
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))
			modeJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(modeJson).To(Equal([]byte(`{"usage":{"counters":{"capture":1,"modify":0,"simulate":0,"synthesize":0}}}`)))
		})

		It("Should get the usage counters with 1 modify request when a request has been made", func() {
			SetHoverflyMode("modify")

			proxyReq := sling.New().Get("http://www.google.com")
			DoRequestThroughProxy(proxyReq)
			req := sling.New().Get(hoverflyAdminUrl + "/api/v2/hoverfly/usage")
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))
			modeJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(modeJson).To(Equal([]byte(`{"usage":{"counters":{"capture":0,"modify":1,"simulate":0,"synthesize":0}}}`)))
		})

		It("Should get the usage counters with 1 modify request when a request has been made", func() {
			SetHoverflyMode("synthesize")

			proxyReq := sling.New().Get("http://www.google.com")
			DoRequestThroughProxy(proxyReq)
			req := sling.New().Get(hoverflyAdminUrl + "/api/v2/hoverfly/usage")
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))
			modeJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(modeJson).To(Equal([]byte(`{"usage":{"counters":{"capture":0,"modify":0,"simulate":0,"synthesize":1}}}`)))
		})
	})

	Context("GET /api/v2/hoverfly/version", func() {

		It("Should get the version", func() {
			req := sling.New().Get(hoverflyAdminUrl + "/api/v2/hoverfly/version")
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))
			modeJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(string(modeJson)).To(MatchRegexp(`{"version":"v\d+.\d+.\d+"}`))
		})
	})

	Context("GET /api/v2/hoverfly/upstream-proxy", func() {

		It("Should get the upstream proxy", func() {
			req := sling.New().Get(hoverflyAdminUrl + "/api/v2/hoverfly/upstream-proxy")
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))
			modeJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(string(modeJson)).To(MatchRegexp(`{"upstream-proxy":""}`))
		})
	})
})

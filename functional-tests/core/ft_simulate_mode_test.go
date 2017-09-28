package hoverfly_test

import (
	"io/ioutil"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/util"
	"github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("When I run Hoverfly in simulate mode", func() {

	var (
		hoverfly *functional_tests.Hoverfly
	)

	BeforeEach(func() {
		hoverfly = functional_tests.NewHoverfly()
		hoverfly.Start()
		hoverfly.SetMode("simulate")
	})

	AfterEach(func() {
		hoverfly.Stop()
	})

	It("should match against the first request matcher in simulation", func() {
		hoverfly.ImportSimulation(functional_tests.JsonPayload)

		resp := hoverfly.Proxy(sling.New().Get("http://test-server.com/path1"))
		Expect(resp.StatusCode).To(Equal(200))

		body, err := ioutil.ReadAll(resp.Body)
		Expect(err).To(BeNil())

		Expect(string(body)).To(Equal("exact match"))
		Expect(resp.Header).To(HaveKeyWithValue("Header", []string{"value1", "value2"}))
	})

	It("should match against the first request matcher in simulation over HTTPS", func() {
		hoverfly.ImportSimulation(functional_tests.JsonPayload)

		resp := hoverfly.Proxy(sling.New().Get("https://test-server.com/path1"))
		Expect(resp.StatusCode).To(Equal(200))

		body, err := ioutil.ReadAll(resp.Body)
		Expect(err).To(BeNil())

		Expect(string(body)).To(Equal("exact match"))
		Expect(resp.Header).To(HaveKeyWithValue("Header", []string{"value1", "value2"}))
	})

	It("should match against the second request matcher in simulation", func() {
		hoverfly.ImportSimulation(functional_tests.JsonPayload)

		slingRequest := sling.New().Get("http://destination-server.com/should-match-regardless")
		response := hoverfly.Proxy(slingRequest)

		body, err := ioutil.ReadAll(response.Body)
		Expect(err).To(BeNil())
		Expect(string(body)).To(Equal("destination matched"))
	})

	It("should match against the second request matcher in simulation over HTTPS", func() {
		hoverfly.ImportSimulation(functional_tests.JsonPayload)

		slingRequest := sling.New().Get("http://destination-server.com/should-match-regardless")
		response := hoverfly.Proxy(slingRequest)

		body, err := ioutil.ReadAll(response.Body)
		Expect(err).To(BeNil())
		Expect(string(body)).To(Equal("destination matched"))
	})

	It("should apply middleware to the cached response", func() {
		hoverfly.SetMiddleware("python", functional_tests.Middleware)
		hoverfly.ImportSimulation(functional_tests.JsonPayload)

		resp := hoverfly.Proxy(sling.New().Get("http://test-server.com/path1"))

		body, err := ioutil.ReadAll(resp.Body)
		Expect(err).To(BeNil())

		Expect(string(body)).To(Equal("CHANGED_RESPONSE_BODY"))
	})

	It("Should perform a strongest match by default", func() {

		hoverfly.ImportSimulation(functional_tests.StrongestMatchProofSimulation)

		slingRequest := sling.New().Get("http://destination.com/should-match-strongest")
		response := hoverfly.Proxy(slingRequest)

		body, err := ioutil.ReadAll(response.Body)
		Expect(err).To(BeNil())
		Expect(string(body)).To(Equal("second and strongest match"))
	})

	It("Should perform a strongest match when set explicitly", func() {
		hoverfly.SetModeWithArgs("simulate", v2.ModeArgumentsView{
			MatchingStrategy: util.StringToPointer("strongest"),
		})

		hoverfly.ImportSimulation(functional_tests.StrongestMatchProofSimulation)

		slingRequest := sling.New().Get("http://destination.com/should-match-strongest")
		response := hoverfly.Proxy(slingRequest)

		body, err := ioutil.ReadAll(response.Body)
		Expect(err).To(BeNil())
		Expect(string(body)).To(Equal("second and strongest match"))
	})

	It("Should perform a strongest match when set explicitly", func() {
		hoverfly.SetModeWithArgs("simulate", v2.ModeArgumentsView{
			MatchingStrategy: util.StringToPointer("first"),
		})

		hoverfly.ImportSimulation(functional_tests.StrongestMatchProofSimulation)

		slingRequest := sling.New().Get("http://destination.com/should-match-strongest")
		response := hoverfly.Proxy(slingRequest)

		body, err := ioutil.ReadAll(response.Body)
		Expect(err).To(BeNil())
		Expect(string(body)).To(Equal("first and weakest match"))
	})

	It("Should respond with the closest miss, once from matchers & once from cache", func() {

		hoverfly.ImportSimulation(functional_tests.ClosestMissProofSimulation)

		slingRequest := sling.New().Get("http://destination.com/closest-miss")
		response := hoverfly.Proxy(slingRequest)

		body, err := ioutil.ReadAll(response.Body)
		Expect(err).To(BeNil())

		expected := `Hoverfly Error!

There was an error when matching

Got error: Could not find a match for request, create or record a valid matcher first!

The following request was made, but was not matched by Hoverfly:

{
    "Path": "/closest-miss",
    "Method": "GET",
    "Destination": "destination.com",
    "Scheme": "http",
    "Query": {},
    "Body": "",
    "Headers": {
        "Accept-Encoding": [
            "gzip"
        ],
        "User-Agent": [
            "Go-http-client/1.1"
        ]
    }
}

Whilst Hoverfly has the following state:

{}

The matcher which came closest was:

{
    "path": {
        "exactMatch": "/closest-miss"
    },
    "destination": {
        "exactMatch": "destination.com"
    },
    "body": {
        "exactMatch": "body"
    }
}

But it did not match on the following fields:

[body]

Which if hit would have given the following response:

{
    "status": 200,
    "body": "",
    "encodedBody": false,
    "templated": false
}`
		Expect(string(body)).To(Equal(expected))

		slingRequest = sling.New().Get("http://destination.com/closest-miss")
		response = hoverfly.Proxy(slingRequest)

		body, err = ioutil.ReadAll(response.Body)
		Expect(err).To(BeNil())

		expected = `Hoverfly Error!

There was an error when matching

Got error: Could not find a match for request, create or record a valid matcher first!

The following request was made, but was not matched by Hoverfly:

{
    "Path": "/closest-miss",
    "Method": "GET",
    "Destination": "destination.com",
    "Scheme": "http",
    "Query": {},
    "Body": "",
    "Headers": {
        "Accept-Encoding": [
            "gzip"
        ],
        "User-Agent": [
            "Go-http-client/1.1"
        ]
    }
}

Whilst Hoverfly has the following state:

{}

The matcher which came closest was:

{
    "path": {
        "exactMatch": "/closest-miss"
    },
    "destination": {
        "exactMatch": "destination.com"
    },
    "body": {
        "exactMatch": "body"
    }
}

But it did not match on the following fields:

[body]

Which if hit would have given the following response:

{
    "status": 200,
    "body": "",
    "encodedBody": false,
    "templated": false
}`
		Expect(string(body)).To(Equal(expected))
	})

	It("should no longer cause issue #607", func() {

		hoverfly.ImportSimulation(functional_tests.Issue607)

		// Match
		i := sling.New().Get("https://domain.com/billing/v1/servicequotes/123456?saleschannel=RETAIL")
		i.Set("Accept", "application/json")
		i.Set("Activityid", "ChangeMSISDN_CR_PushtoBill(Get)-200")
		i.Set("Applicationid", "ACUI")
		i.Set("Authorization", "Bearer token")
		i.Set("Cache-Control", "no-cache")
		i.Set("Channelid", "RETAIL")
		i.Set("Content-Type", "application/json")
		i.Set("Interactionid", "123456787")
		i.Set("Senderid", "ACUI")
		i.Set("User-Agent", "curl/7.54.0")
		i.Set("Workflowid", "CHANGEMSISDN")

		resp := hoverfly.Proxy(i)

		Expect(resp.StatusCode).To(Equal(200))
		Expect(hoverfly.GetCache().Cache).To(BeEmpty()) // Don't cache hits which include header matching

		// Miss
		i = sling.New().Get("https://domain.com/billing/v1/servicequotes/123456?saleschannel=RETAIL")
		i.Set("Accept", "application/json")
		i.Set("Activityid", "ChangeMSISDN_Procedural(Get)-200")
		i.Set("Applicationid", "ACUI")
		i.Set("Authorization", "Bearer token")
		i.Set("Cache-Control", "no-cache")
		i.Set("Channelid", "RETAIL")
		i.Set("Content-Type", "application/json")
		i.Set("Interactionid", "123456787")
		i.Set("Senderid", "ACUI")
		i.Set("User-Agent", "curl/7.54.0")
		i.Set("Workflowid", "CHANGEMSISDN")

		resp = hoverfly.Proxy(i)
		Expect(resp.StatusCode).To(Equal(502))
		Expect(hoverfly.GetCache().Cache).To(BeEmpty()) // Don't cache misses when only headers were not matched

		// Match again
		i = sling.New().Get("https://domain.com/billing/v1/servicequotes/123456?saleschannel=RETAIL")
		i.Set("Accept", "application/json")
		i.Set("Activityid", "ChangeMSISDN_CR_PushtoBill(Get)-200")
		i.Set("Applicationid", "ACUI")
		i.Set("Authorization", "Bearer token")
		i.Set("Cache-Control", "no-cache")
		i.Set("Channelid", "RETAIL")
		i.Set("Content-Type", "application/json")
		i.Set("Interactionid", "123456787")
		i.Set("Senderid", "ACUI")
		i.Set("User-Agent", "curl/7.54.0")
		i.Set("Workflowid", "CHANGEMSISDN")

		resp = hoverfly.Proxy(i)

		body, _ := ioutil.ReadAll(resp.Body)
		GinkgoWriter.Write(body)

		Expect(resp.StatusCode).To(Equal(200))
		Expect(hoverfly.GetCache().Cache).To(BeEmpty()) // Don't cache hits which include header matching
	})

	It("should template response if templating is enabled and cache template not response", func() {
		hoverfly.ImportSimulation(functional_tests.TemplatingEnabled)

		hoverfly.WriteLogsIfError()

		resp := hoverfly.Proxy(sling.New().Get("http://test-server.com?one=foo"))
		Expect(resp.StatusCode).To(Equal(200))

		body, err := ioutil.ReadAll(resp.Body)
		Expect(err).To(BeNil())

		Expect(string(body)).To(Equal("foo"))

		resp = hoverfly.Proxy(sling.New().Get("http://test-server.com?one=bar"))
		Expect(resp.StatusCode).To(Equal(200))

		body, err = ioutil.ReadAll(resp.Body)
		Expect(err).To(BeNil())

		Expect(string(body)).To(Equal("bar"))
	})

	It("should be able to use state in templating", func() {
		hoverfly.ImportSimulation(functional_tests.TemplatingEnabledWithStateInBody)

		resp := hoverfly.Proxy(sling.New().Get("http://test-server.com/one"))
		Expect(resp.StatusCode).To(Equal(200))

		resp = hoverfly.Proxy(sling.New().Get("http://test-server.com/two"))
		Expect(resp.StatusCode).To(Equal(200))
		body, err := ioutil.ReadAll(resp.Body)
		Expect(err).To(BeNil())
		Expect(string(body)).To(Equal("state for eggs"))
	})

	It("should not template response if templating is disabled explicitely", func() {
		hoverfly.ImportSimulation(functional_tests.TemplatingDisabled)

		resp := hoverfly.Proxy(sling.New().Get("http://test-server.com?one=foo"))
		Expect(resp.StatusCode).To(Equal(200))

		body, err := ioutil.ReadAll(resp.Body)
		Expect(err).To(BeNil())

		Expect(string(body)).To(Equal("{{ Request.QueryParam.singular }}"))
	})

	It("should not template response if templating is not explcitely enabled or disabled", func() {
		hoverfly.ImportSimulation(functional_tests.TemplatingDisabledByDefault)

		resp := hoverfly.Proxy(sling.New().Get("http://test-server.com?one=foo"))
		Expect(resp.StatusCode).To(Equal(200))

		body, err := ioutil.ReadAll(resp.Body)
		Expect(err).To(BeNil())

		Expect(string(body)).To(Equal("{{ Request.QueryParam.one }}"))
	})

	It("should not crash when templating a response if templating variable does not exist", func() {
		hoverfly.ImportSimulation(functional_tests.TemplatingEnabled)

		resp := hoverfly.Proxy(sling.New().Get("http://test-server.com?wrong=foo"))
		Expect(resp.StatusCode).To(Equal(200))

		body, err := ioutil.ReadAll(resp.Body)
		Expect(err).To(BeNil())

		Expect(string(body)).To(Equal(""))
	})
})

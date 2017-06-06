package hoverfly_test

import (
	"io/ioutil"

	"github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/util"
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

	It("should match against the second request matcher in simulation", func() {
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
			MatchingStrategy : util.StringToPointer("STRONGEST"),
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
			MatchingStrategy : util.StringToPointer("FIRST"),
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
    "Query": "",
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

The closest miss was the following matcher:

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
    "headers": null
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
    "Query": "",
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

The closest miss was the following matcher:

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
    "headers": null
}`
		Expect(string(body)).To(Equal(expected))
	})
})

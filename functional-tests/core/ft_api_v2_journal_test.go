package hoverfly_test

import (
	"encoding/json"
	"io/ioutil"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("/api/v2/journal", func() {

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

		It("should display an empty journal", func() {
			req := sling.New().Get("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/journal")
			res := functional_tests.DoRequest(req)

			Expect(res.StatusCode).To(Equal(200))

			responseJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())

			var journal []v2.JournalEntryView

			err = json.Unmarshal(responseJson, &journal)
			Expect(err).To(BeNil())

			Expect(journal).To(HaveLen(0))
		})

		It("should display one item in the journal", func() {
			hoverfly.Proxy(sling.New().Get("http://hoverfly.io"))

			req := sling.New().Get("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/journal")
			res := functional_tests.DoRequest(req)

			Expect(res.StatusCode).To(Equal(200))

			responseJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())

			var journal []v2.JournalEntryView

			err = json.Unmarshal(responseJson, &journal)
			Expect(err).To(BeNil())

			Expect(journal).To(HaveLen(1))

			Expect(*journal[0].Request.Scheme).To(Equal("http"))
			Expect(*journal[0].Request.Method).To(Equal("GET"))
			Expect(*journal[0].Request.Destination).To(Equal("hoverfly.io"))
			Expect(*journal[0].Request.Path).To(Equal("/"))
			Expect(*journal[0].Request.Query).To(Equal(""))
			Expect(journal[0].Request.Headers["Accept-Encoding"]).To(ContainElement("gzip"))
			Expect(journal[0].Request.Headers["User-Agent"]).To(ContainElement("Go-http-client/1.1"))

			Expect(journal[0].Response.Status).To(Equal(502))
			Expect(journal[0].Response.Body).To(Equal("Hoverfly Error!\n\nThere was an error when matching\n\nGot error: Could not find a match for request, create or record a valid matcher first!"))
			Expect(journal[0].Response.Headers["Content-Type"]).To(ContainElement("text/plain"))

			Expect(journal[0].Latency).To(BeNumerically("<", 1))
			Expect(journal[0].Mode).To(Equal("simulate"))
		})

		It("should display multiple items in the journal", func() {
			hoverfly.Proxy(sling.New().Get("http://hoverfly.io"))
			hoverfly.Proxy(sling.New().Get("http://github.com/SpectoLabs/hoverfly"))
			hoverfly.Proxy(sling.New().Get("http://specto.io"))

			req := sling.New().Get("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/journal")
			res := functional_tests.DoRequest(req)

			Expect(res.StatusCode).To(Equal(200))

			responseJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())

			var journal []v2.JournalEntryView

			err = json.Unmarshal(responseJson, &journal)
			Expect(err).To(BeNil())

			Expect(journal).To(HaveLen(3))

			Expect(*journal[0].Request.Destination).To(Equal("hoverfly.io"))
			Expect(*journal[0].Request.Path).To(Equal("/"))

			Expect(*journal[1].Request.Destination).To(Equal("github.com"))
			Expect(*journal[1].Request.Path).To(Equal("/SpectoLabs/hoverfly"))

			Expect(*journal[2].Request.Destination).To(Equal("specto.io"))
			Expect(*journal[2].Request.Path).To(Equal("/"))
		})

		It("should display the mode each request was in", func() {
			hoverfly.SetMode("simulate")
			hoverfly.Proxy(sling.New().Get("http://localhost:" + hoverfly.GetAdminPort()))

			hoverfly.SetMode("capture")
			hoverfly.Proxy(sling.New().Get("http://localhost:" + hoverfly.GetAdminPort()))

			req := sling.New().Get("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/journal")
			res := functional_tests.DoRequest(req)

			Expect(res.StatusCode).To(Equal(200))

			responseJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())

			var journal []v2.JournalEntryView

			err = json.Unmarshal(responseJson, &journal)
			Expect(err).To(BeNil())

			Expect(journal).To(HaveLen(2))

			Expect(journal[0].Mode).To(Equal("simulate"))
			Expect(journal[1].Mode).To(Equal("capture"))
		})
	})
})

package hoverfly_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Running Hoverfly", func() {

	var (
		hoverfly   *functional_tests.Hoverfly
		fakeServer *httptest.Server
	)

	BeforeEach(func() {
		hoverfly = functional_tests.NewHoverfly()
		hoverfly.Start()
		hoverfly.SetMode("capture")

		fakeServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			w.Header().Set("date", "date")
			w.Write([]byte("Hello world"))
		}))
	})

	AfterEach(func() {
		hoverfly.Stop()
		fakeServer.Close()
	})

	Context("in capture mode", func() {

		It("Should not capture if destination does not match", func() {
			hoverfly.SetDestination("notlocalhost")

			resp := hoverfly.Proxy(sling.New().Get(fakeServer.URL))

			Expect(resp.StatusCode).To(Equal(200))
			Expect(resp.Header.Get("date")).To(Equal("date"))

			recordsJson, err := ioutil.ReadAll(hoverfly.GetSimulation())
			Expect(err).To(BeNil())
			Expect(recordsJson).ToNot(ContainSubstring(`"destination":[{"matcher":"exact","value":"127.0.0.1`))
		})

		It("Should capture if destination is 127.0.0.1", func() {
			hoverfly.SetDestination("127.0.0.1")

			resp := hoverfly.Proxy(sling.New().Get(fakeServer.URL))

			Expect(resp.StatusCode).To(Equal(200))
			Expect(resp.Header.Get("date")).To(Equal("date"))

			recordsJson, err := ioutil.ReadAll(hoverfly.GetSimulation())
			Expect(err).To(BeNil())
			Expect(recordsJson).To(ContainSubstring(`"destination":[{"matcher":"exact","value":"127.0.0.1`))
		})

		It("Should capture if destination is set to port numbers", func() {
			hoverfly.SetDestination(strings.Replace(fakeServer.URL, "http://127.0.0.1", "", 1))

			resp := hoverfly.Proxy(sling.New().Get(fakeServer.URL))

			Expect(resp.StatusCode).To(Equal(200))
			Expect(resp.Header.Get("date")).To(Equal("date"))

			recordsJson, err := ioutil.ReadAll(hoverfly.GetSimulation())
			Expect(err).To(BeNil())
			Expect(recordsJson).To(ContainSubstring(`"destination":[{"matcher":"exact","value":"127.0.0.1`))
		})

		It("Should capture if destination is set to the path", func() {
			hoverfly.SetDestination("/path")

			resp := hoverfly.Proxy(sling.New().Get(fakeServer.URL + "/path"))

			Expect(resp.StatusCode).To(Equal(200))
			Expect(resp.Header.Get("date")).To(Equal("date"))

			recordsJson, err := ioutil.ReadAll(hoverfly.GetSimulation())
			Expect(err).To(BeNil())
			Expect(recordsJson).To(ContainSubstring(`"path":[{"matcher":"exact","value":"/path"`))
		})

		It("Should not capture if destination is set to the wrong path", func() {
			hoverfly.SetDestination("/wrongpath")

			resp := hoverfly.Proxy(sling.New().Get(fakeServer.URL + "/path"))

			Expect(resp.StatusCode).To(Equal(200))
			Expect(resp.Header.Get("date")).To(Equal("date"))

			recordsJson, err := ioutil.ReadAll(hoverfly.GetSimulation())
			Expect(err).To(BeNil())
			Expect(recordsJson).ToNot(ContainSubstring(`"path":[{"matcher":"exact","value":"/path"`))
		})
	})
})

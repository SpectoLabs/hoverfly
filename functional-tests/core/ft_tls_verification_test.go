package hoverfly_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("When I run Hoverfly", func() {

	var (
		hoverfly *functional_tests.Hoverfly
	)

	BeforeEach(func() {
		hoverfly = functional_tests.NewHoverfly()
	})

	AfterEach(func() {
		hoverfly.Stop()
	})

	Context("don't specify tls-verification", func() {

		BeforeEach(func() {
			hoverfly.Start()
			hoverfly.SetMode("capture")
		})

		It("should not error with https with bad certificate", func() {

			fakeServer := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/plain")
				w.Header().Set("Date", "date")
				w.Write([]byte("Hello world"))
			}))

			fakeServer.StartTLS()
			defer fakeServer.Close()

			resp := hoverfly.Proxy(sling.New().Get(fakeServer.URL))
			Expect(resp.StatusCode).To(Equal(502))

			responseBody, _ := ioutil.ReadAll(resp.Body)

			Expect(string(responseBody)).To(ContainSubstring("Hoverfly Error!"))
			Expect(string(responseBody)).To(ContainSubstring("There was an error when forwarding the request to the intended destination"))
			Expect(string(responseBody)).To(ContainSubstring("x509: certificate signed by unknown authority"))
		})
	})

	Context("tls-verification=true", func() {

		BeforeEach(func() {
			hoverfly.Start("-tls-verification=true")
			hoverfly.SetMode("capture")
		})

		It("should not error with https with bad certificate", func() {

			fakeServer := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/plain")
				w.Header().Set("Date", "date")
				w.Write([]byte("Hello world"))
			}))

			fakeServer.StartTLS()
			defer fakeServer.Close()

			resp := hoverfly.Proxy(sling.New().Get(fakeServer.URL))
			Expect(resp.StatusCode).To(Equal(502))

			responseBody, _ := ioutil.ReadAll(resp.Body)

			Expect(string(responseBody)).To(ContainSubstring("Hoverfly Error!"))
			Expect(string(responseBody)).To(ContainSubstring("There was an error when forwarding the request to the intended destination"))
			Expect(string(responseBody)).To(ContainSubstring("x509: certificate signed by unknown authority"))
		})
	})

	Context("tls-verification=false", func() {

		BeforeEach(func() {
			hoverfly.Start("-tls-verification=false")
			hoverfly.SetMode("capture")
		})

		It("should not error with https with bad certificate", func() {

			fakeServer := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/plain")
				w.Header().Set("Date", "date")
				w.Write([]byte("Hello world"))
			}))

			fakeServer.StartTLS()
			defer fakeServer.Close()

			resp := hoverfly.Proxy(sling.New().Get(fakeServer.URL))
			Expect(resp.StatusCode).To(Equal(200))

			responseBody, _ := ioutil.ReadAll(resp.Body)

			Expect(string(responseBody)).To(ContainSubstring("Hello world"))
		})
	})
})

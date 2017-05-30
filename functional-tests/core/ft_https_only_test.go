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

	Context("and specify https-only on startup", func() {

		BeforeEach(func() {
			hoverfly.Start("-https-only", "-tls-verification=false")
			hoverfly.SetMode("capture")
		})

		It("should error with http requests", func() {

			fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/plain")
				w.Header().Set("Date", "date")
				w.Write([]byte("Hello world"))
			}))

			defer fakeServer.Close()

			resp := hoverfly.Proxy(sling.New().Get(fakeServer.URL))
			Expect(resp.StatusCode).To(Equal(502))

			responseBody, _ := ioutil.ReadAll(resp.Body)
			Expect(string(responseBody)).To(ContainSubstring("This proxy requires TLS (HTTPS)"))
		})

		It("should not error with https requests", func() {

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

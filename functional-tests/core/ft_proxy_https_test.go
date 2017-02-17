package hoverfly_test

import (
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

	Context("and it uses default certificate and key configuration", func() {

		BeforeEach(func() {
			hoverfly.Start()
			hoverfly.SetMode("capture")
		})

		AfterEach(func() {
			hoverfly.Stop()
		})

		It("should respond with HTTPS responses with the default Hoverfly certificate", func() {

			fakeServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/plain")
				w.Header().Set("Date", "date")
				w.Write([]byte("Hello world"))
			}))

			defer fakeServer.Close()

			response := hoverfly.Proxy(sling.New().Get(fakeServer.URL))

			Expect(response.TLS.PeerCertificates[0].Issuer.CommonName).To(Equal("hoverfly.proxy"))
			Expect(response.TLS.PeerCertificates[0].Issuer.Organization).To(ContainElement("Hoverfly Authority"))
			Expect(response.TLS.PeerCertificates[0].Subject.Names[0].Value).To(Equal("GoProxy untrusted MITM proxy Inc"))
		})
	})

	Context("and it uses default certificate and key configuration", func() {

		BeforeEach(func() {
			hoverfly.Start("-cert", "testdata/cert.pem", "-key", "testdata/key.pem")
			hoverfly.SetMode("capture")
		})

		AfterEach(func() {
			hoverfly.Stop()
		})

		It("should respond with HTTPS responses with the default Hoverfly certificate", func() {

			fakeServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/plain")
				w.Header().Set("Date", "date")
				w.Write([]byte("Hello world"))
			}))

			defer fakeServer.Close()

			response := hoverfly.Proxy(sling.New().Get(fakeServer.URL))

			Expect(response.TLS.PeerCertificates[0].Issuer.CommonName).To(Equal("test.cert"))
			Expect(response.TLS.PeerCertificates[0].Issuer.Organization).To(ContainElement("Testdata Certificate Authority"))
			Expect(response.TLS.PeerCertificates[0].Subject.Names[0].Value).To(Equal("GoProxy untrusted MITM proxy Inc"))
		})
	})
})

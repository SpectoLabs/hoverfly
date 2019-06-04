package hoverfly_test

import (
	"github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
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

	Context("and specify cors on startup", func() {

		BeforeEach(func() {
			hoverfly.Start("-cors")
		})

		It("should handle a pre-flight request", func() {

			fakeServer := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/plain")
				w.Header().Set("Date", "date")
				w.Write([]byte("Hello world"))
			}))

			fakeServer.StartTLS()
			defer fakeServer.Close()

			resp := hoverfly.Proxy(sling.New().Options(fakeServer.URL).Add("Origin", "http://some-host.com").Add("Access-Control-Request-Methods", "POST"))
			Expect(resp.StatusCode).To(Equal(200))
			Expect(resp.Header.Get("Access-Control-Allow-Origin")).To(Equal("*"))
			Expect(resp.Header.Get("Access-Control-Allow-Methods")).To(Equal("GET,POST,PUT,PATCH,DELETE,HEAD,OPTIONS"))

			responseBody, _ := ioutil.ReadAll(resp.Body)
			Expect(string(responseBody)).To(ContainSubstring(""))
		})
	})
})

package hoverfly_test

import (
	"net/http"

	"github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("hoverfly -dev", func() {

	var (
		hoverfly *functional_tests.Hoverfly
	)

	Context("unauthenticated Hoverfly", func() {

		BeforeEach(func() {
			hoverfly = functional_tests.NewHoverfly()
			hoverfly.Start("-dev")
		})

		AfterEach(func() {
			hoverfly.Stop()
		})

		It("should add CORS headers to API responses", func() {
			req := sling.New().Get("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/hoverfly")
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(http.StatusOK))

			Expect(res.Header.Get("Access-Control-Allow-Origin")).To(Equal("http://localhost:4200"))
			Expect(res.Header.Get("Access-Control-Allow-Methods")).To(Equal("GET, PUT, POST, OPTIONS, DELETE"))
			Expect(res.Header.Get("Access-Control-Allow-Headers")).To(Equal("Origin, X-Requested-With, Content-Type, Accept, Authorization"))
			Expect(res.Header.Get("Access-Control-Allow-Credentials")).To(Equal("true"))
		})
	})

	Context("authenticated Hoverfly", func() {

		BeforeEach(func() {
			hoverfly = functional_tests.NewHoverfly()
			hoverfly.Start("-dev", "-auth", "-username", "benji", "-password", "password")
		})

		AfterEach(func() {
			hoverfly.Stop()
		})

		It("should add CORS headers to unauthenticated error responses", func() {
			req := sling.New().Get("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/hoverfly")
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(http.StatusUnauthorized))

			Expect(res.Header.Get("Access-Control-Allow-Origin")).To(Equal("http://localhost:4200"))
			Expect(res.Header.Get("Access-Control-Allow-Methods")).To(Equal("GET, PUT, POST, OPTIONS, DELETE"))
			Expect(res.Header.Get("Access-Control-Allow-Headers")).To(Equal("Origin, X-Requested-With, Content-Type, Accept, Authorization"))
			Expect(res.Header.Get("Access-Control-Allow-Credentials")).To(Equal("true"))
		})
	})
})

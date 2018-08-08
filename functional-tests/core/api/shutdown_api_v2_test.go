package api_test

import (
	"net/http"

	"github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("/api/v2/shutdown", func() {

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

	Context("DELETE", func() {

		It("should shutdown after 10 seconds", func() {
			response := functional_tests.DoRequest(sling.New().Delete("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/shutdown"))
			Expect(response.StatusCode).To(Equal(http.StatusOK))

			request, _ := sling.New().Get("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/hoverfly/mode").Request()
			_, err := http.DefaultClient.Do(request)

			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(ContainSubstring("connection refused"))

		})
	})
})

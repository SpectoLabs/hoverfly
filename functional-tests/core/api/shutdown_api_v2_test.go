package api_test

import (
	"net/http"
	"time"

	"github.com/SpectoLabs/hoverfly/v2/functional-tests"
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

		It("should shut down immediately", func() {
			response := functional_tests.DoRequest(sling.New().Delete("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/shutdown"))
			Expect(response.StatusCode).To(Equal(http.StatusOK))
			time.Sleep(2 * time.Second) // in case there are some shut down delay

			request, _ := sling.New().Get("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/hoverfly/mode").Request()
			_, err := http.DefaultClient.Do(request)

			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(ContainSubstring("connection refused"))

		})
	})
})

package api_test

import (
	"io/ioutil"
	"strings"

	"github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("/api/v2/hoverfly/pac", func() {

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

		It("Should get 404 when PAC not set", func() {
			req := sling.New().Get("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/hoverfly/pac")
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(404))
			modeJson, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(modeJson).To(Equal([]byte(`{"error":"Not found"}`)))
		})
	})

	Context("PUT", func() {

		It("Should put the PAC file", func() {
			req := sling.New().Put("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/hoverfly/pac")
			req.Body(strings.NewReader(`PACFILE`))
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))
			responseBody, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(responseBody).To(Equal([]byte(`PACFILE`)))

			req = sling.New().Get("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/hoverfly/pac")
			res = functional_tests.DoRequest(req)
			responseBody, err = ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(responseBody).To(Equal([]byte(`PACFILE`)))
		})
	})

	Context("DELETE", func() {

		It("Should delete the PAC file", func() {
			hoverfly.SetPACFile("PACFILE")
			req := sling.New().Delete("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/hoverfly/pac")
			res := functional_tests.DoRequest(req)
			Expect(res.StatusCode).To(Equal(200))
			responseBody, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(responseBody).To(Equal([]byte(``)))

			req = sling.New().Get("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/hoverfly/pac")
			res = functional_tests.DoRequest(req)
			responseBody, err = ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())
			Expect(responseBody).To(Equal([]byte(`{"error":"Not found"}`)))
		})
	})
})

package hoverctl_end_to_end

import (
	"io/ioutil"

	"github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/antonholmquist/jason"
	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("hoverctl flush cache", func() {

	var (
		hoverfly *functional_tests.Hoverfly
	)

	BeforeEach(func() {
		hoverfly = functional_tests.NewHoverfly()
		hoverfly.Start()
		hoverfly.SetMode("simulate")
		hoverfly.ImportSimulation(functional_tests.JsonPayload)
		hoverfly.Proxy(sling.New().Get("http://template-server.com"))

		WriteConfiguration("localhost", hoverfly.GetAdminPort(), hoverfly.GetProxyPort())
	})

	AfterEach(func() {
		hoverfly.Stop()
	})

	It("should flush cache", func() {
		output := functional_tests.Run(hoverctlBinary, "flush", "--force")

		Expect(output).To(ContainSubstring("Successfully flushed cache"))

		req := sling.New().Get("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/cache")
		res := functional_tests.DoRequest(req)
		Expect(res.StatusCode).To(Equal(200))
		responseJson, err := ioutil.ReadAll(res.Body)
		Expect(err).To(BeNil())

		jsonObject, err := jason.NewObjectFromBytes(responseJson)
		Expect(err).To(BeNil())

		cacheArray, err := jsonObject.GetObjectArray("cache")
		Expect(err).To(BeNil())

		Expect(cacheArray).To(HaveLen(0))
	})

})

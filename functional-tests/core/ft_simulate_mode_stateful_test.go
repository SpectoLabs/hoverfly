package hoverfly_test

import (
	"io/ioutil"

	"github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/SpectoLabs/hoverfly/functional-tests/testdata"
	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("When I run Hoverfly in simulate mode", func() {

	var (
		hoverfly *functional_tests.Hoverfly
		stateURL string
	)

	BeforeEach(func() {
		hoverfly = functional_tests.NewHoverfly()
		hoverfly.Start()
		hoverfly.SetMode("simulate")
		stateURL = "http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/state"
	})

	AfterEach(func() {
		hoverfly.Stop()
	})

	It("should be able to transition through states", func() {
		hoverfly.ImportSimulation(testdata.StatePayload)

		GinkgoWriter.Write([]byte(`hello world`))

		assertState(stateURL, `{"state":{}}`)

		resp := hoverfly.Proxy(sling.New().Get("http://test-server.com/basket"))
		Expect(ioutil.ReadAll(resp.Body)).To(Equal([]byte(`empty`)))

		resp = hoverfly.Proxy(sling.New().Get("http://test-server.com/add-eggs"))
		Expect(ioutil.ReadAll(resp.Body)).To(Equal([]byte(`added eggs`)))

		assertState(stateURL, `{"state": {"eggs":"present"}}`)

		foo := hoverfly.Proxy(sling.New().Get("http://test-server.com/basket"))

		bytes, _ := ioutil.ReadAll(foo.Body)

		Expect(bytes).To(Equal([]byte(`eggs`)))

		resp = hoverfly.Proxy(sling.New().Get("http://test-server.com/add-bacon"))
		Expect(ioutil.ReadAll(resp.Body)).To(Equal([]byte(`added bacon`)))

		assertState(stateURL, `{"state":{"eggs":"present","bacon":"present"}}`)

		resp = hoverfly.Proxy(sling.New().Get("http://test-server.com/basket"))
		Expect(ioutil.ReadAll(resp.Body)).To(Equal([]byte(`eggs, bacon`)))

		resp = hoverfly.Proxy(sling.New().Get("http://test-server.com/remove-eggs"))
		Expect(ioutil.ReadAll(resp.Body)).To(Equal([]byte(`removed eggs`)))

		assertState(stateURL, `{"state": {"bacon":"present"}}`)

		resp = hoverfly.Proxy(sling.New().Get("http://test-server.com/basket"))
		Expect(ioutil.ReadAll(resp.Body)).To(Equal([]byte(`bacon`)))

		resp = hoverfly.Proxy(sling.New().Get("http://test-server.com/remove-bacon"))
		Expect(ioutil.ReadAll(resp.Body)).To(Equal([]byte(`removed bacon`)))

		assertState(stateURL, `{"state":{}}`)

		resp = hoverfly.Proxy(sling.New().Get("http://test-server.com/basket"))
		Expect(ioutil.ReadAll(resp.Body)).To(Equal([]byte(`empty`)))

		// Repeat it all to make sure caching has not broken anything

		assertState(stateURL, `{"state":{}}`)

		resp = hoverfly.Proxy(sling.New().Get("http://test-server.com/basket"))
		Expect(ioutil.ReadAll(resp.Body)).To(Equal([]byte(`empty`)))

		resp = hoverfly.Proxy(sling.New().Get("http://test-server.com/add-eggs"))
		Expect(ioutil.ReadAll(resp.Body)).To(Equal([]byte(`added eggs`)))

		assertState(stateURL, `{"state":{"eggs":"present"}}`)

		GinkgoWriter.Write([]byte(`18`))

		resp = hoverfly.Proxy(sling.New().Get("http://test-server.com/basket"))
		Expect(ioutil.ReadAll(resp.Body)).To(Equal([]byte(`eggs`)))

		GinkgoWriter.Write([]byte(`19`))

		resp = hoverfly.Proxy(sling.New().Get("http://test-server.com/add-bacon"))
		Expect(ioutil.ReadAll(resp.Body)).To(Equal([]byte(`added bacon`)))

		GinkgoWriter.Write([]byte(`20`))

		assertState(stateURL, `{"state": {"eggs":"present","bacon":"present"}}`)

		GinkgoWriter.Write([]byte(`21`))

		resp = hoverfly.Proxy(sling.New().Get("http://test-server.com/basket"))
		Expect(ioutil.ReadAll(resp.Body)).To(Equal([]byte(`eggs, bacon`)))

		GinkgoWriter.Write([]byte(`22`))

		resp = hoverfly.Proxy(sling.New().Get("http://test-server.com/remove-eggs"))
		Expect(ioutil.ReadAll(resp.Body)).To(Equal([]byte(`removed eggs`)))

		GinkgoWriter.Write([]byte(`23`))

		assertState(stateURL, `{"state":{"bacon":"present"}}`)

		GinkgoWriter.Write([]byte(`24`))

		resp = hoverfly.Proxy(sling.New().Get("http://test-server.com/basket"))
		Expect(ioutil.ReadAll(resp.Body)).To(Equal([]byte(`bacon`)))

		GinkgoWriter.Write([]byte(`25`))

		resp = hoverfly.Proxy(sling.New().Get("http://test-server.com/remove-bacon"))
		Expect(ioutil.ReadAll(resp.Body)).To(Equal([]byte(`removed bacon`)))

		GinkgoWriter.Write([]byte(`26`))

		assertState(stateURL, `{"state":{}}`)

		GinkgoWriter.Write([]byte(`27`))

		resp = hoverfly.Proxy(sling.New().Get("http://test-server.com/basket"))
		Expect(ioutil.ReadAll(resp.Body)).To(Equal([]byte(`empty`)))
	})
})

func assertState(stateURL, expectedState string) {
	req := sling.New().Get(stateURL)
	res := functional_tests.DoRequest(req)
	Expect(ioutil.ReadAll(res.Body)).To(MatchJSON(expectedState))
	Expect(res.StatusCode).To(Equal(200))
}

package hoverfly_test

import (
	"io/ioutil"

	"fmt"

	"github.com/SpectoLabs/hoverfly/functional-tests"
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
		hoverfly.ImportSimulation(functional_tests.StatePayload)

		GinkgoWriter.Write([]byte(`hello world`))

		assertState(stateURL, `{"state":{}}`)

		GinkgoWriter.Write([]byte(`1`))

		resp := hoverfly.Proxy(sling.New().Get("http://test-server.com/basket"))
		Expect(ioutil.ReadAll(resp.Body)).To(Equal([]byte(`empty`)))

		GinkgoWriter.Write([]byte(`2`))

		resp = hoverfly.Proxy(sling.New().Get("http://test-server.com/add-eggs"))
		Expect(ioutil.ReadAll(resp.Body)).To(Equal([]byte(`added eggs`)))

		GinkgoWriter.Write([]byte(`3`))

		assertState(stateURL, `{"state": {"eggs":"present"}}`)

		GinkgoWriter.Write([]byte(`4`))

		foo := hoverfly.Proxy(sling.New().Get("http://test-server.com/basket"))

		bytes, _ := ioutil.ReadAll(foo.Body)

		fmt.Println(string(bytes))

		GinkgoWriter.Write(bytes)

		Expect(bytes).To(Equal([]byte(`eggs`)))

		GinkgoWriter.Write([]byte(`5`))

		resp = hoverfly.Proxy(sling.New().Get("http://test-server.com/add-bacon"))
		Expect(ioutil.ReadAll(resp.Body)).To(Equal([]byte(`added bacon`)))

		GinkgoWriter.Write([]byte(`6`))

		assertState(stateURL, `{"state":{"eggs":"present","bacon":"present"}}`)

		GinkgoWriter.Write([]byte(`7`))

		resp = hoverfly.Proxy(sling.New().Get("http://test-server.com/basket"))
		Expect(ioutil.ReadAll(resp.Body)).To(Equal([]byte(`eggs, bacon`)))

		GinkgoWriter.Write([]byte(`8`))

		resp = hoverfly.Proxy(sling.New().Get("http://test-server.com/remove-eggs"))
		Expect(ioutil.ReadAll(resp.Body)).To(Equal([]byte(`removed eggs`)))

		GinkgoWriter.Write([]byte(`9`))

		assertState(stateURL, `{"state": {"bacon":"present"}}`)

		GinkgoWriter.Write([]byte(`10`))

		resp = hoverfly.Proxy(sling.New().Get("http://test-server.com/basket"))
		Expect(ioutil.ReadAll(resp.Body)).To(Equal([]byte(`bacon`)))

		GinkgoWriter.Write([]byte(`11`))

		resp = hoverfly.Proxy(sling.New().Get("http://test-server.com/remove-bacon"))
		Expect(ioutil.ReadAll(resp.Body)).To(Equal([]byte(`removed bacon`)))

		GinkgoWriter.Write([]byte(`12`))

		assertState(stateURL, `{"state":{}}`)

		GinkgoWriter.Write([]byte(`13`))

		resp = hoverfly.Proxy(sling.New().Get("http://test-server.com/basket"))
		Expect(ioutil.ReadAll(resp.Body)).To(Equal([]byte(`empty`)))

		GinkgoWriter.Write([]byte(`14`))

		// Repeat it all to make sure caching has not broken anything

		assertState(stateURL, `{"state":{}}`)

		GinkgoWriter.Write([]byte(`15`))

		resp = hoverfly.Proxy(sling.New().Get("http://test-server.com/basket"))
		Expect(ioutil.ReadAll(resp.Body)).To(Equal([]byte(`empty`)))

		GinkgoWriter.Write([]byte(`16`))

		resp = hoverfly.Proxy(sling.New().Get("http://test-server.com/add-eggs"))
		Expect(ioutil.ReadAll(resp.Body)).To(Equal([]byte(`added eggs`)))

		GinkgoWriter.Write([]byte(`17`))

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

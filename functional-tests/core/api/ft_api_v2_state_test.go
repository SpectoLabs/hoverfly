package api_test

import (
	"io/ioutil"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Working with the current state of Hoverfly via the API", func() {

	var (
		hoverfly *functional_tests.Hoverfly
	)

	BeforeEach(func() {
		hoverfly = functional_tests.NewHoverfly()
	})

	AfterEach(func() {
		hoverfly.Stop()
	})

	It("Should be able to put, get, patch and delete state", func() {

		hoverfly.Start()

		stateURL := "http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/state"

		// PUT
		req := sling.New().Put(stateURL).BodyJSON(v2.StateView{
			State: map[string]string{"foo": "bar", "cheese": "ham"},
		})
		res := functional_tests.DoRequest(req)
		Expect(ioutil.ReadAll(res.Body)).To(MatchJSON(`{"state":{"foo":"bar","cheese":"ham"}}`))
		Expect(res.StatusCode).To(Equal(200))

		// GET RESULTS OF PUT
		req = sling.New().Get(stateURL)
		res = functional_tests.DoRequest(req)
		Expect(res.StatusCode).To(Equal(200))
		Expect(ioutil.ReadAll(res.Body)).To(MatchJSON(`{"state":{"foo":"bar","cheese":"ham"}}`))

		// PATCH
		req = sling.New().Patch(stateURL).BodyJSON(v2.StateView{
			State: map[string]string{"foo": "patched"},
		})
		res = functional_tests.DoRequest(req)
		Expect(res.StatusCode).To(Equal(200))
		Expect(ioutil.ReadAll(res.Body)).To(MatchJSON(`{"state":{"cheese":"ham","foo":"patched"}}`))

		// GET RESULTS OF PATCH
		req = sling.New().Get(stateURL)
		res = functional_tests.DoRequest(req)
		Expect(res.StatusCode).To(Equal(200))
		Expect(ioutil.ReadAll(res.Body)).To(MatchJSON(`{"state":{"cheese":"ham","foo":"patched"}}`))

		// DELETE
		req = sling.New().Delete(stateURL)
		res = functional_tests.DoRequest(req)
		Expect(res.StatusCode).To(Equal(200))

		// GET RESULTS OF DELETE
		req = sling.New().Get(stateURL)
		res = functional_tests.DoRequest(req)
		Expect(res.StatusCode).To(Equal(200))
		Expect(ioutil.ReadAll(res.Body)).To(MatchJSON(`{"state":{}}`))
	})
})

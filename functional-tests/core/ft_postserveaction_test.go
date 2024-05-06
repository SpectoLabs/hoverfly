package hoverfly_test

import (
	"encoding/json"
	v2 "github.com/SpectoLabs/hoverfly/core/handlers/v2"
	functional_tests "github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/SpectoLabs/hoverfly/functional-tests/testdata"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
)

var _ = Describe("Manage post serve actions in hoverfly", func() {

	var (
		hoverfly *functional_tests.Hoverfly
	)

	//var server *httptest.Server

	BeforeEach(func() {
		hoverfly = functional_tests.NewHoverfly()
	})

	AfterEach(func() {
		hoverfly.Stop()
	})

	Context("get post serve action", func() {

		Context("hoverfly with local post-serve-action", func() {

			BeforeEach(func() {
				hoverfly.Start("-post-serve-action", "test-callback python testdata/middleware.py 1300")
			})

			It("Should return post serve action details", func() {
				postServeActionDetails := hoverfly.GetAllPostServeAction()
				Expect(postServeActionDetails).NotTo(BeNil())
				Expect(postServeActionDetails.Actions).To(HaveLen(1))
				Expect(postServeActionDetails.Actions[0].ActionName).To(Equal("test-callback"))
				Expect(postServeActionDetails.Actions[0].Binary).To(Equal("python"))
				Expect(postServeActionDetails.Actions[0].DelayInMs).To(Equal(1300))
			})
		})

		Context("hoverfly with remote post-serve-action", func() {

			BeforeEach(func() {
				hoverfly.Start("-post-serve-action", "test-callback http://localhost:8080 1300")
			})

			It("Should return post serve action details", func() {
				postServeActionDetails := hoverfly.GetAllPostServeAction()
				Expect(postServeActionDetails).NotTo(BeNil())
				Expect(postServeActionDetails.Actions).To(HaveLen(1))
				Expect(postServeActionDetails.Actions[0].ActionName).To(Equal("test-callback"))
				Expect(postServeActionDetails.Actions[0].Remote).To(Equal("http://localhost:8080"))
				Expect(postServeActionDetails.Actions[0].DelayInMs).To(Equal(1300))
			})
		})
	})

	Context("set local post serve action", func() {

		Context("start hoverfly and set post serve action", func() {

			BeforeEach(func() {
				hoverfly.Start()
			})

			AfterEach(func() {
				hoverfly.Stop()
			})

			It("Should set post serve action", func() {
				postServeActionDetails := hoverfly.SetLocalPostServeAction("testing", "python3", "dummy-script", 1400)
				Expect(postServeActionDetails).NotTo(BeNil())
				Expect(postServeActionDetails.Actions).To(HaveLen(1))
				Expect(postServeActionDetails.Actions[0].ActionName).To(Equal("testing"))
				Expect(postServeActionDetails.Actions[0].Binary).To(Equal("python3"))
				Expect(postServeActionDetails.Actions[0].ScriptContent).To(Equal("dummy-script"))
				Expect(postServeActionDetails.Actions[0].DelayInMs).To(Equal(1400))
			})
		})

		Context("start hoverfly and set remote post serve action", func() {

			BeforeEach(func() {
				hoverfly.Start()
			})

			AfterEach(func() {
				hoverfly.Stop()
			})

			It("Should set post serve action", func() {
				postServeActionDetails := hoverfly.SetRemotePostServeAction("testing", "http://localhost", 1400)
				Expect(postServeActionDetails).NotTo(BeNil())
				Expect(postServeActionDetails.Actions).To(HaveLen(1))
				Expect(postServeActionDetails.Actions[0].ActionName).To(Equal("testing"))
				Expect(postServeActionDetails.Actions[0].Remote).To(Equal("http://localhost"))
				Expect(postServeActionDetails.Actions[0].DelayInMs).To(Equal(1400))
			})
		})

		Context("start hoverfly and set fallback remote post serve action", func() {

			BeforeEach(func() {
				hoverfly.Start()
			})

			AfterEach(func() {
				hoverfly.Stop()
			})

			It("Should set post serve action", func() {
				postServeActionDetails := hoverfly.SetRemotePostServeAction("", "http://localhost:8080", 1600)
				Expect(postServeActionDetails).NotTo(BeNil())
				Expect(postServeActionDetails.Actions).NotTo(BeNil())
				Expect(postServeActionDetails.Actions[0].ActionName).To(Equal(""))
				Expect(postServeActionDetails.Actions[0].Remote).To(Equal("http://localhost:8080"))
				Expect(postServeActionDetails.Actions[0].DelayInMs).To(Equal(1600))
			})
		})
	})

	Context("delete post serve action", func() {

		Context("start local post serve acton and delete it", func() {

			BeforeEach(func() {
				hoverfly.Start("-post-serve-action", "test-callback python testdata/middleware.py 1300")
			})

			It("Should return empty post serve action details on deletion", func() {
				postServeActionDetails := hoverfly.DeletePostServeAction("test-callback")
				Expect(postServeActionDetails).NotTo(BeNil())
				Expect(postServeActionDetails.Actions).To(HaveLen(0))
			})
		})

		Context("start remote post serve acton and delete it", func() {

			BeforeEach(func() {
				hoverfly.Start("-post-serve-action", "test-callback http://localhost:8080 1300")
			})

			It("Should return empty post serve action details on deletion", func() {
				postServeActionDetails := hoverfly.DeletePostServeAction("test-callback")
				Expect(postServeActionDetails).NotTo(BeNil())
				Expect(postServeActionDetails.Actions).To(HaveLen(0))
			})
		})
	})

	Context("set post serve action in simulation", func() {

		Context("start hoverfly with local post serve action and set in simulation", func() {

			BeforeEach(func() {
				hoverfly.Start("-post-serve-action", "test-callback python testdata/middleware.py 1300")
			})

			AfterEach(func() {
				hoverfly.Stop()
			})

			It("Should be able to set post-serve-action in simulation", func() {
				hoverfly.ImportSimulation(testdata.SimulationWithPostServeAction)
				body := hoverfly.GetSimulation()

				simulationBytes, err := ioutil.ReadAll(body)
				Expect(err).To(BeNil())

				simulationView := &v2.SimulationViewV5{}
				err = json.Unmarshal(simulationBytes, simulationView)
				Expect(err).To(BeNil())
				Expect(simulationView).NotTo(BeNil())
				Expect(simulationView.DataViewV5.RequestResponsePairs).To(HaveLen(1))
				Expect(simulationView.DataViewV5.RequestResponsePairs[0].Response.PostServeAction).To(Equal("test-callback"))

			})
		})

		Context("start hoverfly with remote post serve action and set in simulation", func() {

			BeforeEach(func() {
				hoverfly.Start("-post-serve-action", "test-callback http://localhost:8080 1300")
			})

			AfterEach(func() {
				hoverfly.Stop()
			})

			It("Should be able to set post-serve-action in simulation", func() {
				hoverfly.ImportSimulation(testdata.SimulationWithPostServeAction)
				body := hoverfly.GetSimulation()

				simulationBytes, err := ioutil.ReadAll(body)
				Expect(err).To(BeNil())

				simulationView := &v2.SimulationViewV5{}
				err = json.Unmarshal(simulationBytes, simulationView)
				Expect(err).To(BeNil())
				Expect(simulationView).NotTo(BeNil())
				Expect(simulationView.DataViewV5.RequestResponsePairs).To(HaveLen(1))
				Expect(simulationView.DataViewV5.RequestResponsePairs[0].Response.PostServeAction).To(Equal("test-callback"))

			})
		})
	})
})

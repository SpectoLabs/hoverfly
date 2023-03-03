package hoverfly_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	v2 "github.com/SpectoLabs/hoverfly/core/handlers/v2"
	functional_tests "github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/SpectoLabs/hoverfly/functional-tests/testdata"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Manage post serve actions in hoverfly", func() {

	var (
		hoverfly *functional_tests.Hoverfly
	)

	BeforeEach(func() {
		hoverfly = functional_tests.NewHoverfly()
	})

	AfterEach(func() {
		hoverfly.Stop()
	})

	Context("get post serve action", func() {

		Context("hoverfly with post-serve-action", func() {

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
	})

	Context("set post serve action", func() {

		Context("start hoverfly and set post serve action", func() {

			BeforeEach(func() {
				hoverfly.Start()
			})

			AfterEach(func() {
				hoverfly.Stop()
			})

			It("Should set post serve action", func() {
				postServeActionDetails := hoverfly.SetPostServeAction("testing", "python3", "dummy-script", 1400)
				Expect(postServeActionDetails).NotTo(BeNil())
				Expect(postServeActionDetails.Actions).To(HaveLen(1))
				Expect(postServeActionDetails.Actions[0].ActionName).To(Equal("testing"))
				Expect(postServeActionDetails.Actions[0].Binary).To(Equal("python3"))
				Expect(postServeActionDetails.Actions[0].ScriptContent).To(Equal("dummy-script"))
				Expect(postServeActionDetails.Actions[0].DelayInMs).To(Equal(1400))
			})
		})
	})

	Context("delete post serve action", func() {

		Context("start post serve acton and delete it", func() {

			BeforeEach(func() {
				hoverfly.Start("-post-serve-action", "test-callback python testdata/middleware.py 1300")
			})

			It("Should return empty post serve action details on deletion", func() {
				postServeActionDetails := hoverfly.DeletePostServeAction("test-callback")
				Expect(postServeActionDetails).NotTo(BeNil())
				Expect(postServeActionDetails.Actions).To(HaveLen(0))
			})
		})
	})

	Context("set post serve action in simulation", func() {

		Context("start hoverfly with post serve action and set in simulation", func() {

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
				fmt.Println(simulationView)
				Expect(err).To(BeNil())
				Expect(simulationView).NotTo(BeNil())
				Expect(simulationView.DataViewV5.RequestResponsePairs).To(HaveLen(1))
				Expect(simulationView.DataViewV5.RequestResponsePairs[0].Response.PostServeAction).To(Equal("test-callback"))

			})
		})
	})
})

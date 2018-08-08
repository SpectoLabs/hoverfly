package hoverfly

import (
	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/SpectoLabs/hoverfly/functional-tests/testdata"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Running Hoverfly with older simulations", func() {

	var (
		hoverfly *functional_tests.Hoverfly
	)

	BeforeEach(func() {
		hoverfly = functional_tests.NewHoverfly()
	})

	AfterEach(func() {
		hoverfly.Stop()
	})

	Context("v1 simulations", func() {

		BeforeEach(func() {
			hoverfly.Start()
		})

		It("should upgrade it to the latest simulation", func() {
			hoverfly.ImportSimulation(testdata.V1JsonPayload)
			upgradedSimulation := hoverfly.ExportSimulation()

			simulation := v2.SimulationViewV5{}

			functional_tests.Unmarshal([]byte(testdata.V5JsonPayload), &simulation)

			Expect(upgradedSimulation.DataViewV5).To(Equal(simulation.DataViewV5))
		})
	})

	Context("v3 simulations", func() {

		BeforeEach(func() {
			hoverfly.Start()
		})

		It("should upgrade it to the latest simulation", func() {
			hoverfly.ImportSimulation(testdata.V3Delays)
			upgradedSimulation := hoverfly.ExportSimulation()

			simulation := v2.SimulationViewV5{}

			functional_tests.Unmarshal([]byte(testdata.Delays), &simulation)

			Expect(upgradedSimulation.DataViewV5).To(Equal(simulation.DataViewV5))
		})

		It("should upgrade it to the latest simulation", func() {
			hoverfly.ImportSimulation(testdata.V3ClosestMissProof)
			upgradedSimulation := hoverfly.ExportSimulation()

			simulation := v2.SimulationViewV5{}

			functional_tests.Unmarshal([]byte(testdata.ClosestMissProof), &simulation)

			Expect(upgradedSimulation.DataViewV5).To(Equal(simulation.DataViewV5))
		})

		It("should upgrade it to the latest simulation", func() {
			hoverfly.ImportSimulation(testdata.V3ExactMatch)
			upgradedSimulation := hoverfly.ExportSimulation()

			simulation := v2.SimulationViewV5{}

			functional_tests.Unmarshal([]byte(testdata.ExactMatch), &simulation)

			Expect(upgradedSimulation.DataViewV5).To(Equal(simulation.DataViewV5))
		})

		It("should upgrade it to the latest simulation", func() {
			hoverfly.ImportSimulation(testdata.V3GlobMatch)
			upgradedSimulation := hoverfly.ExportSimulation()

			simulation := v2.SimulationViewV5{}

			functional_tests.Unmarshal([]byte(testdata.GlobMatch), &simulation)

			Expect(upgradedSimulation.DataViewV5).To(Equal(simulation.DataViewV5))
		})

		It("should upgrade it to the latest simulation", func() {
			hoverfly.ImportSimulation(testdata.V3XmlMatch)
			upgradedSimulation := hoverfly.ExportSimulation()

			simulation := v2.SimulationViewV5{}

			functional_tests.Unmarshal([]byte(testdata.XmlMatch), &simulation)

			Expect(upgradedSimulation.DataViewV5).To(Equal(simulation.DataViewV5))
		})

		It("should upgrade it to the latest simulation", func() {
			hoverfly.ImportSimulation(testdata.V3XpathMatch)
			upgradedSimulation := hoverfly.ExportSimulation()

			simulation := v2.SimulationViewV5{}

			functional_tests.Unmarshal([]byte(testdata.XpathMatch), &simulation)

			Expect(upgradedSimulation.DataViewV5).To(Equal(simulation.DataViewV5))
		})
	})

	Context("v4 simulations", func() {

		BeforeEach(func() {
			hoverfly.Start()
		})

		It("should upgrade it to the latest simulation", func() {
			hoverfly.ImportSimulation(testdata.V4QueryMatchers)
			upgradedSimulation := hoverfly.ExportSimulation()

			simulation := v2.SimulationViewV5{}

			functional_tests.Unmarshal([]byte(testdata.QueryMatchers), &simulation)

			Expect(upgradedSimulation.DataViewV5).To(Equal(simulation.DataViewV5))
		})

		It("should upgrade it to the latest simulation", func() {
			hoverfly.ImportSimulation(testdata.V4HeaderMatchers)
			upgradedSimulation := hoverfly.ExportSimulation()

			simulation := v2.SimulationViewV5{}

			functional_tests.Unmarshal([]byte(testdata.HeaderMatchers), &simulation)

			Expect(upgradedSimulation.DataViewV5).To(Equal(simulation.DataViewV5))
		})
	})
})

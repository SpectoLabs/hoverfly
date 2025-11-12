package hoverfly_test

import (
	"encoding/json"
	"io/ioutil"

	v2 "github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/SpectoLabs/hoverfly/functional-tests/testdata"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("When I run Hoverfly", func() {

	var (
		hoverfly *functional_tests.Hoverfly
	)

	BeforeEach(func() {
		hoverfly = functional_tests.NewHoverfly()
	})

	AfterEach(func() {
		hoverfly.Stop()
	})

	Context("with -no-import-check", func() {

		BeforeEach(func() {
			hoverfly.Start("-no-import-check")
		})

		It("should skip duplicate pair check", func() {

			hoverfly.ImportSimulation(testdata.DuplicatePairs)

			recordsJson, err := ioutil.ReadAll(hoverfly.GetSimulation())
			Expect(err).To(BeNil())

			payload := v2.SimulationViewV5{}

			Expect(json.Unmarshal(recordsJson, &payload)).To(Succeed())
			Expect(payload.RequestResponsePairs).To(HaveLen(2))
		})
	})
})

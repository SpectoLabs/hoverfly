package main

import (
	. "github.com/onsi/gomega"
	"testing"
)

func Test_NewSimulation_CanCreateASimulationFromCompleteKey(t *testing.T) {
	RegisterTestingT(t)

	simulation, err := NewSimulation("testvendor/testname:v1")

	Expect(err).To(BeNil())

	Expect(simulation.Vendor).To(Equal("testvendor"))
	Expect(simulation.Name).To(Equal("testname"))
	Expect(simulation.Version).To(Equal("v1"))
}

func Test_NewSimulation_CanCreateASimulationFromDifferentCompleteKey(t *testing.T) {
	RegisterTestingT(t)

	simulation, err := NewSimulation("another-vendor/test_simulation:v7")

	Expect(err).To(BeNil())

	Expect(simulation.Vendor).To(Equal("another-vendor"))
	Expect(simulation.Name).To(Equal("test_simulation"))
	Expect(simulation.Version).To(Equal("v7"))
}

func Test_NewSimulation_CanCreateASimulationFromKey_WithNoVersion(t *testing.T) {
	RegisterTestingT(t)

	simulation, err := NewSimulation("tester/tested")

	Expect(err).To(BeNil())

	Expect(simulation.Vendor).To(Equal("tester"))
	Expect(simulation.Name).To(Equal("tested"))
	Expect(simulation.Version).To(Equal("latest"))
}

func Test_NewSimulation_CanCreateASimulationFromKey_WithNoVendor(t *testing.T) {
	RegisterTestingT(t)

	simulation, err := NewSimulation("just_a-name")
	Expect(err).To(BeNil())

	Expect(simulation.Vendor).To(Equal(""))
	Expect(simulation.Name).To(Equal("just_a-name"))
	Expect(simulation.Version).To(Equal("latest"))
}

func Test_NewSimulation_WontCreateASimulationFromKey_WithSpecialCharacters(t *testing.T) {
	RegisterTestingT(t)

	simulation, err := NewSimulation("just_@-name")
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("Invalid characters used in simulation name"))
	Expect(simulation).To(Equal(Simulation{}))

	simulation, err = NewSimulation("just_\\-name")
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("Invalid characters used in simulation name"))
	Expect(simulation).To(Equal(Simulation{}))

	simulation, err = NewSimulation("just()an&simulation")
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("Invalid characters used in simulation name"))
	Expect(simulation).To(Equal(Simulation{}))

	simulation, err = NewSimulation("just()an£im%lation")
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("Invalid characters used in simulation name"))
	Expect(simulation).To(Equal(Simulation{}))
}

func Test_Simulation_GetFileName(t *testing.T) {
	RegisterTestingT(t)

	simulation := Simulation{
		Vendor: "vendor",
		Name: "name",
		Version: "version",
	}

	resultFileName := simulation.GetFileName()
	Expect(resultFileName).To(Equal("vendor.name.version.hfile"))
}

func Test_Simulation_String(t *testing.T) {
	RegisterTestingT(t)

	simulation := Simulation{
		Vendor: "vendor",
		Name: "name",
		Version: "version",
	}

	resultString := simulation.String()
	Expect(resultString).To(Equal("vendor/name:version"))

}
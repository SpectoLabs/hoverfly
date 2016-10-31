package hoverfly

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/models"
	. "github.com/onsi/gomega"
	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/util"
	"github.com/SpectoLabs/hoverfly/core/handlers/v1"
)

var (
	pairOne = v2.RequestResponsePairView{
		Request: v2.RequestDetailsView{
			Destination: util.StringToPointer("test.com"),
			Path:        util.StringToPointer("/testing"),
		},
		Response: v2.ResponseDetailsView{
			Body: "test-body",
		},
	}
)

func TestHoverflyGetSimulationReturnsBlankSimulation_ifThereIsNoData(t *testing.T) {
	RegisterTestingT(t)

	server, unit := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	simulation, err := unit.GetSimulation()
	Expect(err).To(BeNil())

	Expect(simulation.DataView.RequestResponsePairs).To(HaveLen(0))
	Expect(simulation.DataView.GlobalActions.Delays).To(HaveLen(0))

	Expect(simulation.MetaView.SchemaVersion).To(Equal("v1"))
	Expect(simulation.MetaView.HoverflyVersion).To(Equal("v0.9.0"))
	Expect(simulation.MetaView.TimeExported).ToNot(BeNil())
}

func TestHoverflyGetSimulationReturnsASingleRequestResponsePair(t *testing.T) {
	RegisterTestingT(t)

	server, unit := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	recording := models.RequestResponsePair{
		Request: models.RequestDetails{
			Destination: "testhost.com",
			Path:        "/test",
		},
		Response: models.ResponseDetails{
			Status: 200,
			Body:   "test",
		},
	}

	recordingBytes, err := recording.Encode()
	Expect(err).To(BeNil())

	unit.RequestCache.Set([]byte("key"), recordingBytes)

	simulation, err := unit.GetSimulation()
	Expect(err).To(BeNil())

	Expect(simulation.DataView.RequestResponsePairs).To(HaveLen(1))

	Expect(*simulation.DataView.RequestResponsePairs[0].Request.Destination).To(Equal("testhost.com"))
	Expect(*simulation.DataView.RequestResponsePairs[0].Request.Path).To(Equal("/test"))
	Expect(*simulation.DataView.RequestResponsePairs[0].Request.RequestType).To(Equal("recording"))

	Expect(simulation.DataView.RequestResponsePairs[0].Response.Status).To(Equal(200))
	Expect(simulation.DataView.RequestResponsePairs[0].Response.Body).To(Equal("test"))

	Expect(nil).To(BeNil())
}

func TestHoverflyGetSimulationReturnsMultipleRequestResponsePairs(t *testing.T) {
	RegisterTestingT(t)

	server, unit := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	recording := models.RequestResponsePair{
		Request: models.RequestDetails{
			Destination: "testhost.com",
			Path:        "/test",
		},
		Response: models.ResponseDetails{
			Status: 200,
			Body:   "test",
		},
	}

	recordingBytes, err := recording.Encode()
	Expect(err).To(BeNil())

	unit.RequestCache.Set([]byte("key"), recordingBytes)
	unit.RequestCache.Set([]byte("key2"), recordingBytes)

	simulation, err := unit.GetSimulation()
	Expect(err).To(BeNil())

	Expect(simulation.DataView.RequestResponsePairs).To(HaveLen(2))

	Expect(*simulation.DataView.RequestResponsePairs[0].Request.Destination).To(Equal("testhost.com"))
	Expect(*simulation.DataView.RequestResponsePairs[0].Request.Path).To(Equal("/test"))
	Expect(*simulation.DataView.RequestResponsePairs[0].Request.RequestType).To(Equal("recording"))

	Expect(simulation.DataView.RequestResponsePairs[0].Response.Status).To(Equal(200))
	Expect(simulation.DataView.RequestResponsePairs[0].Response.Body).To(Equal("test"))

	Expect(*simulation.DataView.RequestResponsePairs[1].Request.Destination).To(Equal("testhost.com"))
	Expect(*simulation.DataView.RequestResponsePairs[1].Request.Path).To(Equal("/test"))
	Expect(*simulation.DataView.RequestResponsePairs[1].Request.RequestType).To(Equal("recording"))

	Expect(simulation.DataView.RequestResponsePairs[1].Response.Status).To(Equal(200))
	Expect(simulation.DataView.RequestResponsePairs[1].Response.Body).To(Equal("test"))
}

func TestHoverflyGetSimulationReturnsMultipleDelays(t *testing.T) {
	RegisterTestingT(t)

	server, unit := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	delay1 := models.ResponseDelay{
		UrlPattern: "test-pattern",
		Delay:      100,
	}

	delay2 := models.ResponseDelay{
		HttpMethod: "test",
		Delay:      200,
	}

	responseDelays := models.ResponseDelayList{delay1, delay2}

	unit.ResponseDelays = &responseDelays

	simulation, err := unit.GetSimulation()
	Expect(err).To(BeNil())

	Expect(simulation.DataView.GlobalActions.Delays).To(HaveLen(2))

	Expect(simulation.DataView.GlobalActions.Delays[0].UrlPattern).To(Equal("test-pattern"))
	Expect(simulation.DataView.GlobalActions.Delays[0].HttpMethod).To(Equal(""))
	Expect(simulation.DataView.GlobalActions.Delays[0].Delay).To(Equal(100))

	Expect(simulation.DataView.GlobalActions.Delays[1].UrlPattern).To(Equal(""))
	Expect(simulation.DataView.GlobalActions.Delays[1].HttpMethod).To(Equal("test"))
	Expect(simulation.DataView.GlobalActions.Delays[1].Delay).To(Equal(200))
}

func TestHoverfly_PutSimulation_ImportsRecordings(t *testing.T) {
	RegisterTestingT(t)

	server, unit := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	simulationToImport := v2.SimulationView {
		DataView: v2.DataView {
			RequestResponsePairs: []v2.RequestResponsePairView{pairOne},
			GlobalActions: v2.GlobalActionsView {
				Delays: []v1.ResponseDelayView {},
			},
		},
		MetaView: v2.MetaView {},
	}

	unit.PutSimulation(simulationToImport)

	importedSimulation, err := unit.GetSimulation()
	Expect(err).To(BeNil())

	Expect(importedSimulation.RequestResponsePairs).To(HaveLen(1))
}

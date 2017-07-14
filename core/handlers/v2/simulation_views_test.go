package v2_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	. "github.com/onsi/gomega"
)

func Test_NewSimulationViewFromResponseBody_CanCreateSimulationFromV3Payload(t *testing.T) {
	RegisterTestingT(t)

	simulation, err := v2.NewSimulationViewFromResponseBody([]byte(`{
		"data": {
			"pairs": [
				{
					"response": {
						"status": 200,
						"body": "exact match",
						"encodedBody": false,
						"headers": {
							"Header": [
								"value"
							]
						},
						"templated":true
					},
					"request": {
						"destination": {
							"exactMatch": "test-server.com"
						}
					}
				}
			],
			"globalActions": {
				"delays": []
			}
		},
		"meta": {
			"schemaVersion": "v3",
			"hoverflyVersion": "v0.11.0",
			"timeExported": "2017-02-23T12:43:48Z"
		}
	}`))

	Expect(err).To(BeNil())

	Expect(simulation.RequestResponsePairs).To(HaveLen(1))

	Expect(simulation.RequestResponsePairs[0].RequestMatcher.Body).To(BeNil())
	Expect(*simulation.RequestResponsePairs[0].RequestMatcher.Destination.ExactMatch).To(Equal("test-server.com"))
	Expect(simulation.RequestResponsePairs[0].RequestMatcher.Headers).To(BeNil())
	Expect(simulation.RequestResponsePairs[0].RequestMatcher.Method).To(BeNil())
	Expect(simulation.RequestResponsePairs[0].RequestMatcher.Path).To(BeNil())
	Expect(simulation.RequestResponsePairs[0].RequestMatcher.Query).To(BeNil())
	Expect(simulation.RequestResponsePairs[0].RequestMatcher.Scheme).To(BeNil())

	Expect(simulation.RequestResponsePairs[0].Response.Body).To(Equal("exact match"))
	Expect(simulation.RequestResponsePairs[0].Response.Templated).To(BeTrue())
	Expect(simulation.RequestResponsePairs[0].Response.EncodedBody).To(BeFalse())
	Expect(simulation.RequestResponsePairs[0].Response.Headers).To(HaveKeyWithValue("Header", []string{"value"}))
	Expect(simulation.RequestResponsePairs[0].Response.Status).To(Equal(200))

	Expect(simulation.SchemaVersion).To(Equal("v3"))
	Expect(simulation.HoverflyVersion).To(Equal("v0.11.0"))
	Expect(simulation.TimeExported).To(Equal("2017-02-23T12:43:48Z"))
}

func Test_NewSimulationViewFromResponseBody_CanCreateSimulationFromV2Payload(t *testing.T) {
	RegisterTestingT(t)

	simulation, err := v2.NewSimulationViewFromResponseBody([]byte(`{
		"data": {
			"pairs": [
				{
					"response": {
						"status": 200,
						"body": "exact match",
						"encodedBody": false,
						"headers": {
							"Header": [
								"value"
							]
						}
					},
					"request": {
						"destination": {
							"exactMatch": "test-server.com"
						}
					}
				}
			],
			"globalActions": {
				"delays": []
			}
		},
		"meta": {
			"schemaVersion": "v2",
			"hoverflyVersion": "v0.11.0",
			"timeExported": "2017-02-23T12:43:48Z"
		}
	}`))

	Expect(err).To(BeNil())

	Expect(simulation.RequestResponsePairs).To(HaveLen(1))

	Expect(simulation.RequestResponsePairs[0].RequestMatcher.Body).To(BeNil())
	Expect(*simulation.RequestResponsePairs[0].RequestMatcher.Destination.ExactMatch).To(Equal("test-server.com"))
	Expect(simulation.RequestResponsePairs[0].RequestMatcher.Headers).To(BeNil())
	Expect(simulation.RequestResponsePairs[0].RequestMatcher.Method).To(BeNil())
	Expect(simulation.RequestResponsePairs[0].RequestMatcher.Path).To(BeNil())
	Expect(simulation.RequestResponsePairs[0].RequestMatcher.Query).To(BeNil())
	Expect(simulation.RequestResponsePairs[0].RequestMatcher.Scheme).To(BeNil())

	Expect(simulation.RequestResponsePairs[0].Response.Body).To(Equal("exact match"))
	Expect(simulation.RequestResponsePairs[0].Response.Templated).To(BeFalse())
	Expect(simulation.RequestResponsePairs[0].Response.EncodedBody).To(BeFalse())
	Expect(simulation.RequestResponsePairs[0].Response.Headers).To(HaveKeyWithValue("Header", []string{"value"}))
	Expect(simulation.RequestResponsePairs[0].Response.Status).To(Equal(200))

	Expect(simulation.SchemaVersion).To(Equal("v3"))
	Expect(simulation.HoverflyVersion).To(Equal("v0.11.0"))
	Expect(simulation.TimeExported).To(Equal("2017-02-23T12:43:48Z"))
}

func Test_NewSimulationViewFromResponseBody_WontCreateSimulationIfThereIsNoSchemaVersion(t *testing.T) {
	RegisterTestingT(t)

	simulation, err := v2.NewSimulationViewFromResponseBody([]byte(`{
		"data": {},
		"meta": {
			"hoverflyVersion": "v0.11.0",
			"timeExported": "2017-02-23T12:43:48Z"
		}
	}`))

	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("Invalid JSON, missing \"meta.schemaVersion\" string"))

	Expect(simulation).ToNot(BeNil())
	Expect(simulation.RequestResponsePairs).To(HaveLen(0))
	Expect(simulation.GlobalActions.Delays).To(HaveLen(0))
}

func Test_NewSimulationViewFromResponseBody_WontBlowUpIfMetaIsMissing(t *testing.T) {
	RegisterTestingT(t)

	simulation, err := v2.NewSimulationViewFromResponseBody([]byte(`{
		"data": {}
	}`))

	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal(`Invalid JSON, missing "meta" object`))

	Expect(simulation).ToNot(BeNil())
	Expect(simulation.RequestResponsePairs).To(HaveLen(0))
	Expect(simulation.GlobalActions.Delays).To(HaveLen(0))
}

func Test_NewSimulationViewFromResponseBody_CanCreateSimulationFromV1Payload(t *testing.T) {
	RegisterTestingT(t)

	simulation, err := v2.NewSimulationViewFromResponseBody([]byte(`{
		"data": {
			"pairs": [
				{
					"response": {
						"status": 200,
						"body": "exact match",
						"encodedBody": false,
						"headers": {
							"Header": [
								"value"
							]
						}
					},
					"request": {
						"destination":"test-server.com"
					}
				}
			],
			"globalActions": {
				"delays": []
			}
		},
		"meta": {
			"schemaVersion": "v1",
			"hoverflyVersion": "v0.11.0",
			"timeExported": "2017-02-23T12:43:48Z"
		}
	}`))

	Expect(err).To(BeNil())

	Expect(simulation.RequestResponsePairs).To(HaveLen(1))

	Expect(simulation.RequestResponsePairs[0].RequestMatcher.Body).To(BeNil())
	Expect(*simulation.RequestResponsePairs[0].RequestMatcher.Destination.ExactMatch).To(Equal("test-server.com"))
	Expect(simulation.RequestResponsePairs[0].RequestMatcher.Headers).To(BeNil())
	Expect(simulation.RequestResponsePairs[0].RequestMatcher.Method).To(BeNil())
	Expect(simulation.RequestResponsePairs[0].RequestMatcher.Path).To(BeNil())
	Expect(simulation.RequestResponsePairs[0].RequestMatcher.Query).To(BeNil())
	Expect(simulation.RequestResponsePairs[0].RequestMatcher.Scheme).To(BeNil())

	Expect(simulation.RequestResponsePairs[0].Response.Body).To(Equal("exact match"))
	Expect(simulation.RequestResponsePairs[0].Response.EncodedBody).To(BeFalse())
	Expect(simulation.RequestResponsePairs[0].Response.Headers).To(HaveKeyWithValue("Header", []string{"value"}))
	Expect(simulation.RequestResponsePairs[0].Response.Status).To(Equal(200))
	Expect(simulation.RequestResponsePairs[0].Response.Templated).To(BeFalse())

	Expect(simulation.SchemaVersion).To(Equal("v3"))
	Expect(simulation.HoverflyVersion).To(Equal("v0.11.0"))
	Expect(simulation.TimeExported).To(Equal("2017-02-23T12:43:48Z"))
}

func Test_NewSimulationViewFromResponseBody_WontCreateSimulationFromInvalidV1Simulation(t *testing.T) {
	RegisterTestingT(t)

	simulation, err := v2.NewSimulationViewFromResponseBody([]byte(`{
		"data": {
			"pairs": [
				{
					
				}
			]
		},
		"meta": {
			"schemaVersion": "v1",
			"hoverflyVersion": "v0.11.0",
			"timeExported": "2017-02-23T12:43:48Z"
		}
	}`))

	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("Invalid v1 simulation: request is required, response is required"))

	Expect(simulation).ToNot(BeNil())
	Expect(simulation.RequestResponsePairs).To(HaveLen(0))
	Expect(simulation.GlobalActions.Delays).To(HaveLen(0))
}

func Test_NewSimulationViewFromResponseBody_WontCreateSimulationFromUnknownSchemaVersion(t *testing.T) {
	RegisterTestingT(t)

	_, err := v2.NewSimulationViewFromResponseBody([]byte(`{
		"data": {
			"pairs": [
				{
					
				}
			]
		},
		"meta": {
			"schemaVersion": "r3",
			"hoverflyVersion": "v0.11.0",
			"timeExported": "2017-02-23T12:43:48Z"
		}
	}`))

	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("Invalid simulation: schema version r3 is not supported by this version of Hoverfly, you may need to update Hoverfly"))
}

func Test_NewSimulationViewFromResponseBody_WontCreateSimulationFromInvalidJson(t *testing.T) {
	RegisterTestingT(t)

	simulation, err := v2.NewSimulationViewFromResponseBody([]byte(`{}{}[^.^]{}{}`))

	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("Invalid JSON"))

	Expect(simulation).ToNot(BeNil())
	Expect(simulation.RequestResponsePairs).To(HaveLen(0))
	Expect(simulation.GlobalActions.Delays).To(HaveLen(0))
}

package v2_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	. "github.com/onsi/gomega"
)

func Test_NewSimulationViewFromRequestBody_CanCreateSimulationFromV3Payload(t *testing.T) {
	RegisterTestingT(t)

	simulation, err := v2.NewSimulationViewFromRequestBody([]byte(`{
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
				"delays": [],
				"delaysLogNormal": [
					{"min": 1, "max": 4, "mean": 3, "median" :2}
				]
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

	Expect(simulation.RequestResponsePairs[0].RequestMatcher.Body).To(HaveLen(0))
	Expect(simulation.RequestResponsePairs[0].RequestMatcher.Destination[0].Matcher).To(Equal("exact"))
	Expect(simulation.RequestResponsePairs[0].RequestMatcher.Destination[0].Value).To(Equal("test-server.com"))
	Expect(simulation.RequestResponsePairs[0].RequestMatcher.Headers).To(HaveLen(0))
	Expect(simulation.RequestResponsePairs[0].RequestMatcher.Method).To(HaveLen(0))
	Expect(simulation.RequestResponsePairs[0].RequestMatcher.Path).To(HaveLen(0))
	Expect(simulation.RequestResponsePairs[0].RequestMatcher.DeprecatedQuery).To(HaveLen(0))
	Expect(simulation.RequestResponsePairs[0].RequestMatcher.Scheme).To(HaveLen(0))

	Expect(simulation.RequestResponsePairs[0].Response.Body).To(Equal("exact match"))
	Expect(simulation.RequestResponsePairs[0].Response.Templated).To(BeTrue())
	Expect(simulation.RequestResponsePairs[0].Response.EncodedBody).To(BeFalse())
	Expect(simulation.RequestResponsePairs[0].Response.Headers).To(HaveKeyWithValue("Header", []string{"value"}))
	Expect(simulation.RequestResponsePairs[0].Response.Status).To(Equal(200))

	Expect(simulation.GlobalActions.DelaysLogNormal[0].Min).To(Equal(1))
	Expect(simulation.GlobalActions.DelaysLogNormal[0].Max).To(Equal(4))
	Expect(simulation.GlobalActions.DelaysLogNormal[0].Mean).To(Equal(3))
	Expect(simulation.GlobalActions.DelaysLogNormal[0].Median).To(Equal(2))

	Expect(simulation.SchemaVersion).To(Equal("v3"))
	Expect(simulation.HoverflyVersion).To(Equal("v0.11.0"))
	Expect(simulation.TimeExported).To(Equal("2017-02-23T12:43:48Z"))
}

func Test_NewSimulationViewFromRequestBody_CanCreateSimulationFromV2Payload(t *testing.T) {
	RegisterTestingT(t)

	simulation, err := v2.NewSimulationViewFromRequestBody([]byte(`{
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
				"delays": [],
				"delaysLogNormal": [
					{"min": 1, "max": 4, "mean": 3, "median" :2}
				]
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

	Expect(simulation.RequestResponsePairs[0].RequestMatcher.Body).To(HaveLen(0))
	Expect(simulation.RequestResponsePairs[0].RequestMatcher.Destination[0].Matcher).To(Equal("exact"))
	Expect(simulation.RequestResponsePairs[0].RequestMatcher.Destination[0].Value).To(Equal("test-server.com"))
	Expect(simulation.RequestResponsePairs[0].RequestMatcher.Headers).To(HaveLen(0))
	Expect(simulation.RequestResponsePairs[0].RequestMatcher.Method).To(HaveLen(0))
	Expect(simulation.RequestResponsePairs[0].RequestMatcher.Path).To(HaveLen(0))
	Expect(simulation.RequestResponsePairs[0].RequestMatcher.DeprecatedQuery).To(HaveLen(0))
	Expect(simulation.RequestResponsePairs[0].RequestMatcher.Scheme).To(HaveLen(0))

	Expect(simulation.RequestResponsePairs[0].Response.Body).To(Equal("exact match"))
	Expect(simulation.RequestResponsePairs[0].Response.Templated).To(BeFalse())
	Expect(simulation.RequestResponsePairs[0].Response.EncodedBody).To(BeFalse())
	Expect(simulation.RequestResponsePairs[0].Response.Headers).To(HaveKeyWithValue("Header", []string{"value"}))
	Expect(simulation.RequestResponsePairs[0].Response.Status).To(Equal(200))

	Expect(simulation.GlobalActions.DelaysLogNormal[0].Min).To(Equal(1))
	Expect(simulation.GlobalActions.DelaysLogNormal[0].Max).To(Equal(4))
	Expect(simulation.GlobalActions.DelaysLogNormal[0].Mean).To(Equal(3))
	Expect(simulation.GlobalActions.DelaysLogNormal[0].Median).To(Equal(2))

	Expect(simulation.SchemaVersion).To(Equal("v3"))
	Expect(simulation.HoverflyVersion).To(Equal("v0.11.0"))
	Expect(simulation.TimeExported).To(Equal("2017-02-23T12:43:48Z"))
}

func Test_NewSimulationViewFromRequestBody_WontCreateSimulationIfThereIsNoSchemaVersion(t *testing.T) {
	RegisterTestingT(t)

	simulation, err := v2.NewSimulationViewFromRequestBody([]byte(`{
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
	Expect(simulation.GlobalActions.DelaysLogNormal).To(HaveLen(0))
}

func Test_NewSimulationViewFromRequestBody_WontBlowUpIfMetaIsMissing(t *testing.T) {
	RegisterTestingT(t)

	simulation, err := v2.NewSimulationViewFromRequestBody([]byte(`{
		"data": {}
	}`))

	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal(`Invalid JSON, missing "meta" object`))

	Expect(simulation).ToNot(BeNil())
	Expect(simulation.RequestResponsePairs).To(HaveLen(0))
	Expect(simulation.GlobalActions.Delays).To(HaveLen(0))
	Expect(simulation.GlobalActions.DelaysLogNormal).To(HaveLen(0))
}

func Test_NewSimulationViewFromRequestBody_CanCreateSimulationFromV1Payload(t *testing.T) {
	RegisterTestingT(t)

	simulation, err := v2.NewSimulationViewFromRequestBody([]byte(`{
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

	Expect(simulation.RequestResponsePairs[0].RequestMatcher.Body).To(HaveLen(0))
	Expect(simulation.RequestResponsePairs[0].RequestMatcher.Destination[0].Matcher).To(Equal("exact"))
	Expect(simulation.RequestResponsePairs[0].RequestMatcher.Destination[0].Value).To(Equal("test-server.com"))
	Expect(simulation.RequestResponsePairs[0].RequestMatcher.Headers).To(HaveLen(0))
	Expect(simulation.RequestResponsePairs[0].RequestMatcher.Method).To(HaveLen(0))
	Expect(simulation.RequestResponsePairs[0].RequestMatcher.Path).To(HaveLen(0))
	Expect(simulation.RequestResponsePairs[0].RequestMatcher.DeprecatedQuery).To(HaveLen(0))
	Expect(simulation.RequestResponsePairs[0].RequestMatcher.Scheme).To(HaveLen(0))

	Expect(simulation.RequestResponsePairs[0].Response.Body).To(Equal("exact match"))
	Expect(simulation.RequestResponsePairs[0].Response.EncodedBody).To(BeFalse())
	Expect(simulation.RequestResponsePairs[0].Response.Headers).To(HaveKeyWithValue("Header", []string{"value"}))
	Expect(simulation.RequestResponsePairs[0].Response.Status).To(Equal(200))
	Expect(simulation.RequestResponsePairs[0].Response.Templated).To(BeFalse())

	Expect(simulation.SchemaVersion).To(Equal("v3"))
	Expect(simulation.HoverflyVersion).To(Equal("v0.11.0"))
	Expect(simulation.TimeExported).To(Equal("2017-02-23T12:43:48Z"))
}

func Test_NewSimulationViewFromRequestBody_WontCreateSimulationFromInvalidV1Simulation(t *testing.T) {
	RegisterTestingT(t)

	simulation, err := v2.NewSimulationViewFromRequestBody([]byte(`{
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
	Expect(err.Error()).To(Equal("Invalid v1 simulation: [Error for <request>: request is required; Error for <response>: response is required]"))

	Expect(simulation).ToNot(BeNil())
	Expect(simulation.RequestResponsePairs).To(HaveLen(0))
	Expect(simulation.GlobalActions.Delays).To(HaveLen(0))
}


func Test_NewSimulationViewFromRequestBody_ReturnErrorMessagesOnInvalidSimulation(t *testing.T) {
	RegisterTestingT(t)

	simulation, err := v2.NewSimulationViewFromRequestBody([]byte(`{
	"data": {
		"pairs": [
			{
				"request": [],
				"response": []
			
			}
		],
		"globalActions": {
			"delays": []
		}
	},
	"meta": {
		"schemaVersion": "v4"
	}
}`))

	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(MatchRegexp(`Invalid v4 simulation: \[Error for <data.pairs.0.request|response>: Invalid type. Expected: object, given: array; Error for <data.pairs.0.response|request>: Invalid type. Expected: object, given: array\]`))

	Expect(simulation).ToNot(BeNil())
	Expect(simulation.RequestResponsePairs).To(HaveLen(0))
	Expect(simulation.GlobalActions.Delays).To(HaveLen(0))
}

func Test_NewSimulationViewFromRequestBody_WontCreateSimulationFromUnknownSchemaVersion(t *testing.T) {
	RegisterTestingT(t)

	_, err := v2.NewSimulationViewFromRequestBody([]byte(`{
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

func Test_NewSimulationViewFromRequestBody_WontCreateSimulationFromInvalidJson(t *testing.T) {
	RegisterTestingT(t)

	simulation, err := v2.NewSimulationViewFromRequestBody([]byte(`{}{}[^.^]{}{}`))

	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("Invalid JSON"))

	Expect(simulation).ToNot(BeNil())
	Expect(simulation.RequestResponsePairs).To(HaveLen(0))
	Expect(simulation.GlobalActions.Delays).To(HaveLen(0))
	Expect(simulation.GlobalActions.DelaysLogNormal).To(HaveLen(0))
}

func Test_SimulationImportResult_AddDeprecatedQueryWarning_AddsWarning(t *testing.T) {
	RegisterTestingT(t)

	unit := v2.SimulationImportResult{}
	unit.AddDeprecatedQueryWarning(15)

	Expect(unit.WarningMessages).To(HaveLen(1))

	Expect(unit.WarningMessages[0].Message).To(ContainSubstring("WARNING"))
	Expect(unit.WarningMessages[0].Message).To(ContainSubstring("deprecatedQuery"))
	Expect(unit.WarningMessages[0].Message).To(ContainSubstring("data.pairs[15].request.deprecatedQuery"))
}

func Test_SimulationImportResult_WriteResponse_IncludesMultipleWarnings(t *testing.T) {
	RegisterTestingT(t)

	unit := v2.SimulationImportResult{}
	unit.AddDeprecatedQueryWarning(15)
	unit.AddDeprecatedQueryWarning(30)
	unit.AddDeprecatedQueryWarning(45)

	Expect(unit.WarningMessages[0].Message).To(ContainSubstring("data.pairs[15].request.deprecatedQuery"))
	Expect(unit.WarningMessages[1].Message).To(ContainSubstring("data.pairs[30].request.deprecatedQuery"))
	Expect(unit.WarningMessages[2].Message).To(ContainSubstring("data.pairs[45].request.deprecatedQuery"))
}

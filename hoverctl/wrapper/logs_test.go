package wrapper

import (
	"testing"

	"time"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/matching/matchers"
	. "github.com/onsi/gomega"
)

func Test_GetLogs_GetsLogsWithCorrect_Text_Plain_AcceptHeader(t *testing.T) {
	RegisterTestingT(t)

	hoverfly.DeleteSimulation()
	hoverfly.PutSimulation(v2.SimulationViewV5{
		v2.DataViewV5{
			RequestResponsePairs: []v2.RequestMatcherResponsePairViewV5{
				{
					RequestMatcher: v2.RequestMatcherViewV5{
						Method: []v2.MatcherViewV5{
							{
								Matcher: matchers.Exact,
								Value:   "GET",
							},
						},
						Path: []v2.MatcherViewV5{
							{
								Matcher: matchers.Exact,
								Value:   "/api/v2/logs",
							},
						},
						Headers: map[string][]v2.MatcherViewV5{
							"Accept": {
								{
									Matcher: matchers.Exact,
									Value:   "text/plain",
								},
							},
						},
					},
					Response: v2.ResponseDetailsViewV5{
						Status: 200,
						Body:   "logs line 1\nlogs line 2",
					},
				},
			},
		},
		v2.MetaView{
			SchemaVersion: "v2",
		},
	})

	logs, err := GetLogs(target, "plain", nil)
	Expect(err).To(BeNil())

	Expect(logs).To(HaveLen(2))
	Expect(logs[0]).To(Equal("logs line 1"))
	Expect(logs[1]).To(Equal("logs line 2"))
}

func Test_GetLogs_CanHandleEmptyTextPlainLogResponse(t *testing.T) {
	RegisterTestingT(t)

	hoverfly.DeleteSimulation()
	hoverfly.PutSimulation(v2.SimulationViewV5{
		v2.DataViewV5{
			RequestResponsePairs: []v2.RequestMatcherResponsePairViewV5{
				{
					RequestMatcher: v2.RequestMatcherViewV5{
						Method: []v2.MatcherViewV5{
							{
								Matcher: matchers.Exact,
								Value:   "GET",
							},
						},
						Path: []v2.MatcherViewV5{
							{
								Matcher: matchers.Exact,
								Value:   "/api/v2/logs",
							},
						},
						Headers: map[string][]v2.MatcherViewV5{
							"Accept": {
								{
									Matcher: matchers.Exact,
									Value:   "text/plain",
								},
							},
						},
					},
					Response: v2.ResponseDetailsViewV5{
						Status: 200,
						Body:   ``,
					},
				},
			},
		},
		v2.MetaView{
			SchemaVersion: "v2",
		},
	})

	logs, err := GetLogs(target, "plain", nil)
	Expect(err).To(BeNil())
	Expect(logs).To(HaveLen(0))
}

func Test_GetLogs_CanHandleEmptyLineAtEndOfTextPlainLogResponse(t *testing.T) {
	RegisterTestingT(t)

	hoverfly.DeleteSimulation()
	hoverfly.PutSimulation(v2.SimulationViewV5{
		v2.DataViewV5{
			RequestResponsePairs: []v2.RequestMatcherResponsePairViewV5{
				{
					RequestMatcher: v2.RequestMatcherViewV5{
						Method: []v2.MatcherViewV5{
							{
								Matcher: matchers.Exact,
								Value:   "GET",
							},
						},
						Path: []v2.MatcherViewV5{
							{
								Matcher: matchers.Exact,
								Value:   "/api/v2/logs",
							},
						},
						Headers: map[string][]v2.MatcherViewV5{
							"Accept": {
								{
									Matcher: matchers.Exact,
									Value:   "text/plain",
								},
							},
						},
					},
					Response: v2.ResponseDetailsViewV5{
						Status: 200,
						Body:   "this is log message one\n",
					},
				},
			},
		},
		v2.MetaView{
			SchemaVersion: "v2",
		},
	})

	logs, err := GetLogs(target, "plain", nil)

	Expect(err).To(BeNil())
	Expect(logs).To(HaveLen(1))
	Expect(logs[0]).To(Equal("this is log message one"))
}

func Test_GetLogs_GetsLogsWithCorrect_JSON_AcceptHeader(t *testing.T) {
	RegisterTestingT(t)

	hoverfly.DeleteSimulation()
	hoverfly.PutSimulation(v2.SimulationViewV5{
		v2.DataViewV5{
			RequestResponsePairs: []v2.RequestMatcherResponsePairViewV5{
				{
					RequestMatcher: v2.RequestMatcherViewV5{
						Method: []v2.MatcherViewV5{
							{
								Matcher: matchers.Exact,
								Value:   "GET",
							},
						},
						Path: []v2.MatcherViewV5{
							{
								Matcher: matchers.Exact,
								Value:   "/api/v2/logs",
							},
						},
						Headers: map[string][]v2.MatcherViewV5{
							"Accept": {
								{
									Matcher: matchers.Exact,
									Value:   "application/json",
								},
							},
						},
					},
					Response: v2.ResponseDetailsViewV5{
						Status: 200,
						Body:   `{"logs":[{"msg": "logs line 1"}]}`,
					},
				},
			},
		},
		v2.MetaView{
			SchemaVersion: "v2",
		},
	})

	logs, err := GetLogs(target, "json", nil)
	Expect(err).To(BeNil())
	Expect(logs[0]).To(Equal(`{"msg":"logs line 1"}`))
}

func Test_GetLogs_ErrorsWhen_HoverflyNotAccessible(t *testing.T) {
	RegisterTestingT(t)

	_, err := GetLogs(inaccessibleTarget, "plain", nil)

	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("Could not connect to Hoverfly at something:1234"))
}

func Test_GetLogs_ErrorsWhen_HoverflyReturnsNon200(t *testing.T) {
	RegisterTestingT(t)

	hoverfly.DeleteSimulation()
	hoverfly.PutSimulation(v2.SimulationViewV5{
		v2.DataViewV5{
			RequestResponsePairs: []v2.RequestMatcherResponsePairViewV5{
				{
					RequestMatcher: v2.RequestMatcherViewV5{
						Method: []v2.MatcherViewV5{
							{
								Matcher: matchers.Exact,
								Value:   "GET",
							},
						},
						Path: []v2.MatcherViewV5{
							{
								Matcher: matchers.Exact,
								Value:   "/api/v2/logs",
							},
						},
					},
					Response: v2.ResponseDetailsViewV5{
						Status: 400,
						Body:   "{\"error\":\"test error\"}",
					},
				},
			},
		},
		v2.MetaView{
			SchemaVersion: "v2",
		},
	})

	_, err := GetLogs(target, "plain", nil)
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("Could not retrieve logs\n\ntest error"))
}

func Test_GetLogs_FiltersByDateWhenFilterTimeProvided(t *testing.T) {
	RegisterTestingT(t)

	hoverfly.DeleteSimulation()
	hoverfly.PutSimulation(v2.SimulationViewV5{
		v2.DataViewV5{
			RequestResponsePairs: []v2.RequestMatcherResponsePairViewV5{
				{
					RequestMatcher: v2.RequestMatcherViewV5{
						Method: []v2.MatcherViewV5{
							{
								Matcher: matchers.Exact,
								Value:   "GET",
							},
						},
						Path: []v2.MatcherViewV5{
							{
								Matcher: matchers.Exact,
								Value:   "/api/v2/logs",
							},
						},
						Query: &v2.QueryMatcherViewV5{
							"from": []v2.MatcherViewV5{
								{
									Matcher: matchers.Exact,
									Value:   "684552180",
								},
							},
						},
					},
					Response: v2.ResponseDetailsViewV5{
						Status: 200,
						Body:   `{"logs":[{"msg": "filtered logs"}]}`,
					},
				},
			},
		},
		v2.MetaView{
			SchemaVersion: "v2",
		},
	})
	fromTime := time.Unix(int64(684552180), 0)

	logs, err := GetLogs(target, "json", &fromTime)
	Expect(err).To(BeNil())
	Expect(logs[0]).To(Equal(`{"msg":"filtered logs"}`))
}

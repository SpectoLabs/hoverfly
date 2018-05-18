package models_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/matching/matchers"
	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/SpectoLabs/hoverfly/core/state"
	. "github.com/onsi/gomega"
)

func Test_Simulation_AddPair_CanAddAPairToTheArray(t *testing.T) {
	RegisterTestingT(t)

	unit := models.NewSimulation()

	unit.AddPair(&models.RequestMatcherResponsePair{
		models.RequestMatcher{
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "space",
				},
			},
		},
		models.ResponseDetails{},
	})

	Expect(unit.GetMatchingPairs()).To(HaveLen(1))
	Expect(unit.GetMatchingPairs()[0].RequestMatcher.Destination[0].Matcher).To(Equal("exact"))
	Expect(unit.GetMatchingPairs()[0].RequestMatcher.Destination[0].Value).To(Equal("space"))
}

func Test_Simulation_AddPair_CanAddAFullPairToTheArray(t *testing.T) {
	RegisterTestingT(t)

	unit := models.NewSimulation()

	unit.AddPair(&models.RequestMatcherResponsePair{
		models.RequestMatcher{
			Body: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "testbody",
				},
			},
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "testdestination",
				},
			},
			Headers: map[string][]string{"testheader": []string{"testvalue"}},
			Method: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "testmethod",
				},
			},
			Path: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "/testpath",
				},
			},
			Query: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "?query=test",
				},
			},
			Scheme: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "http",
				},
			},
		},
		models.ResponseDetails{
			Body:    "testresponsebody",
			Headers: map[string][]string{"testheader": []string{"testvalue"}},
			Status:  200,
		},
	})

	Expect(unit.GetMatchingPairs()).To(HaveLen(1))

	Expect(unit.GetMatchingPairs()[0].RequestMatcher.Body[0].Matcher).To(Equal("exact"))
	Expect(unit.GetMatchingPairs()[0].RequestMatcher.Body[0].Value).To(Equal("testbody"))
	Expect(unit.GetMatchingPairs()[0].RequestMatcher.Destination[0].Matcher).To(Equal("exact"))
	Expect(unit.GetMatchingPairs()[0].RequestMatcher.Destination[0].Value).To(Equal("testdestination"))
	Expect(unit.GetMatchingPairs()[0].RequestMatcher.Headers).To(HaveKeyWithValue("testheader", []string{"testvalue"}))
	Expect(unit.GetMatchingPairs()[0].RequestMatcher.Method[0].Matcher).To(Equal("exact"))
	Expect(unit.GetMatchingPairs()[0].RequestMatcher.Method[0].Value).To(Equal("testmethod"))
	Expect(unit.GetMatchingPairs()[0].RequestMatcher.Path[0].Matcher).To(Equal("exact"))
	Expect(unit.GetMatchingPairs()[0].RequestMatcher.Path[0].Value).To(Equal("/testpath"))
	Expect(unit.GetMatchingPairs()[0].RequestMatcher.Query[0].Matcher).To(Equal("exact"))
	Expect(unit.GetMatchingPairs()[0].RequestMatcher.Query[0].Value).To(Equal("?query=test"))
	Expect(unit.GetMatchingPairs()[0].RequestMatcher.Scheme[0].Matcher).To(Equal("exact"))
	Expect(unit.GetMatchingPairs()[0].RequestMatcher.Scheme[0].Value).To(Equal("http"))

	Expect(unit.GetMatchingPairs()[0].Response.Body).To(Equal("testresponsebody"))
	Expect(unit.GetMatchingPairs()[0].Response.Headers).To(HaveKeyWithValue("testheader", []string{"testvalue"}))
	Expect(unit.GetMatchingPairs()[0].Response.Status).To(Equal(200))
}

func Test_Simulation_AddPairInSequence_CanAddAFullPairToTheArray(t *testing.T) {
	RegisterTestingT(t)

	unit := models.NewSimulation()

	unit.AddPairInSequence(&models.RequestMatcherResponsePair{
		models.RequestMatcher{
			Body: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "testbody",
				},
			},
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "testdestination",
				},
			},
			Headers: map[string][]string{"testheader": []string{"testvalue"}},
			Method: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "testmethod",
				},
			},
			Path: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "/testpath",
				},
			},
			Query: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "?query=test",
				},
			},
			Scheme: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "http",
				},
			},
		},
		models.ResponseDetails{
			Body:    "testresponsebody",
			Headers: map[string][]string{"testheader": []string{"testvalue"}},
			Status:  200,
		},
	}, &state.State{State: map[string]string{}})

	Expect(unit.GetMatchingPairs()).To(HaveLen(1))

	Expect(unit.GetMatchingPairs()[0].RequestMatcher.Body[0].Matcher).To(Equal("exact"))
	Expect(unit.GetMatchingPairs()[0].RequestMatcher.Body[0].Value).To(Equal("testbody"))
	Expect(unit.GetMatchingPairs()[0].RequestMatcher.Destination[0].Matcher).To(Equal("exact"))
	Expect(unit.GetMatchingPairs()[0].RequestMatcher.Destination[0].Value).To(Equal("testdestination"))
	Expect(unit.GetMatchingPairs()[0].RequestMatcher.Headers).To(HaveKeyWithValue("testheader", []string{"testvalue"}))
	Expect(unit.GetMatchingPairs()[0].RequestMatcher.Method[0].Matcher).To(Equal("exact"))
	Expect(unit.GetMatchingPairs()[0].RequestMatcher.Method[0].Value).To(Equal("testmethod"))
	Expect(unit.GetMatchingPairs()[0].RequestMatcher.Path[0].Matcher).To(Equal("exact"))
	Expect(unit.GetMatchingPairs()[0].RequestMatcher.Path[0].Value).To(Equal("/testpath"))
	Expect(unit.GetMatchingPairs()[0].RequestMatcher.Query[0].Matcher).To(Equal("exact"))
	Expect(unit.GetMatchingPairs()[0].RequestMatcher.Query[0].Value).To(Equal("?query=test"))
	Expect(unit.GetMatchingPairs()[0].RequestMatcher.Scheme[0].Matcher).To(Equal("exact"))
	Expect(unit.GetMatchingPairs()[0].RequestMatcher.Scheme[0].Value).To(Equal("http"))

	Expect(unit.GetMatchingPairs()[0].Response.Body).To(Equal("testresponsebody"))
	Expect(unit.GetMatchingPairs()[0].Response.Headers).To(HaveKeyWithValue("testheader", []string{"testvalue"}))
	Expect(unit.GetMatchingPairs()[0].Response.Status).To(Equal(200))
}

func Test_Simulation_AddPairInSequence_CanSequence(t *testing.T) {
	RegisterTestingT(t)

	unit := models.NewSimulation()

	unit.AddPairInSequence(&models.RequestMatcherResponsePair{
		models.RequestMatcher{
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "testdestination",
				},
			},
		},
		models.ResponseDetails{
			Body:    "1",
			Headers: map[string][]string{"testheader": []string{"testvalue"}},
			Status:  200,
		},
	}, &state.State{State: map[string]string{}})

	unit.AddPairInSequence(&models.RequestMatcherResponsePair{
		models.RequestMatcher{
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "testdestination",
				},
			},
		},
		models.ResponseDetails{
			Body:    "2",
			Headers: map[string][]string{"testheader": []string{"testvalue"}},
			Status:  200,
		},
	}, &state.State{State: map[string]string{}})

	unit.AddPairInSequence(&models.RequestMatcherResponsePair{
		models.RequestMatcher{
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "testdestination",
				},
			},
		},
		models.ResponseDetails{
			Body:    "3",
			Headers: map[string][]string{"testheader": []string{"testvalue"}},
			Status:  200,
		},
	}, &state.State{State: map[string]string{}})

	Expect(unit.GetMatchingPairs()).To(HaveLen(3))

	Expect(unit.GetMatchingPairs()[0].RequestMatcher.Destination[0].Matcher).To(Equal("exact"))
	Expect(unit.GetMatchingPairs()[0].RequestMatcher.Destination[0].Value).To(Equal("testdestination"))
	Expect(unit.GetMatchingPairs()[0].RequestMatcher.RequiresState["sequence"]).To(Equal("1"))

	Expect(unit.GetMatchingPairs()[0].Response.Body).To(Equal("1"))
	Expect(unit.GetMatchingPairs()[0].Response.TransitionsState["sequence"]).To(Equal("2"))

	Expect(unit.GetMatchingPairs()[1].RequestMatcher.Destination[0].Matcher).To(Equal("exact"))
	Expect(unit.GetMatchingPairs()[1].RequestMatcher.Destination[0].Value).To(Equal("testdestination"))
	Expect(unit.GetMatchingPairs()[1].RequestMatcher.RequiresState["sequence"]).To(Equal("2"))

	Expect(unit.GetMatchingPairs()[1].Response.Body).To(Equal("2"))
	Expect(unit.GetMatchingPairs()[1].Response.TransitionsState["sequence"]).To(Equal("3"))

	Expect(unit.GetMatchingPairs()[2].RequestMatcher.Destination[0].Matcher).To(Equal("exact"))
	Expect(unit.GetMatchingPairs()[2].RequestMatcher.Destination[0].Value).To(Equal("testdestination"))
	Expect(unit.GetMatchingPairs()[2].RequestMatcher.RequiresState["sequence"]).To(Equal("3"))

	Expect(unit.GetMatchingPairs()[2].Response.Body).To(Equal("3"))
	Expect(unit.GetMatchingPairs()[2].Response.TransitionsState["sequence"]).To(Equal(""))
}

func Test_Simulation_AddPairInSequence_CanBeUsedWithAddPair(t *testing.T) {
	RegisterTestingT(t)

	unit := models.NewSimulation()

	unit.AddPair(&models.RequestMatcherResponsePair{
		models.RequestMatcher{
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "testdestination",
				},
			},
		},
		models.ResponseDetails{
			Body:    "1",
			Headers: map[string][]string{"testheader": []string{"testvalue"}},
			Status:  200,
		},
	})

	unit.AddPairInSequence(&models.RequestMatcherResponsePair{
		models.RequestMatcher{
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "testdestination",
				},
			},
		},
		models.ResponseDetails{
			Body:    "2",
			Headers: map[string][]string{"testheader": []string{"testvalue"}},
			Status:  200,
		},
	}, &state.State{State: map[string]string{}})

	Expect(unit.GetMatchingPairs()).To(HaveLen(2))

	Expect(unit.GetMatchingPairs()[0].RequestMatcher.Destination[0].Matcher).To(Equal("exact"))
	Expect(unit.GetMatchingPairs()[0].RequestMatcher.Destination[0].Value).To(Equal("testdestination"))
	Expect(unit.GetMatchingPairs()[0].RequestMatcher.RequiresState["sequence"]).To(Equal("1"))

	Expect(unit.GetMatchingPairs()[0].Response.Body).To(Equal("1"))
	Expect(unit.GetMatchingPairs()[0].Response.TransitionsState["sequence"]).To(Equal("2"))

	Expect(unit.GetMatchingPairs()[1].RequestMatcher.Destination[0].Matcher).To(Equal("exact"))
	Expect(unit.GetMatchingPairs()[1].RequestMatcher.Destination[0].Value).To(Equal("testdestination"))
	Expect(unit.GetMatchingPairs()[1].RequestMatcher.RequiresState["sequence"]).To(Equal("2"))

	Expect(unit.GetMatchingPairs()[1].Response.Body).To(Equal("2"))
}
func Test_Simulation_AddPair_WillNotSaveDuplicates(t *testing.T) {
	RegisterTestingT(t)

	unit := models.NewSimulation()

	unit.AddPair(&models.RequestMatcherResponsePair{
		models.RequestMatcher{
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "space",
				},
			},
		},
		models.ResponseDetails{},
	})

	unit.AddPair(&models.RequestMatcherResponsePair{
		models.RequestMatcher{
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "space",
				},
			},
		},
		models.ResponseDetails{},
	})

	Expect(unit.GetMatchingPairs()).To(HaveLen(1))
}

func Test_Simulation_AddPair_WillSaveTwoWhenNotDuplicates(t *testing.T) {
	RegisterTestingT(t)

	unit := models.NewSimulation()

	unit.AddPair(&models.RequestMatcherResponsePair{
		models.RequestMatcher{
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "space",
				},
			},
		},
		models.ResponseDetails{},
	})

	unit.AddPair(&models.RequestMatcherResponsePair{
		models.RequestMatcher{
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "again",
				},
			},
		},
		models.ResponseDetails{},
	})

	Expect(unit.GetMatchingPairs()).To(HaveLen(2))
	Expect(unit.GetMatchingPairs()[0].RequestMatcher.Destination[0].Value).To(Equal("space"))
	Expect(unit.GetMatchingPairs()[1].RequestMatcher.Destination[0].Value).To(Equal("again"))
}

func Test_Simulation_GetMatchingPairs(t *testing.T) {
	RegisterTestingT(t)

	unit := models.NewSimulation()

	unit.AddPair(&models.RequestMatcherResponsePair{
		models.RequestMatcher{
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "space",
				},
			},
		},
		models.ResponseDetails{},
	})

	Expect(unit.GetMatchingPairs()).To(HaveLen(1))
	Expect(unit.GetMatchingPairs()[0].RequestMatcher.Destination[0].Value).To(Equal("space"))
}

func Test_Simulation_DeleteMatchingPairs(t *testing.T) {
	RegisterTestingT(t)

	unit := models.NewSimulation()

	unit.AddPair(&models.RequestMatcherResponsePair{
		models.RequestMatcher{
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "space",
				},
			},
		},
		models.ResponseDetails{},
	})

	unit.DeleteMatchingPairs()

	Expect(unit.GetMatchingPairs()).To(HaveLen(0))
}

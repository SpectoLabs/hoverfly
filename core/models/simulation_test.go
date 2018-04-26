package models_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/models"
	. "github.com/onsi/gomega"
)

func Test_Simulation_AddRequestMatcherResponsePair_CanAddAPairToTheArray(t *testing.T) {
	RegisterTestingT(t)

	unit := models.NewSimulation()

	unit.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		models.RequestMatcher{
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: "exact",
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

func Test_Simulation_AddRequestMatcherResponsePair_CanAddAFullPairToTheArray(t *testing.T) {
	RegisterTestingT(t)

	unit := models.NewSimulation()

	unit.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		models.RequestMatcher{
			Body: []models.RequestFieldMatchers{
				{
					Matcher: "exact",
					Value:   "testbody",
				},
			},
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: "exact",
					Value:   "testdestination",
				},
			},
			Headers: map[string][]string{"testheader": []string{"testvalue"}},
			Method: []models.RequestFieldMatchers{
				{
					Matcher: "exact",
					Value:   "testmethod",
				},
			},
			Path: []models.RequestFieldMatchers{
				{
					Matcher: "exact",
					Value:   "/testpath",
				},
			},
			Query: []models.RequestFieldMatchers{
				{
					Matcher: "exact",
					Value:   "?query=test",
				},
			},
			Scheme: []models.RequestFieldMatchers{
				{
					Matcher: "exact",
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

func Test_Simulation_AddRequestMatcherResponsePair_WillNotSaveDuplicates(t *testing.T) {
	RegisterTestingT(t)

	unit := models.NewSimulation()

	unit.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		models.RequestMatcher{
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: "exact",
					Value:   "space",
				},
			},
		},
		models.ResponseDetails{},
	})

	unit.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		models.RequestMatcher{
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: "exact",
					Value:   "space",
				},
			},
		},
		models.ResponseDetails{},
	})

	Expect(unit.GetMatchingPairs()).To(HaveLen(1))
}

func Test_Simulation_AddRequestMatcherResponsePair_WillSaveTwoWhenNotDuplicates(t *testing.T) {
	RegisterTestingT(t)

	unit := models.NewSimulation()

	unit.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		models.RequestMatcher{
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: "exact",
					Value:   "space",
				},
			},
		},
		models.ResponseDetails{},
	})

	unit.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		models.RequestMatcher{
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: "exact",
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

	unit.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		models.RequestMatcher{
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: "exact",
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

	unit.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		models.RequestMatcher{
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: "exact",
					Value:   "space",
				},
			},
		},
		models.ResponseDetails{},
	})

	unit.DeleteMatchingPairs()

	Expect(unit.GetMatchingPairs()).To(HaveLen(0))
}

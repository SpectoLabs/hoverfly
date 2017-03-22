package models_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/SpectoLabs/hoverfly/core/util"
	. "github.com/onsi/gomega"
)

func Test_Simulation_AddRequestMatcherResponsePair_CanAddAPairToTheArray(t *testing.T) {
	RegisterTestingT(t)

	unit := models.NewSimulation()

	unit.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		models.RequestMatcher{
			Destination: &models.RequestFieldMatchers{
				ExactMatch: util.StringToPointer("space"),
			},
		},
		models.ResponseDetails{},
	})

	Expect(unit.MatchingPairs).To(HaveLen(1))
	Expect(*unit.MatchingPairs[0].RequestMatcher.Destination.ExactMatch).To(Equal("space"))
}

func Test_Simulation_AddRequestMatcherResponsePair_CanAddAFullPairToTheArray(t *testing.T) {
	RegisterTestingT(t)

	unit := models.NewSimulation()

	unit.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		models.RequestMatcher{
			Body: &models.RequestFieldMatchers{
				ExactMatch: util.StringToPointer("testbody"),
			},
			Destination: &models.RequestFieldMatchers{
				ExactMatch: util.StringToPointer("testdestination"),
			},
			Headers: map[string][]string{"testheader": []string{"testvalue"}},
			Method: &models.RequestFieldMatchers{
				ExactMatch: util.StringToPointer("testmethod"),
			},
			Path: &models.RequestFieldMatchers{
				ExactMatch: util.StringToPointer("/testpath"),
			},
			Query: &models.RequestFieldMatchers{
				ExactMatch: util.StringToPointer("?query=test"),
			},
			Scheme: &models.RequestFieldMatchers{
				ExactMatch: util.StringToPointer("http"),
			},
		},
		models.ResponseDetails{
			Body:    "testresponsebody",
			Headers: map[string][]string{"testheader": []string{"testvalue"}},
			Status:  200,
		},
	})

	Expect(unit.MatchingPairs).To(HaveLen(1))

	Expect(*unit.MatchingPairs[0].RequestMatcher.Body.ExactMatch).To(Equal("testbody"))
	Expect(*unit.MatchingPairs[0].RequestMatcher.Destination.ExactMatch).To(Equal("testdestination"))
	Expect(unit.MatchingPairs[0].RequestMatcher.Headers).To(HaveKeyWithValue("testheader", []string{"testvalue"}))
	Expect(*unit.MatchingPairs[0].RequestMatcher.Method.ExactMatch).To(Equal("testmethod"))
	Expect(*unit.MatchingPairs[0].RequestMatcher.Path.ExactMatch).To(Equal("/testpath"))
	Expect(*unit.MatchingPairs[0].RequestMatcher.Query.ExactMatch).To(Equal("?query=test"))
	Expect(*unit.MatchingPairs[0].RequestMatcher.Scheme.ExactMatch).To(Equal("http"))

	Expect(unit.MatchingPairs[0].Response.Body).To(Equal("testresponsebody"))
	Expect(unit.MatchingPairs[0].Response.Headers).To(HaveKeyWithValue("testheader", []string{"testvalue"}))
	Expect(unit.MatchingPairs[0].Response.Status).To(Equal(200))
}

func Test_Simulation_AddRequestMatcherResponsePair_WillNotSaveDuplicates(t *testing.T) {
	RegisterTestingT(t)

	unit := models.NewSimulation()

	unit.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		models.RequestMatcher{
			Destination: &models.RequestFieldMatchers{
				ExactMatch: util.StringToPointer("space"),
			},
		},
		models.ResponseDetails{},
	})

	unit.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		models.RequestMatcher{
			Destination: &models.RequestFieldMatchers{
				ExactMatch: util.StringToPointer("space"),
			},
		},
		models.ResponseDetails{},
	})

	Expect(unit.MatchingPairs).To(HaveLen(1))
}

func Test_Simulation_AddRequestMatcherResponsePair_WillSaveTwoWhenNotDuplicates(t *testing.T) {
	RegisterTestingT(t)

	unit := models.NewSimulation()

	unit.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		models.RequestMatcher{
			Destination: &models.RequestFieldMatchers{
				ExactMatch: util.StringToPointer("space"),
			},
		},
		models.ResponseDetails{},
	})

	unit.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		models.RequestMatcher{
			Destination: &models.RequestFieldMatchers{
				ExactMatch: util.StringToPointer("again"),
			},
		},
		models.ResponseDetails{},
	})

	Expect(unit.MatchingPairs).To(HaveLen(2))
	Expect(*unit.MatchingPairs[0].RequestMatcher.Destination.ExactMatch).To(Equal("space"))
	Expect(*unit.MatchingPairs[1].RequestMatcher.Destination.ExactMatch).To(Equal("again"))
}

package models_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/SpectoLabs/hoverfly/core/util"
	. "github.com/onsi/gomega"
)

func Test_Simulation_AddRequestTemplateResponsePair_CanAddAPairToTheArray(t *testing.T) {
	RegisterTestingT(t)

	unit := models.NewSimulation()

	unit.AddRequestTemplateResponsePair(&models.RequestTemplateResponsePair{
		models.RequestTemplate{
			Destination: &models.RequestFieldMatchers{
				ExactMatch: util.StringToPointer("space"),
			},
		},
		models.ResponseDetails{},
	})

	Expect(unit.Templates).To(HaveLen(1))
	Expect(*unit.Templates[0].RequestTemplate.Destination.ExactMatch).To(Equal("space"))
}

func Test_Simulation_AddRequestTemplateResponsePair_CanAddAFullPairToTheArray(t *testing.T) {
	RegisterTestingT(t)

	unit := models.NewSimulation()

	unit.AddRequestTemplateResponsePair(&models.RequestTemplateResponsePair{
		models.RequestTemplate{
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

	Expect(unit.Templates).To(HaveLen(1))

	Expect(*unit.Templates[0].RequestTemplate.Body.ExactMatch).To(Equal("testbody"))
	Expect(*unit.Templates[0].RequestTemplate.Destination.ExactMatch).To(Equal("testdestination"))
	Expect(unit.Templates[0].RequestTemplate.Headers).To(HaveKeyWithValue("testheader", []string{"testvalue"}))
	Expect(*unit.Templates[0].RequestTemplate.Method.ExactMatch).To(Equal("testmethod"))
	Expect(*unit.Templates[0].RequestTemplate.Path.ExactMatch).To(Equal("/testpath"))
	Expect(*unit.Templates[0].RequestTemplate.Query.ExactMatch).To(Equal("?query=test"))
	Expect(*unit.Templates[0].RequestTemplate.Scheme.ExactMatch).To(Equal("http"))

	Expect(unit.Templates[0].Response.Body).To(Equal("testresponsebody"))
	Expect(unit.Templates[0].Response.Headers).To(HaveKeyWithValue("testheader", []string{"testvalue"}))
	Expect(unit.Templates[0].Response.Status).To(Equal(200))
}

func Test_Simulation_AddRequestTemplateResponsePair_WillNotSaveDuplicates(t *testing.T) {
	RegisterTestingT(t)

	unit := models.NewSimulation()

	unit.AddRequestTemplateResponsePair(&models.RequestTemplateResponsePair{
		models.RequestTemplate{
			Destination: &models.RequestFieldMatchers{
				ExactMatch: util.StringToPointer("space"),
			},
		},
		models.ResponseDetails{},
	})

	unit.AddRequestTemplateResponsePair(&models.RequestTemplateResponsePair{
		models.RequestTemplate{
			Destination: &models.RequestFieldMatchers{
				ExactMatch: util.StringToPointer("space"),
			},
		},
		models.ResponseDetails{},
	})

	Expect(unit.Templates).To(HaveLen(1))
}

func Test_Simulation_AddRequestTemplateResponsePair_WillSaveTwoWhenNotDuplicates(t *testing.T) {
	RegisterTestingT(t)

	unit := models.NewSimulation()

	unit.AddRequestTemplateResponsePair(&models.RequestTemplateResponsePair{
		models.RequestTemplate{
			Destination: &models.RequestFieldMatchers{
				ExactMatch: util.StringToPointer("space"),
			},
		},
		models.ResponseDetails{},
	})

	unit.AddRequestTemplateResponsePair(&models.RequestTemplateResponsePair{
		models.RequestTemplate{
			Destination: &models.RequestFieldMatchers{
				ExactMatch: util.StringToPointer("again"),
			},
		},
		models.ResponseDetails{},
	})

	Expect(unit.Templates).To(HaveLen(2))
	Expect(*unit.Templates[0].RequestTemplate.Destination.ExactMatch).To(Equal("space"))
	Expect(*unit.Templates[1].RequestTemplate.Destination.ExactMatch).To(Equal("again"))
}

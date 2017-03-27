package models

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/util"
	. "github.com/onsi/gomega"
)

func Test_Simulation_AddRequestTemplateResponsePair_CanAddAPairToTheArray(t *testing.T) {
	RegisterTestingT(t)

	unit := NewSimulation()

	unit.AddRequestTemplateResponsePair(&RequestTemplateResponsePair{
		RequestTemplate{
			Destination: &RequestFieldMatchers{
				ExactMatch: util.StringToPointer("space"),
			},
		},
		ResponseDetails{},
	})

	Expect(unit.Templates).To(HaveLen(1))
	Expect(*unit.Templates[0].RequestTemplate.Destination.ExactMatch).To(Equal("space"))
}

func Test_Simulation_AddRequestTemplateResponsePair_CanAddAFullPairToTheArray(t *testing.T) {
	RegisterTestingT(t)

	unit := NewSimulation()

	unit.AddRequestTemplateResponsePair(&RequestTemplateResponsePair{
		RequestTemplate{
			Body: util.StringToPointer("testbody"),
			Destination: &RequestFieldMatchers{
				ExactMatch: util.StringToPointer("testdestination"),
			},
			Headers: map[string][]string{"testheader": []string{"testvalue"}},
			Method:  util.StringToPointer("testmethod"),
			Path: &RequestFieldMatchers{
				ExactMatch: util.StringToPointer("/testpath"),
			},
			Query: &RequestFieldMatchers{
				ExactMatch: util.StringToPointer("?query=test"),
			},
			Scheme: &RequestFieldMatchers{
				ExactMatch: util.StringToPointer("http"),
			},
		},
		ResponseDetails{
			Body:    "testresponsebody",
			Headers: map[string][]string{"testheader": []string{"testvalue"}},
			Status:  200,
		},
	})

	Expect(unit.Templates).To(HaveLen(1))

	Expect(*unit.Templates[0].RequestTemplate.Body).To(Equal("testbody"))
	Expect(*unit.Templates[0].RequestTemplate.Destination.ExactMatch).To(Equal("testdestination"))
	Expect(unit.Templates[0].RequestTemplate.Headers).To(HaveKeyWithValue("testheader", []string{"testvalue"}))
	Expect(*unit.Templates[0].RequestTemplate.Method).To(Equal("testmethod"))
	Expect(*unit.Templates[0].RequestTemplate.Path.ExactMatch).To(Equal("/testpath"))
	Expect(*unit.Templates[0].RequestTemplate.Query.ExactMatch).To(Equal("?query=test"))
	Expect(*unit.Templates[0].RequestTemplate.Scheme.ExactMatch).To(Equal("http"))

	Expect(unit.Templates[0].Response.Body).To(Equal("testresponsebody"))
	Expect(unit.Templates[0].Response.Headers).To(HaveKeyWithValue("testheader", []string{"testvalue"}))
	Expect(unit.Templates[0].Response.Status).To(Equal(200))
}

func Test_Simulation_AddRequestTemplateResponsePair_WillNotSaveDuplicates(t *testing.T) {
	RegisterTestingT(t)

	unit := NewSimulation()

	unit.AddRequestTemplateResponsePair(&RequestTemplateResponsePair{
		RequestTemplate{
			Destination: &RequestFieldMatchers{
				ExactMatch: util.StringToPointer("space"),
			},
		},
		ResponseDetails{},
	})

	unit.AddRequestTemplateResponsePair(&RequestTemplateResponsePair{
		RequestTemplate{
			Destination: &RequestFieldMatchers{
				ExactMatch: util.StringToPointer("space"),
			},
		},
		ResponseDetails{},
	})

	Expect(unit.Templates).To(HaveLen(1))
}

func Test_Simulation_AddRequestTemplateResponsePair_WillSaveTwoWhenNotDuplicates(t *testing.T) {
	RegisterTestingT(t)

	unit := NewSimulation()

	unit.AddRequestTemplateResponsePair(&RequestTemplateResponsePair{
		RequestTemplate{
			Destination: &RequestFieldMatchers{
				ExactMatch: util.StringToPointer("space"),
			},
		},
		ResponseDetails{},
	})

	unit.AddRequestTemplateResponsePair(&RequestTemplateResponsePair{
		RequestTemplate{
			Destination: &RequestFieldMatchers{
				ExactMatch: util.StringToPointer("again"),
			},
		},
		ResponseDetails{},
	})

	Expect(unit.Templates).To(HaveLen(2))
	Expect(*unit.Templates[0].RequestTemplate.Destination.ExactMatch).To(Equal("space"))
	Expect(*unit.Templates[1].RequestTemplate.Destination.ExactMatch).To(Equal("again"))
}

package models_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/SpectoLabs/hoverfly/core/util"
	. "github.com/onsi/gomega"
)

func Test_NewRequestFieldMatchersFromView_ReturnsNewStruct(t *testing.T) {
	RegisterTestingT(t)

	unit := models.NewRequestFieldMatchersFromView(&v2.RequestFieldMatchersView{
		ExactMatch: util.StringToPointer("exactly"),
	})

	Expect(unit).ToNot(BeNil())
	Expect(*unit.ExactMatch).To(Equal("exactly"))
}

func Test_NewRequestFieldMatchersFromView_WillReturnNilIfGivenNil(t *testing.T) {
	RegisterTestingT(t)

	unit := models.NewRequestFieldMatchersFromView(nil)

	Expect(unit).To(BeNil())
}

func Test_NewRequestFieldMatchers_BuildView(t *testing.T) {
	RegisterTestingT(t)

	unit := models.RequestFieldMatchers{
		ExactMatch: util.StringToPointer("exactly"),
	}

	view := unit.BuildView()
	Expect(*view.ExactMatch).To(Equal("exactly"))
}

func Test_NewRequestTemplateResponsePairFromView_BuildsPair(t *testing.T) {
	RegisterTestingT(t)

	unit := models.NewRequestTemplateResponsePairFromView(&v2.RequestResponsePairViewV2{
		Request: v2.RequestDetailsViewV2{
			Path: &v2.RequestFieldMatchersView{
				ExactMatch: util.StringToPointer("/"),
			},
		},
		Response: v2.ResponseDetailsView{
			Body: "body",
		},
	})

	Expect(*unit.RequestTemplate.Path.ExactMatch).To(Equal("/"))
	Expect(unit.RequestTemplate.Destination).To(BeNil())

	Expect(unit.Response.Body).To(Equal("body"))
}

func Test_NewRequestTemplateResponsePairFromView_SortsQuery(t *testing.T) {
	RegisterTestingT(t)

	unit := models.NewRequestTemplateResponsePairFromView(&v2.RequestResponsePairViewV2{
		Request: v2.RequestDetailsViewV2{
			Query: &v2.RequestFieldMatchersView{
				ExactMatch: util.StringToPointer("b=b&a=a"),
			},
		},
		Response: v2.ResponseDetailsView{
			Body: "body",
		},
	})

	Expect(*unit.RequestTemplate.Query.ExactMatch).To(Equal("a=a&b=b"))
}

func Test_RequestTemplate_BuildRequestDetailsFromExactMatches_GeneratesARequestDetails(t *testing.T) {
	RegisterTestingT(t)

	unit := models.RequestTemplate{
		Body: &models.RequestFieldMatchers{
			ExactMatch: util.StringToPointer("body"),
		},
		Destination: &models.RequestFieldMatchers{
			ExactMatch: util.StringToPointer("destination"),
		},
		Method: &models.RequestFieldMatchers{
			ExactMatch: util.StringToPointer("method"),
		},
		Path: &models.RequestFieldMatchers{
			ExactMatch: util.StringToPointer("path"),
		},
		Query: &models.RequestFieldMatchers{
			ExactMatch: util.StringToPointer("query"),
		},
		Scheme: &models.RequestFieldMatchers{
			ExactMatch: util.StringToPointer("scheme"),
		},
	}

	Expect(unit.BuildRequestDetailsFromExactMatches()).ToNot(BeNil())
	Expect(unit.BuildRequestDetailsFromExactMatches()).To(Equal(&models.RequestDetails{
		Body:        "body",
		Destination: "destination",
		Method:      "method",
		Path:        "path",
		Query:       "query",
		Scheme:      "scheme",
	}))
}

func Test_RequestTemplate_BuildRequestDetailsFromExactMatches_IncludesHeaders(t *testing.T) {
	RegisterTestingT(t)

	unit := models.RequestTemplate{
		Body: &models.RequestFieldMatchers{
			ExactMatch: util.StringToPointer("body"),
		},
		Destination: &models.RequestFieldMatchers{
			ExactMatch: util.StringToPointer("destination"),
		},
		Headers: map[string][]string{
			"header": []string{"value"},
		},
		Method: &models.RequestFieldMatchers{
			ExactMatch: util.StringToPointer("method"),
		},
		Path: &models.RequestFieldMatchers{
			ExactMatch: util.StringToPointer("path"),
		},
		Query: &models.RequestFieldMatchers{
			ExactMatch: util.StringToPointer("query"),
		},
		Scheme: &models.RequestFieldMatchers{
			ExactMatch: util.StringToPointer("scheme"),
		},
	}

	Expect(unit.BuildRequestDetailsFromExactMatches()).ToNot(BeNil())
	Expect(unit.BuildRequestDetailsFromExactMatches()).To(Equal(&models.RequestDetails{
		Body:        "body",
		Destination: "destination",
		Headers: map[string][]string{
			"header": []string{"value"},
		},
		Method: "method",
		Path:   "path",
		Query:  "query",
		Scheme: "scheme",
	}))
}

func Test_RequestTemplate_BuildRequestDetailsFromExactMatches_ReturnsNilIfEmptyTemplate(t *testing.T) {
	RegisterTestingT(t)

	unit := models.RequestTemplate{}

	Expect(unit.BuildRequestDetailsFromExactMatches()).To(BeNil())
}

func Test_RequestTemplate_BuildRequestDetailsFromExactMatches_ReturnsNilIfMissingAnExactMatch(t *testing.T) {
	RegisterTestingT(t)

	unit := models.RequestTemplate{
		Destination: &models.RequestFieldMatchers{
			ExactMatch: util.StringToPointer("destination"),
		},
		Method: &models.RequestFieldMatchers{
			ExactMatch: util.StringToPointer("method"),
		},
		Path: &models.RequestFieldMatchers{
			ExactMatch: util.StringToPointer("path"),
		},
		Query: &models.RequestFieldMatchers{
			ExactMatch: util.StringToPointer("query"),
		},
		Scheme: &models.RequestFieldMatchers{
			ExactMatch: util.StringToPointer("scheme"),
		},
	}

	Expect(unit.BuildRequestDetailsFromExactMatches()).To(BeNil())
}

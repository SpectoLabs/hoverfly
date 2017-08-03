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

func Test_NewRequestMatcherResponsePairFromView_BuildsPair(t *testing.T) {
	RegisterTestingT(t)

	unit := models.NewRequestMatcherResponsePairFromView(&v2.RequestMatcherResponsePairViewV4{
		RequestMatcher: v2.RequestMatcherViewV4{
			Path: &v2.RequestFieldMatchersView{
				ExactMatch: util.StringToPointer("/"),
			},
		},
		Response: v2.ResponseDetailsViewV4{
			Body: "body",
		},
	})

	Expect(*unit.RequestMatcher.Path.ExactMatch).To(Equal("/"))
	Expect(unit.RequestMatcher.Destination).To(BeNil())

	Expect(unit.Response.Body).To(Equal("body"))
}

func Test_NewRequestMatcherResponsePairFromView_SortsQuery(t *testing.T) {
	RegisterTestingT(t)

	unit := models.NewRequestMatcherResponsePairFromView(&v2.RequestMatcherResponsePairViewV4{
		RequestMatcher: v2.RequestMatcherViewV4{
			Query: &v2.RequestFieldMatchersView{
				ExactMatch: util.StringToPointer("b=b&a=a"),
			},
		},
		Response: v2.ResponseDetailsViewV4{
			Body: "body",
		},
	})

	Expect(*unit.RequestMatcher.Query.ExactMatch).To(Equal("a=a&b=b"))
}

func Test_NewRequestMatcherResponsePairFromView_StoresTemplated(t *testing.T) {
	RegisterTestingT(t)

	unit := models.NewRequestMatcherResponsePairFromView(&v2.RequestMatcherResponsePairViewV4{
		RequestMatcher: v2.RequestMatcherViewV4{
			Query: &v2.RequestFieldMatchersView{
				ExactMatch: util.StringToPointer("b=b&a=a"),
			},
		},
		Response: v2.ResponseDetailsViewV4{
			Body:      "body",
			Templated: true,
		},
	})

	Expect(unit.Response.Templated).To(BeTrue())
}

func Test_RequestMatcher_BuildRequestDetailsFromExactMatches_GeneratesARequestDetails(t *testing.T) {
	RegisterTestingT(t)

	unit := models.RequestMatcher{
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
			ExactMatch: util.StringToPointer("query=two"),
		},
		Scheme: &models.RequestFieldMatchers{
			ExactMatch: util.StringToPointer("scheme"),
		},
	}

	Expect(unit.ToEagerlyCachable()).ToNot(BeNil())
	Expect(unit.ToEagerlyCachable()).To(Equal(&models.RequestDetails{
		Body:        "body",
		Destination: "destination",
		Method:      "method",
		Path:        "path",
		Query:       map[string][]string{"query": []string{"two"}},
		Scheme:      "scheme",
	}))
}

func Test_RequestMatcher_BuildRequestDetailsFromExactMatches_ReturnsNilIfEmpty(t *testing.T) {
	RegisterTestingT(t)

	unit := models.RequestMatcher{}

	Expect(unit.ToEagerlyCachable()).To(BeNil())
}

func Test_RequestMatcher_BuildRequestDetailsFromExactMatches_ReturnsNilIfMissingAnExactMatch(t *testing.T) {
	RegisterTestingT(t)

	unit := models.RequestMatcher{
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

	Expect(unit.ToEagerlyCachable()).To(BeNil())
}

package matching_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/matching"
	"github.com/SpectoLabs/hoverfly/core/models"
	. "github.com/onsi/gomega"
)

func Test_CacheMatcher_GetResponse_WillReturnErrorIfCacheIsNil(t *testing.T) {
	RegisterTestingT(t)
	unit := matching.CacheMatcher{}

	_, err := unit.GetResponse(&models.RequestDetails{})
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("No cache set"))
}

func Test_CacheMatcher_GetAllResponses_WillReturnErrorIfCacheIsNil(t *testing.T) {
	RegisterTestingT(t)
	unit := matching.CacheMatcher{}

	_, err := unit.GetAllResponses()
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("No cache set"))
}

func Test_CacheMatcher_SaveRequestResponsePair_WillReturnErrorIfCacheIsNil(t *testing.T) {
	RegisterTestingT(t)
	unit := matching.CacheMatcher{}

	err := unit.SaveRequestTemplateResponsePair(models.RequestDetails{}, nil)
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("No cache set"))
}

func Test_CacheMatcher_FlushCache_WillReturnErrorIfCacheIsNil(t *testing.T) {
	RegisterTestingT(t)
	unit := matching.CacheMatcher{}

	err := unit.FlushCache()
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("No cache set"))
}

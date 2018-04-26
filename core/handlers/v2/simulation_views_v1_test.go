package v2_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/util"
	. "github.com/onsi/gomega"
)

func Test_RequestDetailsView_GetQuery_SortsQueryString(t *testing.T) {
	RegisterTestingT(t)

	unit := v2.RequestDetailsView{
		Query: util.StringToPointer("b=b&a=a"),
	}
	queryString := unit.GetQuery()
	Expect(queryString).ToNot(BeNil())

	Expect(*queryString).To(Equal("a=a&b=b"))
}

func Test_RequestDetailsView_GetQuery_ReturnsNilIfNil(t *testing.T) {
	RegisterTestingT(t)

	unit := v2.RequestDetailsView{
		Query: nil,
	}
	queryString := unit.GetQuery()
	Expect(queryString).To(BeNil())
}

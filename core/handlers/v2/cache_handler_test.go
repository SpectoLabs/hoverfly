package v2

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/SpectoLabs/hoverfly/core/util"
	. "github.com/onsi/gomega"
)

type HoverflyCacheStub struct {
	GetError    bool
	FlushCalled bool
	FlushError  bool
}

func (this HoverflyCacheStub) GetCache() ([]RequestResponsePairView, error) {
	if this.GetError {
		return nil, errors.New("There was an error")
	}

	return []RequestResponsePairView{
		RequestResponsePairView{
			Request: RequestDetailsView{
				Destination: util.StringToPointer("one"),
			},
			Response: ResponseDetailsView{},
		},
		RequestResponsePairView{
			Request: RequestDetailsView{
				Destination: util.StringToPointer("two"),
			},
			Response: ResponseDetailsView{},
		},
	}, nil
}

func (this *HoverflyCacheStub) FlushCache() error {
	this.FlushCalled = true

	if this.FlushError {
		return errors.New("There was an error")
	}

	return nil
}

func Test_Get_ReturnsTheCache(t *testing.T) {
	RegisterTestingT(t)

	stubHoverfly := &HoverflyCacheStub{}
	unit := CacheHandler{Hoverfly: stubHoverfly}

	request, err := http.NewRequest("GET", "/api/v2/cache", nil)
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Get, request)

	Expect(response.Code).To(Equal(http.StatusOK))

	cacheView, err := unmarshalCacheView(response.Body)
	Expect(err).To(BeNil())

	Expect(cacheView.RequestResponsePairs).To(HaveLen(2))
	Expect(*cacheView.RequestResponsePairs[0].Request.Destination).To(Equal("one"))
	Expect(*cacheView.RequestResponsePairs[1].Request.Destination).To(Equal("two"))
}

func Test_Get_ReturnsNiceErrorMessage(t *testing.T) {
	RegisterTestingT(t)

	stubHoverfly := &HoverflyCacheStub{GetError: true}
	unit := CacheHandler{Hoverfly: stubHoverfly}

	request, err := http.NewRequest("GET", "/api/v2/cache", nil)
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Get, request)

	Expect(response.Code).To(Equal(http.StatusInternalServerError))

	errorView, err := unmarshalErrorView(response.Body)
	Expect(err).To(BeNil())

	Expect(errorView.Error).To(Equal("There was an error"))
}

func Test_Delete_CallsFlushCache(t *testing.T) {
	RegisterTestingT(t)

	stubHoverfly := &HoverflyCacheStub{}
	unit := CacheHandler{Hoverfly: stubHoverfly}

	request, err := http.NewRequest("DELETE", "/api/v2/cache", nil)
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Delete, request)

	Expect(response.Code).To(Equal(http.StatusOK))

	Expect(stubHoverfly.FlushCalled).To(BeTrue())
}

func Test_Delete_ReturnsNiceErrorMessage(t *testing.T) {
	RegisterTestingT(t)

	stubHoverfly := &HoverflyCacheStub{FlushError: true}
	unit := CacheHandler{Hoverfly: stubHoverfly}

	request, err := http.NewRequest("DELETE", "/api/v2/cache", nil)
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Delete, request)

	Expect(response.Code).To(Equal(http.StatusInternalServerError))

	errorView, err := unmarshalErrorView(response.Body)
	Expect(err).To(BeNil())

	Expect(errorView.Error).To(Equal("There was an error"))
}

func unmarshalCacheView(buffer *bytes.Buffer) (CacheView, error) {
	body, err := ioutil.ReadAll(buffer)
	if err != nil {
		return CacheView{}, err
	}

	var cacheView CacheView

	err = json.Unmarshal(body, &cacheView)
	if err != nil {
		return CacheView{}, err
	}

	return cacheView, nil
}

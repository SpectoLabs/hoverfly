package v2

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	. "github.com/onsi/gomega"
)

type HoverflyUpstreamProxyStub struct{}

func (this HoverflyUpstreamProxyStub) GetUpstreamProxy() string {
	return "upstream-proxy.org"
}

func Test_HoverflyUpstreamProxyHandler_GetReturnsUpstreamProxy(t *testing.T) {
	RegisterTestingT(t)

	stubHoverfly := &HoverflyUpstreamProxyStub{}
	unit := HoverflyUpstreamProxyHandler{Hoverfly: stubHoverfly}

	request, err := http.NewRequest("GET", "", nil)
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Get, request)

	Expect(response.Code).To(Equal(http.StatusOK))

	upstreamProxyView, err := unmarshalUpsteamProxyView(response.Body)
	Expect(err).To(BeNil())

	Expect(upstreamProxyView.UpstreamProxy).To(Equal("upstream-proxy.org"))
}

func unmarshalUpsteamProxyView(buffer *bytes.Buffer) (UpstreamProxyView, error) {
	body, err := ioutil.ReadAll(buffer)
	if err != nil {
		return UpstreamProxyView{}, err
	}

	var upstreamProxyView UpstreamProxyView

	err = json.Unmarshal(body, &upstreamProxyView)
	if err != nil {
		return UpstreamProxyView{}, err
	}

	return upstreamProxyView, nil
}

package v2

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	. "github.com/onsi/gomega"
)

type HoverflyCORSStub struct {
}

func (h HoverflyCORSStub) GetCORS() CORSView {
	return CORSView{
		Enabled:     true,
		AllowOrigin: "*",
	}
}

func TestH_HoverflyCORSHandler_GetReturnsTheCorrectDestination(t *testing.T) {
	RegisterTestingT(t)

	stubHoverfly := &HoverflyCORSStub{}
	unit := HoverflyCORSHandler{Hoverfly: stubHoverfly}

	request, err := http.NewRequest("GET", "", nil)
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Get, request)

	Expect(response.Code).To(Equal(http.StatusOK))

	corsView, err := unmarshalCORSView(response.Body)
	Expect(err).To(BeNil())
	Expect(corsView.Enabled).To(BeTrue())
	Expect(corsView.AllowOrigin).To(Equal("*"))
}

func Test_HoverflyCORSHandler_Options_GetsOptions(t *testing.T) {
	RegisterTestingT(t)

	var stubHoverfly HoverflyCORSStub
	unit := HoverflyCORSHandler{Hoverfly: &stubHoverfly}

	request, err := http.NewRequest("OPTIONS", "/api/v2/hoverfly/cors", nil)
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Options, request)

	Expect(response.Code).To(Equal(http.StatusOK))
	Expect(response.Header().Get("Allow")).To(Equal("OPTIONS, GET, PUT"))
}

func unmarshalCORSView(buffer *bytes.Buffer) (CORSView, error) {
	body, err := ioutil.ReadAll(buffer)
	if err != nil {
		return CORSView{}, err
	}

	var corsView CORSView

	err = json.Unmarshal(body, &corsView)
	if err != nil {
		return CORSView{}, err
	}

	return corsView, nil
}

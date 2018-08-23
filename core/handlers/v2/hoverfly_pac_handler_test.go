package v2

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	. "github.com/onsi/gomega"
)

type HoverflyPacStub struct {
	PACFile []byte
}

func (this HoverflyPacStub) GetPACFile() []byte {
	return this.PACFile
}

func (this *HoverflyPacStub) SetPACFile(PACFile []byte) {
	this.PACFile = PACFile
}

func (this *HoverflyPacStub) DeletePACFile() {
	this.PACFile = nil
}

func Test_HoverflyPACHandler_Get_ReturnsPACfile(t *testing.T) {
	RegisterTestingT(t)

	stubHoverfly := &HoverflyPacStub{
		PACFile: []byte("PACFILE"),
	}

	unit := HoverflyPACHandler{Hoverfly: stubHoverfly}

	request, err := http.NewRequest("GET", "", nil)
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Get, request)

	Expect(response.Code).To(Equal(http.StatusOK))
	Expect(response.Header().Get("Content-Type")).To(Equal("application/x-ns-proxy-autoconfig"))

	bodyBytes, err := ioutil.ReadAll(response.Body)
	Expect(err).To(BeNil())

	Expect(string(bodyBytes)).To(Equal("PACFILE"))
}

func Test_HoverflyPACHandler_Get_Returns404IfNotSet(t *testing.T) {
	RegisterTestingT(t)

	stubHoverfly := &HoverflyPacStub{}

	unit := HoverflyPACHandler{Hoverfly: stubHoverfly}

	request, err := http.NewRequest("GET", "", nil)
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Get, request)

	Expect(response.Code).To(Equal(http.StatusNotFound))

	errorViewResponse, err := unmarshalErrorView(response.Body)
	Expect(err).To(BeNil())

	Expect(errorViewResponse.Error).To(Equal("Not found"))
}

func Test_HoverflyPACHandler_Put_ReturnsPACfile(t *testing.T) {
	RegisterTestingT(t)

	stubHoverfly := &HoverflyPacStub{}

	unit := HoverflyPACHandler{Hoverfly: stubHoverfly}

	request, err := http.NewRequest("PUT", "", ioutil.NopCloser(bytes.NewBuffer([]byte("PACFILE"))))
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Put, request)

	Expect(response.Code).To(Equal(http.StatusOK))
	Expect(response.Header().Get("Content-Type")).To(Equal("application/x-ns-proxy-autoconfig"))

	bodyBytes, err := ioutil.ReadAll(response.Body)
	Expect(err).To(BeNil())

	Expect(string(bodyBytes)).To(Equal("PACFILE"))
	Expect(string(stubHoverfly.PACFile)).To(Equal("PACFILE"))
}

func Test_HoverflyPACHandler_Delete_DeletesPACfile(t *testing.T) {
	RegisterTestingT(t)

	stubHoverfly := &HoverflyPacStub{
		PACFile: []byte("PACFILE"),
	}

	unit := HoverflyPACHandler{Hoverfly: stubHoverfly}

	request, err := http.NewRequest("DELETE", "", nil)
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Delete, request)

	Expect(response.Code).To(Equal(http.StatusOK))

	bodyBytes, err := ioutil.ReadAll(response.Body)
	Expect(err).To(BeNil())

	Expect(string(bodyBytes)).To(Equal(""))
	Expect(stubHoverfly.PACFile).To(BeNil())
}

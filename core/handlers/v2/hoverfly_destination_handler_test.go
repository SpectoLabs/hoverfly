package v2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	. "github.com/onsi/gomega"
)

type HoverflyDestinationStub struct {
	Destination string
}

func (this HoverflyDestinationStub) GetDestination() string {
	return this.Destination
}

func (this *HoverflyDestinationStub) SetDestination(destination string) error {
	this.Destination = destination
	if destination == "error" {
		return fmt.Errorf("error")
	}

	return nil
}

func TestHoverflyDestinationHandlerGetReturnsTheCorrectDestination(t *testing.T) {
	RegisterTestingT(t)

	stubHoverfly := &HoverflyDestinationStub{Destination: "testination.com"}
	unit := HoverflyDestinationHandler{Hoverfly: stubHoverfly}

	request, err := http.NewRequest("GET", "", nil)
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Get, request)

	Expect(response.Code).To(Equal(http.StatusOK))

	destinationView, err := unmarshalDestinationView(response.Body)
	Expect(err).To(BeNil())
	Expect(destinationView.Destination).To(Equal("testination.com"))
}

func TestHoverflyDestinationHandlerPutSetsTheNewDestinationAndReplacesTheTestDestination(t *testing.T) {
	RegisterTestingT(t)

	stubHoverfly := &HoverflyDestinationStub{Destination: "testination.com"}
	unit := HoverflyDestinationHandler{Hoverfly: stubHoverfly}

	destinationView := &DestinationView{Destination: "new-domain.com"}

	bodyBytes, err := json.Marshal(destinationView)
	Expect(err).To(BeNil())

	request, err := http.NewRequest("PUT", "", ioutil.NopCloser(bytes.NewBuffer(bodyBytes)))
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Put, request)
	Expect(response.Code).To(Equal(http.StatusOK))
	Expect(stubHoverfly.Destination).To(Equal("new-domain.com"))

	destinationViewResponse, err := unmarshalDestinationView(response.Body)
	Expect(err).To(BeNil())

	Expect(destinationViewResponse.Destination).To(Equal("new-domain.com"))
}

func TestHoverflyDestinationHandlerPutWill422ErrorIfHoverflyErrors(t *testing.T) {
	RegisterTestingT(t)

	var stubHoverfly HoverflyDestinationStub
	unit := HoverflyDestinationHandler{Hoverfly: &stubHoverfly}

	destinationView := &DestinationView{Destination: "error"}

	bodyBytes, err := json.Marshal(destinationView)
	Expect(err).To(BeNil())

	request, err := http.NewRequest("PUT", "", ioutil.NopCloser(bytes.NewBuffer(bodyBytes)))
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Put, request)
	Expect(response.Code).To(Equal(http.StatusUnprocessableEntity))

	errorViewResponse, err := unmarshalErrorView(response.Body)
	Expect(err).To(BeNil())

	Expect(errorViewResponse.Error).To(Equal("error"))
}

func TestHoverflyDestinationeHandlerPutWill400ErrorIfJsonIsBad(t *testing.T) {
	RegisterTestingT(t)

	var stubHoverfly HoverflyDestinationStub
	unit := HoverflyDestinationHandler{Hoverfly: &stubHoverfly}

	bodyBytes := []byte("{{}{}}")

	request, err := http.NewRequest("PUT", "/api/v2/hoverfly/mode", ioutil.NopCloser(bytes.NewBuffer(bodyBytes)))
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Put, request)
	Expect(response.Code).To(Equal(http.StatusBadRequest))

	errorViewResponse, err := unmarshalErrorView(response.Body)
	Expect(err).To(BeNil())

	Expect(errorViewResponse.Error).To(Equal("Malformed JSON"))
}

func Test_HoverflyDestinationHandler_Options_GetsOptions(t *testing.T) {
	RegisterTestingT(t)

	var stubHoverfly HoverflyDestinationStub
	unit := HoverflyDestinationHandler{Hoverfly: &stubHoverfly}

	request, err := http.NewRequest("OPTIONS", "/api/v2/hoverfly/mode", nil)
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Options, request)

	Expect(response.Code).To(Equal(http.StatusOK))
	Expect(response.Header().Get("Allow")).To(Equal("OPTIONS, GET, PUT"))
}

func unmarshalDestinationView(buffer *bytes.Buffer) (DestinationView, error) {
	body, err := ioutil.ReadAll(buffer)
	if err != nil {
		return DestinationView{}, err
	}

	var destinationView DestinationView

	err = json.Unmarshal(body, &destinationView)
	if err != nil {
		return DestinationView{}, err
	}

	return destinationView, nil
}

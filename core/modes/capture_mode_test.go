package modes_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/SpectoLabs/hoverfly/core/modes"
	. "github.com/onsi/gomega"
)

type hoverflyCaptureStub struct {
	SavedRequest  *models.RequestDetails
	SavedResponse *models.ResponseDetails
	SavedHeaders  []string
	MiddlewareSet bool
}

// ApplyMiddleware - Stub implementation of modes.HoverflyCapture interface
func (this hoverflyCaptureStub) ApplyMiddleware(pair models.RequestResponsePair) (models.RequestResponsePair, error) {
	return pair, nil
}

// DoRequest - Stub implementation of modes.HoverflyCapture interface
func (this hoverflyCaptureStub) DoRequest(request *http.Request) (*http.Response, error) {
	response := &http.Response{}
	if request.Host == "error.com" {
		return nil, errors.New("Could not reach error.com")
	}

	response.StatusCode = 200
	response.Body = ioutil.NopCloser(bytes.NewBufferString("test"))

	if request.Host == "trailer.com" {
		response.Header = make(http.Header)
		response.Header.Set("Content-Type", "application/json")
		response.Trailer = make(http.Header)
		response.Trailer.Set("X-Streaming-Error", "Connection closed")
		response.Trailer.Set("X-Bin-Id", "xyz")
	}

	return response, nil
}

// Save - Stub implementation of modes.HoverflyCapture interface
func (this *hoverflyCaptureStub) Save(request *models.RequestDetails, response *models.ResponseDetails, headersToSave []string, recordSequence bool) error {
	this.SavedRequest = request
	this.SavedResponse = response
	this.SavedHeaders = headersToSave

	return nil
}

func Test_CaptureMode_CanSetArguments(t *testing.T) {
	RegisterTestingT(t)

	unit := &modes.CaptureMode{
		Hoverfly: &hoverflyCaptureStub{},
	}

	unit.SetArguments(modes.ModeArguments{
		Headers: []string{"value", "two"},
	})

	Expect(unit.Arguments.Headers).To(ContainElement("value"))
	Expect(unit.Arguments.Headers).To(ContainElement("two"))
}

func Test_CaptureMode_WhenGivenARequestItWillMakeTheRequestAndSaveIt(t *testing.T) {
	RegisterTestingT(t)

	hoverflyStub := &hoverflyCaptureStub{}

	unit := &modes.CaptureMode{
		Hoverfly: hoverflyStub,
	}

	requestDetails := models.RequestDetails{
		Scheme:      "http",
		Destination: "positive-match.com",
	}

	request, err := http.NewRequest("GET", "http://positive-match.com", nil)
	Expect(err).To(BeNil())

	response, err := unit.Process(request, requestDetails)
	Expect(err).To(BeNil())

	Expect(response.StatusCode).To(Equal(200))

	responseBody, err := ioutil.ReadAll(response.Body)
	Expect(err).To(BeNil())

	Expect(string(responseBody)).To(Equal("test"))

	Expect(hoverflyStub.SavedRequest.Destination).To(Equal("positive-match.com"))
	Expect(hoverflyStub.SavedResponse.Body).To(Equal("test"))
}

func Test_CaptureMode_IfHeadersArgumentNotSet_CallsSaveWithEmptyList(t *testing.T) {
	RegisterTestingT(t)

	hoverflyStub := &hoverflyCaptureStub{}

	unit := &modes.CaptureMode{
		Hoverfly: hoverflyStub,
	}

	requestDetails := models.RequestDetails{
		Scheme:      "http",
		Destination: "positive-match.com",
	}

	request, _ := http.NewRequest("GET", "http://positive-match.com", nil)

	_, err := unit.Process(request, requestDetails)
	Expect(err).To(BeNil())

	Expect(hoverflyStub.SavedHeaders).To(HaveLen(0))
}

func Test_CaptureMode_IfHeadersArgumentSetToAll_CallsSaveWithEmptyList(t *testing.T) {
	RegisterTestingT(t)

	hoverflyStub := &hoverflyCaptureStub{}

	unit := &modes.CaptureMode{
		Hoverfly: hoverflyStub,
	}

	requestDetails := models.RequestDetails{
		Scheme:      "http",
		Destination: "positive-match.com",
	}

	unit.SetArguments(modes.ModeArguments{
		Headers: []string{"*"},
	})

	request, _ := http.NewRequest("GET", "http://positive-match.com", nil)

	_, err := unit.Process(request, requestDetails)
	Expect(err).To(BeNil())

	Expect(hoverflyStub.SavedHeaders).To(HaveLen(1))
	Expect(hoverflyStub.SavedHeaders).To(ContainElement("*"))
}

func Test_CaptureMode_IfHeadersArgumentSetToOneHeaders_CallsSaveWithOneHeaderList(t *testing.T) {
	RegisterTestingT(t)

	hoverflyStub := &hoverflyCaptureStub{}

	unit := &modes.CaptureMode{
		Hoverfly: hoverflyStub,
	}

	requestDetails := models.RequestDetails{
		Scheme:      "http",
		Destination: "positive-match.com",
	}

	unit.SetArguments(modes.ModeArguments{
		Headers: []string{"Content-Type"},
	})

	request, _ := http.NewRequest("GET", "http://positive-match.com", nil)

	_, err := unit.Process(request, requestDetails)
	Expect(err).To(BeNil())

	Expect(hoverflyStub.SavedHeaders).To(HaveLen(1))
	Expect(hoverflyStub.SavedHeaders).To(ContainElement("Content-Type"))
}

func Test_CaptureMode_WhenGivenABadRequestItWillError(t *testing.T) {
	RegisterTestingT(t)

	hoverflyStub := &hoverflyCaptureStub{}

	unit := &modes.CaptureMode{
		Hoverfly: hoverflyStub,
	}

	requestDetails := models.RequestDetails{
		Scheme:      "http",
		Destination: "error.com",
	}

	request, err := http.NewRequest("GET", "http://error.com", nil)
	Expect(err).To(BeNil())

	response, err := unit.Process(request, requestDetails)
	Expect(err).ToNot(BeNil())

	Expect(response.StatusCode).To(Equal(http.StatusBadGateway))

	responseBody, err := ioutil.ReadAll(response.Body)
	Expect(err).To(BeNil())

	Expect(string(responseBody)).To(ContainSubstring("There was an error when forwarding the request to the intended destination"))
	Expect(string(responseBody)).To(ContainSubstring("Could not reach error.com"))

	Expect(hoverflyStub.SavedRequest).To(BeNil())
	Expect(hoverflyStub.SavedResponse).To(BeNil())
}

func Test_CaptureMode_SavesResponseTrailersIfPresent(t *testing.T) {
	RegisterTestingT(t)

	hoverflyStub := &hoverflyCaptureStub{}

	unit := &modes.CaptureMode{
		Hoverfly: hoverflyStub,
	}

	requestDetails := models.RequestDetails{
		Scheme:      "http",
		Destination: "trailer.com",
	}

	request, _ := http.NewRequest("GET", "http://trailer.com", nil)

	response, err := unit.Process(request, requestDetails)
	Expect(err).To(BeNil())

	Expect(response.Header).To(HaveLen(1))
	Expect(response.Trailer).To(HaveLen(2))

	Expect(hoverflyStub.SavedResponse.Headers).To(HaveLen(4))
	Expect(hoverflyStub.SavedResponse.Headers["Content-Type"]).To(ConsistOf("application/json"))
	Expect(hoverflyStub.SavedResponse.Headers["Trailer"]).To(ConsistOf("X-Streaming-Error", "X-Bin-Id"))
	Expect(hoverflyStub.SavedResponse.Headers["X-Streaming-Error"]).To(ConsistOf("Connection closed"))
	Expect(hoverflyStub.SavedResponse.Headers["X-Bin-Id"]).To(ConsistOf("xyz"))
}

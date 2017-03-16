package modes_test

import (
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/SpectoLabs/hoverfly/core/modes"
	. "github.com/onsi/gomega"
)

type hoverflySynthesizeStub struct {
	MiddlewareSet bool
}

func (this hoverflySynthesizeStub) ApplyMiddleware(pair models.RequestResponsePair) (models.RequestResponsePair, error) {
	if pair.Request.Destination == "error.com" {
		return pair, errors.New("Middleware failed")
	}
	pair.Response.Body = "modified by middleware"
	pair.Response.Status = 201
	return pair, nil
}

func (this hoverflySynthesizeStub) IsMiddlewareSet() bool {
	return this.MiddlewareSet
}

func Test_SynthesizeMode_WhenGivenARequestItWillUseMiddlewareToGenerateAResponse(t *testing.T) {
	RegisterTestingT(t)

	hoverflyStub := &hoverflySynthesizeStub{
		MiddlewareSet: true,
	}

	unit := &modes.SynthesizeMode{
		Hoverfly: hoverflyStub,
	}

	requestDetails := models.RequestDetails{
		Destination: "positive-match.com",
	}

	request, err := http.NewRequest("GET", "http://positive-match.com", nil)
	Expect(err).To(BeNil())

	response, err := unit.Process(request, requestDetails)
	Expect(err).To(BeNil())

	Expect(response.StatusCode).To(Equal(http.StatusCreated))

	responseBody, err := ioutil.ReadAll(response.Body)
	Expect(err).To(BeNil())

	Expect(string(responseBody)).To(Equal("modified by middleware"))
}

func Test_SynthesizeMode_IfMiddlewareFailsThenModeReturnsNiceError(t *testing.T) {
	RegisterTestingT(t)

	hoverflyStub := &hoverflySynthesizeStub{
		MiddlewareSet: true,
	}

	unit := &modes.SynthesizeMode{
		Hoverfly: hoverflyStub,
	}

	requestDetails := models.RequestDetails{
		Destination: "error.com",
	}

	request, err := http.NewRequest("GET", "http://error.com", nil)
	Expect(err).To(BeNil())

	response, err := unit.Process(request, requestDetails)
	Expect(err).ToNot(BeNil())

	Expect(response.StatusCode).To(Equal(http.StatusBadGateway))

	responseBody, err := ioutil.ReadAll(response.Body)
	Expect(err).To(BeNil())

	Expect(string(responseBody)).To(ContainSubstring("There was an error when executing middleware"))
	Expect(string(responseBody)).To(ContainSubstring("Middleware failed"))
}

func Test_SynthesizeMode_IfMiddlewareNotSetModeReturnsNiceError(t *testing.T) {
	RegisterTestingT(t)

	hoverflyStub := &hoverflySynthesizeStub{
		MiddlewareSet: false,
	}

	unit := &modes.SynthesizeMode{
		Hoverfly: hoverflyStub,
	}

	requestDetails := models.RequestDetails{}

	request, err := http.NewRequest("GET", "http://test.com", nil)
	Expect(err).To(BeNil())

	response, err := unit.Process(request, requestDetails)
	Expect(err).ToNot(BeNil())

	Expect(response.StatusCode).To(Equal(http.StatusBadGateway))

	responseBody, err := ioutil.ReadAll(response.Body)
	Expect(err).To(BeNil())

	Expect(string(responseBody)).To(ContainSubstring("There was an error when creating a synthetic response"))
	Expect(string(responseBody)).To(ContainSubstring("Middleware not set"))
}

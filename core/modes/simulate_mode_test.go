package modes_test

import (
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/SpectoLabs/hoverfly/core/matching"
	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/SpectoLabs/hoverfly/core/modes"
	. "github.com/onsi/gomega"
)

type hoverflySimulateStub struct{}

func (this hoverflySimulateStub) GetResponse(request models.RequestDetails) (*models.ResponseDetails, *matching.MatchingError) {
	if request.Destination == "positive-match.com" {
		return &models.ResponseDetails{
			Status: 200,
		}, nil
	} else {
		return nil, &matching.MatchingError{
			Description: "matching-error",
			StatusCode:  500,
		}
	}
}

func (this hoverflySimulateStub) ApplyMiddleware(pair models.RequestResponsePair) (models.RequestResponsePair, error) {
	if pair.Request.Path == "middleware-error" {
		return pair, errors.New("middleware-error")
	}
	return pair, nil
}

func Test_SimulateMode_WhenGivenAMatchingRequestItReturnsTheCorrectResponse(t *testing.T) {
	RegisterTestingT(t)

	unit := &modes.SimulateMode{
		Hoverfly: hoverflySimulateStub{},
	}

	request := models.RequestDetails{
		Destination: "positive-match.com",
	}

	response, err := unit.Process(nil, request)
	Expect(err).To(BeNil())

	Expect(response.StatusCode).To(Equal(200))
}

func Test_SimulateMode_WhenGivenANonMatchingRequestItReturnsAnError(t *testing.T) {
	RegisterTestingT(t)

	unit := &modes.SimulateMode{
		Hoverfly: hoverflySimulateStub{},
	}

	request := models.RequestDetails{
		Destination: "negative-match.com",
	}

	response, err := unit.Process(&http.Request{}, request)
	Expect(err).ToNot(BeNil())

	Expect(response.StatusCode).To(Equal(http.StatusBadGateway))

	responseBody, err := ioutil.ReadAll(response.Body)
	Expect(err).To(BeNil())

	Expect(string(responseBody)).To(ContainSubstring("There was an error when matching"))
	Expect(string(responseBody)).To(ContainSubstring("matching-error"))
}

func Test_SimulateMode_WhenGivenAMatchingRequesAndMiddlewareFaislItReturnsAnError(t *testing.T) {
	RegisterTestingT(t)

	unit := &modes.SimulateMode{
		Hoverfly: hoverflySimulateStub{},
	}

	request := models.RequestDetails{
		Destination: "positive-match.com",
		Path:        "middleware-error",
	}

	response, err := unit.Process(&http.Request{}, request)
	Expect(err).ToNot(BeNil())

	Expect(response.StatusCode).To(Equal(http.StatusBadGateway))

	responseBody, err := ioutil.ReadAll(response.Body)
	Expect(err).To(BeNil())

	Expect(string(responseBody)).To(ContainSubstring("There was an error when executing middleware"))
	Expect(string(responseBody)).To(ContainSubstring("middleware-error"))
}

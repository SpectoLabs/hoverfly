package middleware

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/models"
	. "github.com/onsi/gomega"
)

func TestChangeBodyMiddleware(t *testing.T) {
	RegisterTestingT(t)

	resp := models.ResponseDetails{Status: 201, Body: "original body"}
	req := models.RequestDetails{Path: "/", Method: "GET", Destination: "hostname-x"}

	originalPair := models.RequestResponsePair{Response: resp, Request: req}

	unit := &Middleware{}

	err := unit.SetBinary("python")
	Expect(err).To(BeNil())

	err = unit.SetScript(pythonModifyResponse)
	Expect(err).To(BeNil())

	newPair, err := unit.executeMiddlewareLocally(originalPair)

	Expect(err).To(BeNil())
	Expect(newPair.Response.Body).To(Equal("body was replaced by middleware"))
}

func TestMalformedRequestResponsePairWithMiddleware(t *testing.T) {
	RegisterTestingT(t)

	resp := models.ResponseDetails{Status: 201, Body: "original body"}
	req := models.RequestDetails{Path: "/", Method: "GET", Destination: "hostname-x"}

	malformedPair := models.RequestResponsePair{Response: resp, Request: req}

	unit := &Middleware{}

	err := unit.SetBinary("ruby")
	Expect(err).To(BeNil())

	err = unit.SetScript(rubyEcho)
	Expect(err).To(BeNil())

	newPair, err := unit.executeMiddlewareLocally(malformedPair)

	Expect(err).To(BeNil())
	Expect(newPair.Response.Body).To(Equal("original body"))
}

func TestReflectBody(t *testing.T) {
	RegisterTestingT(t)

	req := models.RequestDetails{Path: "/", Method: "GET", Destination: "hostname-x", Body: "request_body_here"}

	originalPair := models.RequestResponsePair{Request: req}

	unit := &Middleware{}

	err := unit.SetBinary("python")
	Expect(err).To(BeNil())

	err = unit.SetScript(pythonReflectBody)
	Expect(err).To(BeNil())

	newPair, err := unit.executeMiddlewareLocally(originalPair)

	Expect(err).To(BeNil())
	Expect(newPair.Response.Body).To(Equal(req.Body))
	Expect(newPair.Request.Method).To(Equal(req.Method))
	Expect(newPair.Request.Destination).To(Equal(req.Destination))
}

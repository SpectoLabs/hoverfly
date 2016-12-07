package hoverfly

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SpectoLabs/hoverfly/core/handlers/v1"
	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/gorilla/mux"
	. "github.com/onsi/gomega"
)

func TestChangeBodyMiddleware(t *testing.T) {
	RegisterTestingT(t)

	resp := models.ResponseDetails{Status: 201, Body: "original body"}
	req := models.RequestDetails{Path: "/", Method: "GET", Destination: "hostname-x", Query: ""}

	originalPair := models.RequestResponsePair{Response: resp, Request: req}

	unit := &Middleware{
		Script: "./examples/middleware/modify_response/modify_response.py",
	}

	newPair, err := unit.ExecuteMiddlewareLocally(originalPair)

	Expect(err).To(BeNil())
	Expect(newPair.Response.Body).To(Equal("body was replaced by middleware\n"))
}

func TestMalformedRequestResponsePairWithMiddleware(t *testing.T) {
	RegisterTestingT(t)

	resp := models.ResponseDetails{Status: 201, Body: "original body"}
	req := models.RequestDetails{Path: "/", Method: "GET", Destination: "hostname-x", Query: ""}

	malformedPair := models.RequestResponsePair{Response: resp, Request: req}

	unit := &Middleware{
		Script: "./examples/middleware/ruby_echo/echo.rb",
	}

	newPair, err := unit.ExecuteMiddlewareLocally(malformedPair)

	Expect(err).To(BeNil())
	Expect(newPair.Response.Body).To(Equal("original body"))
}

func TestMakeCustom404(t *testing.T) {
	RegisterTestingT(t)

	resp := models.ResponseDetails{Status: 201, Body: "original body"}
	req := models.RequestDetails{Path: "/", Method: "GET", Destination: "hostname-x", Query: ""}

	originalPair := models.RequestResponsePair{Response: resp, Request: req}

	unit := &Middleware{
		Script: "go run ./examples/middleware/go_example/change_to_custom_404.go",
	}

	newPair, err := unit.ExecuteMiddlewareLocally(originalPair)

	Expect(err).To(BeNil())
	Expect(newPair.Response.Body).To(Equal("Custom body here"))
	Expect(newPair.Response.Status).To(Equal(http.StatusNotFound))
	Expect(newPair.Response.Headers["middleware"][0]).To(Equal("changed response"))
}

func TestReflectBody(t *testing.T) {
	RegisterTestingT(t)

	req := models.RequestDetails{Path: "/", Method: "GET", Destination: "hostname-x", Query: "", Body: "request_body_here"}

	originalPair := models.RequestResponsePair{Request: req}

	unit := &Middleware{
		Script: "./examples/middleware/reflect_body/reflect_body.py",
	}

	newPair, err := unit.ExecuteMiddlewareLocally(originalPair)

	Expect(err).To(BeNil())
	Expect(newPair.Response.Body).To(Equal(req.Body))
	Expect(newPair.Request.Method).To(Equal(req.Method))
	Expect(newPair.Request.Destination).To(Equal(req.Destination))
}

func processHandlerOkay(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)

	var newPairView v1.RequestResponsePairView

	json.Unmarshal(body, &newPairView)

	newPairView.Response.Body = "You got straight up messed with"

	pairViewBytes, _ := json.Marshal(newPairView)
	w.Write(pairViewBytes)
}

func processHandlerOkayButNoResponse(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}

func processHandlerNotOkay(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
}

func TestExecuteMiddlewareRemotely(t *testing.T) {
	RegisterTestingT(t)

	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/process", processHandlerOkay).Methods("POST")
	server := httptest.NewServer(muxRouter)
	defer server.Close()

	originalPair := models.RequestResponsePair{
		Response: models.ResponseDetails{
			Body: "Normal body",
		},
	}

	unit := &Middleware{
		Script: server.URL + "/process",
	}

	newPair, err := unit.ExecuteMiddlewareRemotely(originalPair)
	Expect(err).To(BeNil())

	Expect(newPair).ToNot(Equal(originalPair))
	Expect(newPair.Response.Body).To(Equal("You got straight up messed with"))
}

func TestExecuteMiddlewareRemotely_ReturnsErrorIfDoesntGetA200_AndSameRequestResponsePairs(t *testing.T) {
	RegisterTestingT(t)

	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/process", processHandlerNotOkay).Methods("POST")
	server := httptest.NewServer(muxRouter)
	defer server.Close()

	originalPair := models.RequestResponsePair{
		Response: models.ResponseDetails{
			Body: "Normal body",
		},
	}

	unit := &Middleware{
		Script: server.URL + "/process",
	}

	newPair, err := unit.ExecuteMiddlewareRemotely(originalPair)
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("Error when communicating with remote middleware"))

	Expect(newPair).To(Equal(originalPair))
}

func TestExecuteMiddlewareRemotely_ReturnsErrorIfNoRequestResponsePairOnResponse_TheUntouchedPairIsReturned(t *testing.T) {
	RegisterTestingT(t)

	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/process", processHandlerOkayButNoResponse).Methods("POST")
	server := httptest.NewServer(muxRouter)
	defer server.Close()

	originalPair := models.RequestResponsePair{
		Response: models.ResponseDetails{
			Body: "Normal body",
		},
	}

	unit := &Middleware{
		Script: server.URL + "/process",
	}

	untouchedPair, err := unit.ExecuteMiddlewareRemotely(originalPair)
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("unexpected end of JSON input"))

	Expect(untouchedPair).To(Equal(originalPair))
}

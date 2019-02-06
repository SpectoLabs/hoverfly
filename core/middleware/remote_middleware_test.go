package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/gorilla/mux"
	. "github.com/onsi/gomega"
)

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

	unit := &Middleware{}

	err := unit.SetRemote(server.URL + "/process")
	Expect(err).To(BeNil())

	newPair, err := unit.executeMiddlewareRemotely(originalPair)
	Expect(err).To(BeNil())

	Expect(newPair).ToNot(Equal(originalPair))
	Expect(newPair.Response.Body).To(Equal("You got straight up messed with"))
}

func Test_Middleware_executeMiddlewareRemotely_ReturnsErrorIfDoesntGetA200_AndSameRequestResponsePairs(t *testing.T) {
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

	unit := &Middleware{}

	unit.Remote = server.URL + "/process"

	newPair, err := unit.executeMiddlewareRemotely(originalPair)
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(ContainSubstring("Error when communicating with remote middleware: received 404"))
	Expect(err.Error()).To(ContainSubstring("URL: " + server.URL))
	Expect(err.Error()).To(ContainSubstring("STDIN:"))
	Expect(err.Error()).To(ContainSubstring(`{"response":{"status":0,"body":"Normal body","encodedBody":false},"request":{"path":"","method":"","destination":"","scheme":"","query":"","body":"","headers":null}}`))

	Expect(newPair).To(Equal(originalPair))
}

func Test_Middleware_executeMiddlewareRemotely_ReturnsErrorIfNoRequestResponsePairOnResponse_TheUntouchedPairIsReturned(t *testing.T) {
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

	unit := &Middleware{}

	err := unit.SetRemote(server.URL + "/process")
	Expect(err).To(BeNil())

	untouchedPair, err := unit.executeMiddlewareRemotely(originalPair)
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(ContainSubstring("Error when trying to serialize response from remote middleware"))
	Expect(err.Error()).To(ContainSubstring("URL: " + server.URL))
	Expect(err.Error()).To(ContainSubstring("STDIN:"))
	Expect(err.Error()).To(ContainSubstring(`{"response":{"status":0,"body":"Normal body","encodedBody":false},"request":{"path":"","method":"","destination":"","scheme":"","query":"","body":"","headers":null}}`))

	Expect(untouchedPair).To(Equal(originalPair))
}

func Test_Middleware_executeMiddlewareRemotely_ReturnsError_WebsiteIsUnreachable(t *testing.T) {
	RegisterTestingT(t)

	originalPair := models.RequestResponsePair{
		Response: models.ResponseDetails{
			Body: "Normal body",
		},
	}

	unit := &Middleware{}

	unit.Remote = "[]somemadeupwebsite"

	untouchedPair, err := unit.executeMiddlewareRemotely(originalPair)
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(ContainSubstring("Error when communicating with remote middleware:"))
	Expect(err.Error()).To(ContainSubstring("Post []somemadeupwebsite: unsupported protocol scheme"))
	Expect(err.Error()).To(ContainSubstring("URL: []somemadeupwebsite"))
	Expect(err.Error()).To(ContainSubstring("STDIN:"))
	Expect(err.Error()).To(ContainSubstring(`{"response":{"status":0,"body":"Normal body","encodedBody":false},"request":{"path":"","method":"","destination":"","scheme":"","query":"","body":"","headers":null}}`))

	Expect(untouchedPair).To(Equal(originalPair))

	unit.Remote = "http://localhost:4321/spectolabs/hoverfly"

	untouchedPair, err = unit.executeMiddlewareRemotely(originalPair)
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(ContainSubstring("Error when communicating with remote middleware:"))
	Expect(err.Error()).To(MatchRegexp("Post http://localhost:4321/spectolabs/hoverfly: dial tcp .+:4321: connect: connection refused"))
	Expect(err.Error()).To(ContainSubstring("URL: http://localhost:4321/spectolabs/hoverfly"))
	Expect(err.Error()).To(ContainSubstring("STDIN:"))
	Expect(err.Error()).To(ContainSubstring(`{"response":{"status":0,"body":"Normal body","encodedBody":false},"request":{"path":"","method":"","destination":"","scheme":"","query":"","body":"","headers":null}}`))

	Expect(untouchedPair).To(Equal(originalPair))
}

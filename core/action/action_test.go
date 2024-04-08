package action_test

import (
	"github.com/SpectoLabs/hoverfly/core/journal"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SpectoLabs/hoverfly/core/action"
	"github.com/SpectoLabs/hoverfly/core/models"
	. "github.com/onsi/gomega"
)

const pythonBasicScript = "import sys\nprint(sys.stdin.readlines()[0])"

func Test_NewLocalActionMethod(t *testing.T) {
	RegisterTestingT(t)

	newAction, err := action.NewLocalAction("test-callback", "python3", "dummy-script", 1800)

	Expect(err).To(BeNil())
	Expect(newAction).NotTo(BeNil())
	Expect(newAction.Binary).To(Equal("python3"))
	Expect(newAction.DelayInMs).To(Equal(1800))

	scriptContent, err := newAction.GetScript()

	Expect(err).To(BeNil())
	Expect(scriptContent).To(Equal("dummy-script"))
}

func Test_NewRemoteActionMethodWithEmptyHost(t *testing.T) {
	RegisterTestingT(t)

	newAction, err := action.NewRemoteAction("test-callback", "", 1800)

	Expect(err).NotTo(BeNil())
	Expect(newAction).To(BeNil())
}

func Test_NewRemoteActionMethodWithInvalidHost(t *testing.T) {
	RegisterTestingT(t)

	newAction, err := action.NewRemoteAction("test-callback", "testing", 1800)

	Expect(err).NotTo(BeNil())
	Expect(err.Error()).To(Equal("remote host is invalid"))
	Expect(newAction).To(BeNil())
}

func Test_NewRemoteActionMethodWithHttpHost(t *testing.T) {
	RegisterTestingT(t)

	newAction, err := action.NewRemoteAction("test-callback", "http://localhost", 1800)

	Expect(err).To(BeNil())
	Expect(newAction).NotTo(BeNil())
	Expect(newAction.Remote).To(Equal("http://localhost"))
	Expect(newAction.DelayInMs).To(Equal(1800))
}

func Test_NewRemoteActionMethodWithHttpsHost(t *testing.T) {
	RegisterTestingT(t)

	newAction, err := action.NewRemoteAction("test-callback", "https://test.com", 1800)

	Expect(err).To(BeNil())
	Expect(newAction).NotTo(BeNil())
	Expect(newAction.Remote).To(Equal("https://test.com"))
	Expect(newAction.DelayInMs).To(Equal(1800))
}

func Test_GetLocalActionViewMethod(t *testing.T) {
	RegisterTestingT(t)

	newAction, err := action.NewLocalAction("test-callback", "python3", "dummy-script", 1800)

	Expect(err).To(BeNil())
	actionView := newAction.GetActionView("test-callback")

	Expect(actionView.ActionName).To(Equal("test-callback"))
	Expect(actionView.Binary).To(Equal("python3"))
	Expect(actionView.ScriptContent).To(Equal("dummy-script"))
	Expect(actionView.DelayInMs).To(Equal(1800))
}

func Test_GetRemoteActionViewMethod(t *testing.T) {
	RegisterTestingT(t)

	newAction, err := action.NewRemoteAction("test-callback", "http://localhost:8000", 1800)

	Expect(err).To(BeNil())
	actionView := newAction.GetActionView("test-callback")

	Expect(actionView.ActionName).To(Equal("test-callback"))
	Expect(actionView.Binary).To(Equal(""))
	Expect(actionView.ScriptContent).To(Equal(""))
	Expect(actionView.Remote).To(Equal("http://localhost:8000"))
	Expect(actionView.DelayInMs).To(Equal(1800))
}

func Test_ExecuteLocalPostServeAction(t *testing.T) {
	RegisterTestingT(t)
	newAction, err := action.NewLocalAction("test-callback", "python3", pythonBasicScript, 0)

	Expect(err).To(BeNil())

	resp := models.ResponseDetails{Status: 201, Body: "original body"}
	req := models.RequestDetails{Path: "/", Method: "GET", Destination: "hostname-x"}

	originalPair := models.RequestResponsePair{Response: resp, Request: req}

	//not adding entry as update journal method will be tested in its file
	journalIDChannel := make(chan string)
	newJournal := journal.NewJournal()
	journalIDChannel <- "1"
	err = newAction.Execute(&originalPair, journalIDChannel, newJournal)
	Expect(err).To(BeNil())
}

func Test_ExecuteRemotePostServeAction(t *testing.T) {
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
	//not adding entry as update journal method will be tested in its file
	journalIDChannel := make(chan string)
	newJournal := journal.NewJournal()
	journalIDChannel <- "1"

	newAction, err := action.NewRemoteAction("test-callback", server.URL+"/process", 0)
	Expect(err).To(BeNil())
	err = newAction.Execute(&originalPair, journalIDChannel, newJournal)
	Expect(err).To(BeNil())
}

func Test_ExecuteRemotePostServeAction_WithUnReachableHost(t *testing.T) {
	originalPair := models.RequestResponsePair{
		Response: models.ResponseDetails{
			Body: "Normal body",
		},
	}

	newAction, err := action.NewRemoteAction("test-callback", "http://test", 0)
	Expect(err).To(BeNil())

	//not adding entry as update journal method will be tested in its file
	journalIDChannel := make(chan string)
	newJournal := journal.NewJournal()
	journalIDChannel <- "1"
	err = newAction.Execute(&originalPair, journalIDChannel, newJournal)
	Expect(err).NotTo(BeNil())
}

func processHandlerOkay(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}

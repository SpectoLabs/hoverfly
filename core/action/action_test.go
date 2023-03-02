package action_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/action"
	"github.com/SpectoLabs/hoverfly/core/models"
	. "github.com/onsi/gomega"
)

const pythonBasicScript = "import sys\nprint(sys.stdin.readlines()[0])"

func Test_NewActionMethod(t *testing.T) {
	RegisterTestingT(t)

	newAction, err := action.NewAction("test-callback", "python3", "dummy-script", 1800)

	Expect(err).To(BeNil())
	Expect(newAction).NotTo(BeNil())
	Expect(newAction.Binary).To(Equal("python3"))
	Expect(newAction.DelayInMs).To(Equal(1800))

	scriptContent, err := newAction.GetScript()

	Expect(err).To(BeNil())
	Expect(scriptContent).To(Equal("dummy-script"))
}

func Test_GetActionViewMethod(t *testing.T) {
	RegisterTestingT(t)

	newAction, err := action.NewAction("test-callback", "python3", "dummy-script", 1800)

	Expect(err).To(BeNil())
	actionView := newAction.GetActionView("test-callback")

	Expect(actionView.ActionName).To(Equal("test-callback"))
	Expect(actionView.Binary).To(Equal("python3"))
	Expect(actionView.ScriptContent).To(Equal("dummy-script"))
	Expect(actionView.DelayInMs).To(Equal(1800))
}

func Test_ExecuteLocallyPostServeAction(t *testing.T) {
	RegisterTestingT(t)
	newAction, err := action.NewAction("test-callback", "python3", pythonBasicScript, 0)

	Expect(err).To(BeNil())

	resp := models.ResponseDetails{Status: 201, Body: "original body"}
	req := models.RequestDetails{Path: "/", Method: "GET", Destination: "hostname-x"}

	originalPair := models.RequestResponsePair{Response: resp, Request: req}

	err = newAction.ExecuteLocally(&originalPair)
	Expect(err).To(BeNil())
}

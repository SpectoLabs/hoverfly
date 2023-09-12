package v2

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	. "github.com/onsi/gomega"
)

type HoverflyPostServeActionDetailsStub struct{}

var actionView = ActionView{ActionName: "dummy-action", Binary: "python3", ScriptContent: "Dummy Python Script", DelayInMs: 1000}

func (HoverflyPostServeActionDetailsStub) GetAllPostServeActions() PostServeActionDetailsView {

	actions := []ActionView{actionView}
	return PostServeActionDetailsView{Actions: actions}
}

func (HoverflyPostServeActionDetailsStub) SetPostServeAction(string, string, string, int) error {
	return nil
}

func (HoverflyPostServeActionDetailsStub) DeletePostServeAction(string) error {
	return nil
}

func Test_PostServeActionHandler_DeletePostServeAction(t *testing.T) {
	RegisterTestingT(t)

	stubHoverfly := &HoverflyPostServeActionDetailsStub{}
	unit := HoverflyPostServeActionDetailsHandler{Hoverfly: stubHoverfly}

	request, err := http.NewRequest("DELETE", "/api/v2/hoverfly/post-serve-action/test-action", nil)
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Delete, request)
	Expect(response.Code).To(Equal(http.StatusOK))
}

func Test_PostServeActionHandler_SetPostServeAction(t *testing.T) {
	RegisterTestingT(t)

	stubHoverfly := &HoverflyPostServeActionDetailsStub{}
	unit := HoverflyPostServeActionDetailsHandler{Hoverfly: stubHoverfly}

	actionView := &ActionView{ActionName: "test-action", Binary: "python", ScriptContent: "dummy script", DelayInMs: 1000}

	bodyBytes, err := json.Marshal(actionView)
	Expect(err).To(BeNil())

	request, err := http.NewRequest("PUT", "/api/v2/hoverfly/post-serve-action", ioutil.NopCloser(bytes.NewBuffer(bodyBytes)))
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Put, request)
	Expect(response.Code).To(Equal(http.StatusOK))
}

func Test_PostServeActionHandler_GetAllPostServeActions(t *testing.T) {
	RegisterTestingT(t)
	stubHoverfly := &HoverflyPostServeActionDetailsStub{}
	unit := HoverflyPostServeActionDetailsHandler{Hoverfly: stubHoverfly}

	request, err := http.NewRequest("GET", "/api/v2/hoverfly/post-serve-action", nil)
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Get, request)
	Expect(response.Code).To(Equal(http.StatusOK))

	responseBody := response.Body.String()
	expectedResponseBodyBytes, _ := json.Marshal(PostServeActionDetailsView{Actions: []ActionView{actionView}})

	Expect(responseBody).To(Equal(string(expectedResponseBodyBytes)))

}

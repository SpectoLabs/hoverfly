package hoverfly_end_to_end_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
	"github.com/dghubble/sling"
	"strings"
	"net/http"
	"fmt"
	"encoding/json"
	"io/ioutil"
)

func TestHoverflyEndToEnd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Hoverfly End To End Suite")
}

func SetHoverflyMode(mode string, port int) {
	req := sling.New().Post(fmt.Sprintf("http://localhost:%v/api/state", port)).Body(strings.NewReader(`{"mode":"` + mode +`"}`))
	res := DoRequest(req)
	Expect(res.StatusCode).To(Equal(200))
}

func DoRequest(r *sling.Sling) (*http.Response) {
	req, err := r.Request()
	Expect(err).To(BeNil())
	response, err := http.DefaultClient.Do(req)

	Expect(err).To(BeNil())
	return response
}

func GetHoverflyMode(port int) string {
	currentState := &stateRequest{}
	resp := DoRequest(sling.New().Get(fmt.Sprintf("http://localhost:%v/api/state", port)))

	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).To(BeNil())

	err = json.Unmarshal(body, currentState)
	Expect(err).To(BeNil())

	return currentState.Mode
}

type stateRequest struct {
	Mode        string `json:"mode"`
	Destination string `json:"destination"`
}
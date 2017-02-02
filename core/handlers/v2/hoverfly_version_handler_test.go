package v2

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	. "github.com/onsi/gomega"
)

type HoverflyVersionStub struct{}

func (this HoverflyVersionStub) GetVersion() string {
	return "test-version"
}

func Test_HoverflyUsageHandler_GetReturnsVersion(t *testing.T) {
	RegisterTestingT(t)

	stubHoverfly := &HoverflyVersionStub{}
	unit := HoverflyVersionHandler{Hoverfly: stubHoverfly}

	request, err := http.NewRequest("GET", "", nil)
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Get, request)

	Expect(response.Code).To(Equal(http.StatusOK))

	versionView, err := unmarshalVersionView(response.Body)
	Expect(err).To(BeNil())

	Expect(versionView.Version).To(Equal("test-version"))
}

func unmarshalVersionView(buffer *bytes.Buffer) (VersionView, error) {
	body, err := ioutil.ReadAll(buffer)
	if err != nil {
		return VersionView{}, err
	}

	var versionView VersionView

	err = json.Unmarshal(body, &versionView)
	if err != nil {
		return VersionView{}, err
	}

	return versionView, nil
}

package v2

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	. "github.com/onsi/gomega"
)

type HoverflyJournalIndexStub struct{}

func (HoverflyJournalIndexStub) GetAllIndexes() []JournalIndexView {

	return getJournalIndexes()
}

func getJournalIndexes() []JournalIndexView {
	journalIndexViews := []JournalIndexView{}
	journalIndexViews = append(journalIndexViews, JournalIndexView{Name: "Request.QueryParam.id"})
	return journalIndexViews
}
func (HoverflyJournalIndexStub) AddJournalIndex(string) error {
	return nil
}
func (HoverflyJournalIndexStub) DeleteJournalIndex(string) {

}

func Test_JournalIndexHandler_GetAllJournalIndexes(t *testing.T) {
	RegisterTestingT(t)
	stubHoverfly := HoverflyJournalIndexStub{}
	unit := HoverflyJournalIndexHandler{Hoverfly: stubHoverfly}

	request, err := http.NewRequest("GET", "/api/v2/journal/index", nil)
	Expect(err).To(BeNil())
	response := makeRequestOnHandler(unit.GetAll, request)
	Expect(response.Code).To(Equal(http.StatusOK))

	responseBody := response.Body.String()
	expectedResponseBodyBytes, _ := json.Marshal(getJournalIndexes())
	Expect(responseBody).To(Equal(string(expectedResponseBodyBytes)))

}

func Test_JournalIndexHandler_SetJournalIndex(t *testing.T) {
	RegisterTestingT(t)
	stubHoverfly := HoverflyJournalIndexStub{}
	unit := HoverflyJournalIndexHandler{Hoverfly: stubHoverfly}

	journalIndexRequest := JournalIndexRequestView{
		Name: "Request.QueryParams.id",
	}
	bodyBytes, err := json.Marshal(journalIndexRequest)
	Expect(err).To(BeNil())
	request, err := http.NewRequest("POST", "/api/v2/journal/index", io.NopCloser(bytes.NewBuffer(bodyBytes)))

	Expect(err).To(BeNil())
	response := makeRequestOnHandler(unit.Post, request)
	Expect(response.Code).To(Equal(http.StatusOK))
}

func Test_JournalIndexHandler_DeleteJournalIndex(t *testing.T) {
	RegisterTestingT(t)
	stubHoverfly := HoverflyJournalIndexStub{}
	unit := HoverflyJournalIndexHandler{Hoverfly: stubHoverfly}
	request, err := http.NewRequest("DELETE", "/api/v2/journal/index/Request.QueryParams.id", nil)

	Expect(err).To(BeNil())
	response := makeRequestOnHandler(unit.Delete, request)
	Expect(response.Code).To(Equal(http.StatusOK))
}

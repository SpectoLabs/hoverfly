package v2

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"fmt"

	. "github.com/onsi/gomega"
)

type HoverflyJournalStub struct {
	limit   int
	deleted bool
	error   bool
}

func (this *HoverflyJournalStub) GetEntries() ([]JournalEntryView, error) {
	if this.error {
		return []JournalEntryView{}, fmt.Errorf("entries error")
	}

	if this.deleted {
		return []JournalEntryView{}, nil
	} else {
		return []JournalEntryView{
			JournalEntryView{
				Mode: "test",
			},
		}, nil
	}
}

func (this *HoverflyJournalStub) DeleteEntries() error {
	if this.error {
		return fmt.Errorf("delete error")
	}

	this.deleted = true
	return nil
}

func Test_JournalHandler_Get_ReturnsJournal(t *testing.T) {
	RegisterTestingT(t)

	stubHoverfly := &HoverflyJournalStub{}
	unit := JournalHandler{Hoverfly: stubHoverfly}

	request, err := http.NewRequest("GET", "/api/v2/journal", nil)
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Get, request)

	Expect(response.Code).To(Equal(http.StatusOK))

	journalView, err := unmarshalJournalEntryView(response.Body)
	Expect(err).To(BeNil())

	Expect(journalView).To(HaveLen(1))
	Expect(journalView[0].Mode).To(Equal("test"))
}

func Test_JournalHandler_Get_Error(t *testing.T) {
	RegisterTestingT(t)

	stubHoverfly := HoverflyJournalStub{
		error: true,
	}
	unit := JournalHandler{Hoverfly: &stubHoverfly}

	request, err := http.NewRequest("GET", "/api/v2/journal", nil)
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Get, request)

	Expect(response.Code).To(Equal(http.StatusInternalServerError))

	errorView, err := unmarshalErrorView(response.Body)
	Expect(err).To(BeNil())

	Expect(errorView.Error).To(Equal("entries error"))
}

func Test_JournalHandler_Delete_CallsDelete(t *testing.T) {
	RegisterTestingT(t)

	var stubHoverfly HoverflyJournalStub
	unit := JournalHandler{Hoverfly: &stubHoverfly}

	request, err := http.NewRequest("DELETE", "/api/v2/journal", nil)
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Delete, request)

	Expect(response.Code).To(Equal(http.StatusOK))

	journalView, err := unmarshalJournalEntryView(response.Body)
	Expect(err).To(BeNil())

	Expect(journalView).To(HaveLen(0))

	Expect(stubHoverfly.deleted).To(BeTrue())
}

func Test_JournalHandler_Delete_Error(t *testing.T) {
	RegisterTestingT(t)

	stubHoverfly := HoverflyJournalStub{
		error: true,
	}
	unit := JournalHandler{Hoverfly: &stubHoverfly}

	request, err := http.NewRequest("DELETE", "/api/v2/journal", nil)
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Delete, request)

	Expect(response.Code).To(Equal(http.StatusInternalServerError))

	errorView, err := unmarshalErrorView(response.Body)
	Expect(err).To(BeNil())

	Expect(errorView.Error).To(Equal("delete error"))
}

func Test_JournalHandler_Options_GetsOptions(t *testing.T) {
	RegisterTestingT(t)

	var stubHoverfly HoverflyJournalStub
	unit := JournalHandler{Hoverfly: &stubHoverfly}

	request, err := http.NewRequest("OPTIONS", "/api/v2/journal", nil)
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Options, request)

	Expect(response.Code).To(Equal(http.StatusOK))
	Expect(response.Header().Get("Allow")).To(Equal("OPTIONS, GET, DELETE"))
}

func unmarshalJournalEntryView(buffer *bytes.Buffer) ([]JournalEntryView, error) {
	body, err := ioutil.ReadAll(buffer)
	if err != nil {
		return []JournalEntryView{}, err
	}

	var journalView []JournalEntryView

	err = json.Unmarshal(body, &journalView)
	if err != nil {
		return []JournalEntryView{}, err
	}

	return journalView, nil
}

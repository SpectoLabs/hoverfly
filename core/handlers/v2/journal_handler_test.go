package v2

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"fmt"

	"time"

	"github.com/SpectoLabs/hoverfly/core/matching/matchers"
	. "github.com/onsi/gomega"
)

type HoverflyJournalStub struct {
	deleted                bool
	error                  bool
	limit                  int
	offset                 int
	sort                   string
	from                   *time.Time
	to                     *time.Time
	journalEntryFilterView JournalEntryFilterView
}

func (this *HoverflyJournalStub) GetEntries(offset int, limit int, from *time.Time, to *time.Time, sort string) (JournalView, error) {
	this.offset = offset
	this.limit = limit
	this.from = from
	this.to = to
	this.sort = sort

	journalView := JournalView{
		Journal: []JournalEntryView{},
		Offset:  offset,
		Limit:   limit,
		Total:   0,
	}
	if this.error {
		return journalView, fmt.Errorf("entries error")
	}

	if !this.deleted {
		journalView.Journal = []JournalEntryView{{
			Mode: "test",
		}}
		journalView.Total = 1
	}

	return journalView, nil
}

func (this *HoverflyJournalStub) GetFilteredEntries(journalEntryFilterView JournalEntryFilterView) ([]JournalEntryView, error) {
	if this.error {
		return []JournalEntryView{}, fmt.Errorf("journal error")
	}

	this.journalEntryFilterView = journalEntryFilterView
	return []JournalEntryView{
		{
			Mode: "test",
		},
	}, nil
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

	journalView, err := unmarshalJournalView(response.Body)
	Expect(err).To(BeNil())

	Expect(journalView.Journal).To(HaveLen(1))
	Expect(journalView.Journal[0].Mode).To(Equal("test"))
}

func Test_JournalHandler_Get_SetDefaultPagingQueryIfNotSpecified(t *testing.T) {
	RegisterTestingT(t)

	stubHoverfly := &HoverflyJournalStub{}

	unit := JournalHandler{Hoverfly: stubHoverfly}

	request, err := http.NewRequest("GET", "/api/v2/journal", nil)
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Get, request)

	Expect(response.Code).To(Equal(http.StatusOK))

	Expect(stubHoverfly.limit).To(Equal(25))
	Expect(stubHoverfly.offset).To(Equal(0))
}

func Test_JournalHandler_Get_WithPagingQuery(t *testing.T) {
	RegisterTestingT(t)

	stubHoverfly := &HoverflyJournalStub{}

	unit := JournalHandler{Hoverfly: stubHoverfly}

	request, err := http.NewRequest("GET", "/api/v2/journal?offset=50&limit=25", nil)
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Get, request)

	Expect(response.Code).To(Equal(http.StatusOK))

	Expect(stubHoverfly.limit).To(Equal(25))
	Expect(stubHoverfly.offset).To(Equal(50))

	journalView, err := unmarshalJournalView(response.Body)
	Expect(err).To(BeNil())

	Expect(journalView.Total).To(Equal(1))
	Expect(journalView.Offset).To(Equal(50))
	Expect(journalView.Limit).To(Equal(25))
}

func Test_JournalHandler_Get_WithSortQuery(t *testing.T) {
	RegisterTestingT(t)

	stubHoverfly := &HoverflyJournalStub{}

	unit := JournalHandler{Hoverfly: stubHoverfly}

	request, err := http.NewRequest("GET", "/api/v2/journal?sort=timeStarted:desc", nil)
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Get, request)

	Expect(response.Code).To(Equal(http.StatusOK))

	Expect(stubHoverfly.sort).To(Equal("timeStarted:desc"))
}

func Test_JournalHandler_Get_WithDateTimeQuery(t *testing.T) {
	RegisterTestingT(t)

	stubHoverfly := &HoverflyJournalStub{}

	unit := JournalHandler{Hoverfly: stubHoverfly}

	request, err := http.NewRequest("GET", "/api/v2/journal?from=1517498951&to=1517498986", nil)
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Get, request)

	Expect(response.Code).To(Equal(http.StatusOK))

	Expect(stubHoverfly.from.Unix()).To(Equal(int64(1517498951)))
	Expect(stubHoverfly.to.Unix()).To(Equal(int64(1517498986)))
}

func Test_JournalHandler_Get_DoesNotSetTimeIfDateTimeQueryIsBadTime(t *testing.T) {
	RegisterTestingT(t)

	stubHoverfly := &HoverflyJournalStub{}

	unit := JournalHandler{Hoverfly: stubHoverfly}

	request, err := http.NewRequest("GET", "/api/v2/journal?from=bad-time", nil)
	Expect(err).To(BeNil())

	makeRequestOnHandler(unit.Get, request)

	Expect(stubHoverfly.from).To(BeNil())
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

func Test_JournalHandler_Post_CallsFilter(t *testing.T) {
	RegisterTestingT(t)

	var stubHoverfly HoverflyJournalStub
	unit := JournalHandler{Hoverfly: &stubHoverfly}

	journalEntryFilterView := JournalEntryFilterView{
		Request: &RequestMatcherViewV5{
			Destination: []MatcherViewV5{
				NewMatcherView(matchers.Exact, "hoverfly.io"),
			},
			Path: []MatcherViewV5{
				NewMatcherView(matchers.Glob, "*"),
			},
		},
	}

	body, _ := json.Marshal(journalEntryFilterView)

	request, err := http.NewRequest("POST", "/api/v2/journal", bytes.NewBuffer(body))
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Post, request)

	Expect(response.Code).To(Equal(http.StatusOK))

	journalView, err := unmarshalJournalView(response.Body)
	Expect(err).To(BeNil())

	Expect(journalView.Journal).To(HaveLen(1))
	Expect(journalView.Journal[0].Mode).To(Equal("test"))

	Expect(stubHoverfly.journalEntryFilterView.Request.Destination[0].Value).To(Equal("hoverfly.io"))
	Expect(stubHoverfly.journalEntryFilterView.Request.Path[0].Value).To(Equal("*"))
}

func Test_JournalHandler_Post_MalformedJson(t *testing.T) {
	RegisterTestingT(t)

	var stubHoverfly HoverflyJournalStub

	unit := JournalHandler{Hoverfly: &stubHoverfly}

	request, err := http.NewRequest("POST", "/api/v2/journal", bytes.NewBufferString("werw{{}[][{}"))
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Post, request)

	Expect(response.Code).To(Equal(http.StatusBadRequest))

	errorView, err := unmarshalErrorView(response.Body)
	Expect(err).To(BeNil())

	Expect(errorView.Error).To(Equal("Malformed JSON"))
}

func Test_JournalHandler_Post_MalformedJson_EmptyRequest(t *testing.T) {
	RegisterTestingT(t)

	var stubHoverfly HoverflyJournalStub

	unit := JournalHandler{Hoverfly: &stubHoverfly}

	request, err := http.NewRequest("POST", "/api/v2/journal", bytes.NewBufferString("{}"))
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Post, request)

	Expect(response.Code).To(Equal(http.StatusBadRequest))

	errorView, err := unmarshalErrorView(response.Body)
	Expect(err).To(BeNil())

	Expect(errorView.Error).To(Equal("No \"request\" object in search parameters"))
}

func Test_JournalHandler_Post_JournalError(t *testing.T) {
	RegisterTestingT(t)

	stubHoverfly := HoverflyJournalStub{
		error: true,
	}

	unit := JournalHandler{Hoverfly: &stubHoverfly}

	requestMatcher := JournalEntryFilterView{
		Request: &RequestMatcherViewV5{
			Destination: []MatcherViewV5{
				NewMatcherView(matchers.Exact, "hoverfly.io"),
			},
			Path: []MatcherViewV5{
				NewMatcherView(matchers.Glob, "*"),
			},
		},
	}

	body, _ := json.Marshal(requestMatcher)

	request, err := http.NewRequest("POST", "/api/v2/journal", bytes.NewBuffer(body))
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Post, request)

	Expect(response.Code).To(Equal(http.StatusInternalServerError))

	errorView, err := unmarshalErrorView(response.Body)
	Expect(err).To(BeNil())

	Expect(errorView.Error).To(Equal("journal error"))
}

func Test_JournalHandler_Delete_CallsDelete(t *testing.T) {
	RegisterTestingT(t)

	var stubHoverfly HoverflyJournalStub
	unit := JournalHandler{Hoverfly: &stubHoverfly}

	request, err := http.NewRequest("DELETE", "/api/v2/journal", nil)
	Expect(err).To(BeNil())

	response := makeRequestOnHandler(unit.Delete, request)

	Expect(response.Code).To(Equal(http.StatusOK))

	journalView, err := unmarshalJournalView(response.Body)
	Expect(err).To(BeNil())

	Expect(journalView.Journal).To(HaveLen(0))

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
	Expect(response.Header().Get("Allow")).To(Equal("OPTIONS, GET, DELETE, POST"))
}

func unmarshalJournalView(buffer *bytes.Buffer) (JournalView, error) {
	body, err := ioutil.ReadAll(buffer)
	if err != nil {
		return JournalView{}, err
	}

	var journalView JournalView

	err = json.Unmarshal(body, &journalView)
	if err != nil {
		return JournalView{}, err
	}

	return journalView, nil
}

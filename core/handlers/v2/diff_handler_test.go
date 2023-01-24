package v2

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"bytes"
	"encoding/json"
	"io/ioutil"
	"time"

	. "github.com/onsi/gomega"
)

type DiffHOverflyStub struct {
	*HoverflyStub
}

func (this *DiffHOverflyStub) GetDiff() map[SimpleRequestDefinitionView][]DiffReport {
	return diffView
}

func (this *DiffHOverflyStub) GetFilteredDiff(DiffFilterView) map[SimpleRequestDefinitionView][]DiffReport {
	return diffView
}

func (this *DiffHOverflyStub) ClearDiff() {
	diffView = make(map[SimpleRequestDefinitionView][]DiffReport)
}

var diffView map[SimpleRequestDefinitionView][]DiffReport

func TestDiffHandlerGetReturnsTheCorrectDiff(t *testing.T) {
	RegisterTestingT(t)

	// given
	initializeDiff()
	unit, request, err := createRequest("GET", nil)

	// when
	response := makeRequestOnHandler(unit.Get, request)

	// then
	Expect(err).To(BeNil())
	Expect(response.Code).To(Equal(http.StatusOK))

	assertResponseDiff(response)
}

func TestDiffHandlerReturnsBadRequestIfInvalidFilteredDataIsPassed(t *testing.T) {
	RegisterTestingT(t)

	unit, request, err := createRequest("POST", bytes.NewReader([]byte("Hello Test")))

	response := makeRequestOnHandler(unit.GetFilteredData, request)

	Expect(err).To(BeNil())
	Expect(response.Code).To(Equal(http.StatusBadRequest))

}
func TestDiffHandlerGetFilteredDiffReturnsTheCorrectDiff(t *testing.T) {
	RegisterTestingT(t)
	// given
	initializeDiff()
	unit, request, err := createRequest("POST", bytes.NewReader([]byte(`{"excludedHeaders": ["test"], "excludedResponseFields":["$.test"]}`)))

	response := makeRequestOnHandler(unit.GetFilteredData, request)

	Expect(err).To(BeNil())
	Expect(response.Code).To(Equal(http.StatusOK))

	assertResponseDiff(response)
}

func TestDiffHandlerDeleteCleansAllStoredDiffs(t *testing.T) {
	RegisterTestingT(t)

	// given
	initializeDiff()
	unit, request, err := createRequest("GET", nil)

	// when
	deleteResponse := makeRequestOnHandler(unit.Delete, request)
	getResponse := makeRequestOnHandler(unit.Get, request)

	// then
	Expect(err).To(BeNil())
	Expect(deleteResponse.Code).To(Equal(http.StatusOK))
	Expect(getResponse.Code).To(Equal(http.StatusOK))

	diffView, err := unmarshalDiffView(getResponse.Body)
	Expect(err).To(BeNil())
	Expect(len(diffView.Diff)).To(Equal(0))
}

func TestDiffHandlerOptionsGetsOptions(t *testing.T) {
	RegisterTestingT(t)
	// given
	initializeDiff()
	unit, request, err := createRequest("OPTIONS", nil)

	// when
	response := makeRequestOnHandler(unit.Options, request)

	//then
	Expect(err).To(BeNil())
	Expect(response.Code).To(Equal(http.StatusOK))
	Expect(response.Header().Get("Allow")).To(Equal("OPTIONS, GET, DELETE"))
}

func createRequest(method string, body io.Reader) (DiffHandler, *http.Request, error) {
	var stubHoverfly DiffHOverflyStub
	unit := DiffHandler{Hoverfly: &stubHoverfly}

	request, err := http.NewRequest(method, "", body)

	return unit, request, err
}

func initializeDiff() {
	diffView = map[SimpleRequestDefinitionView][]DiffReport{
		SimpleRequestDefinitionView{
			Host:   "testHost",
			Method: "testMethod",
			Path:   "testPath",
			Query:  "testQuery",
		}: {
			{
				Timestamp: time.Now().Format(time.RFC3339),
				DiffEntries: []DiffReportEntry{
					{"first", "expected1", "actual1"},
				},
			},
			{
				Timestamp: time.Now().Format(time.RFC3339),
				DiffEntries: []DiffReportEntry{
					{"second", "expected2", "actual2"},
				},
			},
		},
	}
}

func unmarshalDiffView(buffer *bytes.Buffer) (DiffView, error) {
	body, err := ioutil.ReadAll(buffer)
	if err != nil {
		return DiffView{}, err
	}

	var diffView DiffView

	err = json.Unmarshal(body, &diffView)
	if err != nil {
		return DiffView{}, err
	}

	return diffView, nil
}

func assertResponseDiff(response *httptest.ResponseRecorder) {
	diffView, err := unmarshalDiffView(response.Body)
	Expect(err).To(BeNil())
	Expect(len(diffView.Diff)).To(Equal(1))

	req := diffView.Diff[0].Request
	Expect(req.Host).To(Equal("testHost"))
	Expect(req.Method).To(Equal("testMethod"))
	Expect(req.Path).To(Equal("testPath"))
	Expect(req.Query).To(Equal("testQuery"))

	report := diffView.Diff[0].DiffReport
	Expect(len(report)).To(Equal(2))
	Expect(report[0].DiffEntries).To(ConsistOf(DiffReportEntry{"first", "expected1", "actual1"}))
	Expect(report[1].DiffEntries).To(ConsistOf(DiffReportEntry{"second", "expected2", "actual2"}))
}

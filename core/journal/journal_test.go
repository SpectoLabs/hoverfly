package journal_test

import (
	"bytes"
	"github.com/SpectoLabs/hoverfly/core/models"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/journal"
	"github.com/SpectoLabs/hoverfly/core/matching/matchers"
	. "github.com/onsi/gomega"
)

func Test_NewJournal_ProducesAJournalWithAnEmptyArray(t *testing.T) {
	RegisterTestingT(t)

	unit := journal.NewJournal()

	journalView, err := unit.GetEntries(0, 25, nil, nil, "")
	entries := journalView.Journal
	indexes := journalView.Index
	Expect(err).To(BeNil())

	Expect(entries).ToNot(BeNil())
	Expect(entries).To(HaveLen(0))
	Expect(indexes).ToNot(BeNil())
	Expect(indexes).To(HaveLen(0))

	Expect(unit.EntryLimit).To(Equal(1000))
}

func Test_Journal_NewEntry_AddsJournalEntryToEntries(t *testing.T) {
	RegisterTestingT(t)

	unit := journal.NewJournal()

	request, _ := http.NewRequest("GET", "http://hoverfly.io", nil)

	nowTime := time.Now()

	_, err := unit.NewEntry(request, &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString("test body")),
		Header: http.Header{
			"test-header": []string{
				"one", "two",
			},
		},
	}, "test-mode", nowTime)
	Expect(err).To(BeNil())

	journalView, err := unit.GetEntries(0, 25, nil, nil, "")
	entries := journalView.Journal
	Expect(err).To(BeNil())

	Expect(entries).ToNot(BeNil())
	Expect(entries).To(HaveLen(1))

	Expect(*entries[0].Request.Method).To(Equal("GET"))
	Expect(*entries[0].Request.Destination).To(Equal("hoverfly.io"))
	Expect(*entries[0].Request.Body).To(Equal(""))

	Expect(entries[0].Response.Status).To(Equal(200))
	Expect(entries[0].Response.Body).To(Equal("test body"))
	Expect(entries[0].Response.Headers["test-header"]).To(ContainElement("one"))
	Expect(entries[0].Response.Headers["test-header"]).To(ContainElement("two"))

	Expect(entries[0].Mode).To(Equal("test-mode"))
	Expect(entries[0].TimeStarted).To(Equal(nowTime.Format(journal.RFC3339Milli)))
	Expect(entries[0].Latency).To(BeNumerically("<", 1))
}

func Test_Journal_NewEntryWithMemoryLimit_TruncateBody(t *testing.T) {
	RegisterTestingT(t)

	unit := journal.NewJournal()
	unit.BodySizeLimit = 15

	request, _ := http.NewRequest("GET", "http://hoverfly.io", io.NopCloser(bytes.NewBufferString("large request body")),)

	nowTime := time.Now()

	_, err := unit.NewEntry(request, &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString("large response body")),
		Header: http.Header{
			"test-header": []string{
				"one", "two",
			},
		},
	}, "test-mode", nowTime)
	Expect(err).To(BeNil())

	journalView, err := unit.GetEntries(0, 25, nil, nil, "")
	entries := journalView.Journal
	Expect(err).To(BeNil())

	Expect(entries).ToNot(BeNil())
	Expect(entries).To(HaveLen(1))

	Expect(*entries[0].Request.Method).To(Equal("GET"))
	Expect(*entries[0].Request.Destination).To(Equal("hoverfly.io"))
	Expect(*entries[0].Request.Body).To(Equal("large reques..."))

	Expect(entries[0].Response.Status).To(Equal(200))
	Expect(entries[0].Response.Body).To(Equal("large respon..."))
}


func Test_Journal_UpdateEntry_AddsRemotePostServeActionToJournalEntry(t *testing.T) {
	RegisterTestingT(t)

	unit := journal.NewJournal()

	request, _ := http.NewRequest("GET", "http://hoverfly.io", nil)

	nowTime := time.Now()

	id, err := unit.NewEntry(request, &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString("test body")),
		Header: http.Header{
			"test-header": []string{
				"one", "two",
			},
		},
	}, "test-mode", nowTime)
	Expect(err).To(BeNil())
	unit.UpdatePostServeActionDetailsInJournal(id, "outbound-http", "dummy-tracing-id", nowTime, nowTime, 200)

	journalView, err := unit.GetEntries(0, 25, nil, nil, "")
	entries := journalView.Journal
	expectedTime := nowTime.Format(journal.RFC3339Milli)
	Expect(err).To(BeNil())

	Expect(entries).ToNot(BeNil())
	Expect(entries).To(HaveLen(1))

	Expect(*entries[0].Request.Method).To(Equal("GET"))
	Expect(*entries[0].Request.Destination).To(Equal("hoverfly.io"))
	Expect(*entries[0].Request.Body).To(Equal(""))

	Expect(entries[0].Response.Status).To(Equal(200))
	Expect(entries[0].Response.Body).To(Equal("test body"))
	Expect(entries[0].Response.Headers["test-header"]).To(ContainElement("one"))
	Expect(entries[0].Response.Headers["test-header"]).To(ContainElement("two"))

	Expect(entries[0].Mode).To(Equal("test-mode"))
	Expect(entries[0].TimeStarted).To(Equal(expectedTime))
	Expect(entries[0].PostServeActionEntry.ActionName).To(Equal("outbound-http"))
	Expect(entries[0].PostServeActionEntry.CorrelationId).To(Equal("dummy-tracing-id"))
	Expect(entries[0].PostServeActionEntry.InvokedTime).To(Equal(expectedTime))
	Expect(entries[0].PostServeActionEntry.CompletedTime).To(Equal(expectedTime))
	Expect(entries[0].PostServeActionEntry.HttpStatus).To(Equal(200))
	Expect(entries[0].Latency).To(BeNumerically("<", 1))
}

func Test_Journal_UpdateEntry_AddsLocalPostServeActionToJournalEntry(t *testing.T) {
	RegisterTestingT(t)

	unit := journal.NewJournal()

	request, _ := http.NewRequest("GET", "http://hoverfly.io", nil)

	nowTime := time.Now()

	id, err := unit.NewEntry(request, &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString("test body")),
		Header: http.Header{
			"test-header": []string{
				"one", "two",
			},
		},
	}, "test-mode", nowTime)
	Expect(err).To(BeNil())
	unit.UpdatePostServeActionDetailsInJournal(id, "outbound-http", "", nowTime, nowTime, 0)

	journalView, err := unit.GetEntries(0, 25, nil, nil, "")
	entries := journalView.Journal
	expectedTime := nowTime.Format(journal.RFC3339Milli)
	Expect(err).To(BeNil())

	Expect(entries).ToNot(BeNil())
	Expect(entries).To(HaveLen(1))

	Expect(*entries[0].Request.Method).To(Equal("GET"))
	Expect(*entries[0].Request.Destination).To(Equal("hoverfly.io"))
	Expect(*entries[0].Request.Body).To(Equal(""))

	Expect(entries[0].Response.Status).To(Equal(200))
	Expect(entries[0].Response.Body).To(Equal("test body"))
	Expect(entries[0].Response.Headers["test-header"]).To(ContainElement("one"))
	Expect(entries[0].Response.Headers["test-header"]).To(ContainElement("two"))

	Expect(entries[0].Mode).To(Equal("test-mode"))
	Expect(entries[0].TimeStarted).To(Equal(expectedTime))
	Expect(entries[0].PostServeActionEntry.ActionName).To(Equal("outbound-http"))
	Expect(entries[0].PostServeActionEntry.CorrelationId).To(Equal(""))
	Expect(entries[0].PostServeActionEntry.InvokedTime).To(Equal(expectedTime))
	Expect(entries[0].PostServeActionEntry.CompletedTime).To(Equal(expectedTime))
	Expect(entries[0].PostServeActionEntry.HttpStatus).To(Equal(0))
	Expect(entries[0].Latency).To(BeNumerically("<", 1))
}

func Test_JournalIndex_NewEntryAfterAddingIndex_AddsJournalIndexEntryToIndexes(t *testing.T) {
	RegisterTestingT(t)

	unit := journal.NewJournal()

	request, _ := http.NewRequest("GET", "http://hoverfly.io?id=1234", nil)

	nowTime := time.Now()

	indexName := "Request.QueryParam.id"
	indexErr := unit.AddIndex(indexName)
	Expect(indexErr).To(BeNil())

	_, err := unit.NewEntry(request, &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString("test body")),
		Header: http.Header{
			"test-header": []string{
				"one", "two",
			},
		},
	}, "test-mode", nowTime)

	Expect(err).To(BeNil())
	indexes := unit.Indexes
	Expect(indexes).ToNot(BeNil())
	Expect(indexes).To(HaveLen(1))
	Expect(indexes[0].Name).To(Equal(indexName))
	Expect(indexes[0].Entries).To(HaveKey("1234"))
}

func Test_JournalIndex_NewEntryBeforeAddingIndex_AddsJournalIndexEntryToIndexes(t *testing.T) {
	RegisterTestingT(t)

	unit := journal.NewJournal()

	request, _ := http.NewRequest("GET", "http://hoverfly.io?id=1234", nil)

	nowTime := time.Now()

	_, err := unit.NewEntry(request, &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString("test body")),
		Header: http.Header{
			"test-header": []string{
				"one", "two",
			},
		},
	}, "test-mode", nowTime)

	Expect(err).To(BeNil())

	indexName := "Request.QueryParam.id"
	indexErr := unit.AddIndex(indexName)
	Expect(indexErr).To(BeNil())

	indexes := unit.Indexes
	Expect(indexes).ToNot(BeNil())
	Expect(indexes).To(HaveLen(1))
	Expect(indexes[0].Name).To(Equal(indexName))
	Expect(indexes[0].Entries).To(HaveKey("1234"))
}

func Test_DeleteJournalIndex(t *testing.T) {
	RegisterTestingT(t)

	unit := journal.NewJournal()

	indexName := "Request.QueryParam.id"
	indexErr := unit.AddIndex(indexName)
	Expect(indexErr).To(BeNil())

	unit.DeleteIndex(indexName)

	indexes := unit.Indexes
	Expect(indexes).ToNot(BeNil())
	Expect(indexes).To(HaveLen(0))
}

func Test_DeleteAndAddIndexBack(t *testing.T) {
	RegisterTestingT(t)

	unit := journal.NewJournal()

	indexName := "Request.QueryParam.id"
	indexErr := unit.AddIndex(indexName)
	Expect(indexErr).To(BeNil())

	unit.DeleteIndex(indexName)
	indexErr = unit.AddIndex(indexName)
	Expect(indexErr).To(BeNil())

	indexes := unit.Indexes
	Expect(indexes).ToNot(BeNil())
	Expect(indexes).To(HaveLen(1))
	Expect(indexes[0].Name).To(Equal(indexName))
}

func Test_Journal_NewEntry_RespectsEntryLimit(t *testing.T) {
	RegisterTestingT(t)

	unit := journal.NewJournal()
	unit.EntryLimit = 5

	request, _ := http.NewRequest("GET", "http://hoverfly.io", nil)

	for i := 1; i < 8; i++ {
		_, err := unit.NewEntry(request, &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString("test body")),
			Header: http.Header{
				"test-header": []string{
					"one", "two",
				},
			},
		}, strconv.Itoa(i), time.Now())
		Expect(err).To(BeNil())
	}

	journalView, err := unit.GetEntries(0, 25, nil, nil, "")
	entries := journalView.Journal
	Expect(err).To(BeNil())

	Expect(entries).ToNot(BeNil())
	Expect(entries).To(HaveLen(5))

	Expect(entries[0].Mode).To(Equal("3"))
	Expect(entries[1].Mode).To(Equal("4"))
	Expect(entries[2].Mode).To(Equal("5"))
	Expect(entries[3].Mode).To(Equal("6"))
	Expect(entries[4].Mode).To(Equal("7"))

}

func Test_Journal_NewEntry_KeepsOrder(t *testing.T) {
	RegisterTestingT(t)

	unit := journal.NewJournal()

	request, _ := http.NewRequest("GET", "http://hoverfly.io", nil)

	nowTime := time.Now()

	_, err := unit.NewEntry(request, &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString("test body")),
		Header: http.Header{
			"test-header": []string{
				"one", "two",
			},
		},
	}, "test-mode", nowTime)
	Expect(err).To(BeNil())

	request.Method = "DELETE"
	_, err = unit.NewEntry(request, &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString("test body")),
		Header: http.Header{
			"test-header": []string{
				"one", "two",
			},
		},
	}, "test-mode", nowTime)
	Expect(err).To(BeNil())

	journalView, err := unit.GetEntries(0, 25, nil, nil, "")
	entries := journalView.Journal
	Expect(err).To(BeNil())

	Expect(entries).ToNot(BeNil())
	Expect(entries).To(HaveLen(2))

	Expect(*entries[0].Request.Method).To(Equal("GET"))
	Expect(*entries[1].Request.Method).To(Equal("DELETE"))
}

func Test_Journal_NewEntry_WhenDisabledReturnsError(t *testing.T) {
	RegisterTestingT(t)

	unit := journal.NewJournal()
	unit.EntryLimit = 0

	request, _ := http.NewRequest("GET", "http://hoverfly.io", nil)
	_, err := unit.NewEntry(request, &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString("test body")),
		Header: http.Header{
			"test-header": []string{
				"one", "two",
			},
		},
	}, "test-mode", time.Now())

	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("Journal disabled"))
}

func Test_Journal_DeleteEntries_DeletesAllEntries(t *testing.T) {
	RegisterTestingT(t)

	unit := journal.NewJournal()

	request, _ := http.NewRequest("GET", "http://hoverfly.io", nil)

	nowTime := time.Now()

	unit.NewEntry(request, &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString("test body")),
		Header: http.Header{
			"test-header": []string{
				"one", "two",
			},
		},
	}, "test-mode", nowTime)

	err := unit.DeleteEntries()
	Expect(err).To(BeNil())

	journalView, err := unit.GetEntries(0, 25, nil, nil, "")
	Expect(err).To(BeNil())

	Expect(journalView.Journal).To(HaveLen(0))
}

func Test_Journal_DeleteEntries_WhenDisabledReturnsError(t *testing.T) {
	RegisterTestingT(t)

	unit := journal.NewJournal()
	unit.EntryLimit = 0

	err := unit.DeleteEntries()
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("Journal disabled"))
}

func Test_Journal_GetEntries_TurnsTimeDurationToMilliseconds(t *testing.T) {
	RegisterTestingT(t)

	unit := journal.NewJournal()

	request, _ := http.NewRequest("GET", "http://hoverfly.io", nil)
	_, err := unit.NewEntry(request, &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString("test body")),
		Header: http.Header{
			"test-header": []string{
				"one", "two",
			},
		},
	}, "test-mode", time.Now())

	Expect(err).To(BeNil())

	journalView, err := unit.GetEntries(0, 25, nil, nil, "")
	entries := journalView.Journal
	Expect(err).To(BeNil())
	Expect(entries).To(HaveLen(1))

	Expect(entries[0].Latency).To(BeNumerically(">", 0))
	Expect(entries[0].Latency).To(BeNumerically("<", 0.1))
}

func Test_Journal_GetEntries_WhenDisabledReturnsError(t *testing.T) {
	RegisterTestingT(t)

	unit := journal.NewJournal()
	unit.EntryLimit = 0

	_, err := unit.GetEntries(0, 25, nil, nil, "")
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("Journal disabled"))
}

func Test_Journal_GetEntries_WhenErrorShouldReturnDefaultJournalView(t *testing.T) {
	RegisterTestingT(t)

	unit := journal.NewJournal()
	unit.EntryLimit = 0

	journalView, err := unit.GetEntries(10, 35, nil, nil, "")
	Expect(err).ToNot(BeNil())
	Expect(journalView.Journal).To(BeEmpty())
	Expect(journalView.Offset).To(Equal(10))
	Expect(journalView.Limit).To(Equal(35))
	Expect(journalView.Total).To(Equal(0))
}

func Test_Journal_GetEntries_WithInvalidSortKeyReturnsError(t *testing.T) {
	RegisterTestingT(t)

	unit := journal.NewJournal()

	_, err := unit.GetEntries(0, 25, nil, nil, "time")
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("'time' is not a valid sort key, use timeStarted or latency"))
}

func Test_Journal_GetEntries_WithInvalidSortOrderReturnError(t *testing.T) {
	RegisterTestingT(t)

	unit := journal.NewJournal()

	_, err := unit.GetEntries(0, 25, nil, nil, "latency:slow")

	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("'slow' is not a valid sort order. use asc or desc"))
}

func Test_Journal_GetEntries_ReturnResultsSortedDescendinglyByTimeStarted(t *testing.T) {
	RegisterTestingT(t)

	unit := journal.NewJournal()

	response := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString("test body")),
	}

	for i := 0; i < 5; i++ {
		request, _ := http.NewRequest("GET", "http://hoverfly.io/path?id="+strconv.Itoa(i), nil)
		unit.NewEntry(request, response, "test-mode", time.Now())
	}

	journalView, err := unit.GetEntries(0, 2, nil, nil, "timeStarted:desc")
	Expect(err).To(BeNil())
	Expect(journalView.Journal).To(HaveLen(2))
	Expect(*journalView.Journal[0].Request.Query).To(Equal("id=4"))
	Expect(*journalView.Journal[1].Request.Query).To(Equal("id=3"))

	journalView, _ = unit.GetEntries(2, 2, nil, nil, "timeStarted:desc")
	Expect(journalView.Journal).To(HaveLen(2))
	Expect(*journalView.Journal[0].Request.Query).To(Equal("id=2"))
	Expect(*journalView.Journal[1].Request.Query).To(Equal("id=1"))

	journalView, _ = unit.GetEntries(4, 2, nil, nil, "timeStarted:desc")
	Expect(journalView.Journal).To(HaveLen(1))
	Expect(*journalView.Journal[0].Request.Query).To(Equal("id=0"))
}

func Test_Journal_GetEntries_ReturnResultsSortedAscendinglyByLatency(t *testing.T) {
	RegisterTestingT(t)

	unit := journal.NewJournal()

	response := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString("test body")),
	}

	for i := 0; i < 3; i++ {
		request, _ := http.NewRequest("GET", "http://hoverfly.io/path?id="+strconv.Itoa(i), nil)
		unit.NewEntry(request, response, "test-mode", time.Now())
	}

	journalView, err := unit.GetEntries(0, 5, nil, nil, "latency:asc")
	Expect(err).To(BeNil())
	Expect(journalView.Journal).To(HaveLen(3))
	Expect(journalView.Journal[0].Latency).Should(BeNumerically("<=", journalView.Journal[1].Latency))
	Expect(journalView.Journal[1].Latency).Should(BeNumerically("<=", journalView.Journal[2].Latency))
}

func Test_Journal_GetEntries_ReturnResultsSortedDescendinglyByLatency(t *testing.T) {
	RegisterTestingT(t)

	unit := journal.NewJournal()

	response := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString("test body")),
	}

	for i := 0; i < 3; i++ {
		request, _ := http.NewRequest("GET", "http://hoverfly.io/path?id="+strconv.Itoa(i), nil)
		unit.NewEntry(request, response, "test-mode", time.Now())
	}

	journalView, err := unit.GetEntries(0, 5, nil, nil, "latency:desc")
	Expect(err).To(BeNil())
	Expect(journalView.Journal).To(HaveLen(3))
	Expect(journalView.Journal[0].Latency).Should(BeNumerically(">=", journalView.Journal[1].Latency))
	Expect(journalView.Journal[1].Latency).Should(BeNumerically(">=", journalView.Journal[2].Latency))
}

func Test_Journal_GetEntries_ReturnPaginationResults(t *testing.T) {
	RegisterTestingT(t)

	unit := journal.NewJournal()

	response := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString("test body")),
	}

	for i := 0; i < 5; i++ {
		request, _ := http.NewRequest("GET", "http://hoverfly.io/path?id="+strconv.Itoa(i), nil)
		unit.NewEntry(request, response, "test-mode", time.Now())
	}

	journalView, err := unit.GetEntries(0, 2, nil, nil, "")
	Expect(err).To(BeNil())
	Expect(journalView.Journal).To(HaveLen(2))
	Expect(*journalView.Journal[0].Request.Query).To(Equal("id=0"))
	Expect(*journalView.Journal[1].Request.Query).To(Equal("id=1"))
	Expect(journalView.Limit).To(Equal(2))
	Expect(journalView.Offset).To(Equal(0))
	Expect(journalView.Total).To(Equal(5))

	journalView, _ = unit.GetEntries(2, 2, nil, nil, "")
	Expect(journalView.Journal).To(HaveLen(2))
	Expect(*journalView.Journal[0].Request.Query).To(Equal("id=2"))
	Expect(*journalView.Journal[1].Request.Query).To(Equal("id=3"))
	Expect(journalView.Limit).To(Equal(2))
	Expect(journalView.Offset).To(Equal(2))
	Expect(journalView.Total).To(Equal(5))

	journalView, _ = unit.GetEntries(4, 2, nil, nil, "")
	Expect(journalView.Journal).To(HaveLen(1))
	Expect(*journalView.Journal[0].Request.Query).To(Equal("id=4"))
	Expect(journalView.Limit).To(Equal(2))
	Expect(journalView.Offset).To(Equal(4))
	Expect(journalView.Total).To(Equal(5))

	journalView, err = unit.GetEntries(0, 2, nil, nil, "")
	Expect(err).To(BeNil())
	Expect(journalView.Journal).To(HaveLen(2))
	Expect(*journalView.Journal[0].Request.Query).To(Equal("id=0"))
	Expect(*journalView.Journal[1].Request.Query).To(Equal("id=1"))
}

func Test_Journal_GetEntries_ReturnEmptyPageIfOffsetIsLargerThanTotalElements(t *testing.T) {
	RegisterTestingT(t)

	unit := journal.NewJournal()

	response := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString("test body")),
	}

	for i := 0; i < 5; i++ {
		request, _ := http.NewRequest("GET", "http://hoverfly.io/path?id="+strconv.Itoa(i), nil)
		unit.NewEntry(request, response, "test-mode", time.Now())
	}

	journalView, err := unit.GetEntries(10, 2, nil, nil, "")
	Expect(err).To(BeNil())
	Expect(journalView.Journal).To(HaveLen(0))
}

func Test_Journal_GetEntries_ReturnRightLimitWhenEmpty(t *testing.T) {
	RegisterTestingT(t)

	unit := journal.NewJournal()
	journalView, err := unit.GetEntries(0, 100, nil, nil, "")
	Expect(err).To(BeNil())
	Expect(journalView.Journal).To(HaveLen(0))
	Expect(journalView.Limit).To(Equal(100))
}

func Test_Journal_GetEntries_FilteredByTimeWindow(t *testing.T) {
	RegisterTestingT(t)

	unit := journal.NewJournal()

	response := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString("test body")),
	}

	for i := 0; i < 5; i++ {
		request, _ := http.NewRequest("GET", "http://hoverfly.io/path?id="+strconv.Itoa(i), nil)
		unit.NewEntry(request, response, "test-mode", time.Date(2018, 2, 1, 2, 0, i, 0, time.UTC))
	}

	fromQuery := time.Date(2018, 2, 1, 2, 0, 1, 0, time.UTC)
	toQuery := time.Date(2018, 2, 1, 2, 0, 3, 0, time.UTC)

	journalView, err := unit.GetEntries(0, 25, &fromQuery, &toQuery, "")
	entries := journalView.Journal
	Expect(err).To(BeNil())
	Expect(entries).To(HaveLen(3))
	Expect(entries[0].TimeStarted).To(Equal("2018-02-01T02:00:01.000Z"))
	Expect(entries[1].TimeStarted).To(Equal("2018-02-01T02:00:02.000Z"))
	Expect(entries[2].TimeStarted).To(Equal("2018-02-01T02:00:03.000Z"))
}

func Test_Journal_GetFilteredEntries_WillReturnEmptySliceIfNoJournalFound(t *testing.T) {
	RegisterTestingT(t)

	unit := journal.NewJournal()

	request, _ := http.NewRequest("GET", "http://hoverfly.io/path/one?one=1&two=2", bytes.NewBufferString(`{"meta:{"field": "value"}}`))
	request.Header.Add("Accept", "application/json")

	unit.NewEntry(request, &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString("test body")),
	}, "test-mode", time.Now())

	// Body
	entries, err := unit.GetFilteredEntries(v2.JournalEntryFilterView{
		Request: &v2.RequestMatcherViewV5{
			Body: []v2.MatcherViewV5{
				{
					Matcher: matchers.Exact,
					Value:   `{"meta:{"field": "other-value"}}`,
				},
			},
		},
	})
	Expect(err).To(BeNil())
	Expect(entries).ToNot(BeNil())
	Expect(entries).To(HaveLen(0))
}

func Test_Journal_GetFilteredEntries_WillFilterOnRequestFields(t *testing.T) {
	RegisterTestingT(t)

	unit := journal.NewJournal()

	request, _ := http.NewRequest("GET", "http://hoverfly.io/path/one?one=1&two=2", bytes.NewBufferString(`{"meta:{"field": "value"}}`))
	request.Header.Add("Accept", "application/json")

	unit.NewEntry(request, &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString("test body")),
	}, "test-mode", time.Now())

	// Body

	Expect(unit.GetFilteredEntries(v2.JournalEntryFilterView{
		Request: &v2.RequestMatcherViewV5{
			Body: []v2.MatcherViewV5{
				{
					Matcher: matchers.Exact,
					Value:   `{"meta:{"field": "value"}}`,
				},
			},
		},
	})).To(HaveLen(1))

	Expect(unit.GetFilteredEntries(v2.JournalEntryFilterView{
		Request: &v2.RequestMatcherViewV5{
			Body: []v2.MatcherViewV5{
				{
					Matcher: matchers.Glob,
					Value:   `{"meta:{"field": "*"}}`,
				},
			},
		},
	})).To(HaveLen(1))

	Expect(unit.GetFilteredEntries(v2.JournalEntryFilterView{
		Request: &v2.RequestMatcherViewV5{
			Body: []v2.MatcherViewV5{
				{
					Matcher: matchers.Json,
					Value:   `{"meta:{"field": "value"}}`,
				},
			},
		},
	})).To(HaveLen(1))

	Expect(unit.GetFilteredEntries(v2.JournalEntryFilterView{
		Request: &v2.RequestMatcherViewV5{
			Body: []v2.MatcherViewV5{
				{
					Matcher: matchers.Exact,
					Value:   `{"meta:{"field": "other-value"}}`,
				},
			},
		},
	})).To(HaveLen(0))

	// Destination

	Expect(unit.GetFilteredEntries(v2.JournalEntryFilterView{
		Request: &v2.RequestMatcherViewV5{
			Destination: []v2.MatcherViewV5{
				{
					Matcher: matchers.Exact,
					Value:   `hoverfly.io`,
				},
			},
		},
	})).To(HaveLen(1))

	Expect(unit.GetFilteredEntries(v2.JournalEntryFilterView{
		Request: &v2.RequestMatcherViewV5{
			Destination: []v2.MatcherViewV5{
				{
					Matcher: matchers.Glob,
					Value:   "*.io",
				},
			},
		},
	})).To(HaveLen(1))

	Expect(unit.GetFilteredEntries(v2.JournalEntryFilterView{
		Request: &v2.RequestMatcherViewV5{
			Destination: []v2.MatcherViewV5{
				{
					Matcher: matchers.Exact,
					Value:   "not-hoverfly.com",
				},
			},
		},
	})).To(HaveLen(0))

	// Destination

	Expect(unit.GetFilteredEntries(v2.JournalEntryFilterView{
		Request: &v2.RequestMatcherViewV5{
			Method: []v2.MatcherViewV5{
				{
					Matcher: matchers.Exact,
					Value:   "GET",
				},
			},
		},
	})).To(HaveLen(1))

	Expect(unit.GetFilteredEntries(v2.JournalEntryFilterView{
		Request: &v2.RequestMatcherViewV5{
			Method: []v2.MatcherViewV5{
				{
					Matcher: matchers.Glob,
					Value:   "*",
				},
			},
		},
	})).To(HaveLen(1))

	Expect(unit.GetFilteredEntries(v2.JournalEntryFilterView{
		Request: &v2.RequestMatcherViewV5{
			Method: []v2.MatcherViewV5{
				{
					Matcher: matchers.Exact,
					Value:   "POST",
				},
			},
		},
	})).To(HaveLen(0))

	// Path

	Expect(unit.GetFilteredEntries(v2.JournalEntryFilterView{
		Request: &v2.RequestMatcherViewV5{
			Path: []v2.MatcherViewV5{
				{
					Matcher: matchers.Exact,
					Value:   "/path/one",
				},
			},
		},
	})).To(HaveLen(1))

	Expect(unit.GetFilteredEntries(v2.JournalEntryFilterView{
		Request: &v2.RequestMatcherViewV5{
			Path: []v2.MatcherViewV5{
				{
					Matcher: matchers.Glob,
					Value:   "/path/*",
				},
			},
		},
	})).To(HaveLen(1))

	Expect(unit.GetFilteredEntries(v2.JournalEntryFilterView{
		Request: &v2.RequestMatcherViewV5{
			Path: []v2.MatcherViewV5{
				{
					Matcher: matchers.Exact,
					Value:   "/path/two",
				},
			},
		},
	})).To(HaveLen(0))

	// Query

	Expect(unit.GetFilteredEntries(v2.JournalEntryFilterView{
		Request: &v2.RequestMatcherViewV5{
			Query: &v2.QueryMatcherViewV5{
				"one": []v2.MatcherViewV5{
					{
						Matcher: matchers.Exact,
						Value:   "1",
					},
				},
				"two": []v2.MatcherViewV5{
					{
						Matcher: matchers.Exact,
						Value:   "2",
					},
				},
			},
		},
	})).To(HaveLen(1))

	Expect(unit.GetFilteredEntries(v2.JournalEntryFilterView{
		Request: &v2.RequestMatcherViewV5{
			Query: &v2.QueryMatcherViewV5{
				"one": []v2.MatcherViewV5{
					{
						Matcher: matchers.Glob,
						Value:   "*",
					},
				},
			},
		},
	})).To(HaveLen(1))

	Expect(unit.GetFilteredEntries(v2.JournalEntryFilterView{
		Request: &v2.RequestMatcherViewV5{
			Query: &v2.QueryMatcherViewV5{
				"three": []v2.MatcherViewV5{
					{
						Matcher: matchers.Glob,
						Value:   "*",
					},
				},
			},
		},
	})).To(HaveLen(0))

	// Scheme

	Expect(unit.GetFilteredEntries(v2.JournalEntryFilterView{
		Request: &v2.RequestMatcherViewV5{
			Scheme: []v2.MatcherViewV5{
				{
					Matcher: matchers.Exact,
					Value:   "http",
				},
			},
		},
	})).To(HaveLen(1))

	Expect(unit.GetFilteredEntries(v2.JournalEntryFilterView{
		Request: &v2.RequestMatcherViewV5{
			Scheme: []v2.MatcherViewV5{
				{
					Matcher: matchers.Glob,
					Value:   "*",
				},
			},
		},
	})).To(HaveLen(1))

	Expect(unit.GetFilteredEntries(v2.JournalEntryFilterView{
		Request: &v2.RequestMatcherViewV5{
			Scheme: []v2.MatcherViewV5{
				{
					Matcher: matchers.Exact,
					Value:   "not-http",
				},
			},
		},
	})).To(HaveLen(0))

	// Headers

	Expect(unit.GetFilteredEntries(v2.JournalEntryFilterView{
		Request: &v2.RequestMatcherViewV5{
			Headers: map[string][]v2.MatcherViewV5{
				"Accept": {
					{
						Matcher: matchers.Exact,
						Value:   "application/json",
					},
				},
			},
		},
	})).To(HaveLen(1))

	Expect(unit.GetFilteredEntries(v2.JournalEntryFilterView{
		Request: &v2.RequestMatcherViewV5{
			Headers: map[string][]v2.MatcherViewV5{
				"Accept": {
					{
						Matcher: matchers.Exact,
						Value:   "application/xml",
					},
				},
			},
		},
	})).To(HaveLen(0))
}

func Test_Journal_GetFilteredEntries_WillReturnEmptyIfRequestMatcherIsEmpty(t *testing.T) {
	RegisterTestingT(t)

	unit := journal.NewJournal()

	request, _ := http.NewRequest("GET", "http://hoverfly.io/path/one?one=1&two=2", bytes.NewBufferString(`{"meta:{"field": "value"}}`))
	request.Header.Add("Accept", "application/json")

	unit.NewEntry(request, &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString("test body")),
	}, "test-mode", time.Now())

	Expect(unit.GetFilteredEntries(v2.JournalEntryFilterView{
		Request: &v2.RequestMatcherViewV5{},
	})).To(HaveLen(0))
}

func Test_Journal_GetFilteredEntries_WillFilterOnFormBodyRequest(t *testing.T) {
	RegisterTestingT(t)

	unit := journal.NewJournal()

	formData := url.Values{
		"field1": {"value1"},
		"field2": {"value2"},
	}

	request, _ := http.NewRequest("POST", "http://hoverfly.io/path/one", strings.NewReader(formData.Encode()))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	unit.NewEntry(request, &http.Response{
		StatusCode: 202,
		Body:       io.NopCloser(bytes.NewBufferString("test body")),
	}, "test-mode", time.Now())

	Expect(unit.GetFilteredEntries(v2.JournalEntryFilterView{
		Request: &v2.RequestMatcherViewV5{
			Body: []v2.MatcherViewV5{
				{
					Matcher: "form",
					Value: map[string][]models.RequestFieldMatchers{
						"field1": {
							{
								Matcher: matchers.Exact,
								Value:   "value1",
							},
						},
					},
				},
			},
		},
	})).To(HaveLen(1))

	Expect(unit.GetFilteredEntries(v2.JournalEntryFilterView{
		Request: &v2.RequestMatcherViewV5{
			Body: []v2.MatcherViewV5{
				{
					Matcher: "form",
					Value: map[string][]models.RequestFieldMatchers{
						"field1": {
							{
								Matcher: matchers.Exact,
								Value:   "value2",
							},
						},
					},
				},
			},
		},
	})).To(HaveLen(0))
}

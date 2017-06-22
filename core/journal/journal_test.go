package journal_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/SpectoLabs/hoverfly/core/journal"
	. "github.com/onsi/gomega"
)

func Test_NewJournal_ProducesAJournalWithAnEmptyArray(t *testing.T) {
	RegisterTestingT(t)

	unit := journal.NewJournal()

	entries, err := unit.GetEntries()
	Expect(err).To(BeNil())

	Expect(entries).ToNot(BeNil())
	Expect(entries).To(HaveLen(0))

	Expect(unit.EntryLimit).To(Equal(1000))
}

func Test_Journal_NewEntry_AddsJournalEntryToEntries(t *testing.T) {
	RegisterTestingT(t)

	unit := journal.NewJournal()

	request, _ := http.NewRequest("GET", "http://hoverfly.io", nil)

	nowTime := time.Now()

	err := unit.NewEntry(request, &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewBufferString("test body")),
		Header: http.Header{
			"test-header": []string{
				"one", "two",
			},
		},
	}, "test-mode", nowTime)
	Expect(err).To(BeNil())

	entries, err := unit.GetEntries()
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
	Expect(entries[0].TimeStarted).To(Equal(nowTime.Format(time.RFC3339)))
	Expect(entries[0].Latency).To(BeNumerically("<", 1))
}

func Test_Journal_NewEntry_RespectsEntryLimit(t *testing.T) {
	RegisterTestingT(t)

	unit := journal.NewJournal()
	unit.EntryLimit = 5

	request, _ := http.NewRequest("GET", "http://hoverfly.io", nil)

	for i := 1; i < 8; i++ {
		err := unit.NewEntry(request, &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewBufferString("test body")),
			Header: http.Header{
				"test-header": []string{
					"one", "two",
				},
			},
		}, strconv.Itoa(i), time.Now())
		Expect(err).To(BeNil())
	}

	entries, err := unit.GetEntries()
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

	err := unit.NewEntry(request, &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewBufferString("test body")),
		Header: http.Header{
			"test-header": []string{
				"one", "two",
			},
		},
	}, "test-mode", nowTime)
	Expect(err).To(BeNil())

	request.Method = "DELETE"
	err = unit.NewEntry(request, &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewBufferString("test body")),
		Header: http.Header{
			"test-header": []string{
				"one", "two",
			},
		},
	}, "test-mode", nowTime)
	Expect(err).To(BeNil())

	entries, err := unit.GetEntries()
	Expect(err).To(BeNil())

	Expect(entries).ToNot(BeNil())
	Expect(entries).To(HaveLen(2))

	Expect(*entries[0].Request.Method).To(Equal("GET"))
	Expect(*entries[1].Request.Method).To(Equal("DELETE"))
}

func Test_Journal_NewEntry_WhenDisabledReturnsError(t *testing.T) {
	RegisterTestingT(t)

	unit := journal.NewDisabledJournal()

	request, _ := http.NewRequest("GET", "http://hoverfly.io", nil)
	err := unit.NewEntry(request, &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewBufferString("test body")),
		Header: http.Header{
			"test-header": []string{
				"one", "two",
			},
		},
	}, "test-mode", time.Now())

	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("No journal set"))
}

func Test_Journal_DeleteEntries_DeletesAllEntries(t *testing.T) {
	RegisterTestingT(t)

	unit := journal.NewJournal()

	request, _ := http.NewRequest("GET", "http://hoverfly.io", nil)

	nowTime := time.Now()

	unit.NewEntry(request, &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewBufferString("test body")),
		Header: http.Header{
			"test-header": []string{
				"one", "two",
			},
		},
	}, "test-mode", nowTime)

	err := unit.DeleteEntries()
	Expect(err).To(BeNil())

	entries, err := unit.GetEntries()
	Expect(err).To(BeNil())

	Expect(entries).To(HaveLen(0))
}

func Test_Journal_DeleteEntries_WhenDisabledReturnsError(t *testing.T) {
	RegisterTestingT(t)

	unit := journal.NewDisabledJournal()

	err := unit.DeleteEntries()
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("No journal set"))
}

func Test_Journal_GetEntries_WhenDisabledReturnsError(t *testing.T) {
	RegisterTestingT(t)

	unit := journal.NewDisabledJournal()

	_, err := unit.GetEntries()
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("No journal set"))
}

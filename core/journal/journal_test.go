package journal_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
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
	Expect(entries[0].TimeStarted).To(Equal(nowTime))
	Expect(entries[0].Latency).To(BeNumerically("<", 1))
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

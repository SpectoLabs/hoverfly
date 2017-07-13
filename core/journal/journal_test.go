package journal_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/journal"
	"github.com/SpectoLabs/hoverfly/core/util"
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
	Expect(entries[0].TimeStarted).To(Equal(nowTime.Format(journal.RFC3339Milli)))
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

	unit := journal.NewJournal()
	unit.EntryLimit = 0

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
	Expect(err.Error()).To(Equal("Journal disabled"))
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
	err := unit.NewEntry(request, &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewBufferString("test body")),
		Header: http.Header{
			"test-header": []string{
				"one", "two",
			},
		},
	}, "test-mode", time.Now())

	Expect(err).To(BeNil())

	entries, err := unit.GetEntries()
	Expect(err).To(BeNil())
	Expect(entries).To(HaveLen(1))

	Expect(entries[0].Latency).To(BeNumerically(">", 0))
	Expect(entries[0].Latency).To(BeNumerically("<", 0.1))
}

func Test_Journal_GetEntries_WhenDisabledReturnsError(t *testing.T) {
	RegisterTestingT(t)

	unit := journal.NewJournal()
	unit.EntryLimit = 0

	_, err := unit.GetEntries()
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("Journal disabled"))
}

func Test_Journal_GetFilteredEntries_WillFilterOnRequestFields(t *testing.T) {
	RegisterTestingT(t)

	unit := journal.NewJournal()

	request, _ := http.NewRequest("GET", "http://hoverfly.io/path/one?one=1&two=2", bytes.NewBufferString(`{"meta:{"field": "value"}}`))
	request.Header.Add("Accept", "application/json")

	unit.NewEntry(request, &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewBufferString("test body")),
	}, "test-mode", time.Now())

	// Body

	Expect(unit.GetFilteredEntries(v2.JournalEntryFilterView{
		Request: &v2.RequestMatcherViewV2{
			Body: &v2.RequestFieldMatchersView{
				ExactMatch: util.StringToPointer(`{"meta:{"field": "value"}}`),
			},
		},
	})).To(HaveLen(1))

	Expect(unit.GetFilteredEntries(v2.JournalEntryFilterView{
		Request: &v2.RequestMatcherViewV2{
			Body: &v2.RequestFieldMatchersView{
				GlobMatch: util.StringToPointer(`{"meta:{"field": "*"}}`),
			},
		},
	})).To(HaveLen(1))

	Expect(unit.GetFilteredEntries(v2.JournalEntryFilterView{
		Request: &v2.RequestMatcherViewV2{
			Body: &v2.RequestFieldMatchersView{
				JsonMatch: util.StringToPointer(`{"meta:{"field": "value"}}`),
			},
		},
	})).To(HaveLen(1))

	Expect(unit.GetFilteredEntries(v2.JournalEntryFilterView{
		Request: &v2.RequestMatcherViewV2{
			Body: &v2.RequestFieldMatchersView{
				ExactMatch: util.StringToPointer(`{"meta:{"field": "other-value"}}`),
			},
		},
	})).To(HaveLen(0))

	// Destination

	Expect(unit.GetFilteredEntries(v2.JournalEntryFilterView{
		Request: &v2.RequestMatcherViewV2{
			Destination: &v2.RequestFieldMatchersView{
				ExactMatch: util.StringToPointer("hoverfly.io"),
			},
		},
	})).To(HaveLen(1))

	Expect(unit.GetFilteredEntries(v2.JournalEntryFilterView{
		Request: &v2.RequestMatcherViewV2{
			Destination: &v2.RequestFieldMatchersView{
				GlobMatch: util.StringToPointer("*.io"),
			},
		},
	})).To(HaveLen(1))

	Expect(unit.GetFilteredEntries(v2.JournalEntryFilterView{
		Request: &v2.RequestMatcherViewV2{
			Destination: &v2.RequestFieldMatchersView{
				ExactMatch: util.StringToPointer("hoverfly.com"),
			},
		},
	})).To(HaveLen(0))

	// Destination

	Expect(unit.GetFilteredEntries(v2.JournalEntryFilterView{
		Request: &v2.RequestMatcherViewV2{
			Method: &v2.RequestFieldMatchersView{
				ExactMatch: util.StringToPointer("GET"),
			},
		},
	})).To(HaveLen(1))

	Expect(unit.GetFilteredEntries(v2.JournalEntryFilterView{
		Request: &v2.RequestMatcherViewV2{
			Method: &v2.RequestFieldMatchersView{
				GlobMatch: util.StringToPointer("*"),
			},
		},
	})).To(HaveLen(1))

	Expect(unit.GetFilteredEntries(v2.JournalEntryFilterView{
		Request: &v2.RequestMatcherViewV2{
			Method: &v2.RequestFieldMatchersView{
				ExactMatch: util.StringToPointer("POST"),
			},
		},
	})).To(HaveLen(0))

	// Path

	Expect(unit.GetFilteredEntries(v2.JournalEntryFilterView{
		Request: &v2.RequestMatcherViewV2{
			Path: &v2.RequestFieldMatchersView{
				ExactMatch: util.StringToPointer("/path/one"),
			},
		},
	})).To(HaveLen(1))

	Expect(unit.GetFilteredEntries(v2.JournalEntryFilterView{
		Request: &v2.RequestMatcherViewV2{
			Path: &v2.RequestFieldMatchersView{
				GlobMatch: util.StringToPointer("/path/*"),
			},
		},
	})).To(HaveLen(1))

	Expect(unit.GetFilteredEntries(v2.JournalEntryFilterView{
		Request: &v2.RequestMatcherViewV2{
			Path: &v2.RequestFieldMatchersView{
				ExactMatch: util.StringToPointer("/path/two"),
			},
		},
	})).To(HaveLen(0))

	// Query

	Expect(unit.GetFilteredEntries(v2.JournalEntryFilterView{
		Request: &v2.RequestMatcherViewV2{
			Query: &v2.RequestFieldMatchersView{
				ExactMatch: util.StringToPointer("one=1&two=2"),
			},
		},
	})).To(HaveLen(1))

	Expect(unit.GetFilteredEntries(v2.JournalEntryFilterView{
		Request: &v2.RequestMatcherViewV2{
			Query: &v2.RequestFieldMatchersView{
				GlobMatch: util.StringToPointer("*one=1*"),
			},
		},
	})).To(HaveLen(1))

	Expect(unit.GetFilteredEntries(v2.JournalEntryFilterView{
		Request: &v2.RequestMatcherViewV2{
			Query: &v2.RequestFieldMatchersView{
				ExactMatch: util.StringToPointer("one=1"),
			},
		},
	})).To(HaveLen(0))

	// Scheme

	Expect(unit.GetFilteredEntries(v2.JournalEntryFilterView{
		Request: &v2.RequestMatcherViewV2{
			Scheme: &v2.RequestFieldMatchersView{
				ExactMatch: util.StringToPointer("http"),
			},
		},
	})).To(HaveLen(1))

	Expect(unit.GetFilteredEntries(v2.JournalEntryFilterView{
		Request: &v2.RequestMatcherViewV2{
			Scheme: &v2.RequestFieldMatchersView{
				GlobMatch: util.StringToPointer("*"),
			},
		},
	})).To(HaveLen(1))

	Expect(unit.GetFilteredEntries(v2.JournalEntryFilterView{
		Request: &v2.RequestMatcherViewV2{
			Scheme: &v2.RequestFieldMatchersView{
				ExactMatch: util.StringToPointer("nothttp"),
			},
		},
	})).To(HaveLen(0))

	// Headers

	Expect(unit.GetFilteredEntries(v2.JournalEntryFilterView{
		Request: &v2.RequestMatcherViewV2{
			Headers: map[string][]string{
				"Accept": []string{"application/json"},
			},
		},
	})).To(HaveLen(1))

	Expect(unit.GetFilteredEntries(v2.JournalEntryFilterView{
		Request: &v2.RequestMatcherViewV2{
			Headers: map[string][]string{
				"Accept": []string{"application/xml"},
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
		Body:       ioutil.NopCloser(bytes.NewBufferString("test body")),
	}, "test-mode", time.Now())

	Expect(unit.GetFilteredEntries(v2.JournalEntryFilterView{
		Request: &v2.RequestMatcherViewV2{},
	})).To(HaveLen(0))
}

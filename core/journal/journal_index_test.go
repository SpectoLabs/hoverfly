package journal

import (
	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/SpectoLabs/hoverfly/core/util"
	. "github.com/onsi/gomega"
	"testing"
)

func TestIndex_AddJournalEntry(t *testing.T) {
	RegisterTestingT(t)

	index := Index{
		Name:    "Request.QueryParam.id",
		Entries: make(map[string]*JournalEntry),
	}

	queryParams := make(map[string][]string)
	queryParams["id"] = []string{"0100"}
	journalEntry := JournalEntry{
		Id: "12345",
		Request: &models.RequestDetails{
			Query: queryParams,
		},
	}
	index.AddJournalEntry(&journalEntry)
	Expect(index.Entries).ToNot(BeNil())
	Expect(index.Entries).To(HaveLen(1))
	Expect(index.Entries["0100"]).To(Equal(&journalEntry))
}

func TestIndex_AddJournalEntryWithJsonRequestBody(t *testing.T) {
	RegisterTestingT(t)

	index := Index{
		Name:    "Request.Body 'jsonpath' '$.name'",
		Entries: make(map[string]*JournalEntry),
	}

	journalEntry := JournalEntry{
		Id: "12345",
		Request: &models.RequestDetails{
			Body: "{\"name\":\"Application Testing\"}",
		},
	}
	index.AddJournalEntry(&journalEntry)
	Expect(index.Entries).ToNot(BeNil())
	Expect(index.Entries).To(HaveLen(1))
	Expect(index.Entries["Application Testing"]).To(Equal(&journalEntry))
}

func TestIndex_AddJournalEntryWithXMLRequestBody(t *testing.T) {
	RegisterTestingT(t)

	index := Index{
		Name:    "Request.Body 'xpath' '/document/id'",
		Entries: make(map[string]*JournalEntry),
	}

	journalEntry := JournalEntry{
		Id: "12345",
		Request: &models.RequestDetails{
			Body: "<document><id>1234</id><name>Test</name></document>",
		},
	}
	index.AddJournalEntry(&journalEntry)
	Expect(index.Entries).ToNot(BeNil())
	Expect(index.Entries).To(HaveLen(1))
	Expect(index.Entries["1234"]).To(Equal(&journalEntry))
}

func TestIndex_ConvertToIndexView(t *testing.T) {
	RegisterTestingT(t)

	journalEntry1 := JournalEntry{
		Id: "1",
	}
	journalEntry2 := JournalEntry{
		Id: "2",
	}

	entries := make(map[string]*JournalEntry)
	entries["a"] = &journalEntry1
	entries["b"] = &journalEntry2
	index := Index{
		Name:    "Request.QueryParam.id",
		Entries: entries,
	}
	indexView := index.getIndexView()
	Expect(indexView).ToNot(BeNil())
	Expect(indexView.Name).To(Equal("Request.QueryParam.id"))
	Expect(indexView.Entries).To(HaveLen(2))
	Expect(indexView.Entries[0].Key).To(Equal("a"))
	Expect(indexView.Entries[0].JournalEntryId).To(Equal("1"))
	Expect(indexView.Entries[1].Key).To(Equal("b"))
	Expect(indexView.Entries[1].JournalEntryId).To(Equal("2"))
}

func TestIndex_FilteredConversionIndexView(t *testing.T) {
	RegisterTestingT(t)
	journalEntry1 := JournalEntry{
		Id: "1",
	}
	journalEntry2 := JournalEntry{
		Id: "2",
	}

	entries := make(map[string]*JournalEntry)
	entries["a"] = &journalEntry1
	entries["b"] = &journalEntry2
	index := Index{
		Name:    "Request.QueryParam.id",
		Entries: entries,
	}
	selectedJournalEntryIds := util.NewHashSet()
	selectedJournalEntryIds.Add("2")
	indexView := index.convertIndex(selectedJournalEntryIds)
	Expect(indexView).ToNot(BeNil())
	Expect(indexView.Name).To(Equal("Request.QueryParam.id"))
	Expect(indexView.Entries).To(HaveLen(1))
	Expect(indexView.Entries[0].Key).To(Equal("b"))
	Expect(indexView.Entries[0].JournalEntryId).To(Equal("2"))

}

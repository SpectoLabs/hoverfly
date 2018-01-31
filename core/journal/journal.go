package journal

import (
	"fmt"
	"net/http"
	"time"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/matching"
	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/SpectoLabs/hoverfly/core/util"
)

var RFC3339Milli = "2006-01-02T15:04:05.000Z07:00"

type JournalEntry struct {
	Request     *models.RequestDetails
	Response    *models.ResponseDetails
	Mode        string
	TimeStarted time.Time
	Latency     time.Duration
}

type Journal struct {
	entries    []JournalEntry
	EntryLimit int
}

func NewJournal() *Journal {
	return &Journal{
		entries:    []JournalEntry{},
		EntryLimit: 1000,
	}
}

func (this *Journal) NewEntry(request *http.Request, response *http.Response, mode string, started time.Time) error {
	if this.EntryLimit == 0 {
		return fmt.Errorf("Journal disabled")
	}

	payloadRequest, _ := models.NewRequestDetailsFromHttpRequest(request)

	respBody, _ := util.GetResponseBody(response)

	payloadResponse := &models.ResponseDetails{
		Status:  response.StatusCode,
		Body:    string(respBody),
		Headers: response.Header,
	}

	if len(this.entries) >= this.EntryLimit {
		this.entries = append(this.entries[:0], this.entries[1:]...)
	}

	this.entries = append(this.entries, JournalEntry{
		Request:     &payloadRequest,
		Response:    payloadResponse,
		Mode:        mode,
		TimeStarted: started,
		Latency:     time.Since(started),
	})

	return nil
}

func (this Journal) GetEntries(offset int, limit int) (v2.JournalView, error) {
	journalView := v2.JournalView{
		Journal: []v2.JournalEntryView{},
		Offset: 0,
		Limit: v2.DefaultJournalLimit,
		Total: 0,
	}

	if this.EntryLimit == 0 {
		return journalView, fmt.Errorf("Journal disabled")
	}

	totalElements := len(this.entries)

	if offset < 0 {
		offset = 0
	} else if offset >= totalElements {
		return journalView, nil
	}


	endIndex := offset + limit
	if endIndex > totalElements {
		endIndex = totalElements
	}

	journalView.Journal = convertJournalEntries(this.entries[offset:endIndex])
	journalView.Offset = offset
	journalView.Limit = limit
	journalView.Total = totalElements
	return journalView, nil
}

func (this Journal) GetFilteredEntries(journalEntryFilterView v2.JournalEntryFilterView) ([]v2.JournalEntryView, error) {
	filteredEntries := []v2.JournalEntryView{}
	if this.EntryLimit == 0 {
		return filteredEntries, fmt.Errorf("Journal disabled")
	}

	requestMatcher := models.RequestMatcher{
		Path:        models.NewRequestFieldMatchersFromView(journalEntryFilterView.Request.Path),
		Method:      models.NewRequestFieldMatchersFromView(journalEntryFilterView.Request.Method),
		Destination: models.NewRequestFieldMatchersFromView(journalEntryFilterView.Request.Destination),
		Scheme:      models.NewRequestFieldMatchersFromView(journalEntryFilterView.Request.Scheme),
		Query:       models.NewRequestFieldMatchersFromView(journalEntryFilterView.Request.Query),
		Body:        models.NewRequestFieldMatchersFromView(journalEntryFilterView.Request.Body),
		Headers:     journalEntryFilterView.Request.Headers,

	}

	allEntries := convertJournalEntries(this.entries)

	for _, entry := range allEntries {
		if requestMatcher.Body == nil && requestMatcher.Destination == nil &&
			requestMatcher.Headers == nil && requestMatcher.Method == nil &&
			requestMatcher.Path == nil && requestMatcher.Query == nil &&
			requestMatcher.Scheme == nil {
			continue
		}
		if !matching.UnscoredFieldMatcher(requestMatcher.Body, *entry.Request.Body).Matched {
			continue
		}
		if !matching.UnscoredFieldMatcher(requestMatcher.Destination, *entry.Request.Destination).Matched {
			continue
		}
		if !matching.UnscoredFieldMatcher(requestMatcher.Method, *entry.Request.Method).Matched {
			continue
		}
		if !matching.UnscoredFieldMatcher(requestMatcher.Path, *entry.Request.Path).Matched {
			continue
		}
		if !matching.UnscoredFieldMatcher(requestMatcher.Query, *entry.Request.Query).Matched {
			continue
		}
		if !matching.UnscoredFieldMatcher(requestMatcher.Scheme, *entry.Request.Scheme).Matched {
			continue
		}
		if !matching.CountlessHeaderMatcher(requestMatcher.Headers, entry.Request.Headers).Matched {
			continue
		}
		filteredEntries = append(filteredEntries, entry)
	}

	return filteredEntries, nil
}

func (this *Journal) DeleteEntries() error {
	if this.EntryLimit == 0 {
		return fmt.Errorf("Journal disabled")
	}

	this.entries = []JournalEntry{}

	return nil
}

func convertJournalEntries(entries []JournalEntry) []v2.JournalEntryView {

	journalEntryViews := []v2.JournalEntryView{}

	for _, journalEntry := range entries {
		journalEntryViews = append(journalEntryViews, v2.JournalEntryView{
			Request:     journalEntry.Request.ConvertToRequestDetailsView(),
			Response:    journalEntry.Response.ConvertToResponseDetailsView(),
			Mode:        journalEntry.Mode,
			TimeStarted: journalEntry.TimeStarted.Format(RFC3339Milli),
			Latency:     journalEntry.Latency.Seconds() * 1e3,
		})
	}

	return journalEntryViews
}

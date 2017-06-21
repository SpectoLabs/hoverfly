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

func (this Journal) GetEntries() ([]v2.JournalEntryView, error) {
	if this.EntryLimit == 0 {
		return []v2.JournalEntryView{}, fmt.Errorf("Journal disabled")
	}

	journalEntryViews := []v2.JournalEntryView{}
	for _, journalEntry := range this.entries {
		journalEntryViews = append(journalEntryViews, v2.JournalEntryView{
			Request:     journalEntry.Request.ConvertToRequestDetailsView(),
			Response:    journalEntry.Response.ConvertToResponseDetailsView(),
			Mode:        journalEntry.Mode,
			TimeStarted: journalEntry.TimeStarted.Format(time.RFC3339),
			Latency:     (journalEntry.Latency / time.Millisecond),
		})
	}
	return journalEntryViews, nil
}

func (this Journal) GetFilteredEntries(journalEntryFilterView v2.JournalEntryFilterView) ([]v2.JournalEntryView, error) {
	requestMatcher := models.RequestMatcher{
		Path:        models.NewRequestFieldMatchersFromView(journalEntryFilterView.Request.Path),
		Method:      models.NewRequestFieldMatchersFromView(journalEntryFilterView.Request.Method),
		Destination: models.NewRequestFieldMatchersFromView(journalEntryFilterView.Request.Destination),
		Scheme:      models.NewRequestFieldMatchersFromView(journalEntryFilterView.Request.Scheme),
		Query:       models.NewRequestFieldMatchersFromView(journalEntryFilterView.Request.Query),
		Body:        models.NewRequestFieldMatchersFromView(journalEntryFilterView.Request.Body),
		Headers:     journalEntryFilterView.Request.Headers,
	}
	allEntries, err := this.GetEntries()
	if err != nil {
		return []v2.JournalEntryView{}, err
	}

	filteredEntries := []v2.JournalEntryView{}

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

package journal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	sorting "sort"
	"strings"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/matching"
	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/SpectoLabs/hoverfly/core/util"
	log "github.com/sirupsen/logrus"
)

const RFC3339Milli = "2006-01-02T15:04:05.000Z07:00"

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
	mutex      sync.Mutex
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
		Body:    respBody,
		Headers: response.Header,
	}

	this.mutex.Lock()
	if len(this.entries) >= this.EntryLimit {
		this.entries = append(this.entries[:0], this.entries[1:]...)
	}

	entry := JournalEntry{
		Request:     &payloadRequest,
		Response:    payloadResponse,
		Mode:        mode,
		TimeStarted: started,
		Latency:     time.Since(started),
	}

	this.entries = append(this.entries, entry)

	if log.IsLevelEnabled(log.DebugLevel) {
		buf := new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		// do not escape characters in HTML like < and >
		enc.SetEscapeHTML(false)
		err := enc.Encode(convertJournalEntry(entry))
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Error("invalid journal entry")
		} else {
			log.WithFields(log.Fields{
				"json": buf.String(),
			}).Debug("journal entry")
		}
	}

	this.mutex.Unlock()

	return nil
}

func (this *Journal) GetEntries(offset int, limit int, from *time.Time, to *time.Time, sort string) (v2.JournalView, error) {

	journalView := v2.JournalView{
		Journal: []v2.JournalEntryView{},
		Offset:  offset,
		Limit:   limit,
		Total:   0,
	}

	if this.EntryLimit == 0 {
		return journalView, fmt.Errorf("Journal disabled")
	}

	sortKey, sortOrder, err := getSortParameters(sort)

	if err != nil {
		return journalView, err
	}

	var selectedEntries []JournalEntry

	// Filtering
	if from != nil || to != nil {
		for _, entry := range this.entries {
			if from != nil && entry.TimeStarted.Before(*from) {
				continue
			}
			if to != nil && entry.TimeStarted.After(*to) {
				continue
			}
			selectedEntries = append(selectedEntries, entry)
		}
	} else {
		selectedEntries = append(selectedEntries, this.entries...)
	}

	// Sorting
	if sortKey == "timestarted" && sortOrder == "desc" {
		sorting.Slice(selectedEntries, func(i, j int) bool {
			return selectedEntries[i].TimeStarted.After(selectedEntries[j].TimeStarted)
		})
	} else if sortKey == "latency" {
		sorting.Slice(selectedEntries, func(i, j int) bool {
			if sortOrder == "desc" {
				return selectedEntries[i].Latency > selectedEntries[j].Latency
			} else {
				return selectedEntries[i].Latency < selectedEntries[j].Latency
			}
		})
	}

	totalElements := len(selectedEntries)

	if offset >= totalElements {
		return journalView, nil
	}

	endIndex := offset + limit
	if endIndex > totalElements {
		endIndex = totalElements
	}

	journalView.Journal = convertJournalEntries(selectedEntries[offset:endIndex])
	journalView.Total = totalElements
	return journalView, nil
}

func (this *Journal) GetFilteredEntries(journalEntryFilterView v2.JournalEntryFilterView) ([]v2.JournalEntryView, error) {
	// init an empty slice to prevent serializing to a null value
	filteredEntries := []v2.JournalEntryView{}
	if this.EntryLimit == 0 {

		return filteredEntries, fmt.Errorf("Journal disabled")
	}

	requestMatcher := models.RequestMatcher{
		Path:            models.NewRequestFieldMatchersFromView(journalEntryFilterView.Request.Path),
		Method:          models.NewRequestFieldMatchersFromView(journalEntryFilterView.Request.Method),
		Destination:     models.NewRequestFieldMatchersFromView(journalEntryFilterView.Request.Destination),
		Scheme:          models.NewRequestFieldMatchersFromView(journalEntryFilterView.Request.Scheme),
		DeprecatedQuery: models.NewRequestFieldMatchersFromView(journalEntryFilterView.Request.DeprecatedQuery),
		Body:            models.NewRequestFieldMatchersFromView(journalEntryFilterView.Request.Body),
		Query:           models.NewQueryRequestFieldMatchersFromMapView(journalEntryFilterView.Request.Query),
		Headers:         models.NewRequestFieldMatchersFromMapView(journalEntryFilterView.Request.Headers),
	}

	allEntries := convertJournalEntries(this.entries)

	for _, entry := range allEntries {
		if requestMatcher.Body == nil && requestMatcher.Destination == nil &&
			requestMatcher.Headers == nil && requestMatcher.Method == nil &&
			requestMatcher.Path == nil && requestMatcher.DeprecatedQuery == nil &&
			requestMatcher.Scheme == nil && requestMatcher.Query == nil {
			continue
		}
		if !matching.FieldMatcher(requestMatcher.Body, *entry.Request.Body).Matched {
			continue
		}
		if !matching.FieldMatcher(requestMatcher.Destination, *entry.Request.Destination).Matched {
			continue
		}
		if !matching.FieldMatcher(requestMatcher.Method, *entry.Request.Method).Matched {
			continue
		}
		if !matching.FieldMatcher(requestMatcher.Path, *entry.Request.Path).Matched {
			continue
		}
		if !matching.FieldMatcher(requestMatcher.DeprecatedQuery, *entry.Request.Query).Matched {
			continue
		}
		if !matching.FieldMatcher(requestMatcher.Scheme, *entry.Request.Scheme).Matched {
			continue
		}
		if !matching.QueryMatching(requestMatcher, entry.Request.QueryMap).Matched {
			continue
		}
		if !matching.HeaderMatching(requestMatcher, entry.Request.Headers).Matched {
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

	var journalEntryViews []v2.JournalEntryView

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

func convertJournalEntry(entry JournalEntry) v2.JournalEntryView {

	return v2.JournalEntryView{
		Request:     entry.Request.ConvertToRequestDetailsView(),
		Response:    entry.Response.ConvertToResponseDetailsView(),
		Mode:        entry.Mode,
		TimeStarted: entry.TimeStarted.Format(RFC3339Milli),
		Latency:     entry.Latency.Seconds() * 1e3,
	}
}

func getSortParameters(sort string) (string, string, error) {
	sortParams := strings.Split(sort, ":")

	sortKey := strings.ToLower(sortParams[0])
	if sortKey == "" {
		sortKey = "timestarted"
	}
	sortOrder := "asc"

	if sortKey != "timestarted" && sortKey != "latency" {
		return sortKey, sortOrder, fmt.Errorf("'%s' is not a valid sort key, use timeStarted or latency", sortKey)
	}

	if len(sortParams) > 1 {
		sortOrder = strings.ToLower(sortParams[1])

		if sortOrder != "asc" && sortOrder != "desc" {
			return sortKey, sortOrder, fmt.Errorf("'%s' is not a valid sort order. use asc or desc", sortOrder)
		}
	}

	return sortKey, sortOrder, nil
}

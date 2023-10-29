package journal

import (
	v2 "github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/SpectoLabs/hoverfly/core/util"
	"github.com/aymerick/raymond"
	"strings"
)

type Index struct {
	Name     string
	template string
	Entries  map[string]*JournalEntry
}

type Request struct {
	QueryParam map[string][]string
	Header     map[string][]string
	Path       []string
	Scheme     string
	FormData   map[string][]string
	Body       func(queryType, query string, options *raymond.Options) interface{}
	BodyStr    string
	Method     string
}

func (index Index) AddJournalEntry(entry *JournalEntry) {

	indexKey := parseIndexKey(index.template, entry.Request)
	//it will be same in the condition when raymond(i.e. templating library) unable to parse it
	if indexKey == "" || indexKey != index.template {
		index.Entries[indexKey] = entry
	}
}

func parseIndexKey(source string, details *models.RequestDetails) string {

	if tpl, err := raymond.Parse("{{" + source + "}}"); err == nil {
		ctx := make(map[string]Request)
		ctx["Request"] = getRequest(details)
		if parsedValue, execErr := tpl.Exec(ctx); execErr == nil {
			return parsedValue
		}
	}
	return source
}

func getRequest(requestDetails *models.RequestDetails) Request {
	return Request{
		Path:       strings.Split(requestDetails.Path, "/")[1:],
		QueryParam: requestDetails.Query,
		Header:     requestDetails.Headers,
		Scheme:     requestDetails.Scheme,
		FormData:   requestDetails.FormData,
		BodyStr:    requestDetails.Body,
		Body:       requestBody,
		Method:     requestDetails.Method,
	}
}

func (index Index) convertIndex(filteredJournalEntryIds util.HashSet) v2.JournalIndexView {
	var journalIndexEntries []v2.JournalIndexEntryView
	for key, journalEntry := range index.Entries {

		if filteredJournalEntryIds.Contains(journalEntry.Id) {
			journalIndexEntry := v2.JournalIndexEntryView{
				Key:            key,
				JournalEntryId: journalEntry.Id,
			}
			journalIndexEntries = append(journalIndexEntries, journalIndexEntry)
		}
	}
	return v2.JournalIndexView{
		Name:    index.Name,
		Entries: journalIndexEntries,
	}
}

func (index Index) getIndexView() v2.JournalIndexView {
	var journalIndexEntries []v2.JournalIndexEntryView
	for key, journalEntry := range index.Entries {
		journalIndexEntry := v2.JournalIndexEntryView{
			Key:            key,
			JournalEntryId: journalEntry.Id,
		}
		journalIndexEntries = append(journalIndexEntries, journalIndexEntry)
	}
	return v2.JournalIndexView{
		Name:    index.Name,
		Entries: journalIndexEntries,
	}
}

func requestBody(queryType, query string, options *raymond.Options) interface{} {
	toMatch := options.Value("Request").(Request).BodyStr
	queryType = strings.ToLower(queryType)
	return util.FetchFromRequestBody(queryType, query, toMatch)
}

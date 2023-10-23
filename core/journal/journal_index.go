package journal

import (
	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/aymerick/raymond"
	"strings"
)

type Index struct {
	Name    string
	Entries map[string]*JournalEntry
}

type Request struct {
	QueryParam map[string][]string
	Header     map[string][]string
	Path       []string
	Scheme     string
	FormData   map[string][]string
	body       string
	Method     string
}

func (index Index) AddJournalEntry(entry *JournalEntry) {

	indexKey := parseIndexKey(index.Name, entry.Request)
	//it will be same in the condition when raymond(i.e. templating library) unable to parse it
	if indexKey != index.Name {
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
		body:       requestDetails.Body,
		Method:     requestDetails.Method,
	}
}

package matching

import (
	"github.com/SpectoLabs/hoverfly/core/models"
	"errors"
	"net/http"
	"strings"
)

type RequestTemplateStore struct {
	RequestTemplates	[]Entry
}

type Entry struct {
	Payload	*models.Payload
	RequestTemplate RequestTemplate
}

type RequestTemplate struct {
	Headers     map[string][]string `json:"headers"`
}

func(this *RequestTemplateStore) GetPayload(req *http.Request) (*models.Payload, error) {
	// iterate through the request templates, looking for template to match request
	for _, entry := range this.RequestTemplates {
		//TODO: not matching by default on URL and body - need to enable this
		// TODO: need to enable regex matches
		if headerMatch(entry.RequestTemplate.Headers, req.Header) {
			// return the first template to match
			return entry.Payload, nil
		}
	}
	return nil, errors.New("No match found")
}

/**
Check keys and corresponding values in template headers are also present in request headers
 */
func headerMatch(templHeaders map[string][]string, reqHeaders http.Header) (bool) {
	for headerName, headerVal := range templHeaders {
		// TODO: why is payload storing a list of strings but http has a single string??
		if (strings.Join(headerVal[:],",") == reqHeaders.Get(headerName)) {
			continue;
		} else {
			return false;
		}
	}
	return true;
}
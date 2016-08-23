package matching

import (
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/SpectoLabs/hoverfly/core/views"
	"reflect"
)

type RequestTemplateStore []RequestTemplateResponsePair

type RequestTemplateResponsePair struct {
	RequestTemplate RequestTemplate        `json:"requestTemplate"`
	Response        models.ResponseDetails `json:"response"`
}

type RequestTemplateResponsePairView struct {
	RequestTemplate RequestTemplate           `json:"requestTemplate"`
	Response        views.ResponseDetailsView `json:"response"`
}

type RequestTemplateResponsePairPayload struct {
	Data *[]RequestTemplateResponsePairView `json:"data"`
}

type RequestTemplate struct {
	Path        *string             `json:"path"`
	Method      *string             `json:"method"`
	Destination *string             `json:"destination"`
	Scheme      *string             `json:"scheme"`
	Query       *string             `json:"query"`
	Body        *string             `json:"body"`
	Headers     map[string][]string `json:"headers"`
}

func (this *RequestTemplateStore) GetResponse(req models.RequestDetails, webserver bool) (*models.ResponseDetails, error) {
	// iterate through the request templates, looking for template to match request
	for _, entry := range *this {
		// TODO: not matching by default on URL and body - need to enable this
		// TODO: need to enable regex matches
		// TODO: enable matching on scheme

		if entry.RequestTemplate.Body != nil && *entry.RequestTemplate.Body != req.Body {
			continue
		}

		if !webserver {
			if entry.RequestTemplate.Destination != nil && *entry.RequestTemplate.Destination != req.Destination {
				continue
			}
		}
		if entry.RequestTemplate.Path != nil && *entry.RequestTemplate.Path != req.Path {
			continue
		}
		if entry.RequestTemplate.Query != nil && *entry.RequestTemplate.Query != req.Query {
			continue
		}
		if !headerMatch(entry.RequestTemplate.Headers, req.Headers) {
			continue
		}
		if entry.RequestTemplate.Method != nil && *entry.RequestTemplate.Method != req.Method {
			continue
		}

		// return the first template to match
		return &entry.Response, nil
	}
	return nil, errors.New("No match found")
}

// ImportPayloads - a function to save given payloads into the database.
func (this *RequestTemplateStore) ImportPayloads(pairPayload RequestTemplateResponsePairPayload) error {
	if len(*pairPayload.Data) > 0 {
		// Convert PayloadView back to Payload for internal storage
		templateStore := pairPayload.ConvertToRequestTemplateStore()
		for _, pl := range templateStore {

			//TODO: add hooks for concsistency with request import
			// note that importing hoverfly is a disallowed circular import

			*this = append(*this, pl)
		}
		log.WithFields(log.Fields{
			"total": len(*this),
		}).Info("payloads imported")
		return nil
	}
	return fmt.Errorf("Bad request. Nothing to import!")
}

func (this *RequestTemplateStore) Wipe() {
	// don't change the pointer here!
	*this = RequestTemplateStore{}
}

/**
Check keys and corresponding values in template headers are also present in request headers
*/
func headerMatch(tmplHeaders, reqHeaders map[string][]string) bool {

	for headerName, headerVal := range tmplHeaders {
		// TODO: case insensitive lookup
		// TODO: is order of values in slice really important?

		reqHeaderVal, ok := reqHeaders[headerName]
		if ok && reflect.DeepEqual(headerVal, reqHeaderVal) {
			continue
		} else {
			return false
		}
	}
	return true
}

func (this *RequestTemplateStore) GetPayload() RequestTemplateResponsePairPayload {
	var pairsPayload []RequestTemplateResponsePairView
	for _, pair := range *this {
		pairsPayload = append(pairsPayload, pair.ConvertToRequestTemplateResponsePairView())
	}
	return RequestTemplateResponsePairPayload{
		Data: &pairsPayload,
	}
}

func (this *RequestTemplateResponsePair) ConvertToRequestTemplateResponsePairView() RequestTemplateResponsePairView {
	return RequestTemplateResponsePairView{
		RequestTemplate: this.RequestTemplate,
		Response:        this.Response.ConvertToResponseDetailsView(),
	}
}

func (this *RequestTemplateResponsePair) ConvertToRequestResponsePairView() views.RequestResponsePairView {

	getValue := func (value *string) string {
		if value == nil {
			return ""
		}
		return *value
	}

	return views.RequestResponsePairView{
		Request: views.RequestDetailsView{
			RequestType: "template",
			Path:        getValue(this.RequestTemplate.Path),
			Method:      getValue(this.RequestTemplate.Method),
			Destination: getValue(this.RequestTemplate.Destination),
			Scheme:      getValue(this.RequestTemplate.Scheme),
			Query:       getValue(this.RequestTemplate.Query),
			Body:        getValue(this.RequestTemplate.Body),
			Headers:     this.RequestTemplate.Headers,
		},
		Response: this.Response.ConvertToResponseDetailsView(),
	}
}

func (this *RequestTemplateResponsePairPayload) ConvertToRequestTemplateStore() RequestTemplateStore {
	var requestTemplateStore RequestTemplateStore
	for _, pair := range *this.Data {
		requestTemplateStore = append(requestTemplateStore, pair.ConvertToRequestTemplateResponsePair())
	}
	return requestTemplateStore
}

func (this *RequestTemplateResponsePairView) ConvertToRequestTemplateResponsePair() RequestTemplateResponsePair {
	return RequestTemplateResponsePair{
		RequestTemplate: this.RequestTemplate,
		Response:        models.NewResponseDetialsFromResponseDetailsView(this.Response),
	}
}

package matching

import (
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/core/handlers/v1"
	"github.com/SpectoLabs/hoverfly/core/models"
	. "github.com/SpectoLabs/hoverfly/core/util"
	"github.com/SpectoLabs/hoverfly/core/views"
	"github.com/ryanuber/go-glob"
	"strings"
)

type RequestTemplateStore []RequestTemplateResponsePair

type RequestTemplateResponsePair struct {
	RequestTemplate RequestTemplate        `json:"requestTemplate"`
	Response        models.ResponseDetails `json:"response"`
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

		if entry.RequestTemplate.Body != nil && !glob.Glob(*entry.RequestTemplate.Body, req.Body) {
			continue
		}

		if !webserver {
			if entry.RequestTemplate.Destination != nil && !glob.Glob(*entry.RequestTemplate.Destination, req.Destination) {
				continue
			}
		}
		if entry.RequestTemplate.Path != nil && !glob.Glob(*entry.RequestTemplate.Path, req.Path) {
			continue
		}
		if entry.RequestTemplate.Query != nil && !glob.Glob(*entry.RequestTemplate.Query, req.Query) {
			continue
		}
		if !headerMatch(entry.RequestTemplate.Headers, req.Headers) {
			continue
		}
		if entry.RequestTemplate.Method != nil && !glob.Glob(*entry.RequestTemplate.Method, req.Method) {
			continue
		}

		// return the first template to match
		return &entry.Response, nil
	}
	return nil, errors.New("No match found")
}

// ImportPayloads - a function to save given payloads into the database.
func (this *RequestTemplateStore) ImportPayloads(pairPayload v1.RequestTemplateResponsePairPayload) error {
	if len(*pairPayload.Data) > 0 {
		// Convert PayloadView back to Payload for internal storage
		templateStore := ConvertPayloadToRequestTemplateStore(pairPayload)
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
func headerMatch(templateHeaders, requestHeaders map[string][]string) bool {

	for templateHeaderKey, templateHeaderValues := range templateHeaders {
		for requestHeaderKey, requestHeaderValues := range requestHeaders {
			delete(requestHeaders, requestHeaderKey)
			requestHeaders[strings.ToLower(requestHeaderKey)] = requestHeaderValues

		}

		requestTemplateValues, templateHeaderMatched := requestHeaders[strings.ToLower(templateHeaderKey)]
		if !templateHeaderMatched {
			return false
		}

		for _, templateHeaderValue := range templateHeaderValues {
			templateValueMatched := false
			for _, requestHeaderValue := range requestTemplateValues {
				if glob.Glob(strings.ToLower(templateHeaderValue), strings.ToLower(requestHeaderValue)) {
					templateValueMatched = true
				}
			}

			if !templateValueMatched {
				return false
			}
		}
	}
	return true
}

func (this RequestTemplateStore) GetPayload() v1.RequestTemplateResponsePairPayload {
	var pairsPayload []v1.RequestTemplateResponsePairView
	for _, pair := range this {
		pairsPayload = append(pairsPayload, pair.ConvertToRequestTemplateResponsePairView())
	}
	return v1.RequestTemplateResponsePairPayload{
		Data: &pairsPayload,
	}
}

func (this *RequestTemplateResponsePair) ConvertToRequestTemplateResponsePairView() v1.RequestTemplateResponsePairView {
	return v1.RequestTemplateResponsePairView{
		RequestTemplate: v1.RequestTemplateView{
			Path:        this.RequestTemplate.Path,
			Method:      this.RequestTemplate.Method,
			Destination: this.RequestTemplate.Destination,
			Scheme:      this.RequestTemplate.Scheme,
			Query:       this.RequestTemplate.Query,
			Body:        this.RequestTemplate.Body,
			Headers:     this.RequestTemplate.Headers,
		},
		Response: this.Response.ConvertToResponseDetailsView(),
	}
}

func (this *RequestTemplateResponsePair) ConvertToRequestResponsePairView() views.RequestResponsePairView {

	return views.RequestResponsePairView{
		Request: views.RequestDetailsView{
			RequestType: StringToPointer("template"),
			Path:        this.RequestTemplate.Path,
			Method:      this.RequestTemplate.Method,
			Destination: this.RequestTemplate.Destination,
			Scheme:      this.RequestTemplate.Scheme,
			Query:       this.RequestTemplate.Query,
			Body:        this.RequestTemplate.Body,
			Headers:     this.RequestTemplate.Headers,
		},
		Response: this.Response.ConvertToResponseDetailsView(),
	}
}

func ConvertPayloadToRequestTemplateStore(payload v1.RequestTemplateResponsePairPayload) RequestTemplateStore {
	var requestTemplateStore RequestTemplateStore
	for _, pair := range *payload.Data {
		requestTemplateStore = append(requestTemplateStore, ConvertToRequestTemplateResponsePair(pair))
	}
	return requestTemplateStore
}

func ConvertToRequestTemplateResponsePair(pairView v1.RequestTemplateResponsePairView) RequestTemplateResponsePair {
	return RequestTemplateResponsePair{
		RequestTemplate: RequestTemplate{
			Path:        pairView.RequestTemplate.Path,
			Method:      pairView.RequestTemplate.Method,
			Destination: pairView.RequestTemplate.Destination,
			Scheme:      pairView.RequestTemplate.Scheme,
			Query:       pairView.RequestTemplate.Query,
			Body:        pairView.RequestTemplate.Body,
			Headers:     pairView.RequestTemplate.Headers,
		},
		Response: models.NewResponseDetailsFromResponseDetailsView(pairView.Response),
	}
}

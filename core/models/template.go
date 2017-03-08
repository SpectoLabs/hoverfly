package models

import (
	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/util"
)

type RequestFieldMatchers struct {
	ExactMatch *string
	XpathMatch *string
}

func NewRequestFieldMatchersFromView(matchers *v2.RequestFieldMatchersView) *RequestFieldMatchers {
	if matchers == nil {
		return nil
	}

	return &RequestFieldMatchers{
		ExactMatch: matchers.ExactMatch,
	}
}

func (this RequestFieldMatchers) BuildView() *v2.RequestFieldMatchersView {
	return &v2.RequestFieldMatchersView{
		ExactMatch: this.ExactMatch,
	}
}

type RequestTemplateResponsePair struct {
	RequestTemplate RequestTemplate
	Response        ResponseDetails
}

func NewRequestTemplateResponsePairFromView(view *v2.RequestResponsePairViewV2) *RequestTemplateResponsePair {
	if view.Request.Query != nil && view.Request.Query.ExactMatch != nil {
		sortedQuery := util.SortQueryString(*view.Request.Query.ExactMatch)
		view.Request.Query.ExactMatch = &sortedQuery
	}

	return &RequestTemplateResponsePair{
		RequestTemplate: RequestTemplate{
			Path:        NewRequestFieldMatchersFromView(view.Request.Path),
			Method:      NewRequestFieldMatchersFromView(view.Request.Method),
			Destination: NewRequestFieldMatchersFromView(view.Request.Destination),
			Scheme:      NewRequestFieldMatchersFromView(view.Request.Scheme),
			Query:       NewRequestFieldMatchersFromView(view.Request.Query),
			Body:        NewRequestFieldMatchersFromView(view.Request.Body),
			Headers:     view.Request.Headers,
		},
		Response: NewResponseDetailsFromResponse(view.Response),
	}
}

func (this *RequestTemplateResponsePair) BuildView() v2.RequestResponsePairViewV2 {

	var path, method, destination, scheme, query, body *v2.RequestFieldMatchersView

	if this.RequestTemplate.Path != nil {
		path = this.RequestTemplate.Path.BuildView()
	}

	if this.RequestTemplate.Method != nil {
		method = this.RequestTemplate.Method.BuildView()
	}

	if this.RequestTemplate.Destination != nil {
		destination = this.RequestTemplate.Destination.BuildView()
	}

	if this.RequestTemplate.Scheme != nil {
		scheme = this.RequestTemplate.Scheme.BuildView()
	}

	if this.RequestTemplate.Query != nil {
		query = this.RequestTemplate.Query.BuildView()
	}

	if this.RequestTemplate.Body != nil {
		body = this.RequestTemplate.Body.BuildView()
	}

	return v2.RequestResponsePairViewV2{
		Request: v2.RequestDetailsViewV2{
			Path:        path,
			Method:      method,
			Destination: destination,
			Scheme:      scheme,
			Query:       query,
			Body:        body,
			Headers:     this.RequestTemplate.Headers,
		},
		Response: this.Response.ConvertToResponseDetailsView(),
	}
}

type RequestTemplate struct {
	Path        *RequestFieldMatchers
	Method      *RequestFieldMatchers
	Destination *RequestFieldMatchers
	Scheme      *RequestFieldMatchers
	Query       *RequestFieldMatchers
	Body        *RequestFieldMatchers
	Headers     map[string][]string
}

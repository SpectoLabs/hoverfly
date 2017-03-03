package models

import (
	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
)

type RequestFieldMatchers struct {
	ExactMatch *string
}

type RequestTemplateResponsePair struct {
	RequestTemplate RequestTemplate
	Response        ResponseDetails
}

func (this *RequestTemplateResponsePair) ConvertToRequestResponsePairView() v2.RequestResponsePairViewV1 {

	var path, destination, scheme, query *string

	if this.RequestTemplate.Path != nil {
		path = this.RequestTemplate.Path.ExactMatch
	}

	if this.RequestTemplate.Destination != nil {
		destination = this.RequestTemplate.Destination.ExactMatch
	}

	if this.RequestTemplate.Scheme != nil {
		scheme = this.RequestTemplate.Scheme.ExactMatch
	}

	if this.RequestTemplate.Query != nil {
		query = this.RequestTemplate.Query.ExactMatch
	}

	return v2.RequestResponsePairViewV1{
		Request: v2.RequestDetailsViewV1{
			Path:        path,
			Method:      this.RequestTemplate.Method,
			Destination: destination,
			Scheme:      scheme,
			Query:       query,
			Body:        this.RequestTemplate.Body,
			Headers:     this.RequestTemplate.Headers,
		},
		Response: this.Response.ConvertToResponseDetailsView(),
	}
}

type RequestTemplate struct {
	Path        *RequestFieldMatchers
	Method      *string
	Destination *RequestFieldMatchers
	Scheme      *RequestFieldMatchers
	Query       *RequestFieldMatchers
	Body        *string
	Headers     map[string][]string
}

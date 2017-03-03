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

	var path, method, destination, scheme, query, body *string

	if this.RequestTemplate.Path != nil {
		path = this.RequestTemplate.Path.ExactMatch
	}

	if this.RequestTemplate.Method != nil {
		method = this.RequestTemplate.Method.ExactMatch
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

	if this.RequestTemplate.Body != nil {
		body = this.RequestTemplate.Body.ExactMatch
	}

	return v2.RequestResponsePairViewV1{
		Request: v2.RequestDetailsViewV1{
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

package models

import (
	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
)

type RequestTemplateResponsePair struct {
	RequestTemplate RequestTemplate
	Response        ResponseDetails
}

func (this *RequestTemplateResponsePair) ConvertToRequestResponsePairView() v2.RequestResponsePairViewV1 {

	return v2.RequestResponsePairViewV1{
		Request: v2.RequestDetailsViewV1{
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

type RequestTemplate struct {
	Path        *string
	Method      *string
	Destination *string
	Scheme      *string
	Query       *string
	Body        *string
	Headers     map[string][]string
}

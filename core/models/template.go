package models

import (
	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
)

type RequestTemplateResponsePair struct {
	RequestTemplate RequestTemplate
	Response        ResponseDetails
}

func (this *RequestTemplateResponsePair) ConvertToRequestResponsePairView() v2.RequestResponsePairView {

	return v2.RequestResponsePairView{
		Request: v2.RequestDetailsView{
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

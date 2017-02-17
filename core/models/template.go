package models

import (
	"github.com/SpectoLabs/hoverfly/core/handlers/v1"
	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	. "github.com/SpectoLabs/hoverfly/core/util"
)

type RequestTemplateResponsePair struct {
	RequestTemplate RequestTemplate
	Response        ResponseDetails
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
		Response: this.Response.ConvertToV1ResponseDetailsView(),
	}
}

// DEPRICATED - Once we remove the v1 API, this will also go
func (this *RequestTemplateResponsePair) ConvertToV1RequestResponsePairView() v1.RequestResponsePairView {

	return v1.RequestResponsePairView{
		Request: v1.RequestDetailsView{
			RequestType: StringToPointer("template"),
			Path:        this.RequestTemplate.Path,
			Method:      this.RequestTemplate.Method,
			Destination: this.RequestTemplate.Destination,
			Scheme:      this.RequestTemplate.Scheme,
			Query:       this.RequestTemplate.Query,
			Body:        this.RequestTemplate.Body,
			Headers:     this.RequestTemplate.Headers,
		},
		Response: this.Response.ConvertToV1ResponseDetailsView(),
	}
}

func (this *RequestTemplateResponsePair) ConvertToRequestResponsePairView() v2.RequestResponsePairView {

	return v2.RequestResponsePairView{
		Request: v2.RequestDetailsView{
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

type RequestTemplate struct {
	Path        *string
	Method      *string
	Destination *string
	Scheme      *string
	Query       *string
	Body        *string
	Headers     map[string][]string
}

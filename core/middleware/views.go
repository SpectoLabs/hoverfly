package middleware

import (
	"github.com/SpectoLabs/hoverfly/core/interfaces"
	"github.com/SpectoLabs/hoverfly/core/util"
)

// This is a JSON serializable representation of the internal
// Hoverfly structs for HTTP requests and responses.
// These structs are only used when serializing requests
// and responses to middleware.
type RequestResponsePairView struct {
	Response ResponseDetailsView `json:"response"`
	Request  RequestDetailsView  `json:"request"`
}

func (this RequestResponsePairView) GetResponse() interfaces.Response { return this.Response }

func (this RequestResponsePairView) GetRequest() interfaces.Request { return this.Request }

type RequestDetailsView struct {
	RequestType *string             `json:"requestType"`
	Path        *string             `json:"path"`
	Method      *string             `json:"method"`
	Destination *string             `json:"destination"`
	Scheme      *string             `json:"scheme"`
	Query       *string             `json:"query"`
	Body        *string             `json:"body"`
	Headers     map[string][]string `json:"headers"`
}

func (this RequestDetailsView) GetPath() *string { return this.Path }

func (this RequestDetailsView) GetMethod() *string { return this.Method }

func (this RequestDetailsView) GetDestination() *string { return this.Destination }

func (this RequestDetailsView) GetScheme() *string { return this.Scheme }

func (this RequestDetailsView) GetQuery() *string {
	if this.Query == nil {
		return this.Query
	}
	queryString := util.SortQueryString(*this.Query)
	return &queryString
}

func (this RequestDetailsView) GetBody() *string { return this.Body }

func (this RequestDetailsView) GetHeaders() map[string][]string { return this.Headers }

type ResponseDetailsView struct {
	Status      int                 `json:"status"`
	Body        string              `json:"body"`
	EncodedBody bool                `json:"encodedBody"`
	Headers     map[string][]string `json:"headers"`
}

func (this ResponseDetailsView) GetStatus() int { return this.Status }

func (this ResponseDetailsView) GetBody() string { return this.Body }

func (this ResponseDetailsView) GetEncodedBody() bool { return this.EncodedBody }

func (this ResponseDetailsView) GetTemplated() bool { return false }

func (this ResponseDetailsView) GetTransitionsState() map[string]string { return nil }

func (this ResponseDetailsView) GetRemovesState() []string { return nil }

func (this ResponseDetailsView) GetHeaders() map[string][]string { return this.Headers }

package v2

import (
	"github.com/SpectoLabs/hoverfly/core/handlers/v1"
	"github.com/SpectoLabs/hoverfly/core/interfaces"
	"github.com/SpectoLabs/hoverfly/core/util"
	valid "github.com/gima/govalid/v1"
)

type SimulationViewV1 struct {
	DataViewV1 `json:"data"`
	MetaView   `json:"meta"`
}

func (this SimulationViewV1) GetValidationSchema() valid.Validator {
	return valid.Object(
		valid.ObjKV("data", valid.Object(
			valid.ObjKV("pairs", valid.Array(valid.ArrEach(valid.Optional(valid.Object(
				valid.ObjKV("request", valid.Object(
					valid.ObjKV("path", valid.Optional(valid.String())),
					valid.ObjKV("method", valid.Optional(valid.String())),
					valid.ObjKV("scheme", valid.Optional(valid.String())),
					valid.ObjKV("query", valid.Optional(valid.String())),
					valid.ObjKV("body", valid.Optional(valid.String())),
					valid.ObjKV("headers", valid.Optional(valid.Object())),
				)),
				valid.ObjKV("response", valid.Object(
					valid.ObjKV("status", valid.Optional(valid.Number())),
					valid.ObjKV("body", valid.Optional(valid.String())),
					valid.ObjKV("encodedBody", valid.Optional(valid.Boolean())),
					valid.ObjKV("headers", valid.Optional(valid.Object())),
				)),
			))))),
			valid.ObjKV("globalActions", valid.Optional(valid.Object(
				valid.ObjKV("delays", valid.Array(valid.ArrEach(valid.Optional(valid.Object(
					valid.ObjKV("urlPattern", valid.Optional(valid.String())),
					valid.ObjKV("httpMethod", valid.Optional(valid.String())),
					valid.ObjKV("delay", valid.Optional(valid.Number())),
				))))),
			))),
		)),
		valid.ObjKV("meta", valid.Object(
			valid.ObjKV("schemaVersion", valid.String()),
		)),
	)
}

type DataViewV1 struct {
	RequestResponsePairViewV1 []RequestResponsePairViewV1 `json:"pairs"`
	GlobalActions             GlobalActionsView           `json:"globalActions"`
}

type RequestResponsePairViewV1 struct {
	Response ResponseDetailsView  `json:"response"`
	Request  RequestDetailsViewV1 `json:"request"`
}

//Gets Response - required for interfaces.RequestResponsePairView
func (this RequestResponsePairViewV1) GetResponse() interfaces.Response { return this.Response }

//Gets Request - required for interfaces.RequestResponsePairView
func (this RequestResponsePairViewV1) GetRequest() interfaces.Request { return this.Request }

// RequestDetailsView is used when marshalling and unmarshalling RequestDetails
type RequestDetailsViewV1 struct {
	Path        *string             `json:"path"`
	Method      *string             `json:"method"`
	Destination *string             `json:"destination"`
	Scheme      *string             `json:"scheme"`
	Query       *string             `json:"query"`
	Body        *string             `json:"body"`
	Headers     map[string][]string `json:"headers"`
}

//Gets Path - required for interfaces.Request
func (this RequestDetailsViewV1) GetPath() *string { return this.Path }

//Gets Method - required for interfaces.Request
func (this RequestDetailsViewV1) GetMethod() *string { return this.Method }

//Gets Destination - required for interfaces.Request
func (this RequestDetailsViewV1) GetDestination() *string { return this.Destination }

//Gets Scheme - required for interfaces.Request
func (this RequestDetailsViewV1) GetScheme() *string { return this.Scheme }

//Gets Query - required for interfaces.Request
func (this RequestDetailsViewV1) GetQuery() *string {
	if this.Query == nil {
		return this.Query
	}
	queryString := util.SortQueryString(*this.Query)
	return &queryString
}

//Gets Body - required for interfaces.Request
func (this RequestDetailsViewV1) GetBody() *string { return this.Body }

//Gets Headers - required for interfaces.Request
func (this RequestDetailsViewV1) GetHeaders() map[string][]string { return this.Headers }

// ResponseDetailsView is used when marshalling and
// unmarshalling requests. This struct's Body may be Base64
// encoded based on the EncodedBody field.
type ResponseDetailsView struct {
	Status      int                 `json:"status"`
	Body        string              `json:"body"`
	EncodedBody bool                `json:"encodedBody"`
	Headers     map[string][]string `json:"headers"`
}

//Gets Status - required for interfaces.Response
func (this ResponseDetailsView) GetStatus() int { return this.Status }

// Gets Body - required for interfaces.Response
func (this ResponseDetailsView) GetBody() string { return this.Body }

// Gets EncodedBody - required for interfaces.Response
func (this ResponseDetailsView) GetEncodedBody() bool { return this.EncodedBody }

// Gets Headers - required for interfaces.Response
func (this ResponseDetailsView) GetHeaders() map[string][]string { return this.Headers }

type GlobalActionsView struct {
	Delays []v1.ResponseDelayView `json:"delays"`
}

type MetaView struct {
	SchemaVersion   string `json:"schemaVersion"`
	HoverflyVersion string `json:"hoverflyVersion"`
	TimeExported    string `json:"timeExported"`
}

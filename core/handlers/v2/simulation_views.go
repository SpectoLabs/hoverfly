package v2

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/SpectoLabs/hoverfly/core/handlers/v1"
	"github.com/SpectoLabs/hoverfly/core/interfaces"
	"github.com/SpectoLabs/hoverfly/core/util"
	valid "github.com/gima/govalid/v1"
)

func NewSimulationViewFromResponseBody(responseBody []byte) (SimulationViewV2, error) {
	var simulationView SimulationViewV2

	jsonMap := make(map[string]interface{})

	if err := json.Unmarshal(responseBody, &jsonMap); err != nil {
		return SimulationViewV2{}, errors.New("Invalid JSON")
	}

	if jsonMap["meta"] == nil {
		return SimulationViewV2{}, errors.New("Invalid JSON, missing \"meta\" object")
	}

	schemaVersion := jsonMap["meta"].(map[string]interface{})["schemaVersion"].(string)

	if schemaVersion == "v2" {
		if path, err := simulationView.GetValidationSchema().Validate(jsonMap); err != nil {

			return SimulationViewV2{}, errors.New("Invalid v2 simulation: " + path)
		}
		err := json.Unmarshal(responseBody, &simulationView)
		if err != nil {
			return SimulationViewV2{}, err
		}
	} else {
		var simulationViewV1 SimulationViewV1

		if path, err := simulationViewV1.GetValidationSchema().Validate(jsonMap); err != nil {
			return SimulationViewV2{}, errors.New("Invalid v1 simulation: " + path)
		}

		err := json.Unmarshal(responseBody, &simulationViewV1)
		if err != nil {
			return SimulationViewV2{}, err
		}

		simulationView = simulationViewV1.Upgrade()
	}

	return simulationView, nil
}

type SimulationViewV2 struct {
	DataViewV2 `json:"data"`
	MetaView   `json:"meta"`
}

func (this SimulationViewV2) GetValidationSchema() valid.Validator {
	return valid.Object(
		valid.ObjKV("data", valid.Object(
			valid.ObjKV("pairs", valid.Array(valid.ArrEach(valid.Optional(valid.Object(
				valid.ObjKV("request", valid.Object(
					valid.ObjKV("path", valid.Optional(valid.Object())),
					valid.ObjKV("method", valid.Optional(valid.Object())),
					valid.ObjKV("scheme", valid.Optional(valid.Object())),
					valid.ObjKV("query", valid.Optional(valid.Object())),
					valid.ObjKV("body", valid.Optional(valid.Object())),
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

type SimulationViewV1 struct {
	DataViewV1 `json:"data"`
	MetaView   `json:"meta"`
}

func (this SimulationViewV1) Upgrade() SimulationViewV2 {
	var pairs []RequestResponsePairViewV2
	for _, pairV1 := range this.RequestResponsePairViewV1 {

		var schemeMatchers, methodMatchers, destinationMatchers, pathMatchers, queryMatchers, bodyMatchers *RequestFieldMatchersView
		if pairV1.Request.Scheme != nil {
			schemeMatchers = &RequestFieldMatchersView{
				ExactMatch: pairV1.Request.Scheme,
			}
		}

		if pairV1.Request.Method != nil {
			methodMatchers = &RequestFieldMatchersView{
				ExactMatch: pairV1.Request.Method,
			}
		}

		if pairV1.Request.Destination != nil {
			destinationMatchers = &RequestFieldMatchersView{
				ExactMatch: pairV1.Request.Destination,
			}
		}

		if pairV1.Request.Path != nil {
			pathMatchers = &RequestFieldMatchersView{
				ExactMatch: pairV1.Request.Path,
			}
		}

		if pairV1.Request.Query != nil {
			queryMatchers = &RequestFieldMatchersView{
				ExactMatch: pairV1.Request.Query,
			}
		}

		if pairV1.Request.Body != nil {
			bodyMatchers = &RequestFieldMatchersView{
				ExactMatch: pairV1.Request.Body,
			}
		}

		pair := RequestResponsePairViewV2{
			Request: RequestDetailsViewV2{
				Scheme:      schemeMatchers,
				Method:      methodMatchers,
				Destination: destinationMatchers,
				Path:        pathMatchers,
				Query:       queryMatchers,
				Body:        bodyMatchers,
				Headers:     pairV1.Request.Headers,
			},
			Response: pairV1.Response,
		}
		pairs = append(pairs, pair)
	}

	return SimulationViewV2{
		DataViewV2{
			RequestResponsePairs: pairs,
		},
		MetaView{
			SchemaVersion:   "v2",
			HoverflyVersion: this.HoverflyVersion,
			TimeExported:    this.TimeExported,
		},
	}
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

type DataViewV2 struct {
	RequestResponsePairs []RequestResponsePairViewV2 `json:"pairs"`
	GlobalActions        GlobalActionsView           `json:"globalActions"`
}

type DataViewV1 struct {
	RequestResponsePairViewV1 []RequestResponsePairViewV1 `json:"pairs"`
	GlobalActions             GlobalActionsView           `json:"globalActions"`
}

type RequestResponsePairViewV2 struct {
	Response ResponseDetailsView  `json:"response"`
	Request  RequestDetailsViewV2 `json:"request"`
}

//Gets Response - required for interfaces.RequestResponsePairView
func (this RequestResponsePairViewV2) GetResponse() interfaces.Response { return this.Response }

type RequestResponsePairViewV1 struct {
	Response ResponseDetailsView  `json:"response"`
	Request  RequestDetailsViewV1 `json:"request"`
}

//Gets Response - required for interfaces.RequestResponsePairView
func (this RequestResponsePairViewV1) GetResponse() interfaces.Response { return this.Response }

//Gets Request - required for interfaces.RequestResponsePairView
func (this RequestResponsePairViewV1) GetRequest() interfaces.Request { return this.Request }

type RequestFieldMatchersView struct {
	ExactMatch *string `json:"exactMatch,omitempty"`
	XpathMatch *string `json:"xpathMatch,omitempty"`
	JsonMatch  *string `json:"jsonMatch,omitempty"`
	RegexMatch *string `json:"regexMatch,omitempty"`
}

// RequestDetailsView is used when marshalling and unmarshalling RequestDetails
type RequestDetailsViewV2 struct {
	Path        *RequestFieldMatchersView `json:"path,omitempty"`
	Method      *RequestFieldMatchersView `json:"method,omitempty"`
	Destination *RequestFieldMatchersView `json:"destination,omitempty"`
	Scheme      *RequestFieldMatchersView `json:"scheme,omitempty"`
	Query       *RequestFieldMatchersView `json:"query,omitempty"`
	Body        *RequestFieldMatchersView `json:"body,omitempty"`
	Headers     map[string][]string       `json:"headers,omitempty"`
}

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

func NewMetaView(version string) *MetaView {
	return &MetaView{
		HoverflyVersion: version,
		SchemaVersion:   "v2",
		TimeExported:    time.Now().Format(time.RFC3339),
	}
}

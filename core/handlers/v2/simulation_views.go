package v2

import (
	"encoding/json"
	"errors"
	"time"

	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/core/handlers/v1"
	"github.com/SpectoLabs/hoverfly/core/interfaces"
	"github.com/SpectoLabs/hoverfly/core/util"
	"github.com/xeipuuv/gojsonschema"
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

	if jsonMap["meta"].(map[string]interface{})["schemaVersion"] == nil {
		return SimulationViewV2{}, errors.New("Invalid JSON, missing \"meta.schemaVersion\" string")
	}

	schemaVersion := jsonMap["meta"].(map[string]interface{})["schemaVersion"].(string)

	if schemaVersion == "v2" {
		err := ValidateSimulation(jsonMap, SimulationViewV2Schema)
		if err != nil {
			return simulationView, errors.New("Invalid v2 simulation:" + err.Error())
		}

		err = json.Unmarshal(responseBody, &simulationView)
		if err != nil {
			return SimulationViewV2{}, err
		}
	} else if schemaVersion == "v1" {
		err := ValidateSimulation(jsonMap, SimulationViewV1Schema)
		if err != nil {
			return simulationView, errors.New("Invalid v1 simulation:" + err.Error())
		}

		var simulationViewV1 SimulationViewV1

		err = json.Unmarshal(responseBody, &simulationViewV1)
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

type SimulationViewV1 struct {
	DataViewV1 `json:"data"`
	MetaView   `json:"meta"`
}

func ValidateSimulation(json, schema map[string]interface{}) error {
	jsonLoader := gojsonschema.NewGoLoader(json)
	schemaLoader := gojsonschema.NewGoLoader(schema)

	result, err := gojsonschema.Validate(schemaLoader, jsonLoader)
	if err != nil {
		log.Error("Error when validating simulaton: " + err.Error())
		return errors.New("Error when validating simulaton")
	}

	if !result.Valid() {
		errorMessage := ""
		for i, parsingError := range result.Errors() {
			message := strings.Split(parsingError.String(), ":")[1]
			var comma string
			if i != 0 {
				comma = ","
			}
			errorMessage = errorMessage + comma + " " + strings.TrimSpace(message)
		}
		return errors.New(errorMessage)
	}

	return nil
}

func (this SimulationViewV1) Upgrade() SimulationViewV2 {
	var pairs []RequestResponsePairViewV2
	for _, pairV1 := range this.RequestResponsePairViewV1 {

		var schemeMatchers, methodMatchers, destinationMatchers, pathMatchers, queryMatchers, bodyMatchers *RequestFieldMatchersView
		var headers map[string][]string

		isNotRecording := pairV1.Request.RequestType != nil && *pairV1.Request.RequestType != "recording"

		if isNotRecording {
			headers = pairV1.Request.Headers
		}
		if pairV1.Request.Scheme != nil {

			if isNotRecording {
				schemeMatchers = &RequestFieldMatchersView{
					GlobMatch: pairV1.Request.Scheme,
				}
			} else {
				schemeMatchers = &RequestFieldMatchersView{
					ExactMatch: pairV1.Request.Scheme,
				}
			}
		}

		if pairV1.Request.Method != nil {

			if isNotRecording {
				methodMatchers = &RequestFieldMatchersView{
					GlobMatch: pairV1.Request.Method,
				}
			} else {
				methodMatchers = &RequestFieldMatchersView{
					ExactMatch: pairV1.Request.Method,
				}
			}
		}

		if pairV1.Request.Destination != nil {
			if isNotRecording {
				destinationMatchers = &RequestFieldMatchersView{
					GlobMatch: pairV1.Request.Destination,
				}
			} else {
				destinationMatchers = &RequestFieldMatchersView{
					ExactMatch: pairV1.Request.Destination,
				}
			}
		}

		if pairV1.Request.Path != nil {
			if isNotRecording {
				pathMatchers = &RequestFieldMatchersView{
					GlobMatch: pairV1.Request.Path,
				}
			} else {
				pathMatchers = &RequestFieldMatchersView{
					ExactMatch: pairV1.Request.Path,
				}
			}
		}

		if pairV1.Request.Query != nil {
			if isNotRecording {
				queryMatchers = &RequestFieldMatchersView{
					GlobMatch: pairV1.Request.Query,
				}
			} else {
				queryMatchers = &RequestFieldMatchersView{
					ExactMatch: pairV1.Request.Query,
				}
			}
		}

		if pairV1.Request.Body != nil {
			if isNotRecording {
				bodyMatchers = &RequestFieldMatchersView{
					GlobMatch: pairV1.Request.Body,
				}
			} else {
				bodyMatchers = &RequestFieldMatchersView{
					ExactMatch: pairV1.Request.Body,
				}
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
				Headers:     headers,
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
	ExactMatch    *string `json:"exactMatch,omitempty"`
	XmlMatch      *string `json:"xmlMatch,omitempty"`
	XpathMatch    *string `json:"xpathMatch,omitempty"`
	JsonMatch     *string `json:"jsonMatch,omitempty"`
	JsonPathMatch *string `json:"jsonPathMatch,omitempty"`
	RegexMatch    *string `json:"regexMatch,omitempty"`
	GlobMatch     *string `json:"globMatch,omitempty"`
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
	RequestType *string             `json:"requestType"`
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

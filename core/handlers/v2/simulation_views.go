package v2

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/core/handlers/v1"
	"github.com/SpectoLabs/hoverfly/core/interfaces"
	"github.com/SpectoLabs/hoverfly/core/util"
	"github.com/xeipuuv/gojsonschema"
)

func NewSimulationViewFromResponseBody(responseBody []byte) (SimulationViewV3, error) {
	var simulationView SimulationViewV3

	jsonMap := make(map[string]interface{})

	if err := json.Unmarshal(responseBody, &jsonMap); err != nil {
		return SimulationViewV3{}, errors.New("Invalid JSON")
	}

	if jsonMap["meta"] == nil {
		return SimulationViewV3{}, errors.New("Invalid JSON, missing \"meta\" object")
	}

	if jsonMap["meta"].(map[string]interface{})["schemaVersion"] == nil {
		return SimulationViewV3{}, errors.New("Invalid JSON, missing \"meta.schemaVersion\" string")
	}

	schemaVersion := jsonMap["meta"].(map[string]interface{})["schemaVersion"].(string)

	if schemaVersion == "v3" || schemaVersion == "v2" {
		err := ValidateSimulation(jsonMap, SimulationViewV3Schema)
		if err != nil {
			return simulationView, errors.New(fmt.Sprintf("Invalid %s simulation:", schemaVersion) + err.Error())
		}

		err = json.Unmarshal(responseBody, &simulationView)
		if err != nil {
			return SimulationViewV3{}, err
		}
	} else if schemaVersion == "v1" {
		err := ValidateSimulation(jsonMap, SimulationViewV1Schema)
		if err != nil {
			return simulationView, errors.New("Invalid v1 simulation:" + err.Error())
		}

		var simulationViewV1 SimulationViewV1

		err = json.Unmarshal(responseBody, &simulationViewV1)
		if err != nil {
			return SimulationViewV3{}, err
		}

		simulationView = simulationViewV1.Upgrade()
	} else {
		return simulationView, fmt.Errorf("Invalid simulation: schema version %v is not supported by this version of Hoverfly, you may need to update Hoverfly", schemaVersion)
	}

	simulationView.MetaView.SchemaVersion = "v3"
	return simulationView, nil
}

type SimulationViewV3 struct {
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

func (this SimulationViewV1) Upgrade() SimulationViewV3 {
	var pairs []RequestMatcherResponsePairViewV3
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

		pair := RequestMatcherResponsePairViewV3{
			RequestMatcher: RequestMatcherViewV2{
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

	return SimulationViewV3{
		DataViewV2{
			RequestResponsePairs: pairs,
		},
		MetaView{
			SchemaVersion:   "v3",
			HoverflyVersion: this.HoverflyVersion,
			TimeExported:    this.TimeExported,
		},
	}
}

type DataViewV2 struct {
	RequestResponsePairs []RequestMatcherResponsePairViewV3 `json:"pairs"`
	GlobalActions        GlobalActionsView                  `json:"globalActions"`
}

type DataViewV1 struct {
	RequestResponsePairViewV1 []RequestResponsePairViewV1 `json:"pairs"`
	GlobalActions             GlobalActionsView           `json:"globalActions"`
}

type RequestMatcherResponsePairViewV3 struct {
	Response       ResponseDetailsView  `json:"response"`
	RequestMatcher RequestMatcherViewV2 `json:"request"`
}

type ClosestMissView struct {
	Response       ResponseDetailsView  `json:"response"`
	RequestMatcher RequestMatcherViewV2 `json:"requestMatcher"`
	MissedFields   []string             `json:"missedFields"`
}

//Gets Response - required for interfaces.RequestResponsePairView
func (this RequestMatcherResponsePairViewV3) GetResponse() interfaces.Response { return this.Response }

type RequestResponsePairViewV1 struct {
	Response ResponseDetailsView `json:"response"`
	Request  RequestDetailsView  `json:"request"`
}

//Gets Response - required for interfaces.RequestResponsePairView
func (this RequestResponsePairViewV1) GetResponse() interfaces.Response { return this.Response }

//Gets RequestMatcher - required for interfaces.RequestResponsePairView
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
type RequestMatcherViewV2 struct {
	Path        *RequestFieldMatchersView `json:"path,omitempty"`
	Method      *RequestFieldMatchersView `json:"method,omitempty"`
	Destination *RequestFieldMatchersView `json:"destination,omitempty"`
	Scheme      *RequestFieldMatchersView `json:"scheme,omitempty"`
	Query       *RequestFieldMatchersView `json:"query,omitempty"`
	Body        *RequestFieldMatchersView `json:"body,omitempty"`
	Headers     map[string][]string       `json:"headers,omitempty"`
}

// RequestDetailsView is used when marshalling and unmarshalling RequestDetails
type RequestDetailsView struct {
	RequestType *string             `json:"requestType,omitempty"`
	Path        *string             `json:"path"`
	Method      *string             `json:"method"`
	Destination *string             `json:"destination"`
	Scheme      *string             `json:"scheme"`
	Query       *string             `json:"query"`
	Body        *string             `json:"body"`
	Headers     map[string][]string `json:"headers"`
}

//Gets Path - required for interfaces.RequestMatcher
func (this RequestDetailsView) GetPath() *string { return this.Path }

//Gets Method - required for interfaces.RequestMatcher
func (this RequestDetailsView) GetMethod() *string { return this.Method }

//Gets Destination - required for interfaces.RequestMatcher
func (this RequestDetailsView) GetDestination() *string { return this.Destination }

//Gets Scheme - required for interfaces.RequestMatcher
func (this RequestDetailsView) GetScheme() *string { return this.Scheme }

//Gets Query - required for interfaces.RequestMatcher
func (this RequestDetailsView) GetQuery() *string {
	if this.Query == nil {
		return this.Query
	}
	queryString := util.SortQueryString(*this.Query)
	return &queryString
}

//Gets Body - required for interfaces.RequestMatcher
func (this RequestDetailsView) GetBody() *string { return this.Body }

//Gets Headers - required for interfaces.RequestMatcher
func (this RequestDetailsView) GetHeaders() map[string][]string { return this.Headers }

// ResponseDetailsView is used when marshalling and
// unmarshalling requests. This struct's Body may be Base64
// encoded based on the EncodedBody field.
type ResponseDetailsView struct {
	Status      int                 `json:"status"`
	Body        string              `json:"body"`
	EncodedBody bool                `json:"encodedBody"`
	Headers     map[string][]string `json:"headers,omitempty"`
	Templated   bool                `json:"templated"`
}

//Gets Status - required for interfaces.Response
func (this ResponseDetailsView) GetStatus() int { return this.Status }

// Gets Body - required for interfaces.Response
func (this ResponseDetailsView) GetBody() string { return this.Body }

// Gets EncodedBody - required for interfaces.Response
func (this ResponseDetailsView) GetEncodedBody() bool { return this.EncodedBody }

func (this ResponseDetailsView) GetTemplated() bool { return this.Templated }

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
		SchemaVersion:   "v3",
		TimeExported:    time.Now().Format(time.RFC3339),
	}
}

package v2

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"strings"

	"github.com/SpectoLabs/hoverfly/core/handlers/v1"
	log "github.com/sirupsen/logrus"
	"github.com/xeipuuv/gojsonschema"
)

func NewSimulationViewFromRequestBody(requestBody []byte) (SimulationViewV5, error) {
	var simulationView SimulationViewV5

	jsonMap := make(map[string]interface{})

	if err := json.Unmarshal(requestBody, &jsonMap); err != nil {
		return SimulationViewV5{}, errors.New("Invalid JSON")
	}

	if jsonMap["meta"] == nil {
		return SimulationViewV5{}, errors.New("Invalid JSON, missing \"meta\" object")
	}

	if jsonMap["meta"].(map[string]interface{})["schemaVersion"] == nil {
		return SimulationViewV5{}, errors.New("Invalid JSON, missing \"meta.schemaVersion\" string")
	}

	schemaVersion := jsonMap["meta"].(map[string]interface{})["schemaVersion"].(string)

	if schemaVersion == "v5" || schemaVersion == "v5.1" {

		err := ValidateSimulation(jsonMap, SimulationViewV5Schema)
		if err != nil {
			return simulationView, errors.New(fmt.Sprintf("Invalid %s simulation: ", schemaVersion) + err.Error())
		}

		err = json.Unmarshal(requestBody, &simulationView)
		if err != nil {
			return SimulationViewV5{}, err
		}
	} else if schemaVersion == "v4" || schemaVersion == "v3" {
		err := ValidateSimulation(jsonMap, SimulationViewV4Schema)
		if err != nil {
			return simulationView, errors.New(fmt.Sprintf("Invalid %s simulation: ", schemaVersion) + err.Error())
		}

		var simulationViewV4 SimulationViewV4

		err = json.Unmarshal(requestBody, &simulationViewV4)
		if err != nil {
			return SimulationViewV5{}, err
		}

		simulationView = upgradeV4(simulationViewV4)
	} else if schemaVersion == "v2" {
		err := ValidateSimulation(jsonMap, SimulationViewV2Schema)
		if err != nil {
			return simulationView, errors.New(fmt.Sprintf("Invalid %s simulation: ", schemaVersion) + err.Error())
		}

		var simulationViewV2 SimulationViewV2

		err = json.Unmarshal(requestBody, &simulationViewV2)
		if err != nil {
			return SimulationViewV5{}, err
		}

		simulationView = upgradeV2(simulationViewV2)
	} else if schemaVersion == "v1" {
		err := ValidateSimulation(jsonMap, SimulationViewV1Schema)
		if err != nil {
			return simulationView, errors.New("Invalid v1 simulation: " + err.Error())
		}

		var simulationViewV1 SimulationViewV1

		err = json.Unmarshal(requestBody, &simulationViewV1)
		if err != nil {
			return SimulationViewV5{}, err
		}

		simulationView = upgradeV1(simulationViewV1)
	} else {
		return simulationView, fmt.Errorf("Invalid simulation: schema version %v is not supported by this version of Hoverfly, you may need to update Hoverfly", schemaVersion)
	}

	return simulationView, nil
}

func ValidateSimulation(json, schema map[string]interface{}) error {
	jsonLoader := gojsonschema.NewGoLoader(json)
	schemaLoader := gojsonschema.NewGoLoader(schema)

	result, err := gojsonschema.Validate(schemaLoader, jsonLoader)
	if err != nil {
		log.Error("Error when validating simulation: " + err.Error())
		return errors.New("Error when validating simulation" + err.Error())
	}

	if !result.Valid() {
		// TODO return as an array in a custom error struct
		var resultDetails []string
		for _, parsingError := range result.Errors() {
			resultDetails = append(resultDetails, fmt.Sprintf("Error for <%s>: %s", parsingError.Field(), parsingError.Description()))
		}

		errorMessage := fmt.Sprintf("[%s]", strings.Join(resultDetails, "; "))
		return errors.New(errorMessage)
	}

	return nil
}

type GlobalActionsView struct {
	Delays          []v1.ResponseDelayView          `json:"delays"`
	DelaysLogNormal []v1.ResponseDelayLogNormalView `json:"delaysLogNormal"`
}

type MetaView struct {
	SchemaVersion   string `json:"schemaVersion"`
	HoverflyVersion string `json:"hoverflyVersion"`
	TimeExported    string `json:"timeExported"`
}

func NewMetaView(version string) *MetaView {
	return &MetaView{
		HoverflyVersion: version,
		SchemaVersion:   "v5.1",
		TimeExported:    time.Now().Format(time.RFC3339),
	}
}

func BuildSimulationView(
	pairViews []RequestMatcherResponsePairViewV5,
	delayView v1.ResponseDelayPayloadView,
	delayLogNormalView v1.ResponseDelayLogNormalPayloadView,
	version string,
) SimulationViewV5 {
	return SimulationViewV5{
		DataViewV5{
			RequestResponsePairs: pairViews,
			GlobalActions: GlobalActionsView{
				Delays:          delayView.Data,
				DelaysLogNormal: delayLogNormalView.Data,
			},
		},
		*NewMetaView(version),
	}
}

const deprecatedQueryMessage = "Usage of deprecated field `deprecatedQuery` on data.pairs[%v].request.deprecatedQuery, please update your simulation to use `query` field"
const deprecatedQueryDocs = "https://hoverfly.readthedocs.io/en/latest/pages/troubleshooting/troubleshooting.html#why-does-my-simulation-have-a-deprecatedquery-field"
const ContentLengthAndTransferEncodingMessage = "Response contains both Content-Length and Transfer-Encoding headers on data.pairs[%v].response, please remove one of these headers"
const BodyAndBodyFileMessage = "Response contains both `body` and `bodyFile` in data.pairs[%v].response, please remove one of them otherwise `body` is used if non empty"
const ContentLengthMismatchMessage = "Response contains incorrect Content-Length header on data.pairs[%v].response, please correct or remove header"
const pairIgnoredMessage = "data.pairs[%v] is not added due to a conflict with the existing simulation"

type SimulationImportResult struct {
	Err             error                     `json:"error,omitempty"`
	WarningMessages []SimulationImportWarning `json:"warnings,omitempty"`
}

type SimulationImportWarning struct {
	Message  string `json:"message,omitempty"`
	DocsLink string `json:"documentation,omitempty"`
}

func (s *SimulationImportResult) SetError(err error) {
	s.Err = err
}

func (s SimulationImportResult) GetError() error {
	return s.Err
}

func (s *SimulationImportResult) AddDeprecatedQueryWarning(requestNumber int) {
	warning := fmt.Sprintf("WARNING: %s", fmt.Sprintf(deprecatedQueryMessage, requestNumber))
	if s.WarningMessages == nil {
		s.WarningMessages = []SimulationImportWarning{}
	}
	s.WarningMessages = append(s.WarningMessages, SimulationImportWarning{Message: warning, DocsLink: deprecatedQueryDocs})
}

func (s *SimulationImportResult) AddContentLengthAndTransferEncodingWarning(requestNumber int) {
	warning := fmt.Sprintf("WARNING: %s", fmt.Sprintf(ContentLengthAndTransferEncodingMessage, requestNumber))
	if s.WarningMessages == nil {
		s.WarningMessages = []SimulationImportWarning{}
	}
	s.WarningMessages = append(s.WarningMessages, SimulationImportWarning{Message: warning})
}

func (s *SimulationImportResult) AddBodyAndBodyFileWarning(requestNumber int) {
	warning := fmt.Sprintf("WARNING: %s", fmt.Sprintf(BodyAndBodyFileMessage, requestNumber))
	if s.WarningMessages == nil {
		s.WarningMessages = []SimulationImportWarning{}
	}
	s.WarningMessages = append(s.WarningMessages, SimulationImportWarning{Message: warning})
}

func (s *SimulationImportResult) AddContentLengthMismatchWarning(requestNumber int) {
	warning := fmt.Sprintf("WARNING: %s", fmt.Sprintf(ContentLengthMismatchMessage, requestNumber))
	if s.WarningMessages == nil {
		s.WarningMessages = []SimulationImportWarning{}
	}
	s.WarningMessages = append(s.WarningMessages, SimulationImportWarning{Message: warning})
}

func (s *SimulationImportResult) AddPairIgnoredWarning(requestNumber int) {
	warning := fmt.Sprintf("WARNING: %s", fmt.Sprintf(pairIgnoredMessage, requestNumber))
	if s.WarningMessages == nil {
		s.WarningMessages = []SimulationImportWarning{}
	}
	s.WarningMessages = append(s.WarningMessages, SimulationImportWarning{Message: warning})
}

package v2

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/core/handlers/v1"
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

	if schemaVersion == "v3" {
		err := ValidateSimulation(jsonMap, SimulationViewV3Schema)
		if err != nil {
			return simulationView, errors.New(fmt.Sprintf("Invalid %s simulation:", schemaVersion) + err.Error())
		}

		err = json.Unmarshal(responseBody, &simulationView)
		if err != nil {
			return SimulationViewV3{}, err
		}
	} else if schemaVersion == "v2" {
		err := ValidateSimulation(jsonMap, SimulationViewV3Schema)
		if err != nil {
			return simulationView, errors.New(fmt.Sprintf("Invalid %s simulation:", schemaVersion) + err.Error())
		}

		var simulationViewV2 SimulationViewV2

		err = json.Unmarshal(responseBody, &simulationViewV2)
		if err != nil {
			return SimulationViewV3{}, err
		}

		simulationView = simulationViewV2.Upgrade()
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

package wrapper

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/hoverctl/configuration"
)

func ExportSimulation(target configuration.Target) ([]byte, error) {
	response, err := doRequest(target, "GET", v2ApiSimulation, "", nil)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	err = handleResponseError(response, "Could not retrieve simulation")
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Debug(err.Error())
		return nil, errors.New("Could not export from Hoverfly")
	}

	var jsonBytes bytes.Buffer
	err = json.Indent(&jsonBytes, body, "", "\t")
	if err != nil {
		log.Debug(err.Error())
		return nil, errors.New("Could not export from Hoverfly")
	}

	return jsonBytes.Bytes(), nil
}

func ImportSimulation(target configuration.Target, simulationData string) error {
	response, err := doRequest(target, "PUT", v2ApiSimulation, simulationData, nil)
	if err != nil {
		return err
	}

	err = handleResponseError(response, "Could not import simulation")
	if err != nil {
		return err
	}

	return nil
}

// Wipe will call the records endpoint in Hoverfly with a DELETE request, triggering Hoverfly to wipe the database
func DeleteSimulations(target configuration.Target) error {
	response, err := doRequest(target, "DELETE", v2ApiSimulation, "", nil)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	err = handleResponseError(response, "Could not delete simulation")
	if err != nil {
		return err
	}

	return nil
}

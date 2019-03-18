package wrapper

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"

	"fmt"
	"net/url"

	log "github.com/sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/hoverctl/configuration"
)

func ExportSimulation(target configuration.Target, urlPattern string) ([]byte, error) {

	requestUrl := v2ApiSimulation
	if len(urlPattern) > 0 {
		requestUrl = fmt.Sprintf("%s?urlPattern=%s", requestUrl, url.QueryEscape(urlPattern))
	}
	response, err := doRequest(target, "GET", requestUrl, "", nil)
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

	responseBytes, _ := ioutil.ReadAll(response.Body)

	result := &v2.SimulationImportResult{}
	json.Unmarshal(responseBytes, result)

	for _, warning := range result.WarningMessages {
		fmt.Println(warning.Message)
		fmt.Println(warning.DocsLink + "\n")
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

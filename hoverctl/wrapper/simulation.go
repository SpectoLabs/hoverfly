package wrapper

import (
	"encoding/json"
	"io/ioutil"

	"fmt"
	"net/url"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/hoverctl/configuration"
)

func ExportSimulation(target configuration.Target, urlPattern string) (v2.SimulationViewV6, error) {
	view := v2.SimulationViewV6{}
	requestUrl := v2ApiSimulation
	if len(urlPattern) > 0 {
		requestUrl = fmt.Sprintf("%s?urlPattern=%s", requestUrl, url.QueryEscape(urlPattern))
	}
	response, err := doRequest(target, "GET", requestUrl, "", nil)
	if err != nil {
		return view, err
	}

	defer response.Body.Close()

	err = handleResponseError(response, "Could not retrieve simulation")
	if err != nil {
		return view, err
	}

	err = json.NewDecoder(response.Body).Decode(&view)
	return view, err
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

func AddSimulation(target configuration.Target, simulationData string) error {
	response, err := doRequest(target, "POST", v2ApiSimulation, simulationData, nil)
	if err != nil {
		return err
	}

	err = handleResponseError(response, "Could not add simulation")
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

package wrapper

import (
	"encoding/json"
	"io/ioutil"

	"fmt"
	"net/url"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/hoverctl/configuration"
	log "github.com/sirupsen/logrus"
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

	var view v2.SimulationViewV6
	if err := json.NewDecoder(response.Body).Decode(&view); err != nil {
		log.Debug(err.Error())
		return nil, err
	}

	for i, pair := range view.DataViewV6.RequestResponsePairs {
		bodyFile := pair.Response.GetBodyFile()
		if len(bodyFile) == 0 {
			continue
		}

		if err := configuration.WriteFile(bodyFile, []byte(pair.Response.GetBody())); err != nil {
			log.Debug(err.Error())
			return nil, err
		}

		view.DataViewV6.RequestResponsePairs[i].Response.Body = ""
	}

	return json.MarshalIndent(view, "", "\t")
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

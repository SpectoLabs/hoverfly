package wrapper

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/hoverctl/configuration"
)

func GetLogs(target configuration.Target, format string) ([]string, error) {
	headers := map[string]string{
		"Accept": "text/plain",
	}

	if format == "json" {
		headers["Accept"] = "application/json"
	}

	response, err := doRequest(target, "GET", v2ApiLogs, "", headers)
	if err != nil {
		return []string{}, err
	}

	defer response.Body.Close()

	err = handleResponseError(response, "Could not retrieve logs")
	if err != nil {
		return nil, err
	}

	responseBody, _ := ioutil.ReadAll(response.Body)
	if format == "json" {
		var logsView v2.LogsView

		err := json.Unmarshal(responseBody, &logsView)
		if err != nil {
			return nil, err
		}

		var logs []string
		for _, log := range logsView.Logs {
			jsonLog, err := json.Marshal(log)
			if err != nil {
				return nil, err
			}

			logs = append(logs, string(jsonLog))
		}

		return logs, nil
	} else {
		return strings.Split(string(responseBody), "\n"), nil
	}
}

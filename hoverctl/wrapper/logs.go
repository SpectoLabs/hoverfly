package wrapper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/hoverctl/configuration"
)

func GetLogs(target configuration.Target, format string, filterTime *time.Time) ([]string, error) {
	headers := map[string]string{
		"Accept": "text/plain",
	}

	if format == "json" {
		headers["Accept"] = "application/json"
	}

	url := v2ApiLogs
	if filterTime != nil {
		url = fmt.Sprintf("%v%v%v", url, "?from=", filterTime.Unix())
	}

	response, err := doRequest(target, "GET", url, "", headers)
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
		if string(responseBody) == "" {
			return []string{}, err
		}
		logs := strings.Split(string(responseBody), "\n")
		lastLogPosition := len(logs) - 1
		if logs[lastLogPosition] == "" {
			return logs[:lastLogPosition], nil
		}

		return logs, nil
	}
}

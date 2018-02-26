package wrapper

import (
	"encoding/json"
	"io/ioutil"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/hoverctl/configuration"
)

func GetAllDiffs(target configuration.Target) ([]v2.ResponseDiffForRequestView, error) {

	res, err := doRequest(target, "GET", v2ApiDiff, "", nil)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	diffs := &v2.DiffView{}

	err = json.Unmarshal(responseBody, diffs)

	if err != nil {
		return nil, err
	}

	return diffs.Diff, nil
}

func DeleteAllDiffs(target configuration.Target) error {

	_, err := doRequest(target, "DELETE", v2ApiDiff, "", nil)

	return err
}

package wrapper

import (
	"encoding/json"

	v2 "github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/hoverctl/configuration"
)

func GetAllJournalIndexes(target configuration.Target) ([]v2.JournalIndexView, error) {

	response, err := doRequest(target, "GET", v2JournalIndex, "", nil)
	if err != nil {
		return []v2.JournalIndexView{}, err
	}

	defer response.Body.Close()

	err = handleResponseError(response, "Could not retrieve all journal indexes")
	if err != nil {
		return []v2.JournalIndexView{}, err
	}

	var journalIndexes []v2.JournalIndexView
	err = UnmarshalToInterface(response, &journalIndexes)
	if err != nil {
		return []v2.JournalIndexView{}, err
	}

	return journalIndexes, nil
}

func SetJournalIndex(indexName string, target configuration.Target) error {

	journalIndexRequestView := v2.JournalIndexRequestView{
		Name: indexName,
	}
	journalIndexRequestViewData, err := json.Marshal(journalIndexRequestView)
	if err != nil {
		return err
	}

	response, err := doRequest(target, "POST", v2JournalIndex, string(journalIndexRequestViewData), nil)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	err = handleResponseError(response, "Could not set journal index")
	if err != nil {
		return err
	}
	return nil
}

func DeleteJournalIndex(indexName string, target configuration.Target) error {

	response, err := doRequest(target, "DELETE", v2JournalIndex+"/"+indexName, "", nil)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	err = handleResponseError(response, "Could not delete journal index")
	if err != nil {
		return err
	}
	return nil

}

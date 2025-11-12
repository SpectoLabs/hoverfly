package wrapper

import (
	"encoding/json"

	v2 "github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/hoverctl/configuration"
)

func GetAllTemplateDataSources(target configuration.Target) (v2.TemplateDataSourceView, error) {

	response, err := doRequest(target, "GET", v2ApiTemplateDataSourceAction, "", nil)
	if err != nil {
		return v2.TemplateDataSourceView{}, err
	}

	defer response.Body.Close()

	err = handleResponseError(response, "Could not retrieve all template data sources")
	if err != nil {
		return v2.TemplateDataSourceView{}, err
	}

	var templateDataSourceView v2.TemplateDataSourceView

	err = UnmarshalToInterface(response, &templateDataSourceView)
	if err != nil {
		return v2.TemplateDataSourceView{}, err
	}

	return templateDataSourceView, nil
}

func SetCsvTemplateDataSource(dataSourceName, scriptContent string, target configuration.Target) error {

	csvDataSource := v2.CSVDataSourceView{
		Data: scriptContent,
		Name: dataSourceName,
	}
	marshalledCsvDataSource, err := json.Marshal(csvDataSource)
	if err != nil {
		return err
	}
	response, err := doRequest(target, "PUT", v2ApiTemplateDataSourceAction, string(marshalledCsvDataSource), nil)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	err = handleResponseError(response, "Could not set csv data source")
	if err != nil {
		return err
	}
	return nil
}

func DeleteCsvDataSource(dataSourceName string, target configuration.Target) error {

	response, err := doRequest(target, "DELETE", v2ApiTemplateDataSourceAction+"/"+dataSourceName, "", nil)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	err = handleResponseError(response, "Could not delete data source")
	if err != nil {
		return err
	}

	return nil
}

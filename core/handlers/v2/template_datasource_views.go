package v2

type TemplateDataSourceView struct {
	DataSources []CSVDataSourceView `json:"csvDataSources,omitempty"`
}

type CSVDataSourceView struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

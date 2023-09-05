package templating

import (
	"encoding/csv"
	"strings"
)

type DataSource struct {
	SourceType string
	Name       string
	Data       [][]string
}

func NewCsvDataSource(fileName, fileContent string) (*DataSource, error) {

	csvReader := csv.NewReader(strings.NewReader(fileContent))
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	return &DataSource{"csv", fileName, records}, nil
}

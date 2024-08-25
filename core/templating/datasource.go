package templating

import (
	"encoding/csv"
	"strings"
	"sync"

	v2 "github.com/SpectoLabs/hoverfly/core/handlers/v2"
	log "github.com/sirupsen/logrus"
)

type DataSource struct {
	SourceType string
	Name       string
	Data       [][]string
	mu         sync.Mutex
}

func NewCsvDataSource(fileName, fileContent string) (*DataSource, error) {

	csvReader := csv.NewReader(strings.NewReader(fileContent))
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	return &DataSource{"csv", fileName, records, sync.Mutex{}}, nil
}

func (dataSource *DataSource) GetDataSourceView() (v2.CSVDataSourceView, error) {

	content, err := getData(dataSource)
	if err != nil {
		return v2.CSVDataSourceView{}, err
	}
	return v2.CSVDataSourceView{Name: dataSource.Name, Data: content}, nil
}

func getData(source *DataSource) (string, error) {

	var csvData strings.Builder
	csvWriter := csv.NewWriter(&csvData)
	for _, row := range source.Data {
		if err := csvWriter.Write(row); err != nil {
			log.Error("error writing csv")
			return "", err
		}
	}
	csvWriter.Flush()
	if err := csvWriter.Error(); err != nil {
		log.Error(err)
		return "", err
	}
	return csvData.String(), nil
}

func dataSourceExists(dataSources map[string]*DataSource, name string) bool {
	for _, dataSource := range dataSources {
		if dataSource.Name == name {
			return true
		}
	}
	return false
}

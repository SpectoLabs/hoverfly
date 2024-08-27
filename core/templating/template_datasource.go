package templating

import (
	"sync"
)

type TemplateDataSource struct {
	DataSources map[string]*DataSource
	RWMutex     sync.RWMutex
}

func NewTemplateDataSource() *TemplateDataSource {

	return &TemplateDataSource{
		DataSources: make(map[string]*DataSource),
	}
}

func (templateDataSource *TemplateDataSource) SetDataSource(dataSourceName string, dataSource *DataSource) {

	templateDataSource.RWMutex.Lock()
	templateDataSource.DataSources[dataSourceName] = dataSource
	templateDataSource.RWMutex.Unlock()
}

func (templateDataSource *TemplateDataSource) DeleteDataSource(dataSourceName string) {

	templateDataSource.RWMutex.Lock()

	if _, ok := templateDataSource.DataSources[dataSourceName]; ok {
		delete(templateDataSource.DataSources, dataSourceName)
	}
	templateDataSource.RWMutex.Unlock()
}

func (templateDataSource *TemplateDataSource) GetAllDataSources() map[string]*DataSource {

	return templateDataSource.DataSources
}

func (templateDataSource *TemplateDataSource) DataSourceExists(name string) bool {
	templateDataSource.RWMutex.Lock()
	defer templateDataSource.RWMutex.Unlock()

	for _, dataSource := range templateDataSource.DataSources {
		if dataSource.Name == name {
			return true
		}
	}
	return false
}

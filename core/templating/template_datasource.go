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

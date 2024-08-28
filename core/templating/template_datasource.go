package templating

import (
	"sync"
)

type TemplateDataSource struct {
	dataSources map[string]*DataSource
	rwMutex     sync.RWMutex
}

func NewTemplateDataSource() *TemplateDataSource {

	return &TemplateDataSource{
		dataSources: make(map[string]*DataSource),
	}
}

func (t *TemplateDataSource) SetDataSource(dataSourceName string, dataSource *DataSource) {

	t.rwMutex.Lock()
	defer t.rwMutex.Unlock()

	t.dataSources[dataSourceName] = dataSource
}

func (t *TemplateDataSource) DeleteDataSource(dataSourceName string) {

	t.rwMutex.Lock()
	defer t.rwMutex.Unlock()

	delete(t.dataSources, dataSourceName)
}

func (t *TemplateDataSource) GetAllDataSources() map[string]*DataSource {

	t.rwMutex.RLock()
	defer t.rwMutex.RUnlock()

	return t.dataSources
}

func (t *TemplateDataSource) GetDataSource(name string) (*DataSource, bool) {
	t.rwMutex.RLock()
	defer t.rwMutex.RUnlock()

	source, exits := t.dataSources[name]
	return source, exits
}

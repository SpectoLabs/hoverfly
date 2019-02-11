package cache

// FastCache - cache interface used to store and retrieve any data type which requires no serialization
type FastCache interface {
	Set(key, value interface{}) error
	Get(key interface{}) (interface{}, bool)
	GetAllEntries() (map[interface{}]interface{}, error)
	RecordsCount() (int, error)
	DeleteData() error
}



package hoverfly

// Metadata - interface to store and retrieve any metadata that is related to Hoverfly
type Metadata interface {
	Set(key, value string) error
	Get(key string) (string, error)
	Delete(key string) error
	GetAll() (map[string]string, error)
	DeleteData() error
	CloseDB()
}

// NewMetadata - default metadata store
func NewMetadata(cache Cache) *Meta {
	return &Meta{
		DS: cache,
	}
}

// Meta - metadata backend that uses Cache interface
type Meta struct {
	DS Cache
}

// CloseDB - closes database
func (m *Meta) CloseDB() {
	m.DS.CloseDB()
}

// Set - saves given key and value pair to BoltDB
func (m *Meta) Set(key, value string) error {
	return m.DS.Set([]byte(key), []byte(value))
}

// Get - gets value for given key
func (m *Meta) Get(key string) (value string, err error) {
	val, err := m.DS.Get([]byte(key))
	if err != nil {
		return "", nil
	}
	return string(val), err
}

// GetAll - returns all key/value pairs
func (m *Meta) GetAll() (map[string]string, error) {
	entries, err := m.DS.GetAllEntries()
	newEntries := make(map[string]string)
	if err != nil {
		return newEntries, err
	}
	for k, v := range entries {
		newEntries[k] = string(v)
	}
	return newEntries, nil
}

// Delete - deletes given metadata key
func (m *Meta) Delete(key string) error {
	return m.DS.Delete([]byte(key))
}

// DeleteData - deletes bucket with all saved data
func (m *Meta) DeleteData() error {
	return m.DS.DeleteData()
}

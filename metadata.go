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

// NewBoltDBMetadata - default metadata store
func NewBoltDBMetadata(cache Cache) *BoltMeta {
	return &BoltMeta{
		DS: cache,
	}
}

// BoltMeta - metadata backend that uses BoltDB
type BoltMeta struct {
	DS             Cache
	MetadataBucket []byte
}

// CloseDB - closes database
func (m *BoltMeta) CloseDB() {
	m.DS.CloseDB()
}

// Set - saves given key and value pair to BoltDB
func (m *BoltMeta) Set(key, value string) error {
	return m.DS.Set([]byte(key), []byte(value))
}

// Get - gets value for given key
func (m *BoltMeta) Get(key string) (value string, err error) {
	val, err := m.DS.Get([]byte(key))
	if err != nil {
		return "", nil
	}
	return string(val), err
}

// GetAll - returns all key/value pairs
func (m *BoltMeta) GetAll() (map[string]string, error) {
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
func (m *BoltMeta) Delete(key string) error {
	return m.DS.Delete([]byte(key))
}

// DeleteData - deletes bucket with all saved data
func (m *BoltMeta) DeleteData() error {
	return m.DS.DeleteData()
}

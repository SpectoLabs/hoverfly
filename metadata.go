package hoverfly

import (
	"bytes"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
)

// Metadata - interface to store and retrieve any metadata that is related to Hoverfly
type Metadata interface {
	Set(key, value []byte) error
	Get(key []byte) ([]byte, error)
	Delete(key []byte) error
	GetAll() ([]MetaObject, error)
	DeleteData() error
	CloseDB()
}

// NewBoltDBMetadata - default metadata store
func NewBoltDBMetadata(db *bolt.DB, bucket []byte) *BoltMeta {
	return &BoltMeta{
		DS:             db,
		MetadataBucket: []byte(bucket),
	}
}

// MetadataBucketName - default bucket name for storing metadata in boltdb
const MetadataBucketName = "metadataBucket"

// BoltMeta - metadata backend that uses BoltDB
type BoltMeta struct {
	DS             *bolt.DB
	MetadataBucket []byte
}

// CloseDB - closes database
func (m *BoltMeta) CloseDB() {
	m.DS.Close()
}

// Set - saves given key and value pair to BoltDB
func (m *BoltMeta) Set(key, value []byte) error {
	err := m.DS.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(m.MetadataBucket)
		if err != nil {
			return err
		}
		err = bucket.Put(key, value)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

// Get - gets value for given key
func (m *BoltMeta) Get(key []byte) (value []byte, err error) {
	err = m.DS.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(m.MetadataBucket)
		if bucket == nil {
			return fmt.Errorf("Bucket %q not found!", m.MetadataBucket)
		}
		var buffer bytes.Buffer
		val := bucket.Get(key)

		// If it doesn't exist then it will return nil
		if val == nil {
			return fmt.Errorf("key %q not found \n", key)
		}

		buffer.Write(val)
		value = buffer.Bytes()
		return nil
	})

	return
}

// MetaObject - container to store both keys and values of captured objects
type MetaObject struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// GetAll - returns all key/value pairs
func (m *BoltMeta) GetAll() (objects []MetaObject, err error) {
	err = m.DS.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(m.MetadataBucket)
		if b == nil {
			// bucket doesn't exist
			return nil
		}
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			obj := &MetaObject{Key: string(k), Value: string(v)}
			objects = append(objects, *obj)
		}
		return nil
	})
	return
}

// Delete - deletes given metadata key
func (m *BoltMeta) Delete(key []byte) error {
	err := m.DS.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(m.MetadataBucket)
		if err != nil {
			return err
		}
		err = bucket.Delete(key)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

// DeleteData - deletes bucket with all saved data
func (m *BoltMeta) DeleteData() error {
	err := m.deleteBucket(m.MetadataBucket)
	return err
}

// DeleteBucket - deletes bucket with all saved data
func (m *BoltMeta) deleteBucket(name []byte) (err error) {
	err = m.DS.Update(func(tx *bolt.Tx) error {
		err = tx.DeleteBucket(name)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
				"name":  string(name),
			}).Warning("Failed to delete bucket")

		}
		return err
	})
	return
}

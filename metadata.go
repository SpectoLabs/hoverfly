package hoverfly

import (
	"bytes"
	"fmt"

	"github.com/boltdb/bolt"
)

// Metadata - interface to store and retrieve any metadata that is related to Hoverfly
type Metadata interface {
	Set(key, value []byte) error
	Get(key []byte) ([]byte, error)
	Delete(key []byte) error
	GetAll() ([]MetaObject, error)
	CloseDB()
}

// NewBoltDBMetadata - default metadata store
func NewBoltDBMetadata(db *bolt.DB, bucket []byte) *BoltCache {
	return &BoltCache{
		DS:             db,
		RequestsBucket: []byte(bucket),
	}
}

const MetadataBucketName = []byte("metadataBucket")

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

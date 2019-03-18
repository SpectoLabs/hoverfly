package cache

import (
	"bytes"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/boltdb/bolt"
)

// NewBoltDBCache - returns new BoltCache instance
func NewBoltDBCache(db *bolt.DB, bucket []byte) *BoltCache {
	return &BoltCache{
		DS:            db,
		CurrentBucket: []byte(bucket),
	}
}

// BoltCache - container to implement Cache instance with BoltDB backend for storage
type BoltCache struct {
	DS            *bolt.DB
	CurrentBucket []byte
}

// Set - saves given key and value pair to cache
func (c *BoltCache) Set(key, value []byte) error {
	err := c.DS.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(c.CurrentBucket)
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

// Get - searches for given key in the cache and returns value if found
func (c *BoltCache) Get(key []byte) (value []byte, err error) {

	err = c.DS.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(c.CurrentBucket)
		if bucket == nil {
			return fmt.Errorf("Bucket %q not found!", c.CurrentBucket)
		}
		// "Byte slices returned from Bolt are only valid during a transaction."
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

// GetAllValues - returns all values
func (c *BoltCache) GetAllValues() (values [][]byte, err error) {
	err = c.DS.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(c.CurrentBucket)
		if b == nil {
			// bucket doesn't exist
			return nil
		}
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var buffer bytes.Buffer
			buffer.Write(v)
			values = append(values, buffer.Bytes())
		}
		return nil
	})
	return
}

// GetAllEntries - returns all keys/values
func (c *BoltCache) GetAllEntries() (values map[string][]byte, err error) {
	err = c.DS.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(c.CurrentBucket)
		if b == nil {
			// bucket doesn't exist
			return nil
		}
		c := b.Cursor()

		values = make(map[string][]byte)

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var buffer bytes.Buffer
			buffer.Write(v)
			values[string(k)] = buffer.Bytes()
		}
		return nil
	})
	return
}

// RecordsCount - returns records count
func (c *BoltCache) RecordsCount() (count int, err error) {
	err = c.DS.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(c.CurrentBucket)
		if b == nil {
			// bucket doesn't exist
			return nil
		}

		count = b.Stats().KeyN

		return nil
	})
	return
}

// Delete - deletes specified key
func (c *BoltCache) Delete(key []byte) error {
	err := c.DS.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(c.CurrentBucket)
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
func (c *BoltCache) DeleteData() error {
	err := c.DeleteBucket(c.CurrentBucket)
	return err
}

// DeleteBucket - deletes bucket with all saved data
func (c *BoltCache) DeleteBucket(name []byte) (err error) {
	err = c.DS.Update(func(tx *bolt.Tx) error {
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

// GetAllKeys - gets all current keys
func (c *BoltCache) GetAllKeys() (keys map[string]bool, err error) {
	err = c.DS.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(c.CurrentBucket)

		keys = make(map[string]bool)

		if b == nil {
			// bucket doesn't exist
			return nil
		}
		c := b.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			keys[string(k)] = true
		}
		return nil
	})
	return
}

// GetDB - returns open BoltDB database with read/write permissions or goes down in flames if
// something bad happends
func GetDB(name string) *bolt.DB {
	log.WithFields(log.Fields{
		"databaseName": name,
	}).Info("Initiating database")
	db, err := bolt.Open(name, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

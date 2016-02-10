package hoverfly

import (
	"bytes"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
)

// RequestsBucketName - default name for BoltDB bucket
const RequestsBucketName = "rqbucket"

// Cache - provides access to BoltDB and holds current bucket name
type Cache struct {
	DS             *bolt.DB
	RequestsBucket []byte
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

// Set - saves given key and value pair to cache
func (c *Cache) Set(key, value []byte) error {
	err := c.DS.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(c.RequestsBucket)
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
func (c *Cache) Get(key []byte) (value []byte, err error) {

	err = c.DS.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(c.RequestsBucket)
		if bucket == nil {
			return fmt.Errorf("Bucket %q not found!", c.RequestsBucket)
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

// GetAllRequests - returns all captured requests/responses
func (c *Cache) GetAllRequests() (payloads []Payload, err error) {
	err = c.DS.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(c.RequestsBucket)
		if b == nil {
			// bucket doesn't exist
			return nil
		}
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			pl, err := decodePayload(v)
			if err != nil {
				log.WithFields(log.Fields{
					"error": err.Error(),
					"json":  v,
				}).Warning("Failed to deserialize bytes to payload.")
			} else {
				payloads = append(payloads, *pl)
			}
		}
		return nil
	})
	return
}

// RecordsCount - returns records count
func (c *Cache) RecordsCount() (count int, err error) {
	err = c.DS.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(c.RequestsBucket)
		if b == nil {
			// bucket doesn't exist
			return nil
		}

		count = b.Stats().KeyN

		return nil
	})
	return
}

// DeleteBucket - deletes bucket with all saved data
func (c *Cache) DeleteBucket(name []byte) (err error) {
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
func (c *Cache) GetAllKeys() (keys map[string]bool, err error) {
	err = c.DS.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(c.RequestsBucket)

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

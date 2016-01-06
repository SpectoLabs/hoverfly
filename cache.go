// redis cache backend implementation for persistent storage

package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
	"github.com/garyburd/redigo/redis"
)

const prefix = "genproxy:"
const requestsBucketName = "rqbucket"

type Cache struct {
	pool   *redis.Pool
	prefix string

	db             *bolt.DB
	requestsBucket []byte
}

// set records a key in cache (redis)
func (c *Cache) set(key string, value interface{}) error {

	if c.prefix == "" {
		c.prefix = prefix
	}

	client := c.pool.Get()
	defer client.Close()

	_, err := client.Do("SET", fmt.Sprintf(c.prefix+key), value)

	if err != nil {
		// getting payload for debugging
		var payload string
		switch v := value.(type) {
		case bytes.Buffer: // return as is
			payload = v.String() // Here v is of type bytes.Buffer
		}
		log.WithFields(log.Fields{
			"error":   err.Error(),
			"key":     fmt.Sprintf(c.prefix + key),
			"payload": payload,
		}).Error("Failed to SET key...")
	} else {
		log.WithFields(log.Fields{
			"key": fmt.Sprintf(c.prefix + key),
		}).Info("Key/value SET successfuly!")
	}

	return err
}

func (c *Cache) Set(key, value []byte) error {
	err := c.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(c.requestsBucket)
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
// getAllKeys returns all keys for specified (or default) prefix
func (c *Cache) getAllKeys() ([]string, error) {

	client := c.pool.Get()
	defer client.Close()

	values, err := redis.Strings(client.Do("KEYS", fmt.Sprintf(c.prefix+"*")))

	return values, err
}

// getAllValues returns values for specified keys
func (c *Cache) getAllValues(keys []string) ([]string, error) {

	log.WithFields(log.Fields{
		"keys": keys,
	}).Debug("Getting all supplied values")

	client := c.pool.Get()
	defer client.Close()

	// preparing keys
	var args []interface{}
	for _, k := range keys {
		args = append(args, k)
	}

	jsonStr, err := redis.Strings(client.Do("MGET", args...))

	log.WithFields(log.Fields{
		"keys": keys,
	}).Debug("Returning supplied values")

	return jsonStr, err

}

// get returns key from cache
func (c *Cache) get(key string) ([]byte, error) {

	client := c.pool.Get()
	defer client.Close()

	value, err := redis.Bytes(client.Do("GET", fmt.Sprintf(c.prefix+key)))

	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
			"key":   fmt.Sprintf(c.prefix + key),
		}).Error("Failed to GET key...")
	} else {
		log.WithFields(log.Fields{
			"key": fmt.Sprintf(c.prefix + key),
		}).Debug("Key found!")
	}

	return value, err
}

// delete removes specified entry from cache
func (c *Cache) delete(key string) error {
	if c.prefix == "" {
		c.prefix = prefix
	}

	client := c.pool.Get()
	defer client.Close()

	_, err := client.Do("DEL", key)

	return err

}

// getRedisPool returns thread safe Redis connection pool
func getRedisPool() *redis.Pool {

	// getting redis connection
	maxConnections := 10
	mc := os.Getenv("MaxConnections")
	if mc != "" {
		maxCons, err := strconv.Atoi(mc)
		if err != nil {
			maxConnections = 10
		} else {
			maxConnections = maxCons
		}
	}

	// getting redis client for state storing
	redisPool := redis.NewPool(func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", AppConfig.redisAddress)

		if err != nil {
			log.WithFields(log.Fields{
				"Error": err.Error(),
			}).Warn("Failed to create Redis connection pool!")
			return nil, err
		}
		if AppConfig.redisPassword != "" {
			if _, err := c.Do("AUTH", AppConfig.redisPassword); err != nil {
				log.WithFields(log.Fields{
					"Error":        err.Error(),
					"PasswordUsed": AppConfig.redisPassword,
				}).Warn("Failed to authenticate to Redis!")
				c.Close()
				return nil, err
			} else {
				log.Debug("Authenticated to Redis successfully! ")
			}
		}

		return c, err
	}, maxConnections)

	return redisPool
}

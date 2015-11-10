package main

import (
	"os"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/garyburd/redigo/redis"
)

type Cache struct {
	pool *redis.Pool
}

// set records a key in cache (redis)
func (c *Cache) set(key string, value []byte) error {
	client := c.pool.Get()
	defer client.Close()

	_, err := client.Do("SET", key, value)

	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Failed to record request...")
	} else {
		log.WithFields(log.Fields{}).Info("Request recorded!")
	}
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
			log.WithFields(log.Fields{"Error": err.Error()}).Panic("Failed to create Redis connection pool!")
			return nil, err
		}
		if AppConfig.redisPassword != "" {
			if _, err := c.Do("AUTH", AppConfig.redisPassword); err != nil {
				log.WithFields(log.Fields{
					"Error":        err.Error(),
					"PasswordUsed": AppConfig.redisPassword,
				}).Panic("Failed to authenticate to Redis!")
				c.Close()
				return nil, err
			} else {
				log.Info("Authenticated to Redis successfully! ")
			}
		}

		return c, err
	}, maxConnections)

	defer redisPool.Close()

	return redisPool
}

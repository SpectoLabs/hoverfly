package main

import (
	"os"
	"strconv"
	"net/http"

	"github.com/garyburd/redigo/redis"
	log "github.com/Sirupsen/logrus"
)

type DBClient struct {
	pool *redis.Pool
}

func (d *DBClient) recordRequest(r *http.Request) {
	log.Debug("Recording request")
}

func (d *DBClient) getResponse(r *http.Request) *http.Response {
	log.Debug("Returning response")
	return nil
}


// getRedisPool returns thread safe Redis connection pool
func getRedisPool() *redis.Pool {

	// getting redis connection
	maxConnections := 10
	mc := os.Getenv("MaxConnections")
	if (mc != "") {
		maxCons, err := strconv.Atoi(mc)
		if (err != nil) {
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
					"Error": err.Error(),
					"PasswordUsed": AppConfig.redisPassword,
				}).Panic("Failed to authenticate to Redis!")
				c.Close()
				return nil, err
			} else {
				log.Info("Authenticated to Redis successfully! ")
			}
		}

		return c, err

		return c, err
	}, maxConnections)

	defer redisPool.Close()

	return redisPool
}

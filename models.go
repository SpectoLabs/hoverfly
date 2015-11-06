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
	http *http.Client
}

// res structure hold response body from external service, body is not decoded and is supposed
// to be bytes, however headers should provide all required information for later decoding
// by the client.
type res struct {
	Status  int `json:"status"`
	Body    []byte `json:"body"`
	Headers map[string]string `json:"headers"`
}

// recordRequest saves request for later playback
func (d *DBClient) recordRequest(req *http.Request) (*http.Response, error) {
	log.Info("Recording request")

	// forwarding request
	resp, err := d.doRequest(req)

	// record request here

	// return new response or error here
	return resp, err

	//	c := d.pool.Get()
	//	defer c.Close()
	//
	//	_, err := c.Do("SET", r.URL.Path, r.Body)
	//
	//	if err != nil {
	//		log.WithFields(log.Fields{
	//			"error": err.Error(),
	//		}).Error("Failed to record request...")
	//	} else {
	//		log.WithFields(log.Fields{
	//		}).Info("Request recorded!")
	//	}
}


func (d *DBClient) getResponse(r *http.Request) *http.Response {
	log.Info("Returning response")
	return nil
}

// doRequest performs original request and returns response that should be returned to client and error (if there is one)
func (d *DBClient) doRequest(request *http.Request) (*http.Response, error) {
	// We can't have this set. And it only contains "/pkg/net/http/" anyway
	request.RequestURI = ""

	resp, err := d.http.Do(request)

	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
			"host": request.Host,
			"method": request.Method,
			"path": request.URL.Path,
		}).Error("Could not forward request.")
		return nil, err
	}

	log.WithFields(log.Fields{
	}).Info("Request forwarded!")

	resp.Header.Set("Gen-proxy", "was here")
	return resp, nil


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

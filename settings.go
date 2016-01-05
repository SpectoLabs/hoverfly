package main

import (
	"os"
)

// Initial structure of configuration
type Configuration struct {
	redisAddress   string
	redisPassword  string
	adminInterface string
	cachePrefix    string
	mode           string
	destination    string
	middleware     string
}

// AppCondig stores application configuration
var AppConfig Configuration

func initSettings() {
	// getting redis connection
	redisAddress := os.Getenv("RedisAddress")
	if redisAddress == "" {
		redisAddress = "127.0.0.1:6379"
	}
	AppConfig.redisAddress = redisAddress
	// getting redis password
	AppConfig.redisPassword = os.Getenv("RedisPassword")

	// admin interface port
	AppConfig.adminInterface = ":8888"

	// cache prefix
	AppConfig.cachePrefix = "genproxy:"

	// getting destination information
	//	AppConfig.destination = "get this from cache"

	// proxy state
	// should be taken from cache if we want to make it horizontally scalable (currently not needed)

	// middleware configuration
	AppConfig.middleware = os.Getenv("HoverflyMiddleware")

}

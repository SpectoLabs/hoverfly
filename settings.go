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
	recordState    bool
}

// AppCondig stores application configuration
var AppConfig Configuration

func initSettings() {
	// getting redis connection
	redisAddress := os.Getenv("RedisAddress")
	if redisAddress == "" {
		redisAddress = ":6379"
	}
	AppConfig.redisAddress = redisAddress
	// getting redis password
	AppConfig.redisPassword = os.Getenv("RedisPassword")

	// admin interface port
	AppConfig.adminInterface = ":8888"

	// cache prefix
	AppConfig.cachePrefix = "genproxy:"

	// proxy state
	// should be taken from cache
}

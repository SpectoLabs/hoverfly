package main

import (
	"os"
)
// Initial structure of configuration
type Configuration struct {
	redisAddress  string
	redisPassword string
}

// AppCondig stores application configuration
var AppConfig Configuration

func init() {
	// getting redis connection
	redisAddress := os.Getenv("RedisAddress")
	if (redisAddress == "") {
		redisAddress = ":6379"
	}
	AppConfig.redisAddress = redisAddress
	// getting redis password
	AppConfig.redisPassword = os.Getenv("RedisPassword")
}
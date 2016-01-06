package main

import (
	"os"
)

// Initial structure of configuration
type Configuration struct {
	adminInterface string
	mode           string
	destination    string
	middleware     string
	databaseName   string
}

// AppCondig stores application configuration
var AppConfig Configuration

func initSettings() {
	// admin interface port
	AppConfig.adminInterface = ":8888"

	databaseName := os.Getenv("HOVERFLY_DB")
	if databaseName == "" {
		databaseName = "requests.db"
	}

	// getting destination information
	//	AppConfig.destination = "get this from cache"

	// proxy state
	// should be taken from cache if we want to make it horizontally scalable (currently not needed)

	// middleware configuration
	AppConfig.middleware = os.Getenv("HoverflyMiddleware")

}

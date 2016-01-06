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
	AppConfig.databaseName = databaseName

	// middleware configuration
	AppConfig.middleware = os.Getenv("HoverflyMiddleware")

}

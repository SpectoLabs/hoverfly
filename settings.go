package main

import (
	"fmt"
	"os"
)

// Configuration - initial structure of configuration
type Configuration struct {
	adminInterface string
	proxyPort      string
	mode           string
	destination    string
	middleware     string
	databaseName   string
	verbose        bool
}

const DefaultPort = ":8500"      // default proxy port
const DefaultAdminPort = ":8888" // default admin interface port

// initSettings gets and returns initial configuration from env
// variables or sets defaults
func InitSettings() (AppConfig Configuration) {

	// getting default admin interface port
	adminPort := os.Getenv("AdminPort")
	if adminPort == "" {
		adminPort = DefaultAdminPort
	} else {
		adminPort = fmt.Sprintf(":%s", adminPort)
	}
	AppConfig.adminInterface = adminPort

	// getting default database
	port := os.Getenv("ProxyPort")
	if port == "" {
		port = DefaultPort
	} else {
		port = fmt.Sprintf(":%s", port)
	}

	AppConfig.proxyPort = port

	databaseName := os.Getenv("HoverflyDB")
	if databaseName == "" {
		databaseName = "requests.db"
	}
	AppConfig.databaseName = databaseName

	// middleware configuration
	AppConfig.middleware = os.Getenv("HoverflyMiddleware")

	return AppConfig
}

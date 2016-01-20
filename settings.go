package main

import (
	"os"
	"sync"
)

// Configuration - initial structure of configuration
type Configuration struct {
	adminPort    string
	proxyPort    string
	mode         string
	destination  string
	middleware   string
	databaseName string
	verbose      bool

	mu sync.Mutex
}

// SetMode - provides safe way to set new mode
func (c *Configuration) SetMode(mode string) {
	c.mu.Lock()
	c.mode = mode
	c.mu.Unlock()
}

// GetMode - provides safe way to get current mode
func (c *Configuration) GetMode() (mode string) {
	c.mu.Lock()
	mode = c.mode
	c.mu.Unlock()
	return
}

// DefaultPort - default proxy port
const DefaultPort = "8500"

// DefaultAdminPort - default admin interface port
const DefaultAdminPort = "8888"

// DefaultDatabaseName - default database name that will be created
// or used by Hoverfly
const DefaultDatabaseName = "requests.db"

// InitSettings gets and returns initial configuration from env
// variables or sets defaults
func InitSettings() *Configuration {

	var appConfig Configuration
	// getting default admin interface port
	if os.Getenv("AdminPort") != "" {
		appConfig.adminPort = os.Getenv("AdminPort")
	} else {
		appConfig.adminPort = DefaultAdminPort
	}

	// getting proxy port
	if os.Getenv("ProxyPort") != "" {
		appConfig.proxyPort = os.Getenv("ProxyPort")
	} else {
		appConfig.proxyPort = DefaultPort
	}

	databaseName := os.Getenv("HoverflyDB")
	if databaseName == "" {
		databaseName = DefaultDatabaseName
	}
	appConfig.databaseName = databaseName

	// middleware configuration
	appConfig.middleware = os.Getenv("HoverflyMiddleware")

	return &appConfig
}

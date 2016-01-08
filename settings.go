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

func (c *Configuration) SetMode(mode string) {
	c.mu.Lock()
	c.mode = mode
	c.mu.Unlock()
}

func (c *Configuration) GetMode() (mode string) {
	c.mu.Lock()
	mode = c.mode
	c.mu.Unlock()
	return
}

const DefaultPort = "8500"      // default proxy port
const DefaultAdminPort = "8888" // default admin interface port

// initSettings gets and returns initial configuration from env
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
		databaseName = "requests.db"
	}
	appConfig.databaseName = databaseName

	// middleware configuration
	appConfig.middleware = os.Getenv("HoverflyMiddleware")

	return &appConfig
}

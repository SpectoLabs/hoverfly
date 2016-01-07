package main

import (
	"fmt"
	"os"
	"sync"
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

const DefaultPort = ":8500"      // default proxy port
const DefaultAdminPort = ":8888" // default admin interface port

// initSettings gets and returns initial configuration from env
// variables or sets defaults
func InitSettings() *Configuration {

	var appConfig Configuration
	// getting default admin interface port
	adminPort := os.Getenv("AdminPort")
	if adminPort == "" {
		adminPort = DefaultAdminPort
	} else {
		adminPort = fmt.Sprintf(":%s", adminPort)
	}
	appConfig.adminInterface = adminPort

	// getting default database
	port := os.Getenv("ProxyPort")
	if port == "" {
		port = DefaultPort
	} else {
		port = fmt.Sprintf(":%s", port)
	}

	appConfig.proxyPort = port

	databaseName := os.Getenv("HoverflyDB")
	if databaseName == "" {
		databaseName = "requests.db"
	}
	appConfig.databaseName = databaseName

	// middleware configuration
	appConfig.middleware = os.Getenv("HoverflyMiddleware")

	return &appConfig
}

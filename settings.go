package hoverfly

import (
	"os"
	"strconv"
	"sync"

	log "github.com/Sirupsen/logrus"
)

// Configuration - initial structure of configuration
type Configuration struct {
	AdminPort          string
	ProxyPort          string
	Mode               string
	Destination        string
	Middleware         string
	DatabaseName       string
	Verbose            bool
	Development        bool
	SecretKey          []byte
	JWTExpirationDelta int

	mu sync.Mutex
}

// SetMode - provides safe way to set new mode
func (c *Configuration) SetMode(mode string) {
	c.mu.Lock()
	c.Mode = mode
	c.mu.Unlock()
}

// GetMode - provides safe way to get current mode
func (c *Configuration) GetMode() (mode string) {
	c.mu.Lock()
	mode = c.Mode
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

const DefaultJWTExpirationDelta = 72

// InitSettings gets and returns initial configuration from env
// variables or sets defaults
func InitSettings() *Configuration {

	var appConfig Configuration
	// getting default admin interface port
	if os.Getenv("AdminPort") != "" {
		appConfig.AdminPort = os.Getenv("AdminPort")
	} else {
		appConfig.AdminPort = DefaultAdminPort
	}

	// getting proxy port
	if os.Getenv("ProxyPort") != "" {
		appConfig.ProxyPort = os.Getenv("ProxyPort")
	} else {
		appConfig.ProxyPort = DefaultPort
	}

	databaseName := os.Getenv("HoverflyDB")
	if databaseName == "" {
		databaseName = DefaultDatabaseName
	}
	appConfig.DatabaseName = databaseName

	if os.Getenv("HoverflySecret") != "" {
		appConfig.SecretKey = []byte(os.Getenv("HoverflySecret"))
	} else {
		appConfig.SecretKey = GetRandomName(10)
	}

	if os.Getenv("HoverflyTokenExpiration") != "" {

		exp, err := strconv.Atoi(os.Getenv("HoverflyTokenExpiration"))
		if err != nil {
			log.WithFields(log.Fields{
				"error":                   err.Error(),
				"HoverflyTokenExpiration": os.Getenv("HoverflyTokenExpiration"),
			}).Error("failed to get token exipration delta, using default value")
			exp = DefaultJWTExpirationDelta
		}
		appConfig.JWTExpirationDelta = exp

	} else {
		appConfig.JWTExpirationDelta = DefaultJWTExpirationDelta
	}

	// middleware configuration
	appConfig.Middleware = os.Getenv("HoverflyMiddleware")

	return &appConfig
}

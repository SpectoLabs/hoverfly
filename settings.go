package hoverfly

import (
	"os"
	"strconv"
	"sync"

	log "github.com/Sirupsen/logrus"
)

// Configuration - initial structure of configuration
type Configuration struct {
	AdminPort    string
	ProxyPort    string
	Mode         string
	Destination  string
	Middleware   string
	DatabaseName string

	Verbose     bool
	Development bool

	SecretKey          []byte
	JWTExpirationDelta int
	AuthEnabled        bool

	ProxyControlWG sync.WaitGroup

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

// DefaultJWTExpirationDelta - default token expiration if environment variable is no provided
const DefaultJWTExpirationDelta = 72

// Environment variables
const (
	HoverflyAuthEnabledEV     = "HoverflyAuthEnabled"
	HoverflySecretEV          = "HoverflySecret"
	HoverflyTokenExpirationEV = "HoverflyTokenExpiration"

	HoverflyAdminPortEV = "AdminPort"
	HoverflyProxyPortEV = "ProxyPort"

	HoverflyDBEV         = "HoverflyDB"
	HoverflyMiddlewareEV = "HoverflyMiddleware"
)

// InitSettings gets and returns initial configuration from env
// variables or sets defaults
func InitSettings() *Configuration {

	var appConfig Configuration
	// getting default admin interface port
	if os.Getenv(HoverflyAdminPortEV) != "" {
		appConfig.AdminPort = os.Getenv(HoverflyAdminPortEV)
	} else {
		appConfig.AdminPort = DefaultAdminPort
	}

	// getting proxy port
	if os.Getenv(HoverflyProxyPortEV) != "" {
		appConfig.ProxyPort = os.Getenv(HoverflyProxyPortEV)
	} else {
		appConfig.ProxyPort = DefaultPort
	}

	databaseName := os.Getenv(HoverflyDBEV)
	if databaseName == "" {
		databaseName = DefaultDatabaseName
	}
	appConfig.DatabaseName = databaseName

	if os.Getenv(HoverflySecretEV) != "" {
		appConfig.SecretKey = []byte(os.Getenv(HoverflySecretEV))
	} else {
		appConfig.SecretKey = GetRandomName(10)
	}

	if os.Getenv(HoverflyTokenExpirationEV) != "" {

		exp, err := strconv.Atoi(os.Getenv(HoverflyTokenExpirationEV))
		if err != nil {
			log.WithFields(log.Fields{
				"error":                   err.Error(),
				"HoverflyTokenExpiration": os.Getenv(HoverflyTokenExpirationEV),
			}).Error("failed to get token exipration delta, using default value")
			exp = DefaultJWTExpirationDelta
		}
		appConfig.JWTExpirationDelta = exp

	} else {
		appConfig.JWTExpirationDelta = DefaultJWTExpirationDelta
	}

	if os.Getenv(HoverflyAuthEnabledEV) == "true" {
		appConfig.AuthEnabled = true
	} else {
		appConfig.AuthEnabled = false
	}

	// middleware configuration
	appConfig.Middleware = os.Getenv(HoverflyMiddlewareEV)

	return &appConfig
}

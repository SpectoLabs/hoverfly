package hoverfly

import (
	"os"
	"strconv"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/core/models"
	"regexp"
)

// Configuration - initial structure of configuration
type Configuration struct {
	AdminPort    string
	ProxyPort    string
	Mode         string
	Destination  string
	Middleware   string
	DatabasePath string
	Webserver    bool

	TLSVerification bool

	ResponseDelays []models.ResponseDelay

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
func (c *Configuration) GetMode() (string) {
	c.mu.Lock()
	mode := c.Mode
	c.mu.Unlock()
	return mode
}

func (c *Configuration) GetDelay(host string) (delay *models.ResponseDelay) {
	for _, val := range c.ResponseDelays {
		match := regexp.MustCompile(val.HostPattern).MatchString(host)
		if match {
			log.Info("Found response delay setting for this request host: ", delay)
			return &val
		}
	}
	return delay
}

// DefaultPort - default proxy port
const DefaultPort = "8500"

// DefaultAdminPort - default admin interface port
const DefaultAdminPort = "8888"

// DefaultDatabasePath - default database name that will be created
// or used by Hoverfly
const DefaultDatabasePath = "requests.db"

// DefaultJWTExpirationDelta - default token expiration if environment variable is no provided
const DefaultJWTExpirationDelta = 1 * 24 * 60 * 60

// Environment variables
const (
	HoverflyAuthEnabledEV     = "HoverflyAuthEnabled"
	HoverflySecretEV          = "HoverflySecret"
	HoverflyTokenExpirationEV = "HoverflyTokenExpiration"

	HoverflyAdminPortEV = "AdminPort"
	HoverflyProxyPortEV = "ProxyPort"

	HoverflyDBEV         = "HoverflyDB"
	HoverflyMiddlewareEV = "HoverflyMiddleware"

	HoverflyTLSVerification = "HoverflyTlsVerification"

	HoverflyAdminUsernameEV = "HoverflyAdmin"
	HoverflyAdminPasswordEV = "HoverflyAdminPass"

	HoverflyImportRecordsEV = "HoverflyImport"
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

	databasePath := os.Getenv(HoverflyDBEV)
	if databasePath == "" {
		appConfig.DatabasePath = DefaultDatabasePath
	}

	appConfig.Webserver = false

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
			}).Fatal("failed to get token exipration delta, using default value")
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

	if os.Getenv(HoverflyTLSVerification) == "false" {
		appConfig.TLSVerification = false
	} else {
		appConfig.TLSVerification = true
	}

	return &appConfig
}

package hoverfly

import (
	"os"
	"strconv"
	"sync"

	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/core/middleware"
)

// Configuration - initial structure of configuration
type Configuration struct {
	AdminPort    string
	ProxyPort    string
	ListenOnHost string
	Mode         string
	Destination  string
	Middleware   middleware.Middleware
	DatabasePath string
	Webserver    bool

	TLSVerification bool

	UpstreamProxy string
	PACFile       []byte

	Verbose bool

	DisableCache bool
	CacheSize int

	SecretKey          []byte
	JWTExpirationDelta int
	AuthEnabled        bool

	ProxyAuthorizationHeader string

	HttpsOnly bool

	PlainHttpTunneling bool

	ClientAuthenticationDestination string
	ClientAuthenticationClientCert  string
	ClientAuthenticationClientKey   string
	ClientAuthenticationCACert      string

	ProxyControlWG sync.WaitGroup

	mu sync.Mutex
}

// SetMode - provides safe way to set new mode
func (c *Configuration) SetMode(mode string) {
	c.mu.Lock()
	c.Mode = mode
	c.mu.Unlock()
}

func (c *Configuration) SetUpstreamProxy(upstreamProxy string) {
	if !strings.HasPrefix(upstreamProxy, "http://") && !strings.HasPrefix(upstreamProxy, "https://") {
		upstreamProxy = "http://" + upstreamProxy
	}
	c.UpstreamProxy = upstreamProxy
}

// GetMode - provides safe way to get current mode
func (c *Configuration) GetMode() string {
	c.mu.Lock()
	mode := c.Mode
	c.mu.Unlock()
	return mode
}

// DefaultPort - default proxy port
const DefaultPort = "8500"

// DefaultAdminPort - default admin interface port
const DefaultAdminPort = "8888"

const DefaultListenOnHost = "127.0.0.1"

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

	HoverflyUpstreamProxyPortEV = "UpstreamProxy"
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

	appConfig.ListenOnHost = DefaultListenOnHost

	// getting external proxy
	if os.Getenv(HoverflyUpstreamProxyPortEV) != "" {
		appConfig.UpstreamProxy = os.Getenv(HoverflyUpstreamProxyPortEV)
	} else {
		appConfig.UpstreamProxy = ""
	}

	databasePath := os.Getenv(HoverflyDBEV)
	if databasePath == "" {
		appConfig.DatabasePath = DefaultDatabasePath
	}

	appConfig.Webserver = false

	if os.Getenv(HoverflySecretEV) != "" {
		appConfig.SecretKey = []byte(os.Getenv(HoverflySecretEV))
	} else {
		appConfig.SecretKey = getRandomName(10)
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
	newMiddleware, _ := middleware.ConvertToNewMiddleware(os.Getenv(HoverflyMiddlewareEV))

	appConfig.Middleware = *newMiddleware

	if os.Getenv(HoverflyTLSVerification) == "false" {
		appConfig.TLSVerification = false
	} else {
		appConfig.TLSVerification = true
	}

	appConfig.Mode = "simulate"

	appConfig.ProxyAuthorizationHeader = "Proxy-Authorization"

	appConfig.CacheSize = 1000

	return &appConfig
}

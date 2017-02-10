package wrapper

import (
	"io/ioutil"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

type Flags []string

type Config struct {
	HoverflyHost          string `yaml:"hoverfly.host"`
	HoverflyAdminPort     string `yaml:"hoverfly.admin.port"`
	HoverflyProxyPort     string `yaml:"hoverfly.proxy.port"`
	HoverflyDbType        string `yaml:"hoverfly.db.type"`
	HoverflyUsername      string `yaml:"hoverfly.username"`
	HoverflyPassword      string `yaml:"hoverfly.password"`
	HoverflyWebserver     bool   `yaml:"hoverfly.webserver"`
	HoverflyCertificate   string `yaml:"hoverfly.tls.certificate"`
	HoverflyKey           string `yaml:"hoverfly.tls.key"`
	HoverflyDisableTls    bool   `yaml:"hoverfly.tls.disable"`
	HoverflyUpstreamProxy string `yaml:"hoverfly.upstream.proxy"`
}

func GetConfig() *Config {
	err := viper.ReadInConfig()
	if err != nil {
		log.Debug(err.Error())
	}

	return &Config{
		HoverflyHost:        viper.GetString("hoverfly.host"),
		HoverflyAdminPort:   viper.GetString("hoverfly.admin.port"),
		HoverflyProxyPort:   viper.GetString("hoverfly.proxy.port"),
		HoverflyDbType:      viper.GetString("hoverfly.db.type"),
		HoverflyUsername:    viper.GetString("hoverfly.username"),
		HoverflyPassword:    viper.GetString("hoverfly.password"),
		HoverflyWebserver:   viper.GetBool("hoverfly.webserver"),
		HoverflyCertificate: viper.GetString("hoverfly.tls.certificate"),
		HoverflyKey:         viper.GetString("hoverfly.tls.key"),
		HoverflyDisableTls:  viper.GetBool("hoverfly.tls.disable"),
	}
}

func (this *Config) SetHost(host string) *Config {
	if len(host) > 0 {
		this.HoverflyHost = host
	}
	return this
}

func (this *Config) SetAdminPort(adminPort string) *Config {
	if len(adminPort) > 0 {
		this.HoverflyAdminPort = adminPort
	}
	return this
}

func (this *Config) SetProxyPort(proxyPort string) *Config {
	if len(proxyPort) > 0 {
		this.HoverflyProxyPort = proxyPort
	}
	return this
}

func (this *Config) SetDbType(dbType string) *Config {
	if dbType == "memory" {
		this.HoverflyDbType = dbType
	}
	if dbType == "boltdb" {
		this.HoverflyDbType = dbType
	}
	return this
}

func (this *Config) SetUsername(username string) *Config {
	if len(username) > 0 {
		this.HoverflyUsername = username
	}
	return this
}

func (this *Config) SetPassword(password string) *Config {
	if len(password) > 0 {
		this.HoverflyPassword = password
	}
	return this
}

func (this *Config) SetWebserver(hoverflyType string) *Config {
	if hoverflyType == "webserver" {
		this.HoverflyWebserver = true
	}

	if hoverflyType == "proxy" {
		this.HoverflyWebserver = false
	}

	return this
}

func (this *Config) SetCertificate(certificate string) *Config {
	if len(certificate) > 0 {
		this.HoverflyCertificate = certificate
	}
	return this
}

func (this *Config) SetKey(key string) *Config {
	if len(key) > 0 {
		this.HoverflyKey = key
	}
	return this
}

func (this *Config) SetUpstreamProxy(upstreamProxy string) *Config {
	if len(upstreamProxy) > 0 {
		this.HoverflyUpstreamProxy = upstreamProxy
	}
	return this
}

func (c *Config) GetFilepath() string {
	return viper.ConfigFileUsed()
}

func (this *Config) DisableTls(disableTls bool) *Config {
	if this.HoverflyDisableTls || disableTls {
		this.HoverflyDisableTls = true
	}
	return this
}

func (c *Config) WriteToFile(hoverflyDirectory HoverflyDirectory) error {
	data, err := yaml.Marshal(c)

	if err != nil {
		log.Debug(err.Error())
		return err
	}

	filepath := filepath.Join(hoverflyDirectory.Path, "config.yaml")

	err = ioutil.WriteFile(filepath, data, 0644)

	if err != nil {
		log.Debug(err.Error())
		return err
	}

	return nil
}

func (this Config) BuildFlags() Flags {
	flags := Flags{}

	if this.HoverflyAdminPort != "" {
		flags = append(flags, "-ap="+this.HoverflyAdminPort)
	}

	if this.HoverflyProxyPort != "" {
		flags = append(flags, "-pp="+this.HoverflyProxyPort)
	}

	if this.HoverflyDbType != "" {
		flags = append(flags, "-db="+this.HoverflyDbType)
	}

	if this.HoverflyWebserver {
		flags = append(flags, "-webserver")
	}

	if this.HoverflyCertificate != "" {
		flags = append(flags, "-cert="+this.HoverflyCertificate)
	}

	if this.HoverflyKey != "" {
		flags = append(flags, "-key="+this.HoverflyKey)
	}

	if this.HoverflyDisableTls {
		flags = append(flags, "-tls-verification=false")
	}

	if this.HoverflyUpstreamProxy != "" {
		flags = append(flags, "-upstream-proxy="+this.HoverflyUpstreamProxy)
	}

	return flags
}

func SetConfigurationPaths() {
	viper.AddConfigPath("./.hoverfly")
	viper.AddConfigPath("$HOME/.hoverfly")
}

func SetConfigurationDefaults() {
	viper.SetDefault("hoverfly.host", "localhost")
	viper.SetDefault("hoverfly.admin.port", "8888")
	viper.SetDefault("hoverfly.proxy.port", "8500")
	viper.SetDefault("hoverfly.db.type", "memory")
	viper.SetDefault("hoverfly.username", "")
	viper.SetDefault("hoverfly.password", "")
	viper.SetDefault("hoverfly.webserver", false)
	viper.SetDefault("hoverfly.tls.certificate", "")
	viper.SetDefault("hoverfly.tls.key", "")
	viper.SetDefault("hoverfly.tls.disable", false)
	viper.SetDefault("hoverfly.upsream.proxy", "")
}

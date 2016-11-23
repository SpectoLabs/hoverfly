package main

import (
	"io/ioutil"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

type Config struct {
	HoverflyHost        string `yaml:"hoverfly.host"`
	HoverflyAdminPort   string `yaml:"hoverfly.admin.port"`
	HoverflyProxyPort   string `yaml:"hoverfly.proxy.port"`
	HoverflyUsername    string `yaml:"hoverfly.username"`
	HoverflyPassword    string `yaml:"hoverfly.password"`
	HoverflyWebserver   bool   `yaml:"hoverfly.webserver"`
	HoverflyCertificate string `yaml:"hoverfly.tls.certificate"`
	HoverflyKey         string `yaml:"hoverfly.tls.key"`
}

func GetConfig() *Config {
	err := viper.ReadInConfig()
	if err != nil {
		log.Debug(err.Error())
	}

	return &Config{
		HoverflyHost:      viper.GetString("hoverfly.host"),
		HoverflyAdminPort: viper.GetString("hoverfly.admin.port"),
		HoverflyProxyPort: viper.GetString("hoverfly.proxy.port"),
		HoverflyUsername:  viper.GetString("hoverfly.username"),
		HoverflyPassword:  viper.GetString("hoverfly.password"),
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

func (this *Config) SetWebserver(webserver bool) *Config {
	this.HoverflyWebserver = webserver
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

func (c *Config) GetFilepath() string {
	return viper.ConfigFileUsed()
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

func SetConfigurationPaths() {
	viper.AddConfigPath("./.hoverfly")
	viper.AddConfigPath("$HOME/.hoverfly")
}

func SetConfigurationDefaults() {
	viper.SetDefault("hoverfly.host", "localhost")
	viper.SetDefault("hoverfly.admin.port", "8888")
	viper.SetDefault("hoverfly.proxy.port", "8500")
	viper.SetDefault("hoverfly.username", "")
	viper.SetDefault("hoverfly.password", "")
	viper.SetDefault("hoverfly.webserver", "false")
	viper.SetDefault("hoverfly.tls.certificate", "")
	viper.SetDefault("hoverfly.tls.key", "")
}

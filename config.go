package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"path/filepath"
	"io/ioutil"
)

type Config struct {
	HoverflyHost      string `yaml:"hoverfly.host"`
	HoverflyAdminPort string `yaml:"hoverfly.admin.port"`
	HoverflyProxyPort string `yaml:"hoverfly.proxy.port"`
	SpectoLabAPIKey   string `yaml:"specto.lab.api.key"`
}

func GetConfig(hoverflyHostOverride, hoverflyAdminPortOverride, hoverflyProxyPortOverride string) (Config) {
	err := viper.ReadInConfig()
	if err != nil {
		log.Debug(err.Error())
	}

	config := Config{
		HoverflyHost: viper.GetString("hoverfly.host"),
		HoverflyAdminPort: viper.GetString("hoverfly.admin.port"),
		HoverflyProxyPort: viper.GetString("hoverfly.proxy.port"),
		SpectoLabAPIKey: viper.GetString("specto.lab.api.key"),
	}

	if len(hoverflyHostOverride) > 0 {
		config.HoverflyHost = hoverflyHostOverride
	}

	if len(hoverflyAdminPortOverride) > 0 {
		config.HoverflyAdminPort = hoverflyAdminPortOverride
	}

	if len(hoverflyProxyPortOverride) > 0 {
		config.HoverflyProxyPort = hoverflyProxyPortOverride
	}

	return config
}

func (c *Config) GetFilepath() (string) {
	return viper.ConfigFileUsed()
}

func (c *Config) WriteToFile(hoverflyDirectory HoverflyDirectory) (error) {
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
	viper.SetDefault("specto.lab.api.key", "")
}
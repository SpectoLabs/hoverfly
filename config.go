package main

import (
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"path/filepath"
	"io/ioutil"
)

type Config struct {
	HoverflyHost      string `yaml:"hoverfly.host"`
	HoverflyAdminPort string `yaml:"hoverfly.admin.port"`
	HoverflyProxyPort string `yaml:"hoverfly.proxy.port"`
	SpectoLabHost     string `yaml:"specto.lab.host"`
	SpectoLabPort     string `yaml:"specto.lab.port"`
}

func NewConfig() Config {
	return Config{
		HoverflyHost: viper.GetString("hoverfly.host"),
		HoverflyAdminPort: viper.GetString("hoverfly.admin.port"),
		HoverflyProxyPort: viper.GetString("hoverfly.proxy.port"),
		SpectoLabHost: viper.GetString("specto.lab.host"),
		SpectoLabPort: viper.GetString("specto.lab.port"),
	}
}

func (c *Config) WriteToFile(path string) {
	data, err := yaml.Marshal(c)

	if err != nil {
		failAndExit(err)
	}

	filepath := filepath.Join(path, "config.yaml")

	err = ioutil.WriteFile(filepath, data, 0644)

	if err != nil {
		failAndExit(err)
	}
}
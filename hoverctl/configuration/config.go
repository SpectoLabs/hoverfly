package configuration

import (
	"io/ioutil"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

type Flags []string

type Config struct {
	DefaultTarget string            `mapstructure:"default" yaml:"default"`	// viper uses mapstructure to unmarshall
	Targets       map[string]Target `yaml:"targets"`
}

func GetConfig() *Config {
	err := viper.ReadInConfig()
	if err != nil {
		log.Debug("Error reading config")
		if err.Error() == `Unsupported Config Type ""` {
			log.Debug("viper not properly configured, if this is the first time executing hoverctl, please try the command again")
		} else {
			log.Debug(err.Error())
		}
	}

	return parseConfig()
}

func parseConfig() *Config {
	config := &Config{}
	err := viper.Unmarshal(config)
	if err != nil {
		log.Debug("Error parsing config")
		log.Debug(err.Error())
	}
	if config.DefaultTarget == "" {
		config.DefaultTarget = viper.GetString("default")
	}
	if config.Targets == nil {
		config.Targets = map[string]Target{}
	}
	defaultTarget := NewDefaultTarget()

	// Initialize local target
	if config.Targets["local"] == (Target{}) {
		config.Targets["local"] = *defaultTarget
	}

	// Assume default value for any required config
	for key, target := range config.Targets {
		if target.Host == "" {
			target.Host = defaultTarget.Host
		}

		if target.AdminPort == 0 {
			target.AdminPort = defaultTarget.AdminPort
		}

		if target.ProxyPort == 0 {
			target.ProxyPort = defaultTarget.ProxyPort
		}
		config.Targets[key] = target
	}

	return config
}

func (this *Config) GetTarget(targetName string) *Target {
	if targetName == "" {
		targetName = this.DefaultTarget
	}

	for key, target := range this.Targets {
		if key == targetName {
			return &target
		}
	}

	return nil
}

func (this *Config) NewTarget(target Target) {
	this.Targets[target.Name] = target
}

func (this *Config) DeleteTarget(targetToDelete Target) {
	targets := map[string]Target{}

	for key, target := range this.Targets {
		if key != targetToDelete.Name {
			targets[key] = target
		}
	}

	this.Targets = targets
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

	configPath := filepath.Join(hoverflyDirectory.Path, "config.yaml")

	err = ioutil.WriteFile(configPath, data, 0644)

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
	viper.SetDefault("default", "local")
}

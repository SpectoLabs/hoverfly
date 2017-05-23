package configuration

import (
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

type Target struct {
	Name      string
	Host      string `yaml:"host,omitempty"`
	AdminPort int    `yaml:"admin.port,omitempty"`
	ProxyPort int    `yaml:"proxy.port,omitempty"`
	AuthToken string `yaml:"auth.token,omitempty"`
	Pid       int    `yaml:"pid,omitempty"`

	Webserver    bool   `yaml:",omitempty"`
	CachePath    string `yaml:",omitempty"`
	DisableCache bool   `yaml:",omitempty"`

	CertificatePath string `yaml:",omitempty"`
	KeyPath         string `yaml:",omitempty"`
	DisableTls      bool   `yaml:",omitempty"`

	UpstreamProxyUrl string `yaml:",omitempty"`
	HttpsOnly        bool   `yaml:",omitempty"`

	AuthEnabled bool
	Username    string
	Password    string
}

func NewDefaultTarget() *Target {
	return &Target{
		Name:      "local",
		Host:      "localhost",
		AdminPort: 8888,
		ProxyPort: 8500,
	}
}

func NewTarget(name, host string, adminPort, proxyPort int) *Target {
	target := NewDefaultTarget()
	if name != "" {
		target.Name = name
	}

	if host != "" {
		target.Host = host
	}

	if adminPort != 0 {
		target.AdminPort = adminPort
	}

	if proxyPort != 0 {
		target.ProxyPort = proxyPort
	}

	return target
}

func getTargetsFromConfig(configTargets map[string]interface{}) map[string]Target {
	targets := map[string]Target{}

	for key, target := range configTargets {
		targetMap := target.(map[interface{}]interface{})

		targetHoverfly := Target{}

		targetHoverfly.Name = key

		if targetMap["host"] != nil {
			targetHoverfly.Host = targetMap["host"].(string)
		}

		if targetMap["admin.port"] != nil {
			targetHoverfly.AdminPort = targetMap["admin.port"].(int)
		}

		if targetMap["proxy.port"] != nil {
			targetHoverfly.ProxyPort = targetMap["proxy.port"].(int)
		}

		if targetMap["auth.token"] != nil {
			targetHoverfly.AuthToken = targetMap["auth.token"].(string)
		}

		targets[key] = targetHoverfly
	}

	if targets["local"] == (Target{}) {
		localTarget := NewDefaultTarget()
		targets["local"] = *localTarget
	}

	return targets
}

func (this Target) BuildFlags() Flags {
	flags := Flags{}

	if this.AdminPort != 0 {
		flags = append(flags, "-ap="+strconv.Itoa(this.AdminPort))
	}

	if this.ProxyPort != 0 {
		flags = append(flags, "-pp="+strconv.Itoa(this.ProxyPort))
	}

	if this.Webserver {
		flags = append(flags, "-webserver")
	}

	if this.CachePath != "" {
		flags = append(flags, "-db=boltdb", "-db-path="+this.CachePath)
	}

	if this.DisableCache {
		flags = append(flags, "-disable-cache")
	}

	if this.CertificatePath != "" {
		flags = append(flags, "-cert="+this.CertificatePath)
	}

	if this.KeyPath != "" {
		flags = append(flags, "-key="+this.KeyPath)
	}

	if this.DisableTls {
		flags = append(flags, "-tls-verification=false")
	}

	if this.UpstreamProxyUrl != "" {
		flags = append(flags, "-upstream-proxy="+this.UpstreamProxyUrl)
	}

	if this.HttpsOnly {
		flags = append(flags, "-https-only")
	}

	if this.AuthEnabled {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(this.Password), 10)
		flags = append(flags, "-auth", "-username", this.Username, "-password-hash", string(hashedPassword))
	}

	return flags
}

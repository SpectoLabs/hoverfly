package configuration

import (
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

type Target struct {
	Name      string
	Host      string `yaml:"host,omitempty"`
	AdminPort int    `mapstructure:"admin.port,omitempty" yaml:"admin.port,omitempty"`
	ProxyPort int    `mapstructure:"proxy.port,omitempty" yaml:"proxy.port,omitempty"`
	AuthToken string `mapstructure:"auth.token,omitempty" yaml:"auth.token,omitempty"`
	Pid       int    `yaml:"pid,omitempty"`

	Webserver    bool   `yaml:",omitempty"`
	CachePath    string `yaml:",omitempty"`
	DisableCache bool   `yaml:",omitempty"`
	ListenOnHost string `yaml:",omitempty"`

	CertificatePath string `yaml:",omitempty"`
	KeyPath         string `yaml:",omitempty"`
	DisableTls      bool   `yaml:",omitempty"`

	UpstreamProxyUrl string `yaml:",omitempty"`
	PACFile          string `yaml:",omitempty"`
	HttpsOnly        bool   `yaml:",omitempty"`

	ClientAuthenticationDestination string `yaml:",omitempty"`
	ClientAuthenticationClientCert  string `yaml:",omitempty"`
	ClientAuthenticationClientKey   string `yaml:",omitempty"`
	ClientAuthenticationCACert      string `yaml:",omitempty"`

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

	if this.ListenOnHost != "" {
		flags = append(flags, "-listen-on-host="+this.ListenOnHost)
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

	if this.ClientAuthenticationDestination != "" {
		flags = append(flags, "-client-authentication-destination="+this.ClientAuthenticationDestination)
	}

	if this.ClientAuthenticationClientCert != "" {
		flags = append(flags, "-client-authentication-client-cert="+this.ClientAuthenticationClientCert)
	}

	if this.ClientAuthenticationClientKey != "" {
		flags = append(flags, "-client-authentication-client-key="+this.ClientAuthenticationClientKey)
	}

	if this.ClientAuthenticationCACert != "" {
		flags = append(flags, "-client-authentication-ca-cert="+this.ClientAuthenticationCACert)
	}

	return flags
}

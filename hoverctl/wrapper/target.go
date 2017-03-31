package wrapper

type Target struct {
	Name      string
	Host      string `yaml:"host"`
	AdminPort int    `yaml:"admin.port"`
	AuthToken string `yaml:"auth.token"`
	Pid       int    `yaml:"pid"`

	Webserver bool
	ProxyPort int
	CachePath string

	CertificatePath  string
	KeyPath          string
	DisableTls       bool
	UpstreamProxyUrl string
}

func NewDefaultTarget() *Target {
	return &Target{
		Name:      "default",
		Host:      "localhost",
		AdminPort: 8888,
		ProxyPort: 8500,
	}
}

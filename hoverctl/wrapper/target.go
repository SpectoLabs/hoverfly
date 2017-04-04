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

		if targetMap["auth.token"] != nil {
			targetHoverfly.AuthToken = targetMap["auth.token"].(string)
		}

		if targetMap["pid"] != nil {
			targetHoverfly.Pid = targetMap["pid"].(int)
		}

		targets[key] = targetHoverfly
	}

	return targets
}

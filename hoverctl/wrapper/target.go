package wrapper

type TargetHoverfly struct {
	Name      string
	Host      string `yaml:"host"`
	AdminPort int    `yaml:"admin.port"`
	AuthToken string `yaml:"auth.token"`
	Pid       int    `yaml:"pid"`
}

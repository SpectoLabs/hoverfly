package wrapper

import (
	"os"
	"testing"

	hf "github.com/SpectoLabs/hoverfly/core"
	"github.com/SpectoLabs/hoverfly/hoverctl/configuration"
)

var hoverfly *hf.Hoverfly

var target = configuration.Target{
	Host:      "localhost",
	AdminPort: 8500,
}

var inaccessibleTarget = configuration.Target{
	Host:      "something",
	AdminPort: 1234,
}

func TestMain(m *testing.M) {
	hoverfly = hf.NewHoverfly()
	hoverfly.Cfg.Webserver = true
	hoverfly.StartProxy()

	returnCode := m.Run()
	hoverfly.StopProxy()
	os.Exit(returnCode)
}

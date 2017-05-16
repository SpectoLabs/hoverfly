package wrapper

import (
	"os"
	"testing"

	hf "github.com/SpectoLabs/hoverfly/core"
)

var hoverfly *hf.Hoverfly

func TestMain(m *testing.M) {
	hoverfly = hf.NewHoverfly()
	hoverfly.Cfg.Webserver = true
	hoverfly.StartProxy()

	returnCode := m.Run()
	hoverfly.StopProxy()
	os.Exit(returnCode)
}

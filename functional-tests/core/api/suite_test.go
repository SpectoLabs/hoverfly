package api_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"os/exec"
	"testing"

	"github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/phayes/freeport"
)

var (
	hoverflyCmd *exec.Cmd

	adminPort = freeport.GetPort()

	proxyPort = freeport.GetPort()
)

func TestCoreAPI(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Core API Suite")
}

var _ = BeforeSuite(func() {
	functional_tests.BinaryPrefix = ".."
})

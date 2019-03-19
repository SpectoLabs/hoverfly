package hoverctl_suite

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/SpectoLabs/hoverfly/functional-tests"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

const (
	simulate          = "simulate"
	capture           = "capture"
	synthesize        = "synthesize"
	modify            = "modify"
	spy            	  = "spy"
	diff              = "diff"
	generatedTestData = "testdata-gen"
)

var (
	hoverctlBinary   string
	workingDirectory string
)

func TestHoverctlFunctionalTestSuite(t *testing.T) {
	os.Mkdir(generatedTestData, os.ModePerm)

	RegisterFailHandler(Fail)
	RunSpecs(t, "Hoverctl functional test suite")

	os.RemoveAll(generatedTestData)
}

var _ = BeforeSuite(func() {
	workingDirectory, _ = os.Getwd()

	binDirectory := filepath.Join(workingDirectory, "bin")

	os.Setenv("PATH", fmt.Sprintf("%v:%v", binDirectory, os.Getenv("PATH")))

	var err error
	hoverctlBinary, err = gexec.Build("github.com/SpectoLabs/hoverfly/hoverctl")
	Expect(err).ShouldNot(HaveOccurred())
})

var _ = BeforeEach(func() {
	WipeConfig()

})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})

func WipeConfig() {
	configPath := filepath.Join(workingDirectory, ".hoverfly", "config.yaml")
	ioutil.WriteFile(configPath, []byte(""), 0644)
}

func KillHoverflyTargets(table string) {
	targets := functional_tests.TableToSliceMapStringString(table)
	for _, target := range targets {
		adminPort, _ := strconv.Atoi(target["ADMIN PORT"])

		functional_tests.NewHoverflyWithAdminPort(adminPort).StopAPIAuthenticated(functional_tests.HoverflyUsername, functional_tests.HoverflyPassword)
	}
}

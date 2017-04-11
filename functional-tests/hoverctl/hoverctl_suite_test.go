package hoverctl_suite

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

const (
	simulate          = "simulate"
	capture           = "capture"
	synthesize        = "synthesize"
	modify            = "modify"
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
	filepath := filepath.Join(workingDirectory, ".hoverfly", "config.yaml")
	ioutil.WriteFile(filepath, []byte(""), 0644)
}

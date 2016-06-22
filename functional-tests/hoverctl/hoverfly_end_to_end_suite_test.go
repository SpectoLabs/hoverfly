package hoverfly_end_to_end_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
	"github.com/dghubble/sling"
	"strings"
	"net/http"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"os"
	"os/exec"
	"strconv"
	"time"
	"gopkg.in/yaml.v2"
)

var (
	hoverctlBinary string
	hoverctlCacheDir string
	workingDirectory string
)

func TestHoverflyEndToEnd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Hoverfly End To End Suite")
}

var _ = BeforeSuite(func() {
	workingDirectory, _ := os.Getwd()

	hoverctlCacheDir = filepath.Join(workingDirectory, ".hoverfly/cache")

	hoverctlBinary = filepath.Join(workingDirectory, "bin/hoverctl")

	binDirectory := filepath.Join(workingDirectory, "bin")

	os.Setenv("PATH", fmt.Sprintf("%v:%v", binDirectory, os.Getenv("PATH")))
})

func SetHoverflyMode(mode string, port int) {
	req := sling.New().Post(fmt.Sprintf("http://localhost:%v/api/state", port)).Body(strings.NewReader(`{"mode":"` + mode +`"}`))
	res := DoRequest(req)
	Expect(res.StatusCode).To(Equal(200))
}

func DoRequest(r *sling.Sling) (*http.Response) {
	req, err := r.Request()
	Expect(err).To(BeNil())
	response, err := http.DefaultClient.Do(req)

	Expect(err).To(BeNil())
	return response
}

func GetHoverflyMode(port int) string {
	currentState := &stateRequest{}
	resp := DoRequest(sling.New().Get(fmt.Sprintf("http://localhost:%v/api/state", port)))

	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).To(BeNil())

	err = json.Unmarshal(body, currentState)
	Expect(err).To(BeNil())

	return currentState.Mode
}

type stateRequest struct {
	Mode        string `json:"mode"`
	Destination string `json:"destination"`
}

func startHoverfly(adminPort, proxyPort int, workingDir string) * exec.Cmd {
	hoverflyBinaryUri := filepath.Join(workingDir, "bin/hoverfly")
	hoverflyCmd := exec.Command(hoverflyBinaryUri, "-db", "memory", "-ap", strconv.Itoa(adminPort), "-pp", strconv.Itoa(proxyPort))

	err := hoverflyCmd.Start()

	if err != nil {
		fmt.Println("Unable to start Hoverfly")
		fmt.Println(hoverflyBinaryUri)
		fmt.Println("Is the binary there?")
		os.Exit(1)
	}

	Eventually(func() int {
		resp, err := http.Get(fmt.Sprintf("http://localhost:%v/api/state", adminPort))
		if err == nil {
			return resp.StatusCode
		} else {
			fmt.Println(err.Error())
			return 0
		}
	}, time.Second * 3).Should(BeNumerically("==", http.StatusOK))

	return hoverflyCmd
}

type testConfig struct {
	HoverflyHost      string `yaml:"hoverfly.host"`
	HoverflyAdminPort string `yaml:"hoverfly.admin.port"`
	HoverflyProxyPort string `yaml:"hoverfly.proxy.port"`
}

func WriteConfiguration(host, adminPort, proxyPort string) {
	configHost := "localhost"
	configAdminPort := "8888"
	configProxyPort := "8500"

	if len(adminPort) > 0 {
		configHost = host
	}

	if len(adminPort) > 0 {
		configAdminPort = adminPort
	}

	if len(proxyPort) > 0 {
		configProxyPort = proxyPort
	}

	testConfig := testConfig{
		HoverflyHost:configHost,
		HoverflyAdminPort: configAdminPort,
		HoverflyProxyPort: configProxyPort,
	}

	data, _ := yaml.Marshal(testConfig)

	filepath := filepath.Join(workingDirectory, ".hoverfly", "config.yaml")

	ioutil.WriteFile(filepath, data, 0644)

}
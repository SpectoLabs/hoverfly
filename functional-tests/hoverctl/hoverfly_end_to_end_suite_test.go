package hoverctl_end_to_end

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"
)

const (
	simulate          = "simulate"
	capture           = "capture"
	synthesize        = "synthesize"
	modify            = "modify"
	generatedTestData = "testdata-gen"
)

var (
	hoverfly *functional_tests.Hoverfly

	hoverctlBinary   string
	hoverctlCacheDir string
	workingDirectory string
)

func TestHoverflyEndToEnd(t *testing.T) {
	os.Mkdir(generatedTestData, os.ModePerm)

	RegisterFailHandler(Fail)
	RunSpecs(t, "Hoverfly End To End Suite")

	os.RemoveAll(generatedTestData)
}

var _ = BeforeSuite(func() {
	workingDirectory, _ = os.Getwd()

	hoverctlCacheDir = filepath.Join(workingDirectory, ".hoverfly/cache")

	hoverctlBinary = filepath.Join(workingDirectory, "bin/hoverctl")

	binDirectory := filepath.Join(workingDirectory, "bin")

	os.Setenv("PATH", fmt.Sprintf("%v:%v", binDirectory, os.Getenv("PATH")))

})

func SetHoverflyMode(mode string, port int) {
	req := sling.New().Post(fmt.Sprintf("http://localhost:%v/api/state", port)).Body(strings.NewReader(`{"mode":"` + mode + `"}`))
	res := functional_tests.DoRequest(req)
	Expect(res.StatusCode).To(Equal(200))
}

func GetHoverflyMode(port int) string {
	currentState := &stateRequest{}
	resp := functional_tests.DoRequest(sling.New().Get(fmt.Sprintf("http://localhost:%v/api/state", port)))

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

func startHoverfly(adminPort, proxyPort int, workingDir string) *exec.Cmd {
	hoverflyBinaryUri := filepath.Join(workingDir, "bin/hoverfly")
	hoverflyCmd := exec.Command(hoverflyBinaryUri, "-ap", strconv.Itoa(adminPort), "-pp", strconv.Itoa(proxyPort))

	err := hoverflyCmd.Start()

	functional_tests.BinaryErrorCheck(err, hoverflyBinaryUri)
	functional_tests.Healthcheck(adminPort)

	return hoverflyCmd
}

func startHoverflyWithMiddleware(adminPort, proxyPort int, workingDir, binary, script string) *exec.Cmd {
	hoverflyBinaryUri := filepath.Join(workingDir, "bin/hoverfly")
	hoverflyCmd := exec.Command(hoverflyBinaryUri, "-ap", strconv.Itoa(adminPort), "-pp", strconv.Itoa(proxyPort))

	err := hoverflyCmd.Start()

	functional_tests.BinaryErrorCheck(err, hoverflyBinaryUri)
	functional_tests.Healthcheck(adminPort)

	request := sling.New().Put(fmt.Sprintf("http://localhost:%v/api/v2/hoverfly/middleware", adminPort)).BodyJSON(v2.MiddlewareView{Binary: binary, Script: script})

	functional_tests.DoRequest(request)

	return hoverflyCmd
}

func startHoverflyWithAuth(adminPort, proxyPort int, workingDir, username, password string) *exec.Cmd {
	os.Remove(filepath.Join(workingDir, "requests.db"))

	hoverflyBinaryUri := filepath.Join(workingDir, "bin/hoverfly")

	hoverflyAddUserCmd := exec.Command(hoverflyBinaryUri, "-db", "boltdb", "-add", "-username", username, "-password", password, "-ap", strconv.Itoa(adminPort), "-pp", strconv.Itoa(proxyPort))
	err := hoverflyAddUserCmd.Run()

	if err != nil {
		fmt.Println("Unable to start Hoverfly to add user")
		fmt.Println(hoverflyBinaryUri)
		fmt.Println("Is the binary there?")
		os.Exit(1)
	}

	hoverflyCmd := exec.Command(hoverflyBinaryUri, "-db", "boltdb", "-ap", strconv.Itoa(adminPort), "-pp", strconv.Itoa(proxyPort), "-auth", "true")
	err = hoverflyCmd.Start()

	functional_tests.BinaryErrorCheck(err, hoverflyBinaryUri)
	functional_tests.Healthcheck(adminPort)

	return hoverflyCmd
}

func startHoverflyWebserver(adminPort, proxyPort int, workingDir string) *exec.Cmd {
	hoverflyBinaryUri := filepath.Join(workingDir, "bin/hoverfly")
	hoverflyCmd := exec.Command(hoverflyBinaryUri, "-ap", strconv.Itoa(adminPort), "-pp", strconv.Itoa(proxyPort), "-webserver")

	err := hoverflyCmd.Start()

	functional_tests.BinaryErrorCheck(err, hoverflyBinaryUri)
	functional_tests.Healthcheck(adminPort)

	return hoverflyCmd
}

type testConfig struct {
	HoverflyHost      string `yaml:"hoverfly.host"`
	HoverflyAdminPort string `yaml:"hoverfly.admin.port"`
	HoverflyProxyPort string `yaml:"hoverfly.proxy.port"`
	HoverflyUsername  string `yaml:"hoverfly.username"`
	HoverflyPassword  string `yaml:"hoverfly.password"`
	HoverflyWebserver bool   `yaml:"hoverfly.webserver"`
}

func WriteConfiguration(host, adminPort, proxyPort string) {
	WriteConfigurationWithAuth(host, adminPort, proxyPort, false, "", "")
}

func WriteConfigurationWithAuth(host, adminPort, proxyPort string, webserver bool, username, password string) {
	configHost := "localhost"
	configAdminPort := "8888"
	configProxyPort := "8500"
	configUsername := ""
	configPassword := ""

	if len(host) > 0 {
		configHost = host
	}

	if len(adminPort) > 0 {
		configAdminPort = adminPort
	}

	if len(proxyPort) > 0 {
		configProxyPort = proxyPort
	}

	if len(username) > 0 {
		configUsername = username
	}

	if len(password) > 0 {
		configPassword = password
	}

	testConfig := testConfig{
		HoverflyHost:      configHost,
		HoverflyAdminPort: configAdminPort,
		HoverflyProxyPort: configProxyPort,
		HoverflyWebserver: webserver,
		HoverflyUsername:  configUsername,
		HoverflyPassword:  configPassword,
	}

	data, _ := yaml.Marshal(testConfig)

	filepath := filepath.Join(workingDirectory, ".hoverfly", "config.yaml")

	ioutil.WriteFile(filepath, data, 0644)

}

func generateFileName() string {

	rb := make([]byte, 6)
	rand.Read(rb)

	rs := base64.URLEncoding.EncodeToString(rb)

	return "testdata-gen/" + rs + ".json"
}

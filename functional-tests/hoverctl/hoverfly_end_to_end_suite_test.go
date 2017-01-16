package hoverctl_end_to_end

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
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
	res := DoRequest(req)
	Expect(res.StatusCode).To(Equal(200))
}

func DoRequest(r *sling.Sling) *http.Response {
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

func startHoverfly(adminPort, proxyPort int, workingDir string) *exec.Cmd {
	hoverflyBinaryUri := filepath.Join(workingDir, "bin/hoverfly")
	hoverflyCmd := exec.Command(hoverflyBinaryUri, "-ap", strconv.Itoa(adminPort), "-pp", strconv.Itoa(proxyPort))

	err := hoverflyCmd.Start()

	binaryErrorCheck(err, hoverflyBinaryUri)
	healthcheck(adminPort)

	return hoverflyCmd
}

func startHoverflyWithMiddleware(adminPort, proxyPort int, workingDir, binary, script string) *exec.Cmd {
	hoverflyBinaryUri := filepath.Join(workingDir, "bin/hoverfly")
	hoverflyCmd := exec.Command(hoverflyBinaryUri, "-ap", strconv.Itoa(adminPort), "-pp", strconv.Itoa(proxyPort))

	err := hoverflyCmd.Start()

	binaryErrorCheck(err, hoverflyBinaryUri)
	healthcheck(adminPort)

	request := sling.New().Put(fmt.Sprintf("http://localhost:%v/api/v2/hoverfly/middleware", adminPort)).BodyJSON(v2.MiddlewareView{Binary: binary, Script: script})

	DoRequest(request)

	return hoverflyCmd
}

func startHoverflyWithAuth(adminPort, proxyPort int, workingDir, username, password string) *exec.Cmd {
	os.Remove(filepath.Join(workingDir, "requests.db"))

	hoverflyBinaryUri := filepath.Join(workingDir, "bin/hoverfly")

	hoverflyAddUserCmd := exec.Command(hoverflyBinaryUri,"-db", "boltdb", "-add", "-username", username, "-password", password, "-ap", strconv.Itoa(adminPort), "-pp", strconv.Itoa(proxyPort))
	err := hoverflyAddUserCmd.Run()

	if err != nil {
		fmt.Println("Unable to start Hoverfly to add user")
		fmt.Println(hoverflyBinaryUri)
		fmt.Println("Is the binary there?")
		os.Exit(1)
	}

	hoverflyCmd := exec.Command(hoverflyBinaryUri,"-db", "boltdb", "-ap", strconv.Itoa(adminPort), "-pp", strconv.Itoa(proxyPort), "-auth", "true")
	err = hoverflyCmd.Start()

	binaryErrorCheck(err, hoverflyBinaryUri)
	healthcheck(adminPort)

	return hoverflyCmd
}

func startHoverflyWebserver(adminPort, proxyPort int, workingDir string) *exec.Cmd {
	hoverflyBinaryUri := filepath.Join(workingDir, "bin/hoverfly")
	hoverflyCmd := exec.Command(hoverflyBinaryUri, "-ap", strconv.Itoa(adminPort), "-pp", strconv.Itoa(proxyPort), "-webserver")

	err := hoverflyCmd.Start()

	binaryErrorCheck(err, hoverflyBinaryUri)
	healthcheck(adminPort)

	return hoverflyCmd
}

func binaryErrorCheck(err error, binaryPath string) {
	if err != nil {
		fmt.Println("Unable to start Hoverfly")
		fmt.Println(binaryPath)
		fmt.Println("Is the binary there?")
		os.Exit(1)
	}
}

func healthcheck(adminPort int) {
	var err error
	var resp *http.Response

	hasPassed := Eventually(func() int {
		resp, err = http.Get(fmt.Sprintf("http://localhost:%v/api/health", adminPort))
		if err == nil {
			return resp.StatusCode
		} else {
			return 0
		}
	}, time.Second*3).Should(BeNumerically("==", http.StatusOK))

	if !hasPassed {
		fmt.Println(err.Error())
	}
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

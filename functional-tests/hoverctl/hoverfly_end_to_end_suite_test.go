package hoverctl_end_to_end

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

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

	hoverctlBinary = filepath.Join(workingDirectory, "bin/hoverctl")

	binDirectory := filepath.Join(workingDirectory, "bin")

	os.Setenv("PATH", fmt.Sprintf("%v:%v", binDirectory, os.Getenv("PATH")))

})

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

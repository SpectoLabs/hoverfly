package hoverfly_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
	"github.com/phayes/freeport"
	"fmt"
	"net/http"
	"github.com/dghubble/sling"
	"strconv"
	"os"
	"time"
	"net/url"
	"strings"
	"io"
	"net/http/httptest"
	"os/exec"
	"path/filepath"
	"io/ioutil"
)

var (
	hoverflyAdminUrl string
	hoverflyProxyUrl string

	hoverflyCmd *exec.Cmd

	adminPort = freeport.GetPort()
	adminPortAsString = strconv.Itoa(adminPort)

	proxyPort = freeport.GetPort()
	proxyPortAsString = strconv.Itoa(proxyPort)
)

func TestHoverfly(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Hoverfly Suite")
}

var _ = BeforeSuite(func() {
	hoverflyAdminUrl = fmt.Sprintf("http://localhost:%v", adminPort)
	hoverflyProxyUrl = fmt.Sprintf("http://localhost:%v", proxyPort)

	os.Setenv("HTTP_PROXY", hoverflyProxyUrl)
	os.Setenv("HTTPS_PROXY", hoverflyProxyUrl)
})

var _ = AfterSuite(func() {
	os.Setenv("HTTP_PROXY", "")
	os.Setenv("HTTPS_PROXY", "")

	stopHoverfly()
})

func startHoverfly(adminPort, proxyPort int) * exec.Cmd {
	hoverflyBinaryUri := buildBinaryPath()
	hoverflyCmd := exec.Command(hoverflyBinaryUri, "-db", "memory", "-ap", strconv.Itoa(adminPort), "-pp", strconv.Itoa(proxyPort))

	err := hoverflyCmd.Start()

	binaryErrorCheck(err, hoverflyBinaryUri)
	healthcheck(adminPort)

	return hoverflyCmd
}

func startHoverflyWebServer(adminPort, proxyPort int) * exec.Cmd {
	hoverflyBinaryUri := buildBinaryPath()
	hoverflyCmd := exec.Command(hoverflyBinaryUri, "-db", "memory", "-ap", strconv.Itoa(adminPort), "-pp", strconv.Itoa(proxyPort), "-webserver")

	err := hoverflyCmd.Start()

	binaryErrorCheck(err, hoverflyBinaryUri)
	healthcheck(adminPort)

	return hoverflyCmd
}

func startHoverflyWithMiddleware(adminPort, proxyPort int, middlewarePath string) * exec.Cmd {
	hoverflyBinaryUri := buildBinaryPath()
	hoverflyCmd := exec.Command(hoverflyBinaryUri, "-db", "memory", "-ap", strconv.Itoa(adminPort), "-pp", strconv.Itoa(proxyPort), "-middleware", middlewarePath)
	hoverflyCmd.Stdout = os.Stdout
	hoverflyCmd.Stderr = os.Stderr

	err := hoverflyCmd.Start()

	binaryErrorCheck(err, hoverflyBinaryUri)
	healthcheck(adminPort)

	return hoverflyCmd
}

func buildBinaryPath() (string) {
	workingDirectory, _ := os.Getwd()
	return filepath.Join(workingDirectory, "bin/hoverfly")
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
	Eventually(func() int {
		resp, err := http.Get(fmt.Sprintf("http://localhost:%v/api/health", adminPort))
		if err == nil {
			return resp.StatusCode
		} else {
			fmt.Println(err.Error())
			return 0
		}
	}, time.Second * 3).Should(BeNumerically("==", http.StatusOK))
}

func stopHoverfly() {
	hoverflyCmd.Process.Kill()
}

func DoRequest(r *sling.Sling) (*http.Response) {
	req, err := r.Request()
	Expect(err).To(BeNil())
	response, err := http.DefaultClient.Do(req)

	Expect(err).To(BeNil())
	return response
}

func DoRequestThroughProxy(r *sling.Sling) (*http.Response) {
	req, err := r.Request()
	Expect(err).To(BeNil())

	proxy, err := url.Parse(hoverflyProxyUrl)
	proxyHttpClient := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxy)}}
	response, err := proxyHttpClient.Do(req)

	Expect(err).To(BeNil())

	return response
}

func SetHoverflyMode(mode string) {
	req := sling.New().Post(hoverflyAdminUrl + "/api/state").Body(strings.NewReader(`{"mode":"` + mode +`"}`))
	res := DoRequest(req)
	Expect(res.StatusCode).To(Equal(200))
}

func EraseHoverflyRecords() {
	req := sling.New().Delete(hoverflyAdminUrl + "/api/records")
	res := DoRequest(req)
	Expect(res.StatusCode).To(Equal(200))
}

func ExportHoverflyRecords() (io.Reader) {
	res := sling.New().Get(hoverflyAdminUrl + "/api/records")
	req := DoRequest(res)
	Expect(req.StatusCode).To(Equal(200))
	return req.Body
}

func ImportHoverflyRecords(payload io.Reader) {
	req := sling.New().Post(hoverflyAdminUrl + "/api/records").Body(payload)
	res := DoRequest(req)
	fmt.Println(hoverflyAdminUrl)
	Expect(res.StatusCode).To(Equal(200))
}

func CallFakeServerThroughProxy(server * httptest.Server) *http.Response {
	return DoRequestThroughProxy(sling.New().Get(server.URL))
}

func SetHoverflyResponseDelays(path string) {
	delaysConf, err := ioutil.ReadFile(path)
	if err != nil {
		Fail("can't read delay config file")
	}
	req := sling.New().Put(hoverflyAdminUrl + "/api/delays").Body(strings.NewReader(string(delaysConf)))
	res := DoRequest(req)
	Expect(res.StatusCode).To(Equal(201))
}
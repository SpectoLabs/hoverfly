package hoverfly_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"

	"github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/dghubble/sling"
	"github.com/phayes/freeport"
)

var (
	hoverflyAdminUrl string
	hoverflyProxyUrl string

	hoverflyCmd *exec.Cmd

	adminPort         = freeport.GetPort()
	adminPortAsString = strconv.Itoa(adminPort)

	proxyPort         = freeport.GetPort()
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

func startHoverfly(adminPort, proxyPort int) *exec.Cmd {
	return startHoverflyInternal("-ap", strconv.Itoa(adminPort), "-pp", strconv.Itoa(proxyPort))
}

func startHoverflyWebServer(adminPort, proxyPort int) *exec.Cmd {
	return startHoverflyInternal("-ap", strconv.Itoa(adminPort), "-pp", strconv.Itoa(proxyPort), "-webserver")
}

func startHoverflyWithDatabase(adminPort, proxyPort int) *exec.Cmd {
	return startHoverflyInternal("-db", "boltdb", "-ap", strconv.Itoa(adminPort), "-pp", strconv.Itoa(proxyPort))
}

func startHoverflyWebServerWithDatabase(adminPort, proxyPort int) *exec.Cmd {
	return startHoverflyInternal("-db", "boltdb", "-ap", strconv.Itoa(adminPort), "-pp", strconv.Itoa(proxyPort), "-webserver")
}

func startHoverflyWithMiddleware(adminPort, proxyPort int, middlewarePath string) *exec.Cmd {
	hoverflyCmd := startHoverflyInternal("-ap", strconv.Itoa(adminPort), "-pp", strconv.Itoa(proxyPort), "-middleware", middlewarePath)
	hoverflyCmd.Stdout = os.Stdout
	hoverflyCmd.Stderr = os.Stderr
	return hoverflyCmd
}

func startHoverflyInternal(commands ...string) *exec.Cmd {
	hoverflyBinaryUri := functional_tests.BuildBinaryPath()
	hoverflyCmd := exec.Command(hoverflyBinaryUri, commands...)

	err := hoverflyCmd.Start()

	functional_tests.BinaryErrorCheck(err, hoverflyBinaryUri)
	functional_tests.Healthcheck(adminPort)

	return hoverflyCmd
}

func stopHoverfly() {
	hoverflyCmd.Process.Kill()
}

func DoRequestThroughProxy(r *sling.Sling) *http.Response {
	req, err := r.Request()
	Expect(err).To(BeNil())

	proxy, err := url.Parse(hoverflyProxyUrl)
	proxyHttpClient := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxy)}, CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse }}
	response, err := proxyHttpClient.Do(req)

	Expect(err).To(BeNil())

	return response
}

func SetHoverflyMode(mode string) {
	req := sling.New().Put(hoverflyAdminUrl + "/api/v2/hoverfly/mode").Body(strings.NewReader(`{"mode":"` + mode + `"}`))
	res := functional_tests.DoRequest(req)
	Expect(res.StatusCode).To(Equal(200))
}

func SetHoverflyDestination(destination string) {
	req := sling.New().Put(hoverflyAdminUrl + "/api/v2/hoverfly/destination").Body(strings.NewReader(`{"destination":"` + destination + `"}`))
	res := functional_tests.DoRequest(req)
	Expect(res.StatusCode).To(Equal(200))
}

func EraseHoverflyRecords() {
	req := sling.New().Delete(hoverflyAdminUrl + "/api/records")
	res := functional_tests.DoRequest(req)
	Expect(res.StatusCode).To(Equal(200))
}

func ExportHoverflyRecords() io.Reader {
	res := sling.New().Get(hoverflyAdminUrl + "/api/records")
	req := functional_tests.DoRequest(res)
	Expect(req.StatusCode).To(Equal(200))
	return req.Body
}

func ImportHoverflyRecords(payload io.Reader) {
	req := sling.New().Post(hoverflyAdminUrl + "/api/records").Body(payload)
	res := functional_tests.DoRequest(req)
	Expect(res.StatusCode).To(Equal(200))
}

func ImportHoverflySimulation(payload io.Reader) *http.Response {
	req := sling.New().Put(hoverflyAdminUrl + "/api/v2/simulation").Body(payload)
	return functional_tests.DoRequest(req)
}

func CallFakeServerThroughProxy(server *httptest.Server) *http.Response {
	return DoRequestThroughProxy(sling.New().Get(server.URL))
}

func SetHoverflyResponseDelays(path string) {
	delaysConf, err := ioutil.ReadFile(path)
	if err != nil {
		Fail("can't read delay config file")
	}
	req := sling.New().Put(hoverflyAdminUrl + "/api/delays").Body(strings.NewReader(string(delaysConf)))
	res := functional_tests.DoRequest(req)
	Expect(res.StatusCode).To(Equal(201))
}

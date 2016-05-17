package hoverfly_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
	"github.com/phayes/freeport"
	"fmt"
	"net/http"
	//"go/build"
	//"strings"
	"github.com/SpectoLabs/hoverfly"
	"github.com/SpectoLabs/hoverfly/authentication/backends"
	"github.com/SpectoLabs/hoverfly/cache"
	"github.com/dghubble/sling"
	"strconv"
	"os"

	"github.com/Sirupsen/logrus"
	"time"
	"github.com/boltdb/bolt"
	"net/url"
	"strings"
	"io"
	"net/http/httptest"
)

var (
	hoverflyAdminUrl string
	hoverflyProxyUrl string
	cfg *hoverfly.Configuration
	db * bolt.DB
	requestCache cache.Cache
	hf * hoverfly.Hoverfly
)


func TestHoverfly(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Hoverfly Suite")
}

var _ = BeforeSuite(func() {

	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{})

	adminPort := strconv.Itoa(freeport.GetPort())
	proxyPort := strconv.Itoa(freeport.GetPort())

	fmt.Println("Admin: " + adminPort)
	fmt.Println("Proxy: " + proxyPort)
	hoverflyAdminUrl = fmt.Sprintf("http://localhost:%v", adminPort)
	hoverflyProxyUrl = fmt.Sprintf("http://localhost:%v", proxyPort)

	cfg = hoverfly.InitSettings()
	cfg.AdminPort = adminPort
	cfg.ProxyPort = proxyPort
	cfg.SetMode(hoverfly.SimulateMode)
	cfg.Verbose = true
	db = cache.GetDB(cfg.DatabasePath)


	requestCache = cache.NewBoltDBCache(db, []byte("requestsBucket"))
	metadataCache := cache.NewBoltDBCache(db, []byte("metadataBucket"))
	tokenCache := cache.NewBoltDBCache(db, []byte(backends.TokenBucketName))
	userCache := cache.NewBoltDBCache(db, []byte(backends.UserBucketName))
	authBackend := backends.NewCacheBasedAuthBackend(tokenCache, userCache)
	hf = hoverfly.GetNewHoverfly(cfg, requestCache, metadataCache, authBackend)

	err := hf.StartProxy()

	if err != nil {
		panic(err)
	}

	go hf.StartAdminInterface()

	os.Setenv("HTTP_PROXY", hoverflyProxyUrl)
	os.Setenv("HTTPS_PROXY", hoverflyProxyUrl)

	Eventually(func() int {
		resp, err := http.Get(hoverflyAdminUrl + "/api/state")
		if err == nil {
			return resp.StatusCode
		} else {
			return 0
		}
	}, time.Second * 3).Should(BeNumerically("==", http.StatusOK))
})

var _ = AfterSuite(func() {
	os.Setenv("HTTP_PROXY", "")
	os.Setenv("HTTPS_PROXY", "")
	db.Close()
})

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
	Expect(res.StatusCode).To(Equal(200))
}

func CallFakeServerThroughProxy(server * httptest.Server) *http.Response {
	return DoRequestThroughProxy(sling.New().Get(server.URL))
}
package hoverfly

import (
	"bytes"
	"fmt"
	"github.com/SpectoLabs/hoverfly/core/authentication/backends"
	"github.com/SpectoLabs/hoverfly/core/cache"
	"github.com/SpectoLabs/hoverfly/core/testutil"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"net/url"
	"github.com/SpectoLabs/hoverfly/core/models"
)

func TestGetNewHoverflyCheckConfig(t *testing.T) {

	cfg := InitSettings()

	db := cache.GetDB("testing2.db")
	requestCache := cache.NewBoltDBCache(db, []byte("requestBucket"))
	metaCache := cache.NewBoltDBCache(db, []byte("metaBucket"))
	tokenCache := cache.NewBoltDBCache(db, []byte("tokenBucket"))
	userCache := cache.NewBoltDBCache(db, []byte("userBucket"))
	backend := backends.NewCacheBasedAuthBackend(tokenCache, userCache)

	dbClient := GetNewHoverfly(cfg, requestCache, metaCache, backend)

	testutil.Expect(t, dbClient.Cfg, cfg)

	// deleting this database
	os.Remove("testing2.db")
}

func TestGetNewHoverfly(t *testing.T) {
	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	dbClient.Cfg.ProxyPort = "6666"

	err := dbClient.StartProxy()
	testutil.Expect(t, err, nil)

	newResponse, err := http.Get(fmt.Sprintf("http://localhost:%s/", dbClient.Cfg.ProxyPort))
	testutil.Expect(t, err, nil)
	testutil.Expect(t, newResponse.StatusCode, 500)

}

func TestProcessCaptureRequest(t *testing.T) {
	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()

	r, err := http.NewRequest("GET", "http://somehost.com", nil)
	testutil.Expect(t, err, nil)

	dbClient.Cfg.SetMode("capture")

	req, resp := dbClient.processRequest(r)

	testutil.Refute(t, req, nil)
	testutil.Refute(t, resp, nil)
	testutil.Expect(t, resp.StatusCode, 201)
}

func TestProcessSimulateRequest(t *testing.T) {
	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()

	r, err := http.NewRequest("GET", "http://somehost.com", nil)
	testutil.Expect(t, err, nil)

	// capturing
	dbClient.Cfg.SetMode("capture")
	req, resp := dbClient.processRequest(r)

	testutil.Refute(t, req, nil)
	testutil.Refute(t, resp, nil)
	testutil.Expect(t, resp.StatusCode, 201)

	// virtualizing
	dbClient.Cfg.SetMode(SimulateMode)
	newReq, newResp := dbClient.processRequest(r)

	testutil.Refute(t, newReq, nil)
	testutil.Refute(t, newResp, nil)
	testutil.Expect(t, newResp.StatusCode, 201)
}

func TestProcessSynthesizeRequest(t *testing.T) {
	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()

	// getting reflect middleware
	dbClient.Cfg.Middleware = "./examples/middleware/reflect_body/reflect_body.py"

	bodyBytes := []byte("request_body_here")

	r, err := http.NewRequest("GET", "http://somehost.com", ioutil.NopCloser(bytes.NewBuffer(bodyBytes)))
	testutil.Expect(t, err, nil)

	dbClient.Cfg.SetMode(SynthesizeMode)
	newReq, newResp := dbClient.processRequest(r)

	testutil.Refute(t, newReq, nil)
	testutil.Refute(t, newResp, nil)
	testutil.Expect(t, newResp.StatusCode, 200)
	b, err := ioutil.ReadAll(newResp.Body)
	testutil.Expect(t, err, nil)
	testutil.Expect(t, string(b), string(bodyBytes))
}

func TestProcessModifyRequest(t *testing.T) {
	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	// getting reflect middleware
	dbClient.Cfg.Middleware = "./examples/middleware/modify_request/modify_request.py"

	r, err := http.NewRequest("POST", "http://somehost.com", nil)
	testutil.Expect(t, err, nil)

	dbClient.Cfg.SetMode(ModifyMode)
	newReq, newResp := dbClient.processRequest(r)

	testutil.Refute(t, newReq, nil)
	testutil.Refute(t, newResp, nil)

	testutil.Expect(t, newResp.StatusCode, 202)
}

func TestURLToStringWorksAsExpected(t *testing.T) {
	testUrl := url.URL {
		Scheme: "http",
		Host: "test.com",
		Path: "/args/1",
		RawQuery: "query=val",
	}
	testutil.Expect(t, testUrl.String(), "http://test.com/args/1?query=val")
}

type ResponseDelayListStub struct {
	gotDelays int;
}

func (this *ResponseDelayListStub) Json() []byte {
	return nil
}

func (this *ResponseDelayListStub) Len() int {
	return this.Len()
}

func (this *ResponseDelayListStub) GetDelay(urlPattern, httpMethod string) (*models.ResponseDelay){
	this.gotDelays++;
	return nil;
}

func TestDelayAppliedToSuccessfulSimulateRequest(t *testing.T) {
	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()

	r, err := http.NewRequest("GET", "http://somehost.com", nil)
	testutil.Expect(t, err, nil)

	// capturing
	dbClient.Cfg.SetMode("capture")
	req, resp := dbClient.processRequest(r)

	testutil.Refute(t, req, nil)
	testutil.Refute(t, resp, nil)
	testutil.Expect(t, resp.StatusCode, 201)

	// virtualizing
	dbClient.Cfg.SetMode(SimulateMode)

	stub := ResponseDelayListStub{}
	dbClient.ResponseDelays = &stub

	newReq, newResp := dbClient.processRequest(r)

	testutil.Refute(t, newReq, nil)
	testutil.Refute(t, newResp, nil)
	testutil.Expect(t, newResp.StatusCode, 201)

	testutil.Expect(t, stub.gotDelays, 1)
}
package hoverfly

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/SpectoLabs/hoverfly/core/authentication/backends"
	"github.com/SpectoLabs/hoverfly/core/cache"
	"github.com/SpectoLabs/hoverfly/core/handlers/v1"
	"github.com/SpectoLabs/hoverfly/core/models"
	. "github.com/onsi/gomega"
)

func TestGetNewHoverflyCheckConfig(t *testing.T) {
	RegisterTestingT(t)

	cfg := InitSettings()

	db := cache.GetDB("testing2.db")
	requestCache := cache.NewBoltDBCache(db, []byte("requestBucket"))
	metaCache := cache.NewBoltDBCache(db, []byte("metaBucket"))
	tokenCache := cache.NewBoltDBCache(db, []byte("tokenBucket"))
	userCache := cache.NewBoltDBCache(db, []byte("userBucket"))
	backend := backends.NewCacheBasedAuthBackend(tokenCache, userCache)

	dbClient := GetNewHoverfly(cfg, requestCache, metaCache, backend)

	Expect(dbClient.Cfg).To(Equal(cfg))

	// deleting this database
	os.Remove("testing2.db")
}

func TestGetNewHoverfly(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	dbClient.Cfg.ProxyPort = "6666"

	err := dbClient.StartProxy()
	Expect(err).To(BeNil())

	newResponse, err := http.Get(fmt.Sprintf("http://localhost:%s/", dbClient.Cfg.ProxyPort))
	Expect(err).To(BeNil())
	Expect(newResponse.StatusCode).To(Equal(http.StatusInternalServerError))

}

func Test_Hoverfly_processRequest_CaptureModeReturnsResponseAndSavesIt(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()

	r, err := http.NewRequest("GET", "http://somehost.com", nil)
	Expect(err).To(BeNil())

	dbClient.Cfg.SetMode("capture")

	resp := dbClient.processRequest(r)

	Expect(resp).ToNot(BeNil())
	Expect(resp.StatusCode).To(Equal(http.StatusCreated))

	count, err := dbClient.RequestCache.RecordsCount()
	Expect(err).To(BeNil())

	Expect(count).To(Equal(1))
}

func Test_Hoverfly_processRequest_CanSimulateRequest(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()

	r, err := http.NewRequest("GET", "http://somehost.com", nil)
	Expect(err).To(BeNil())

	// capturing
	dbClient.Cfg.SetMode("capture")
	resp := dbClient.processRequest(r)

	Expect(resp).ToNot(BeNil())
	Expect(resp.StatusCode).To(Equal(http.StatusCreated))

	// virtualizing
	dbClient.Cfg.SetMode("simulate")
	newResp := dbClient.processRequest(r)

	Expect(newResp).ToNot(BeNil())
	Expect(newResp.StatusCode).To(Equal(http.StatusCreated))
}

func Test_Hoverfly_processRequest_CanUseMiddlewareToSynthesizeRequest(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()

	// getting reflect middleware
	err := dbClient.Cfg.Middleware.SetBinary("python")
	Expect(err).To(BeNil())

	err = dbClient.Cfg.Middleware.SetScript(pythonReflectBody)
	Expect(err).To(BeNil())

	bodyBytes := []byte("request_body_here")

	r, err := http.NewRequest("GET", "http://somehost.com", ioutil.NopCloser(bytes.NewBuffer(bodyBytes)))
	Expect(err).To(BeNil())

	dbClient.Cfg.SetMode("synthesize")
	newResp := dbClient.processRequest(r)

	Expect(newResp).ToNot(BeNil())
	Expect(newResp.StatusCode).To(Equal(http.StatusCreated))
	b, err := ioutil.ReadAll(newResp.Body)
	Expect(err).To(BeNil())
	Expect(string(b)).To(Equal(string(bodyBytes)))
}

func Test_Hoverfly_processRequest_CanModifyRequest(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	err := dbClient.Cfg.Middleware.SetBinary("python")
	Expect(err).To(BeNil())

	err = dbClient.Cfg.Middleware.SetScript(pythonModifyResponse)
	Expect(err).To(BeNil())

	r, err := http.NewRequest("POST", "http://somehost.com", nil)
	Expect(err).To(BeNil())

	dbClient.Cfg.SetMode("modify")
	newResp := dbClient.processRequest(r)

	Expect(newResp).ToNot(BeNil())

	Expect(newResp.StatusCode).To(Equal(http.StatusCreated))
}

type ResponseDelayListStub struct {
	gotDelays int
}

func (this *ResponseDelayListStub) Json() []byte {
	return nil
}

func (this *ResponseDelayListStub) Len() int {
	return this.Len()
}

func (this *ResponseDelayListStub) GetDelay(request models.RequestDetails) *models.ResponseDelay {
	this.gotDelays++
	return nil
}

func (this ResponseDelayListStub) ConvertToResponseDelayPayloadView() v1.ResponseDelayPayloadView {
	return v1.ResponseDelayPayloadView{}
}

func TestDelayAppliedToSuccessfulSimulateRequest(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()

	r, err := http.NewRequest("GET", "http://somehost.com", nil)
	Expect(err).To(BeNil())

	// capturing
	dbClient.Cfg.SetMode("capture")
	resp := dbClient.processRequest(r)

	Expect(resp.StatusCode).To(Equal(http.StatusCreated))

	// virtualizing
	dbClient.Cfg.SetMode("simulate")

	stub := ResponseDelayListStub{}
	dbClient.ResponseDelays = &stub

	newResp := dbClient.processRequest(r)

	Expect(newResp.StatusCode).To(Equal(http.StatusCreated))

	Expect(stub.gotDelays, Equal(1))
}

func TestDelayNotAppliedToFailedSimulateRequest(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	r, err := http.NewRequest("GET", "http://somehost.com", nil)
	Expect(err).To(BeNil())

	// virtualizing
	dbClient.Cfg.SetMode("simulate")

	stub := ResponseDelayListStub{}
	dbClient.ResponseDelays = &stub

	newResp := dbClient.processRequest(r)

	Expect(newResp.StatusCode).To(Equal(http.StatusBadGateway))

	Expect(stub.gotDelays).To(Equal(0))
}

func TestDelayNotAppliedToCaptureRequest(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()

	r, err := http.NewRequest("GET", "http://somehost.com", nil)
	Expect(err).To(BeNil())

	dbClient.Cfg.SetMode("capture")

	stub := ResponseDelayListStub{}
	dbClient.ResponseDelays = &stub

	resp := dbClient.processRequest(r)

	Expect(resp.StatusCode).To(Equal(http.StatusCreated))

	Expect(stub.gotDelays).To(Equal(0))
}

func TestDelayAppliedToSynthesizeRequest(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()

	err := dbClient.Cfg.Middleware.SetBinary("python")
	Expect(err).To(BeNil())

	err = dbClient.Cfg.Middleware.SetScript(pythonReflectBody)
	Expect(err).To(BeNil())

	bodyBytes := []byte("request_body_here")

	r, err := http.NewRequest("GET", "http://somehost.com", ioutil.NopCloser(bytes.NewBuffer(bodyBytes)))
	Expect(err).To(BeNil())

	dbClient.Cfg.SetMode("synthesize")

	stub := ResponseDelayListStub{}
	dbClient.ResponseDelays = &stub
	newResp := dbClient.processRequest(r)

	Expect(newResp.StatusCode).To(Equal(http.StatusCreated))

	Expect(stub.gotDelays).To(Equal(1))
}

func TestDelayNotAppliedToFailedSynthesizeRequest(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()

	err := dbClient.Cfg.Middleware.SetBinary("python")
	Expect(err).To(BeNil())

	err = dbClient.Cfg.Middleware.SetScript(pythonMiddlewareBad)
	Expect(err).To(BeNil())

	bodyBytes := []byte("request_body_here")

	r, err := http.NewRequest("GET", "http://somehost.com", ioutil.NopCloser(bytes.NewBuffer(bodyBytes)))
	Expect(err).To(BeNil())

	dbClient.Cfg.SetMode("synthesize")

	stub := ResponseDelayListStub{}
	dbClient.ResponseDelays = &stub
	newResp := dbClient.processRequest(r)

	Expect(newResp.StatusCode).To(Equal(http.StatusBadGateway))

	Expect(stub.gotDelays).To(Equal(0))
}

func TestDelayAppliedToSuccessfulMiddleware(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	err := dbClient.Cfg.Middleware.SetBinary("python")
	Expect(err).To(BeNil())

	err = dbClient.Cfg.Middleware.SetScript(pythonModifyResponse)
	Expect(err).To(BeNil())

	r, err := http.NewRequest("POST", "http://somehost.com", nil)
	Expect(err).To(BeNil())

	dbClient.Cfg.SetMode("modify")

	stub := ResponseDelayListStub{}
	dbClient.ResponseDelays = &stub
	newResp := dbClient.processRequest(r)

	Expect(newResp.StatusCode).To(Equal(http.StatusCreated))

	Expect(stub.gotDelays).To(Equal(1))
}

func TestDelayNotAppliedToFailedModifyRequest(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	err := dbClient.Cfg.Middleware.SetBinary("python")
	Expect(err).To(BeNil())

	err = dbClient.Cfg.Middleware.SetScript(pythonMiddlewareBad)
	Expect(err).To(BeNil())

	r, err := http.NewRequest("POST", "http://somehost.com", nil)
	Expect(err).To(BeNil())

	dbClient.Cfg.SetMode("modify")

	stub := ResponseDelayListStub{}
	dbClient.ResponseDelays = &stub
	newResp := dbClient.processRequest(r)

	Expect(newResp.StatusCode).To(Equal(http.StatusBadGateway))

	Expect(stub.gotDelays).To(Equal(0))
}

func Test_Hoverfly_DoRequest_DoesNotPanicWhenCannotMakeRequest(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	ioutil.NopCloser(bytes.NewBuffer([]byte("")))
	request, err := http.NewRequest("GET", "w.specto.fake", ioutil.NopCloser(bytes.NewBuffer([]byte(""))))
	Expect(err).To(BeNil())

	response, err := dbClient.DoRequest(request)
	Expect(response).To(BeNil())
	Expect(err).ToNot(BeNil())
}

func Test_Hoverfly_DoRequest_FailedHTTP(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(200, `{'message': 'here'}`)
	// stopping server
	server.Close()

	requestBody := []byte("fizz=buzz")

	body := ioutil.NopCloser(bytes.NewBuffer(requestBody))

	req, err := http.NewRequest("POST", "http://capture_body.com", body)
	Expect(err).To(BeNil())

	_, err = dbClient.DoRequest(req)
	Expect(err).ToNot(BeNil())
}

// TestCaptureHeader tests whether request gets new header assigned
func Test_DoRequest_AddsHoverflyHeaderOnSuccessfulRequest(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()

	req, err := http.NewRequest("GET", "http://example.com", ioutil.NopCloser(bytes.NewBuffer([]byte(""))))
	Expect(err).To(BeNil())

	response, err := dbClient.DoRequest(req)

	Expect(response.Header.Get("hoverfly")).To(Equal("Was-Here"))
}

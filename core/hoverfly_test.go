package hoverfly

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
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

func TestProcessCaptureRequest(t *testing.T) {
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
}

func TestProcessSimulateRequest(t *testing.T) {
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
	dbClient.Cfg.SetMode(SimulateMode)
	newResp := dbClient.processRequest(r)

	Expect(newResp).ToNot(BeNil())
	Expect(newResp.StatusCode).To(Equal(http.StatusCreated))
}

func TestProcessSynthesizeRequest(t *testing.T) {
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

	dbClient.Cfg.SetMode(SynthesizeMode)
	newResp := dbClient.processRequest(r)

	Expect(newResp).ToNot(BeNil())
	Expect(newResp.StatusCode).To(Equal(http.StatusCreated))
	b, err := ioutil.ReadAll(newResp.Body)
	Expect(err).To(BeNil())
	Expect(string(b)).To(Equal(string(bodyBytes)))
}

func TestProcessModifyRequest(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	err := dbClient.Cfg.Middleware.SetBinary("python")
	Expect(err).To(BeNil())

	err = dbClient.Cfg.Middleware.SetScript(pythonModifyResponse)
	Expect(err).To(BeNil())

	r, err := http.NewRequest("POST", "http://somehost.com", nil)
	Expect(err).To(BeNil())

	dbClient.Cfg.SetMode(ModifyMode)
	newResp := dbClient.processRequest(r)

	Expect(newResp).ToNot(BeNil())

	Expect(newResp.StatusCode).To(Equal(http.StatusCreated))
}

func TestURLToStringWorksAsExpected(t *testing.T) {
	RegisterTestingT(t)

	testUrl := url.URL{
		Scheme:   "http",
		Host:     "test.com",
		Path:     "/args/1",
		RawQuery: "query=val",
	}
	Expect(testUrl.String()).To(Equal("http://test.com/args/1?query=val"))
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
	dbClient.Cfg.SetMode(SimulateMode)

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
	dbClient.Cfg.SetMode(SimulateMode)

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

	dbClient.Cfg.SetMode(SynthesizeMode)

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

	dbClient.Cfg.SetMode(SynthesizeMode)

	stub := ResponseDelayListStub{}
	dbClient.ResponseDelays = &stub
	newResp := dbClient.processRequest(r)

	Expect(newResp.StatusCode).To(Equal(http.StatusServiceUnavailable))

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

	dbClient.Cfg.SetMode(ModifyMode)

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

	dbClient.Cfg.SetMode(ModifyMode)

	stub := ResponseDelayListStub{}
	dbClient.ResponseDelays = &stub
	newResp := dbClient.processRequest(r)

	Expect(newResp.StatusCode).To(Equal(503))

	Expect(stub.gotDelays).To(Equal(0))
}

func Test_Hoverfly_captureRequest_DoesNotPanicWhenCannotMakeRequest(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	request, err := http.NewRequest("GET", "w.specto.fake", nil)
	Expect(err).To(BeNil())

	response, err := dbClient.captureRequest(request)
	Expect(response).To(BeNil())
	Expect(err).ToNot(BeNil())
}

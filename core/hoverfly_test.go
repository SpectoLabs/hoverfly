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
	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/models"
	. "github.com/onsi/gomega"
)

const pythonMiddlewareBasic = "import sys\nprint(sys.stdin.readlines()[0])"

const pythonModifyResponse = "#!/usr/bin/env python\n" +
	"import sys\n" +
	"import json\n" +

	"def main():\n" +
	"	data = sys.stdin.readlines()\n" +
	"	payload = data[0]\n" +

	"	payload_dict = json.loads(payload)\n" +

	"	payload_dict['response']['status'] = 201\n" +
	"	payload_dict['response']['body'] = \"body was replaced by middleware\"\n" +

	"	print(json.dumps(payload_dict))\n" +

	"if __name__ == \"__main__\":\n" +
	"	main()\n"

const rubyModifyResponse = "#!/usr/bin/env ruby\n" +
	"# encoding: utf-8\n\n" +

	"require 'rubygems'\n" +
	"require 'json'\n\n" +

	"while payload = STDIN.gets\n" +
	"  next unless payload\n\n" +

	"  jsonPayload = JSON.parse(payload)\n\n" +

	"  jsonPayload[\"response\"][\"body\"] = \"body was replaced by middleware\\n\"\n\n" +

	"  STDOUT.puts jsonPayload.to_json\n\n" +

	"end"

const pythonReflectBody = "#!/usr/bin/env python\n" +
	"import sys\n" +
	"import json\n" +

	"def main():\n" +
	"	data = sys.stdin.readlines()\n" +
	"	payload = data[0]\n" +

	"	payload_dict = json.loads(payload)\n" +

	"	payload_dict['response']['status'] = 201\n" +
	"	payload_dict['response']['body'] = payload_dict['request']['body']\n" +

	"	print(json.dumps(payload_dict))\n" +

	"if __name__ == \"__main__\":\n" +
	"	main()\n"

const pythonMiddlewareBad = "this shouldn't work"

const rubyEcho = "#!/usr/bin/env ruby\n" +
	"# encoding: utf-8\n" +
	"while payload = STDIN.gets\n" +
	"  next unless payload\n" +
	"\n" +
	"  STDOUT.puts payload\n" +
	"\n" +
	"  STDERR.puts \"Payload data: #{payload}\"\n" +
	"\n" +
	"end"

// TestMain prepares database for testing and then performs a cleanup
func TestMain(m *testing.M) {
	setup()
	retCode := m.Run()

	// delete test database
	teardown()
	// call with result of m.Run()
	os.Exit(retCode)
}
func Test_NewHoverflyWithConfiguration_DoesNotCreateCacheIfCfgIsDisabled(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{
		DisableCache: true,
	})

	Expect(unit.CacheMatcher.RequestCache).To(BeNil())
}

func TestGetNewHoverflyCheckConfig(t *testing.T) {
	RegisterTestingT(t)

	cfg := InitSettings()

	db := cache.GetDB("testing2.db")
	requestCache := cache.NewDefaultLRUCache()
	tokenCache := cache.NewBoltDBCache(db, []byte("tokenBucket"))
	userCache := cache.NewBoltDBCache(db, []byte("userBucket"))
	backend := backends.NewCacheBasedAuthBackend(tokenCache, userCache)

	unit := GetNewHoverfly(cfg, requestCache, backend)

	Expect(unit.Cfg).To(Equal(cfg))

	// deleting this database
	os.Remove("testing2.db")
}

func Test_NewHoverflyWithConfiguration(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	unit.Cfg.ProxyPort = "6666"

	err := unit.StartProxy()
	Expect(err).To(BeNil())

	newResponse, err := http.Get(fmt.Sprintf("http://localhost:%s/", unit.Cfg.ProxyPort))
	Expect(err).To(BeNil())
	Expect(newResponse.StatusCode).To(Equal(http.StatusInternalServerError))

}

func Test_Hoverfly_processRequest_CaptureModeReturnsResponseAndSavesIt(t *testing.T) {
	RegisterTestingT(t)

	server, unit := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	r, err := http.NewRequest("GET", "http://somehost.com", nil)
	Expect(err).To(BeNil())

	unit.Cfg.SetMode("capture")

	resp := unit.processRequest(r)

	Expect(resp).ToNot(BeNil())
	Expect(resp.StatusCode).To(Equal(http.StatusCreated))

	Expect(unit.Simulation.GetMatchingPairs()).To(HaveLen(1))
}

func Test_Hoverfly_processRequest_CanSimulateRequest(t *testing.T) {
	RegisterTestingT(t)

	server, unit := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	r, err := http.NewRequest("GET", "http://somehost.com", nil)
	Expect(err).To(BeNil())

	// capturing
	unit.Cfg.SetMode("capture")
	resp := unit.processRequest(r)

	Expect(resp).ToNot(BeNil())
	Expect(resp.StatusCode).To(Equal(http.StatusCreated))

	// virtualizing
	unit.Cfg.SetMode("simulate")
	newResp := unit.processRequest(r)

	Expect(newResp).ToNot(BeNil())
	Expect(newResp.StatusCode).To(Equal(http.StatusCreated))
}

func Test_Hoverfly_processRequest_CanUseMiddlewareToSynthesizeRequest(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(InitSettings())

	// getting reflect middleware
	err := unit.Cfg.Middleware.SetBinary("python")
	Expect(err).To(BeNil())

	err = unit.Cfg.Middleware.SetScript(pythonReflectBody)
	Expect(err).To(BeNil())

	bodyBytes := []byte("request_body_here")

	r, err := http.NewRequest("GET", "http://somehost.com", ioutil.NopCloser(bytes.NewBuffer(bodyBytes)))
	Expect(err).To(BeNil())

	unit.Cfg.SetMode("synthesize")
	newResp := unit.processRequest(r)

	Expect(newResp).ToNot(BeNil())
	Expect(newResp.StatusCode).To(Equal(http.StatusCreated))
	b, err := ioutil.ReadAll(newResp.Body)
	Expect(err).To(BeNil())
	Expect(string(b)).To(Equal(string(bodyBytes)))
}

func Test_Hoverfly_processRequest_CanModifyRequest(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(InitSettings())

	err := unit.Cfg.Middleware.SetBinary("python")
	Expect(err).To(BeNil())

	err = unit.Cfg.Middleware.SetScript(pythonModifyResponse)
	Expect(err).To(BeNil())

	r, err := http.NewRequest("POST", "http://somehost.com", nil)
	Expect(err).To(BeNil())

	unit.Cfg.SetMode("modify")
	newResp := unit.processRequest(r)

	Expect(newResp).ToNot(BeNil())

	Expect(newResp.StatusCode).To(Equal(http.StatusCreated))
	Expect(newResp.Header).To(HaveKeyWithValue("Hoverfly", []string{"Was-Here"}))
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

type ResponseDelayLogNormalListStub struct {
	gotDelays int
}

func (this *ResponseDelayLogNormalListStub) Json() []byte {
	return nil
}

func (this *ResponseDelayLogNormalListStub) Len() int {
	return this.Len()
}

func (this *ResponseDelayLogNormalListStub) GetDelay(request models.RequestDetails) *models.ResponseDelayLogNormal {
	this.gotDelays++
	return nil
}

func (this ResponseDelayLogNormalListStub) ConvertToResponseDelayLogNormalPayloadView() v1.ResponseDelayLogNormalPayloadView {
	return v1.ResponseDelayLogNormalPayloadView{}
}

func Test_Hoverfly_processRequest_DelayAppliedToSuccessfulSimulateRequest(t *testing.T) {
	RegisterTestingT(t)

	server, unit := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	r, err := http.NewRequest("GET", "http://somehost.com", nil)
	Expect(err).To(BeNil())

	// capturing
	unit.Cfg.SetMode("capture")
	resp := unit.processRequest(r)

	Expect(resp.StatusCode).To(Equal(http.StatusCreated))

	// virtualizing
	unit.Cfg.SetMode("simulate")

	stub := ResponseDelayListStub{}
	unit.Simulation.ResponseDelays = &stub
	stubLogNormal := ResponseDelayLogNormalListStub{}
	unit.Simulation.ResponseDelaysLogNormal = &stubLogNormal

	newResp := unit.processRequest(r)

	Expect(newResp.StatusCode).To(Equal(http.StatusCreated))

	Expect(stub.gotDelays, Equal(1))
	Expect(stubLogNormal.gotDelays, Equal(1))
}

func Test_Hoverfly_processRequest_DelayNotAppliedToFailedSimulateRequest(t *testing.T) {
	RegisterTestingT(t)

	server, unit := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	r, err := http.NewRequest("GET", "http://somehost.com", nil)
	Expect(err).To(BeNil())

	// virtualizing
	unit.Cfg.SetMode("simulate")

	stub := ResponseDelayListStub{}
	unit.Simulation.ResponseDelays = &stub
	stubLogNormal := ResponseDelayLogNormalListStub{}
	unit.Simulation.ResponseDelaysLogNormal = &stubLogNormal

	newResp := unit.processRequest(r)

	Expect(newResp.StatusCode).To(Equal(http.StatusBadGateway))

	Expect(stub.gotDelays).To(Equal(0))
	Expect(stubLogNormal.gotDelays).To(Equal(0))
}

func Test_Hoverfly_processRequest_DelayNotAppliedToCaptureRequest(t *testing.T) {
	RegisterTestingT(t)

	server, unit := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	r, err := http.NewRequest("GET", "http://somehost.com", nil)
	Expect(err).To(BeNil())

	unit.Cfg.SetMode("capture")

	stub := ResponseDelayListStub{}
	unit.Simulation.ResponseDelays = &stub
	stubLogNormal := ResponseDelayLogNormalListStub{}
	unit.Simulation.ResponseDelaysLogNormal = &stubLogNormal

	resp := unit.processRequest(r)

	Expect(resp.StatusCode).To(Equal(http.StatusCreated))

	Expect(stub.gotDelays).To(Equal(0))
	Expect(stubLogNormal.gotDelays).To(Equal(0))
}

func Test_Hoverfly_processRequest_DelayAppliedToSynthesizeRequest(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	err := unit.Cfg.Middleware.SetBinary("python")
	Expect(err).To(BeNil())

	err = unit.Cfg.Middleware.SetScript(pythonReflectBody)
	Expect(err).To(BeNil())

	bodyBytes := []byte("request_body_here")

	r, err := http.NewRequest("GET", "http://somehost.com", ioutil.NopCloser(bytes.NewBuffer(bodyBytes)))
	Expect(err).To(BeNil())

	unit.Cfg.SetMode("synthesize")

	stub := ResponseDelayListStub{}
	unit.Simulation.ResponseDelays = &stub
	stubLogNormal := ResponseDelayLogNormalListStub{}
	unit.Simulation.ResponseDelaysLogNormal = &stubLogNormal
	newResp := unit.processRequest(r)

	Expect(newResp.StatusCode).To(Equal(http.StatusCreated))

	Expect(stub.gotDelays).To(Equal(1))
	Expect(stubLogNormal.gotDelays).To(Equal(1))
}

func Test_Hoverfly_processRequest_DelayNotAppliedToFailedSynthesizeRequest(t *testing.T) {
	RegisterTestingT(t)

	server, unit := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	err := unit.Cfg.Middleware.SetBinary("python")
	Expect(err).To(BeNil())

	err = unit.Cfg.Middleware.SetScript(pythonMiddlewareBad)
	Expect(err).To(BeNil())

	bodyBytes := []byte("request_body_here")

	r, err := http.NewRequest("GET", "http://somehost.com", ioutil.NopCloser(bytes.NewBuffer(bodyBytes)))
	Expect(err).To(BeNil())

	unit.Cfg.SetMode("synthesize")

	stub := ResponseDelayListStub{}
	unit.Simulation.ResponseDelays = &stub
	stubLogNormal := ResponseDelayLogNormalListStub{}
	unit.Simulation.ResponseDelaysLogNormal = &stubLogNormal
	newResp := unit.processRequest(r)

	Expect(newResp.StatusCode).To(Equal(http.StatusBadGateway))

	Expect(stub.gotDelays).To(Equal(0))
	Expect(stubLogNormal.gotDelays).To(Equal(0))
}

func Test_Hoverfly_processRequest_DelayAppliedToSuccessfulMiddleware(t *testing.T) {
	RegisterTestingT(t)

	server, unit := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	err := unit.Cfg.Middleware.SetBinary("python")
	Expect(err).To(BeNil())

	err = unit.Cfg.Middleware.SetScript(pythonModifyResponse)
	Expect(err).To(BeNil())

	r, err := http.NewRequest("POST", "http://somehost.com", nil)
	Expect(err).To(BeNil())

	unit.Cfg.SetMode("modify")

	stub := ResponseDelayListStub{}
	unit.Simulation.ResponseDelays = &stub
	stubLogNormal := ResponseDelayLogNormalListStub{}
	unit.Simulation.ResponseDelaysLogNormal = &stubLogNormal
	newResp := unit.processRequest(r)

	Expect(newResp.StatusCode).To(Equal(http.StatusCreated))

	Expect(stub.gotDelays).To(Equal(1))
	Expect(stubLogNormal.gotDelays).To(Equal(1))
}

func Test_Hoverfly_processRequest_DelayNotAppliedToFailedModifyRequest(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	err := unit.Cfg.Middleware.SetBinary("python")
	Expect(err).To(BeNil())

	err = unit.Cfg.Middleware.SetScript(pythonMiddlewareBad)
	Expect(err).To(BeNil())

	r, err := http.NewRequest("POST", "http://somehost.com", nil)
	Expect(err).To(BeNil())

	unit.Cfg.SetMode("modify")

	stub := ResponseDelayListStub{}
	unit.Simulation.ResponseDelays = &stub
	stubLogNormal := ResponseDelayLogNormalListStub{}
	unit.Simulation.ResponseDelaysLogNormal = &stubLogNormal
	newResp := unit.processRequest(r)

	Expect(newResp.StatusCode).To(Equal(http.StatusBadGateway))

	Expect(stub.gotDelays).To(Equal(0))
	Expect(stubLogNormal.gotDelays).To(Equal(0))
}

func Test_Hoverfly_processRequest_CanHandleResponseDiff(t *testing.T) {
	RegisterTestingT(t)

	server, expectedUnit := testTools(201, `{'message': 'expected'}`)

	r, err := http.NewRequest("GET", "http://somehost.com", nil)
	Expect(err).To(BeNil())

	// capturing
	expectedUnit.Cfg.SetMode("capture")
	resp := expectedUnit.processRequest(r)

	Expect(resp).ToNot(BeNil())
	Expect(resp.StatusCode).To(Equal(http.StatusCreated))

	server.Close()
	server, actualUnit := testTools(201, `{'message': 'actual'}`)
	defer server.Close()

	// comparing
	actualUnit.Cfg.SetMode("diff")
	actualUnit.Simulation = expectedUnit.Simulation
	newResp := actualUnit.processRequest(r)

	Expect(newResp).ToNot(BeNil())
	Expect(newResp.StatusCode).To(Equal(http.StatusCreated))
	Expect(len(actualUnit.responsesDiff)).To(Equal(1))
	requestDef := v2.SimpleRequestDefinitionView{
		Method: "GET",
		Host:   "somehost.com"}
	Expect(len(actualUnit.responsesDiff[requestDef])).To(Equal(1))
	Expect(actualUnit.responsesDiff[requestDef][0].DiffEntries).NotTo(BeEmpty())
	Expect(actualUnit.responsesDiff[requestDef][0].DiffEntries).To(ContainElement(
		v2.DiffReportEntry{Field: "body/message", Expected: "expected", Actual: "actual"}))

}

func TestMatchOnRequestBody(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	// preparing and saving requests/responses with unique bodies
	for i := 0; i < 5; i++ {
		req := &models.RequestDetails{
			Method:      "POST",
			Scheme:      "http",
			Destination: "capture_body.com",
			Body:        fmt.Sprintf("fizz=buzz, number=%d", i),
		}

		resp := &models.ResponseDetails{
			Status: 200,
			Body:   fmt.Sprintf("body here, number=%d", i),
		}

		unit.Save(req, resp, nil, false)
	}

	// now getting responses
	for i := 0; i < 5; i++ {
		requestBody := []byte(fmt.Sprintf("fizz=buzz, number=%d", i))
		body := ioutil.NopCloser(bytes.NewBuffer(requestBody))

		request, err := http.NewRequest("POST", "http://capture_body.com", body)
		Expect(err).To(BeNil())

		requestDetails, err := models.NewRequestDetailsFromHttpRequest(request)
		Expect(err).To(BeNil())

		response, err := unit.GetResponse(requestDetails)
		Expect(err).To(BeNil())

		Expect(response.Body).To(Equal(fmt.Sprintf("body here, number=%d", i)))

	}

}

// TODO: Fix by implementing Middleware check in Modify mode

// func TestModifyRequestNoMiddleware(t *testing.T) {
// 	RegisterTestingT(t)

// 	server, unit := testTools(201, `{'message': 'here'}`)
// 	defer server.Close()

// 	unit.SetMode("modify")

// 	unit.Cfg.Middleware.Binary = ""
// 	unit.Cfg.Middleware.Script = nil
// 	unit.Cfg.Middleware.Remote = ""

// 	req, err := http.NewRequest("GET", "http://very-interesting-website.com/q=123", nil)
// 	Expect(err).To(BeNil())

// 	response := unit.processRequest(req)

// 	responseBody, err := ioutil.ReadAll(response.Body)

// 	Expect(responseBody).To(Equal("THIS TEST IS BROKEN AND NEEDS FIXING"))

// 	Expect(response.StatusCode).To(Equal(http.StatusBadGateway))
// }

func Test_Hoverfly_StartProxy_StartProxyWOPort(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})
	unit.Cfg.ProxyPort = ""

	err := unit.StartProxy()
	Expect(err).ToNot(BeNil())
}

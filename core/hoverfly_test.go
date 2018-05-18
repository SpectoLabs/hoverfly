package hoverfly

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"encoding/json"

	"github.com/SpectoLabs/hoverfly/core/authentication/backends"
	"github.com/SpectoLabs/hoverfly/core/cache"
	"github.com/SpectoLabs/hoverfly/core/handlers/v1"
	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/matching"
	"github.com/SpectoLabs/hoverfly/core/matching/matchers"
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
	requestCache := cache.NewBoltDBCache(db, []byte("requestBucket"))
	tokenCache := cache.NewBoltDBCache(db, []byte("tokenBucket"))
	userCache := cache.NewBoltDBCache(db, []byte("userBucket"))
	backend := backends.NewCacheBasedAuthBackend(tokenCache, userCache)

	dbClient := GetNewHoverfly(cfg, requestCache, backend)

	Expect(dbClient.Cfg).To(Equal(cfg))

	// deleting this database
	os.Remove("testing2.db")
}

func TestGetNewHoverfly(t *testing.T) {
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

	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()

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
	Expect(newResp.Header).To(HaveKeyWithValue("Hoverfly", []string{"Was-Here"}))
}

func Test_Hoverfly_GetResponse_CanReturnResponseFromCache(t *testing.T) {
	RegisterTestingT(t)

	server, unit := testTools(201, `{'message': 'here'}`)
	server.Close()

	unit.CacheMatcher.SaveRequestMatcherResponsePair(models.RequestDetails{
		Destination: "somehost.com",
		Method:      "POST",
		Scheme:      "http",
	}, &models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "somehost.com",
				},
			},
			Method: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "POST",
				},
			},
			Scheme: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "http",
				},
			},
		},
		Response: models.ResponseDetails{
			Status: 200,
			Body:   "cached response",
		},
	}, nil)

	response, err := unit.GetResponse(models.RequestDetails{
		Destination: "somehost.com",
		Method:      "POST",
		Scheme:      "http",
	})

	Expect(err).To(BeNil())
	Expect(response).ToNot(BeNil())

	Expect(response.Status).To(Equal(http.StatusOK))
	Expect(response.Body).To(Equal("cached response"))
}

func Test_Hoverfly_GetResponse_CanReturnResponseFromSimulationAndNotCache(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	unit.Simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "somehost.com",
				},
			},
			Method: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "POST",
				},
			},
			Scheme: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "http",
				},
			},
		},
		Response: models.ResponseDetails{
			Status: 200,
			Body:   "response body",
		},
	})

	response, err := unit.GetResponse(models.RequestDetails{
		Destination: "somehost.com",
		Method:      "POST",
		Scheme:      "http",
	})

	Expect(err).To(BeNil())
	Expect(response).ToNot(BeNil())

	Expect(response.Status).To(Equal(http.StatusOK))
	Expect(response.Body).To(Equal("response body"))
}

func Test_Hoverfly_GetResponse_WillCacheResponseIfNotInCache(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	unit.Simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "somehost.com",
				},
			},
			Method: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "POST",
				},
			},
			Scheme: []models.RequestFieldMatchers{
				{
					Matcher: matchers.Exact,
					Value:   "http",
				},
			},
		},
		Response: models.ResponseDetails{
			Status: 200,
			Body:   "response body",
		},
	})

	unit.GetResponse(models.RequestDetails{
		Destination: "somehost.com",
		Method:      "POST",
		Scheme:      "http",
	})

	Expect(unit.CacheMatcher.RequestCache.RecordsCount()).Should(Equal(1))

	pairBytes, err := unit.CacheMatcher.RequestCache.Get([]byte("75b4ae6efa2a3f6d3ee6b9fed4d8c8c5"))
	Expect(err).To(BeNil())

	cachedRequestResponsePair, err := models.NewCachedResponseFromBytes(pairBytes)
	Expect(err).To(BeNil())

	Expect(cachedRequestResponsePair.MatchingPair.Response.Body).To(Equal("response body"))

	unit.Simulation = models.NewSimulation()
	response, err := unit.GetResponse(models.RequestDetails{
		Destination: "somehost.com",
		Method:      "POST",
		Scheme:      "http",
	})

	Expect(err).To(BeNil())
	Expect(response).ToNot(BeNil())

	Expect(response.Status).To(Equal(http.StatusOK))
	Expect(response.Body).To(Equal("response body"))
}

func Test_Hoverfly_GetResponse_WillReturnCachedResponseIfHeaderMatchIsFalse(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	requestDetails := models.RequestDetails{
		Destination: "somehost.com",
		Method:      "POST",
		Scheme:      "http",
	}

	unit.CacheMatcher.SaveRequestMatcherResponsePair(requestDetails, &models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{},
		Response: models.ResponseDetails{
			Body: "cached response",
		},
	}, nil)

	response, err := unit.GetResponse(requestDetails)
	Expect(err).To(BeNil())

	Expect(response.Body).To(Equal("cached response"))
}

func Test_Hoverfly_GetResponse_WillCacheMisses(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	requestDetails := models.RequestDetails{
		Destination: "somehost.com",
		Method:      "POST",
		Scheme:      "http",
	}

	_, err := unit.GetResponse(requestDetails)
	Expect(err.Error()).To(Equal("Could not find a match for request, create or record a valid matcher first!"))

	cachedResponse, err := unit.CacheMatcher.GetCachedResponse(&requestDetails)
	Expect(err).To(BeNil())

	Expect(cachedResponse.MatchingPair).To(BeNil())
	Expect(cachedResponse.ClosestMiss).To(BeNil())
}

func Test_Hoverfly_GetResponse_WillCacheClosestMiss(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})
	unit.PutSimulation(v2.SimulationViewV5{
		v2.DataViewV5{
			RequestResponsePairs: []v2.RequestMatcherResponsePairViewV5{
				{
					RequestMatcher: v2.RequestMatcherViewV5{
						Method: []v2.MatcherViewV5{
							{
								Matcher: matchers.Exact,
								Value:   "closest",
							},
						},
					},
					Response: v2.ResponseDetailsViewV5{
						Body: "closest",
					},
				},
			},
		},
		v2.MetaView{
			SchemaVersion: "v3",
		},
	})

	requestDetails := models.RequestDetails{
		Destination: "somehost.com",
		Method:      "POST",
		Scheme:      "http",
	}

	_, err := unit.GetResponse(requestDetails)
	Expect(err.Error()).ToNot(BeNil())

	cachedResponse, err := unit.CacheMatcher.GetCachedResponse(&requestDetails)
	Expect(err).To(BeNil())

	Expect(cachedResponse.MatchingPair).To(BeNil())
	Expect(cachedResponse.ClosestMiss.RequestMatcher.Method[0].Matcher).To(Equal("exact"))
	Expect(cachedResponse.ClosestMiss.RequestMatcher.Method[0].Value).To(Equal("closest"))

	Expect(cachedResponse.ClosestMiss.Response.Body).To(Equal("closest"))
	Expect(cachedResponse.ClosestMiss.MissedFields).To(ConsistOf("method"))
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

	r, err := http.NewRequest("GET", "http://somehost.com", nil)
	Expect(err).To(BeNil())

	// capturing
	dbClient.Cfg.SetMode("capture")
	resp := dbClient.processRequest(r)

	Expect(resp.StatusCode).To(Equal(http.StatusCreated))

	// virtualizing
	dbClient.Cfg.SetMode("simulate")

	stub := ResponseDelayListStub{}
	dbClient.Simulation.ResponseDelays = &stub

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
	dbClient.Simulation.ResponseDelays = &stub

	newResp := dbClient.processRequest(r)

	Expect(newResp.StatusCode).To(Equal(http.StatusBadGateway))

	Expect(stub.gotDelays).To(Equal(0))
}

func TestDelayNotAppliedToCaptureRequest(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	r, err := http.NewRequest("GET", "http://somehost.com", nil)
	Expect(err).To(BeNil())

	dbClient.Cfg.SetMode("capture")

	stub := ResponseDelayListStub{}
	dbClient.Simulation.ResponseDelays = &stub

	resp := dbClient.processRequest(r)

	Expect(resp.StatusCode).To(Equal(http.StatusCreated))

	Expect(stub.gotDelays).To(Equal(0))
}

func TestDelayAppliedToSynthesizeRequest(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	err := dbClient.Cfg.Middleware.SetBinary("python")
	Expect(err).To(BeNil())

	err = dbClient.Cfg.Middleware.SetScript(pythonReflectBody)
	Expect(err).To(BeNil())

	bodyBytes := []byte("request_body_here")

	r, err := http.NewRequest("GET", "http://somehost.com", ioutil.NopCloser(bytes.NewBuffer(bodyBytes)))
	Expect(err).To(BeNil())

	dbClient.Cfg.SetMode("synthesize")

	stub := ResponseDelayListStub{}
	dbClient.Simulation.ResponseDelays = &stub
	newResp := dbClient.processRequest(r)

	Expect(newResp.StatusCode).To(Equal(http.StatusCreated))

	Expect(stub.gotDelays).To(Equal(1))
}

func TestDelayNotAppliedToFailedSynthesizeRequest(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	err := dbClient.Cfg.Middleware.SetBinary("python")
	Expect(err).To(BeNil())

	err = dbClient.Cfg.Middleware.SetScript(pythonMiddlewareBad)
	Expect(err).To(BeNil())

	bodyBytes := []byte("request_body_here")

	r, err := http.NewRequest("GET", "http://somehost.com", ioutil.NopCloser(bytes.NewBuffer(bodyBytes)))
	Expect(err).To(BeNil())

	dbClient.Cfg.SetMode("synthesize")

	stub := ResponseDelayListStub{}
	dbClient.Simulation.ResponseDelays = &stub
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
	dbClient.Simulation.ResponseDelays = &stub
	newResp := dbClient.processRequest(r)

	Expect(newResp.StatusCode).To(Equal(http.StatusCreated))

	Expect(stub.gotDelays).To(Equal(1))
}

func TestDelayNotAppliedToFailedModifyRequest(t *testing.T) {
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
	newResp := unit.processRequest(r)

	Expect(newResp.StatusCode).To(Equal(http.StatusBadGateway))

	Expect(stub.gotDelays).To(Equal(0))
}

func Test_Hoverfly_DoRequest_DoesNotPanicWhenCannotMakeRequest(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	ioutil.NopCloser(bytes.NewBuffer([]byte("")))
	request, err := http.NewRequest("GET", "w.specto.fake", ioutil.NopCloser(bytes.NewBuffer([]byte(""))))
	Expect(err).To(BeNil())

	response, err := unit.DoRequest(request)
	Expect(response).To(BeNil())
	Expect(err).ToNot(BeNil())
}

func Test_Hoverfly_DoRequest_FailedHTTP(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	requestBody := []byte("fizz=buzz")

	body := ioutil.NopCloser(bytes.NewBuffer(requestBody))

	req, err := http.NewRequest("POST", "http://capture_body.com", body)
	Expect(err).To(BeNil())

	_, err = unit.DoRequest(req)
	Expect(err).ToNot(BeNil())
}

func Test_Hoverfly_Save_SavesRequestAndResponseToSimulation(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	unit.Save(&models.RequestDetails{
		Body:        "testbody",
		Destination: "testdestination",
		Headers:     map[string][]string{"testheader": []string{"testvalue"}},
		Method:      "testmethod",
		Path:        "/testpath",
		Query: map[string][]string{
			"query": []string{"test"},
		},
		Scheme: "http",
	}, &models.ResponseDetails{
		Body:    "testresponsebody",
		Headers: map[string][]string{"testheader": []string{"testvalue"}},
		Status:  200,
	}, nil, false)

	Expect(unit.Simulation.GetMatchingPairs()).To(HaveLen(1))

	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Body).To(HaveLen(1))
	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Body[0].Matcher).To(Equal("exact"))
	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Body[0].Value).To(Equal("testbody"))
	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Destination).To(HaveLen(1))
	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Destination[0].Matcher).To(Equal("exact"))
	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Destination[0].Value).To(Equal("testdestination"))
	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Method).To(HaveLen(1))
	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Method[0].Matcher).To(Equal("exact"))
	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Method[0].Value).To(Equal("testmethod"))
	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Path).To(HaveLen(1))
	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Path[0].Matcher).To(Equal("exact"))
	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Path[0].Value).To(Equal("/testpath"))
	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Query).To(HaveLen(1))
	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Query[0].Matcher).To(Equal("exact"))
	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Query[0].Value).To(Equal("query=test"))
	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Scheme).To(HaveLen(1))
	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Scheme[0].Matcher).To(Equal("exact"))
	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Scheme[0].Value).To(Equal("http"))

	Expect(unit.Simulation.GetMatchingPairs()[0].Response.Body).To(Equal("testresponsebody"))
	Expect(unit.Simulation.GetMatchingPairs()[0].Response.Headers).To(HaveKeyWithValue("testheader", []string{"testvalue"}))
	Expect(unit.Simulation.GetMatchingPairs()[0].Response.Status).To(Equal(200))
}

func Test_Hoverfly_Save_DoesNotSaveRequestHeadersWhenGivenHeadersArrayIsNil(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	unit.Save(&models.RequestDetails{
		Headers: map[string][]string{"testheader": []string{"testvalue"}},
	}, &models.ResponseDetails{
		Body:    "testresponsebody",
		Headers: map[string][]string{"testheader": []string{"testvalue"}},
		Status:  200,
	}, nil, false)

	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Headers).To(BeEmpty())
}

func Test_Hoverfly_Save_SavesAllRequestHeadersWhenGivenAnAsterisk(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	unit.Save(&models.RequestDetails{
		Headers: map[string][]string{
			"testheader":  []string{"testvalue"},
			"testheader2": []string{"testvalue2"},
		},
	}, &models.ResponseDetails{
		Body:    "testresponsebody",
		Headers: map[string][]string{"testheader": []string{"testvalue"}},
		Status:  200,
	}, []string{"*"}, false)

	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Headers).To(HaveLen(2))
	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Headers).To(HaveKeyWithValue("testheader", []string{"testvalue"}))
	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Headers).To(HaveKeyWithValue("testheader2", []string{"testvalue2"}))
}

func Test_Hoverfly_Save_SavesSpecificRequestHeadersWhenSpecifiedInHeadersArray(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	unit.Save(&models.RequestDetails{
		Headers: map[string][]string{
			"testheader":  []string{"testvalue"},
			"testheader2": []string{"testvalue2"},
		},
	}, &models.ResponseDetails{
		Body:    "testresponsebody",
		Headers: map[string][]string{"testheader": []string{"testvalue"}},
		Status:  200,
	}, []string{"testheader"}, false)

	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Headers).To(HaveLen(1))
	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Headers).To(HaveKeyWithValue("testheader", []string{"testvalue"}))
}

func Test_Hoverfly_Save_DoesNotSaveAnyRequestHeaderIfItDoesNotMatchEntryInHeadersArray(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	unit.Save(&models.RequestDetails{
		Headers: map[string][]string{
			"testheader":  []string{"testvalue"},
			"testheader2": []string{"testvalue2"},
		},
	}, &models.ResponseDetails{
		Body:    "testresponsebody",
		Headers: map[string][]string{"testheader": []string{"testvalue"}},
		Status:  200,
	}, []string{"nonmatch"}, false)

	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Headers).To(BeEmpty())
}

func Test_Hoverfly_Save_SavesMultipleRequestHeadersWhenMultiplesSpecifiedInHeadersArray(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	unit.Save(&models.RequestDetails{
		Headers: map[string][]string{
			"testheader":  []string{"testvalue"},
			"testheader2": []string{"testvalue2"},
			"nonmatch":    []string{"nonmatchvalue"},
		},
	}, &models.ResponseDetails{
		Body:    "testresponsebody",
		Headers: map[string][]string{"testheader": []string{"testvalue"}},
		Status:  200,
	}, []string{"testheader", "nonmatch"}, false)

	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Headers).To(HaveLen(2))
	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Headers).To(HaveKeyWithValue("testheader", []string{"testvalue"}))
	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Headers).To(HaveKeyWithValue("nonmatch", []string{"nonmatchvalue"}))
}

func Test_Hoverfly_Save_SavesIncompleteRequestAndResponseToSimulation(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	unit.Save(&models.RequestDetails{
		Destination: "testdestination",
	}, &models.ResponseDetails{
		Body:    "testresponsebody",
		Headers: map[string][]string{"testheader": []string{"testvalue"}},
		Status:  200,
	}, nil, false)

	Expect(unit.Simulation.GetMatchingPairs()).To(HaveLen(1))

	// Expect(unit.Simulation.MatchingPairs[0].RequestMatcher.Body).To(HaveLen(1))
	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Destination).To(HaveLen(1))
	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Destination[0].Matcher).To(Equal("exact"))
	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Destination[0].Value).To(Equal("testdestination"))
	// Expect(unit.Simulation.MatchingPairs[0].RequestMatcher.Headers).To(BeNil())
	// Expect(*unit.Simulation.MatchingPairs[0].RequestMatcher.Method).To(BeNil())
	// Expect(*unit.Simulation.MatchingPairs[0].RequestMatcher.Path).To(BeNil())
	// Expect(*unit.Simulation.MatchingPairs[0].RequestMatcher.Query).To(BeNil())
	// Expect(*unit.Simulation.MatchingPairs[0].RequestMatcher.Scheme).To(BeNil())

	Expect(unit.Simulation.GetMatchingPairs()[0].Response.Body).To(Equal("testresponsebody"))
	Expect(unit.Simulation.GetMatchingPairs()[0].Response.Headers).To(HaveKeyWithValue("testheader", []string{"testvalue"}))
	Expect(unit.Simulation.GetMatchingPairs()[0].Response.Status).To(Equal(200))
}

func Test_Hoverfly_Save_SavesRequestBodyAsJsonPathIfContentTypeIsJson(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	unit.Save(&models.RequestDetails{
		Body: `{"test": []}`,
		Headers: map[string][]string{
			"Content-Type": []string{"application/json"},
		},
	}, &models.ResponseDetails{}, nil, false)

	Expect(unit.Simulation.GetMatchingPairs()).To(HaveLen(1))

	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Body).To(HaveLen(1))
	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Body[0].Matcher).To(Equal("json"))
	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Body[0].Value).To(Equal(`{"test": []}`))
}

func Test_Hoverfly_Save_SavesRequestBodyAsXmlPathIfContentTypeIsXml(t *testing.T) {
	RegisterTestingT(t)

	unit := NewHoverflyWithConfiguration(&Configuration{})

	unit.Save(&models.RequestDetails{
		Body: `<xml>`,
		Headers: map[string][]string{
			"Content-Type": {"application/xml"},
		},
	}, &models.ResponseDetails{}, nil, false)

	Expect(unit.Simulation.GetMatchingPairs()).To(HaveLen(1))

	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Body).To(HaveLen(1))
	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Body[0].Matcher).To(Equal("xml"))
	Expect(unit.Simulation.GetMatchingPairs()[0].RequestMatcher.Body[0].Value).To(Equal(`<xml>`))
}

func Test_TransitioningBetweenStatesWhenSimulating(t *testing.T) {
	RegisterTestingT(t)

	simulation := `{
		"data": {
			"pairs": [{
					"request": {
						"path": [
							{
								"matcher": "exact",
								"value": "/basket"
							}
						]
					},
					"response": {
						"status": 200,
						"body": "empty"
					}
				},
				{
					"request": {
						"path": [
							{
								"matcher": "exact",
								"value": "/basket"
							}
						],
						"requiresState": {
							"eggs": "present"
						}
					},
					"response": {
						"status": 200,
						"body": "eggs"
					}
				},
				{
					"request": {
						"path": [
							{
								"matcher": "exact",
								"value": "/basket"
							}
						],
						"requiresState": {
							"bacon": "present"
						}
					},
					"response": {
						"status": 200,
						"body": "bacon"
					}
				},
				{
					"request": {
						"path": [
							{
								"matcher": "exact",
								"value": "/basket"
							}
						],
						"requiresState": {
							"eggs": "present",
							"bacon": "present"
						}
					},
					"response": {
						"status": 200,
						"body": "eggs, bacon"
					}
				},
				{
					"request": {
						"path": [
							{
								"matcher": "exact",
								"value": "/add-eggs"
							}
						]
					},
					"response": {
						"status": 200,
						"body": "added eggs",
						"transitionsState": {
							"eggs": "present"
						}
					}
				},
				{
					"request": {
						"path": [
							{
								"matcher": "exact",
								"value": "/add-bacon"
							}
						]
					},
					"response": {
						"status": 200,
						"body": "added bacon",
						"transitionsState": {
							"bacon": "present"
						}
					}
				},
				{
					"request": {
						"path": [
							{
								"matcher": "exact",
								"value": "/remove-eggs"
							}
						]
					},
					"response": {
						"status": 200,
						"body": "removed eggs",
						"removesState": ["eggs"]
					}
				},
				{
					"request": {
						"path": [
							{
								"matcher": "exact",
								"value": "/remove-bacon"
							}
						]
					},
					"response": {
						"status": 200,
						"body": "removed bacon",
						"removesState": ["bacon"]
					}
				}
			],
			"globalActions": {
				"delays": []
			}
		},
		"meta": {
			"schemaVersion": "v5",
			"hoverflyVersion": "v0.10.2",
			"timeExported": "2017-02-23T12:43:48Z"
		}
	}`

	v5 := &v2.SimulationViewV5{}

	json.Unmarshal([]byte(simulation), v5)

	hoverfly := NewHoverfly()
	hoverfly.CacheMatcher = matching.CacheMatcher{
		RequestCache: cache.NewInMemoryCache(),
	}
	hoverfly.PutSimulation(*v5)

	hoverfly.SetMode("simulate")

	response, _ := hoverfly.GetResponse(models.RequestDetails{
		Path: "/basket",
	})
	Expect(string(response.Body)).To(Equal(`empty`))

	response, _ = hoverfly.GetResponse(models.RequestDetails{
		Path: "/add-eggs",
	})
	Expect(string(response.Body)).To(Equal(`added eggs`))

	response, _ = hoverfly.GetResponse(models.RequestDetails{
		Path: "/basket",
	})
	Expect(string(response.Body)).To(Equal(`eggs`))

	response, _ = hoverfly.GetResponse(models.RequestDetails{
		Path: "/add-bacon",
	})
	Expect(string(response.Body)).To(Equal(`added bacon`))

	response, _ = hoverfly.GetResponse(models.RequestDetails{
		Path: "/basket",
	})
	Expect(string(response.Body)).To(Equal(`eggs, bacon`))

	response, _ = hoverfly.GetResponse(models.RequestDetails{
		Path: "/remove-eggs",
	})
	Expect(string(response.Body)).To(Equal(`removed eggs`))

	response, _ = hoverfly.GetResponse(models.RequestDetails{
		Path: "/basket",
	})
	Expect(string(response.Body)).To(Equal(`bacon`))

	response, _ = hoverfly.GetResponse(models.RequestDetails{
		Path: "/remove-bacon",
	})
	Expect(string(response.Body)).To(Equal(`removed bacon`))

	response, _ = hoverfly.GetResponse(models.RequestDetails{
		Path: "/basket",
	})
	Expect(string(response.Body)).To(Equal(`empty`))

	response, _ = hoverfly.GetResponse(models.RequestDetails{
		Path: "/basket",
	})
	Expect(string(response.Body)).To(Equal(`empty`))

	response, _ = hoverfly.GetResponse(models.RequestDetails{
		Path: "/add-eggs",
	})
	Expect(string(response.Body)).To(Equal(`added eggs`))

	response, _ = hoverfly.GetResponse(models.RequestDetails{
		Path: "/basket",
	})
	Expect(string(response.Body)).To(Equal(`eggs`))

	response, _ = hoverfly.GetResponse(models.RequestDetails{
		Path: "/add-bacon",
	})
	Expect(string(response.Body)).To(Equal(`added bacon`))

	response, _ = hoverfly.GetResponse(models.RequestDetails{
		Path: "/basket",
	})

	Expect(string(response.Body)).To(Equal(`eggs, bacon`))

	response, _ = hoverfly.GetResponse(models.RequestDetails{
		Path: "/remove-eggs",
	})
	Expect(string(response.Body)).To(Equal(`removed eggs`))

	response, _ = hoverfly.GetResponse(models.RequestDetails{
		Path: "/basket",
	})
	Expect(string(response.Body)).To(Equal(`bacon`))

	response, _ = hoverfly.GetResponse(models.RequestDetails{
		Path: "/remove-bacon",
	})
	Expect(string(response.Body)).To(Equal(`removed bacon`))

	response, _ = hoverfly.GetResponse(models.RequestDetails{
		Path: "/basket",
	})
	Expect(string(response.Body)).To(Equal(`empty`))
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

	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()

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

		dbClient.Save(req, resp, nil, false)
	}

	// now getting responses
	for i := 0; i < 5; i++ {
		requestBody := []byte(fmt.Sprintf("fizz=buzz, number=%d", i))
		body := ioutil.NopCloser(bytes.NewBuffer(requestBody))

		request, err := http.NewRequest("POST", "http://capture_body.com", body)
		Expect(err).To(BeNil())

		requestDetails, err := models.NewRequestDetailsFromHttpRequest(request)
		Expect(err).To(BeNil())

		response, err := dbClient.GetResponse(requestDetails)
		Expect(err).To(BeNil())

		Expect(response.Body).To(Equal(fmt.Sprintf("body here, number=%d", i)))

	}

}

func TestGetNotRecordedRequest(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()

	request, err := http.NewRequest("POST", "http://capture_body.com", nil)
	Expect(err).To(BeNil())

	requestDetails, err := models.NewRequestDetailsFromHttpRequest(request)
	Expect(err).To(BeNil())

	response, err := dbClient.GetResponse(requestDetails)
	Expect(err).ToNot(BeNil())

	Expect(response).To(BeNil())
}

// TODO: Fix by implementing Middleware check in Modify mode

// func TestModifyRequestNoMiddleware(t *testing.T) {
// 	RegisterTestingT(t)

// 	server, dbClient := testTools(201, `{'message': 'here'}`)
// 	defer server.Close()

// 	dbClient.SetMode("modify")

// 	dbClient.Cfg.Middleware.Binary = ""
// 	dbClient.Cfg.Middleware.Script = nil
// 	dbClient.Cfg.Middleware.Remote = ""

// 	req, err := http.NewRequest("GET", "http://very-interesting-website.com/q=123", nil)
// 	Expect(err).To(BeNil())

// 	response := dbClient.processRequest(req)

// 	responseBody, err := ioutil.ReadAll(response.Body)

// 	Expect(responseBody).To(Equal("THIS TEST IS BROKEN AND NEEDS FIXING"))

// 	Expect(response.StatusCode).To(Equal(http.StatusBadGateway))
// }

// func TestGetResponseCorruptedRequestResponsePair(t *testing.T) {
// 	RegisterTestingT(t)

// 	server, dbClient := testTools(200, `{'message': 'here'}`)
// 	defer server.Close()

// 	requestBody := []byte("fizz=buzz")

// 	body := ioutil.NopCloser(bytes.NewBuffer(requestBody))

// 	req, err := http.NewRequest("POST", "http://capture_body.com", body)
// 	Expect(err).To(BeNil())

// 	_, err = dbClient.captureRequest(req)
// 	Expect(err).To(BeNil())

// 	fp := matching.GetRequestFingerprint(req, requestBody, false)

// 	dbClient.RequestCache.Set([]byte(fp), []byte("you shall not decode me!"))

// 	// repeating process
// 	bodyNew := ioutil.NopCloser(bytes.NewBuffer(requestBody))

// 	reqNew, err := http.NewRequest("POST", "http://capture_body.com", bodyNew)
// 	Expect(err).To(BeNil())

// 	requestDetails, err := models.NewRequestDetailsFromHttpRequest(reqNew)
// 	Expect(err).To(BeNil())

// 	response, err := dbClient.GetResponse(requestDetails)
// 	Expect(err).ToNot(BeNil())

// 	Expect(response).To(BeNil())
// }

func TestStartProxyWOPort(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(200, `{'message': 'here'}`)
	// stopping server
	server.Close()

	dbClient.Cfg.ProxyPort = ""

	err := dbClient.StartProxy()
	Expect(err).ToNot(BeNil())
}

func TestSetDestination(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(200, `{'message': 'here'}`)
	// stopping server
	server.Close()
	dbClient.Cfg.ProxyPort = "5556"
	err := dbClient.StartProxy()
	Expect(err).To(BeNil())
	dbClient.SetDestination("newdest")

	Expect(dbClient.Cfg.Destination).To(Equal("newdest"))
}

func TestUpdateDestinationEmpty(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(200, `{'message': 'here'}`)
	// stopping server
	server.Close()
	dbClient.Cfg.ProxyPort = "5557"
	dbClient.StartProxy()
	err := dbClient.SetDestination("e^^**#")
	Expect(err).ToNot(BeNil())
}

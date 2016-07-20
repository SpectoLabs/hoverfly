package hoverfly

import (
	"bytes"
	"fmt"
	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/SpectoLabs/hoverfly/core/testutil"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"github.com/SpectoLabs/hoverfly/core/matching"
)

// TestMain prepares database for testing and then performs a cleanup
func TestMain(m *testing.M) {

	setup()

	retCode := m.Run()

	// delete test database
	teardown()

	// call with result of m.Run()
	os.Exit(retCode)
}

// TestCaptureHeader tests whether request gets new header assigned
func TestCaptureHeader(t *testing.T) {

	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()

	req, err := http.NewRequest("GET", "http://example.com", nil)
	testutil.Expect(t, err, nil)

	response, err := dbClient.captureRequest(req)

	testutil.Expect(t, response.Header.Get("hoverfly"), "Was-Here")
}

// TestRequestBodyCaptured tests whether request body is recorded
func TestRequestBodyCaptured(t *testing.T) {

	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()

	requestBody := []byte("fizz=buzz")

	body := ioutil.NopCloser(bytes.NewBuffer(requestBody))

	req, err := http.NewRequest("POST", "http://capture_body.com", body)
	testutil.Expect(t, err, nil)

	_, err = dbClient.captureRequest(req)
	testutil.Expect(t, err, nil)

	fp := matching.GetRequestFingerprint(req, requestBody, false)

	payloadBts, err := dbClient.RequestCache.Get([]byte(fp))
	testutil.Expect(t, err, nil)

	payload, err := models.NewPayloadFromBytes(payloadBts)
	testutil.Expect(t, err, nil)
	testutil.Expect(t, payload.Request.Body, "fizz=buzz")
}

func TestRequestBodySentToMiddleware(t *testing.T) {
	// sends a request with fizz=buzz body, server responds with {'message': 'here'}
	// then, since it's modify mode - middleware is applied again, this time
	// middleware takes original request body and replaces response body with it.
	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()

	requestBody := []byte("fizz=buzz")

	body := ioutil.NopCloser(bytes.NewBuffer(requestBody))

	req, err := http.NewRequest("POST", "http://capture_body.com", body)
	testutil.Expect(t, err, nil)

	resp, err := dbClient.modifyRequestResponse(req, "./examples/middleware/reflect_body/reflect_body.py")

	// body from the request should be in response body, instead of server's response
	responseBody, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	testutil.Expect(t, err, nil)
	testutil.Expect(t, string(responseBody), string(requestBody))

}

func TestMatchOnRequestBody(t *testing.T) {

	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()

	// preparing and saving requests/responses with unique bodies
	for i := 0; i < 5; i++ {
		requestBody := []byte(fmt.Sprintf("fizz=buzz, number=%d", i))
		body := ioutil.NopCloser(bytes.NewBuffer(requestBody))

		request, err := http.NewRequest("POST", "http://capture_body.com", body)
		testutil.Expect(t, err, nil)

		resp := models.ResponseDetails{
			Status: 200,
			Body:   fmt.Sprintf("body here, number=%d", i),
		}
		payload := models.Payload{Response: resp}

		// creating response
		c := NewConstructor(request, payload)
		response := c.ReconstructResponse()

		dbClient.save(request, requestBody, response, []byte(resp.Body))
	}

	// now getting responses
	for i := 0; i < 5; i++ {
		requestBody := []byte(fmt.Sprintf("fizz=buzz, number=%d", i))
		body := ioutil.NopCloser(bytes.NewBuffer(requestBody))

		request, _ := http.NewRequest("POST", "http://capture_body.com", body)

		response := dbClient.getResponse(request)

		responseBody, err := ioutil.ReadAll(response.Body)
		response.Body.Close()

		testutil.Expect(t, err, nil)
		testutil.Expect(t, string(responseBody), fmt.Sprintf("body here, number=%d", i))

	}

}

func TestGetNotRecordedRequest(t *testing.T) {
	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()

	request, _ := http.NewRequest("POST", "http://capture_body.com", nil)

	response := dbClient.getResponse(request)

	testutil.Expect(t, response.StatusCode, http.StatusPreconditionFailed)
}

// TestRequestFingerprint tests whether we get correct request ID
func TestRequestFingerprint(t *testing.T) {

	req, err := http.NewRequest("GET", "http://example.com", nil)
	testutil.Expect(t, err, nil)

	fp := matching.GetRequestFingerprint(req, []byte(""), false)

	testutil.Expect(t, fp, "92a65ed4ca2b7100037a4cba9afd15ea")
}

// TestRequestFingerprintBody tests where request body is also used to create unique request ID
func TestRequestFingerprintBody(t *testing.T) {

	req, err := http.NewRequest("GET", "http://example.com", nil)
	testutil.Expect(t, err, nil)

	fp := matching.GetRequestFingerprint(req, []byte("some huge XML or JSON here"), false)

	testutil.Expect(t, fp, "b3918a54eb6e42652e29e14c21ba8f81")
}

func TestScheme(t *testing.T) {

	req, err := http.NewRequest("GET", "http://example.com", nil)
	testutil.Expect(t, err, nil)

	originalFp := matching.GetRequestFingerprint(req, []byte(""), false)

	httpsReq, err := http.NewRequest("GET", "https://example.com", nil)
	testutil.Expect(t, err, nil)

	newFp := matching.GetRequestFingerprint(httpsReq, []byte(""), false)

	// fingerprint should be the same
	testutil.Expect(t, originalFp, newFp)
}

func TestDeleteAllRecords(t *testing.T) {

	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	// inserting some payloads
	for i := 0; i < 5; i++ {
		req, err := http.NewRequest("GET", fmt.Sprintf("http://delete_all_records.com/q=%d", i), nil)
		testutil.Expect(t, err, nil)
		dbClient.captureRequest(req)
	}
	err := dbClient.RequestCache.DeleteData()
	testutil.Expect(t, err, nil)
}

func TestPayloadEncodeDecode(t *testing.T) {
	resp := models.ResponseDetails{
		Status: 200,
		Body:   "body here",
	}

	payload := models.Payload{Response: resp}

	bts, err := payload.Encode()
	testutil.Expect(t, err, nil)

	pl, err := models.NewPayloadFromBytes(bts)
	testutil.Expect(t, err, nil)
	testutil.Expect(t, pl.Response.Body, resp.Body)
	testutil.Expect(t, pl.Response.Status, resp.Status)

}

func TestPayloadEncodeEmpty(t *testing.T) {
	payload := models.Payload{}

	bts, err := payload.Encode()
	testutil.Expect(t, err, nil)

	_, err = models.NewPayloadFromBytes(bts)
	testutil.Expect(t, err, nil)
}

func TestDecodeRandomBytes(t *testing.T) {
	bts := []byte("some random stuff here")
	_, err := models.NewPayloadFromBytes(bts)
	testutil.Refute(t, err, nil)
}

func TestModifyRequest(t *testing.T) {
	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	dbClient.Cfg.Middleware = "./examples/middleware/modify_request/modify_request.py"

	req, err := http.NewRequest("GET", "http://very-interesting-website.com/q=123", nil)
	testutil.Expect(t, err, nil)

	response, err := dbClient.modifyRequestResponse(req, dbClient.Cfg.Middleware)
	testutil.Expect(t, err, nil)

	// response should be changed to 202
	testutil.Expect(t, response.StatusCode, 202)

}

func TestModifyRequestWODestination(t *testing.T) {
	// tests modify mode but uses different middleware to not supply destination
	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	dbClient.Cfg.Middleware = "./examples/middleware/modify_response/modify_response.py"

	req, err := http.NewRequest("GET", "http://very-interesting-website.com/q=123", nil)
	testutil.Expect(t, err, nil)

	response, err := dbClient.modifyRequestResponse(req, dbClient.Cfg.Middleware)
	testutil.Expect(t, err, nil)

	// response should be changed to 201
	testutil.Expect(t, response.StatusCode, 201)

}

func TestModifyRequestNoMiddleware(t *testing.T) {
	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	dbClient.Cfg.Middleware = ""

	req, err := http.NewRequest("GET", "http://very-interesting-website.com/q=123", nil)
	testutil.Expect(t, err, nil)

	_, err = dbClient.modifyRequestResponse(req, dbClient.Cfg.Middleware)
	testutil.Refute(t, err, nil)
}

func TestGetResponseCorruptedPayload(t *testing.T) {

	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()

	requestBody := []byte("fizz=buzz")

	body := ioutil.NopCloser(bytes.NewBuffer(requestBody))

	req, err := http.NewRequest("POST", "http://capture_body.com", body)
	testutil.Expect(t, err, nil)

	_, err = dbClient.captureRequest(req)
	testutil.Expect(t, err, nil)

	fp := matching.GetRequestFingerprint(req, requestBody, false)

	dbClient.RequestCache.Set([]byte(fp), []byte("you shall not decode me!"))

	// repeating process
	bodyNew := ioutil.NopCloser(bytes.NewBuffer(requestBody))

	reqNew, err := http.NewRequest("POST", "http://capture_body.com", bodyNew)
	testutil.Expect(t, err, nil)
	response := dbClient.getResponse(reqNew)

	testutil.Expect(t, response.StatusCode, http.StatusInternalServerError)

}

func TestDoRequestWFailedMiddleware(t *testing.T) {

	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()

	// adding middleware which doesn't exist, doRequest should return error
	dbClient.Cfg.Middleware = "./should/not/exist.go"

	requestBody := []byte("fizz=buzz")

	body := ioutil.NopCloser(bytes.NewBuffer(requestBody))

	req, err := http.NewRequest("POST", "http://capture_body.com", body)
	testutil.Expect(t, err, nil)

	_, _, err = dbClient.doRequest(req)
	testutil.Refute(t, err, nil)
}

func TestDoRequestFailedHTTP(t *testing.T) {
	server, dbClient := testTools(200, `{'message': 'here'}`)
	// stopping server
	server.Close()

	requestBody := []byte("fizz=buzz")

	body := ioutil.NopCloser(bytes.NewBuffer(requestBody))

	req, err := http.NewRequest("POST", "http://capture_body.com", body)
	testutil.Expect(t, err, nil)

	_, _, err = dbClient.doRequest(req)
	testutil.Refute(t, err, nil)

}

func TestStartProxyWOPort(t *testing.T) {
	server, dbClient := testTools(200, `{'message': 'here'}`)
	// stopping server
	server.Close()

	dbClient.Cfg.ProxyPort = ""

	err := dbClient.StartProxy()
	testutil.Refute(t, err, nil)
}

func TestUpdateDestination(t *testing.T) {
	server, dbClient := testTools(200, `{'message': 'here'}`)
	// stopping server
	server.Close()
	dbClient.Cfg.ProxyPort = "5556"
	err := dbClient.StartProxy()
	testutil.Expect(t, err, nil)
	dbClient.UpdateDestination("newdest")

	testutil.Expect(t, dbClient.Cfg.Destination, "newdest")
}

func TestUpdateDestinationEmpty(t *testing.T) {
	server, dbClient := testTools(200, `{'message': 'here'}`)
	// stopping server
	server.Close()
	dbClient.Cfg.ProxyPort = "5557"
	dbClient.StartProxy()
	err := dbClient.UpdateDestination("e^^**#")
	testutil.Refute(t, err, nil)
}

func TestJSONMinifier(t *testing.T) {

	// body can be nil here, it's not reading it from request anyway
	req, err := http.NewRequest("GET", "http://example.com", nil)
	testutil.Expect(t, err, nil)
	req.Header.Add("Content-Type", "application/json")

	fpOne := matching.GetRequestFingerprint(req, []byte(`{"foo": "bar"}`), false)
	fpTwo := matching.GetRequestFingerprint(req, []byte(`{     "foo":           "bar"}`), false)

	testutil.Expect(t, fpOne, fpTwo)
}

func TestJSONMinifierWOHeader(t *testing.T) {

	// body can be nil here, it's not reading it from request anyway
	req, err := http.NewRequest("GET", "http://example.com", nil)
	testutil.Expect(t, err, nil)

	// application/json header is not set, shouldn't be equal
	fpOne := matching.GetRequestFingerprint(req, []byte(`{"foo": "bar"}`), false)
	fpTwo := matching.GetRequestFingerprint(req, []byte(`{     "foo":           "bar"}`), false)

	testutil.Refute(t, fpOne, fpTwo)
}

var xmlBody = `<project xmlns="http://maven.apache.org/POM/4.0.0" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
		  xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 http://maven.apache.org/maven-v4_0_0.xsd">
		  <modelVersion>4.0.0</modelVersion>
		  <groupId>some ID here</groupId>
	       </project>`

var xmlBodyTwo = `<project xmlns="http://maven.apache.org/POM/4.0.0" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
		  xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 http://maven.apache.org/maven-v4_0_0.xsd">


		  <modelVersion>4.0.0</modelVersion>


		  <groupId>some ID here</groupId>
		  
	       </project>`

func TestXMLMinifier(t *testing.T) {

	// body can be nil here, it's not reading it from request anyway
	req, err := http.NewRequest("GET", "http://example.com", nil)
	testutil.Expect(t, err, nil)

	req.Header.Add("Content-Type", "application/xml")

	fpOne := matching.GetRequestFingerprint(req, []byte(xmlBody), false)
	fpTwo := matching.GetRequestFingerprint(req, []byte(xmlBodyTwo), false)
	testutil.Expect(t, fpOne, fpTwo)
}

func TestXMLMinifierWOHeader(t *testing.T) {

	// body can be nil here, it's not reading it from request anyway
	req, err := http.NewRequest("GET", "http://example.com", nil)
	testutil.Expect(t, err, nil)

	// application/xml header is not set, shouldn't be equal
	fpOne := matching.GetRequestFingerprint(req, []byte(xmlBody), false)
	fpTwo := matching.GetRequestFingerprint(req, []byte(xmlBodyTwo), false)
	testutil.Refute(t, fpOne, fpTwo)
}

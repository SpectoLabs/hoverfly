package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

// TestMain prepares database for testing and then performs a cleanup
func TestMain(m *testing.M) {

	setup()

	retCode := m.Run()

	// your func
	teardown()

	// call with result of m.Run()
	os.Exit(retCode)
}

// TestCaptureHeader tests whether request gets new header assigned
func TestCaptureHeader(t *testing.T) {

	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()

	req, err := http.NewRequest("GET", "http://example.com", nil)
	expect(t, err, nil)

	response, err := dbClient.captureRequest(req)

	expect(t, response.Header.Get("hoverfly"), "Was-Here")
}

// TestRequestBodyCaptured tests whether request body is recorded
func TestRequestBodyCaptured(t *testing.T) {

	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()

	requestBody := []byte("fizz=buzz")

	body := ioutil.NopCloser(bytes.NewBuffer(requestBody))

	req, err := http.NewRequest("POST", "http://capture_body.com", body)
	expect(t, err, nil)

	_, err = dbClient.captureRequest(req)
	expect(t, err, nil)

	fp := getRequestFingerprint(req, requestBody)

	payloadBts, err := dbClient.cache.Get([]byte(fp))
	expect(t, err, nil)

	payload, err := decodePayload(payloadBts)
	expect(t, err, nil)
	expect(t, payload.Request.Body, "fizz=buzz")
}

func TestMatchOnRequestBody(t *testing.T) {

	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()

	// preparing and saving requests/responses with unique bodies
	for i := 0; i < 5; i++ {
		requestBody := []byte(fmt.Sprintf("fizz=buzz, number=%d", i))
		body := ioutil.NopCloser(bytes.NewBuffer(requestBody))

		request, err := http.NewRequest("POST", "http://capture_body.com", body)
		expect(t, err, nil)

		resp := response{
			Status: 200,
			Body:   fmt.Sprintf("body here, number=%d", i),
		}
		payload := Payload{Response: resp}

		// creating response
		c := NewConstructor(request, payload)
		response := c.reconstructResponse()

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

		expect(t, err, nil)
		expect(t, string(responseBody), fmt.Sprintf("body here, number=%d", i))

	}

}

func TestGetNotRecordedRequest(t *testing.T) {
	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()

	request, _ := http.NewRequest("POST", "http://capture_body.com", nil)

	response := dbClient.getResponse(request)

	expect(t, response.StatusCode, http.StatusPreconditionFailed)
}

// TestRequestFingerprint tests whether we get correct request ID
func TestRequestFingerprint(t *testing.T) {

	req, err := http.NewRequest("GET", "http://example.com", nil)
	expect(t, err, nil)

	fp := getRequestFingerprint(req, []byte(""))

	expect(t, fp, "92a65ed4ca2b7100037a4cba9afd15ea")
}

// TestRequestFingerprintBody tests where request body is also used to create unique request ID
func TestRequestFingerprintBody(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com", nil)
	expect(t, err, nil)

	fp := getRequestFingerprint(req, []byte("some huge XML or JSON here"))

	expect(t, fp, "b3918a54eb6e42652e29e14c21ba8f81")
}

func TestScheme(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com", nil)
	expect(t, err, nil)

	originalFp := getRequestFingerprint(req, []byte(""))

	httpsReq, err := http.NewRequest("GET", "https://example.com", nil)
	expect(t, err, nil)

	newFp := getRequestFingerprint(httpsReq, []byte(""))

	// fingerprint should be the same
	expect(t, originalFp, newFp)
}

func TestDeleteAllRecords(t *testing.T) {

	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	// inserting some payloads
	for i := 0; i < 5; i++ {
		req, err := http.NewRequest("GET", fmt.Sprintf("http://delete_all_records.com/q=%d", i), nil)
		expect(t, err, nil)
		dbClient.captureRequest(req)
	}
	err := dbClient.cache.DeleteBucket(dbClient.cache.requestsBucket)
	expect(t, err, nil)
}

func TestPayloadEncodeDecode(t *testing.T) {
	resp := response{
		Status: 200,
		Body:   "body here",
	}

	payload := Payload{Response: resp}

	bts, err := payload.encode()
	expect(t, err, nil)

	pl, err := decodePayload(bts)
	expect(t, err, nil)
	expect(t, pl.Response.Body, resp.Body)
	expect(t, pl.Response.Status, resp.Status)

}

func TestPayloadEncodeEmpty(t *testing.T) {
	payload := Payload{}

	bts, err := payload.encode()
	expect(t, err, nil)

	_, err = decodePayload(bts)
	expect(t, err, nil)
}

func TestDecodeRandomBytes(t *testing.T) {
	bts := []byte("some random stuff here")
	_, err := decodePayload(bts)
	refute(t, err, nil)
}

func TestModifyRequest(t *testing.T) {
	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	dbClient.cfg.middleware = "./examples/middleware/modify_request/modify_request.py"

	req, err := http.NewRequest("GET", "http://very-interesting-website.com/q=123", nil)
	expect(t, err, nil)

	response, err := dbClient.modifyRequestResponse(req, dbClient.cfg.middleware)
	expect(t, err, nil)

	// response should be changed to 202
	expect(t, response.StatusCode, 202)

}

func TestModifyRequestWODestination(t *testing.T) {
	// tests modify mode but uses different middleware to not supply destination
	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	dbClient.cfg.middleware = "./examples/middleware/modify_response/modify_response.py"

	req, err := http.NewRequest("GET", "http://very-interesting-website.com/q=123", nil)
	expect(t, err, nil)

	response, err := dbClient.modifyRequestResponse(req, dbClient.cfg.middleware)
	expect(t, err, nil)

	// response should be changed to 201
	expect(t, response.StatusCode, 201)

}

func TestModifyRequestNoMiddleware(t *testing.T) {
	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	dbClient.cfg.middleware = ""

	req, err := http.NewRequest("GET", "http://very-interesting-website.com/q=123", nil)
	expect(t, err, nil)

	_, err = dbClient.modifyRequestResponse(req, dbClient.cfg.middleware)
	refute(t, err, nil)
}

func TestGetResponseCorruptedPayload(t *testing.T) {

	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()

	requestBody := []byte("fizz=buzz")

	body := ioutil.NopCloser(bytes.NewBuffer(requestBody))

	req, err := http.NewRequest("POST", "http://capture_body.com", body)
	expect(t, err, nil)

	_, err = dbClient.captureRequest(req)
	expect(t, err, nil)

	fp := getRequestFingerprint(req, requestBody)

	dbClient.cache.Set([]byte(fp), []byte("you shall not decode me!"))

	// repeating process
	bodyNew := ioutil.NopCloser(bytes.NewBuffer(requestBody))

	reqNew, err := http.NewRequest("POST", "http://capture_body.com", bodyNew)
	expect(t, err, nil)
	response := dbClient.getResponse(reqNew)

	expect(t, response.StatusCode, http.StatusInternalServerError)

}

func TestDoRequestWFailedMiddleware(t *testing.T) {

	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()

	// adding middleware which doesn't exist, doRequest should return error
	dbClient.cfg.middleware = "./should/not/exist.go"

	requestBody := []byte("fizz=buzz")

	body := ioutil.NopCloser(bytes.NewBuffer(requestBody))

	req, err := http.NewRequest("POST", "http://capture_body.com", body)
	expect(t, err, nil)

	_, err = dbClient.doRequest(req)
	refute(t, err, nil)
}

func TestDoRequestFailedHTTP(t *testing.T) {
	server, dbClient := testTools(200, `{'message': 'here'}`)
	// stopping server
	server.Close()

	requestBody := []byte("fizz=buzz")

	body := ioutil.NopCloser(bytes.NewBuffer(requestBody))

	req, err := http.NewRequest("POST", "http://capture_body.com", body)
	expect(t, err, nil)

	_, err = dbClient.doRequest(req)
	refute(t, err, nil)

}

package hoverfly

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/SpectoLabs/hoverfly/core/models"
	. "github.com/onsi/gomega"
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

func TestRequestBodySentToMiddleware(t *testing.T) {
	RegisterTestingT(t)

	// sends a request with fizz=buzz body, server responds with {'message': 'here'}
	// then, since it's modify mode - middleware is applied again, this time
	// middleware takes original request body and replaces response body with it.
	server, dbClient := testTools(200, `{'message': 'here'}`)
	defer server.Close()

	dbClient.SetMode("modify")

	requestBody := []byte("fizz=buzz")

	body := ioutil.NopCloser(bytes.NewBuffer(requestBody))

	req, err := http.NewRequest("POST", "http://capture_body.com", body)
	Expect(err).To(BeNil())

	err = dbClient.Cfg.Middleware.SetBinary("python")
	Expect(err).To(BeNil())

	err = dbClient.Cfg.Middleware.SetScript(pythonReflectBody)
	Expect(err).To(BeNil())

	resp := dbClient.processRequest(req)

	// body from the request should be in response body, instead of server's response
	responseBody, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	Expect(err).To(BeNil())
	Expect(string(responseBody)).To(Equal(string(requestBody)))

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

		dbClient.Save(req, resp)
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

func TestRequestResponsePairEncodeDecode(t *testing.T) {
	RegisterTestingT(t)

	resp := models.ResponseDetails{
		Status: 200,
		Body:   "body here",
	}

	pair := models.RequestResponsePair{Response: resp}

	pairBytes, err := pair.Encode()
	Expect(err).To(BeNil())

	pairFromBytes, err := models.NewRequestResponsePairFromBytes(pairBytes)
	Expect(err).To(BeNil())
	Expect(pairFromBytes.Response.Body).To(Equal(resp.Body))
	Expect(pairFromBytes.Response.Status).To(Equal(resp.Status))
}

func TestRequestResponsePairEncodeEmpty(t *testing.T) {
	RegisterTestingT(t)

	pair := models.RequestResponsePair{}

	pairBytes, err := pair.Encode()
	Expect(err).To(BeNil())

	_, err = models.NewRequestResponsePairFromBytes(pairBytes)
	Expect(err).To(BeNil())
}

func TestDecodeRandomBytes(t *testing.T) {
	RegisterTestingT(t)

	bytes := []byte("some random stuff here")
	_, err := models.NewRequestResponsePairFromBytes(bytes)
	Expect(err).ToNot(BeNil())
}

func TestModifyRequest(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	dbClient.SetMode("modify")

	err := dbClient.Cfg.Middleware.SetBinary("python")
	Expect(err).To(BeNil())

	err = dbClient.Cfg.Middleware.SetScript(pythonReflectBody)
	Expect(err).To(BeNil())

	req, err := http.NewRequest("GET", "http://very-interesting-website.com/q=123", nil)
	Expect(err).To(BeNil())

	response := dbClient.processRequest(req)

	// response should be changed to 201
	Expect(response.StatusCode).To(Equal(http.StatusCreated))

}

func TestModifyRequestWODestination(t *testing.T) {
	RegisterTestingT(t)

	// tests modify mode but uses different middleware to not supply destination
	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	dbClient.SetMode("modify")

	err := dbClient.Cfg.Middleware.SetBinary("python")
	Expect(err).To(BeNil())

	err = dbClient.Cfg.Middleware.SetScript(pythonModifyResponse)
	Expect(err).To(BeNil())

	req, err := http.NewRequest("GET", "http://very-interesting-website.com/q=123", nil)
	Expect(err).To(BeNil())

	response := dbClient.processRequest(req)

	// response should be changed to 201
	Expect(response.StatusCode).To(Equal(http.StatusCreated))

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

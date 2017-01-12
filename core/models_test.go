package hoverfly

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/SpectoLabs/hoverfly/core/matching"
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

// TestRequestFingerprint tests whether we get correct request ID
func TestRequestFingerprint(t *testing.T) {
	RegisterTestingT(t)

	req, err := http.NewRequest("GET", "http://example.com", nil)
	Expect(err).To(BeNil())

	fp := matching.GetRequestFingerprint(req, []byte(""), false)

	Expect(fp).To(Equal("92a65ed4ca2b7100037a4cba9afd15ea"))
}

// TestRequestFingerprintBody tests where request body is also used to create unique request ID
func TestRequestFingerprintBody(t *testing.T) {
	RegisterTestingT(t)

	req, err := http.NewRequest("GET", "http://example.com", nil)
	Expect(err).To(BeNil())

	fp := matching.GetRequestFingerprint(req, []byte("some huge XML or JSON here"), false)

	Expect(fp).To(Equal("b3918a54eb6e42652e29e14c21ba8f81"))
}

func TestScheme(t *testing.T) {
	RegisterTestingT(t)

	req, err := http.NewRequest("GET", "http://example.com", nil)
	Expect(err).To(BeNil())

	originalFp := matching.GetRequestFingerprint(req, []byte(""), false)

	httpsReq, err := http.NewRequest("GET", "https://example.com", nil)
	Expect(err).To(BeNil())

	newFp := matching.GetRequestFingerprint(httpsReq, []byte(""), false)

	// fingerprint should be the same
	Expect(originalFp).To(Equal(newFp))
}

func TestDeleteAllRecords(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	// inserting some payloads
	for i := 0; i < 5; i++ {
		dbClient.Save(&models.RequestDetails{
			Destination: "delete_all_records.com",
			Query:       fmt.Sprintf("q=%i", i),
		}, &models.ResponseDetails{
			Status: 200,
			Body:   "ok",
		})
	}
	err := dbClient.RequestCache.DeleteData()
	Expect(err).To(BeNil())

	count, err := dbClient.RequestCache.RecordsCount()
	Expect(err).To(BeNil())

	Expect(count).To(BeZero())

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

func TestJSONMinifier(t *testing.T) {
	RegisterTestingT(t)

	// body can be nil here, it's not reading it from request anyway
	req, err := http.NewRequest("GET", "http://example.com", nil)
	Expect(err).To(BeNil())
	req.Header.Add("Content-Type", "application/json")

	fpOne := matching.GetRequestFingerprint(req, []byte(`{"foo": "bar"}`), false)
	fpTwo := matching.GetRequestFingerprint(req, []byte(`{     "foo":           "bar"}`), false)

	Expect(fpOne).To(Equal(fpTwo))
}

func TestJSONMinifierWOHeader(t *testing.T) {
	RegisterTestingT(t)

	// body can be nil here, it's not reading it from request anyway
	req, err := http.NewRequest("GET", "http://example.com", nil)
	Expect(err).To(BeNil())

	// application/json header is not set, shouldn't be equal
	fpOne := matching.GetRequestFingerprint(req, []byte(`{"foo": "bar"}`), false)
	fpTwo := matching.GetRequestFingerprint(req, []byte(`{     "foo":           "bar"}`), false)

	Expect(fpOne).ToNot(Equal(fpTwo))
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
	RegisterTestingT(t)

	// body can be nil here, it's not reading it from request anyway
	req, err := http.NewRequest("GET", "http://example.com", nil)
	Expect(err).To(BeNil())

	req.Header.Add("Content-Type", "application/xml")

	fpOne := matching.GetRequestFingerprint(req, []byte(xmlBody), false)
	fpTwo := matching.GetRequestFingerprint(req, []byte(xmlBodyTwo), false)
	Expect(fpOne).To(Equal(fpTwo))
}

func TestXMLMinifierWOHeader(t *testing.T) {
	RegisterTestingT(t)

	// body can be nil here, it's not reading it from request anyway
	req, err := http.NewRequest("GET", "http://example.com", nil)
	Expect(err).To(BeNil())

	// application/xml header is not set, shouldn't be equal
	fpOne := matching.GetRequestFingerprint(req, []byte(xmlBody), false)
	fpTwo := matching.GetRequestFingerprint(req, []byte(xmlBodyTwo), false)
	Expect(fpOne).ToNot(Equal(fpTwo))
}

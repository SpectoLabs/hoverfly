package models

import (
	"os"
	"testing"
	. "github.com/onsi/gomega"
	"compress/gzip"
	"bytes"
	"io/ioutil"
)

func TestConvertToResponseDetailsView_WithPlainTextResponseDetails(t *testing.T) {
	RegisterTestingT(t)

	statusCode := 200
	body := "hello_world"
	headers := map[string][]string{"test_header": []string{"true"}}

	originalResp := ResponseDetails{Status: statusCode, Body: body, Headers: headers}

	respView := originalResp.ConvertToResponseDetailsView()

	Expect(respView.Status).To(Equal(statusCode))
	Expect(respView.Headers).To(Equal(headers))

	Expect(respView.EncodedBody).To(Equal(false))
	Expect(respView.Body).To(Equal(body))
}

func TestConvertToResponseDetailsView_WithGzipContentEncodedHeader(t *testing.T) {
	RegisterTestingT(t)

	originalBody := "hello_world"

	statusCode := 200
	body := GzipString(originalBody)
	headers := map[string][]string{"Content-Encoding": []string{"gzip"}}

	originalResp := ResponseDetails{Status: statusCode, Body: body, Headers:headers}

	respView := originalResp.ConvertToResponseDetailsView()

	Expect(respView.Status).To(Equal(statusCode))
	Expect(respView.Headers).To(Equal(headers))

	Expect(respView.EncodedBody).To(Equal(true))
	Expect(respView.Body).NotTo(Equal(body))
	Expect(respView.Body).NotTo(Equal(originalBody))

	base64EncodedBody := "H4sIAAAJbogA/w=="

	Expect(respView.Body).To(Equal(base64EncodedBody))
}

func TestConvertToResponseDetailsView_WithDeflateContentEncodedHeader(t *testing.T) {
	RegisterTestingT(t)

	originalBody := "this_should_be_encoded_but_its_not_important"

	statusCode := 200
	headers := map[string][]string{"Content-Encoding": []string{"deflate"}}

	originalResp := ResponseDetails{Status: statusCode, Body: originalBody, Headers:headers}

	respView := originalResp.ConvertToResponseDetailsView()

	Expect(respView.Status).To(Equal(statusCode))
	Expect(respView.Headers).To(Equal(headers))

	Expect(respView.EncodedBody).To(Equal(true))
	Expect(respView.Body).NotTo(Equal(originalBody))

	base64EncodedBody := "dGhpc19zaG91bGRfYmVfZW5jb2RlZF9idXRfaXRzX25vdF9pbXBvcnRhbnQ="

	Expect(respView.Body).To(Equal(base64EncodedBody))
}

func TestConvertToResponseDetailsView_WithImageBody(t *testing.T) {
	RegisterTestingT(t)

	imageUri := "/testdata/1x1.png"

	file, _ := os.Open(".." + imageUri)
	defer file.Close()

	originalImageBytes, _ := ioutil.ReadAll(file)

	originalResp := ResponseDetails{
		Status: 200,
		Body: string(originalImageBytes),
	}

	respView := originalResp.ConvertToResponseDetailsView()

	base64EncodedBody := "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAAAAAA6fptVAAAACklEQVR4nGP6DwABBQECz6AuzQAAAABJRU5ErkJggg=="
	Expect(respView).To(Equal(ResponseDetailsView{
		Status: 200,
		Body: base64EncodedBody,
		EncodedBody: true,
	}))
}
func TestPayload_ConvertToPayloadView_WithPlainTextResponse(t *testing.T) {
	RegisterTestingT(t)

	respBody := "hello_world"

	originalPayload := Payload{
		Response: ResponseDetails{
			Status: 200,
			Body: respBody,
			Headers: map[string][]string{"test_header": []string{"true"}}},
		Request: RequestDetails{
			Path: "/",
			Method: "GET",
			Destination: "/",
			Scheme: "scheme",
			Query: "",
			Body: "",
			RemoteAddr: "localhost",
			Headers: map[string][]string{"test_header": []string{"true"}}},
	}

	payloadView := originalPayload.ConvertToPayloadView()

	Expect(*payloadView).To(Equal(PayloadView{
		Response: ResponseDetailsView{
			Status: 200,
			Body: respBody,
			Headers: map[string][]string{"test_header": []string{"true"}},
			EncodedBody: false},
		Request: RequestDetailsView{
			Path: "/",
			Method: "GET",
			Destination: "/",
			Scheme: "scheme",
			Query: "",
			Body: "",
			RemoteAddr: "localhost",
			Headers: map[string][]string{"test_header": []string{"true"}}},
	}))
}

func TestPayload_ConvertToPayloadView_WithGzippedResponse(t *testing.T) {
	RegisterTestingT(t)

	originalPayload := Payload{
		Response: ResponseDetails{
			Status: 200,
			Body: GzipString("hello_world"),
			Headers: map[string][]string{"Content-Encoding": []string{"gzip"}}},
		Request: RequestDetails{
			Path: "/",
			Method: "GET",
			Destination: "/",
			Scheme: "scheme",
			Query: "",
			Body: "",
			RemoteAddr: "localhost",
			Headers: map[string][]string{"Content-Encoding": []string{"gzip"}},
		},
	}

	payloadView := originalPayload.ConvertToPayloadView()

	Expect(*payloadView).To(Equal(PayloadView{
		Response: ResponseDetailsView{
			Status: 200,
			Body: "H4sIAAAJbogA/w==",
			Headers: map[string][]string{"Content-Encoding": []string{"gzip"}},
			EncodedBody: true},
		Request: RequestDetailsView{
			Path: "/",
			Method: "GET",
			Destination: "/",
			Scheme: "scheme",
			Query: "",
			Body: "",
			RemoteAddr: "localhost",
			Headers: map[string][]string{"Content-Encoding": []string{"gzip"}},
		},
	}))
}

func TestRequestDetailsView_ConvertToRequestDetails(t *testing.T) {
	RegisterTestingT(t)

	requestDetailsView := RequestDetailsView{
		Path: "/",
		Method: "GET",
		Destination: "/",
		Scheme: "scheme",
		Query: "", Body: "",
		RemoteAddr: "localhost",
		Headers: map[string][]string{"Content-Encoding": []string{"gzip"}}}

	requestDetails := requestDetailsView.ConvertToRequestDetails()

	Expect(requestDetails.Path).To(Equal(requestDetailsView.Path))
	Expect(requestDetails.Method).To(Equal(requestDetailsView.Method))
	Expect(requestDetails.Destination).To(Equal(requestDetailsView.Destination))
	Expect(requestDetails.Scheme).To(Equal(requestDetailsView.Scheme))
	Expect(requestDetails.Query).To(Equal(requestDetailsView.Query))
	Expect(requestDetails.RemoteAddr).To(Equal(requestDetailsView.RemoteAddr))
	Expect(requestDetails.Headers).To(Equal(requestDetailsView.Headers))
}

func TestRequestDetails_ConvertToRequestDetailsView(t *testing.T) {
	RegisterTestingT(t)

	requestDetails := RequestDetails{
		Path: "/",
		Method: "GET",
		Destination: "/",
		Scheme: "scheme",
		Query: "", Body: "",
		RemoteAddr: "localhost",
		Headers: map[string][]string{"Content-Encoding": []string{"gzip"}}}

	requestDetailsView := requestDetails.ConvertToRequestDetailsView()

	Expect(requestDetailsView.Path).To(Equal(requestDetails.Path))
	Expect(requestDetailsView.Method).To(Equal(requestDetails.Method))
	Expect(requestDetailsView.Destination).To(Equal(requestDetails.Destination))
	Expect(requestDetailsView.Scheme).To(Equal(requestDetails.Scheme))
	Expect(requestDetailsView.Query).To(Equal(requestDetails.Query))
	Expect(requestDetailsView.RemoteAddr).To(Equal(requestDetails.RemoteAddr))
	Expect(requestDetailsView.Headers).To(Equal(requestDetails.Headers))
}

// Helper function for gzipping strings
func GzipString(s string) (string) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	gz.Write([]byte(s))
	return b.String()
}
package models_test

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"os"
	"testing"

	"net/http"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/models"
	. "github.com/SpectoLabs/hoverfly/core/util"
	. "github.com/onsi/gomega"
)

func TestResponseDetails_ConvertToResponseDetailsView_WithPlainTextResponseDetails(t *testing.T) {
	RegisterTestingT(t)

	statusCode := 200
	body := "hello_world"
	headers := map[string][]string{"test_header": {"true"}}

	originalResp := models.ResponseDetails{Status: statusCode, Body: body, Headers: headers}

	respView := originalResp.ConvertToResponseDetailsView()

	Expect(respView.Status).To(Equal(statusCode))
	Expect(respView.Headers).To(Equal(headers))

	Expect(respView.EncodedBody).To(Equal(false))
	Expect(respView.Body).To(Equal(body))
}

func TestResponseDetails_ConvertToResponseDetailsView_WithGzipContentEncodedHeader(t *testing.T) {
	RegisterTestingT(t)

	originalBody := "hello_world"

	statusCode := 200
	body := GzipString(originalBody)
	headers := map[string][]string{"Content-Encoding": {"gzip"}}

	originalResp := models.ResponseDetails{Status: statusCode, Body: body, Headers: headers}

	respView := originalResp.ConvertToResponseDetailsView()

	Expect(respView.Status).To(Equal(statusCode))
	Expect(respView.Headers).To(Equal(headers))

	Expect(respView.EncodedBody).To(Equal(true))
	Expect(respView.Body).NotTo(Equal(body))
	Expect(respView.Body).NotTo(Equal(originalBody))

	base64EncodedBody := "H4sIAAAAAAAA/w=="

	Expect(respView.Body).To(Equal(base64EncodedBody))
}

func TestResponseDetails_ConvertToResponseDetailsView_WithDeflateContentEncodedHeader(t *testing.T) {
	RegisterTestingT(t)

	originalBody := "this_should_be_encoded_but_its_not_important"

	statusCode := 200
	headers := map[string][]string{"Content-Encoding": {"deflate"}}

	originalResp := models.ResponseDetails{Status: statusCode, Body: originalBody, Headers: headers}

	respView := originalResp.ConvertToResponseDetailsView()

	Expect(respView.Status).To(Equal(statusCode))
	Expect(respView.Headers).To(Equal(headers))

	Expect(respView.EncodedBody).To(Equal(true))
	Expect(respView.Body).NotTo(Equal(originalBody))

	base64EncodedBody := "dGhpc19zaG91bGRfYmVfZW5jb2RlZF9idXRfaXRzX25vdF9pbXBvcnRhbnQ="

	Expect(respView.Body).To(Equal(base64EncodedBody))
}

func TestResponseDetails_ConvertToResponseDetailsView_WithImageBody(t *testing.T) {
	RegisterTestingT(t)

	imageUri := "/testdata/1x1.png"

	file, _ := os.Open("../../functional-tests/core" + imageUri)
	defer file.Close()

	originalImageBytes, _ := ioutil.ReadAll(file)

	originalResp := models.ResponseDetails{
		Status: 200,
		Body:   string(originalImageBytes),
	}

	respView := originalResp.ConvertToResponseDetailsView()

	base64EncodedBody := "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAAAAAA6fptVAAAACklEQVR4nGP6DwABBQECz6AuzQAAAABJRU5ErkJggg=="
	Expect(respView).To(Equal(v2.ResponseDetailsView{
		Status:      200,
		Body:        base64EncodedBody,
		EncodedBody: true,
	}))
}

func TestRequestResponsePair_ConvertToRequestResponsePairView_WithPlainTextResponse(t *testing.T) {
	RegisterTestingT(t)

	respBody := "hello_world"

	requestResponsePair := models.RequestResponsePair{
		Response: models.ResponseDetails{
			Status:  200,
			Body:    respBody,
			Headers: map[string][]string{"test_header": {"true"}}},
		Request: models.RequestDetails{
			Path:        "/",
			Method:      "GET",
			Destination: "/",
			Scheme:      "scheme",
			Query:       map[string][]string{},
			Body:        "",
			Headers:     map[string][]string{"test_header": {"true"}}},
	}

	pairView := requestResponsePair.ConvertToRequestResponsePairView()

	Expect(pairView).To(Equal(v2.RequestResponsePairViewV1{
		Response: v2.ResponseDetailsView{
			Status:      200,
			Body:        respBody,
			Headers:     map[string][]string{"test_header": {"true"}},
			EncodedBody: false},
		Request: v2.RequestDetailsView{
			Path:        StringToPointer("/"),
			Method:      StringToPointer("GET"),
			Destination: StringToPointer("/"),
			Scheme:      StringToPointer("scheme"),
			Query:       StringToPointer(""),
			QueryMap:    map[string][]string{},
			Body:        StringToPointer(""),
			Headers:     map[string][]string{"test_header": {"true"}}},
	}))
}

func TestRequestResponsePair_ConvertToRequestResponsePairView_WithGzippedResponse(t *testing.T) {
	RegisterTestingT(t)

	requestResponsePair := models.RequestResponsePair{
		Response: models.ResponseDetails{
			Status:  200,
			Body:    GzipString("hello_world"),
			Headers: map[string][]string{"Content-Encoding": {"gzip"}}},
		Request: models.RequestDetails{
			Path:        "/",
			Method:      "GET",
			Destination: "/",
			Scheme:      "scheme",
			Query:       map[string][]string{},
			Body:        "",
			Headers:     map[string][]string{"Content-Encoding": {"gzip"}},
		},
	}

	pairView := requestResponsePair.ConvertToRequestResponsePairView()

	Expect(pairView).To(Equal(v2.RequestResponsePairViewV1{
		Response: v2.ResponseDetailsView{
			Status:      200,
			Body:        "H4sIAAAAAAAA/w==",
			Headers:     map[string][]string{"Content-Encoding": {"gzip"}},
			EncodedBody: true},
		Request: v2.RequestDetailsView{
			Path:        StringToPointer("/"),
			Method:      StringToPointer("GET"),
			Destination: StringToPointer("/"),
			Scheme:      StringToPointer("scheme"),
			Query:       StringToPointer(""),
			QueryMap:    map[string][]string{},
			Body:        StringToPointer(""),
			Headers:     map[string][]string{"Content-Encoding": {"gzip"}},
		},
	}))
}

func TestRequestDetails_ConvertToRequestDetailsView(t *testing.T) {
	RegisterTestingT(t)

	requestDetails := models.RequestDetails{
		Path:        "/",
		Method:      "GET",
		Destination: "/",
		Scheme:      "scheme",
		Query:       map[string][]string{},
		Body:        "",
		Headers:     map[string][]string{"Content-Encoding": {"gzip"}}}

	requestDetailsView := requestDetails.ConvertToRequestDetailsView()

	Expect(requestDetailsView.Path).To(Equal(StringToPointer(requestDetails.Path)))
	Expect(requestDetailsView.Method).To(Equal(StringToPointer(requestDetails.Method)))
	Expect(requestDetailsView.Destination).To(Equal(StringToPointer(requestDetails.Destination)))
	Expect(requestDetailsView.Scheme).To(Equal(StringToPointer(requestDetails.Scheme)))
	Expect(requestDetailsView.Query).To(Equal(StringToPointer("")))
	Expect(requestDetailsView.Headers).To(Equal(requestDetails.Headers))
}

// Helper function for gzipping strings
func GzipString(s string) string {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	gz.Write([]byte(s))
	return b.String()
}

func Test_NewRequestDetailsFromHttpRequest_SortsQueryString(t *testing.T) {
	RegisterTestingT(t)
	request, _ := http.NewRequest("GET", "http://test.org/?a=b&a=a", nil)
	requestDetails, err := models.NewRequestDetailsFromHttpRequest(request)
	Expect(err).To(BeNil())

	Expect(requestDetails.Query["a"]).To(ContainElement("a"))
	Expect(requestDetails.Query["a"]).To(ContainElement("b"))
	Expect(requestDetails.QueryString()).To(Equal("a=a&a=b"))
}

func Test_NewRequestDetailsFromHttpRequest_StripsArbitaryGolangColonEscaping(t *testing.T) {
	RegisterTestingT(t)
	request, _ := http.NewRequest("GET", "http://test.org/?a=b:c", nil)
	requestDetails, err := models.NewRequestDetailsFromHttpRequest(request)
	Expect(err).To(BeNil())

	Expect(requestDetails.Query["a"]).To(ContainElement("b:c"))
}

func Test_NewRequestDetailsFromHttpRequest_UsesRawPathIfAvailable(t *testing.T) {
	RegisterTestingT(t)
	request, _ := http.NewRequest("GET", "http://test.org/hoverfly%20rocks", nil)
	request.URL.RawPath = "/hoverfly%20rocks"
	requestDetails, err := models.NewRequestDetailsFromHttpRequest(request)
	Expect(err).To(BeNil())

	Expect(requestDetails.Path).To(Equal("/hoverfly%20rocks"))
}

func Test_NewRequestDetailsFromHttpRequest_LowerCaseDestination(t *testing.T) {
	RegisterTestingT(t)

	request, _ := http.NewRequest("GET", "http://TEST.ORG/?a=b&a=a", nil)
	requestDetails, err := models.NewRequestDetailsFromHttpRequest(request)
	Expect(err).To(BeNil())

	Expect(requestDetails.Destination).To(Equal("test.org"))
}


func Test_NewRequestDetailsFromHttpRequest_HandleNonAbsoluteURL(t *testing.T) {
	RegisterTestingT(t)
	request, _ := http.NewRequest("GET", "/hello", nil)
	requestDetails, err := models.NewRequestDetailsFromHttpRequest(request)
	Expect(err).To(BeNil())

	Expect(requestDetails.Scheme).To(Equal("http"))
	Expect(requestDetails.Path).To(Equal("/hello"))
}

func TestRequestResponsePairView_ConvertToRequestResponsePairWithoutEncoding(t *testing.T) {
	RegisterTestingT(t)

	view := v2.RequestResponsePairViewV1{
		Request: v2.RequestDetailsView{
			Path:        StringToPointer("A"),
			Method:      StringToPointer("A"),
			Destination: StringToPointer("A"),
			Scheme:      StringToPointer("A"),
			Query:       StringToPointer("A"),
			Body:        StringToPointer("A"),
			Headers: map[string][]string{
				"A": {"B"},
				"C": {"D"},
			},
		},
		Response: v2.ResponseDetailsView{
			Status:      1,
			Body:        "1",
			EncodedBody: false,
			Headers: map[string][]string{
				"1": {"2"},
				"3": {"4"},
			},
		},
	}

	requestResponsePair := models.NewRequestResponsePairFromRequestResponsePairView(view)

	Expect(requestResponsePair).To(Equal(models.RequestResponsePair{
		Request: models.RequestDetails{
			Path:        "A",
			Method:      "A",
			Destination: "A",
			Scheme:      "A",
			Query: map[string][]string{
				"A": {""},
			},
			Body: "A",
			Headers: map[string][]string{
				"A": {"B"},
				"C": {"D"},
			},
		},
		Response: models.ResponseDetails{
			Status: 1,
			Body:   "1",
			Headers: map[string][]string{
				"1": {"2"},
				"3": {"4"},
			},
		},
	}))
}

func TestRequestResponsePairView_ConvertToRequestResponsePairWithEncoding(t *testing.T) {
	RegisterTestingT(t)

	view := v2.RequestResponsePairViewV1{
		Request: v2.RequestDetailsView{
			Query: StringToPointer("somehthing=something"),
		},
		Response: v2.ResponseDetailsView{
			Body:        "ZW5jb2RlZA==",
			EncodedBody: true,
		},
	}

	pair := models.NewRequestResponsePairFromRequestResponsePairView(view)

	Expect(pair.Response.Body).To(Equal("encoded"))
}

func TestRequestDetailsView_ConvertToRequestDetails(t *testing.T) {
	RegisterTestingT(t)

	requestDetailsView := v2.RequestDetailsView{
		Path:        StringToPointer("/"),
		Method:      StringToPointer("GET"),
		Destination: StringToPointer("/"),
		Scheme:      StringToPointer("scheme"),
		Query:       StringToPointer(""),
		Body:        StringToPointer(""),
		Headers:     map[string][]string{"Content-Encoding": {"gzip"}}}

	requestDetails := models.NewRequestDetailsFromRequest(requestDetailsView)

	Expect(requestDetails.Path).To(Equal(*requestDetailsView.Path))
	Expect(requestDetails.Method).To(Equal(*requestDetailsView.Method))
	Expect(requestDetails.Destination).To(Equal(*requestDetailsView.Destination))
	Expect(requestDetails.Scheme).To(Equal(*requestDetailsView.Scheme))
	Expect(requestDetails.Query).To(Equal(map[string][]string{}))
	Expect(requestDetails.Headers).To(Equal(requestDetailsView.Headers))
}

func Test_RequestDetails_Hash_ItHashes(t *testing.T) {
	RegisterTestingT(t)

	unit := models.RequestDetails{
		Method:      "GET",
		Scheme:      "http",
		Destination: "test.com",
		Path:        "/testing",
		Query: map[string][]string{
			"query": {"true"},
		},
	}

	hashedUnit := unit.Hash()

	Expect(hashedUnit).To(Equal("70c4fd58c2db41298071ea0446af0793"))
}

func Test_RequestDetails_Hash_TheHashIgnoresHeaders(t *testing.T) {
	RegisterTestingT(t)

	unit := models.RequestDetails{
		Method:      "GET",
		Scheme:      "http",
		Destination: "test.com",
		Path:        "/testing",
		Query: map[string][]string{
			"query": {"true"},
		},
		Headers: map[string][]string{"Content-Encoding": {"gzip"}},
	}

	hashedUnit := unit.Hash()

	Expect(hashedUnit).To(Equal("70c4fd58c2db41298071ea0446af0793"))
}

func Test_RequestDetails_Hash_TheHashIncludesTheBody(t *testing.T) {
	RegisterTestingT(t)

	unit := models.RequestDetails{
		Method:      "GET",
		Scheme:      "http",
		Destination: "test.com",
		Path:        "/testing",
		Query: map[string][]string{
			"query": {"true"},
		},
		Body: "tidy text",
	}

	hashedUnit := unit.Hash()

	Expect(hashedUnit).To(Equal("51834bfe5334158be38ef5209f2b8e29"))
}

func Test_RequestDetails_QueryString_ConvertsMapToString(t *testing.T) {
	RegisterTestingT(t)

	requestDetails := models.RequestDetails{
		Query: map[string][]string{
			"test": {"value"},
		},
	}

	Expect(requestDetails.QueryString()).To(Equal("test=value"))
}

func Test_RequestDetails_QueryString_HandlesMultipleValuesAsOneValue(t *testing.T) {
	RegisterTestingT(t)

	requestDetails := models.RequestDetails{
		Query: map[string][]string{
			"test": {"value,value2"},
		},
	}

	Expect(requestDetails.QueryString()).To(Equal("test=value,value2"))
}

func Test_RequestDetails_QueryString_HandlesSpaces(t *testing.T) {
	RegisterTestingT(t)

	requestDetails := models.RequestDetails{
		Query: map[string][]string{
			"test": {"val ue"},
		},
	}

	Expect(requestDetails.QueryString()).To(Equal("test=val ue"))
}

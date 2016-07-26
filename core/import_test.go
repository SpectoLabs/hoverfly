package hoverfly

import (
	. "github.com/onsi/gomega"
	"encoding/base64"
	"github.com/SpectoLabs/hoverfly/core/cache"
	"github.com/SpectoLabs/hoverfly/core/matching"
	"github.com/SpectoLabs/hoverfly/core/models"
	"io/ioutil"
	"os"
	"testing"
)

func TestIsURLHTTP(t *testing.T) {
	url := "http://somehost.com"

	b := isURL(url)
	Expect(b).To(BeTrue())
}

func TestIsURLEmpty(t *testing.T) {
	b := isURL("")
	Expect(b).To(BeFalse())
}

func TestIsURLHTTPS(t *testing.T) {
	url := "https://somehost.com"

	b := isURL(url)
	Expect(b).To(BeTrue())
}

func TestIsURLWrong(t *testing.T) {
	url := "somehost.com"

	b := isURL(url)
	Expect(b).To(BeFalse())
}

func TestIsURLWrongTLD(t *testing.T) {
	url := "http://somehost."

	b := isURL(url)
	Expect(b).To(BeFalse())
}

func TestFileExists(t *testing.T) {
	fp := "examples/exports/readthedocs.json"

	ex, err := exists(fp)
	Expect(err).To(BeNil())
	Expect(ex).To(BeTrue())
}

func TestFileDoesNotExist(t *testing.T) {
	fp := "shouldnotbehere.yaml"

	ex, err := exists(fp)
	Expect(err).To(BeNil())
	Expect(ex).To(BeFalse())
}

func TestImportFromDisk(t *testing.T) {
	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()

	err := dbClient.Import("examples/exports/readthedocs.json")
	Expect(err).To(BeNil())

	recordsCount, err := dbClient.RequestCache.RecordsCount()
	Expect(err).To(BeNil())

	Expect(recordsCount).To(Equal(5))
}

func TestImportFromDiskBlankPath(t *testing.T) {
	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()

	err := dbClient.ImportFromDisk("")
	Expect(err).ToNot(BeNil())
}

func TestImportFromDiskWrongJson(t *testing.T) {
	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()

	err := dbClient.ImportFromDisk("examples/exports/README.md")
	Expect(err).ToNot(BeNil())
}

func TestImportFromURL(t *testing.T) {
	// reading file and preparing json payload
	payloadsFile, err := os.Open("examples/exports/readthedocs.json")
	Expect(err).To(BeNil())
	bts, err := ioutil.ReadAll(payloadsFile)
	Expect(err).To(BeNil())

	// pretending this is the endpoint with given json
	server, dbClient := testTools(200, string(bts))
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()

	// importing payloads
	err = dbClient.Import("http://thiswillbeintercepted.json")
	Expect(err).To(BeNil())

	recordsCount, err := dbClient.RequestCache.RecordsCount()
	Expect(err).To(BeNil())
	Expect(recordsCount).To(Equal(5))
}

func TestImportFromURLHTTPFail(t *testing.T) {
	// this tests simulates unreachable server
	server, dbClient := testTools(200, `this shouldn't matter anyway`)
	// closing it immediately
	server.Close()
	defer dbClient.RequestCache.DeleteData()

	err := dbClient.ImportFromURL("somepath")
	Expect(err).ToNot(BeNil())
}

func TestImportFromURLMalformedJSON(t *testing.T) {
	// testing behaviour when there is no json on the other end
	server, dbClient := testTools(200, `i am not json :(`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()

	// importing payloads
	err := dbClient.Import("http://thiswillbeintercepted.json")
	// we should get error
	Expect(err).ToNot(BeNil())
}

func TestImportPayloads_CanImportASinglePayload(t *testing.T) {
	cache := cache.NewInMemoryCache()
	cfg := Configuration{Webserver: false}
	requestMatcher := matching.RequestMatcher{RequestCache: cache, Webserver: &cfg.Webserver}
	hv := Hoverfly{RequestCache: cache, Cfg: &cfg, RequestMatcher: requestMatcher}

	RegisterTestingT(t)

	originalPayload := models.PayloadView{
		Response: models.ResponseDetailsView{
			Status:      200,
			Body:        "hello_world",
			EncodedBody: false,
			Headers:     map[string][]string{"Content-Type": []string{"text/plain"}}},
		Request: models.RequestDetailsView{
			Path:        "/",
			Method:      "GET",
			Destination: "/",
			Scheme:      "scheme",
			Query:       "", Body: "",
			Headers: map[string][]string{"Hoverfly": []string{"testing"}}}}

	hv.ImportPayloads([]models.PayloadView{originalPayload})
	value, _ := cache.Get([]byte("9b114df98da7f7e2afdc975883dab4f2"))
	decodedPayload, _ := models.NewPayloadFromBytes(value)
	Expect(*decodedPayload).To(Equal(models.Payload{
		Response: models.ResponseDetails{
			Status:  200,
			Body:    "hello_world",
			Headers: map[string][]string{"Content-Type": []string{"text/plain"}},
		},
		Request: models.RequestDetails{
			Path:        "/",
			Method:      "GET",
			Destination: "/",
			Scheme:      "scheme",
			Query:       "", Body: "",
			Headers: map[string][]string{
				"Content-Type": []string{"text/plain; charset=utf-8"},
				"Hoverfly":     []string{"testing"},
			},
		},
	}))
}

func TestImportPayloads_CanImportAMultiplePayload(t *testing.T) {
	cache := cache.NewInMemoryCache()
	cfg := Configuration{Webserver: false}
	requestMatcher := matching.RequestMatcher{RequestCache: cache, Webserver: &cfg.Webserver}
	hv := Hoverfly{RequestCache: cache, Cfg: &cfg, RequestMatcher: requestMatcher}

	RegisterTestingT(t)

	originalPayload1 := models.PayloadView{
		Response: models.ResponseDetailsView{
			Status:      200,
			Body:        "hello_world",
			EncodedBody: false,
			Headers:     map[string][]string{"Hoverfly": []string{"testing"}},
		},
		Request: models.RequestDetailsView{
			Path:        "/",
			Method:      "GET",
			Destination: "/",
			Scheme:      "scheme",
			Query:       "", Body: "",
			Headers: map[string][]string{"Hoverfly": []string{"testing"}}}}

	originalPayload2 := originalPayload1

	originalPayload2.Request.Path = "/new/path"

	originalPayload3 := originalPayload1

	originalPayload3.Request.Path = "/newer/path"

	hv.ImportPayloads([]models.PayloadView{originalPayload1, originalPayload2, originalPayload3})
	value, err := cache.Get([]byte("9b114df98da7f7e2afdc975883dab4f2"))
	Expect(err).To(BeNil())
	decodedPayload1, err := models.NewPayloadFromBytes(value)
	Expect(err).To(BeNil())
	Expect(*decodedPayload1).To(Equal(originalPayload1.ConvertToPayload()))

	value, err = cache.Get([]byte("9c03e4af1f30542ff079a712bddad602"))
	Expect(err).To(BeNil())
	decodedPayload2, err := models.NewPayloadFromBytes(value)
	Expect(err).To(BeNil())
	Expect(*decodedPayload2).To(Equal(originalPayload2.ConvertToPayload()))

	value, err = cache.Get([]byte("fd099332afee48101edb7441b098cd4a"))
	Expect(err).To(BeNil())
	decodedPayload3, err := models.NewPayloadFromBytes(value)
	Expect(err).To(BeNil())
	Expect(*decodedPayload3).To(Equal(originalPayload3.ConvertToPayload()))
}

// Helper function for base64 encoding
func base64String(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func TestImportPayloads_CanImportASingleBase64EncodedPayload(t *testing.T) {
	cache := cache.NewInMemoryCache()
	cfg := Configuration{Webserver: false}
	requestMatcher := matching.RequestMatcher{RequestCache: cache, Webserver: &cfg.Webserver}
	hv := Hoverfly{RequestCache: cache, Cfg: &cfg, RequestMatcher: requestMatcher}

	RegisterTestingT(t)

	encodedPayload := models.PayloadView{
		Response: models.ResponseDetailsView{
			Status:      200,
			Body:        base64String("hello_world"),
			EncodedBody: true,
			Headers:     map[string][]string{"Content-Encoding": []string{"gzip"}}},
		Request: models.RequestDetailsView{
			Path:        "/",
			Method:      "GET",
			Destination: "/",
			Scheme:      "scheme",
			Query:       "", Body: "",
			Headers: map[string][]string{"Hoverfly": []string{"testing"}}}}

	hv.ImportPayloads([]models.PayloadView{encodedPayload})

	value, err := cache.Get([]byte("9b114df98da7f7e2afdc975883dab4f2"))
	Expect(err).To(BeNil())

	decodedPayload, err := models.NewPayloadFromBytes(value)
	Expect(err).To(BeNil())

	Expect(decodedPayload).ToNot(Equal(models.Payload{
		Response: models.ResponseDetails{
			Status:  200,
			Body:    "hello_world",
			Headers: map[string][]string{"Content-Encoding": []string{"gzip"}}},
		Request: models.RequestDetails{
			Path:        "/",
			Method:      "GET",
			Destination: "/",
			Scheme:      "scheme",
			Query:       "", Body: "",
			Headers: map[string][]string{"Hoverfly": []string{"testing"}}}}))
}

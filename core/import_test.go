package hoverfly

import (
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/SpectoLabs/hoverfly/core/cache"
	"github.com/SpectoLabs/hoverfly/core/handlers/v1"
	"github.com/SpectoLabs/hoverfly/core/interfaces"
	"github.com/SpectoLabs/hoverfly/core/matching"
	"github.com/SpectoLabs/hoverfly/core/models"
	. "github.com/SpectoLabs/hoverfly/core/util"
	. "github.com/onsi/gomega"
)

func TestIsURLHTTP(t *testing.T) {
	RegisterTestingT(t)

	url := "http://somehost.com"

	b := isURL(url)
	Expect(b).To(BeTrue())
}

func TestIsURLEmpty(t *testing.T) {
	RegisterTestingT(t)

	b := isURL("")
	Expect(b).To(BeFalse())
}

func TestIsURLHTTPS(t *testing.T) {
	RegisterTestingT(t)

	url := "https://somehost.com"

	b := isURL(url)
	Expect(b).To(BeTrue())
}

func TestIsURLWrong(t *testing.T) {
	RegisterTestingT(t)

	url := "somehost.com"

	b := isURL(url)
	Expect(b).To(BeFalse())
}

func TestIsURLWrongTLD(t *testing.T) {
	RegisterTestingT(t)

	url := "http://somehost."

	b := isURL(url)
	Expect(b).To(BeFalse())
}

func TestFileExists(t *testing.T) {
	RegisterTestingT(t)

	fp := "examples/exports/readthedocs.json"

	ex, err := exists(fp)
	Expect(err).To(BeNil())
	Expect(ex).To(BeTrue())
}

func TestFileDoesNotExist(t *testing.T) {
	RegisterTestingT(t)

	fp := "shouldnotbehere.yaml"

	ex, err := exists(fp)
	Expect(err).To(BeNil())
	Expect(ex).To(BeFalse())
}

func TestImportFromDisk(t *testing.T) {
	RegisterTestingT(t)

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
	RegisterTestingT(t)

	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()

	err := dbClient.ImportFromDisk("")
	Expect(err).ToNot(BeNil())
}

func TestImportFromDiskWrongJson(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()

	err := dbClient.ImportFromDisk("examples/exports/README.md")
	Expect(err).ToNot(BeNil())
}

func TestImportFromURL(t *testing.T) {
	RegisterTestingT(t)

	// reading file and preparing json payload
	pairFile, err := os.Open("examples/exports/readthedocs.json")
	Expect(err).To(BeNil())
	pairFileBytes, err := ioutil.ReadAll(pairFile)
	Expect(err).To(BeNil())

	// pretending this is the endpoint with given json
	server, dbClient := testTools(200, string(pairFileBytes))
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()

	// importing payloads
	err = dbClient.Import(server.URL)
	Expect(err).To(BeNil())

	recordsCount, err := dbClient.RequestCache.RecordsCount()
	Expect(err).To(BeNil())
	Expect(recordsCount).To(Equal(5))
}

func TestImportFromURLRedirect(t *testing.T) {
	RegisterTestingT(t)

	// reading file and preparing json payload
	pairFile, err := os.Open("examples/exports/readthedocs.json")
	Expect(err).To(BeNil())
	pairFileBytes, err := ioutil.ReadAll(pairFile)
	Expect(err).To(BeNil())

	// pretending this is the endpoint with given json
	server, dbClient := testTools(200, string(pairFileBytes))
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()

	dbClient.HTTP = GetDefaultHoverflyHTTPClient(false, "")

	redirectServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Location", server.URL)
		w.WriteHeader(301)
	}))
	defer redirectServer.Close()

	// importing payloads
	err = dbClient.Import(redirectServer.URL)
	Expect(err).To(BeNil())

	recordsCount, err := dbClient.RequestCache.RecordsCount()
	Expect(err).To(BeNil())
	Expect(recordsCount).To(Equal(5))
}

func TestImportFromURLHTTPFail(t *testing.T) {
	RegisterTestingT(t)

	// this tests simulates unreachable server
	server, dbClient := testTools(200, `this shouldn't matter anyway`)
	// closing it immediately
	server.Close()
	defer dbClient.RequestCache.DeleteData()

	err := dbClient.ImportFromURL("somepath")
	Expect(err).ToNot(BeNil())
}

func TestImportFromURLMalformedJSON(t *testing.T) {
	RegisterTestingT(t)

	// testing behaviour when there is no json on the other end
	server, dbClient := testTools(200, `i am not json :(`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()

	// importing payloads
	err := dbClient.Import("http://thiswillbeintercepted.json")
	// we should get error
	Expect(err).ToNot(BeNil())
}

func TestImportRequestResponsePairs_CanImportASinglePair(t *testing.T) {
	RegisterTestingT(t)

	cache := cache.NewInMemoryCache()
	cfg := Configuration{Webserver: false}
	requestMatcher := matching.RequestMatcher{RequestCache: cache, Webserver: &cfg.Webserver}
	hv := Hoverfly{RequestCache: cache, Cfg: &cfg, RequestMatcher: requestMatcher}

	RegisterTestingT(t)

	originalPair := v1.RequestResponsePairView{
		Response: v1.ResponseDetailsView{
			Status:      200,
			Body:        "hello_world",
			EncodedBody: false,
			Headers:     map[string][]string{"Content-Type": []string{"text/plain"}}},
		Request: v1.RequestDetailsView{
			Path:        StringToPointer("/"),
			Method:      StringToPointer("GET"),
			Destination: StringToPointer("/"),
			Scheme:      StringToPointer("scheme"),
			Query:       StringToPointer(""),
			Body:        StringToPointer(""),
			Headers:     map[string][]string{"Hoverfly": []string{"testing"}}}}

	hv.ImportRequestResponsePairViews([]interfaces.RequestResponsePair{originalPair})
	value, _ := cache.Get([]byte("9b114df98da7f7e2afdc975883dab4f2"))
	decodedPair, _ := models.NewRequestResponsePairFromBytes(value)
	Expect(*decodedPair).To(Equal(models.RequestResponsePair{
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

func TestImportImportRequestResponsePairs_CanImportAMultiplePairs(t *testing.T) {
	RegisterTestingT(t)

	cache := cache.NewInMemoryCache()
	cfg := Configuration{Webserver: false}
	requestMatcher := matching.RequestMatcher{RequestCache: cache, Webserver: &cfg.Webserver}
	hv := Hoverfly{RequestCache: cache, Cfg: &cfg, RequestMatcher: requestMatcher}

	RegisterTestingT(t)

	originalPair1 := v1.RequestResponsePairView{
		Response: v1.ResponseDetailsView{
			Status:      200,
			Body:        "hello_world",
			EncodedBody: false,
			Headers:     map[string][]string{"Hoverfly": []string{"testing"}},
		},
		Request: v1.RequestDetailsView{
			Path:        StringToPointer("/"),
			Method:      StringToPointer("GET"),
			Destination: StringToPointer("/"),
			Scheme:      StringToPointer("scheme"),
			Query:       StringToPointer(""),
			Body:        StringToPointer(""),
			Headers:     map[string][]string{"Hoverfly": []string{"testing"}}}}

	originalPair2 := originalPair1
	originalPair2.Request.Path = StringToPointer("/new/path")

	originalPair3 := originalPair1
	originalPair3.Request.Path = StringToPointer("/newer/path")

	hv.ImportRequestResponsePairViews([]interfaces.RequestResponsePair{originalPair1, originalPair2, originalPair3})

	pairBytes, err := cache.Get([]byte("9b114df98da7f7e2afdc975883dab4f2"))
	Expect(err).To(BeNil())
	decodedPair1, err := models.NewRequestResponsePairFromBytes(pairBytes)
	Expect(err).To(BeNil())
	Expect(*decodedPair1).To(Equal(models.NewRequestResponsePairFromRequestResponsePairView(originalPair1)))

	pairBytes, err = cache.Get([]byte("9c03e4af1f30542ff079a712bddad602"))
	Expect(err).To(BeNil())
	decodedPair2, err := models.NewRequestResponsePairFromBytes(pairBytes)
	Expect(err).To(BeNil())
	Expect(*decodedPair2).To(Equal(models.NewRequestResponsePairFromRequestResponsePairView(originalPair2)))

	pairBytes, err = cache.Get([]byte("fd099332afee48101edb7441b098cd4a"))
	Expect(err).To(BeNil())
	decodedPair3, err := models.NewRequestResponsePairFromBytes(pairBytes)
	Expect(err).To(BeNil())
	Expect(*decodedPair3).To(Equal(models.NewRequestResponsePairFromRequestResponsePairView(originalPair3)))
}

func TestImportImportRequestResponsePairs_CanImportARequestTemplateResponsePair(t *testing.T) {
	RegisterTestingT(t)

	cache := cache.NewInMemoryCache()
	cfg := Configuration{Webserver: false}
	requestMatcher := matching.RequestMatcher{RequestCache: cache, Webserver: &cfg.Webserver}
	hv := Hoverfly{RequestCache: cache, Cfg: &cfg, RequestMatcher: requestMatcher, Simulation: &models.Simulation{}}

	RegisterTestingT(t)

	requestTemplate := v1.RequestDetailsView{
		RequestType: StringToPointer("template"),
		Method:      StringToPointer("GET"),
	}

	responseView := v1.ResponseDetailsView{
		Status:      200,
		Body:        "hello_world",
		EncodedBody: false,
		Headers:     map[string][]string{"Hoverfly": []string{"testing"}},
	}

	templatePair := v1.RequestResponsePairView{
		Response: responseView,
		Request:  requestTemplate,
	}

	hv.ImportRequestResponsePairViews([]interfaces.RequestResponsePair{templatePair})

	Expect(len(hv.Simulation.Templates)).To(Equal(1))

	Expect(hv.Simulation.Templates[0].RequestTemplate.Method).To(Equal(StringToPointer("GET")))

	Expect(hv.Simulation.Templates[0].Response.Status).To(Equal(200))
	Expect(hv.Simulation.Templates[0].Response.Body).To(Equal("hello_world"))
	Expect(hv.Simulation.Templates[0].Response.Headers).To(Equal(map[string][]string{"Hoverfly": []string{"testing"}}))
}

// Helper function for base64 encoding
func base64String(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func TestImportImportRequestResponsePairs_CanImportASingleBase64EncodedPair(t *testing.T) {
	RegisterTestingT(t)

	cache := cache.NewInMemoryCache()
	cfg := Configuration{Webserver: false}
	requestMatcher := matching.RequestMatcher{RequestCache: cache, Webserver: &cfg.Webserver}
	hv := Hoverfly{RequestCache: cache, Cfg: &cfg, RequestMatcher: requestMatcher}

	RegisterTestingT(t)

	encodedPair := v1.RequestResponsePairView{
		Response: v1.ResponseDetailsView{
			Status:      200,
			Body:        base64String("hello_world"),
			EncodedBody: true,
			Headers:     map[string][]string{"Content-Encoding": []string{"gzip"}}},
		Request: v1.RequestDetailsView{
			Path:        StringToPointer("/"),
			Method:      StringToPointer("GET"),
			Destination: StringToPointer("/"),
			Scheme:      StringToPointer("scheme"),
			Query:       StringToPointer(""),
			Body:        StringToPointer(""),
			Headers:     map[string][]string{"Hoverfly": []string{"testing"}}}}

	hv.ImportRequestResponsePairViews([]interfaces.RequestResponsePair{encodedPair})

	value, err := cache.Get([]byte("9b114df98da7f7e2afdc975883dab4f2"))
	Expect(err).To(BeNil())

	decodedPair, err := models.NewRequestResponsePairFromBytes(value)
	Expect(err).To(BeNil())

	Expect(decodedPair).ToNot(Equal(models.RequestResponsePair{
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

package hoverfly

import (
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/SpectoLabs/hoverfly/core/cache"
	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
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

	err := dbClient.Import("examples/exports/readthedocs.json")
	Expect(err).To(BeNil())

	Expect(dbClient.Simulation.Templates).To(HaveLen(5))
}

func TestImportFromDiskBlankPath(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()

	err := dbClient.ImportFromDisk("")
	Expect(err).ToNot(BeNil())
}

func TestImportFromDiskWrongJson(t *testing.T) {
	RegisterTestingT(t)

	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()

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

	// importing payloads
	err = dbClient.Import(server.URL)
	Expect(err).To(BeNil())

	Expect(dbClient.Simulation.Templates).To(HaveLen(5))
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

	dbClient.HTTP = GetDefaultHoverflyHTTPClient(false, "")

	redirectServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Location", server.URL)
		w.WriteHeader(301)
	}))
	defer redirectServer.Close()

	// importing payloads
	err = dbClient.Import(redirectServer.URL)
	Expect(err).To(BeNil())

	Expect(dbClient.Simulation.Templates).To(HaveLen(5))
}

func TestImportFromURLHTTPFail(t *testing.T) {
	RegisterTestingT(t)

	// this tests simulates unreachable server
	server, dbClient := testTools(200, `this shouldn't matter anyway`)
	// closing it immediately
	server.Close()

	err := dbClient.ImportFromURL("somepath")
	Expect(err).ToNot(BeNil())
}

func TestImportFromURLMalformedJSON(t *testing.T) {
	RegisterTestingT(t)

	// testing behaviour when there is no json on the other end
	server, dbClient := testTools(200, `i am not json :(`)
	defer server.Close()

	// importing payloads
	err := dbClient.Import("http://thiswillbeintercepted.json")
	// we should get error
	Expect(err).ToNot(BeNil())
}

func TestImportRequestResponsePairs_CanImportASinglePair(t *testing.T) {
	RegisterTestingT(t)

	cache := cache.NewInMemoryCache()
	cfg := Configuration{Webserver: false}
	cacheMatcher := matching.CacheMatcher{RequestCache: cache, Webserver: &cfg.Webserver}
	hv := Hoverfly{Cfg: &cfg, CacheMatcher: cacheMatcher, Simulation: models.NewSimulation()}

	RegisterTestingT(t)

	originalPair := v2.RequestResponsePairView{
		Response: v2.ResponseDetailsView{
			Status:      200,
			Body:        "hello_world",
			EncodedBody: false,
			Headers:     map[string][]string{"Content-Type": []string{"text/plain"}}},
		Request: v2.RequestDetailsView{
			Path:        StringToPointer("/"),
			Method:      StringToPointer("GET"),
			Destination: StringToPointer("/"),
			Scheme:      StringToPointer("scheme"),
			Query:       StringToPointer(""),
			Body:        StringToPointer(""),
			Headers:     map[string][]string{"Hoverfly": []string{"testing"}}}}

	hv.ImportRequestResponsePairViews([]interfaces.RequestResponsePair{originalPair})

	Expect(hv.Simulation.Templates[0]).To(Equal(models.RequestTemplateResponsePair{
		Response: models.ResponseDetails{
			Status:  200,
			Body:    "hello_world",
			Headers: map[string][]string{"Content-Type": []string{"text/plain"}},
		},
		RequestTemplate: models.RequestTemplate{
			Path:        StringToPointer("/"),
			Method:      StringToPointer("GET"),
			Destination: StringToPointer("/"),
			Scheme:      StringToPointer("scheme"),
			Query:       StringToPointer(""),
			Body:        StringToPointer(""),
			Headers: map[string][]string{
				"Hoverfly": []string{"testing"},
			},
		},
	}))
}

func TestImportImportRequestResponsePairs_CanImportAMultiplePairs(t *testing.T) {
	RegisterTestingT(t)

	cache := cache.NewInMemoryCache()
	cfg := Configuration{Webserver: false}
	cacheMatcher := matching.CacheMatcher{RequestCache: cache, Webserver: &cfg.Webserver}
	hv := Hoverfly{Cfg: &cfg, CacheMatcher: cacheMatcher, Simulation: models.NewSimulation()}

	RegisterTestingT(t)

	originalPair1 := v2.RequestResponsePairView{
		Response: v2.ResponseDetailsView{
			Status:      200,
			Body:        "hello_world",
			EncodedBody: false,
			Headers:     map[string][]string{"Hoverfly": []string{"testing"}},
		},
		Request: v2.RequestDetailsView{
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

	Expect(hv.Simulation.Templates[0]).To(Equal(models.RequestTemplateResponsePair{
		Response: models.ResponseDetails{
			Status:  200,
			Body:    "hello_world",
			Headers: map[string][]string{"Hoverfly": []string{"testing"}},
		},
		RequestTemplate: models.RequestTemplate{
			Path:        StringToPointer("/"),
			Method:      StringToPointer("GET"),
			Destination: StringToPointer("/"),
			Scheme:      StringToPointer("scheme"),
			Query:       StringToPointer(""),
			Body:        StringToPointer(""),
			Headers:     map[string][]string{"Hoverfly": []string{"testing"}},
		},
	}))

	Expect(hv.Simulation.Templates[1]).To(Equal(models.RequestTemplateResponsePair{
		Response: models.ResponseDetails{
			Status:  200,
			Body:    "hello_world",
			Headers: map[string][]string{"Hoverfly": []string{"testing"}},
		},
		RequestTemplate: models.RequestTemplate{
			Path:        StringToPointer("/new/path"),
			Method:      StringToPointer("GET"),
			Destination: StringToPointer("/"),
			Scheme:      StringToPointer("scheme"),
			Query:       StringToPointer(""),
			Body:        StringToPointer(""),
			Headers:     map[string][]string{"Hoverfly": []string{"testing"}},
		},
	}))

	Expect(hv.Simulation.Templates[2]).To(Equal(models.RequestTemplateResponsePair{
		Response: models.ResponseDetails{
			Status:  200,
			Body:    "hello_world",
			Headers: map[string][]string{"Hoverfly": []string{"testing"}},
		},
		RequestTemplate: models.RequestTemplate{
			Path:        StringToPointer("/newer/path"),
			Method:      StringToPointer("GET"),
			Destination: StringToPointer("/"),
			Scheme:      StringToPointer("scheme"),
			Query:       StringToPointer(""),
			Body:        StringToPointer(""),
			Headers:     map[string][]string{"Hoverfly": []string{"testing"}},
		},
	}))
}

func TestImportImportRequestResponsePairs_CanImportARequestTemplateResponsePair(t *testing.T) {
	RegisterTestingT(t)

	cache := cache.NewInMemoryCache()
	cfg := Configuration{Webserver: false}
	cacheMatcher := matching.CacheMatcher{RequestCache: cache, Webserver: &cfg.Webserver}
	hv := Hoverfly{Cfg: &cfg, CacheMatcher: cacheMatcher, Simulation: models.NewSimulation()}

	RegisterTestingT(t)

	requestTemplate := v2.RequestDetailsView{
		Method: StringToPointer("GET"),
	}

	responseView := v2.ResponseDetailsView{
		Status:      200,
		Body:        "hello_world",
		EncodedBody: false,
		Headers:     map[string][]string{"Hoverfly": []string{"testing"}},
	}

	templatePair := v2.RequestResponsePairView{
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
	cacheMatcher := matching.CacheMatcher{RequestCache: cache, Webserver: &cfg.Webserver}
	hv := Hoverfly{Cfg: &cfg, CacheMatcher: cacheMatcher, Simulation: models.NewSimulation()}

	RegisterTestingT(t)

	encodedPair := v2.RequestResponsePairView{
		Response: v2.ResponseDetailsView{
			Status:      200,
			Body:        base64String("hello_world"),
			EncodedBody: true,
			Headers:     map[string][]string{"Content-Encoding": []string{"gzip"}}},
		Request: v2.RequestDetailsView{
			Path:        StringToPointer("/"),
			Method:      StringToPointer("GET"),
			Destination: StringToPointer("/"),
			Scheme:      StringToPointer("scheme"),
			Query:       StringToPointer(""),
			Body:        StringToPointer(""),
			Headers:     map[string][]string{"Hoverfly": []string{"testing"}}}}

	hv.ImportRequestResponsePairViews([]interfaces.RequestResponsePair{encodedPair})

	Expect(hv.Simulation.Templates[0]).ToNot(Equal(models.RequestResponsePair{
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

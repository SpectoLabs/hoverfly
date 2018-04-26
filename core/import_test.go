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
	"github.com/SpectoLabs/hoverfly/core/matching"
	"github.com/SpectoLabs/hoverfly/core/models"
	. "github.com/onsi/gomega"
)

const hoverfly_io_simulation_path = "../examples/simulations/hoverfly.io.json"

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

	ex, err := exists(hoverfly_io_simulation_path)
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

	err := dbClient.Import(hoverfly_io_simulation_path)
	Expect(err).To(BeNil())

	Expect(dbClient.Simulation.GetMatchingPairs()).To(HaveLen(2))
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
	pairFile, err := os.Open(hoverfly_io_simulation_path)
	Expect(err).To(BeNil())
	pairFileBytes, err := ioutil.ReadAll(pairFile)
	Expect(err).To(BeNil())

	// pretending this is the endpoint with given json
	server, dbClient := testTools(200, string(pairFileBytes))
	defer server.Close()

	// importing payloads
	err = dbClient.Import(server.URL)
	Expect(err).To(BeNil())

	Expect(dbClient.Simulation.GetMatchingPairs()).To(HaveLen(2))
}

func TestImportFromURLRedirect(t *testing.T) {
	RegisterTestingT(t)

	// reading file and preparing json payload
	pairFile, err := os.Open(hoverfly_io_simulation_path)
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

	Expect(dbClient.Simulation.GetMatchingPairs()).To(HaveLen(2))
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
	cacheMatcher := matching.CacheMatcher{RequestCache: cache, Webserver: cfg.Webserver}
	hv := Hoverfly{Cfg: &cfg, CacheMatcher: cacheMatcher, Simulation: models.NewSimulation()}

	RegisterTestingT(t)

	originalPair := v2.RequestMatcherResponsePairViewV5{
		Response: v2.ResponseDetailsViewV5{
			Status:      200,
			Body:        "hello_world",
			EncodedBody: false,
			Headers:     map[string][]string{"Content-Type": []string{"text/plain"}},
			Templated:   true,
		},
		RequestMatcher: v2.RequestMatcherViewV5{
			Path: []v2.MatcherViewV5{
				{
					Matcher: "exact",
					Value:   "/",
				},
			},
			Method: []v2.MatcherViewV5{
				{
					Matcher: "exact",
					Value:   "GET",
				},
			},
			Destination: []v2.MatcherViewV5{
				{
					Matcher: "exact",
					Value:   "/",
				},
			},
			Scheme: []v2.MatcherViewV5{
				{
					Matcher: "exact",
					Value:   "scheme",
				},
			},
			Query: []v2.MatcherViewV5{
				{
					Matcher: "exact",
					Value:   "",
				},
			},
			Body: []v2.MatcherViewV5{
				{
					Matcher: "exact",
					Value:   "",
				},
			},
			Headers: map[string][]string{"Hoverfly": []string{"testing"}}}}

	hv.ImportRequestResponsePairViews([]v2.RequestMatcherResponsePairViewV5{originalPair})

	Expect(hv.Simulation.GetMatchingPairs()[0]).To(Equal(models.RequestMatcherResponsePair{
		Response: models.ResponseDetails{
			Status:    200,
			Body:      "hello_world",
			Headers:   map[string][]string{"Content-Type": []string{"text/plain"}},
			Templated: true,
		},
		RequestMatcher: models.RequestMatcher{
			Path: []models.RequestFieldMatchers{
				{
					Matcher: "exact",
					Value:   "/",
				},
			},
			Method: []models.RequestFieldMatchers{
				{
					Matcher: "exact",
					Value:   "GET",
				},
			},
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: "exact",
					Value:   "/",
				},
			},
			Scheme: []models.RequestFieldMatchers{
				{
					Matcher: "exact",
					Value:   "scheme",
				},
			},
			Query: []models.RequestFieldMatchers{
				{
					Matcher: "exact",
					Value:   "",
				},
			},
			Body: []models.RequestFieldMatchers{
				{
					Matcher: "exact",
					Value:   "",
				},
			},
			Headers: map[string][]string{
				"Hoverfly": []string{"testing"},
			},
		},
	}))
}

func TestImportImportRequestResponsePairs_CanImportAMultiplePairsAndSetTemplateExplicitlyOrExplicitly(t *testing.T) {
	RegisterTestingT(t)

	cache := cache.NewInMemoryCache()
	cfg := Configuration{Webserver: false}
	cacheMatcher := matching.CacheMatcher{RequestCache: cache, Webserver: cfg.Webserver}
	hv := Hoverfly{Cfg: &cfg, CacheMatcher: cacheMatcher, Simulation: models.NewSimulation()}

	RegisterTestingT(t)

	originalPair1 := v2.RequestMatcherResponsePairViewV5{
		Response: v2.ResponseDetailsViewV5{
			Status:      200,
			Body:        "hello_world",
			EncodedBody: false,
			Headers:     map[string][]string{"Hoverfly": []string{"testing"}},
		},
		RequestMatcher: v2.RequestMatcherViewV5{
			Path: []v2.MatcherViewV5{
				{
					Matcher: "exact",
					Value:   "/",
				},
			},
			Method: []v2.MatcherViewV5{
				{
					Matcher: "exact",
					Value:   "GET",
				},
			},
			Destination: []v2.MatcherViewV5{
				{
					Matcher: "exact",
					Value:   "/",
				},
			},
			Scheme: []v2.MatcherViewV5{
				{
					Matcher: "exact",
					Value:   "scheme",
				},
			},
			Query: []v2.MatcherViewV5{
				{
					Matcher: "exact",
					Value:   "",
				},
			},
			Body: []v2.MatcherViewV5{
				{
					Matcher: "exact",
					Value:   "",
				},
			},
			Headers: map[string][]string{"Hoverfly": []string{"testing"}}}}

	originalPair2 := originalPair1
	originalPair2.Response.Templated = false
	originalPair2.RequestMatcher.Path = []v2.MatcherViewV5{
		{
			Matcher: "exact",
			Value:   "/new/path",
		},
	}

	originalPair3 := originalPair1
	originalPair3.RequestMatcher.Path = []v2.MatcherViewV5{
		{
			Matcher: "exact",
			Value:   "/newer/path",
		},
	}
	originalPair3.Response.Templated = true

	hv.ImportRequestResponsePairViews([]v2.RequestMatcherResponsePairViewV5{originalPair1, originalPair2, originalPair3})

	Expect(hv.Simulation.GetMatchingPairs()).To(HaveLen(3))
	Expect(hv.Simulation.GetMatchingPairs()[0]).To(Equal(models.RequestMatcherResponsePair{
		Response: models.ResponseDetails{
			Status:    200,
			Body:      "hello_world",
			Headers:   map[string][]string{"Hoverfly": []string{"testing"}},
			Templated: false,
		},
		RequestMatcher: models.RequestMatcher{
			Path: []models.RequestFieldMatchers{
				{
					Matcher: "exact",
					Value:   "/",
				},
			},
			Method: []models.RequestFieldMatchers{
				{
					Matcher: "exact",
					Value:   "GET",
				},
			},
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: "exact",
					Value:   "/",
				},
			},
			Scheme: []models.RequestFieldMatchers{
				{
					Matcher: "exact",
					Value:   "scheme",
				},
			},
			Query: []models.RequestFieldMatchers{
				{
					Matcher: "exact",
					Value:   "",
				},
			},
			Body: []models.RequestFieldMatchers{
				{
					Matcher: "exact",
					Value:   "",
				},
			},
			Headers: map[string][]string{"Hoverfly": []string{"testing"}},
		},
	}))

	Expect(hv.Simulation.GetMatchingPairs()[1]).To(Equal(models.RequestMatcherResponsePair{
		Response: models.ResponseDetails{
			Status:    200,
			Body:      "hello_world",
			Headers:   map[string][]string{"Hoverfly": []string{"testing"}},
			Templated: false,
		},
		RequestMatcher: models.RequestMatcher{
			Path: []models.RequestFieldMatchers{
				{
					Matcher: "exact",
					Value:   "/new/path",
				},
			},
			Method: []models.RequestFieldMatchers{
				{
					Matcher: "exact",
					Value:   "GET",
				},
			},
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: "exact",
					Value:   "/",
				},
			},
			Scheme: []models.RequestFieldMatchers{
				{
					Matcher: "exact",
					Value:   "scheme",
				},
			},
			Query: []models.RequestFieldMatchers{
				{
					Matcher: "exact",
					Value:   "",
				},
			},
			Body: []models.RequestFieldMatchers{
				{
					Matcher: "exact",
					Value:   "",
				},
			},
			Headers: map[string][]string{"Hoverfly": []string{"testing"}},
		},
	}))

	Expect(hv.Simulation.GetMatchingPairs()[2]).To(Equal(models.RequestMatcherResponsePair{
		Response: models.ResponseDetails{
			Status:    200,
			Body:      "hello_world",
			Headers:   map[string][]string{"Hoverfly": []string{"testing"}},
			Templated: true,
		},
		RequestMatcher: models.RequestMatcher{
			Path: []models.RequestFieldMatchers{
				{
					Matcher: "exact",
					Value:   "/newer/path",
				},
			},
			Method: []models.RequestFieldMatchers{
				{
					Matcher: "exact",
					Value:   "GET",
				},
			},
			Destination: []models.RequestFieldMatchers{
				{
					Matcher: "exact",
					Value:   "/",
				},
			},
			Scheme: []models.RequestFieldMatchers{
				{
					Matcher: "exact",
					Value:   "scheme",
				},
			},
			Query: []models.RequestFieldMatchers{
				{
					Matcher: "exact",
					Value:   "",
				},
			},
			Body: []models.RequestFieldMatchers{
				{
					Matcher: "exact",
					Value:   "",
				},
			},
			Headers: map[string][]string{"Hoverfly": []string{"testing"}},
		},
	}))
}

func TestImportImportRequestResponsePairs_CanImportARequesResponsePairView(t *testing.T) {
	RegisterTestingT(t)

	cache := cache.NewInMemoryCache()
	cfg := Configuration{Webserver: false}
	cacheMatcher := matching.CacheMatcher{RequestCache: cache, Webserver: cfg.Webserver}
	hv := Hoverfly{Cfg: &cfg, CacheMatcher: cacheMatcher, Simulation: models.NewSimulation()}

	RegisterTestingT(t)

	request := v2.RequestMatcherViewV5{
		Method: []v2.MatcherViewV5{
			{
				Matcher: "exact",
				Value:   "GET",
			},
		},
	}

	responseView := v2.ResponseDetailsViewV5{
		Status:      200,
		Body:        "hello_world",
		EncodedBody: false,
		Headers:     map[string][]string{"Hoverfly": []string{"testing"}},
	}

	requestResponsePair := v2.RequestMatcherResponsePairViewV5{
		Response:       responseView,
		RequestMatcher: request,
	}

	hv.ImportRequestResponsePairViews([]v2.RequestMatcherResponsePairViewV5{requestResponsePair})

	Expect(len(hv.Simulation.GetMatchingPairs())).To(Equal(1))

	Expect(hv.Simulation.GetMatchingPairs()[0].RequestMatcher.Method).To(HaveLen(1))
	Expect(hv.Simulation.GetMatchingPairs()[0].RequestMatcher.Method[0].Matcher).To(Equal("exact"))
	Expect(hv.Simulation.GetMatchingPairs()[0].RequestMatcher.Method[0].Value).To(Equal("GET"))

	Expect(hv.Simulation.GetMatchingPairs()[0].Response.Status).To(Equal(200))
	Expect(hv.Simulation.GetMatchingPairs()[0].Response.Body).To(Equal("hello_world"))
	Expect(hv.Simulation.GetMatchingPairs()[0].Response.Headers).To(Equal(map[string][]string{"Hoverfly": []string{"testing"}}))
}

// Helper function for base64 encoding
func base64String(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func TestImportImportRequestResponsePairs_CanImportASingleBase64EncodedPair(t *testing.T) {
	RegisterTestingT(t)

	cache := cache.NewInMemoryCache()
	cfg := Configuration{Webserver: false}
	cacheMatcher := matching.CacheMatcher{RequestCache: cache, Webserver: cfg.Webserver}
	hv := Hoverfly{Cfg: &cfg, CacheMatcher: cacheMatcher, Simulation: models.NewSimulation()}

	RegisterTestingT(t)

	encodedPair := v2.RequestMatcherResponsePairViewV5{
		Response: v2.ResponseDetailsViewV5{
			Status:      200,
			Body:        base64String("hello_world"),
			EncodedBody: true,
			Headers:     map[string][]string{"Content-Encoding": []string{"gzip"}}},
		RequestMatcher: v2.RequestMatcherViewV5{
			Path: []v2.MatcherViewV5{
				{
					Matcher: "exact",
					Value:   "/",
				},
			},
			Method: []v2.MatcherViewV5{
				{
					Matcher: "exact",
					Value:   "GET",
				},
			},
			Destination: []v2.MatcherViewV5{
				{
					Matcher: "exact",
					Value:   "/",
				},
			},
			Scheme: []v2.MatcherViewV5{
				{
					Matcher: "exact",
					Value:   "scheme",
				},
			},
			Query: []v2.MatcherViewV5{
				{
					Matcher: "exact",
					Value:   "",
				},
			},
			Body: []v2.MatcherViewV5{
				{
					Matcher: "exact",
					Value:   "",
				},
			},
			Headers: map[string][]string{
				"Hoverfly": []string{
					"testing",
				},
			},
		},
	}

	hv.ImportRequestResponsePairViews([]v2.RequestMatcherResponsePairViewV5{encodedPair})

	Expect(hv.Simulation.GetMatchingPairs()[0]).ToNot(Equal(models.RequestResponsePair{
		Response: models.ResponseDetails{
			Status:  200,
			Body:    "hello_world",
			Headers: map[string][]string{"Content-Encoding": []string{"gzip"}}},
		Request: models.RequestDetails{
			Path:        "/",
			Method:      "GET",
			Destination: "/",
			Scheme:      "scheme",
			Query:       map[string][]string{},
			Body:        "",
			Headers:     map[string][]string{"Hoverfly": []string{"testing"}}}}))
}

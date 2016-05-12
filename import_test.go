package hoverfly

import (
	"github.com/SpectoLabs/hoverfly/testutil"
	"io/ioutil"
	"os"
	"testing"
	"github.com/SpectoLabs/hoverfly/cache"
	. "github.com/onsi/gomega"
	"encoding/base64"
)

func TestIsURLHTTP(t *testing.T) {
	url := "http://somehost.com"

	b := isURL(url)
	testutil.Expect(t, b, true)
}

func TestIsURLEmpty(t *testing.T) {
	b := isURL("")
	testutil.Expect(t, b, false)
}

func TestIsURLHTTPS(t *testing.T) {
	url := "https://somehost.com"

	b := isURL(url)
	testutil.Expect(t, b, true)
}

func TestIsURLWrong(t *testing.T) {
	url := "somehost.com"

	b := isURL(url)
	testutil.Expect(t, b, false)
}

func TestIsURLWrongTLD(t *testing.T) {
	url := "http://somehost."

	b := isURL(url)
	testutil.Expect(t, b, false)
}

func TestFileExists(t *testing.T) {
	fp := "examples/exports/readthedocs.json"

	ex, err := exists(fp)
	testutil.Expect(t, ex, true)
	testutil.Expect(t, err, nil)
}

func TestFileDoesNotExist(t *testing.T) {
	fp := "shouldnotbehere.yaml"

	ex, err := exists(fp)
	testutil.Expect(t, ex, false)
	testutil.Expect(t, err, nil)
}

func TestImportFromDisk(t *testing.T) {
	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()

	err := dbClient.Import("examples/exports/readthedocs.json")
	testutil.Expect(t, err, nil)

	recordsCount, err := dbClient.RequestCache.RecordsCount()
	testutil.Expect(t, err, nil)
	testutil.Expect(t, recordsCount, 5)
}

func TestImportFromDiskBlankPath(t *testing.T) {
	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()

	err := dbClient.ImportFromDisk("")
	testutil.Refute(t, err, nil)
}

func TestImportFromDiskWrongJson(t *testing.T) {
	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()

	err := dbClient.ImportFromDisk("examples/exports/README.md")
	testutil.Refute(t, err, nil)
}

func TestImportFromURL(t *testing.T) {
	// reading file and preparing json payload
	payloadsFile, err := os.Open("examples/exports/readthedocs.json")
	testutil.Expect(t, err, nil)
	bts, err := ioutil.ReadAll(payloadsFile)
	testutil.Expect(t, err, nil)

	// pretending this is the endpoint with given json
	server, dbClient := testTools(200, string(bts))
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()

	// importing payloads
	err = dbClient.Import("http://thiswillbeintercepted.json")
	testutil.Expect(t, err, nil)

	recordsCount, err := dbClient.RequestCache.RecordsCount()
	testutil.Expect(t, err, nil)
	testutil.Expect(t, recordsCount, 5)
}

func TestImportFromURLHTTPFail(t *testing.T) {
	// this tests simulates unreachable server
	server, dbClient := testTools(200, `this shouldn't matter anyway`)
	// closing it immediately
	server.Close()
	defer dbClient.RequestCache.DeleteData()

	err := dbClient.ImportFromURL("somepath")
	testutil.Refute(t, err, nil)
}

func TestImportFromURLMalformedJSON(t *testing.T) {
	// testing behaviour when there is no json on the other end
	server, dbClient := testTools(200, `i am not json :(`)
	defer server.Close()
	defer dbClient.RequestCache.DeleteData()

	// importing payloads
	err := dbClient.Import("http://thiswillbeintercepted.json")
	// we should get error
	testutil.Refute(t, err, nil)
}

func TestImportPayloads_CanImportASinglePayload(t *testing.T) {
	cache := cache.NewInMemoryCache()
	hv := Hoverfly{RequestCache: cache, MIN: GetNewMinifiers()}

	RegisterTestingT(t)

	originalPayload := PayloadView{
		Response: SerializableResponseDetails{
			Status: 200,
			Body: "hello_world",
			EncodedBody: false,
			Headers: map[string][]string{"Hoverfly": []string {"testing"}}},
		Request:  RequestDetails{
			Path: "/",
			Method: "GET",
			Destination: "/",
			Scheme: "scheme",
			Query: "", Body: "",
			RemoteAddr: "localhost",
			Headers: map[string][]string{"Hoverfly": []string {"testing"}}},
		ID: "9b114df98da7f7e2afdc975883dab4f2"}

	hv.ImportPayloads([]PayloadView{originalPayload})

	value, err := cache.Get([]byte(originalPayload.ID))
	Expect(err).To(BeNil())
	decodedPayload, err := decodePayload(value)
	Expect(err).To(BeNil())
	Expect(*decodedPayload).To(Equal(originalPayload.ConvertToPayload()))
}

func TestImportPayloads_CanImportAMultiplePayload(t *testing.T) {
	cache := cache.NewInMemoryCache()
	hv := Hoverfly{RequestCache: cache, MIN: GetNewMinifiers()}

	RegisterTestingT(t)

	originalPayload1 := PayloadView{
		Response: SerializableResponseDetails{
			Status: 200,
			Body: "hello_world",
			EncodedBody: false,
			Headers: map[string][]string{"Hoverfly": []string {"testing"}}},
		Request:  RequestDetails{
			Path: "/",
			Method: "GET",
			Destination: "/",
			Scheme: "scheme",
			Query: "", Body: "",
			RemoteAddr: "localhost",
			Headers: map[string][]string{"Hoverfly": []string {"testing"}}},
		ID: "9b114df98da7f7e2afdc975883dab4f2"}

	originalPayload2 := originalPayload1

	originalPayload2.ID = "9c03e4af1f30542ff079a712bddad602"
	originalPayload2.Request.Path = "/new/path"

	originalPayload3 := originalPayload1

	originalPayload3.ID = "fd099332afee48101edb7441b098cd4a"
	originalPayload3.Request.Path = "/newer/path"

	hv.ImportPayloads([]PayloadView{originalPayload1, originalPayload2, originalPayload3})

	value, err := cache.Get([]byte(originalPayload1.ID))
	Expect(err).To(BeNil())
	decodedPayload1, err := decodePayload(value)
	Expect(err).To(BeNil())
	Expect(*decodedPayload1).To(Equal(originalPayload1.ConvertToPayload()))

	value, err = cache.Get([]byte(originalPayload2.ID))
	Expect(err).To(BeNil())
	decodedPayload2, err := decodePayload(value)
	Expect(err).To(BeNil())
	Expect(*decodedPayload2).To(Equal(originalPayload2.ConvertToPayload()))

	value, err = cache.Get([]byte(originalPayload3.ID))
	Expect(err).To(BeNil())
	decodedPayload3, err := decodePayload(value)
	Expect(err).To(BeNil())
	Expect(*decodedPayload3).To(Equal(originalPayload3.ConvertToPayload()))
}

// Helper function for base64 encoding
func base64String(s string) (string) {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func TestImportPayloads_CanImportASingleBase64EncodedPayload(t *testing.T) {
	cache := cache.NewInMemoryCache()
	hv := Hoverfly{RequestCache: cache, MIN: GetNewMinifiers()}

	RegisterTestingT(t)

	encodedPayload := PayloadView{
		Response: SerializableResponseDetails{
			Status: 200,
			Body: base64String("hello_world"),
			EncodedBody: true,
			Headers: map[string][]string{"Content-Encoding": []string {"gzip"}}},
		Request:  RequestDetails{
			Path: "/",
			Method: "GET",
			Destination: "/",
			Scheme: "scheme",
			Query: "", Body: "",
			RemoteAddr: "localhost",
			Headers: map[string][]string{"Hoverfly": []string {"testing"}}},
		ID: "9b114df98da7f7e2afdc975883dab4f2"}

	originalPayload := encodedPayload
	originalPayload.Response.Body = "hello_world"

	hv.ImportPayloads([]PayloadView{encodedPayload})

	value, err := cache.Get([]byte(encodedPayload.ID))
	Expect(err).To(BeNil())
	decodedPayload, err := decodePayload(value)
	Expect(err).To(BeNil())
	Expect(*decodedPayload).ToNot(Equal(encodedPayload))
	Expect(*decodedPayload).To(Equal(originalPayload.ConvertToPayload()))
}

package hoverfly

import (
	"github.com/SpectoLabs/hoverfly/testutil"
	"io/ioutil"
	"os"
	"testing"
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

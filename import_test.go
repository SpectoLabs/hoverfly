package hoverfly

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestIsURLHTTP(t *testing.T) {
	url := "http://somehost.com"

	b := isURL(url)
	expect(t, b, true)
}

func TestIsURLEmpty(t *testing.T) {
	b := isURL("")
	expect(t, b, false)
}

func TestIsURLHTTPS(t *testing.T) {
	url := "https://somehost.com"

	b := isURL(url)
	expect(t, b, true)
}

func TestIsURLWrong(t *testing.T) {
	url := "somehost.com"

	b := isURL(url)
	expect(t, b, false)
}

func TestIsURLWrongTLD(t *testing.T) {
	url := "http://somehost."

	b := isURL(url)
	expect(t, b, false)
}

func TestFileExists(t *testing.T) {
	fp := "examples/exports/readthedocs.json"

	ex, err := exists(fp)
	expect(t, ex, true)
	expect(t, err, nil)
}

func TestFileDoesNotExist(t *testing.T) {
	fp := "shouldnotbehere.yaml"

	ex, err := exists(fp)
	expect(t, ex, false)
	expect(t, err, nil)
}

func TestImportFromDisk(t *testing.T) {
	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.Cache.DeleteData()

	err := dbClient.Import("examples/exports/readthedocs.json")
	expect(t, err, nil)

	recordsCount, err := dbClient.Cache.RecordsCount()
	expect(t, err, nil)
	expect(t, recordsCount, 5)
}

func TestImportFromDiskBlankPath(t *testing.T) {
	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.Cache.DeleteData()

	err := dbClient.ImportFromDisk("")
	refute(t, err, nil)
}

func TestImportFromDiskWrongJson(t *testing.T) {
	server, dbClient := testTools(201, `{'message': 'here'}`)
	defer server.Close()
	defer dbClient.Cache.DeleteData()

	err := dbClient.ImportFromDisk("examples/exports/README.md")
	refute(t, err, nil)
}

func TestImportFromURL(t *testing.T) {
	// reading file and preparing json payload
	payloadsFile, err := os.Open("examples/exports/readthedocs.json")
	expect(t, err, nil)
	bts, err := ioutil.ReadAll(payloadsFile)
	expect(t, err, nil)

	// pretending this is the endpoint with given json
	server, dbClient := testTools(200, string(bts))
	defer server.Close()
	defer dbClient.Cache.DeleteData()

	// importing payloads
	err = dbClient.Import("http://thiswillbeintercepted.json")
	expect(t, err, nil)

	recordsCount, err := dbClient.Cache.RecordsCount()
	expect(t, err, nil)
	expect(t, recordsCount, 5)
}

func TestImportFromURLHTTPFail(t *testing.T) {
	// this tests simulates unreachable server
	server, dbClient := testTools(200, `this shouldn't matter anyway`)
	// closing it immediately
	server.Close()
	defer dbClient.Cache.DeleteData()

	err := dbClient.ImportFromURL("somepath")
	refute(t, err, nil)
}

func TestImportFromURLMalformedJSON(t *testing.T) {
	// testing behaviour when there is no json on the other end
	server, dbClient := testTools(200, `i am not json :(`)
	defer server.Close()
	defer dbClient.Cache.DeleteData()

	// importing payloads
	err := dbClient.Import("http://thiswillbeintercepted.json")
	// we should get error
	refute(t, err, nil)
}

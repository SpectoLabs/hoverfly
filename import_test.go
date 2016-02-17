package hoverfly

import (
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

func TestImportFromFile(t *testing.T) {
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

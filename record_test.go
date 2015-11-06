package main

import (
	"net/http"
   "testing"
)

func TestRecord(t *testing.T){
	server, client := testTools(200, `{'message': 'here'}`)
	defer server.Close()

	req, err := http.NewRequest("GET", "http://example.com", nil)
	expect(t, err, nil)

	response, err := client.recordRequest(req)

	expect(t, response.Header.Get("Gen-proxy"), "Was-Here")
}
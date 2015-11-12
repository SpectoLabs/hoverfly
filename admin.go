package main

import (
	"bytes"
	"github.com/elazarl/goproxy"
	"io/ioutil"
	"net/http"
)

// AllRecordsHandler returns JSON content type http response
func (d *DBClient) AllRecordsHandler(req *http.Request) *http.Response {

	records, err := d.getAllRecordsRaw()

	if err == nil {
		newResponse := &http.Response{}
		newResponse.Request = req

		newResponse.Header.Set("Content-Type", "application/json")
		// adding body
		buf := bytes.NewBufferString(records)
		newResponse.ContentLength = int64(buf.Len())
		newResponse.Body = ioutil.NopCloser(buf)

		newResponse.StatusCode = 200

		return newResponse

	} else {

		return goproxy.NewResponse(req,
			goproxy.ContentTypeText, http.StatusInternalServerError,
			"Failed to retrieve records from cache!")
	}

}

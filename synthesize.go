package main

import (
	"bufio"
	"fmt"
	"net/http"
)

// synthesizeResponse calls middleware to populate response data, nothing gets pass proxy
func synthesizeResponse(req *http.Request) *http.Response {

	b := bufio.NewScanner(req.Body)

	bodyStr := b.Text()

	request := requestDetails{
		Path:        req.URL.Path,
		Method:      req.Method,
		Destination: req.URL.Path,
		Query:       req.URL.RawQuery,
		Body:        bodyStr,
		RemoteAddr:  req.RemoteAddr,
		Headers:     req.Header,
	}
	payload := Payload{Request: request}

	c := NewConstructor(req, payload)

	err := c.ApplyMiddleware(AppConfig.middleware)
	if err != nil {
		var errorPayload Payload
		errorPayload.Response.Status = 503
		errorPayload.Response.Body = fmt.Sprintf("Middleware error: %s", err.Error())
		c.payload = errorPayload
	}

	response := c.reconstructResponse()
	return response

}

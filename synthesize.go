package main

import (
	"bufio"
	"fmt"
	"net/http"
)

// synthesizeResponse calls middleware to populate response data, nothing gets pass proxy
func synthesizeResponse(req *http.Request, middleware string) *http.Response {

	b := bufio.NewScanner(req.Body)

	bodyStr := b.Text()

	request := requestDetails{
		Path:        req.URL.Path,
		Method:      req.Method,
		Destination: req.Host,
		Query:       req.URL.RawQuery,
		Body:        bodyStr,
		RemoteAddr:  req.RemoteAddr,
		Headers:     req.Header,
	}
	payload := Payload{Request: request}

	c := NewConstructor(req, payload)

	if middleware != "" {
		err := c.ApplyMiddleware(middleware)
		if err != nil {
			var errorPayload Payload
			errorPayload.Response.Status = 503
			errorPayload.Response.Body = fmt.Sprintf("Middleware error: %s", err.Error())
			c.payload = errorPayload
		}
	} else {
		c.payload.Response.Body = "Precondition failed: middleware not provided."
		c.payload.Response.Status = 428
	}

	response := c.reconstructResponse()
	return response

}

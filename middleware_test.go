package main

import (
	"testing"
)

func TestChangeBodyMiddleware(t *testing.T) {
	command := "./examples/middleware/modify_response/modify_response.py"

	resp := response{Status: 201, Body: "original body"}
	req := requestDetails{Path: "/", Method: "GET", Destination: "hostname-x", Query: ""}

	payload := Payload{Response: resp, Request: req}

	newPayload, err := ExecuteMiddleware(command, payload)

	expect(t, err, nil)
	expect(t, newPayload.Response.Body, "body was replaced by middleware")
}

func TestMalformedPayloadMiddleware(t *testing.T) {
	command := "./examples/middleware/ruby_echo/echo.rb"

	resp := response{Status: 201, Body: "original body"}
	req := requestDetails{Path: "/", Method: "GET", Destination: "hostname-x", Query: ""}

	payload := Payload{Response: resp, Request: req}

	newPayload, err := ExecuteMiddleware(command, payload)

	expect(t, err, nil)
	expect(t, newPayload.Response.Body, "original body")
}


package hoverfly

import (
	"fmt"
	"testing"
)

func TestChangeBodyMiddleware(t *testing.T) {
	command := "./examples/middleware/modify_response/modify_response.py"

	resp := ResponseDetails{Status: 201, Body: "original body"}
	req := RequestDetails{Path: "/", Method: "GET", Destination: "hostname-x", Query: ""}

	payload := Payload{Response: resp, Request: req}

	newPayload, err := ExecuteMiddleware(command, payload)

	expect(t, err, nil)
	expect(t, newPayload.Response.Body, "body was replaced by middleware\n")
}

func TestMalformedPayloadMiddleware(t *testing.T) {
	command := "./examples/middleware/ruby_echo/echo.rb"

	resp := ResponseDetails{Status: 201, Body: "original body"}
	req := RequestDetails{Path: "/", Method: "GET", Destination: "hostname-x", Query: ""}

	payload := Payload{Response: resp, Request: req}

	newPayload, err := ExecuteMiddleware(command, payload)

	expect(t, err, nil)
	expect(t, newPayload.Response.Body, "original body")
}

func TestMakeCustom404(t *testing.T) {
	command := "go run ./examples/middleware/go_example/change_to_custom_404.go"

	resp := ResponseDetails{Status: 201, Body: "original body"}
	req := RequestDetails{Path: "/", Method: "GET", Destination: "hostname-x", Query: ""}

	payload := Payload{Response: resp, Request: req}

	newPayload, err := ExecuteMiddleware(command, payload)

	expect(t, err, nil)
	expect(t, newPayload.Response.Body, "Custom body here")
	expect(t, newPayload.Response.Status, 404)
	expect(t, newPayload.Response.Headers["middleware"][0], "changed response")
}

func TestReflectBody(t *testing.T) {
	command := "./examples/middleware/reflect_body/reflect_body.py"

	req := RequestDetails{Path: "/", Method: "GET", Destination: "hostname-x", Query: "", Body: "request_body_here"}

	payload := Payload{Request: req}

	newPayload, err := ExecuteMiddleware(command, payload)

	expect(t, err, nil)
	expect(t, newPayload.Response.Body, req.Body)
	expect(t, newPayload.Request.Method, req.Method)
	expect(t, newPayload.Request.Destination, req.Destination)
}

func TestPipe(t *testing.T) {
	middlewares := "./examples/middleware/modify_response/modify_response.py | ./examples/middleware/modify_status_code/modify_status_code.py "

	req := RequestDetails{Path: "/", Method: "GET", Destination: "hostname-x", Query: "", Body: "request_body_here"}

	payload := Payload{Request: req}

	newPayload, err := ExecuteMiddleware(middlewares, payload)

	expect(t, err, nil)
	expect(t, newPayload.Response.Body, "body was replaced by middleware\n")
	expect(t, newPayload.Response.Status, 301)
	expect(t, newPayload.Request.Method, req.Method)
}

func TestPipeOrder(t *testing.T) {
	middlewares := "./examples/middleware/modify_status_code/modify_status_code.py | ./examples/middleware/modify_response/modify_response.py "

	req := RequestDetails{Path: "/", Method: "GET", Destination: "hostname-x", Query: "", Body: "request_body_here"}

	payload := Payload{Request: req}

	newPayload, err := ExecuteMiddleware(middlewares, payload)

	expect(t, err, nil)
	expect(t, newPayload.Response.Body, "body was replaced by middleware\n")
	expect(t, newPayload.Response.Status, 201)
	expect(t, newPayload.Request.Method, req.Method)
}

func TestPipeMixed(t *testing.T) {

	m1 := "./examples/middleware/modify_status_code/modify_status_code.py"
	m2 := "./examples/middleware/modify_response/modify_response.py"
	m3 := "go run ./examples/middleware/go_example/change_to_custom_404.go"

	middlewares := fmt.Sprintf("%s | %s | %s", m1, m2, m3)

	req := RequestDetails{Path: "/", Method: "GET", Destination: "hostname-x", Query: "", Body: "request_body_here"}

	payload := Payload{Request: req}

	newPayload, err := ExecuteMiddleware(middlewares, payload)

	expect(t, err, nil)
	expect(t, newPayload.Response.Body, "Custom body here")
	expect(t, newPayload.Response.Status, 404)
	expect(t, newPayload.Request.Method, req.Method)
}

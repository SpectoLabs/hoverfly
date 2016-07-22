package hoverfly

import (
	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/SpectoLabs/hoverfly/core/testutil"
	"testing"
)

func TestChangeBodyMiddleware(t *testing.T) {
	command := "./examples/middleware/modify_response/modify_response.py"

	resp := models.ResponseDetails{Status: 201, Body: "original body"}
	req := models.RequestDetails{Path: "/", Method: "GET", Destination: "hostname-x", Query: ""}

	payload := models.Payload{Response: resp, Request: req}

	newPayload, err := ExecuteMiddlewareLocally(command, payload)

	testutil.Expect(t, err, nil)
	testutil.Expect(t, newPayload.Response.Body, "body was replaced by middleware\n")
}

func TestMalformedPayloadMiddleware(t *testing.T) {
	command := "./examples/middleware/ruby_echo/echo.rb"

	resp := models.ResponseDetails{Status: 201, Body: "original body"}
	req := models.RequestDetails{Path: "/", Method: "GET", Destination: "hostname-x", Query: ""}

	payload := models.Payload{Response: resp, Request: req}

	newPayload, err := ExecuteMiddlewareLocally(command, payload)

	testutil.Expect(t, err, nil)
	testutil.Expect(t, newPayload.Response.Body, "original body")
}

func TestMakeCustom404(t *testing.T) {
	command := "go run ./examples/middleware/go_example/change_to_custom_404.go"

	resp := models.ResponseDetails{Status: 201, Body: "original body"}
	req := models.RequestDetails{Path: "/", Method: "GET", Destination: "hostname-x", Query: ""}

	payload := models.Payload{Response: resp, Request: req}

	newPayload, err := ExecuteMiddlewareLocally(command, payload)

	testutil.Expect(t, err, nil)
	testutil.Expect(t, newPayload.Response.Body, "Custom body here")
	testutil.Expect(t, newPayload.Response.Status, 404)
	testutil.Expect(t, newPayload.Response.Headers["middleware"][0], "changed response")
}

func TestReflectBody(t *testing.T) {
	command := "./examples/middleware/reflect_body/reflect_body.py"

	req := models.RequestDetails{Path: "/", Method: "GET", Destination: "hostname-x", Query: "", Body: "request_body_here"}

	payload := models.Payload{Request: req}

	newPayload, err := ExecuteMiddlewareLocally(command, payload)

	testutil.Expect(t, err, nil)
	testutil.Expect(t, newPayload.Response.Body, req.Body)
	testutil.Expect(t, newPayload.Request.Method, req.Method)
	testutil.Expect(t, newPayload.Request.Destination, req.Destination)
}

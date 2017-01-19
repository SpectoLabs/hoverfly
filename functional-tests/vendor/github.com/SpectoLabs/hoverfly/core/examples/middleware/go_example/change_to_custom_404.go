package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
)

// requestDetails stores information about request, it's used for creating unique hash and also as a payload structure
type requestDetails struct {
	Path        string `json:"path"`
	Method      string `json:"method"`
	Destination string `json:"destination"`
	Query       string `json:"query"`
}

// res structure hold response body from external service, body is not decoded and is supposed
// to be bytes, however headers should provide all required information for later decoding
// by the client.
type responseDetails struct {
	Status  int                 `json:"status"`
	Body    string              `json:"body"`
	Headers map[string][]string `json:"headers"`
}

// Payload structure holds request and response structure
type Payload struct {
	Response responseDetails `json:"response"`
	Request  requestDetails  `json:"request"`
	ID       string          `json:"id"`
}

func main() {
	// logging to stderr
	l := log.New(os.Stderr, "", 0)

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {

		var payload Payload

		err := json.Unmarshal(s.Bytes(), &payload)
		if err != nil {
			l.Println("Failed to unmarshal payload from hoverfly")
		}

		newHeaders := make(map[string][]string)

		vals := []string{"changed response"}

		newHeaders["middleware"] = vals

		payload.Response.Status = 404
		payload.Response.Body = "Custom body here"
		payload.Response.Headers = newHeaders

		bts, err := json.Marshal(payload)
		if err != nil {
			l.Println("Failed to marshal new payload")
		}
		os.Stdout.Write(bts)

	}
}

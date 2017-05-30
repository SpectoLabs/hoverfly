package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
)

func main() {
	l := log.New(os.Stderr, "", 0)

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {

		var payload v2.RequestResponsePairViewV1

		err := json.Unmarshal(s.Bytes(), &payload)
		if err != nil {
			l.Println("Failed to unmarshal payload from hoverfly")
		}

		payload.Response.Body = "body was replaced by middleware\n"

		bts, err := json.Marshal(payload)
		if err != nil {
			l.Println("Failed to marshal new payload")
		}

		os.Stdout.Write(bts)
	}
}

package main

import (
	"bytes"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

type Constructor struct {
	request *http.Request
	payload Payload
}

func NewConstructor(req *http.Request, payload Payload) Constructor {
	c := Constructor{request: req, payload: payload}
	return c
}


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

func (c *Constructor) ApplyMiddleware(middleware string) error {

	newPayload, err := ExecuteMiddleware(middleware, c.payload)

	if err != nil {
		log.WithFields(log.Fields{
			"error":      err.Error(),
			"middleware": AppConfig.middleware,
		}).Error("Error during middleware transformation, not modifying payload!")

		return err
	} else {

		log.WithFields(log.Fields{
			"middleware": AppConfig.middleware,
			"newPayload": newPayload,
		}).Info("Middleware transformation complete!")
		// override payload with transformed new payload
		c.payload = newPayload

		return nil
	}
}

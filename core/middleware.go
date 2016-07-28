package hoverfly

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/core/models"
	"io/ioutil"
	"net/http"
	"errors"
	"github.com/SpectoLabs/hoverfly/core/views"
)

// Pipeline - to provide input to the pipeline, assign an io.Reader to the first's Stdin.
func Pipeline(cmds ...*exec.Cmd) (pipeLineOutput, collectedStandardError []byte, pipeLineError error) {
	// Require at least one command
	if len(cmds) < 1 {
		return nil, nil, nil
	}

	// Collect the output from the command(s)
	var output bytes.Buffer
	var stderr bytes.Buffer

	last := len(cmds) - 1
	for i, cmd := range cmds[:last] {
		// Connect each command's stdin to the previous command's stdout
		var err error
		if cmds[i+1].Stdin, err = cmd.StdoutPipe(); err != nil {
			return nil, nil, err
		}
		// Connect each command's stderr to a buffer
		cmd.Stderr = &stderr
	}

	// Connect the output and error for the last command
	cmds[last].Stdout, cmds[last].Stderr = &output, &stderr

	// Start each command
	for _, cmd := range cmds {
		if err := cmd.Start(); err != nil {
			return output.Bytes(), stderr.Bytes(), err
		}
	}

	// Wait for each command to complete
	for _, cmd := range cmds {
		if err := cmd.Wait(); err != nil {
			return output.Bytes(), stderr.Bytes(), err
		}
	}

	// Return the pipeline output and the collected standard error
	return output.Bytes(), stderr.Bytes(), nil
}

// ExecuteMiddleware - takes command (middleware string) and payload, which is passed to middleware
func ExecuteMiddlewareLocally(middlewares string, payload models.Payload) (models.Payload, error) {

	mws := strings.Split(middlewares, "|")
	var cmdList []*exec.Cmd

	for _, v := range mws {
		commands := strings.Split(strings.TrimSpace(v), " ")

		cmd := exec.Command(commands[0], commands[1:]...)
		cmdList = append(cmdList, cmd)
	}

	// getting payload
	bts, err := json.Marshal(payload.ConvertToPayloadView())

	if log.GetLevel() == log.DebugLevel {
		log.WithFields(log.Fields{
			"middlewares": mws,
			"count":       len(mws),
			"payload":     string(bts),
		}).Debug("preparing to modify payload")
	}

	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Failed to marshal json")
		return payload, err
	}

	//
	cmdList[0].Stdin = bytes.NewReader(bts)

	// Run the pipeline
	mwOutput, stderr, err := Pipeline(cmdList...)

	// middleware failed to execute
	if err != nil {
		if len(stderr) > 0 {
			log.WithFields(log.Fields{
				"sdtderr": string(stderr),
				"error":   err.Error(),
			}).Error("Middleware error")
		} else {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Error("Middleware error")
		}
		return payload, err
	}

	// log stderr, middleware executed successfully
	if len(stderr) > 0 {
		log.WithFields(log.Fields{
			"sdtderr": string(stderr),
		}).Info("Information from middleware")
	}

	if len(mwOutput) > 0 {
		var newPayloadView views.PayloadView

		err = json.Unmarshal(mwOutput, &newPayloadView)

		if err != nil {
			log.WithFields(log.Fields{
				"mwOutput": string(mwOutput),
				"error":    err.Error(),
			}).Error("Failed to unmarshal JSON from middleware")
		} else {
			if log.GetLevel() == log.DebugLevel {
				log.WithFields(log.Fields{
					"middlewares": middlewares,
					"count":       len(middlewares),
					"payload":     string(mwOutput),
				}).Debug("payload after modifications")
			}
			// payload unmarshalled into Payload struct, returning it
			return models.NewPayloadFromPayloadView(newPayloadView), nil
		}
	} else {

		log.WithFields(log.Fields{
			"mwOutput": string(mwOutput),
		}).Warn("No response from middleware.")
	}

	return payload, nil

}

func ExecuteMiddlewareRemotely(middleware string, payload models.Payload) (models.Payload, error) {
	bts, err := json.Marshal(payload.ConvertToPayloadView())

	req, err := http.NewRequest("POST", middleware, bytes.NewBuffer(bts))
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Error when building request to remote middleware")
		return payload, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Error when communicating with remote middleware")
		return payload, err
	}

	if resp.StatusCode != 200 {
		log.Error("Remote middleware did not process payload")
		return payload, errors.New("Error when communicating with remote middleware")
	}

	newPayloadBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Error when process response from remote middleware")
		return payload, err
	}

	var newPayloadView views.PayloadView

	err = json.Unmarshal(newPayloadBytes, &newPayloadView)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Error when trying to serialize response from remote middleware")
		return payload, err
	}
	return models.NewPayloadFromPayloadView(newPayloadView), nil
}
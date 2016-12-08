package hoverfly

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"strings"

	"errors"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/models"
)

type Middleware struct {
	FullCommand string
}

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
func (this Middleware) ExecuteMiddlewareLocally(pair models.RequestResponsePair) (models.RequestResponsePair, error) {

	mws := strings.Split(this.FullCommand, "|")
	var cmdList []*exec.Cmd

	for _, v := range mws {
		commands := strings.Split(strings.TrimSpace(v), " ")

		cmd := exec.Command(commands[0], commands[1:]...)
		cmdList = append(cmdList, cmd)
	}

	// getting payload
	pairViewBytes, err := json.Marshal(pair.ConvertToRequestResponsePairView())

	if log.GetLevel() == log.DebugLevel {
		log.WithFields(log.Fields{
			"middlewares": mws,
			"count":       len(mws),
			"payload":     string(pairViewBytes),
		}).Debug("preparing to modify payload")
	}

	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Failed to marshal json")
		return pair, err
	}

	//
	cmdList[0].Stdin = bytes.NewReader(pairViewBytes)

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
		return pair, err
	}

	// log stderr, middleware executed successfully
	if len(stderr) > 0 {
		log.WithFields(log.Fields{
			"sdtderr": string(stderr),
		}).Info("Information from middleware")
	}

	if len(mwOutput) > 0 {
		var newPairView v2.RequestResponsePairView

		err = json.Unmarshal(mwOutput, &newPairView)

		if err != nil {
			log.WithFields(log.Fields{
				"mwOutput": string(mwOutput),
				"error":    err.Error(),
			}).Error("Failed to unmarshal JSON from middleware")
		} else {
			if log.GetLevel() == log.DebugLevel {
				log.WithFields(log.Fields{
					"middleware": this.FullCommand,
					"count":      len(this.FullCommand),
					"payload":    string(mwOutput),
				}).Debug("payload after modifications")
			}
			// payload unmarshalled into RequestResponsePair struct, returning it
			return models.NewRequestResponsePairFromRequestResponsePairView(newPairView), nil
		}
	} else {

		log.WithFields(log.Fields{
			"mwOutput": string(mwOutput),
		}).Warn("No response from middleware.")
	}

	return pair, nil

}

func (this Middleware) ExecuteMiddlewareRemotely(pair models.RequestResponsePair) (models.RequestResponsePair, error) {
	pairViewBytes, err := json.Marshal(pair.ConvertToRequestResponsePairView())

	req, err := http.NewRequest("POST", this.FullCommand, bytes.NewBuffer(pairViewBytes))
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Error when building request to remote middleware")
		return pair, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Error when communicating with remote middleware")
		return pair, err
	}

	if resp.StatusCode != 200 {
		log.Error("Remote middleware did not process payload")
		return pair, errors.New("Error when communicating with remote middleware")
	}

	returnedPairViewBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Error when process response from remote middleware")
		return pair, err
	}

	var newPairView v2.RequestResponsePairView

	err = json.Unmarshal(returnedPairViewBytes, &newPairView)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Error when trying to serialize response from remote middleware")
		return pair, err
	}
	return models.NewRequestResponsePairFromRequestResponsePairView(newPairView), nil
}

func (this Middleware) IsLocal() bool {
	return !strings.HasPrefix(this.FullCommand, "http")
}

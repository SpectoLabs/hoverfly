package hoverfly

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"strings"

	log "github.com/Sirupsen/logrus"
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
func ExecuteMiddleware(command string, payload Payload) (Payload, error) {
	commands := strings.Split(command, " ")

	cmds := exec.Command(commands[0], commands[1:]...)

	// getting payload
	bts, err := json.Marshal(payload)

	if log.GetLevel() == log.DebugLevel {
		log.WithFields(log.Fields{
			"commands":   commands,
			"noCommands": len(commands),
			"payload":    string(bts),
		}).Debug("preparing to modify payload")
	}

	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Failed to marshal json")
		return payload, err
	}
	cmds.Stdin = bytes.NewReader(bts)

	// Run the pipeline
	mwOutput, stderr, err := Pipeline(cmds)

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
		var newPayload Payload

		err = json.Unmarshal(mwOutput, &newPayload)

		if err != nil {
			log.WithFields(log.Fields{
				"mwOutput": string(mwOutput),
				"error":    err.Error(),
			}).Error("Failed to unmarshal JSON from middleware")
		} else {
			if log.GetLevel() == log.DebugLevel {
				log.WithFields(log.Fields{
					"commands":   commands,
					"noCommands": len(commands),
					"payload":    string(mwOutput),
				}).Debug("payload after modifications")
			}
			// payload unmarshalled into Payload struct, returning it
			return newPayload, nil
		}
	} else {

		log.WithFields(log.Fields{
			"mwOutput": string(mwOutput),
		}).Warn("No response from middleware.")
	}

	return payload, nil

}

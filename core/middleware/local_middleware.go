package middleware

import (
	"bytes"
	"encoding/json"
	"os/exec"

	log "github.com/sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/core/models"
)

// ExecuteMiddleware - takes command (middleware string) and payload, which is passed to middleware
func (this Middleware) executeMiddlewareLocally(pair models.RequestResponsePair) (models.RequestResponsePair, error) {
	commandAndArgs := []string{this.Binary, this.Script.Name()}

	middlewareCommand := exec.Command(commandAndArgs[0], commandAndArgs[1:]...)

	pairViewBytes, err := json.Marshal(pair.ConvertToRequestResponsePairView())
	if err != nil {
		return pair, &MiddlewareError{
			OriginalError: err,
			Message:       "Failed to marshal request to JSON",
		}
	}

	log.WithFields(log.Fields{
		"command": this.toString(),
		"stdin":   string(pairViewBytes),
	}).Info("Preparing to execute local middleware")

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	// Redirect standard streams
	middlewareCommand.Stdin = bytes.NewReader(pairViewBytes)
	middlewareCommand.Stdout = &stdout
	middlewareCommand.Stderr = &stderr

	if err := middlewareCommand.Start(); err != nil {
		log.WithFields(log.Fields{
			"sdtdout": string(stdout.Bytes()),
			"sdtderr": string(stderr.Bytes()),
			"error":   err.Error(),
		}).Error("Middleware failed to start")
		return pair, &MiddlewareError{
			OriginalError: err,
			Message:       "Middleware failed to start",
			Command:       this.toString(),
			Stdin:         string(pairViewBytes),
			Stdout:        string(stdout.Bytes()),
			Stderr:        string(stderr.Bytes()),
		}
	}

	if err := middlewareCommand.Wait(); err != nil {
		log.WithFields(log.Fields{
			"command": this.toString(),
			"stdin":   string(pairViewBytes),
			"sdtdout": string(stdout.Bytes()),
			"sdtderr": string(stderr.Bytes()),
			"error":   err.Error(),
		}).Error("Middleware failed")
		return pair, &MiddlewareError{
			OriginalError: err,
			Message:       "Middleware failed",
			Command:       this.toString(),
			Stdin:         string(pairViewBytes),
			Stdout:        string(stdout.Bytes()),
			Stderr:        string(stderr.Bytes()),
		}
	}

	// log stderr, middleware executed successfully
	if len(stderr.Bytes()) > 0 {
		log.WithFields(log.Fields{
			"sdtderr": string(stderr.Bytes()),
		}).Info("Information from middleware")
	}

	if len(stdout.Bytes()) > 0 {
		var newPairView RequestResponsePairView

		err = json.Unmarshal(stdout.Bytes(), &newPairView)

		if err != nil {
			return pair, &MiddlewareError{
				OriginalError: err,
				Message:       "Failed to unmarshal JSON from middleware",
				Command:       this.toString(),
				Stdin:         string(pairViewBytes),
				Stdout:        string(stdout.Bytes()),
				Stderr:        string(stderr.Bytes()),
			}
		} else {
			if log.GetLevel() == log.DebugLevel {
				log.WithFields(log.Fields{
					"middleware": this.toString(),
					"payload":    string(stdout.Bytes()),
				}).Debug("payload after modifications")
			}
			// payload unmarshalled into RequestResponsePair struct, returning it
			return models.NewRequestResponsePairFromRequestResponsePairView(newPairView), nil
		}
	} else {
		log.WithFields(log.Fields{
			"stdout": string(stdout.Bytes()),
		}).Warn("No response from middleware.")
	}

	return pair, nil

}

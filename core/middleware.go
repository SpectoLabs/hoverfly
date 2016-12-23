package hoverfly

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"strings"

	"errors"
	"io/ioutil"
	"net/http"

	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/models"
)

type Middleware struct {
	Binary      string
	Script      *os.File
	Remote      string
	FullCommand string
}

func (this *Middleware) SetScript(scriptContent string) error {
	script, err := ioutil.TempFile(os.TempDir(), "hoverfly_")
	if err != nil {
		return err
	}

	_, err = script.Write([]byte(scriptContent))
	if err != nil {
		return err
	}

	this.DeleteScript()

	this.Script = script

	return nil
}

func (this Middleware) GetScript() (string, error) {
	if this.Script == nil {
		return "", nil
	}
	contents, err := ioutil.ReadFile(this.Script.Name())
	if err != nil {
		return "", err
	}

	return string(contents), nil
}

func (this *Middleware) DeleteScript() error {
	if this.Script != nil {
		err := os.Remove(this.Script.Name())
		if err != nil {
			return err
		}
		this.Script = nil
	}
	return nil
}

func (this *Middleware) SetBinary(binary string) error {
	if binary == "" {
		this.Binary = ""
		return nil
	}
	testCommand := exec.Command(binary)
	if err := testCommand.Start(); err != nil {
		return err
	}

	testCommand.Process.Kill()

	this.Binary = binary
	return nil
}

func (this *Middleware) SetRemote(remoteUrl string) error {
	if remoteUrl == "" {
		this.Remote = ""
		return nil
	}

	response, err := http.Post(remoteUrl, "", nil)
	if err != nil || response.StatusCode != 200 {
		return fmt.Errorf("Could not reach remote middleware")
	}
	this.Remote = remoteUrl
	return nil
}

func (this Middleware) IsSet() bool {
	return this.Binary != "" || this.Remote != ""
}

func (this *Middleware) Execute(pair models.RequestResponsePair) (models.RequestResponsePair, error) {
	if strings.HasPrefix(this.FullCommand, "http") {
		this.Remote = this.FullCommand
	}

	if this.Remote == "" {
		return this.executeMiddlewareLocally(pair)
	} else {
		return this.executeMiddlewareRemotely(pair)
	}
}

// ExecuteMiddleware - takes command (middleware string) and payload, which is passed to middleware
func (this Middleware) executeMiddlewareLocally(pair models.RequestResponsePair) (models.RequestResponsePair, error) {
	var commandAndArgs []string
	if this.Binary == "" {
		commandAndArgs = strings.Split(strings.TrimSpace(this.FullCommand), " ")
	} else {
		commandAndArgs = []string{this.Binary, this.Script.Name()}
	}

	middlewareCommand := exec.Command(commandAndArgs[0], commandAndArgs[1:]...)

	// getting payload
	pairViewBytes, err := json.Marshal(pair.ConvertToRequestResponsePairView())

	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Failed to marshal json")
		return pair, err
	}

	if log.GetLevel() == log.DebugLevel {
		log.WithFields(log.Fields{
			"middleware": this.FullCommand,
			"stdin":      string(pairViewBytes),
		}).Debug("preparing to modify payload")
	}

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
		return pair, err
	}

	if err := middlewareCommand.Wait(); err != nil {
		log.WithFields(log.Fields{
			"sdtdout": string(stdout.Bytes()),
			"sdtderr": string(stderr.Bytes()),
			"error":   err.Error(),
		}).Error("Middleware failed to stop successfully")
		return pair, err
	}

	// log stderr, middleware executed successfully
	if len(stderr.Bytes()) > 0 {
		log.WithFields(log.Fields{
			"sdtderr": string(stderr.Bytes()),
		}).Info("Information from middleware")
	}

	if len(stdout.Bytes()) > 0 {
		var newPairView v2.RequestResponsePairView

		err = json.Unmarshal(stdout.Bytes(), &newPairView)

		if err != nil {
			log.WithFields(log.Fields{
				"stdout": string(stdout.Bytes()),
				"error":  err.Error(),
			}).Error("Failed to unmarshal JSON from middleware")
		} else {
			if log.GetLevel() == log.DebugLevel {
				log.WithFields(log.Fields{
					"middleware": this.FullCommand,
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

func (this Middleware) executeMiddlewareRemotely(pair models.RequestResponsePair) (models.RequestResponsePair, error) {
	pairViewBytes, err := json.Marshal(pair.ConvertToRequestResponsePairView())

	if this.Remote == "" {
		return pair, fmt.Errorf("Error when communicating with remote middleware")
	}

	req, err := http.NewRequest("POST", this.Remote, bytes.NewBuffer(pairViewBytes))
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

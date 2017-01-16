package hoverfly

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path"
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
	Binary string
	Script *os.File
	Remote string
}

func ConvertToNewMiddleware(middleware string) (*Middleware, error) {
	newMiddleware := &Middleware{}
	if strings.HasPrefix(middleware, "http") {

		err := newMiddleware.SetRemote(middleware)
		if err != nil {
			return nil, err
		}

		return newMiddleware, nil
	} else if strings.Contains(middleware, " ") {
		splitMiddleware := strings.Split(middleware, " ")
		fileContents, _ := ioutil.ReadFile(splitMiddleware[1])

		newMiddleware.SetBinary(splitMiddleware[0])
		newMiddleware.SetScript(string(fileContents))

		return newMiddleware, nil

	} else {
		err := newMiddleware.SetBinary(middleware)
		if err != nil {
			return nil, err
		}
		return newMiddleware, nil
	}

	return nil, nil
}

func (this *Middleware) SetScript(scriptContent string) error {
	tempDir := path.Join(os.TempDir(), "hoverfly")
	this.DeleteScripts(tempDir)

	//We ignore the error it outputs as this directory may already exist
	os.Mkdir(tempDir, 0777)

	script, err := ioutil.TempFile(tempDir, "hoverfly_")
	if err != nil {
		return err
	}

	_, err = script.Write([]byte(scriptContent))
	if err != nil {
		return err
	}

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

func (this *Middleware) DeleteScripts(path string) error {
	err := os.RemoveAll(path)
	if err != nil {
		return err
	}
	this.Script = nil

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

func (this *Middleware) Execute(pair models.RequestResponsePair) (models.RequestResponsePair, error) {
	if !this.IsSet() {
		return pair, fmt.Errorf("Cannot execute middleware as middleware has not been correctly set")
	}

	if this.Remote == "" {
		return this.executeMiddlewareLocally(pair)
	} else {
		return this.executeMiddlewareRemotely(pair)
	}
}

// ExecuteMiddleware - takes command (middleware string) and payload, which is passed to middleware
func (this Middleware) executeMiddlewareLocally(pair models.RequestResponsePair) (models.RequestResponsePair, error) {
	commandAndArgs := []string{this.Binary, this.Script.Name()}

	middlewareCommand := exec.Command(commandAndArgs[0], commandAndArgs[1:]...)

	// getting payload
	pairViewBytes, err := json.Marshal(pair.ConvertToRequestResponsePairView())

	if err != nil {
		return pair, errors.New("Failed to marshal request to JSON")
	}

	log.WithFields(log.Fields{
		"middleware": this.toString(),
		"stdin":      string(pairViewBytes),
	}).Debug("preparing to modify payload")

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
			return pair, errors.New("Failed to unmarshal JSON from middleware")
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

func (this Middleware) IsSet() bool {
	return this.Binary != "" || this.Remote != ""
}

func (this Middleware) toString() string {
	if this.Remote != "" {
		return this.Remote
	} else {
		if this.Script != nil {
			return this.Binary + " " + this.Script.Name()
		}
		return this.Binary
	}
}

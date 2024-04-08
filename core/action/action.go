package action

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	v2 "github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/journal"
	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/pborman/uuid"
	log "github.com/sirupsen/logrus"
)

type Action struct {
	Binary    string
	Script    *os.File
	Remote    string
	DelayInMs int
}

func NewLocalAction(actionName, binary, scriptContent string, delayInMs int) (*Action, error) {

	scriptInfo := &Action{}
	if strings.TrimSpace(actionName) == "" {
		return nil, errors.New("empty action name passed")
	}

	scriptInfo.DelayInMs = delayInMs

	if err := setBinary(scriptInfo, binary); err != nil {
		return nil, err
	}

	if err := setScript(scriptInfo, scriptContent); err != nil {
		return nil, err
	}
	return scriptInfo, nil
}

func NewRemoteAction(actionName, host string, delayInMs int) (*Action, error) {

	if strings.TrimSpace(actionName) == "" {
		return nil, errors.New("empty action name passed")
	}

	if !isValidURL(host) {
		return nil, errors.New("remote host is invalid")
	}

	return &Action{Remote: host, DelayInMs: delayInMs}, nil
}

func setBinary(action *Action, binary string) error {
	action.Binary = binary
	return nil
}

func setScript(action *Action, scriptContent string) error {
	tempDir := path.Join(os.TempDir(), "hoverfly")
	os.Mkdir(tempDir, 0777)

	newScript, err := ioutil.TempFile(tempDir, "hoverfly_")
	if err != nil {
		return err
	}

	_, err = newScript.Write([]byte(scriptContent))
	if err != nil {
		return err
	}

	action.Script = newScript
	return nil
}

func (action *Action) DeleteScript() error {
	if action.Script == nil {
		return nil
	}
	err := os.Remove(action.Script.Name())
	if err != nil {
		return err
	}
	action.Script = nil

	return nil
}

func (action *Action) GetScript() (string, error) {
	if action.Script == nil {
		return "", nil
	}
	contents, err := ioutil.ReadFile(action.Script.Name())
	if err != nil {
		return "", err
	}

	return string(contents), nil
}

func (action *Action) GetActionView(actionName string) v2.ActionView {

	scriptContent, _ := action.GetScript()
	return v2.ActionView{
		ActionName:    actionName,
		Binary:        action.Binary,
		ScriptContent: scriptContent,
		Remote:        action.Remote,
		DelayInMs:     action.DelayInMs,
	}
}

func (action *Action) Execute(pair *models.RequestResponsePair, journalIDChannel chan string, journal *journal.Journal) error {

	pairViewBytes, err := json.Marshal(pair.ConvertToRequestResponsePairView())
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"stdin": string(pairViewBytes),
	}).Info("Delaying to execute post serve action")

	//adding 200 ms to include some buffer for it to return response
	time.Sleep(time.Duration(200+action.DelayInMs) * time.Millisecond)

	journalID := <-journalIDChannel
	close(journalIDChannel)
	log.Info("Journal ID received ", journalID)

	//if it is remote callback
	if action.Remote != "" {

		req, err := http.NewRequest("POST", action.Remote, bytes.NewBuffer(pairViewBytes))
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Error("Error when building request to remote post serve action")
			return err
		}

		correlationID := uuid.New()
		invokedTime := time.Now()
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("X-CORRELATION-ID", correlationID)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Error("Error when communicating with remote post serve action")
			return err
		}

		completionTime := time.Now()
		journal.UpdatePostServeActionDetailsInJournal(journalID, pair.Response.PostServeAction, correlationID, invokedTime, completionTime, resp.StatusCode)
		if resp.StatusCode != 200 {
			log.Error("Remote post serve action did not process payload")

			return nil
		}
		log.Info("Remote post serve action invoked successfully")
		return nil
	}

	invokedTime := time.Now()
	actionCommand := exec.Command(action.Binary, action.Script.Name())
	actionCommand.Stdin = bytes.NewReader(pairViewBytes)
	completionTime := time.Now()

	journal.UpdatePostServeActionDetailsInJournal(journalID, pair.Response.PostServeAction, "", invokedTime, completionTime, 0)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	actionCommand.Stdout = &stdout
	actionCommand.Stderr = &stderr
	if err := actionCommand.Start(); err != nil {
		return err
	}

	if err := actionCommand.Wait(); err != nil {
		return err
	}

	if len(stderr.Bytes()) > 0 {
		log.Error("Error occurred while executing script " + stderr.String())
	}

	if len(stdout.Bytes()) > 0 {
		log.Info("Output from post serve action " + stdout.String())
	}
	return nil
}

func isValidURL(host string) bool {

	if _, err := url.ParseRequestURI(host); err == nil {
		return true
	}
	return false
}

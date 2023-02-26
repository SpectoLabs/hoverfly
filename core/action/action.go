package action

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	v2 "github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/models"
	log "github.com/sirupsen/logrus"
)

type Action struct {
	Binary    string
	Script    *os.File
	DelayInMs int
}

func NewAction(actionName, binary, scriptContent string, delayInMs int) (*Action, error) {

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
		DelayInMs:     action.DelayInMs,
	}
}

func (action *Action) ExecuteLocally(pair *models.RequestResponsePair) error {

	pairViewBytes, err := json.Marshal(pair.ConvertToRequestResponsePairView())
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"stdin": string(pairViewBytes),
	}).Info("Delaying to execute post serve action")

	//adding 200 ms to include some buffer for it to return response
	time.Sleep(time.Duration(200+action.DelayInMs) * time.Millisecond)

	actionCommand := exec.Command(action.Binary, action.Script.Name())
	actionCommand.Stdin = bytes.NewReader(pairViewBytes)
	var stdout bytes.Buffer
	actionCommand.Stdout = &stdout
	if err := actionCommand.Start(); err != nil {
		return err
	}

	if err := actionCommand.Wait(); err != nil {
		return err
	}

	if len(stdout.Bytes()) > 0 {
		log.Info("Output from post serve action " + stdout.String())
	}
	return nil
}

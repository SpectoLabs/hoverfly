package hook

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	v2 "github.com/SpectoLabs/hoverfly/core/handlers/v2"
	log "github.com/sirupsen/logrus"
)

type Hook struct {
	Binary              string
	Script              *os.File
	DelayInMilliSeconds int
}

func NewHook(hookName, binary, scriptContent string, delayInMilliSeconds int) (*Hook, error) {

	scriptInfo := &Hook{}
	if strings.TrimSpace(hookName) == "" {
		return nil, errors.New("empty hook name passed")
	}

	scriptInfo.DelayInMilliSeconds = delayInMilliSeconds

	if err := setBinary(scriptInfo, binary); err != nil {
		return nil, err
	}

	if err := setScript(scriptInfo, hookName, scriptContent); err != nil {
		return nil, err
	}
	return scriptInfo, nil
}

func setBinary(hook *Hook, binary string) error {
	if binary == "" {
		hook.Binary = ""
		return nil
	}
	hook.Binary = binary
	return nil
}

func setScript(hook *Hook, hookName, scriptContent string) error {
	tempDir := path.Join(os.TempDir(), "hoverfly")
	os.Mkdir(tempDir, 0777)

	newScript, err := ioutil.TempFile(tempDir, "hoverfly_"+hookName)
	if err != nil {
		return err
	}

	_, err = newScript.Write([]byte(scriptContent))
	if err != nil {
		return err
	}

	hook.Script = newScript
	return nil
}

func (hook *Hook) DeleteScript() error {
	if hook.Script == nil {
		return nil
	}
	err := os.Remove(hook.Script.Name())
	if err != nil {
		return err
	}
	hook.Script = nil

	return nil
}

func (hook *Hook) GetScript() (string, error) {
	if hook.Script == nil {
		return "", nil
	}
	contents, err := ioutil.ReadFile(hook.Script.Name())
	if err != nil {
		return "", err
	}

	return string(contents), nil
}

func (hook *Hook) GetHookView(hookName string) v2.HookView {

	scriptContent, _ := hook.GetScript()
	return v2.HookView{
		HookName:            hookName,
		Binary:              hook.Binary,
		ScriptContent:       scriptContent,
		DelayInMilliSeconds: hook.DelayInMilliSeconds,
	}
}

func (hook *Hook) ExecuteLocally() error {

	//adding 200 ms to include some buffer for it to return response
	time.Sleep(time.Duration(200+hook.DelayInMilliSeconds) * time.Millisecond)

	hookCommand := exec.Command(hook.Binary, hook.Script.Name())
	var stdout bytes.Buffer
	hookCommand.Stdout = &stdout
	if err := hookCommand.Start(); err != nil {
		return err
	}

	if err := hookCommand.Wait(); err != nil {
		return err
	}

	if len(stdout.Bytes()) > 0 {
		log.Info("Output from post serve action hook" + stdout.String())
	}
	return nil
}

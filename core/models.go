package hoverfly

import (
	"github.com/SpectoLabs/hoverfly/core/models"
	"io/ioutil"
	"net/http"
	"time"
)

var emptyResp = &http.Response{}

func extractBody(resp *http.Response) (extract []byte, err error) {
	save := resp.Body
	savecl := resp.ContentLength

	save, resp.Body, err = models.CopyBody(resp.Body)

	if err != nil {
		return
	}
	defer resp.Body.Close()
	extract, err = ioutil.ReadAll(resp.Body)

	resp.Body = save
	resp.ContentLength = savecl
	if err != nil {
		return nil, err
	}
	return extract, nil
}
// ActionType - action type can be things such as "RequestCaptured", "GotResponse" - anything
type ActionType string

// ActionTypeRequestCaptured - default action type name for identifying
const ActionTypeRequestCaptured = "requestCaptured"

// ActionTypeWipeDB - default action type for wiping database
const ActionTypeWipeDB = "wipeDatabase"

// ActionTypeConfigurationChanged - default action name for identifying configuration changes
const ActionTypeConfigurationChanged = "configurationChanged"

// Entry - holds information about action, based on action type - other clients will be able to decode
// the data field.
type Entry struct {
	// Contains encoded data
	Data []byte

	// Time at which the action entry was fired
	Time time.Time

	ActionType ActionType

	// Message, can carry additional information
	Message string
}

// Hook - an interface to add dynamic hooks to extend functionality
type Hook interface {
	ActionTypes() []ActionType
	Fire(*Entry) error
}

// ActionTypeHooks type for storing the hooks
type ActionTypeHooks map[ActionType][]Hook

// Add a hook
func (hooks ActionTypeHooks) Add(hook Hook) {
	for _, ac := range hook.ActionTypes() {
		hooks[ac] = append(hooks[ac], hook)
	}
}

// Fire all the hooks for the passed ActionType
func (hooks ActionTypeHooks) Fire(ac ActionType, entry *Entry) error {
	for _, hook := range hooks[ac] {
		if err := hook.Fire(entry); err != nil {
			return err
		}
	}

	return nil
}

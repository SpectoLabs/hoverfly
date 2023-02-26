package action

import (
	"fmt"
	"sync"
)

type PostServeActionDetails struct {
	Actions map[string]Action
	RWMutex sync.RWMutex
}

func NewPostServeActionDetails() *PostServeActionDetails {

	return &PostServeActionDetails{
		Actions: make(map[string]Action),
	}
}

func (postServeActionDetails *PostServeActionDetails) SetAction(actionName string, newAction *Action) error {

	postServeActionDetails.RWMutex.Lock()
	//cleanup
	if existingAction, ok := postServeActionDetails.Actions[actionName]; ok {
		existingAction.DeleteScript()
		delete(postServeActionDetails.Actions, actionName)
	}
	postServeActionDetails.Actions[actionName] = *newAction
	postServeActionDetails.RWMutex.Unlock()
	return nil
}

func (postServeActionDetails *PostServeActionDetails) DeleteAction(actionName string) error {

	if _, ok := postServeActionDetails.Actions[actionName]; !ok {
		return fmt.Errorf("invalid action name passed")
	}

	postServeActionDetails.RWMutex.Lock()
	action := postServeActionDetails.Actions[actionName]
	action.DeleteScript()
	delete(postServeActionDetails.Actions, actionName)
	postServeActionDetails.RWMutex.Unlock()
	return nil
}

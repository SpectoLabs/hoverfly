package hook

import (
	"fmt"
	"sync"
)

type PostServeActionDetails struct {
	Hooks   map[string]Hook
	RWMutex sync.RWMutex
}

func NewPostServeActionDetails() *PostServeActionDetails {

	return &PostServeActionDetails{
		Hooks: make(map[string]Hook),
	}
}

func (postServeActionDetails *PostServeActionDetails) AddHook(hookName string, hook *Hook) error {

	if _, ok := postServeActionDetails.Hooks[hookName]; ok {
		return fmt.Errorf("hook with this name already exists")
	}

	postServeActionDetails.RWMutex.Lock()
	postServeActionDetails.Hooks[hookName] = *hook
	postServeActionDetails.RWMutex.Unlock()
	return nil
}

func (postServeActionDetails *PostServeActionDetails) DeleteHook(hookName string) error {

	if _, ok := postServeActionDetails.Hooks[hookName]; !ok {
		return fmt.Errorf("invalid hookname passed")
	}

	postServeActionDetails.RWMutex.Lock()
	hook := postServeActionDetails.Hooks[hookName]
	hook.DeleteScript()
	delete(postServeActionDetails.Hooks, hookName)
	postServeActionDetails.RWMutex.Unlock()
	return nil
}

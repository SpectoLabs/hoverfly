package state

import (
	"fmt"
	"strings"
	"sync"
)

type State struct {
	State   map[string]string
	RWMutex sync.RWMutex
}

func NewState() *State {
	return &State{
		State: map[string]string{},
	}
}

func NewStateFromState(incomingState map[string]string) *State {
	state := &State{
		State: map[string]string{},
	}

	for stateKey := range incomingState {
		if strings.Contains(stateKey, "sequence:") {
			state.State[stateKey] = "1"
		}
	}

	return state
}

func (s *State) GetState(key string) string {
	s.RWMutex.RLock()
	val := s.State[key]
	s.RWMutex.RUnlock()
	return val
}

func (s *State) SetState(state map[string]string) {
	s.RWMutex.Lock()
	s.State = state
	s.RWMutex.Unlock()
}

func (s *State) PatchState(toPatch map[string]string) {
	s.RWMutex.Lock()
	for k, v := range toPatch {
		s.State[k] = v
	}
	s.RWMutex.Unlock()
}

func (s *State) RemoveState(toRemove []string) {
	s.RWMutex.Lock()
	for _, key := range toRemove {
		delete(s.State, key)
	}
	s.RWMutex.Unlock()
}

func (s *State) GetNewSequenceKey() string {
	returnKey := ""
	i := 1
	s.RWMutex.RLock()
	for returnKey == "" {
		tempKey := fmt.Sprintf("sequence:%v", i)
		if s.State[tempKey] == "" {
			returnKey = tempKey
		} else {
			i = i + 1
		}
	}
	s.RWMutex.RUnlock()
	return returnKey
}

package state

import (
	"fmt"
	"strings"
)

type State struct {
	State map[string]string
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

	for stateKey, _ := range incomingState {
		if strings.Contains(stateKey, "sequence:") {
			state.State[stateKey] = "1"
		}
	}

	return state
}

func (s *State) GetState(key string) string {
	return s.State[key]
}

func (s *State) SetState(state map[string]string) {
	s.State = state
}

func (s *State) PatchState(toPatch map[string]string) {
	for k, v := range toPatch {
		s.State[k] = v
	}
}

func (s *State) RemoveState(toRemove []string) {
	for _, key := range toRemove {
		delete(s.State, key)
	}
}

func (s *State) GetNewSequenceKey() string {
	returnKey := ""
	i := 0
	for returnKey == "" {
		tempKey := fmt.Sprintf("sequence:%v", i)
		if s.State[tempKey] == "" {
			returnKey = tempKey
		} else {
			i = i + 1
		}
	}
	return returnKey
}

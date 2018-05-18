package state

type State struct {
	State map[string]string
}

func NewState() *State {
	return &State{
		State: map[string]string{},
	}
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

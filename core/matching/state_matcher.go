package matching

import (
	"github.com/SpectoLabs/hoverfly/core/state"
	"github.com/SpectoLabs/hoverfly/core/util"
)

func StateMatcher(currentState *state.State, requiredState map[string]string) *FieldMatch {

	score := 0
	matched := true

	if requiredState == nil || len(requiredState) == 0 {
		return &FieldMatch{
			Matched: true,
			Score:   0,
		}
	}

	currentState.RWMutex.RLock()
	copyState := util.CopyMap(currentState.State)
	currentState.RWMutex.RUnlock()
	for key, value := range requiredState {
		if _, ok := copyState[key]; !ok {
			matched = false
		}
		if copyState[key] != value {
			matched = false
		} else {
			score++
		}
	}

	return &FieldMatch{
		Matched: matched,
		Score:   score,
	}
}

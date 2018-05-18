package matching

import "github.com/SpectoLabs/hoverfly/core/state"

func StateMatcher(currentState *state.State, requiredState map[string]string) *FieldMatch {

	score := 0
	matched := true

	if requiredState == nil || len(requiredState) == 0 {
		return &FieldMatch{
			Matched: true,
			Score:   0,
		}
	}

	for key, value := range requiredState {
		if _, ok := currentState.State[key]; !ok {
			matched = false
		}
		if currentState.State[key] != value {
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

package matching

func StateMatcher(currentState, requiredState map[string]string) *FieldMatch {

	score := 0
	matched := true

	if requiredState == nil || len(requiredState) == 0 {
		return FieldMatchWithNoScore(true)
	}

	for key, value := range requiredState {
		if _, ok := currentState[key]; !ok {
			matched = false
		}
		if currentState[key] != value {
			matched = false
		} else {
			score++
		}
	}

	return &FieldMatch{
		Matched:    matched,
		MatchScore: score,
	}
}

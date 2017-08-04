package matching

func UnscoredStateMatcher(currentState, requiredState map[string]string) *FieldMatch {
	if requiredState == nil || len(requiredState) == 0 {
		return FieldMatchWithNoScore(true)
	}

	if currentState == nil || len(currentState) == 0 {
		return FieldMatchWithNoScore(false)
	}

	if len(requiredState) > len(currentState) {
		return FieldMatchWithNoScore(false)
	}

	for key, value := range requiredState {
		if _, ok := currentState[key]; !ok {
			return FieldMatchWithNoScore(false)
		}
		if currentState[key] != value {
			return FieldMatchWithNoScore(false)
		}
	}

	return FieldMatchWithNoScore(true)
}

func ScoredStateMatcher(currentState, requiredState map[string]string) *FieldMatch {

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

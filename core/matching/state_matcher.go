package matching

func StateMatcher(copyState map[string]string, requiredState map[string]string) *FieldMatch {

	score := 0
	matched := true

	if requiredState == nil || len(requiredState) == 0 {
		return &FieldMatch{
			Matched: true,
			Score:   0,
		}
	}

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

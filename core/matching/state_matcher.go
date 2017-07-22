package matching


func UnscoredStateMatcher(currentState map[string]string, requiredState map[string]string) *FieldMatch {
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

//func ScoredFieldMatcher(currentState map[string]string, requiredState map[string]string) *FieldMatch {
//
//}
package matching

import "github.com/NodePrime/jsonpath"

func JsonPathMatch(matchingString string, toMatch string) bool {
	paths, err := jsonpath.ParsePaths(matchingString)
	if err != nil {
		return false
	}

	eval, err := jsonpath.EvalPathsInBytes([]byte(toMatch), paths)
	if err != nil {
		return false
	}

	_, ok := eval.Next()

	return ok
}

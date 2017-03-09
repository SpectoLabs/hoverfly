package jsonpath

import "io"

func EvalPathsInBytes(input []byte, paths []*Path) (*Eval, error) {
	lexer := NewSliceLexer(input, JSON)
	eval := newEvaluation(lexer, paths...)
	return eval, nil
}

func EvalPathsInReader(r io.Reader, paths []*Path) (*Eval, error) {
	lexer := NewReaderLexer(r, JSON)
	eval := newEvaluation(lexer, paths...)
	return eval, nil
}

func ParsePaths(pathStrings ...string) ([]*Path, error) {
	paths := make([]*Path, len(pathStrings))
	for x, p := range pathStrings {
		path, err := parsePath(p)
		if err != nil {
			return nil, err
		}
		paths[x] = path
	}
	return paths, nil
}

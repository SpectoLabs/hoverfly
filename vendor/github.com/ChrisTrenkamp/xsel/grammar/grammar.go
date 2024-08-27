package grammar

import (
	"bytes"
	"fmt"

	"github.com/ChrisTrenkamp/xsel/grammar/lexer"
	"github.com/ChrisTrenkamp/xsel/grammar/parser"
	"github.com/ChrisTrenkamp/xsel/grammar/parser/bsr"
)

type Grammar struct {
	BSR *bsr.BSR
	lex *lexer.Lexer
}

func (g *Grammar) Next(bsr *bsr.BSR) *Grammar {
	return &Grammar{
		BSR: bsr,
		lex: g.lex,
	}
}

func (g *Grammar) GetString() string {
	return g.lex.GetString(g.BSR.LeftExtent(), g.BSR.RightExtent()-1)
}

func (g *Grammar) GetStringExtents(left, right int) string {
	return g.lex.GetString(left, right-1)
}

type errPos struct {
	line, col int
}

// Creates an XPath query.
func Build(xpath string) (Grammar, error) {
	lex := lexer.New([]rune(xpath))
	parse, err := parser.Parse(lex)

	if len(err) > 0 {
		errorPositions := make(map[errPos][]parser.Error)

		for _, e := range err {
			pos := errPos{line: e.Line, col: e.Column}
			errorPositions[pos] = append(errorPositions[pos], *e)
		}

		errBuf := bytes.Buffer{}

		for k, v := range errorPositions {
			errBuf.WriteString(fmt.Sprintf("Error on line %d, column %d. ", k.line, k.col))
			errBuf.WriteString("Expected one of: ")
			expected := make(map[string]bool)

			for _, e := range v {
				for _, i := range e.Expected {
					expected[i] = true
				}
			}

			for i := range expected {
				errBuf.WriteString(i)
				errBuf.WriteString(" ")
			}

			errBuf.WriteString("\n")
		}

		return Grammar{}, fmt.Errorf(errBuf.String())
	}

	roots := parse.GetRoots()

	if len(roots) == 0 {
		return Grammar{}, fmt.Errorf("could not build expression tree")
	}

	return Grammar{&roots[0], lex}, nil
}

// Like Build, but panics if an error is thrown.
func MustBuild(xpath string) Grammar {
	grammar, err := Build(xpath)

	if err != nil {
		panic(err)
	}

	return grammar
}

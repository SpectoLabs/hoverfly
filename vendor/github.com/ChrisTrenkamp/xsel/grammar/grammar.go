package grammar

import (
	"bytes"
	"fmt"
	"strings"

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

// Creates an XPath query.
func Build(xpath string) (Grammar, error) {
	lex := lexer.New([]rune(xpath))
	parse, err := parser.Parse(lex)

	if err != nil {
		errBuf := bytes.Buffer{}

		for _, e := range err {
			printError(&errBuf, e)
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

func printError(buf *bytes.Buffer, err *parser.Error) {
	if strings.HasPrefix(err.Token.Type().String(), "T_") {
		// This error message isn't useful
		return
	}

	buf.WriteString(fmt.Sprintf("Error on column %d, ", err.Column))
	buf.WriteString(fmt.Sprintf("on token type '%s' - ", err.Slot.Index().NT))
	buf.WriteString(fmt.Sprintf("received %s - ", err.Token.Type()))
	buf.WriteString("Expected: ")

	expected := make([]string, 0, len(err.Expected))

	for _, i := range err.Expected {
		expected = append(expected, i)
	}

	buf.WriteString(strings.Join(expected, ", "))
	buf.WriteString("\n")
}

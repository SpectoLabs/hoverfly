package js // import "github.com/tdewolff/parse/js"

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func helperStringify(t *testing.T, input string) string {
	s := ""
	l := NewLexer(bytes.NewBufferString(input))
	for i := 0; i < 10; i++ {
		tt, data := l.Next()
		if tt == ErrorToken {
			s += tt.String() + "('" + l.Err().Error() + "')"
			break
		} else if tt == WhitespaceToken {
			continue
		} else {
			s += tt.String() + "('" + string(data) + "') "
		}
	}
	return s
}

////////////////////////////////////////////////////////////////

type TTs []TokenType

func TestTokens(t *testing.T) {
	var tokenTests = []struct {
		js       string
		expected TTs
	}{
		{" \t\v\f\u00A0\uFEFF\u2000", TTs{}}, // WhitespaceToken
		{"\n\r\r\n\u2028\u2029", TTs{LineTerminatorToken}},
		{"5.2 .04 0x0F 5e99", TTs{NumericToken, NumericToken, NumericToken, NumericToken}},
		{"a = 'string'", TTs{IdentifierToken, PunctuatorToken, StringToken}},
		{"/*comment*/ //comment", TTs{CommentToken, CommentToken}},
		{"{ } ( ) [ ]", TTs{PunctuatorToken, PunctuatorToken, PunctuatorToken, PunctuatorToken, PunctuatorToken, PunctuatorToken}},
		{". ; , < > <=", TTs{PunctuatorToken, PunctuatorToken, PunctuatorToken, PunctuatorToken, PunctuatorToken, PunctuatorToken}},
		{">= == != === !==", TTs{PunctuatorToken, PunctuatorToken, PunctuatorToken, PunctuatorToken, PunctuatorToken}},
		{"+ - * % ++ --", TTs{PunctuatorToken, PunctuatorToken, PunctuatorToken, PunctuatorToken, PunctuatorToken, PunctuatorToken}},
		{"<< >> >>> & | ^", TTs{PunctuatorToken, PunctuatorToken, PunctuatorToken, PunctuatorToken, PunctuatorToken, PunctuatorToken}},
		{"! ~ && || ? :", TTs{PunctuatorToken, PunctuatorToken, PunctuatorToken, PunctuatorToken, PunctuatorToken, PunctuatorToken}},
		{"= += -= *= %= <<=", TTs{PunctuatorToken, PunctuatorToken, PunctuatorToken, PunctuatorToken, PunctuatorToken, PunctuatorToken}},
		{">>= >>>= &= |= ^= =>", TTs{PunctuatorToken, PunctuatorToken, PunctuatorToken, PunctuatorToken, PunctuatorToken, PunctuatorToken}},
		{"a = /.*/g;", TTs{IdentifierToken, PunctuatorToken, RegexpToken, PunctuatorToken}},

		{"/*co\nm\u2028m/*ent*/ //co//mment\u2029//comment", TTs{CommentToken, CommentToken, LineTerminatorToken, CommentToken}},
		{"$ _\u200C \\u2000 \u200C", TTs{IdentifierToken, IdentifierToken, IdentifierToken, UnknownToken}},
		{">>>=>>>>=", TTs{PunctuatorToken, PunctuatorToken, PunctuatorToken}},
		{"/", TTs{PunctuatorToken}},
		{"/=", TTs{PunctuatorToken}},
		{"010xF", TTs{NumericToken, NumericToken, IdentifierToken}},
		{"50e+-0", TTs{NumericToken, IdentifierToken, PunctuatorToken, PunctuatorToken, NumericToken}},
		{"'str\\i\\'ng'", TTs{StringToken}},
		{"'str\\\\'abc", TTs{StringToken, IdentifierToken}},
		{"'str\\\ni\\\\u00A0ng'", TTs{StringToken}},
		{"a = /[a-z/]/g", TTs{IdentifierToken, PunctuatorToken, RegexpToken}},
		{"a=/=/g1", TTs{IdentifierToken, PunctuatorToken, RegexpToken}},
		{"a = /'\\\\/\n", TTs{IdentifierToken, PunctuatorToken, RegexpToken, LineTerminatorToken}},
		{"a=/\\//g1", TTs{IdentifierToken, PunctuatorToken, RegexpToken}},
		{"new RegExp(a + /\\d{1,2}/.source)", TTs{IdentifierToken, IdentifierToken, PunctuatorToken, IdentifierToken, PunctuatorToken, RegexpToken, PunctuatorToken, IdentifierToken, PunctuatorToken}},

		{"0b0101 0o0707 0b17", TTs{NumericToken, NumericToken, NumericToken, NumericToken}},
		{"`template`", TTs{TemplateToken}},
		{"`a${x+y}b`", TTs{TemplateToken, IdentifierToken, PunctuatorToken, IdentifierToken, TemplateToken}},
		{"`temp\nlate`", TTs{TemplateToken}},

		// early endings
		{"'string", TTs{StringToken}},
		{"'\n '\u2028", TTs{UnknownToken, LineTerminatorToken, UnknownToken, LineTerminatorToken}},
		{"'str\\\U00100000ing\\0'", TTs{StringToken}},
		{"'strin\\00g'", TTs{StringToken}},
		{"/*comment", TTs{CommentToken}},
		{"a=/regexp", TTs{IdentifierToken, PunctuatorToken, RegexpToken}},
		{"\\u002", TTs{UnknownToken, IdentifierToken}},

		// coverage
		{"Ø a〉", TTs{IdentifierToken, IdentifierToken, UnknownToken}},
		{"0xg 0.f", TTs{NumericToken, IdentifierToken, NumericToken, PunctuatorToken, IdentifierToken}},
		{"0bg 0og", TTs{NumericToken, IdentifierToken, NumericToken, IdentifierToken}},
		{"\u00A0\uFEFF\u2000", TTs{}},
		{"\u2028\u2029", TTs{LineTerminatorToken}},
		{"\\u0029ident", TTs{IdentifierToken}},
		{"\\u{0029FEF}ident", TTs{IdentifierToken}},
		{"\\u{}", TTs{UnknownToken, IdentifierToken, PunctuatorToken, PunctuatorToken}},
		{"\\ugident", TTs{UnknownToken, IdentifierToken}},
		{"'str\u2028ing'", TTs{UnknownToken, IdentifierToken, LineTerminatorToken, IdentifierToken, StringToken}},
		{"a=/\\\n", TTs{IdentifierToken, PunctuatorToken, PunctuatorToken, UnknownToken, LineTerminatorToken}},
		{"a=/x/\u200C\u3009", TTs{IdentifierToken, PunctuatorToken, RegexpToken, UnknownToken}},
		{"a=/x\n", TTs{IdentifierToken, PunctuatorToken, PunctuatorToken, IdentifierToken, LineTerminatorToken}},

		// go fuzz
		{"`", TTs{UnknownToken}},
	}
	for _, test := range tokenTests {
		expected := []TokenType(test.expected)
		stringify := helperStringify(t, test.js)
		l := NewLexer(bytes.NewBufferString(test.js))
		i := 0
		for {
			tt, _ := l.Next()
			if tt == ErrorToken {
				assert.Equal(t, io.EOF, l.Err(), "error must be EOF in "+stringify)
				assert.Equal(t, len(expected), i, "when error occurred we must be at the end in "+stringify)
				break
			} else if tt == WhitespaceToken {
				continue
			}
			assert.False(t, i >= len(expected), "index must not exceed expected token types size in "+stringify)
			if i < len(expected) {
				assert.Equal(t, expected[i], tt, "token types must match at index "+strconv.Itoa(i)+" in "+stringify)
			}
			i++
		}
	}

	assert.Equal(t, "Whitespace", WhitespaceToken.String())
	assert.Equal(t, "Invalid(100)", TokenType(100).String())
}

////////////////////////////////////////////////////////////////

func ExampleNewLexer() {
	l := NewLexer(bytes.NewBufferString("var x = 'lorem ipsum';"))
	out := ""
	for {
		tt, data := l.Next()
		if tt == ErrorToken {
			break
		}
		out += string(data)
		l.Free(len(data))
	}
	fmt.Println(out)
	// Output: var x = 'lorem ipsum';
}

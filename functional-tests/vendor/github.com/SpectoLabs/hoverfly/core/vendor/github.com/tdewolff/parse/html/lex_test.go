package html // import "github.com/tdewolff/parse/html"

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
		} else if tt == AttributeToken {
			s += tt.String() + "('" + string(data) + "=" + string(l.AttrVal()) + "') "
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
		html     string
		expected TTs
	}{
		{"<html></html>", TTs{StartTagToken, StartTagCloseToken, EndTagToken}},
		{"<img/>", TTs{StartTagToken, StartTagVoidToken}},
		{"<!-- comment -->", TTs{CommentToken}},
		{"<!-- comment --!>", TTs{CommentToken}},
		{"<p>text</p>", TTs{StartTagToken, StartTagCloseToken, TextToken, EndTagToken}},
		{"<input type='button'/>", TTs{StartTagToken, AttributeToken, StartTagVoidToken}},
		{"<input  type='button'  value=''/>", TTs{StartTagToken, AttributeToken, AttributeToken, StartTagVoidToken}},
		{"<input type='=/>' \r\n\t\f value=\"'\" name=x checked />", TTs{StartTagToken, AttributeToken, AttributeToken, AttributeToken, AttributeToken, StartTagVoidToken}},
		{"<!doctype>", TTs{DoctypeToken}},
		{"<!doctype html>", TTs{DoctypeToken}},
		{"<?bogus>", TTs{CommentToken}},
		{"</0bogus>", TTs{CommentToken}},
		{"<!bogus>", TTs{CommentToken}},
		{"< ", TTs{TextToken}},
		{"</", TTs{TextToken}},

		// raw tags
		{"<title><p></p></title>", TTs{StartTagToken, StartTagCloseToken, TextToken, EndTagToken}},
		{"<TITLE><p></p></TITLE>", TTs{StartTagToken, StartTagCloseToken, TextToken, EndTagToken}},
		{"<plaintext></plaintext>", TTs{StartTagToken, StartTagCloseToken, TextToken}},
		{"<script></script>", TTs{StartTagToken, StartTagCloseToken, EndTagToken}},
		{"<script>var x='</script>';</script>", TTs{StartTagToken, StartTagCloseToken, TextToken, EndTagToken, TextToken, EndTagToken}},
		{"<script><!--var x='</script>';--></script>", TTs{StartTagToken, StartTagCloseToken, TextToken, EndTagToken, TextToken, EndTagToken}},
		{"<script><!--var x='<script></script>';--></script>", TTs{StartTagToken, StartTagCloseToken, TextToken, EndTagToken}},
		{"<script><!--var x='<script>';--></script>", TTs{StartTagToken, StartTagCloseToken, TextToken, EndTagToken}},
		{"<![CDATA[ test ]]>", TTs{TextToken}},

		// early endings
		{"<!-- comment", TTs{CommentToken}},
		{"<? bogus comment", TTs{CommentToken}},
		{"<foo", TTs{StartTagToken}},
		{"</foo", TTs{EndTagToken}},
		{"<foo x", TTs{StartTagToken, AttributeToken}},
		{"<foo x=", TTs{StartTagToken, AttributeToken}},
		{"<foo x='", TTs{StartTagToken, AttributeToken}},
		{"<foo x=''", TTs{StartTagToken, AttributeToken}},
		{"<!DOCTYPE note SYSTEM", TTs{DoctypeToken}},
		{"<![CDATA[ test", TTs{TextToken}},
		{"<script>", TTs{StartTagToken, StartTagCloseToken}},
		{"<script><!--", TTs{StartTagToken, StartTagCloseToken, TextToken}},
		{"<script><!--var x='<script></script>';-->", TTs{StartTagToken, StartTagCloseToken, TextToken}},

		// go-fuzz
		{"</>", TTs{EndTagToken}},
	}
	for _, test := range tokenTests {
		expected := []TokenType(test.expected)
		stringify := helperStringify(t, test.html)
		l := NewLexer(bytes.NewBufferString(test.html))
		i := 0
		for {
			tt, _ := l.Next()
			if tt == ErrorToken {
				assert.Equal(t, io.EOF, l.Err(), "error must be EOF in "+stringify)
				assert.Equal(t, len(expected), i, "when error occurred we must be at the end in "+stringify)
				break
			}
			assert.False(t, i >= len(expected), "index must not exceed expected token types size in "+stringify)
			if i < len(expected) {
				assert.Equal(t, expected[i], tt, "token types must match at index "+strconv.Itoa(i)+" in "+stringify)
			}
			i++
		}
	}

	assert.Equal(t, "Invalid(100)", TokenType(100).String())
}

func TestTags(t *testing.T) {
	var tagTests = []struct {
		html     string
		expected string
	}{
		{"<foo:bar.qux-norf/>", "foo:bar.qux-norf"},
		{"<foo?bar/qux>", "foo?bar/qux"},
		{"<!DOCTYPE note SYSTEM \"Note.dtd\">", " note SYSTEM \"Note.dtd\""},
		{"</foo >", "foo"},

		// early endings
		{"<foo ", "foo"},
	}
	for _, test := range tagTests {
		stringify := helperStringify(t, test.html)
		l := NewLexer(bytes.NewBufferString(test.html))
		for {
			tt, _ := l.Next()
			if tt == ErrorToken {
				assert.Equal(t, io.EOF, l.Err(), "error must be EOF in "+stringify)
				assert.Fail(t, "when error occurred we must be at the end in "+stringify)
				break
			} else if tt == StartTagToken || tt == EndTagToken || tt == DoctypeToken {
				assert.Equal(t, test.expected, string(l.Text()), "tags must match in "+stringify)
				break
			}
		}
	}
}

func TestAttributes(t *testing.T) {
	var attributeTests = []struct {
		attr     string
		expected []string
	}{
		{"<foo a=\"b\" />", []string{"a", "\"b\""}},
		{"<foo \nchecked \r\n value\r=\t'=/>\"' />", []string{"checked", "", "value", "'=/>\"'"}},
		{"<foo bar=\" a \n\t\r b \" />", []string{"bar", "\" a \n\t\r b \""}},
		{"<foo a/>", []string{"a", ""}},
		{"<foo /=/>", []string{"/", "/"}},

		// early endings
		{"<foo x", []string{"x", ""}},
		{"<foo x=", []string{"x", ""}},
		{"<foo x='", []string{"x", "'"}},
	}
	for _, test := range attributeTests {
		stringify := helperStringify(t, test.attr)
		l := NewLexer(bytes.NewBufferString(test.attr))
		i := 0
		for {
			tt, _ := l.Next()
			if tt == ErrorToken {
				assert.Equal(t, io.EOF, l.Err(), "error must be EOF in "+stringify)
				assert.Equal(t, len(test.expected), i, "when error occurred we must be at the end in "+stringify)
				break
			} else if tt == AttributeToken {
				assert.False(t, i+1 >= len(test.expected), "index must not exceed expected attributes size in "+stringify)
				if i+1 < len(test.expected) {
					assert.Equal(t, test.expected[i], string(l.Text()), "attribute keys must match at index "+strconv.Itoa(i)+" in "+stringify)
					assert.Equal(t, test.expected[i+1], string(l.AttrVal()), "attribute values must match at index "+strconv.Itoa(i)+" in "+stringify)
					i += 2
				}
			}
		}
	}
}

////////////////////////////////////////////////////////////////

var J int
var ss = [][]byte{
	[]byte(" style"),
	[]byte("style"),
	[]byte(" \r\n\tstyle"),
	[]byte("      style"),
	[]byte(" x"),
	[]byte("x"),
}

func BenchmarkWhitespace1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, s := range ss {
			j := 0
			for {
				if c := s[j]; c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\f' {
					j++
				} else {
					break
				}
			}
			J += j
		}
	}
}

func BenchmarkWhitespace2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, s := range ss {
			j := 0
			for {
				if c := s[j]; c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\f' {
					j++
					continue
				}
				break
			}
			J += j
		}
	}
}

func BenchmarkWhitespace3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, s := range ss {
			j := 0
			for {
				if c := s[j]; c != ' ' && c != '\t' && c != '\n' && c != '\r' && c != '\f' {
					break
				}
				j++
			}
			J += j
		}
	}
}

////////////////////////////////////////////////////////////////

func ExampleNewLexer() {
	l := NewLexer(bytes.NewBufferString("<span class='user'>John Doe</span>"))
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
	// Output: <span class='user'>John Doe</span>
}

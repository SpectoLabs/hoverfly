package xml // import "github.com/tdewolff/parse/xml"

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
		xml      string
		expected TTs
	}{
		{"", TTs{}},
		{"<!-- comment -->", TTs{CommentToken}},
		{"<!-- comment \n multi \r line -->", TTs{CommentToken}},
		{"<foo/>", TTs{StartTagToken, StartTagCloseVoidToken}},
		{"<foo \t\r\n/>", TTs{StartTagToken, StartTagCloseVoidToken}},
		{"<foo:bar.qux-norf/>", TTs{StartTagToken, StartTagCloseVoidToken}},
		{"<foo></foo>", TTs{StartTagToken, StartTagCloseToken, EndTagToken}},
		{"<foo>text</foo>", TTs{StartTagToken, StartTagCloseToken, TextToken, EndTagToken}},
		{"<foo/> text", TTs{StartTagToken, StartTagCloseVoidToken, TextToken}},
		{"<a> <b> <c>text</c> </b> </a>", TTs{StartTagToken, StartTagCloseToken, TextToken, StartTagToken, StartTagCloseToken, TextToken, StartTagToken, StartTagCloseToken, TextToken, EndTagToken, TextToken, EndTagToken, TextToken, EndTagToken}},
		{"<foo a='a' b=\"b\" c=c/>", TTs{StartTagToken, AttributeToken, AttributeToken, AttributeToken, StartTagCloseVoidToken}},
		{"<foo a=\"\"/>", TTs{StartTagToken, AttributeToken, StartTagCloseVoidToken}},
		{"<foo a-b=\"\"/>", TTs{StartTagToken, AttributeToken, StartTagCloseVoidToken}},
		{"<foo \nchecked \r\n value\r=\t'=/>\"' />", TTs{StartTagToken, AttributeToken, AttributeToken, StartTagCloseVoidToken}},
		{"<?xml?>", TTs{StartTagPIToken, StartTagClosePIToken}},
		{"<?xml a=\"a\" ?>", TTs{StartTagPIToken, AttributeToken, StartTagClosePIToken}},
		{"<?xml a=a?>", TTs{StartTagPIToken, AttributeToken, StartTagClosePIToken}},
		{"<![CDATA[ test ]]>", TTs{CDATAToken}},
		{"<!DOCTYPE>", TTs{DOCTYPEToken}},
		{"<!DOCTYPE note SYSTEM \"Note.dtd\">", TTs{DOCTYPEToken}},
		{`<!DOCTYPE note [<!ENTITY nbsp "&#xA0;"><!ENTITY writer "Writer: Donald Duck."><!ENTITY copyright "Copyright:]> W3Schools.">]>`, TTs{DOCTYPEToken}},
		{"<!foo>", TTs{StartTagToken, StartTagCloseToken}},

		// early endings
		{"<!-- comment", TTs{CommentToken}},
		{"<foo", TTs{StartTagToken}},
		{"</foo", TTs{EndTagToken}},
		{"<foo x", TTs{StartTagToken, AttributeToken}},
		{"<foo x=", TTs{StartTagToken, AttributeToken}},
		{"<foo x='", TTs{StartTagToken, AttributeToken}},
		{"<foo x=''", TTs{StartTagToken, AttributeToken}},
		{"<?xml", TTs{StartTagPIToken}},
		{"<![CDATA[ test", TTs{CDATAToken}},
		{"<!DOCTYPE note SYSTEM", TTs{DOCTYPEToken}},

		// go fuzz
		{"</", TTs{EndTagToken}},
		{"</\n", TTs{EndTagToken}},
	}
	for _, test := range tokenTests {
		expected := []TokenType(test.expected)
		stringify := helperStringify(t, test.xml)
		l := NewLexer(bytes.NewBufferString(test.xml))
		i := 0
		for {
			tt, _ := l.Next()
			if tt == ErrorToken {
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
		xml      string
		expected string
	}{
		{"<foo:bar.qux-norf/>", "foo:bar.qux-norf"},
		{"<?xml?>", "xml"},
		{"<foo?bar/qux>", "foo?bar/qux"},
		{"<!DOCTYPE note SYSTEM \"Note.dtd\">", " note SYSTEM \"Note.dtd\""},

		// early endings
		{"<foo ", "foo"},
	}
	for _, test := range tagTests {
		stringify := helperStringify(t, test.xml)
		l := NewLexer(bytes.NewBufferString(test.xml))
		for {
			tt, _ := l.Next()
			if tt == ErrorToken {
				assert.Equal(t, io.EOF, l.Err(), "error must be EOF in "+stringify)
				assert.Fail(t, "when error occurred we must be at the end in "+stringify)
				break
			} else if tt == StartTagToken || tt == StartTagPIToken || tt == EndTagToken || tt == DOCTYPEToken {
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
		{"<foo bar=\" a \n\t\r b \" />", []string{"bar", "\" a     b \""}},
		{"<?xml a=b?>", []string{"a", "b"}},
		{"<foo /=? >", []string{"/", "?"}},

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

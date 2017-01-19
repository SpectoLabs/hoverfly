package xml // import "github.com/tdewolff/parse/xml"

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertAttrVal(t *testing.T, input, expected string) {
	s := []byte(input)
	if len(s) > 1 && (s[0] == '"' || s[0] == '\'') && s[0] == s[len(s)-1] {
		s = s[1 : len(s)-1]
	}
	buf := make([]byte, len(s))
	assert.Equal(t, expected, string(EscapeAttrVal(&buf, []byte(s))))
}

func assertCDATAVal(t *testing.T, input, expected string) {
	s := []byte(input[len("<![CDATA[") : len(input)-len("]]>")])
	var buf []byte
	data, useText := EscapeCDATAVal(&buf, s)
	text := string(data)
	if !useText {
		text = "<![CDATA[" + text + "]]>"
	}
	assert.Equal(t, expected, text)
}

////////////////////////////////////////////////////////////////

func TestAttrVal(t *testing.T) {
	assertAttrVal(t, "xyz", "\"xyz\"")
	assertAttrVal(t, "", "\"\"")
	assertAttrVal(t, "x&amp;z", "\"x&amp;z\"")
	assertAttrVal(t, "x'z", "\"x'z\"")
	assertAttrVal(t, "x\"z", "'x\"z'")
	assertAttrVal(t, "a'b=\"\"", "'a&#39;b=\"\"'")
	assertAttrVal(t, "'x&#39;\"&#39;z'", "\"x'&#34;'z\"")
	assertAttrVal(t, "\"x&#34;'&#34;z\"", "'x\"&#39;\"z'")
	assertAttrVal(t, "a&#39;b=\"\"", "'a&#39;b=\"\"'")
}

func TestCDATAVal(t *testing.T) {
	assertCDATAVal(t, "<![CDATA[<b>]]>", "&lt;b>")
	assertCDATAVal(t, "<![CDATA[abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz]]>", "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz")
	assertCDATAVal(t, "<![CDATA[ <b> ]]>", " &lt;b> ")
	assertCDATAVal(t, "<![CDATA[<<<<<]]>", "<![CDATA[<<<<<]]>")
	assertCDATAVal(t, "<![CDATA[&]]>", "&amp;")
	assertCDATAVal(t, "<![CDATA[&&&&]]>", "<![CDATA[&&&&]]>")
	assertCDATAVal(t, "<![CDATA[ a ]]>", " a ")
	assertCDATAVal(t, "<![CDATA[]]>", "")
}

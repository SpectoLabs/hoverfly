package html // import "github.com/tdewolff/parse/html"

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEscapeAttrVal(t *testing.T) {
	var escapeAttrValTests = []struct {
		attrVal  string
		expected string
	}{
		{"xyz", "xyz"},
		{"", ""},
		{"x&amp;z", "x&amp;z"},
		{"x/z", "x/z"},
		{"x'z", "\"x'z\""},
		{"x\"z", "'x\"z'"},
		{"'x\"z'", "'x\"z'"},
		{"'x&#39;\"&#39;z'", "\"x'&#34;'z\""},
		{"\"x&#34;'&#34;z\"", "'x\"&#39;\"z'"},
		{"\"x&#x27;z\"", "\"x'z\""},
		{"'x&#x00022;z'", "'x\"z'"},
		{"'x\"&gt;'", "'x\"&gt;'"},
		{"You&#039;re encouraged to log in; however, it&#039;s not mandatory. [o]", "\"You're encouraged to log in; however, it's not mandatory. [o]\""},
		{"a'b=\"\"", "'a&#39;b=\"\"'"},
		{"x<z", "\"x<z\""},
		{"'x\"&#39;\"z'", "'x\"&#39;\"z'"},
	}
	for _, tt := range escapeAttrValTests {
		s := []byte(tt.attrVal)
		orig := s
		if len(s) > 1 && (s[0] == '"' || s[0] == '\'') && s[0] == s[len(s)-1] {
			s = s[1 : len(s)-1]
		}
		buf := make([]byte, len(s))
		assert.Equal(t, tt.expected, string(EscapeAttrVal(&buf, orig, s)))
	}
}

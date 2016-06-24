package css // import "github.com/tdewolff/parse/css"

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tdewolff/test"
)

////////////////////////////////////////////////////////////////

func TestParse(t *testing.T) {
	var parseTests = []struct {
		inline   bool
		css      string
		expected string
	}{
		{true, " x : y ; ", "x:y;"},
		{true, "color: red;", "color:red;"},
		{true, "color : red;", "color:red;"},
		{true, "color: red; border: 0;", "color:red;border:0;"},
		{true, "color: red !important;", "color:red!important;"},
		{true, "color: red ! important;", "color:red!important;"},
		{true, "white-space: -moz-pre-wrap;", "white-space:-moz-pre-wrap;"},
		{true, "display: -moz-inline-stack;", "display:-moz-inline-stack;"},
		{true, "x: 10px / 1em;", "x:10px/1em;"},
		{true, "x: 1em/1.5em \"Times New Roman\", Times, serif;", "x:1em/1.5em \"Times New Roman\",Times,serif;"},
		{true, "x: hsla(100,50%, 75%, 0.5);", "x:hsla(100,50%,75%,0.5);"},
		{true, "x: hsl(100,50%, 75%);", "x:hsl(100,50%,75%);"},
		{true, "x: rgba(255, 238 , 221, 0.3);", "x:rgba(255,238,221,0.3);"},
		{true, "x: 50vmax;", "x:50vmax;"},
		{true, "color: linear-gradient(to right, black, white);", "color:linear-gradient(to right,black,white);"},
		{true, "color: calc(100%/2 - 1em);", "color:calc(100%/2 - 1em);"},
		{true, "color: calc(100%/2--1em);", "color:calc(100%/2--1em);"},
		{false, "<!-- @charset; -->", "<!--@charset;-->"},
		{false, "@media print, screen { }", "@media print,screen{}"},
		{false, "@media { @viewport ; }", "@media{@viewport;}"},
		{false, "@keyframes 'diagonal-slide' {  from { left: 0; top: 0; } to { left: 100px; top: 100px; } }", "@keyframes 'diagonal-slide'{from{left:0;top:0;}to{left:100px;top:100px;}}"},
		{false, "@keyframes movingbox{0%{left:90%;}50%{left:10%;}100%{left:90%;}}", "@keyframes movingbox{0%{left:90%;}50%{left:10%;}100%{left:90%;}}"},
		{false, ".foo { color: #fff;}", ".foo{color:#fff;}"},
		{false, ".foo { *color: #fff;}", ".foo{*color:#fff;}"},
		{false, ".foo { ; _color: #fff;}", ".foo{_color:#fff;}"},
		{false, "a { color: red; border: 0; }", "a{color:red;border:0;}"},
		{false, "a { color: red; border: 0; } b { padding: 0; }", "a{color:red;border:0;}b{padding:0;}"},

		// extraordinary
		{true, "color: red;;", "color:red;"},
		{true, "color:#c0c0c0", "color:#c0c0c0;"},
		{true, "background:URL(x.png);", "background:URL(x.png);"},
		{true, "filter: progid : DXImageTransform.Microsoft.BasicImage(rotation=1);", "filter:progid:DXImageTransform.Microsoft.BasicImage(rotation=1);"},
		{true, "/*a*/\n/*c*/\nkey: value;", "key:value;"},
		{true, "@-moz-charset;", "@-moz-charset;"},
		{false, "@import;@import;", "@import;@import;"},
		{false, ".a .b#c, .d<.e { x:y; }", ".a .b#c,.d<.e{x:y;}"},
		{false, ".a[b~=c]d { x:y; }", ".a[b~=c]d{x:y;}"},
		// {false, "{x:y;}", "{x:y;}"},
		{false, "a{}", "a{}"},
		{false, "a,.b/*comment*/ {x:y;}", "a,.b{x:y;}"},
		{false, "a,.b/*comment*/.c {x:y;}", "a,.b.c{x:y;}"},
		{false, "a{x:; z:q;}", "a{x:;z:q;}"},
		{false, "@font-face { x:y; }", "@font-face{x:y;}"},
		{false, "a:not([controls]){x:y;}", "a:not([controls]){x:y;}"},
		{false, "@document regexp('https:.*') { p { color: red; } }", "@document regexp('https:.*'){p{color:red;}}"},
		{false, "@media all and ( max-width:400px ) { }", "@media all and (max-width:400px){}"},
		{false, "@media (max-width:400px) { }", "@media(max-width:400px){}"},
		{false, "@media (max-width:400px)", "@media(max-width:400px);"},
		{false, "@font-face { ; font:x; }", "@font-face{font:x;}"},
		{false, "@-moz-font-face { ; font:x; }", "@-moz-font-face{font:x;}"},
		{false, "@unknown abc { {} lala }", "@unknown abc{{}lala}"},
		{false, "a[x={}]{x:y;}", "a[x={}]{x:y;}"},
		{false, "a[x=,]{x:y;}", "a[x=,]{x:y;}"},
		{false, "a[x=+]{x:y;}", "a[x=+]{x:y;}"},
		{false, ".cla .ss > #id { x:y; }", ".cla .ss>#id{x:y;}"},
		{false, ".cla /*a*/ /*b*/ .ss{}", ".cla .ss{}"},
		{false, "a{x:f(a(),b);}", "a{x:f(a(),b);}"},
		{false, "a{x:y!z;}", "a{x:y!z;}"},
		{false, "[class*=\"column\"]+[class*=\"column\"]:last-child{a:b;}", "[class*=\"column\"]+[class*=\"column\"]:last-child{a:b;}"},
		{false, "@media { @viewport }", "@media{@viewport;}"},
		{false, "table { @unknown }", "table{@unknown;}"},

		// early endings
		{true, "~color:red;", ""},
		{false, "selector{", "selector{"},
		{false, "@media{selector{", "@media{selector{"},

		// issues
		{false, "@media print {.class{width:5px;}}", "@media print{.class{width:5px;}}"},                  // #6
		{false, ".class{width:calc((50% + 2em)/2 + 14px);}", ".class{width:calc((50% + 2em)/2 + 14px);}"}, // #7
		{false, ".class [c=y]{}", ".class [c=y]{}"},                                                       // tdewolff/minify#16
		{false, "table{font-family:Verdana}", "table{font-family:Verdana;}"},                              // tdewolff/minify#22

		// go-fuzz
		{false, "@-webkit-", "@-webkit-;"},
	}
	for _, test := range parseTests {
		output := ""
		p := NewParser(bytes.NewBufferString(test.css), test.inline)
		for {
			gt, _, data := p.Next()
			if gt == ErrorGrammar {
				err := p.Err()
				if err != nil {
					assert.Equal(t, io.EOF, err, "parser must not return error '"+err.Error()+"' in "+test.css)
				}
				break
			} else if gt == AtRuleGrammar || gt == BeginAtRuleGrammar || gt == BeginRulesetGrammar || gt == DeclarationGrammar {
				if gt == DeclarationGrammar {
					data = append(data, ":"...)
				}
				for _, val := range p.Values() {
					data = append(data, val.Data...)
				}
				if gt == BeginAtRuleGrammar || gt == BeginRulesetGrammar {
					data = append(data, "{"...)
				} else if gt == AtRuleGrammar || gt == DeclarationGrammar {
					data = append(data, ";"...)
				}
			}
			output += string(data)
		}
		assert.Equal(t, test.expected, output, "parsed string must match expected result in "+test.css)
	}

	assert.Equal(t, "Error", ErrorGrammar.String())
	assert.Equal(t, "AtRule", AtRuleGrammar.String())
	assert.Equal(t, "BeginAtRule", BeginAtRuleGrammar.String())
	assert.Equal(t, "EndAtRule", EndAtRuleGrammar.String())
	assert.Equal(t, "BeginRuleset", BeginRulesetGrammar.String())
	assert.Equal(t, "EndRuleset", EndRulesetGrammar.String())
	assert.Equal(t, "Declaration", DeclarationGrammar.String())
	assert.Equal(t, "Token", TokenGrammar.String())
	assert.Equal(t, "Invalid(100)", GrammarType(100).String())
}

func TestParseError(t *testing.T) {
	var parseErrorTests = []struct {
		inline   bool
		css      string
		expected error
	}{
		{false, "selector", ErrBadQualifiedRule},
		{true, "color 0", ErrBadDeclaration},
	}
	for _, test := range parseErrorTests {
		p := NewParser(bytes.NewBufferString(test.css), test.inline)
		for {
			gt, _, _ := p.Next()
			if gt == ErrorGrammar {
				assert.Equal(t, test.expected, p.Err(), "parser must return error '"+test.expected.Error()+"' in "+test.css)
				break
			}
		}
	}
}

func TestReader(t *testing.T) {
	input := "x:a;"
	p := NewParser(test.NewPlainReader(bytes.NewBufferString(input)), true)
	for {
		gt, _, _ := p.Next()
		if gt == ErrorGrammar {
			break
		}
	}
}

////////////////////////////////////////////////////////////////

func ExampleNewParser() {
	p := NewParser(bytes.NewBufferString("color: red;"), true) // false because this is the content of an inline style attribute
	out := ""
	for {
		gt, _, data := p.Next()
		if gt == ErrorGrammar {
			break
		} else if gt == AtRuleGrammar || gt == BeginAtRuleGrammar || gt == BeginRulesetGrammar || gt == DeclarationGrammar {
			out += string(data)
			if gt == DeclarationGrammar {
				out += ":"
			}
			for _, val := range p.Values() {
				out += string(val.Data)
			}
			if gt == BeginAtRuleGrammar || gt == BeginRulesetGrammar {
				out += "{"
			} else if gt == AtRuleGrammar || gt == DeclarationGrammar {
				out += ";"
			}
		} else {
			out += string(data)
		}
	}
	fmt.Println(out)
	// Output: color:red;
}

package goxpath

import (
	"fmt"
	"testing"
)

func TestStrLit(t *testing.T) {
	p := `'strlit'`
	x := `<?xml version="1.0" encoding="UTF-8"?><test><path/></test>`
	exp := "strlit"
	execVal(p, x, exp, nil, t)
}

func TestNumLit(t *testing.T) {
	p := `123`
	x := `<?xml version="1.0" encoding="UTF-8"?><test><path/></test>`
	exp := "123"
	execVal(p, x, exp, nil, t)
	p = `123.456`
	exp = "123.456"
	execVal(p, x, exp, nil, t)
}

func TestLast(t *testing.T) {
	p := `/p1/*[last()]`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1><p2/><p3/><p4/></p1>`
	exp := []string{"<p4></p4>"}
	execPath(p, x, exp, nil, t)
	p = `/p1/p5[last()]`
	exp = []string{}
	execPath(p, x, exp, nil, t)
	p = `/p1[last()]`
	exp = []string{"<p1><p2></p2><p3></p3><p4></p4></p1>"}
	execPath(p, x, exp, nil, t)
}

func TestCount(t *testing.T) {
	p := `count(//p1)`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1><p2><?test?></p2><p3/><p1/></p1>`
	exp := "2"
	execVal(p, x, exp, nil, t)
}

func TestCount2(t *testing.T) {
	x := `<?xml version="1.0"?>
<test>
    <x a="1">
      <x a="2">
        <x>
          <y>y31</y>
          <y>y32</y>
        </x>
      </x>
    </x>
    <x a="1">
      <x a="2">
        <y>y21</y>
        <y>y22</y>
      </x>
    </x>
    <x a="1">
      <y>y11</y>
      <y>y12</y>
    </x>
    <x>
      <y>y03</y>
      <y>y04</y>
    </x>
</test>
`
	execVal(`count(//x)`, x, "7", nil, t)
	execVal(`count(//x[1])`, x, "4", nil, t)
	execVal(`count(//x/y)`, x, "8", nil, t)
	execVal(`count(//x/y[1])`, x, "4", nil, t)
	execVal(`count(//x[1]/y[1])`, x, "2", nil, t)
}

func TestNames(t *testing.T) {
	x := `<?xml version="1.0" encoding="UTF-8"?><test xmlns="http://foo.com" xmlns:bar="http://bar.com" bar:attr="val"><?pi pival?><!--comment--></test>`
	testMap := make(map[string]map[string]string)
	testMap["/*"] = make(map[string]string)
	testMap["/*"]["local-name"] = "test"
	testMap["/*"]["namespace-uri"] = "http://foo.com"
	testMap["/*"]["name"] = "{http://foo.com}test"

	testMap["/none"] = make(map[string]string)
	testMap["/none"]["local-name"] = ""
	testMap["/none"]["namespace-uri"] = ""
	testMap["/none"]["name"] = ""

	testMap["/*/@*:attr"] = make(map[string]string)
	testMap["/*/@*:attr"]["local-name"] = "attr"
	testMap["/*/@*:attr"]["namespace-uri"] = "http://bar.com"
	testMap["/*/@*:attr"]["name"] = "{http://bar.com}attr"

	testMap["//processing-instruction()"] = make(map[string]string)
	testMap["//processing-instruction()"]["local-name"] = "pi"
	testMap["//processing-instruction()"]["namespace-uri"] = ""
	testMap["//processing-instruction()"]["name"] = "pi"

	testMap["//comment()"] = make(map[string]string)
	testMap["//comment()"]["local-name"] = ""
	testMap["//comment()"]["namespace-uri"] = ""
	testMap["//comment()"]["name"] = ""

	for path, i := range testMap {
		for nt, res := range i {
			p := fmt.Sprintf("%s(%s)", nt, path)
			exp := res
			execVal(p, x, exp, nil, t)
		}
	}

	x = `<?xml version="1.0" encoding="UTF-8"?><test xmlns="http://foo.com" />`
	execPath("/*[local-name() = 'test']", x, []string{`<test xmlns="http://foo.com"></test>`}, nil, t)
	execPath("/*[namespace-uri() = 'http://foo.com']", x, []string{`<test xmlns="http://foo.com"></test>`}, nil, t)
	execPath("/*[name() = '{http://foo.com}test']", x, []string{`<test xmlns="http://foo.com"></test>`}, nil, t)
}

func TestBoolean(t *testing.T) {
	x := `<?xml version="1.0" encoding="UTF-8"?><p1><p2/><p3/><p4/></p1>`
	execVal(`true()`, x, "true", nil, t)
	execVal(`false()`, x, "false", nil, t)
	p := `boolean(/p1/p2)`
	exp := "true"
	execVal(p, x, exp, nil, t)
	p = `boolean(/p1/p5)`
	exp = "false"
	execVal(p, x, exp, nil, t)
	p = `boolean('123')`
	exp = "true"
	execVal(p, x, exp, nil, t)
	p = `boolean(123)`
	exp = "true"
	execVal(p, x, exp, nil, t)
	p = `boolean('')`
	exp = "false"
	execVal(p, x, exp, nil, t)
	p = `boolean(0)`
	exp = "false"
	execVal(p, x, exp, nil, t)
}

func TestNot(t *testing.T) {
	x := `<?xml version="1.0" encoding="UTF-8"?><p1><p2/><p3/><p4/></p1>`
	execVal(`not(false())`, x, "true", nil, t)
	execVal(`not(true())`, x, "false", nil, t)
}

func TestConversions(t *testing.T) {
	x := `<?xml version="1.0" encoding="UTF-8"?><p1>foo</p1>`
	execVal(`number(true())`, x, "1", nil, t)
	execVal(`number(false())`, x, "0", nil, t)
	execVal(`string(/p2)`, x, "", nil, t)
	execVal(`/p1[string() = 'foo']`, x, "foo", nil, t)
	execVal(`number('abc')`, x, "NaN", nil, t)
}

func TestLang(t *testing.T) {
	x := `<?xml version="1.0" encoding="UTF-8"?>
<p1>
	<p xml:lang="en">I went up a floor.</p>
	<p xml:lang="en-GB">I took the lift.</p>
	<p xml:lang="en-US">I rode the elevator.</p>
</p1>`
	execVal(`count(//p[lang('en')])`, x, "3", nil, t)
	execVal(`count(//text()[lang('en-GB')])`, x, "1", nil, t)
	execVal(`count(//p[lang('en-US')])`, x, "1", nil, t)
	execVal(`count(//p[lang('de')])`, x, "0", nil, t)
	execVal(`count(/p1[lang('en')])`, x, "0", nil, t)
}

func TestString(t *testing.T) {
	x := `<?xml version="1.0" encoding="UTF-8"?><p1>text</p1>`
	execVal(`string(2 + 2)`, x, "4", nil, t)
	execVal(`string(/p1)`, x, "text", nil, t)
}

func TestConcat(t *testing.T) {
	x := `<?xml version="1.0" encoding="UTF-8"?><p1/>`
	execVal(`concat('abc', 'def', 'hij', '123')`, x, "abcdefhij123", nil, t)
}

func TestStartsWith(t *testing.T) {
	x := `<?xml version="1.0" encoding="UTF-8"?><p1/>`
	execVal(`starts-with('abcd', 'ab')`, x, "true", nil, t)
	execVal(`starts-with('abcd', 'abd')`, x, "false", nil, t)
}

func TestContains(t *testing.T) {
	x := `<?xml version="1.0" encoding="UTF-8"?><p1/>`
	execVal(`contains('abcd', 'bcd')`, x, "true", nil, t)
	execVal(`contains('abcd', 'bd')`, x, "false", nil, t)
}

func TestSubstrBefore(t *testing.T) {
	x := `<?xml version="1.0" encoding="UTF-8"?><p1/>`
	execVal(`substring-before("1999/04/01","/")`, x, "1999", nil, t)
	execVal(`substring-before("1999/04/01","2")`, x, "", nil, t)
}

func TestSubstrAfter(t *testing.T) {
	x := `<?xml version="1.0" encoding="UTF-8"?><p1/>`
	execVal(`substring-after("1999/04/01","/")`, x, "04/01", nil, t)
	execVal(`substring-after("1999/04/01","19")`, x, "99/04/01", nil, t)
	execVal(`substring-after("1999/04/01","a")`, x, "", nil, t)
}

func TestSubstring(t *testing.T) {
	x := `<?xml version="1.0" encoding="UTF-8"?><p1/>`
	execVal(`substring("12345", 2, 3)`, x, "234", nil, t)
	execVal(`substring("12345", 2)`, x, "2345", nil, t)
	execVal(`substring('abcd', -2, 5)`, x, "ab", nil, t)
	execVal(`substring('abcd', 0)`, x, "abcd", nil, t)
	execVal(`substring('abcd', 1, 4)`, x, "abcd", nil, t)
	execVal(`substring("12345", 1.5, 2.6)`, x, "234", nil, t)
	execVal(`substring("12345", 0 div 0, 3)`, x, "", nil, t)
	execVal(`substring("12345", 1, 0 div 0)`, x, "", nil, t)
	execVal(`substring("12345", -42, 1 div 0)`, x, "12345", nil, t)
	execVal(`substring("12345", -1 div 0, 1 div 0)`, x, "", nil, t)
}

func TestStrLength(t *testing.T) {
	x := `<?xml version="1.0" encoding="UTF-8"?><p1>abc</p1>`
	execVal(`string-length('def')`, x, "3", nil, t)
	execVal(`/p1[string-length() = 3]`, x, "abc", nil, t)
}

func TestNormalizeSpace(t *testing.T) {
	x := `<?xml version="1.0" encoding="UTF-8"?><p1>
	  a   b
</p1>`
	execVal(`normalize-space(/p1)`, x, "a b", nil, t)
	execVal(`/p1[normalize-space(/p1) = 'a b']`, x, `
	  a   b
`, nil, t)
	execVal(`/p1[normalize-space() = 'a b']`, x, `
	  a   b
`, nil, t)
}

func TestTranslate(t *testing.T) {
	x := `<?xml version="1.0" encoding="UTF-8"?><p1/>`
	execVal(`translate("bar","abc","ABC")`, x, "BAr", nil, t)
	execVal(`translate("--aaa--","abc-","ABC")`, x, "AAA", nil, t)
}

package goxpath

import (
	"bytes"
	"runtime/debug"
	"testing"

	"github.com/ChrisTrenkamp/goxpath/tree/xmltree"
)

func execVal(xp, x string, exp string, ns map[string]string, t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Error("Panicked: from XPath expr: '" + xp)
			t.Error(r)
			t.Error(string(debug.Stack()))
		}
	}()
	res := MustParse(xp).MustExec(xmltree.MustParseXML(bytes.NewBufferString(x)), func(o *Opts) { o.NS = ns })

	if res.String() != exp {
		t.Error("Incorrect result:'" + res.String() + "' from XPath expr: '" + xp + "'.  Expecting: '" + exp + "'")
		return
	}
}

func TestNodeVal(t *testing.T) {
	p := `/test`
	x := `<?xml version="1.0" encoding="UTF-8"?><test>test<path>path</path>test2</test>`
	exp := "testpathtest2"
	execVal(p, x, exp, nil, t)
}

func TestAttrVal(t *testing.T) {
	p := `/p1/@test`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1 test="foo" foo="test"><p2/></p1>`
	exp := "foo"
	execVal(p, x, exp, nil, t)
}

func TestCommentVal(t *testing.T) {
	p := `//comment()`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1><!-- comment --></p1>`
	exp := ` comment `
	execVal(p, x, exp, nil, t)
}

func TestProcInstVal(t *testing.T) {
	p := `//processing-instruction()`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1><?proc test?></p1>`
	exp := `test`
	execVal(p, x, exp, nil, t)
}

func TestNodeNamespaceVal(t *testing.T) {
	p := `/test:p1/namespace::test`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1 xmlns:test="http://test"/>`
	exp := `http://test`
	execVal(p, x, exp, nil, t)
}

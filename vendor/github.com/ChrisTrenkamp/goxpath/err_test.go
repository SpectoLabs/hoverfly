package goxpath

import (
	"bytes"
	"encoding/xml"
	"runtime/debug"
	"testing"

	"github.com/ChrisTrenkamp/goxpath/tree"
	"github.com/ChrisTrenkamp/goxpath/tree/xmltree"
	"github.com/ChrisTrenkamp/goxpath/tree/xmltree/xmlele"
)

type dummyType string

func (d dummyType) String() string {
	return string(d)
}
func dummyFunc(c tree.Ctx, args ...tree.Result) (tree.Result, error) {
	return dummyType(""), nil
}

var custFns = map[xml.Name]tree.Wrap{
	{Local: "dummy"}:                               {Fn: dummyFunc},
	{Space: "http://foo.com", Local: "spaceDummy"}: {Fn: dummyFunc},
}

func execErr(xp, x string, errStr string, ns map[string]string, t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Error("Panicked: from XPath expr: '" + xp)
			t.Error(r)
			t.Error(string(debug.Stack()))
		}
	}()
	_, err := ParseExec(xp, xmltree.MustParseXML(bytes.NewBufferString(x)), func(o *Opts) { o.NS = ns; o.Funcs = custFns })

	if err.Error() != errStr {
		t.Error("Incorrect result:'" + err.Error() + "' from XPath expr: '" + xp + "'.  Expecting: '" + errStr + "'")
		return
	}
}

func TestBadAxis(t *testing.T) {
	x := `<?xml version="1.0" encoding="UTF-8"?><p1/>`
	execErr(`/test/chil::p2`, x, "Invalid Axis specifier, chil", nil, t)
}

func TestIncompleteStep(t *testing.T) {
	x := `<?xml version="1.0" encoding="UTF-8"?><p1/>`
	execErr(`/child::+2`, x, "Step is not complete", nil, t)
	execErr(`/foo:`, x, "Step is not complete", nil, t)
}

func TestParseErr(t *testing.T) {
	_, err := xmltree.ParseXML(bytes.NewBufferString("<p1/>"))
	if err.Error() != "Malformed XML file" {
		t.Error("Incorrect error message:", err.Error())
	}

	_, err = xmltree.ParseXML(bytes.NewBufferString(""))
	if err.Error() != "EOF" {
		t.Error("Incorrect error message:", err.Error())
	}

	_, err = xmltree.ParseXML(bytes.NewBufferString("<p1/>"), func(s *xmltree.ParseOptions) {
		s.Strict = false
	})
	if err != nil {
		t.Error("Error not nil:", err.Error())
	}
}

func TestBadNodeType(t *testing.T) {
	x := `<?xml version="1.0" encoding="UTF-8"?><p1/>`
	execErr(`/test/foo()`, x, "Invalid node-type foo", nil, t)
}

func TestXPathErr(t *testing.T) {
	x := `<?xml version="1.0" encoding="UTF-8"?><p1/>`
	execErr(`/test/chil::p2`, x, "Invalid Axis specifier, chil", nil, t)
}

func TestNodeSetConvErr(t *testing.T) {
	x := `<?xml version="1.0" encoding="UTF-8"?><p1/>`
	for _, i := range []string{"sum", "count", "local-name", "namespace-uri", "name"} {
		execErr("/p1["+i+"(1)]", x, "Cannot convert object to a node-set", nil, t)
	}
}

func TestNodeSetConvUnionErr(t *testing.T) {
	x := `<?xml version="1.0" encoding="UTF-8"?><p1/>`
	execErr(`/p1 | 'invalid'`, x, "Cannot convert data type to node-set", nil, t)
}

func TestUnknownFunction(t *testing.T) {
	x := `<?xml version="1.0" encoding="UTF-8"?><p1/>`
	execErr(`invFunc()`, x, "Unknown function: invFunc", nil, t)
}

func TestUnterminatedString(t *testing.T) {
	x := `<?xml version="1.0" encoding="UTF-8"?><p1/>`
	execErr(`"asdf`, x, "Unexpected end of string literal.", nil, t)
}

func TestUnterminatedParenths(t *testing.T) {
	x := `<?xml version="1.0" encoding="UTF-8"?><p1/>`
	execErr(`(1 + 2 * 3`, x, "Missing end )", nil, t)
}

func TestUnterminatedNTQuotes(t *testing.T) {
	x := `<?xml version="1.0" encoding="UTF-8"?><p1><?foo bar?></p1>`
	execErr(`//processing-instruction('foo)`, x, "Unexpected end of string literal.", nil, t)
}

func TestUnterminatedNTParenths(t *testing.T) {
	x := `<?xml version="1.0" encoding="UTF-8"?><p1><?foo bar?></p1>`
	execErr(`//processing-instruction('foo'`, x, "Missing ) at end of NodeType declaration.", nil, t)
}

func TestUnterminatedFnParenths(t *testing.T) {
	x := `<?xml version="1.0" encoding="UTF-8"?><p1/>`
	execErr(`true(`, x, "Missing ) at end of function declaration.", nil, t)
}

func TestEmptyPred(t *testing.T) {
	x := `<?xml version="1.0" encoding="UTF-8"?><p1>text</p1>`
	execErr(`/p1[ ]`, x, "Missing content in predicate.", nil, t)
}

func TestUnterminatedPred(t *testing.T) {
	x := `<?xml version="1.0" encoding="UTF-8"?><p1>text</p1>`
	execErr(`/p1[. = 'text'`, x, "Missing ] at end of predicate.", nil, t)
}

func TestNotEnoughArgs(t *testing.T) {
	x := `<?xml version="1.0" encoding="UTF-8"?><p1>text</p1>`
	execErr(`concat('test')`, x, "Invalid number of arguments", nil, t)
}

func TestMarshalErr(t *testing.T) {
	x := `<?xml version="1.0" encoding="UTF-8"?><p1><p2/></p1>`
	n := xmltree.MustParseXML(bytes.NewBufferString(x))
	f := tree.FindNodeByPos(n, 3).(*xmlele.XMLEle)
	f.Name.Local = ""
	buf := &bytes.Buffer{}
	err := Marshal(n, buf)
	if err == nil {
		t.Error("No error")
	}
}

func TestParsePanic(t *testing.T) {
	errs := 0
	defer func() {
		if errs != 1 {
			t.Error("Err not 1")
		}
	}()
	defer func() {
		if r := recover(); r != nil {
			errs++
		}
	}()
	MustParse(`/foo()`)
}

func TestExecPanic(t *testing.T) {
	errs := 0
	defer func() {
		if errs != 1 {
			t.Error("Err not 1")
		}
	}()
	defer func() {
		if r := recover(); r != nil {
			errs++
		}
	}()
	MustParse("foo()").MustExec(xmltree.MustParseXML(bytes.NewBufferString(xml.Header + "<root/>")))
}

func TestParseXMLPanic(t *testing.T) {
	errs := 0
	defer func() {
		if errs != 1 {
			t.Error("Err not 1")
		}
	}()
	defer func() {
		if r := recover(); r != nil {
			errs++
		}
	}()
	xmltree.MustParseXML(bytes.NewBufferString("<root/>"))
}

func TestDummyType(t *testing.T) {
	ns := map[string]string{"foo": "http://foo.com"}
	x := `<?xml version="1.0" encoding="UTF-8"?><p1><p2/></p1>`
	execErr(`dummy() = 1`, x, "Cannot convert data type to number", nil, t)
	execErr(`dummy() = true()`, x, "Cannot convert argument to boolean", nil, t)
	execErr(`dummy() and true()`, x, "Cannot convert argument to boolean", nil, t)
	execErr(`not(dummy()) = 1`, x, "Cannot convert object to a boolean", nil, t)
	execErr(`1 = not(dummy())`, x, "Cannot convert object to a boolean", nil, t)
	execErr(`not(dummy() = 1)`, x, "Cannot convert data type to number", nil, t)
	execErr(`/p1[dummy()]`, x, "Cannot convert argument to boolean", nil, t)
	for _, i := range []string{"boolean", "not"} {
		execErr(i+`(dummy())`, x, "Cannot convert object to a boolean", nil, t)
	}
	for _, i := range []string{"number", "floor", "ceiling", "round"} {
		execErr(i+`(dummy())`, x, "Cannot convert object to a number", nil, t)
	}
	execErr(`substring("12345", dummy(), 2)`, x, "Cannot convert object to a number", nil, t)
	execErr(`substring("12345", 2, dummy())`, x, "Cannot convert object to a number", nil, t)
	execErr(`foo:spaceDummy() = 1`, x, "Cannot convert data type to number", ns, t)
}

func TestGoxpathBool(t *testing.T) {
	opts := func(o *Opts) { o.Funcs = custFns }
	x := xmltree.MustParseXML(bytes.NewBufferString(`<?xml version="1.0" encoding="UTF-8"?><p1><p2/></p1>`))
	_, err := MustParse(`dummy() = 1`).ExecBool(x, opts)
	if err == nil {
		t.Error("Error not nil")
	}
	_, err = MustParse(`dummy()`).ExecBool(x, opts)
	if err == nil {
		t.Error("Error not nil")
	}
	b, err := MustParse(`/p1`).ExecBool(x, opts)
	if !b || err != nil {
		t.Error("Incorrect result")
	}
}

func TestGoxpathNum(t *testing.T) {
	opts := func(o *Opts) { o.Funcs = custFns }
	x := xmltree.MustParseXML(bytes.NewBufferString(`<?xml version="1.0" encoding="UTF-8"?><p1><p2/>3</p1>`))
	_, err := MustParse(`dummy() = 1`).ExecNum(x, opts)
	if err == nil {
		t.Error("Error not nil")
	}
	_, err = MustParse(`dummy()`).ExecNum(x, opts)
	if err == nil {
		t.Error("Error not nil")
	}
	n, err := MustParse(`/p1`).ExecNum(x, opts)
	if n != 3 || err != nil {
		t.Error("Incorrect result")
	}
}

func TestGoxpathNode(t *testing.T) {
	opts := func(o *Opts) { o.Funcs = custFns }
	x := xmltree.MustParseXML(bytes.NewBufferString(`<?xml version="1.0" encoding="UTF-8"?><p1><p2/></p1>`))
	_, err := MustParse(`dummy() = 1`).ExecNode(x, opts)
	if err == nil {
		t.Error("Error not nil")
	}
	_, err = MustParse(`dummy()`).ExecNode(x, opts)
	if err == nil {
		t.Error("Error not nil")
	}
	n, err := MustParse(`/p1`).ExecNode(x, opts)
	if len(n) != 1 || err != nil {
		t.Error("Incorrect result")
	}
}

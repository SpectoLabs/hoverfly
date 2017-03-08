package goxpath

import (
	"bytes"
	"runtime/debug"
	"testing"

	"github.com/ChrisTrenkamp/goxpath/tree"
	"github.com/ChrisTrenkamp/goxpath/tree/xmltree"
)

func execPath(xp, x string, exp []string, ns map[string]string, t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Error("Panicked: from XPath expr: '" + xp)
			t.Error(r)
			t.Error(string(debug.Stack()))
		}
	}()
	res := MustParse(xp).MustExec(xmltree.MustParseXML(bytes.NewBufferString(x)), func(o *Opts) { o.NS = ns }).(tree.NodeSet)

	if len(res) != len(exp) {
		t.Error("Result length not valid in XPath expression '"+xp+"':", len(res), ", expecting", len(exp))
		for i := range res {
			t.Error(MarshalStr(res[i].(tree.Node)))
		}
		return
	}

	for i := range res {
		r, err := MarshalStr(res[i].(tree.Node))
		if err != nil {
			t.Error(err.Error())
			return
		}
		valid := false
		for j := range exp {
			if r == exp[j] {
				valid = true
			}
		}
		if !valid {
			t.Error("Incorrect result in XPath expression '" + xp + "':" + r)
			t.Error("Expecting one of:")
			for j := range exp {
				t.Error(exp[j])
			}
			return
		}
	}
}

func TestRoot(t *testing.T) {
	p := `/`
	x := `<?xml version="1.0" encoding="UTF-8"?><test><path/></test>`
	exp := []string{"<test><path></path></test>"}
	execPath(p, x, exp, nil, t)
}

func TestRoot2(t *testing.T) {
	x := xmltree.MustParseXML(bytes.NewBufferString(`<?xml version="1.0" encoding="UTF-8"?><test><path/></test>`))
	x = MustParse("/test/path").MustExec(x).(tree.NodeSet)[0]
	for _, i := range []string{"/", "//test"} {
		res, err := MarshalStr(MustParse(i).MustExec(x).(tree.NodeSet)[0])
		if err != nil {
			t.Error(err)
		}
		if res != "<test><path></path></test>" {
			t.Error("Incorrect result", res)
		}
	}
}

func TestAbsPath(t *testing.T) {
	p := ` / test / path`
	x := `<?xml version="1.0" encoding="UTF-8"?><test><path/></test>`
	exp := []string{"<path></path>"}
	execPath(p, x, exp, nil, t)
}

func TestNonAbsPath(t *testing.T) {
	p := ` test `
	x := `<?xml version="1.0" encoding="UTF-8"?><test><path/></test>`
	exp := []string{"<test><path></path></test>"}
	execPath(p, x, exp, nil, t)
}

func TestRelPath(t *testing.T) {
	p := ` // path `
	x := `<?xml version="1.0" encoding="UTF-8"?><test><path/></test>`
	exp := []string{"<path></path>"}
	execPath(p, x, exp, nil, t)
}

func TestRelPath2(t *testing.T) {
	p := ` // path `
	x := `<?xml version="1.0" encoding="UTF-8"?><path><test><path>test</path></test></path>`
	exp := []string{"<path><test><path>test</path></test></path>", `<path>test</path>`}
	execPath(p, x, exp, nil, t)
}

func TestRelNonRootPath(t *testing.T) {
	p := ` / test // path `
	x := `<?xml version="1.0" encoding="UTF-8"?><test><p1><p2><path/></p2></p1></test>`
	exp := []string{"<path></path>"}
	execPath(p, x, exp, nil, t)
}

func TestParent(t *testing.T) {
	p := `/test/path/ parent:: test`
	x := `<?xml version="1.0" encoding="UTF-8"?><test><path/></test>`
	exp := []string{"<test><path></path></test>"}
	execPath(p, x, exp, nil, t)
}

func TestAncestor(t *testing.T) {
	p := `/p1/p2/p3/p1/ancestor::p1`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1><p2><p3><p1></p1></p3></p2></p1>`
	exp := []string{`<p1><p2><p3><p1></p1></p3></p2></p1>`}
	execPath(p, x, exp, nil, t)
}

func TestAncestorOrSelf(t *testing.T) {
	p := `/p1/p2/p3/p1/ancestor-or-self::p1`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1><p2><p3><p1></p1></p3></p2></p1>`
	exp := []string{`<p1></p1>`, `<p1><p2><p3><p1></p1></p3></p2></p1>`}
	execPath(p, x, exp, nil, t)
}

func TestDescendant(t *testing.T) {
	p := `/p1/descendant::p1`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1><p2><p1/></p2></p1>`
	exp := []string{`<p1></p1>`}
	execPath(p, x, exp, nil, t)
}

func TestDescendantOrSelf(t *testing.T) {
	p := `/p1/descendant-or-self::p1`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1><p2><p1/></p2></p1>`
	exp := []string{`<p1><p2><p1></p1></p2></p1>`, `<p1></p1>`}
	execPath(p, x, exp, nil, t)
}

func TestAttribute(t *testing.T) {
	p := `/p1/attribute::test`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1 test="foo"></p1>`
	exp := []string{`<?attribute test="foo"?>`}
	execPath(p, x, exp, nil, t)
}

func TestAttributeSelf(t *testing.T) {
	p := `/p1/attribute::test/self::node()`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1 test="foo"></p1>`
	exp := []string{`<?attribute test="foo"?>`}
	execPath(p, x, exp, nil, t)
}

func TestAttributeAbbr(t *testing.T) {
	p := `/p1/ @ test`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1 test="foo"></p1>`
	exp := []string{`<?attribute test="foo"?>`}
	execPath(p, x, exp, nil, t)
}

func TestAttributeNoChild(t *testing.T) {
	p := `/p1/attribute::test/p2`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1 test="foo"><p2/></p1>`
	exp := []string{}
	execPath(p, x, exp, nil, t)
}

func TestAttributeParent(t *testing.T) {
	p := `/p1/attribute::test/ ..`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1 test="foo"><p2/></p1>`
	exp := []string{`<p1 test="foo"><p2></p2></p1>`}
	execPath(p, x, exp, nil, t)
}

func TestNodeTypeNode(t *testing.T) {
	p := `/p1/child:: node( ) `
	x := `<?xml version="1.0" encoding="UTF-8"?><p1 test="foo"><p2/></p1>`
	exp := []string{`<p2></p2>`}
	execPath(p, x, exp, nil, t)
}

func TestNodeTypeNodeAbbr(t *testing.T) {
	p := `/p1/ .`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1 test="foo"><p2/></p1>`
	exp := []string{`<p1 test="foo"><p2></p2></p1>`}
	execPath(p, x, exp, nil, t)
}

func TestNodeTypeParent(t *testing.T) {
	p := `/p1/p2/parent::node()`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1 test="foo"><p2/></p1>`
	exp := []string{`<p1 test="foo"><p2></p2></p1>`}
	execPath(p, x, exp, nil, t)
}

func TestNodeTypeParentAbbr(t *testing.T) {
	p := `/p1/p2/..`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1 test="foo"><p2/></p1>`
	exp := []string{`<p1 test="foo"><p2></p2></p1>`}
	execPath(p, x, exp, nil, t)
}

func TestFollowing(t *testing.T) {
	p := `//p3/following::node()`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1><p2><p31/><p3/><p4/></p2><p5><p6/></p5></p1>`
	exp := []string{`<p4></p4>`, `<p5><p6></p6></p5>`, `<p6></p6>`}
	execPath(p, x, exp, nil, t)
}

func TestPreceding(t *testing.T) {
	p := `//p6/preceding::node()`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1><p2><p3/><p4/></p2><p5><p6/><p7/></p5></p1>`
	exp := []string{`<p2><p3></p3><p4></p4></p2>`, `<p3></p3>`, `<p4></p4>`}
	execPath(p, x, exp, nil, t)
}

func TestPrecedingSibling(t *testing.T) {
	p := `//p4/preceding-sibling::node()`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1><p2><p3><p31/></p3><p4/><p51/></p2><p5><p6/></p5></p1>`
	exp := []string{`<p3><p31></p31></p3>`}
	execPath(p, x, exp, nil, t)
}

func TestPrecedingSibling2(t *testing.T) {
	p := `/preceding-sibling::node()`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1></p1>`
	exp := []string{}
	execPath(p, x, exp, nil, t)
}

func TestFollowingSibling(t *testing.T) {
	p := `//p2/following-sibling::node()`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1><p21/><p2><p3/><p4/></p2><p5><p6/></p5></p1>`
	exp := []string{`<p5><p6></p6></p5>`}
	execPath(p, x, exp, nil, t)
}

func TestFollowingSibling2(t *testing.T) {
	p := `/following-sibling::node()`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1></p1>`
	exp := []string{}
	execPath(p, x, exp, nil, t)
}

func TestComment(t *testing.T) {
	p := `// comment()`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1><!-- comment --></p1>`
	exp := []string{`<!-- comment -->`}
	execPath(p, x, exp, nil, t)
}

func TestCommentInvPath(t *testing.T) {
	p := `//comment()/ self :: processing-instruction ( ) `
	x := `<?xml version="1.0" encoding="UTF-8"?><p1><!-- comment --></p1>`
	exp := []string{}
	execPath(p, x, exp, nil, t)
}

func TestCommentParent(t *testing.T) {
	p := `//comment()/..`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1><!-- comment --></p1>`
	exp := []string{`<p1><!-- comment --></p1>`}
	execPath(p, x, exp, nil, t)
}

func TestText(t *testing.T) {
	p := `//text()`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1>text</p1>`
	exp := []string{`text`}
	execPath(p, x, exp, nil, t)
}

func TestTextParent(t *testing.T) {
	p := `//text()/..`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1>text</p1>`
	exp := []string{`<p1>text</p1>`}
	execPath(p, x, exp, nil, t)
}

func TestProcInst(t *testing.T) {
	p := `//processing-instruction()`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1><?proc?></p1>`
	exp := []string{`<?proc?>`}
	execPath(p, x, exp, nil, t)
}

func TestProcInstParent(t *testing.T) {
	p := `//processing-instruction()/..`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1><?proc?></p1>`
	exp := []string{`<p1><?proc?></p1>`}
	execPath(p, x, exp, nil, t)
}

func TestProcInstInvPath(t *testing.T) {
	p := `//processing-instruction()/self::comment()`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1><?proc?></p1>`
	exp := []string{}
	execPath(p, x, exp, nil, t)
}

func TestProcInst2(t *testing.T) {
	p := `//processing-instruction( 'proc2' ) `
	x := `<?xml version="1.0" encoding="UTF-8"?><p1><?proc1?><?proc2?><?proc3?></p1>`
	exp := []string{`<?proc2?>`}
	execPath(p, x, exp, nil, t)
}

func TestNamespace(t *testing.T) {
	p := `/*:p1/*:p2/namespace::*`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1 xmlns="http://foo.bar"><p2 xmlns:foo="http://test"></p2></p1>`
	exp := []string{`<?namespace http://foo.bar?>`, `<?namespace http://test?>`, `<?namespace http://www.w3.org/XML/1998/namespace?>`}
	execPath(p, x, exp, nil, t)
}

func TestRootNamespace(t *testing.T) {
	p := `/namespace::*`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1 xmlns="http://foo.bar"><p2 xmlns:foo="http://test"></p2></p1>`
	exp := []string{}
	execPath(p, x, exp, nil, t)
}

func TestNamespaceXml(t *testing.T) {
	p := `/*:p1/ * : p2 /namespace::xml`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1 xmlns="http://foo.bar"><p2 xmlns:foo="http://test"></p2></p1>`
	exp := []string{`<?namespace http://www.w3.org/XML/1998/namespace?>`}
	execPath(p, x, exp, nil, t)
}

func TestNamespaceNodeType(t *testing.T) {
	p := `/*:p1/*:p2/namespace::foo/ self :: node ( ) `
	x := `<?xml version="1.0" encoding="UTF-8"?><p1 xmlns="http://foo.bar"><p2 xmlns:foo="http://test"></p2></p1>`
	exp := []string{`<?namespace http://test?>`}
	execPath(p, x, exp, nil, t)
}

func TestNamespaceNodeWildcard(t *testing.T) {
	p := `/*:p1/*:p2/namespace::*:foo`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1 xmlns="http://foo.bar"><p2 xmlns:foo="http://test"></p2></p1>`
	exp := []string{`<?namespace http://test?>`}
	execPath(p, x, exp, nil, t)
	p = `/*:p1/*:p2/namespace::invalid:foo`
	exp = []string{}
	execPath(p, x, exp, nil, t)
}

func TestNamespaceChild(t *testing.T) {
	p := `/*:p1/namespace::*/ * : p2`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1 xmlns="http://foo.bar"><p2 xmlns:foo="http://test"></p2></p1>`
	exp := []string{}
	execPath(p, x, exp, nil, t)
}

func TestNamespaceParent(t *testing.T) {
	p := `/*:p1/*:p2/namespace::foo/..`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1 xmlns="http://foo.bar"><p2 xmlns:foo="http://test"></p2></p1>`
	exp := []string{`<p2 xmlns="http://foo.bar"></p2>`}
	execPath(p, x, exp, nil, t)
}

func TestNamespace2(t *testing.T) {
	p := `//namespace::test`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1 xmlns:foo="http://test"><p2 xmlns:test="http://foo.bar"></p2></p1>`
	exp := []string{`<?namespace http://foo.bar?>`}
	execPath(p, x, exp, nil, t)
}

func TestNamespace3(t *testing.T) {
	p := `/p1/p2/p3/namespace :: foo`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1 xmlns:foo="http://test"><p2 xmlns:foo="http://foo.bar"><p3/></p2></p1>`
	exp := []string{`<?namespace http://foo.bar?>`}
	execPath(p, x, exp, nil, t)
}

func TestNamespace4(t *testing.T) {
	p := `/test:p1/p2/p3/namespace::*`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1 xmlns="http://test"><p2 xmlns=""><p3/></p2></p1>`
	exp := []string{`<?namespace http://www.w3.org/XML/1998/namespace?>`}
	execPath(p, x, exp, map[string]string{"test": "http://test"}, t)
}

func TestNodeNamespace(t *testing.T) {
	p := `/ test : p1 `
	x := `<?xml version="1.0" encoding="UTF-8"?><p1 xmlns="http://test"/>`
	exp := []string{`<p1 xmlns="http://test"></p1>`}
	execPath(p, x, exp, map[string]string{"test": "http://test"}, t)
}

func TestNodeNamespace2(t *testing.T) {
	p := `/test:p1/test2:p2`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1 xmlns="http://test"><p2 xmlns="http://test2"></p2></p1>`
	exp := []string{`<p2 xmlns="http://test2"></p2>`}
	execPath(p, x, exp, map[string]string{"test": "http://test", "test2": "http://test2"}, t)
}

func TestNodeNamespace3(t *testing.T) {
	p := `/test:p1/foo:p2`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1 xmlns="http://test" xmlns:foo="http://foo.bar"><foo:p2></foo:p2></p1>`
	exp := []string{`<p2 xmlns="http://foo.bar"></p2>`}
	execPath(p, x, exp, map[string]string{"test": "http://test", "foo": "http://foo.bar"}, t)
}

func TestAttrNamespace(t *testing.T) {
	p := `/test:p1/@foo:test`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1 xmlns="http://test" xmlns:foo="http://foo.bar" foo:test="foo"></p1>`
	exp := []string{`<?attribute test="foo" xmlns="http://foo.bar"?>`}
	execPath(p, x, exp, map[string]string{"test": "http://test", "foo": "http://foo.bar"}, t)
	p = `/test:p1/@test:test`
	exp = []string{}
	execPath(p, x, exp, map[string]string{"test": "http://test", "foo": "http://foo.bar"}, t)
}

func TestWildcard(t *testing.T) {
	p := `/p1/*`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1 xmlns:test="http://test" xmlns:foo="http://foo.bar"><p1/><test:p2/><foo:p3/></p1>`
	exp := []string{`<p1></p1>`, `<p2 xmlns="http://test"></p2>`, `<p3 xmlns="http://foo.bar"></p3>`}
	execPath(p, x, exp, nil, t)
}

func TestWildcardNS(t *testing.T) {
	p := `//*:p1`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1 xmlns="http://test" xmlns:foo="http://foo.bar"><foo:p1/></p1>`
	exp := []string{`<p1 xmlns="http://test"><p1 xmlns="http://foo.bar"></p1></p1>`, `<p1 xmlns="http://foo.bar"></p1>`}
	execPath(p, x, exp, map[string]string{"test": "http://test", "foo": "http://foo.bar"}, t)
}

func TestWildcardLocal(t *testing.T) {
	p := `//foo:*`
	x := `<?xml version="1.0" encoding="UTF-8"?><p3 xmlns="http://test" xmlns:foo="http://foo.bar"><foo:p1/><foo:p2/></p3>`
	exp := []string{`<p1 xmlns="http://foo.bar"></p1>`, `<p2 xmlns="http://foo.bar"></p2>`}
	execPath(p, x, exp, map[string]string{"test": "http://test", "foo": "http://foo.bar"}, t)
}

func TestWildcardNSAttr(t *testing.T) {
	p := `/p1/@*:attr`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1 xmlns:foo="http://foo.bar" foo:attr="test"/>`
	exp := []string{`<?attribute attr="test" xmlns="http://foo.bar"?>`}
	execPath(p, x, exp, map[string]string{"test": "http://test", "foo": "http://foo.bar"}, t)
}

func TestWildcardLocalAttr(t *testing.T) {
	p := `/p1/@foo:*`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1 xmlns:foo="http://foo.bar" foo:attr="test" foo:bar="foobar"/>`
	exp := []string{`<?attribute attr="test" xmlns="http://foo.bar"?>`, `<?attribute bar="foobar" xmlns="http://foo.bar"?>`}
	execPath(p, x, exp, map[string]string{"test": "http://test", "foo": "http://foo.bar"}, t)
}

func TestPredicate1(t *testing.T) {
	p := `/p1/p2 [ p3 ] `
	x := `<?xml version="1.0" encoding="UTF-8"?><p1><p2/><p2><p3/></p2></p1>`
	exp := []string{`<p2><p3></p3></p2>`}
	execPath(p, x, exp, nil, t)
}

func TestPredicate2(t *testing.T) {
	p := `/p1/p2 [ 2 ] `
	x := `<?xml version="1.0" encoding="UTF-8"?><p1><p2/><p2><p3/></p2><p2><p4/></p2></p1>`
	exp := []string{`<p2><p3></p3></p2>`}
	execPath(p, x, exp, nil, t)
	p = `/p1/p2 [ position() = 2 ] `
	execPath(p, x, exp, nil, t)
}

func TestPredicate3(t *testing.T) {
	p := `/p1/p2 [ last() ] /p3`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1><p2><p3/></p2><p2><p3/></p2><p2><p3 foo="bar"/></p2></p1>`
	exp := []string{`<p3 foo="bar"></p3>`}
	execPath(p, x, exp, nil, t)
}

func TestPredicate4(t *testing.T) {
	p := `/p1/p2[not(@test)]`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1><p2/><p2 test="foo"/><p2/></p1>`
	exp := []string{`<p2></p2>`, `<p2></p2>`}
	execPath(p, x, exp, nil, t)
}

func TestPredicate5(t *testing.T) {
	p := `/p1/p2[.//p1]`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1><p2/><p2 test="t"/><p2><p1>lkj</p1></p2></p1>`
	exp := []string{`<p2><p1>lkj</p1></p2>`}
	execPath(p, x, exp, nil, t)
}

func TestPredicate6(t *testing.T) {
	p := `/p1/p2[//p1]`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1><p2/><p2 test="t"/><p2><p1>lkj</p1></p2></p1>`
	exp := []string{`<p2></p2>`, `<p2 test="t"></p2>`, `<p2><p1>lkj</p1></p2>`}
	execPath(p, x, exp, nil, t)
}

func TestPredicate7(t *testing.T) {
	p := `/p3[//p1]`
	x := `<?xml version="1.0" encoding="UTF-8"?><p1><p2/><p2 test="t"/><p2><p1>lkj</p1></p2></p1>`
	exp := []string{}
	execPath(p, x, exp, nil, t)
}

func TestPredicate8(t *testing.T) {
	p := `/test[@attr1='a' and @attr2='b']/a`
	x := `<?xml version="1.0" encoding="UTF-8"?><test attr1="a" attr2="b"><a>aaaa</a><b>bbb</b></test>`
	exp := []string{`<a>aaaa</a>`}
	execPath(p, x, exp, nil, t)
}

func TestUnion(t *testing.T) {
	p := `/test/test2 | /test/test3`
	x := `<?xml version="1.0" encoding="UTF-8"?><test><test2>foobar</test2><test3>hamneggs</test3></test>`
	exp := []string{"<test2>foobar</test2>", "<test3>hamneggs</test3>"}
	execPath(p, x, exp, nil, t)
}

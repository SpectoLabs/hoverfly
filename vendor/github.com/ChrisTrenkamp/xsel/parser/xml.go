package parser

import (
	"encoding/xml"
	"io"

	"github.com/ChrisTrenkamp/xsel/node"
	"golang.org/x/net/html/charset"
)

type XmlElement struct {
	space, local string
}

func (x XmlElement) Space() string {
	return x.space
}

func (x XmlElement) Local() string {
	return x.local
}

type XmlNamespace struct {
	prefix, value string
}

func (x XmlNamespace) Prefix() string {
	return x.prefix
}

func (x XmlNamespace) NamespaceValue() string {
	return x.value
}

type XmlAttribute struct {
	space, local, value string
}

func (x XmlAttribute) Space() string {
	return x.space
}

func (x XmlAttribute) Local() string {
	return x.local
}

func (x XmlAttribute) AttributeValue() string {
	return x.value
}

type XmlCharData struct {
	value string
}

func (x XmlCharData) CharDataValue() string {
	return x.value
}

type XmlComment struct {
	value string
}

func (x XmlComment) CommentValue() string {
	return x.value
}

type XmlProcInst struct {
	target, value string
}

func (x XmlProcInst) Target() string {
	return x.target
}

func (x XmlProcInst) ProcInstValue() string {
	return x.value
}

var emptyAttrs = make([]XmlAttribute, 0)
var emptyNamespaces = make([]XmlNamespace, 0)

type XmlParseOptions func(d *xml.Decoder)

// Creates a Parser that reads the given XML document.
func ReadXml(in io.Reader, opts ...XmlParseOptions) Parser {
	xmlReader := xml.NewDecoder(in)
	xmlReader.CharsetReader = charset.NewReaderLabel

	for _, i := range opts {
		i(xmlReader)
	}

	namespaces := emptyNamespaces
	nsPos := 0

	attrs := emptyAttrs
	attrPos := 0

	return func() (node.Node, bool, error) {
		if nsPos < len(namespaces) {
			n := namespaces[nsPos]
			nsPos++

			return n, false, nil
		}

		if attrPos < len(attrs) {
			a := attrs[attrPos]
			attrPos++

			return a, false, nil
		}

		attrs = emptyAttrs
		attrPos = 0
		namespaces = emptyNamespaces
		nsPos = 0
		tok, err := xmlReader.Token()

		if err != nil {
			return nil, false, err
		}

		switch n := tok.(type) {
		case xml.StartElement:
			namespaces = createNamespaces(n.Attr)
			attrs = createAttrs(n.Attr)
			return XmlElement{
				space: n.Name.Space,
				local: n.Name.Local,
			}, false, nil
		case xml.CharData:
			return XmlCharData{
				value: (string)(n),
			}, false, nil
		case xml.Comment:
			return XmlComment{
				value: (string)(n),
			}, false, nil
		case xml.ProcInst:
			return XmlProcInst{
				target: n.Target,
				value:  string(n.Inst),
			}, false, nil
		}

		//case xml.EndElement:
		return nil, true, nil
	}
}

const XMLNS = "xmlns"

func createNamespaces(attrs []xml.Attr) []XmlNamespace {
	ret := make([]XmlNamespace, 0, 1)
	ns := XmlNamespace{
		prefix: "xml",
		value:  "http://www.w3.org/XML/1998/namespace",
	}

	ret = append(ret, ns)

	for _, i := range attrs {
		if i.Name.Space == "" && i.Name.Local == XMLNS {
			ns = XmlNamespace{
				prefix: "",
				value:  i.Value,
			}

			ret = append(ret, ns)
		}

		if i.Name.Local == XMLNS {
			ns = XmlNamespace{
				prefix: i.Name.Space,
				value:  i.Value,
			}

			ret = append(ret, ns)
		}
	}

	return ret
}

func createAttrs(attrs []xml.Attr) []XmlAttribute {
	ret := make([]XmlAttribute, 0, len(attrs))

	for _, i := range attrs {
		if i.Name.Space == XMLNS || i.Name.Local == XMLNS {
			continue
		}

		next := XmlAttribute{
			space: i.Name.Space,
			local: i.Name.Local,
			value: i.Value,
		}

		ret = append(ret, next)
	}

	return ret
}

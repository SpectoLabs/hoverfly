package parser

import (
	"fmt"
	"io"
	"strings"

	"github.com/ChrisTrenkamp/xsel/node"
	"golang.org/x/net/html"
)

type HtmlElement struct {
	local string
}

func (h HtmlElement) Space() string {
	return ""
}

func (h HtmlElement) Local() string {
	return h.local
}

type HtmlAttribute struct {
	local, value string
}

func (h HtmlAttribute) Space() string {
	return ""
}

func (h HtmlAttribute) Local() string {
	return h.local
}

func (h HtmlAttribute) AttributeValue() string {
	return h.value
}

type HtmlCharData struct {
	value string
}

func (h HtmlCharData) CharDataValue() string {
	return h.value
}

type HtmlComment struct {
	value string
}

func (h HtmlComment) CommentValue() string {
	return h.value
}

var emptyHtmlAttrs = make([]HtmlAttribute, 0)

type htmlParser struct {
	node               *html.Node
	attrs              []HtmlAttribute
	attrPos            int
	emitSelfClosingTag bool
	nodeEmitted        bool
	crawlToParent      bool
}

func (x *htmlParser) Pull() (node.Node, bool, error) {
	if x.attrPos < len(x.attrs) {
		a := x.attrs[x.attrPos]
		x.attrPos++

		return a, false, nil
	}

	x.attrs = emptyHtmlAttrs
	x.attrPos = 0

	if x.emitSelfClosingTag {
		x.emitSelfClosingTag = false

		return nil, true, nil
	}

	if x.nodeEmitted {
		x.nodeEmitted = false

		if x.node.FirstChild != nil {
			x.node = x.node.FirstChild
		} else if x.node.NextSibling != nil {
			x.node = x.node.NextSibling
		} else {
			x.crawlToParent = true
		}
	}

	if x.crawlToParent {
		x.crawlToParent = false

		if x.node.Parent == nil {
			return nil, false, io.EOF
		}

		x.node = x.node.Parent

		if x.node.NextSibling == nil {
			x.crawlToParent = true
		} else {
			x.node = x.node.NextSibling
		}

		return nil, true, nil
	}

	switch x.node.Type {
	case html.ErrorNode:
		return nil, false, fmt.Errorf("error parsing html document")
	case html.RawNode:
		return nil, false, fmt.Errorf("encountered raw node")
	case html.DocumentNode:
		x.node = x.node.FirstChild

		if x.node.Type != html.DoctypeNode {
			return nil, false, fmt.Errorf("doctype declaration not found")
		}

		return x.Pull()
	case html.DoctypeNode:
		x.node = x.node.NextSibling
		return x.Pull()
	case html.ElementNode:
		local := getLocalName(x.node.Data)
		x.attrs = createHtmlAttrs(x.node.Attr)
		x.nodeEmitted = true
		if x.node.FirstChild == nil {
			x.emitSelfClosingTag = true
		}

		return HtmlElement{
			local: local,
		}, false, nil
	case html.TextNode:
		x.nodeEmitted = true
		return HtmlCharData{
			value: (string)(x.node.Data),
		}, false, nil
	case html.CommentNode:
		x.nodeEmitted = true
		return HtmlComment{
			value: (string)(x.node.Data),
		}, false, nil
	}

	return nil, true, nil
}

func createHtmlAttrs(attrs []html.Attribute) []HtmlAttribute {
	ret := make([]HtmlAttribute, 0)

	for _, i := range attrs {
		name := i.Key

		if name == xmlns {
			continue
		}

		if strings.HasPrefix(name, xmlns+":") {
			continue
		}

		name = getLocalName(name)

		attr := HtmlAttribute{
			local: name,
			value: i.Val,
		}

		ret = append(ret, attr)
	}

	return ret
}

func getLocalName(name string) string {
	if strings.Contains(name, ":") {
		spl := strings.SplitN(name, ":", 2)
		return spl[1]
	}

	return name
}

// Creates a Parser that reads the given HTML document.
func ReadHtml(in io.Reader) (Parser, error) {
	documentNode, err := html.Parse(in)

	if err != nil {
		return nil, err
	}

	return &htmlParser{
		node:               documentNode,
		attrs:              emptyHtmlAttrs,
		attrPos:            0,
		emitSelfClosingTag: false,
		nodeEmitted:        false,
		crawlToParent:      false,
	}, nil
}

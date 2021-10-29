package exec

import (
	"math"
	"strconv"
	"strings"

	"github.com/ChrisTrenkamp/xsel/node"
	"github.com/ChrisTrenkamp/xsel/store"
)

type Result interface {
	String() string
	Number() float64
	Bool() bool
}

type Bool bool

func (b Bool) String() string {
	if b {
		return "true"
	}

	return "false"
}

func (b Bool) Number() float64 {
	if b {
		return 1.0
	}

	return 0.0
}

func (b Bool) Bool() bool {
	return bool(b)
}

type Number float64

func (n Number) String() string {
	if math.IsInf(float64(n), 1) {
		return "Infinity"
	}

	if math.IsInf(float64(n), -1) {
		return "-Infinity"
	}

	return strconv.FormatFloat(float64(n), 'f', -1, 64)
}

func (n Number) Number() float64 {
	return float64(n)
}

func (n Number) Bool() bool {
	return n != 0
}

type String string

func (n String) String() string {
	return string(n)
}

func (n String) Number() float64 {
	ret, err := strconv.ParseFloat(string(n), 64)

	if err != nil {
		return math.NaN()
	}

	return ret
}

func (n String) Bool() bool {
	return len(n) > 0
}

type NodeSet []store.Cursor

func (n NodeSet) String() string {
	if len(n) == 0 {
		return ""
	}

	return getCursorString(n[0])
}

func (n NodeSet) Number() float64 {
	return getStringNumber(n.String())
}

func (n NodeSet) Bool() bool {
	return len(n) > 0
}

func getStringNumber(str string) float64 {
	ret, err := strconv.ParseFloat(str, 64)

	if err != nil {
		return math.NaN()
	}

	return ret
}

func getCursorString(c store.Cursor) string {
	buf := strings.Builder{}
	getCursorStringValue(&buf, c)
	return buf.String()
}

func getCursorStringValue(buf *strings.Builder, c store.Cursor) {
	n := c.Node()

	switch v := n.(type) {
	case node.Namespace:
		buf.WriteString(v.NamespaceValue())
	case node.Attribute:
		buf.WriteString(v.AttributeValue())
	case node.CharData:
		writeCharData(buf, v)
	case node.Comment:
		buf.WriteString(v.CommentValue())
	case node.ProcInst:
		buf.WriteString(v.ProcInstValue())
	case node.Element:
		getElementStringValue(buf, c)
	case node.Root:
		getElementStringValue(buf, c)
	}

}

func getElementStringValue(buf *strings.Builder, c store.Cursor) {
	children := c.Children()
	pos := 0

	for _, n := range children {
		switch v := n.Node().(type) {
		case node.Element:
			getElementStringValue(buf, n)
		case node.CharData:
			writeCharData(buf, v)
		}

		pos++
	}
}

func writeCharData(buf *strings.Builder, c node.CharData) {
	buf.WriteString(c.CharDataValue())
}

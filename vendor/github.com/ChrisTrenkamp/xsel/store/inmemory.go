package store

import (
	"errors"
	"io"

	"github.com/ChrisTrenkamp/xsel/node"
	"github.com/ChrisTrenkamp/xsel/parser"
)

type rootInMemoryNode struct{}

type InMemory struct {
	node       node.Node
	pos        int
	parent     *InMemory
	namespaces []Cursor
	attributes []Cursor
	nodes      []Cursor
}

func initElement() InMemory {
	return InMemory{
		namespaces: make([]Cursor, 0),
		attributes: make([]Cursor, 0),
		nodes:      make([]Cursor, 0),
	}
}

var emptyCursor = make([]Cursor, 0)

func initNonElement() InMemory {
	return InMemory{
		namespaces: emptyCursor,
		attributes: emptyCursor,
		nodes:      emptyCursor,
	}
}

// Gathers and stores node.Node's in memory.
func CreateInMemory(parse parser.Parser) (*InMemory, error) {
	root := initElement()
	root.node = rootInMemoryNode{}
	root.pos = 0
	root.parent = &root
	err := createInMemory(&root, parse, 0)
	return &root, err
}

func createInMemory(cursor *InMemory, parse parser.Parser, pos int) error {
	n, isEnd, err := parse()

	if errors.Is(err, io.EOF) {
		return nil
	}

	if err != nil {
		return err
	}

	if isEnd {
		return createInMemory(cursor.parent, parse, pos)
	}

	switch v := n.(type) {
	case node.Namespace:
		pos = addNamespace(v, cursor, pos)
	case node.Attribute:
		pos++
		cursor.attributes = append(cursor.attributes, createNonElement(v, cursor, pos))
	case node.Element:
		pos++
		next, pos := createElement(v, cursor, pos)
		cursor.nodes = append(cursor.nodes, next)
		return createInMemory(next, parse, pos)
	default:
		pos++
		cursor.nodes = append(cursor.nodes, createNonElement(v, cursor, pos))
	}

	return createInMemory(cursor, parse, pos)
}

func addNamespace(ns node.Namespace, cursor *InMemory, pos int) int {
	toReplace := -1

	for pos, i := range cursor.namespaces {
		nsTest := i.(*InMemory).node.(node.Namespace)

		if nsTest.Prefix() == ns.Prefix() {
			toReplace = pos
			break
		}
	}

	if toReplace < 0 {
		cursor.namespaces = append(cursor.namespaces, createNonElement(ns, cursor, pos))
		return pos + 1
	}

	nsPos := cursor.namespaces[toReplace].(*InMemory).pos
	cursor.namespaces[toReplace] = createNonElement(ns, cursor, nsPos)
	return pos
}

func createNonElement(node node.Node, parent *InMemory, pos int) *InMemory {
	next := initNonElement()
	next.node = node
	next.pos = pos
	next.parent = parent

	return &next
}

func createElement(node node.Node, parent *InMemory, pos int) (*InMemory, int) {
	next := initElement()
	next.node = node
	next.pos = pos
	next.parent = parent

	ns := make([]Cursor, len(parent.namespaces))
	copy(ns, parent.namespaces)

	next.namespaces = ns

	for _, i := range next.namespaces {
		pos++
		i.(*InMemory).pos = pos
	}

	return &next, pos + len(next.namespaces)
}

func (c *InMemory) Pos() int {
	return c.pos
}

func (c *InMemory) Node() node.Node {
	return c.node
}

func (c *InMemory) Namespaces() []Cursor {
	return c.namespaces
}

func (c *InMemory) Attributes() []Cursor {
	return c.attributes
}

func (c *InMemory) Children() []Cursor {
	return c.nodes
}

func (c *InMemory) Parent() Cursor {
	return c.parent
}

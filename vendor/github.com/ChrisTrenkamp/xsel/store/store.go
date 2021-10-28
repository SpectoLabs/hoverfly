package store

import "github.com/ChrisTrenkamp/xsel/node"

// Cursor's are used for tracking the position and child-parent relationships of node.Node's.
type Cursor interface {
	// Pos returns the position of the Node.
	// Implementations MUST store Cursor's in order by the following rules:
	// Positions MUST be unique.
	// The node.Root MUST have a position of 0.
	// Namespaces, Attributes, and Children MUST be in order, from least to greatest.
	// Namespace positions MUST be less than the Attribute and Children positions.
	// Attribute positions MUST be less than Children positions.
	Pos() int
	// Node returns the node.Node stored in this Cursor.
	Node() node.Node
	Namespaces() []Cursor
	Attributes() []Cursor
	Children() []Cursor
	Parent() Cursor
}

// A convenience method for retrieving a node.Attribute.
func GetAttribute(c Cursor, space, local string) (node.Attribute, bool) {
	for _, a := range c.Attributes() {
		attr := a.Node().(node.Attribute)

		if attr.Space() == space && attr.Local() == local {
			return attr, true
		}
	}

	return nil, false
}

package parser

import "github.com/ChrisTrenkamp/xsel/node"

// Feeds a stream of nodes to be stored in a store.Cursor.
// Implementations of this method must not emit a node.Root; store.Cursor will create it.
// For every node.Element emitted, it must be terminated by returning a false value.
// When emitting a node, it will either be added to the root node, or the last node.Element that was not terminated.
// node.Namespace's MUST be emitted before node.Attributes.
// node.Attribute's MUST be emitted before the element's children.
// When finished emitting nodes, return io.EOF.
type Parser func() (node.Node, bool, error)

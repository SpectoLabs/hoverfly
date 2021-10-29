package node

// The catch-all interface for all nodes.
type Node interface{}

// The container for all nodes in the tree.  There must only be one, and it must exist at the root.
type Root interface {
	Node
}

// Implemented by Element's and Attributes.  It is a convenience interface for querying by names.
type NamedNode interface {
	Space() string
	Local() string
}

type Element interface {
	Node
	NamedNode
}

type Namespace interface {
	Node
	Prefix() string
	NamespaceValue() string
}

type Attribute interface {
	Node
	NamedNode
	AttributeValue() string
}

type CharData interface {
	Node
	CharDataValue() string
}

type Comment interface {
	Node
	CommentValue() string
}

type ProcInst interface {
	Node
	Target() string
	ProcInstValue() string
}

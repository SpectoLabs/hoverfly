package exec

import "github.com/ChrisTrenkamp/xsel/store"

// ContextSettings allows you to add namespace mappings, create new functions,
// and add variable bindings to your XPath query.
type ContextSettings struct {
	NamespaceDecls  map[string]string
	FunctionLibrary map[XmlName]Function
	Variables       map[XmlName]Result
	// Context is the initial position to run XPath queries.  This is initialized to the root of
	// the document.  This may be overridden to point to a different position in the document.
	// When overriding the Context, it must be a Cursor contained within the root, or bad things can happen!
	Context store.Cursor
}

type ContextApply func(c *ContextSettings)

type exprContext struct {
	root             store.Cursor
	result           Result
	contextPosition  int
	builtinFunctions map[XmlName]Function
	ContextSettings
}

type Context interface {
	Result() Result
	ContextPosition() int
}

func (c *exprContext) Result() Result {
	return c.result
}

func (c *exprContext) ContextPosition() int {
	return c.contextPosition
}

func (e *exprContext) copy() exprContext {
	return exprContext{
		root:             e.root,
		result:           e.result,
		contextPosition:  e.contextPosition,
		builtinFunctions: builtinFunctions,
		ContextSettings:  e.ContextSettings,
	}
}

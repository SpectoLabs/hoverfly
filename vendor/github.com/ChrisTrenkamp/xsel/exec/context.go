package exec

import "github.com/ChrisTrenkamp/xsel/store"

// ContextSettings allows you to add namespace mappings, create new functions,
// and add variable bindings to your XPath query.
type ContextSettings struct {
	NamespaceDecls  map[string]string
	FunctionLibrary map[XmlName]Function
	Variables       map[XmlName]Result
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

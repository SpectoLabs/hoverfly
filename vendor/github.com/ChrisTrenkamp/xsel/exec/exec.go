package exec

import (
	"fmt"

	"github.com/ChrisTrenkamp/xsel/grammar"
	"github.com/ChrisTrenkamp/xsel/store"
	"github.com/pkg/errors"
)

// Executes an XPath query against the given Cursor and returns the result.
func Exec(cursor store.Cursor, expr *grammar.Grammar, settings ...ContextApply) (Result, error) {
	contextSettings := ContextSettings{
		Variables:       make(map[XmlName]Result),
		FunctionLibrary: make(map[XmlName]Function),
		NamespaceDecls:  make(map[string]string),
	}

	for _, i := range settings {
		i(&contextSettings)
	}

	context := &exprContext{
		root:             cursor,
		result:           Result(NodeSet{cursor}),
		contextPosition:  0,
		builtinFunctions: builtinFunctions,
		ContextSettings:  contextSettings,
	}

	err := execRecover(context, expr)

	if err != nil {
		return nil, err
	}

	return context.result, nil
}

func execRecover(context *exprContext, expr *grammar.Grammar) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Wrapf(fmt.Errorf("xpath query panic"), "%s", r)
		}
	}()

	err = execContext(context, expr)

	return
}

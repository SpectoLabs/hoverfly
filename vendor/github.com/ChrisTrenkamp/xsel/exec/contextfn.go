package exec

import (
	"fmt"
	"strings"

	"github.com/ChrisTrenkamp/xsel/grammar"
	"github.com/ChrisTrenkamp/xsel/grammar/parser/bsr"
	"github.com/ChrisTrenkamp/xsel/grammar/parser/symbols"
)

type contextFn func(context *exprContext, expr *grammar.Grammar) error

var contextFunctions = map[symbols.NT]contextFn{}

func init() {
	contextFunctions[symbols.NT_Literal] = execLiteral
	contextFunctions[symbols.NT_UnionExprUnion] = execUnionExprUnion
	contextFunctions[symbols.NT_FunctionCall] = execFunctionCall
	contextFunctions[symbols.NT_VariableReference] = execVariableReference
}

func execLiteral(context *exprContext, expr *grammar.Grammar) error {
	literal := expr.GetString()
	literal = literal[1 : len(literal)-1]

	context.result = String(literal)
	return nil
}

func execUnionExprUnion(context *exprContext, expr *grammar.Grammar) error {
	left, right, err := leftRightIndependentResult(context, expr)

	if err != nil {
		return err
	}

	leftNodeSet, lok := left.(NodeSet)
	rightNodeSet, rok := right.(NodeSet)

	if !lok || !rok {
		return fmt.Errorf("cannot union non-NodeSet's")
	}

	context.result = unionCleanup(append(leftNodeSet, rightNodeSet...))
	return nil
}

func unionCleanup(nextResult NodeSet) NodeSet {
	return unique(nextResult)
}

func execFunctionCall(context *exprContext, expr *grammar.Grammar) error {
	children := make([]*bsr.BSR, 0, 2)

	for _, cn := range expr.BSR.GetAllNTChildren() {
		for _, c := range cn {
			children = append(children, &c)
		}
	}

	bsrs := make([]*bsr.BSR, 0)
	gatherFunctionArgs(children[1], &bsrs)

	args := make([]Result, 0, len(bsrs))

	for _, i := range bsrs {
		nextContext := context.copy()
		nextExpr := expr.Next(i)

		if err := execContext(&nextContext, nextExpr); err != nil {
			return err
		}

		args = append(args, nextContext.result)
	}

	qname, err := GetQName(expr.Next(children[0]).GetString(), context.NamespaceDecls)

	if err != nil {
		return err
	}

	fn := context.FunctionLibrary[qname]

	if fn == nil {
		fn = context.builtinFunctions[qname]
	}

	if fn == nil {
		return fmt.Errorf("could not find function %s", qname)
	}

	result, err := fn(context, args...)

	if err != nil {
		return fmt.Errorf("error invoking function %s: %s", qname, err)
	}

	context.result = result

	return nil
}

func gatherFunctionArgs(b *bsr.BSR, args *[]*bsr.BSR) {
	children := make([]*bsr.BSR, 0, 2)

	for _, cn := range b.GetAllNTChildren() {
		for _, c := range cn {
			children = append(children, &c)
		}
	}

	name := b.Label.Slot().NT

	if (name == symbols.NT_FunctionSignature || name == symbols.NT_FunctionCallArgumentList) && len(children) == 1 {
		gatherFunctionArgs(children[0], args)
		return
	}

	if len(children) >= 1 {
		*args = append(*args, children[0])
	}

	if len(children) == 2 {
		gatherFunctionArgs(children[1], args)
	}
}

func execVariableReference(context *exprContext, expr *grammar.Grammar) error {
	variableStr := expr.Next(expr.BSR).GetString()
	variableStr = strings.TrimSpace(variableStr)

	if strings.HasPrefix(variableStr, "$") {
		variableStr = variableStr[1:]
	}

	qname, err := GetQName(variableStr, context.NamespaceDecls)

	if err != nil {
		return err
	}

	variable := context.Variables[qname]

	if variable == nil {
		return fmt.Errorf("could not find variable %s", qname)
	}

	context.result = variable
	return nil
}

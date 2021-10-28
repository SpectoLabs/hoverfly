package exec

import (
	"github.com/ChrisTrenkamp/xsel/grammar"
	"github.com/ChrisTrenkamp/xsel/grammar/parser/symbols"
)

func init() {
	contextFunctions[symbols.NT_OrExprOr] = execOrExprOr
	contextFunctions[symbols.NT_AndExprAnd] = execAndExprAnd
	contextFunctions[symbols.NT_EqualityExprEqual] = execEqualityExprEqual
	contextFunctions[symbols.NT_EqualityExprNotEqual] = execEqualityExprNotEqual
	contextFunctions[symbols.NT_RelationalExprLessThan] = execRelationalExprLessThan
	contextFunctions[symbols.NT_RelationalExprGreaterThan] = execRelationalExprGreaterThan
	contextFunctions[symbols.NT_RelationalExprLessThanOrEqual] = execRelationalExprLessThanOrEqual
	contextFunctions[symbols.NT_RelationalExprGreaterThanOrEqual] = execRelationalExprGreaterThanOrEqual
}

func execOrExprOr(context *exprContext, expr *grammar.Grammar) error {
	left, right, err := leftRightIndependentResult(context, expr)

	if err != nil {
		return err
	}

	leftBool := left.Bool()
	rightBool := right.Bool()

	context.result = Bool(leftBool || rightBool)
	return nil
}

func execAndExprAnd(context *exprContext, expr *grammar.Grammar) error {
	left, right, err := leftRightIndependentResult(context, expr)

	if err != nil {
		return err
	}

	leftBool := left.Bool()
	rightBool := right.Bool()

	context.result = Bool(leftBool && rightBool)
	return nil
}

func execEqualityExprEqual(context *exprContext, expr *grammar.Grammar) error {
	left, right, err := leftRightIndependentResult(context, expr)

	if err != nil {
		return err
	}

	leftNodeSet, leftNodeSetOk := left.(NodeSet)
	rightNodeSet, rightNodeSetOk := right.(NodeSet)

	if leftNodeSetOk && rightNodeSetOk {
		for _, leftNode := range leftNodeSet {
			for _, rightNode := range rightNodeSet {
				if getCursorString(leftNode) == getCursorString(rightNode) {
					context.result = Bool(true)
					return nil
				}
			}
		}

		context.result = Bool(false)
		return nil
	}

	leftNumber, leftNumberOk := left.(Number)

	if leftNumberOk && rightNodeSetOk {
		for _, rightNode := range rightNodeSet {
			if leftNumber == Number(getStringNumber(getCursorString(rightNode))) {
				context.result = Bool(true)
				return nil
			}
		}

		context.result = Bool(false)
		return nil
	}

	rightNumber, rightNumberOk := right.(Number)

	if leftNodeSetOk && rightNumberOk {
		for _, leftNode := range leftNodeSet {
			if Number(getStringNumber(getCursorString(leftNode))) == rightNumber {
				context.result = Bool(true)
				return nil
			}
		}

		context.result = Bool(false)
		return nil
	}

	leftString, leftStringOk := left.(String)

	if leftStringOk && rightNodeSetOk {
		for _, rightNode := range rightNodeSet {
			if leftString == String(getCursorString(rightNode)) {
				context.result = Bool(true)
				return nil
			}
		}

		context.result = Bool(false)
		return nil
	}

	rightString, rightStringOk := right.(String)

	if leftNodeSetOk && rightStringOk {
		for _, leftNode := range leftNodeSet {
			if String(getCursorString(leftNode)) == rightString {
				context.result = Bool(true)
				return nil
			}
		}

		context.result = Bool(false)
		return nil
	}

	leftBool, leftBoolOk := left.(Bool)

	if leftBoolOk && rightNodeSetOk {
		context.result = Bool(bool(leftBool) == rightNodeSet.Bool())
		return nil
	}

	rightBool, rightBoolOk := right.(Bool)

	if leftNodeSetOk && rightBoolOk {
		context.result = Bool(leftNodeSet.Bool() == bool(rightBool))
		return nil
	}

	if leftBoolOk || rightBoolOk {
		context.result = Bool(left.Bool() == right.Bool())
		return nil
	}

	if leftNumberOk || rightNumberOk {
		context.result = Bool(left.Number() == right.Number())
		return nil
	}

	context.result = Bool(left.String() == right.String())
	return nil
}

func execEqualityExprNotEqual(context *exprContext, expr *grammar.Grammar) error {
	left, right, err := leftRightIndependentResult(context, expr)

	if err != nil {
		return err
	}

	leftNodeSet, leftNodeSetOk := left.(NodeSet)
	rightNodeSet, rightNodeSetOk := right.(NodeSet)

	if leftNodeSetOk && rightNodeSetOk {
		for _, leftNode := range leftNodeSet {
			for _, rightNode := range rightNodeSet {
				if getCursorString(leftNode) != getCursorString(rightNode) {
					context.result = Bool(true)
					return nil
				}
			}
		}

		context.result = Bool(false)
		return nil
	}

	leftNumber, leftNumberOk := left.(Number)

	if leftNumberOk && rightNodeSetOk {
		for _, rightNode := range rightNodeSet {
			if leftNumber != Number(getStringNumber(getCursorString(rightNode))) {
				context.result = Bool(true)
				return nil
			}
		}

		context.result = Bool(false)
		return nil
	}

	rightNumber, rightNumberOk := right.(Number)

	if leftNodeSetOk && rightNumberOk {
		for _, leftNode := range leftNodeSet {
			if Number(getStringNumber(getCursorString(leftNode))) != rightNumber {
				context.result = Bool(true)
				return nil
			}
		}

		context.result = Bool(false)
		return nil
	}

	leftString, leftStringOk := left.(String)

	if leftStringOk && rightNodeSetOk {
		for _, rightNode := range rightNodeSet {
			if leftString != String(getCursorString(rightNode)) {
				context.result = Bool(true)
				return nil
			}
		}

		context.result = Bool(false)
		return nil
	}

	rightString, rightStringOk := right.(String)

	if leftNodeSetOk && rightStringOk {
		for _, leftNode := range leftNodeSet {
			if String(getCursorString(leftNode)) != rightString {
				context.result = Bool(true)
				return nil
			}
		}

		context.result = Bool(false)
		return nil
	}

	leftBool, leftBoolOk := left.(Bool)

	if leftBoolOk && rightNodeSetOk {
		context.result = Bool(bool(leftBool) != rightNodeSet.Bool())
		return nil
	}

	rightBool, rightBoolOk := right.(Bool)

	if leftNodeSetOk && rightBoolOk {
		context.result = Bool(leftNodeSet.Bool() != bool(rightBool))
		return nil
	}

	if leftBoolOk || rightBoolOk {
		context.result = Bool(left.Bool() != right.Bool())
		return nil
	}

	if leftNumberOk || rightNumberOk {
		context.result = Bool(left.Number() != right.Number())
		return nil
	}

	context.result = Bool(left.String() != right.String())
	return nil
}

func execRelationalExprLessThan(context *exprContext, expr *grammar.Grammar) error {
	left, right, err := leftRightIndependentResult(context, expr)

	if err != nil {
		return err
	}

	leftNodeSet, leftNodeSetOk := left.(NodeSet)
	rightNodeSet, rightNodeSetOk := right.(NodeSet)

	if leftNodeSetOk && rightNodeSetOk {
		for _, leftNode := range leftNodeSet {
			for _, rightNode := range rightNodeSet {
				if getCursorString(leftNode) < getCursorString(rightNode) {
					context.result = Bool(true)
					return nil
				}
			}
		}

		context.result = Bool(false)
		return nil
	}

	leftNumber, leftNumberOk := left.(Number)

	if leftNumberOk && rightNodeSetOk {
		for _, rightNode := range rightNodeSet {
			if leftNumber < Number(getStringNumber(getCursorString(rightNode))) {
				context.result = Bool(true)
				return nil
			}
		}

		context.result = Bool(false)
		return nil
	}

	rightNumber, rightNumberOk := right.(Number)

	if leftNodeSetOk && rightNumberOk {
		for _, leftNode := range leftNodeSet {
			if Number(getStringNumber(getCursorString(leftNode))) < rightNumber {
				context.result = Bool(true)
				return nil
			}
		}

		context.result = Bool(false)
		return nil
	}

	leftString, leftStringOk := left.(String)

	if leftStringOk && rightNodeSetOk {
		for _, rightNode := range rightNodeSet {
			if leftString < String(getCursorString(rightNode)) {
				context.result = Bool(true)
				return nil
			}
		}

		context.result = Bool(false)
		return nil
	}

	rightString, rightStringOk := right.(String)

	if leftNodeSetOk && rightStringOk {
		for _, leftNode := range leftNodeSet {
			if String(getCursorString(leftNode)) < rightString {
				context.result = Bool(true)
				return nil
			}
		}

		context.result = Bool(false)
		return nil
	}

	context.result = Bool(left.Number() < right.Number())
	return nil
}

func execRelationalExprLessThanOrEqual(context *exprContext, expr *grammar.Grammar) error {
	left, right, err := leftRightIndependentResult(context, expr)

	if err != nil {
		return err
	}

	leftNodeSet, leftNodeSetOk := left.(NodeSet)
	rightNodeSet, rightNodeSetOk := right.(NodeSet)

	if leftNodeSetOk && rightNodeSetOk {
		for _, leftNode := range leftNodeSet {
			for _, rightNode := range rightNodeSet {
				if getCursorString(leftNode) <= getCursorString(rightNode) {
					context.result = Bool(true)
					return nil
				}
			}
		}

		context.result = Bool(false)
		return nil
	}

	leftNumber, leftNumberOk := left.(Number)

	if leftNumberOk && rightNodeSetOk {
		for _, rightNode := range rightNodeSet {
			if leftNumber <= Number(getStringNumber(getCursorString(rightNode))) {
				context.result = Bool(true)
				return nil
			}
		}

		context.result = Bool(false)
		return nil
	}

	rightNumber, rightNumberOk := right.(Number)

	if leftNodeSetOk && rightNumberOk {
		for _, leftNode := range leftNodeSet {
			if Number(getStringNumber(getCursorString(leftNode))) <= rightNumber {
				context.result = Bool(true)
				return nil
			}
		}

		context.result = Bool(false)
		return nil
	}

	leftString, leftStringOk := left.(String)

	if leftStringOk && rightNodeSetOk {
		for _, rightNode := range rightNodeSet {
			if leftString <= String(getCursorString(rightNode)) {
				context.result = Bool(true)
				return nil
			}
		}

		context.result = Bool(false)
		return nil
	}

	rightString, rightStringOk := right.(String)

	if leftNodeSetOk && rightStringOk {
		for _, leftNode := range leftNodeSet {
			if String(getCursorString(leftNode)) <= rightString {
				context.result = Bool(true)
				return nil
			}
		}

		context.result = Bool(false)
		return nil
	}

	context.result = Bool(left.Number() <= right.Number())
	return nil
}

func execRelationalExprGreaterThan(context *exprContext, expr *grammar.Grammar) error {
	left, right, err := leftRightIndependentResult(context, expr)

	if err != nil {
		return err
	}

	leftNodeSet, leftNodeSetOk := left.(NodeSet)
	rightNodeSet, rightNodeSetOk := right.(NodeSet)

	if leftNodeSetOk && rightNodeSetOk {
		for _, leftNode := range leftNodeSet {
			for _, rightNode := range rightNodeSet {
				if getCursorString(leftNode) > getCursorString(rightNode) {
					context.result = Bool(true)
					return nil
				}
			}
		}

		context.result = Bool(false)
		return nil
	}

	leftNumber, leftNumberOk := left.(Number)

	if leftNumberOk && rightNodeSetOk {
		for _, rightNode := range rightNodeSet {
			if leftNumber > Number(getStringNumber(getCursorString(rightNode))) {
				context.result = Bool(true)
				return nil
			}
		}

		context.result = Bool(false)
		return nil
	}

	rightNumber, rightNumberOk := right.(Number)

	if leftNodeSetOk && rightNumberOk {
		for _, leftNode := range leftNodeSet {
			if Number(getStringNumber(getCursorString(leftNode))) > rightNumber {
				context.result = Bool(true)
				return nil
			}
		}

		context.result = Bool(false)
		return nil
	}

	leftString, leftStringOk := left.(String)

	if leftStringOk && rightNodeSetOk {
		for _, rightNode := range rightNodeSet {
			if leftString > String(getCursorString(rightNode)) {
				context.result = Bool(true)
				return nil
			}
		}

		context.result = Bool(false)
		return nil
	}

	rightString, rightStringOk := right.(String)

	if leftNodeSetOk && rightStringOk {
		for _, leftNode := range leftNodeSet {
			if String(getCursorString(leftNode)) > rightString {
				context.result = Bool(true)
				return nil
			}
		}

		context.result = Bool(false)
		return nil
	}

	context.result = Bool(left.Number() > right.Number())
	return nil
}

func execRelationalExprGreaterThanOrEqual(context *exprContext, expr *grammar.Grammar) error {
	left, right, err := leftRightIndependentResult(context, expr)

	if err != nil {
		return err
	}

	leftNodeSet, leftNodeSetOk := left.(NodeSet)
	rightNodeSet, rightNodeSetOk := right.(NodeSet)

	if leftNodeSetOk && rightNodeSetOk {
		for _, leftNode := range leftNodeSet {
			for _, rightNode := range rightNodeSet {
				if getCursorString(leftNode) >= getCursorString(rightNode) {
					context.result = Bool(true)
					return nil
				}
			}
		}

		context.result = Bool(false)
		return nil
	}

	leftNumber, leftNumberOk := left.(Number)

	if leftNumberOk && rightNodeSetOk {
		for _, rightNode := range rightNodeSet {
			if leftNumber >= Number(getStringNumber(getCursorString(rightNode))) {
				context.result = Bool(true)
				return nil
			}
		}

		context.result = Bool(false)
		return nil
	}

	rightNumber, rightNumberOk := right.(Number)

	if leftNodeSetOk && rightNumberOk {
		for _, leftNode := range leftNodeSet {
			if Number(getStringNumber(getCursorString(leftNode))) >= rightNumber {
				context.result = Bool(true)
				return nil
			}
		}

		context.result = Bool(false)
		return nil
	}

	leftString, leftStringOk := left.(String)

	if leftStringOk && rightNodeSetOk {
		for _, rightNode := range rightNodeSet {
			if leftString >= String(getCursorString(rightNode)) {
				context.result = Bool(true)
				return nil
			}
		}

		context.result = Bool(false)
		return nil
	}

	rightString, rightStringOk := right.(String)

	if leftNodeSetOk && rightStringOk {
		for _, leftNode := range leftNodeSet {
			if String(getCursorString(leftNode)) >= rightString {
				context.result = Bool(true)
				return nil
			}
		}

		context.result = Bool(false)
		return nil
	}

	context.result = Bool(left.Number() >= right.Number())
	return nil
}

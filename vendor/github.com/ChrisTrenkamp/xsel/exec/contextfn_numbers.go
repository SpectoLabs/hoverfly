package exec

import (
	"math"
	"strconv"

	"github.com/ChrisTrenkamp/xsel/grammar"
	"github.com/ChrisTrenkamp/xsel/grammar/parser/symbols"
)

func init() {
	contextFunctions[symbols.NT_Number] = execNumber
	contextFunctions[symbols.NT_AdditiveExprAdd] = execAdditiveExprAdd
	contextFunctions[symbols.NT_AdditiveExprSubtract] = execAdditiveExprSubtract
	contextFunctions[symbols.NT_MultiplicativeExprMultiply] = execMultiplicativeExprMultiply
	contextFunctions[symbols.NT_MultiplicativeExprDivide] = execMultiplicativeExprDivide
	contextFunctions[symbols.NT_MultiplicativeExprMod] = execMultiplicativeExprMod
	contextFunctions[symbols.NT_UnaryExprNegate] = execUnaryExprNegate
}

func execNumber(context *exprContext, expr *grammar.Grammar) error {
	numStr := expr.GetString()
	numResult, err := strconv.ParseFloat(numStr, 64)

	context.result = Number(numResult)
	return err
}

func execAdditiveExprAdd(context *exprContext, expr *grammar.Grammar) error {
	left, right, err := leftRightIndependentNumber(context, expr)

	if err != nil {
		return err
	}

	context.result = Number(left + right)
	return nil
}

func execAdditiveExprSubtract(context *exprContext, expr *grammar.Grammar) error {
	left, right, err := leftRightIndependentNumber(context, expr)

	if err != nil {
		return err
	}

	context.result = Number(left - right)
	return nil
}

func execMultiplicativeExprMultiply(context *exprContext, expr *grammar.Grammar) error {
	left, right, err := leftRightIndependentNumber(context, expr)

	if err != nil {
		return err
	}

	context.result = Number(left * right)
	return nil
}

func execMultiplicativeExprDivide(context *exprContext, expr *grammar.Grammar) error {
	left, right, err := leftRightIndependentNumber(context, expr)

	if err != nil {
		return err
	}

	if right == 0 {
		if left == 0 {
			context.result = Number(math.NaN())
		} else if left > 0 {
			context.result = Number(math.Inf(1))
		} else {
			context.result = Number(math.Inf(-1))
		}

		return nil
	}

	context.result = Number(left / right)
	return nil
}

func execMultiplicativeExprMod(context *exprContext, expr *grammar.Grammar) error {
	left, right, err := leftRightIndependentNumber(context, expr)

	if err != nil {
		return err
	}

	if right == 0 {
		context.result = Number(math.NaN())
		return nil
	}

	context.result = Number(int(left) % int(right))
	return nil
}

func execUnaryExprNegate(context *exprContext, expr *grammar.Grammar) error {
	left, err := leftOnlyIndependentResult(context, expr)

	if err != nil {
		return err
	}

	leftNum := left.Number()

	context.result = Number(-leftNum)
	return nil
}

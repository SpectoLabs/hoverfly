package jsonpath

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var expressionTests = []lexTest{
	{"empty", "", []int{exprEOF}},
	{"spaces", "     \t\r\n", []int{exprEOF}},
	{"numbers", " 1.3e10 ", []int{exprNumber, exprEOF}},
	// {"numbers with signs", "+1 -2.23", []int{exprNumber, exprOpPlus, exprNumber, exprEOF}},
	{"paths", " @.aKey[2].bKey ", []int{exprPath, exprEOF}},
	{"addition with mixed sign", "4+-19", []int{exprNumber, exprOpPlus, exprOpMinusUn, exprNumber, exprEOF}},
	{"addition", "4+19", []int{exprNumber, exprOpPlus, exprNumber, exprEOF}},
	{"subtraction", "4-19", []int{exprNumber, exprOpMinus, exprNumber, exprEOF}},

	{"parens", "( () + () )", []int{exprParenLeft, exprParenLeft, exprParenRight, exprOpPlus, exprParenLeft, exprParenRight, exprParenRight, exprEOF}},
	{"equals", "true ==", []int{exprBool, exprOpEq, exprEOF}},
	{"numerical comparisons", "3.4 <", []int{exprNumber, exprOpLt, exprEOF}},
}

func TestExpressionTokens(t *testing.T) {
	as := assert.New(t)
	for _, test := range expressionTests {
		lexer := NewSliceLexer([]byte(test.input), EXPRESSION)
		items := readerToArray(lexer)
		types := itemsToTypes(items)

		for _, i := range items {
			if i.typ == exprError {
				fmt.Println(string(i.val))
			}
		}

		as.EqualValues(types, test.tokenTypes, "Testing of %s: \nactual\n\t%+v\nexpected\n\t%v", test.name, typesDescription(types, exprTokenNames), typesDescription(test.tokenTypes, exprTokenNames))
	}
}

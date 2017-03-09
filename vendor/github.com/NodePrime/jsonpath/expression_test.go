package jsonpath

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var exprTests = []struct {
	input         string
	fields        map[string]Item
	expectedValue interface{}
}{
	// &&
	{"true && true", nil, true},
	{"false && true", nil, false},
	{"false && false", nil, false},

	// ||
	{"true || true", nil, true},
	{"true || false", nil, true},
	{"false ||  false", nil, false},

	// LT
	{"10 < 20", nil, true},
	{"10 < 10", nil, false},
	{"100 < 20", nil, false},
	{"@a < 50", map[string]Item{"@a": genValue(`49`, jsonNumber)}, true},
	{"@a < 50", map[string]Item{"@a": genValue(`50`, jsonNumber)}, false},
	{"@a < 50", map[string]Item{"@a": genValue(`51`, jsonNumber)}, false},

	// LE
	{"10 <= 20", nil, true},
	{"10 <= 10", nil, true},
	{"100 <= 20", nil, false},
	{"@a <= 54", map[string]Item{"@a": genValue(`53`, jsonNumber)}, true},
	{"@a <= 54", map[string]Item{"@a": genValue(`54`, jsonNumber)}, true},
	{"@a <= 54", map[string]Item{"@a": genValue(`55`, jsonNumber)}, false},

	// GT
	{"30 > 20", nil, true},
	{"20 > 20", nil, false},
	{"10 > 20", nil, false},
	{"@a > 50", map[string]Item{"@a": genValue(`49`, jsonNumber)}, false},
	{"@a > 50", map[string]Item{"@a": genValue(`50`, jsonNumber)}, false},
	{"@a > 50", map[string]Item{"@a": genValue(`51`, jsonNumber)}, true},

	// GE
	{"30 >= 20", nil, true},
	{"20 >= 20", nil, true},
	{"10 >= 20", nil, false},
	{"@a >= 50", map[string]Item{"@a": genValue(`49`, jsonNumber)}, false},
	{"@a >= 50", map[string]Item{"@a": genValue(`50`, jsonNumber)}, true},
	{"@a >= 50", map[string]Item{"@a": genValue(`51`, jsonNumber)}, true},

	// EQ
	{"20 == 20", nil, true},
	{"20 == 21", nil, false},
	{"true == true", nil, true},
	{"true == false", nil, false},
	{"@a == @b", map[string]Item{"@a": genValue(`"one"`, jsonString), "@b": genValue(`"one"`, jsonString)}, true},
	{"@a == @b", map[string]Item{"@a": genValue(`"one"`, jsonString), "@b": genValue(`"two"`, jsonString)}, false},
	{`"fire" == "fire"`, nil, true},
	{`"fire" == "water"`, nil, false},
	{`@a == "toronto"`, map[string]Item{"@a": genValue(`"toronto"`, jsonString)}, true},
	{`@a == "toronto"`, map[string]Item{"@a": genValue(`"los angeles"`, jsonString)}, false},
	{`@a == 3.4`, map[string]Item{"@a": genValue(`3.4`, jsonNumber)}, true},
	{`@a == 3.4`, map[string]Item{"@a": genValue(`3.41`, jsonNumber)}, false},
	{`@a == null`, map[string]Item{"@a": genValue(`null`, jsonNull)}, true},

	// NEQ
	{"20 != 20", nil, false},
	{"20 != 21", nil, true},
	{"true != true", nil, false},
	{"true != false", nil, true},
	{"@a != @b", map[string]Item{"@a": genValue(`"one"`, jsonString), "@b": genValue(`"one"`, jsonString)}, false},
	{"@a != @b", map[string]Item{"@a": genValue(`"one"`, jsonString), "@b": genValue(`"two"`, jsonString)}, true},
	{`"fire" != "fire"`, nil, false},
	{`"fire" != "water"`, nil, true},
	{`@a != "toronto"`, map[string]Item{"@a": genValue(`"toronto"`, jsonString)}, false},
	{`@a != "toronto"`, map[string]Item{"@a": genValue(`"los angeles"`, jsonString)}, true},
	{`@a != 3.4`, map[string]Item{"@a": genValue(`3.4`, jsonNumber)}, false},
	{`@a != 3.4`, map[string]Item{"@a": genValue(`3.41`, jsonNumber)}, true},
	{`@a != null`, map[string]Item{"@a": genValue(`null`, jsonNull)}, false},

	// Plus
	{"20 + 7", nil, 27},
	{"20 + 6.999999", nil, 26.999999},

	// Minus
	{"20 - 7", nil, 13},
	{"20 - 7.11111", nil, 12.88889},

	// Minus Unary
	{"-27", nil, -27},
	{"30 - -3", nil, 33},
	{"30 + -3", nil, 27},

	{"2 +++++ 3", nil, 5},
	{"2+--3", nil, 5},
	// Star
	{"20 * 7", nil, 140},
	{"20 * 6.999999", nil, 139.99998},
	{"20 * -7", nil, -140},
	{"-20 * -7", nil, 140},

	// Slash
	{"20 / 5", nil, 4},
	{"20 / 6.999999 - 2.85714326531 <= 0.00000001", nil, true},

	// Hat
	{"7 ^ 4", nil, 2401},
	{"2 ^ -2", nil, 0.25},
	{"((7 ^ -4) - 0.00041649312) <= 0.0001", nil, true},

	// Mod
	{"7.5 % 4", nil, 3.5},
	{"2 % -2", nil, 0},
	{"11 % 22", nil, 11},

	// Negate
	{"!true", nil, false},
	{"!false", nil, true},

	// Mix
	{"20 >= 20 || 2 == 2", nil, true},
	{"20 > @.test && @.test < 13 && @.test > 1.99994", map[string]Item{"@.test": genValue(`10.23423`, jsonNumber)}, true},
	{"20 > @.test && @.test < 13 && @.test > 1.99994", map[string]Item{"@.test": genValue(`15.3423`, jsonNumber)}, false},
}

func genValue(val string, typ int) Item {
	return Item{
		val: []byte(val),
		typ: typ,
	}
}

func TestExpressions(t *testing.T) {
	as := assert.New(t)
	emptyFields := map[string]Item{}

	for _, test := range exprTests {
		if test.fields == nil {
			test.fields = emptyFields
		}

		lexer := NewSliceLexer([]byte(test.input), EXPRESSION)
		items := readerToArray(lexer)
		// trim EOF
		items = items[0 : len(items)-1]
		items_post, err := infixToPostFix(items)
		if as.NoError(err, "Could not transform to postfix\nTest: %q", test.input) {
			val, err := evaluatePostFix(items_post, test.fields)
			if as.NoError(err, "Could not evaluate postfix\nTest Input: %q\nTest Values:%q\nError:%q", test.input, test.fields, err) {
				as.EqualValues(test.expectedValue, val, "\nTest: %q\nActual: %v \nExpected %v\n", test.input, val, test.expectedValue)
			}
		}
	}
}

var exprErrorTests = []struct {
	input                  string
	fields                 map[string]Item
	expectedErrorSubstring string
}{
	{"@a == @b", map[string]Item{"@a": genValue(`"one"`, jsonString), "@b": genValue("3.4", jsonNumber)}, "cannot be compared"},
	{")(", nil, "Mismatched parentheses"},
	{")123", nil, "Mismatched parentheses"},
	{"20 == null", nil, "cannot be compared"},
	{`"toronto" == null`, nil, "cannot be compared"},
	{`false == 20`, nil, "cannot be compared"},
	{`"nick" == 20`, nil, "cannot be compared"},
	{"20 != null", nil, "cannot be compared"},
	{`"toronto" != null`, nil, "cannot be compared"},
	{`false != 20`, nil, "cannot be compared"},
	{`"nick" != 20`, nil, "cannot be compared"},
	{``, nil, "Bad Expression"},
	{`==`, nil, "Bad Expression"},
	{`!=`, nil, "Not enough operands"},

	{`!23`, nil, "cannot be compared"},
	{`"nick" || true`, nil, "cannot be compared"},
	{`"nick" >=  3.2`, nil, "cannot be compared"},
	{`"nick" >3.2`, nil, "cannot be compared"},
	{`"nick" <=  3.2`, nil, "cannot be compared"},
	{`"nick" <  3.2`, nil, "cannot be compared"},
	{`"nick" +  3.2`, nil, "cannot be compared"},
	{`"nick" -  3.2`, nil, "cannot be compared"},
	{`"nick" /  3.2`, nil, "cannot be compared"},
	{`"nick" *  3.2`, nil, "cannot be compared"},
	{`"nick" %  3.2`, nil, "cannot be compared"},
	{`"nick"+`, nil, "cannot be compared"},
	{`"nick"-`, nil, "cannot be compared"},
	{`"nick"^3.2`, nil, "cannot be compared"},

	{`@a == null`, map[string]Item{"@a": genValue(`3.41`, jsonNumber)}, "cannot be compared"},
}

func TestBadExpressions(t *testing.T) {
	as := assert.New(t)
	emptyFields := map[string]Item{}

	for _, test := range exprErrorTests {
		if test.fields == nil {
			test.fields = emptyFields
		}

		lexer := NewSliceLexer([]byte(test.input), EXPRESSION)
		items := readerToArray(lexer)
		// trim EOF
		items = items[0 : len(items)-1]
		items_post, err := infixToPostFix(items)
		if err != nil {
			as.True(strings.Contains(err.Error(), test.expectedErrorSubstring), "Test Input: %q\nError %q does not contain %q", test.input, err.Error(), test.expectedErrorSubstring)
			continue
		}
		if as.NoError(err, "Could not transform to postfix\nTest: %q", test.input) {
			_, err := evaluatePostFix(items_post, test.fields)
			if as.Error(err, "Could not evaluate postfix\nTest Input: %q\nTest Values:%q\nError:%s", test.input, test.fields, err) {
				as.True(strings.Contains(err.Error(), test.expectedErrorSubstring), "Test Input: %q\nError %s does not contain %q", test.input, err.Error(), test.expectedErrorSubstring)
			}

		}
	}
}

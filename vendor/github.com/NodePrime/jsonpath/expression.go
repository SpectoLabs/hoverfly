package jsonpath

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
)

const (
	exprErrorMismatchedParens   = "Mismatched parentheses"
	exprErrorBadExpression      = "Bad Expression"
	exprErrorFinalValueNotBool  = "Expression evaluated to a non-bool: %v"
	exprErrorNotEnoughOperands  = "Not enough operands for operation %q"
	exprErrorValueNotFound      = "Value for %q not found"
	exprErrorBadValue           = "Bad value %q for type %q"
	exprErrorPathValueNotScalar = "Path value must be scalar value"
	exprErrorBadOperandType     = "Operand type expected to be %q for operation %q"
)

type exprErrorBadTypeComparison struct {
	valueType    string
	expectedType string
}

func (e exprErrorBadTypeComparison) Error() string {
	return fmt.Sprintf("Type %s cannot be compared to type %s", e.valueType, e.expectedType)
}

// Lowest priority = lowest #
var opa = map[int]struct {
	prec   int
	rAssoc bool
}{
	exprOpAnd:     {1, false},
	exprOpOr:      {1, false},
	exprOpEq:      {2, false},
	exprOpNeq:     {2, false},
	exprOpLt:      {3, false},
	exprOpLe:      {3, false},
	exprOpGt:      {3, false},
	exprOpGe:      {3, false},
	exprOpPlus:    {4, false},
	exprOpMinus:   {4, false},
	exprOpSlash:   {5, false},
	exprOpStar:    {5, false},
	exprOpPercent: {5, false},
	exprOpHat:     {6, false},
	exprOpNot:     {7, true},
	exprOpPlusUn:  {7, true},
	exprOpMinusUn: {7, true},
}

// Shunting-yard Algorithm (infix -> postfix)
// http://rosettacode.org/wiki/Parsing/Shunting-yard_algorithm#Go
func infixToPostFix(items []Item) (out []Item, err error) {
	stack := newStack()

	for _, i := range items {
		switch i.typ {
		case exprParenLeft:
			stack.push(i) // push "(" to stack
		case exprParenRight:
			found := false
			for {
				// pop item ("(" or operator) from stack
				op_interface, ok := stack.pop()
				if !ok {
					return nil, errors.New(exprErrorMismatchedParens)
				}
				op := op_interface.(Item)
				if op.typ == exprParenLeft {
					found = true
					break // discard "("
				}
				out = append(out, op) // add operator to result
			}
			if !found {
				return nil, errors.New(exprErrorMismatchedParens)
			}
		default:
			if o1, isOp := opa[i.typ]; isOp {
				// token is an operator
				for stack.len() > 0 {
					// consider top item on stack
					op_int, _ := stack.peek()
					op := op_int.(Item)
					if o2, isOp := opa[op.typ]; !isOp || o1.prec > o2.prec ||
						o1.prec == o2.prec && o1.rAssoc {
						break
					}
					// top item is an operator that needs to come off
					stack.pop()           // pop it
					out = append(out, op) // add it to result
				}
				// push operator (the new one) to stack
				stack.push(i)
			} else { // token is an operand
				out = append(out, i) // add operand to result
			}
		}
	}
	// drain stack to result
	for stack.len() > 0 {
		op_int, _ := stack.pop()
		op := op_int.(Item)
		if op.typ == exprParenLeft {
			return nil, errors.New(exprErrorMismatchedParens)
		}
		out = append(out, op)
	}
	return
}

func evaluatePostFix(postFixItems []Item, pathValues map[string]Item) (interface{}, error) {
	s := newStack()

	if len(postFixItems) == 0 {
		return false, errors.New(exprErrorBadExpression)
	}

	for _, item := range postFixItems {
		switch item.typ {

		// VALUES
		case exprBool:
			val, err := strconv.ParseBool(string(item.val))
			if err != nil {
				return false, fmt.Errorf(exprErrorBadValue, string(item.val), exprTokenNames[exprBool])
			}
			s.push(val)
		case exprNumber:
			val, err := strconv.ParseFloat(string(item.val), 64)
			if err != nil {
				return false, fmt.Errorf(exprErrorBadValue, string(item.val), exprTokenNames[exprNumber])
			}
			s.push(val)
		case exprPath:
			// TODO: Handle datatypes of JSON
			i, ok := pathValues[string(item.val)]
			if !ok {
				return false, fmt.Errorf(exprErrorValueNotFound, string(item.val))
			}
			switch i.typ {
			case jsonNull:
				s.push(nil)
			case jsonNumber:
				val_float, err := strconv.ParseFloat(string(i.val), 64)
				if err != nil {
					return false, fmt.Errorf(exprErrorBadValue, string(item.val), jsonTokenNames[jsonNumber])
				}
				s.push(val_float)
			case jsonKey, jsonString:
				s.push(i.val)
			default:
				return false, fmt.Errorf(exprErrorPathValueNotScalar)
			}
		case exprString:
			s.push(item.val)
		case exprNull:
			s.push(nil)

		// OPERATORS
		case exprOpAnd:
			a, b, err := take2Bool(s, item.typ)
			if err != nil {
				return false, err
			}

			s.push(a && b)
		case exprOpEq:
			p, ok := s.peek()
			if !ok {
				return false, fmt.Errorf(exprErrorNotEnoughOperands, exprTokenNames[item.typ])
			}
			switch p.(type) {
			case nil:
				err := take2Null(s, item.typ)
				if err != nil {
					return false, err
				} else {
					s.push(true)
				}
			case bool:
				a, b, err := take2Bool(s, item.typ)
				if err != nil {
					return false, err
				}
				s.push(a == b)
			case float64:
				a, b, err := take2Float(s, item.typ)
				if err != nil {
					return false, err
				}
				s.push(a == b)
			case []byte:
				a, b, err := take2ByteSlice(s, item.typ)
				if err != nil {
					return false, err
				}
				s.push(byteSlicesEqual(a, b))
			}
		case exprOpNeq:
			p, ok := s.peek()
			if !ok {
				return false, fmt.Errorf(exprErrorNotEnoughOperands, exprTokenNames[item.typ])
			}
			switch p.(type) {
			case nil:
				err := take2Null(s, item.typ)
				if err != nil {
					return true, err
				} else {
					s.push(false)
				}
			case bool:
				a, b, err := take2Bool(s, item.typ)
				if err != nil {
					return false, err
				}
				s.push(a != b)
			case float64:
				a, b, err := take2Float(s, item.typ)
				if err != nil {
					return false, err
				}
				s.push(a != b)
			case []byte:
				a, b, err := take2ByteSlice(s, item.typ)
				if err != nil {
					return false, err
				}
				s.push(!byteSlicesEqual(a, b))
			}
		case exprOpNot:
			a, err := take1Bool(s, item.typ)
			if err != nil {
				return false, err
			}

			s.push(!a)
		case exprOpOr:
			a, b, err := take2Bool(s, item.typ)
			if err != nil {
				return false, err
			}

			s.push(a || b)
		case exprOpGt:
			a, b, err := take2Float(s, item.typ)
			if err != nil {
				return false, err
			}

			s.push(b > a)
		case exprOpGe:
			a, b, err := take2Float(s, item.typ)
			if err != nil {
				return false, err
			}

			s.push(b >= a)
		case exprOpLt:
			a, b, err := take2Float(s, item.typ)
			if err != nil {
				return false, err
			}

			s.push(b < a)
		case exprOpLe:
			a, b, err := take2Float(s, item.typ)
			if err != nil {
				return false, err
			}

			s.push(b <= a)
		case exprOpPlus:
			a, b, err := take2Float(s, item.typ)
			if err != nil {
				return false, err
			}
			s.push(b + a)
		case exprOpPlusUn:
			a, err := take1Float(s, item.typ)
			if err != nil {
				return false, err
			}
			s.push(a)
		case exprOpMinus:
			a, b, err := take2Float(s, item.typ)
			if err != nil {
				return false, err
			}
			s.push(b - a)
		case exprOpMinusUn:
			a, err := take1Float(s, item.typ)
			if err != nil {
				return false, err
			}
			s.push(0 - a)
		case exprOpSlash:
			a, b, err := take2Float(s, item.typ)
			if err != nil {
				return false, err
			}

			if a == 0.0 {
				return false, errors.New("Cannot divide by zero")
			}
			s.push(b / a)
		case exprOpStar:
			a, b, err := take2Float(s, item.typ)
			if err != nil {
				return false, err
			}

			s.push(b * a)
		case exprOpPercent:
			a, b, err := take2Float(s, item.typ)
			if err != nil {
				return false, err
			}

			s.push(math.Mod(b, a))
		case exprOpHat:
			a, b, err := take2Float(s, item.typ)
			if err != nil {
				return false, err
			}

			s.push(math.Pow(b, a))
		case exprOpExclam:
			a, err := take1Bool(s, item.typ)
			if err != nil {
				return false, err
			}
			s.push(!a)
		// Other
		default:
			return false, fmt.Errorf("Token not supported in evaluator: %v", exprTokenNames[item.typ])
		}
	}

	if s.len() != 1 {
		return false, fmt.Errorf(exprErrorBadExpression)
	}
	end_int, _ := s.pop()
	return end_int, nil
}

func take1Bool(s *stack, op int) (bool, error) {
	t := exprBool
	val, ok := s.pop()
	if !ok {
		return false, fmt.Errorf(exprErrorNotEnoughOperands, exprTokenNames[op])
	}

	b, ok := val.(bool)
	if !ok {
		return false, exprErrorBadTypeComparison{exprTokenNames[t], (reflect.TypeOf(val)).String()}
	}
	return b, nil
}

func take2Bool(s *stack, op int) (bool, bool, error) {
	a, a_err := take1Bool(s, op)
	b, b_err := take1Bool(s, op)
	return a, b, firstError(a_err, b_err)
}

func take1Float(s *stack, op int) (float64, error) {
	t := exprNumber
	val, ok := s.pop()
	if !ok {
		return 0.0, fmt.Errorf(exprErrorNotEnoughOperands, exprTokenNames[op])
	}

	b, ok := val.(float64)
	if !ok {
		return 0.0, exprErrorBadTypeComparison{exprTokenNames[t], (reflect.TypeOf(val)).String()}
	}
	return b, nil
}

func take2Float(s *stack, op int) (float64, float64, error) {
	a, a_err := take1Float(s, op)
	b, b_err := take1Float(s, op)
	return a, b, firstError(a_err, b_err)
}

func take1ByteSlice(s *stack, op int) ([]byte, error) {
	t := exprNumber
	val, ok := s.pop()
	if !ok {
		return nil, fmt.Errorf(exprErrorNotEnoughOperands, exprTokenNames[op])
	}

	b, ok := val.([]byte)
	if !ok {
		return nil, exprErrorBadTypeComparison{exprTokenNames[t], (reflect.TypeOf(val)).String()}
	}
	return b, nil
}

func take2ByteSlice(s *stack, op int) ([]byte, []byte, error) {
	a, a_err := take1ByteSlice(s, op)
	b, b_err := take1ByteSlice(s, op)
	return a, b, firstError(a_err, b_err)
}

func take1Null(s *stack, op int) error {
	t := exprNull
	val, ok := s.pop()
	if !ok {
		return fmt.Errorf(exprErrorNotEnoughOperands, exprTokenNames[op])
	}

	if v := reflect.TypeOf(val); v != nil {
		return exprErrorBadTypeComparison{exprTokenNames[t], v.String()}
	}
	return nil
}

func take2Null(s *stack, op int) error {
	a_err := take1Null(s, op)
	b_err := take1Null(s, op)
	return firstError(a_err, b_err)
}

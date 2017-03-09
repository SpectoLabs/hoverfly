package jsonpath

const (
	exprError = iota
	exprEOF
	exprParenLeft
	exprParenRight
	exprNumber
	exprPath
	exprBool
	exprNull
	exprString

	exprOperators
	exprOpEq
	exprOpNeq
	exprOpNot
	exprOpLt
	exprOpLe
	exprOpGt
	exprOpGe
	exprOpAnd
	exprOpOr
	exprOpPlus
	exprOpPlusUn
	exprOpMinus
	exprOpMinusUn
	exprOpSlash
	exprOpStar
	exprOpHat
	exprOpPercent
	exprOpExclam
)

var exprTokenNames = map[int]string{
	exprError: "error",
	exprEOF:   "EOF",

	exprParenLeft:  "(",
	exprParenRight: ")",
	exprNumber:     "number",
	exprPath:       "path",
	exprBool:       "bool",
	exprNull:       "null",
	exprString:     "string",
	exprOpEq:       "==",
	exprOpNeq:      "!=",
	exprOpNot:      "!",
	exprOpLt:       "<",
	exprOpLe:       "<=",
	exprOpGt:       ">",
	exprOpGe:       ">=",
	exprOpAnd:      "&&",
	exprOpOr:       "||",
	exprOpPlus:     "+",
	exprOpPlusUn:   "(+)",
	exprOpMinus:    "-",
	exprOpMinusUn:  "(-)",
	exprOpSlash:    "/",
	exprOpStar:     "*",
	exprOpHat:      "^",
	exprOpPercent:  "%",
	exprOpExclam:   "!",
}

var EXPRESSION = lexExprText

func lexExprText(l lexer, state *intStack) stateFn {
	ignoreSpaceRun(l)
	cur := l.peek()
	var next stateFn
	switch cur {
	case '(':
		l.take()
		state.push(exprParenLeft)
		l.emit(exprParenLeft)
		next = lexExprText
	case ')':
		if top, ok := state.peek(); ok && top != exprParenLeft {
			next = l.errorf("Received %#U but has no matching (", cur)
			break
		}
		state.pop()
		l.take()
		l.emit(exprParenRight)

		next = lexOneValue
	case '!':
		l.take()
		l.emit(exprOpNot)
		next = lexExprText
	case '+':
		l.take()
		l.emit(exprOpPlusUn)
		next = lexExprText
	case '-':
		l.take()
		l.emit(exprOpMinusUn)
		next = lexExprText
	case '@': //, '$': // Only support current key location
		l.take()
		takePath(l)
		l.emit(exprPath)
		next = lexOneValue
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		takeNumeric(l)
		l.emit(exprNumber)
		next = lexOneValue
	case 't':
		takeExactSequence(l, bytesTrue)
		l.emit(exprBool)
		next = lexOneValue
	case 'f':
		takeExactSequence(l, bytesFalse)
		l.emit(exprBool)
		next = lexOneValue
	case 'n':
		takeExactSequence(l, bytesNull)
		l.emit(exprNull)
		next = lexOneValue
	case '"':
		err := l.takeString()
		if err != nil {
			return l.errorf("Could not take string because %q", err)
		}
		l.emit(exprString)
		next = lexOneValue
	case eof:
		l.emit(exprEOF)
		// next = nil
	default:
		return l.errorf("Unrecognized sequence in expression: %#U", cur)
	}
	return next
}

func lexOneValue(l lexer, state *intStack) stateFn {
	var next stateFn
	cur := l.peek()
	switch cur {
	case '+':
		l.take()
		l.emit(exprOpPlus)
		next = lexExprText
	case '-':
		l.take()
		l.emit(exprOpMinus)
		next = lexExprText
	case '*':
		l.take()
		l.emit(exprOpStar)
		next = lexExprText
	case '/':
		l.take()
		l.emit(exprOpSlash)
		next = lexExprText
	case '%':
		l.take()
		l.emit(exprOpPercent)
		next = lexExprText
	case '^':
		l.take()
		l.emit(exprOpHat)
		next = lexExprText
	case '<':
		l.take()
		cur = l.peek()
		if cur == '=' {
			l.take()
			l.emit(exprOpLe)
		} else {
			l.emit(exprOpLt)
		}
		next = lexExprText
	case '>':
		l.take()
		cur = l.peek()
		if cur == '=' {
			l.take()
			l.emit(exprOpGe)
		} else {
			l.emit(exprOpGt)
		}
		next = lexExprText
	case '&':
		l.take()
		cur = l.take()
		if cur != '&' {
			return l.errorf("Expected double & instead of %#U", cur)
		}
		l.emit(exprOpAnd)
		next = lexExprText
	case '|':
		l.take()
		cur = l.take()
		if cur != '|' {
			return l.errorf("Expected double | instead of %#U", cur)
		}
		l.emit(exprOpOr)
		next = lexExprText
	case '=':
		l.take()
		cur = l.take()
		if cur != '=' {
			return l.errorf("Expected double = instead of %#U", cur)
		}
		l.emit(exprOpEq)
		next = lexExprText
	case '!':
		l.take()
		cur = l.take()
		if cur != '=' {
			return l.errorf("Expected = for != instead of %#U", cur)
		}
		l.emit(exprOpNeq)
		next = lexExprText
	case ')':
		if top, ok := state.peek(); ok && top != exprParenLeft {
			next = l.errorf("Received %#U but has no matching (", cur)
			break
		}
		state.pop()
		l.take()
		l.emit(exprParenRight)

		next = lexOneValue
	case eof:
		l.emit(exprEOF)
	default:
		return l.errorf("Unrecognized sequence in expression: %#U", cur)
	}
	return next
}

func takeNumeric(l lexer) {
	takeDigits(l)
	if l.peek() == '.' {
		l.take()
		takeDigits(l)
	}
	if l.peek() == 'e' || l.peek() == 'E' {
		l.take()
		if l.peek() == '+' || l.peek() == '-' {
			l.take()
			takeDigits(l)
		} else {
			takeDigits(l)
		}
	}
}

func takePath(l lexer) {
	inQuotes := false
	var prev int = 0
	// capture until end of path - ugly
takeLoop:
	for {
		cur := l.peek()
		switch cur {
		case '"':
			if prev != '\\' {
				inQuotes = !inQuotes
			}
			l.take()
		case ' ':
			if !inQuotes {
				break takeLoop
			}
			l.take()
		case eof:
			break takeLoop
		default:
			l.take()
		}

		prev = cur
	}
}

func lexExprEnd(l lexer, state *intStack) stateFn {
	cur := l.take()
	if cur != eof {
		return l.errorf("Expected EOF but received %#U", cur)
	}
	l.emit(exprEOF)
	return nil
}

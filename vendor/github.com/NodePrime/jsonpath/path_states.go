package jsonpath

const (
	pathError = iota
	pathEOF

	pathRoot
	pathCurrent
	pathKey
	pathBracketLeft
	pathBracketRight
	pathIndex
	pathOr
	pathIndexRange
	pathLength
	pathWildcard
	pathPeriod
	pathValue
	pathWhere
	pathExpression
)

var pathTokenNames = map[int]string{
	pathError: "ERROR",
	pathEOF:   "EOF",

	pathRoot:         "$",
	pathCurrent:      "@",
	pathKey:          "KEY",
	pathBracketLeft:  "[",
	pathBracketRight: "]",
	pathIndex:        "INDEX",
	pathOr:           "|",
	pathIndexRange:   ":",
	pathLength:       "LENGTH",
	pathWildcard:     "*",
	pathPeriod:       ".",
	pathValue:        "+",
	pathWhere:        "?",
	pathExpression:   "EXPRESSION",
}

var PATH = lexPathStart

func lexPathStart(l lexer, state *intStack) stateFn {
	ignoreSpaceRun(l)
	cur := l.take()
	switch cur {
	case '$':
		l.emit(pathRoot)
	case '@':
		l.emit(pathCurrent)
	default:
		return l.errorf("Expected $ or @ at start of path instead of  %#U", cur)
	}

	return lexPathAfterKey
}

func lexPathAfterKey(l lexer, state *intStack) stateFn {
	cur := l.take()
	switch cur {
	case '.':
		l.emit(pathPeriod)
		return lexKey
	case '[':
		l.emit(pathBracketLeft)
		return lexPathBracketOpen
	case '+':
		l.emit(pathValue)
		return lexPathAfterValue
	case '?':
		l.emit(pathWhere)
		return lexPathExpression
	case eof:
		l.emit(pathEOF)
	default:
		return l.errorf("Unrecognized rune after path element %#U", cur)
	}
	return nil
}

func lexPathExpression(l lexer, state *intStack) stateFn {
	cur := l.take()
	if cur != '(' {
		return l.errorf("Expected $ at start of path instead of  %#U", cur)
	}

	parenLeftCount := 1
	for {
		cur = l.take()
		switch cur {
		case '(':
			parenLeftCount++
		case ')':
			parenLeftCount--
		case eof:
			return l.errorf("Unexpected EOF within expression")
		}

		if parenLeftCount == 0 {
			break
		}
	}
	l.emit(pathExpression)
	return lexPathAfterKey
}

func lexPathBracketOpen(l lexer, state *intStack) stateFn {
	switch l.peek() {
	case '*':
		l.take()
		l.emit(pathWildcard)
		return lexPathBracketClose
	case '"':
		l.takeString()
		l.emit(pathKey)
		return lexPathBracketClose
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		l.take()
		takeDigits(l)
		l.emit(pathIndex)
		return lexPathIndexRange
	case eof:
		l.emit(pathEOF)
	}
	return nil
}

func lexPathBracketClose(l lexer, state *intStack) stateFn {
	cur := l.take()
	if cur != ']' {
		return l.errorf("Expected ] instead of  %#U", cur)
	}
	l.emit(pathBracketRight)
	return lexPathAfterKey
}

func lexKey(l lexer, state *intStack) stateFn {
	// TODO: Support globbing of keys
	switch l.peek() {
	case '*':
		l.take()
		l.emit(pathWildcard)
		return lexPathAfterKey
	case '"':
		l.takeString()
		l.emit(pathKey)
		return lexPathAfterKey
	case eof:
		l.take()
		l.emit(pathEOF)
		return nil
	default:
		for {
			v := l.peek()
			if v == '.' || v == '[' || v == '+' || v == '?' || v == eof {
				break
			}
			l.take()
		}
		l.emit(pathKey)
		return lexPathAfterKey
	}
}

func lexPathIndexRange(l lexer, state *intStack) stateFn {
	// TODO: Expand supported operations
	// Currently only supports single index or wildcard (1 or all)
	cur := l.peek()
	switch cur {
	case ':':
		l.take()
		l.emit(pathIndexRange)
		return lexPathIndexRangeSecond
	case ']':
		return lexPathBracketClose
	default:
		return l.errorf("Expected digit or ] instead of  %#U", cur)
	}
}

func lexPathIndexRangeSecond(l lexer, state *intStack) stateFn {
	cur := l.peek()
	switch cur {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		takeDigits(l)
		l.emit(pathIndex)
		return lexPathBracketClose
	case ']':
		return lexPathBracketClose
	default:
		return l.errorf("Expected digit or ] instead of  %#U", cur)
	}
}

func lexPathArrayClose(l lexer, state *intStack) stateFn {
	cur := l.take()
	if cur != ']' {
		return l.errorf("Expected ] instead of  %#U", cur)
	}
	l.emit(pathBracketRight)
	return lexPathAfterKey
}

func lexPathAfterValue(l lexer, state *intStack) stateFn {
	cur := l.take()
	if cur != eof {
		return l.errorf("Expected EOF instead of %#U", cur)
	}
	l.emit(pathEOF)
	return nil
}

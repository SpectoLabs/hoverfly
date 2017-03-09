package jsonpath

const (
	jsonError = iota
	jsonEOF

	jsonBraceLeft
	jsonBraceRight
	jsonBracketLeft
	jsonBracketRight
	jsonColon
	jsonComma
	jsonNumber
	jsonString
	jsonNull
	jsonKey
	jsonBool
)

var trueBytes = []byte{'t', 'r', 'u', 'e'}
var falseBytes = []byte{'f', 'a', 'l', 's', 'e'}
var nullBytes = []byte{'n', 'u', 'l', 'l'}

var jsonTokenNames = map[int]string{
	jsonEOF:   "EOF",
	jsonError: "ERROR",

	jsonBraceLeft:    "{",
	jsonBraceRight:   "}",
	jsonBracketLeft:  "[",
	jsonBracketRight: "]",
	jsonColon:        ":",
	jsonComma:        ",",
	jsonNumber:       "NUMBER",
	jsonString:       "STRING",
	jsonNull:         "NULL",
	jsonKey:          "KEY",
	jsonBool:         "BOOL",
}

var JSON = lexJsonRoot

func lexJsonRoot(l lexer, state *intStack) stateFn {
	ignoreSpaceRun(l)
	cur := l.peek()
	var next stateFn
	switch cur {
	case '{':
		next = stateJsonObjectOpen
	case '[':
		next = stateJsonArrayOpen
	default:
		next = l.errorf("Expected '{' or '[' at root of JSON instead of %#U", cur)
	}
	return next
}

func stateJsonObjectOpen(l lexer, state *intStack) stateFn {
	cur := l.take()
	if cur != '{' {
		return l.errorf("Expected '{' as start of object instead of %#U", cur)
	}
	l.emit(jsonBraceLeft)
	state.push(jsonBraceLeft)

	return stateJsonObject
}

func stateJsonArrayOpen(l lexer, state *intStack) stateFn {
	cur := l.take()
	if cur != '[' {
		return l.errorf("Expected '[' as start of array instead of %#U", cur)
	}
	l.emit(jsonBracketLeft)
	state.push(jsonBracketLeft)

	return stateJsonArray
}

func stateJsonObject(l lexer, state *intStack) stateFn {
	var next stateFn
	cur := l.peek()
	switch cur {
	case '}':
		if top, ok := state.peek(); ok && top != jsonBraceLeft {
			next = l.errorf("Received %#U but has no matching '{'", cur)
			break
		}
		l.take()
		l.emit(jsonBraceRight)
		state.pop()
		next = stateJsonAfterValue
	case '"':
		next = stateJsonKey
	default:
		next = l.errorf("Expected '}' or \" within an object instead of %#U", cur)
	}
	return next
}

func stateJsonArray(l lexer, state *intStack) stateFn {
	var next stateFn
	cur := l.peek()
	switch cur {
	case ']':
		if top, ok := state.peek(); ok && top != jsonBracketLeft {
			next = l.errorf("Received %#U but has no matching '['", cur)
			break
		}
		l.take()
		l.emit(jsonBracketRight)
		state.pop()
		next = stateJsonAfterValue
	default:
		next = stateJsonValue
	}
	return next
}

func stateJsonAfterValue(l lexer, state *intStack) stateFn {
	cur := l.take()
	top, ok := state.peek()
	topVal := noValue
	if ok {
		topVal = top
	}

	switch cur {
	case ',':
		l.emit(jsonComma)
		switch topVal {
		case jsonBraceLeft:
			return stateJsonKey
		case jsonBracketLeft:
			return stateJsonValue
		case noValue:
			return l.errorf("Found %#U outside of array or object", cur)
		default:
			return l.errorf("Unexpected character in lexer stack: %#U", cur)
		}
	case '}':
		l.emit(jsonBraceRight)
		state.pop()
		switch topVal {
		case jsonBraceLeft:
			return stateJsonAfterValue
		case jsonBracketLeft:
			return l.errorf("Unexpected %#U in array", cur)
		case noValue:
			return stateJsonAfterRoot
		}
	case ']':
		l.emit(jsonBracketRight)
		state.pop()
		switch topVal {
		case jsonBraceLeft:
			return l.errorf("Unexpected %#U in object", cur)
		case jsonBracketLeft:
			return stateJsonAfterValue
		case noValue:
			return stateJsonAfterRoot
		}
	case eof:
		if state.len() == 0 {
			l.emit(jsonEOF)
			return nil
		} else {
			return l.errorf("Unexpected EOF instead of value")
		}
	default:
		return l.errorf("Unexpected character after json value token: %#U", cur)
	}
	return nil
}

func stateJsonKey(l lexer, state *intStack) stateFn {
	if err := l.takeString(); err != nil {
		return l.errorf(err.Error())
	}
	l.emit(jsonKey)

	return stateJsonColon
}

func stateJsonColon(l lexer, state *intStack) stateFn {
	cur := l.take()
	if cur != ':' {
		return l.errorf("Expected ':' after key string instead of %#U", cur)
	}
	l.emit(jsonColon)

	return stateJsonValue
}

func stateJsonValue(l lexer, state *intStack) stateFn {
	cur := l.peek()

	switch cur {
	case eof:
		return l.errorf("Unexpected EOF instead of value")
	case '"':
		return stateJsonString
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return stateJsonNumber
	case 't', 'f':
		return stateJsonBool
	case 'n':
		return stateJsonNull
	case '{':
		return stateJsonObjectOpen
	case '[':
		return stateJsonArrayOpen
	default:
		return l.errorf("Unexpected character as start of value: %#U", cur)
	}
}

func stateJsonString(l lexer, state *intStack) stateFn {
	if err := l.takeString(); err != nil {
		return l.errorf(err.Error())
	}
	l.emit(jsonString)
	return stateJsonAfterValue
}

func stateJsonNumber(l lexer, state *intStack) stateFn {
	if err := takeJSONNumeric(l); err != nil {
		return l.errorf(err.Error())
	}
	l.emit(jsonNumber)
	return stateJsonAfterValue
}

func stateJsonBool(l lexer, state *intStack) stateFn {
	cur := l.peek()
	var match []byte
	switch cur {
	case 't':
		match = trueBytes
	case 'f':
		match = falseBytes
	}

	if !takeExactSequence(l, match) {
		return l.errorf("Could not parse value as 'true' or 'false'")
	}
	l.emit(jsonBool)
	return stateJsonAfterValue
}

func stateJsonNull(l lexer, state *intStack) stateFn {
	if !takeExactSequence(l, nullBytes) {
		return l.errorf("Could not parse value as 'null'")
	}
	l.emit(jsonNull)
	return stateJsonAfterValue
}

func stateJsonAfterRoot(l lexer, state *intStack) stateFn {
	cur := l.take()
	if cur != eof {
		return l.errorf("Expected EOF instead of %#U", cur)
	}
	l.emit(jsonEOF)
	return nil
}

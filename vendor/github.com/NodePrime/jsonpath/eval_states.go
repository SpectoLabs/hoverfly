package jsonpath

import (
	"errors"
	"fmt"
)

func evalRoot(e *Eval, i *Item) evalStateFn {
	switch i.typ {
	case jsonBraceLeft:
		e.levelStack.push(i.typ)
		return evalObjectAfterOpen
	case jsonBracketLeft:
		e.levelStack.push(i.typ)
		return evalArrayAfterOpen
	case jsonError:
		return evalError(e, i)
	default:
		e.Error = errors.New(UnexpectedToken)
	}
	return nil
}

func evalObjectAfterOpen(e *Eval, i *Item) evalStateFn {
	switch i.typ {
	case jsonKey:
		c := i.val[1 : len(i.val)-1]
		if e.copyValues {
			d := make([]byte, len(c))
			copy(d, c)
			c = d
		}
		e.nextKey = c
		return evalObjectColon
	case jsonBraceRight:
		return rightBraceOrBracket(e)
	case jsonError:
		return evalError(e, i)
	default:
		e.Error = errors.New(UnexpectedToken)
	}
	return nil
}

func evalObjectColon(e *Eval, i *Item) evalStateFn {
	switch i.typ {
	case jsonColon:
		return evalObjectValue
	case jsonError:
		return evalError(e, i)
	default:
		e.Error = errors.New(UnexpectedToken)
	}

	return nil
}

func evalObjectValue(e *Eval, i *Item) evalStateFn {
	e.location.push(e.nextKey)

	switch i.typ {
	case jsonNull, jsonNumber, jsonString, jsonBool:
		return evalObjectAfterValue
	case jsonBraceLeft:
		e.levelStack.push(i.typ)
		return evalObjectAfterOpen
	case jsonBracketLeft:
		e.levelStack.push(i.typ)
		return evalArrayAfterOpen
	case jsonError:
		return evalError(e, i)
	default:
		e.Error = errors.New(UnexpectedToken)
	}
	return nil
}

func evalObjectAfterValue(e *Eval, i *Item) evalStateFn {
	e.location.pop()
	switch i.typ {
	case jsonComma:
		return evalObjectAfterOpen
	case jsonBraceRight:
		return rightBraceOrBracket(e)
	case jsonError:
		return evalError(e, i)
	default:
		e.Error = errors.New(UnexpectedToken)
	}
	return nil
}

func rightBraceOrBracket(e *Eval) evalStateFn {
	e.levelStack.pop()

	lowerTyp, ok := e.levelStack.peek()
	if !ok {
		return evalRootEnd
	} else {
		switch lowerTyp {
		case jsonBraceLeft:
			return evalObjectAfterValue
		case jsonBracketLeft:
			return evalArrayAfterValue
		}
	}
	return nil
}

func evalArrayAfterOpen(e *Eval, i *Item) evalStateFn {
	e.prevIndex = -1

	switch i.typ {
	case jsonNull, jsonNumber, jsonString, jsonBool, jsonBraceLeft, jsonBracketLeft:
		return evalArrayValue(e, i)
	case jsonBracketRight:
		setPrevIndex(e)
		return rightBraceOrBracket(e)
	case jsonError:
		return evalError(e, i)
	default:
		e.Error = errors.New(UnexpectedToken)
	}
	return nil
}

func evalArrayValue(e *Eval, i *Item) evalStateFn {
	e.prevIndex++
	e.location.push(e.prevIndex)

	switch i.typ {
	case jsonNull, jsonNumber, jsonString, jsonBool:
		return evalArrayAfterValue
	case jsonBraceLeft:
		e.levelStack.push(i.typ)
		return evalObjectAfterOpen
	case jsonBracketLeft:
		e.levelStack.push(i.typ)
		return evalArrayAfterOpen
	case jsonError:
		return evalError(e, i)
	default:
		e.Error = errors.New(UnexpectedToken)
	}
	return nil
}

func evalArrayAfterValue(e *Eval, i *Item) evalStateFn {
	switch i.typ {
	case jsonComma:
		if val, ok := e.location.pop(); ok {
			if valIndex, ok := val.(int); ok {
				e.prevIndex = valIndex
			}
		}
		return evalArrayValue
	case jsonBracketRight:
		e.location.pop()
		setPrevIndex(e)
		return rightBraceOrBracket(e)
	case jsonError:
		return evalError(e, i)
	default:
		e.Error = errors.New(UnexpectedToken)
	}
	return nil
}

func setPrevIndex(e *Eval) {
	e.prevIndex = -1
	peeked, ok := e.location.peek()
	if ok {
		if peekedIndex, intOk := peeked.(int); intOk {
			e.prevIndex = peekedIndex
		}
	}
}

func evalRootEnd(e *Eval, i *Item) evalStateFn {
	if i.typ != jsonEOF {
		if i.typ == jsonError {
			evalError(e, i)
		} else {
			e.Error = errors.New(BadStructure)
		}
	}
	return nil
}

func evalError(e *Eval, i *Item) evalStateFn {
	e.Error = fmt.Errorf("%s at byte index %d", string(i.val), i.pos)
	return nil
}

package jsonpath

import (
	"errors"
	"fmt"
)

func takeExponent(l lexer) error {
	r := l.peek()
	if r != 'e' && r != 'E' {
		return nil
	}
	l.take()
	r = l.take()
	switch r {
	case '+', '-':
		// Check digit immediately follows sign
		if d := l.peek(); !(d >= '0' && d <= '9') {
			return fmt.Errorf("Expected digit after numeric sign instead of %#U", d)
		}
		takeDigits(l)
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		takeDigits(l)
	default:
		return fmt.Errorf("Expected digit after 'e' instead of %#U", r)
	}
	return nil
}

func takeJSONNumeric(l lexer) error {
	cur := l.take()
	switch cur {
	case '-':
		// Check digit immediately follows sign
		if d := l.peek(); !(d >= '0' && d <= '9') {
			return fmt.Errorf("Expected digit after dash instead of %#U", d)
		}
		takeDigits(l)
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		takeDigits(l)
	default:
		return fmt.Errorf("Expected digit or dash instead of %#U", cur)
	}

	// fraction or exponent
	cur = l.peek()
	switch cur {
	case '.':
		l.take()
		// Check digit immediately follows period
		if d := l.peek(); !(d >= '0' && d <= '9') {
			return fmt.Errorf("Expected digit after '.' instead of %#U", d)
		}
		takeDigits(l)
		if err := takeExponent(l); err != nil {
			return err
		}
	case 'e', 'E':
		if err := takeExponent(l); err != nil {
			return err
		}
	}

	return nil
}

func takeDigits(l lexer) {
	for {
		d := l.peek()
		if d >= '0' && d <= '9' {
			l.take()
		} else {
			break
		}
	}
}

// Only used at the very beginning of parsing. After that, the emit() function
// automatically skips whitespace.
func ignoreSpaceRun(l lexer) {
	for {
		r := l.peek()
		if r == ' ' || r == '\t' || r == '\r' || r == '\n' {
			l.take()
		} else {
			break
		}
	}
	l.ignore()
}

func takeExactSequence(l lexer, str []byte) bool {
	for _, r := range str {
		v := l.take()
		if v != int(r) {
			return false
		}
	}
	return true
}

func readerToArray(tr tokenReader) []Item {
	vals := make([]Item, 0)
	for {
		i, ok := tr.next()
		if !ok {
			break
		}
		v := *i
		s := make([]byte, len(v.val))
		copy(s, v.val)
		v.val = s
		vals = append(vals, v)
	}
	return vals
}

func findErrors(items []Item) (Item, bool) {
	for _, i := range items {
		if i.typ == lexError {
			return i, true
		}
	}
	return Item{}, false
}

func byteSlicesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func firstError(errors ...error) error {
	for _, e := range errors {
		if e != nil {
			return e
		}
	}
	return nil
}

func abs(x int) int {
	switch {
	case x < 0:
		return -x
	case x == 0:
		return 0 // return correctly abs(-0)
	}
	return x
}

//TODO: Kill the need for this
func getJsonTokenType(val []byte) (int, error) {
	if len(val) == 0 {
		return -1, errors.New("No Value")
	}
	switch val[0] {
	case '{':
		return jsonBraceLeft, nil
	case '"':
		return jsonString, nil
	case '[':
		return jsonBracketLeft, nil
	case 'n':
		return jsonNull, nil
	case 't', 'b':
		return jsonBool, nil
	case '-', '+', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return jsonNumber, nil
	default:
		return -1, errors.New("Unrecognized Json Value")
	}
}

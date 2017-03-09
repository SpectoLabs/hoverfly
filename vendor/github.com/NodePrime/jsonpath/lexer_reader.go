package jsonpath

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
)

type readerLexer struct {
	lex
	bufInput *bufio.Reader
	input    io.Reader
	pos      Pos
	nextByte int
	lexeme   *bytes.Buffer
}

func NewReaderLexer(rr io.Reader, initial stateFn) *readerLexer {
	l := readerLexer{
		input:    rr,
		bufInput: bufio.NewReader(rr),
		nextByte: noValue,
		lex:      newLex(initial),
		lexeme:   bytes.NewBuffer(make([]byte, 0, 100)),
	}
	return &l
}

func (l *readerLexer) take() int {
	if l.nextByte == noValue {
		l.peek()
	}

	nr := l.nextByte
	l.nextByte = noValue
	l.lexeme.WriteByte(byte(nr))
	return nr
}

func (l *readerLexer) takeString() error {
	cur := l.take()
	if cur != '"' {
		return fmt.Errorf("Expected \" as start of string instead of %#U", cur)
	}

	var previous byte
looper:
	for {
		curByte, err := l.bufInput.ReadByte()
		if err == io.EOF {
			return errors.New("Unexpected EOF in string")
		}
		l.lexeme.WriteByte(curByte)

		if curByte == '"' {
			if previous != '\\' {
				break looper
			} else {
				curByte, err = l.bufInput.ReadByte()
				if err == io.EOF {
					return errors.New("Unexpected EOF in string")
				}
				l.lexeme.WriteByte(curByte)
			}
		}

		previous = curByte
	}
	return nil
}

func (l *readerLexer) peek() int {
	if l.nextByte != noValue {
		return l.nextByte
	}

	r, err := l.bufInput.ReadByte()
	if err == io.EOF {
		l.nextByte = eof
		return eof
	}

	l.nextByte = int(r)
	return l.nextByte
}

func (l *readerLexer) emit(t int) {
	l.setItem(t, l.pos, l.lexeme.Bytes())
	l.pos += Pos(l.lexeme.Len())
	l.hasItem = true

	if t == lexEOF {
		// Do not capture eof character to match slice_lexer
		l.item.val = []byte{}
	}

	// Ignore whitespace after this token
	if l.nextByte == noValue {
		l.peek()
	}

	// ignore white space
	for l.nextByte != eof {
		if l.nextByte == ' ' || l.nextByte == '\t' || l.nextByte == '\r' || l.nextByte == '\n' {
			l.pos++
			r, err := l.bufInput.ReadByte()
			if err == io.EOF {
				l.nextByte = eof
			} else {
				l.nextByte = int(r)
			}
		} else {
			break
		}
	}
}

func (l *readerLexer) setItem(typ int, pos Pos, val []byte) {
	l.item.typ = typ
	l.item.pos = pos
	l.item.val = val
}

func (l *readerLexer) ignore() {
	l.pos += Pos(l.lexeme.Len())
	l.lexeme.Reset()
}

func (l *readerLexer) next() (*Item, bool) {
	l.lexeme.Reset()
	for {
		if l.currentStateFn == nil {
			break
		}

		l.currentStateFn = l.currentStateFn(l, &l.stack)

		if l.hasItem {
			l.hasItem = false
			return &l.item, true
		}
	}
	return &l.item, false
}

func (l *readerLexer) errorf(format string, args ...interface{}) stateFn {
	l.setItem(lexError, l.pos, []byte(fmt.Sprintf(format, args...)))
	l.lexeme.Truncate(0)
	l.hasItem = true
	return nil
}

func (l *readerLexer) reset() {
	l.bufInput.Reset(l.input)
	l.lexeme.Reset()
	l.nextByte = noValue
	l.pos = 0
	l.lex = newLex(l.initialState)
}

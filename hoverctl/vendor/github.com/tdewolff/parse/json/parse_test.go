package json // import "github.com/tdewolff/parse/json"

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func helperStringify(t *testing.T, input string) string {
	s := ""
	p := NewParser(bytes.NewBufferString(input))
	for i := 0; i < 10; i++ {
		tt, text := p.Next()
		if tt == ErrorGrammar {
			if p.Err() != nil {
				s += tt.String() + "('" + p.Err().Error() + "')"
			} else {
				s += tt.String() + "('')"
			}
			break
		} else if tt == WhitespaceGrammar {
			continue
		} else {
			s += tt.String() + "('" + string(text) + "') "
		}
	}
	return s
}

////////////////////////////////////////////////////////////////

type GTs []GrammarType

func TestGrammars(t *testing.T) {
	var grammarTests = []struct {
		json     string
		expected GTs
	}{
		{" \t\n\r", GTs{}}, // WhitespaceGrammar
		{"null", GTs{LiteralGrammar}},
		{"[]", GTs{StartArrayGrammar, EndArrayGrammar}},
		{"[15.2, 0.4, 5e9, -4E-3]", GTs{StartArrayGrammar, NumberGrammar, NumberGrammar, NumberGrammar, NumberGrammar, EndArrayGrammar}},
		{"[true, false, null]", GTs{StartArrayGrammar, LiteralGrammar, LiteralGrammar, LiteralGrammar, EndArrayGrammar}},
		{`["", "abc", "\"", "\\"]`, GTs{StartArrayGrammar, StringGrammar, StringGrammar, StringGrammar, StringGrammar, EndArrayGrammar}},
		{"{}", GTs{StartObjectGrammar, EndObjectGrammar}},
		{`{"a": "b", "c": "d"}`, GTs{StartObjectGrammar, StringGrammar, StringGrammar, StringGrammar, StringGrammar, EndObjectGrammar}},
		{`{"a": [1, 2], "b": {"c": 3}}`, GTs{StartObjectGrammar, StringGrammar, StartArrayGrammar, NumberGrammar, NumberGrammar, EndArrayGrammar, StringGrammar, StartObjectGrammar, StringGrammar, NumberGrammar, EndObjectGrammar, EndObjectGrammar}},
		{"[null,]", GTs{StartArrayGrammar, LiteralGrammar, EndArrayGrammar}},
		{"[\"x\\\x00y\", 0]", GTs{StartArrayGrammar, StringGrammar, NumberGrammar, EndArrayGrammar}},
	}
	for _, test := range grammarTests {
		stringify := helperStringify(t, test.json)
		p := NewParser(bytes.NewBufferString(test.json))
		i := 0
		for {
			tt, _ := p.Next()
			if tt == ErrorGrammar {
				assert.Equal(t, len(test.expected), i, "when error occurred we must be at the end in "+test.json)
				break
			} else if tt == WhitespaceGrammar {
				continue
			}
			assert.False(t, i >= len(test.expected), "index must not exceed expected grammar types size in "+stringify)
			if i < len(test.expected) {
				assert.Equal(t, test.expected[i], tt, "grammartypes must match at index "+strconv.Itoa(i)+" in "+stringify)
			}
			i++
		}
	}

	assert.Equal(t, "Whitespace", WhitespaceGrammar.String())
	assert.Equal(t, "Invalid(100)", GrammarType(100).String())
	assert.Equal(t, "Value", ValueState.String())
	assert.Equal(t, "ObjectKey", ObjectKeyState.String())
	assert.Equal(t, "ObjectValue", ObjectValueState.String())
	assert.Equal(t, "Array", ArrayState.String())
	assert.Equal(t, "Invalid(100)", State(100).String())
}

func TestGrammarsError(t *testing.T) {
	var grammarErrorTests = []struct {
		json     string
		expected error
	}{
		{"true, false", ErrBadComma},
		{"[true false]", ErrNoComma},
		{"]", ErrBadArrayEnding},
		{"}", ErrBadObjectEnding},
		{"{0: 1}", ErrBadObjectKey},
		{"{\"a\" 1}", ErrBadObjectDeclaration},
		{"1.", ErrNoComma},
		{"1e+", ErrNoComma},
		{`{"":"`, io.EOF},
		{"\"a\\", io.EOF},
	}
	for _, test := range grammarErrorTests {
		stringify := helperStringify(t, test.json)
		p := NewParser(bytes.NewBufferString(test.json))
		for {
			tt, _ := p.Next()
			if tt == ErrorGrammar {
				assert.Equal(t, test.expected, p.Err(), "parser must return error '"+test.expected.Error()+"' in "+stringify)
				break
			}
		}
	}
}

func TestStates(t *testing.T) {
	var stateTests = []struct {
		json     string
		expected []State
	}{
		{"null", []State{ValueState}},
		{"[null]", []State{ArrayState, ArrayState, ValueState}},
		{"{\"\":null}", []State{ObjectKeyState, ObjectValueState, ObjectKeyState, ValueState}},
	}
	for _, test := range stateTests {
		stringify := helperStringify(t, test.json)
		p := NewParser(bytes.NewBufferString(test.json))
		i := 0
		for {
			tt, _ := p.Next()
			state := p.State()
			if tt == ErrorGrammar {
				assert.Equal(t, len(test.expected), i, "when error occurred we must be at the end in "+stringify)
				break
			} else if tt == WhitespaceGrammar {
				continue
			}
			assert.False(t, i >= len(test.expected), "index must not exceed expected states size in "+stringify)
			if i < len(test.expected) {
				assert.Equal(t, test.expected[i], state, "states must match at index "+strconv.Itoa(i)+" in "+stringify)
			}
			i++
		}
	}
}

////////////////////////////////////////////////////////////////

func ExampleNewParser() {
	p := NewParser(bytes.NewBufferString(`{"key": 5}`))
	out := ""
	for {
		state := p.State()
		tt, data := p.Next()
		if tt == ErrorGrammar {
			break
		}
		out += string(data)
		if state == ObjectKeyState && tt != EndObjectGrammar {
			out += ":"
		}
		// not handling comma insertion
	}
	fmt.Println(out)
	// Output: {"key":5}
}

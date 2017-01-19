package buffer // import "github.com/tdewolff/buffer"

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/tdewolff/test"
)

func TestShifter(t *testing.T) {
	s := `Lorem ipsum dolor sit amet, consectetur adipiscing elit.`
	var z = NewShifter(bytes.NewBufferString(s))

	test.That(t, z.IsEOF(), "buffer must be fully in memory")
	test.Error(t, z.Err(), nil, "buffer is at EOF but must not return EOF until we reach that")
	test.That(t, z.Pos() == 0, "buffer must start at position 0")
	test.That(t, z.Peek(0) == 'L', "first character must be 'L'")
	test.That(t, z.Peek(1) == 'o', "second character must be 'o'")

	z.Move(1)
	test.That(t, z.Peek(0) == 'o', "must be 'o' at position 1")
	test.That(t, z.Peek(1) == 'r', "must be 'r' at position 1")
	z.MoveTo(6)
	test.That(t, z.Peek(0) == 'i', "must be 'i' at position 6")
	test.That(t, z.Peek(1) == 'p', "must be 'p' at position 7")

	test.Bytes(t, z.Bytes(), []byte("Lorem "), "buffered string must now read 'Lorem ' when at position 6")
	test.Bytes(t, z.Shift(), []byte("Lorem "), "shift must return the buffered string")
	test.That(t, z.Pos() == 0, "after shifting position must be 0")
	test.That(t, z.Peek(0) == 'i', "must be 'i' at position 0 after shifting")
	test.That(t, z.Peek(1) == 'p', "must be 'p' at position 1 after shifting")
	test.Error(t, z.Err(), nil, "error must be nil at this point")

	z.Move(len(s) - len("Lorem ") - 1)
	test.Error(t, z.Err(), nil, "error must be nil just before the end of the buffer")
	z.Skip()
	test.That(t, z.Pos() == 0, "after skipping position must be 0")
	z.Move(1)
	test.Error(t, z.Err(), io.EOF, "error must be EOF when past the buffer")
	z.Move(-1)
	test.Error(t, z.Err(), nil, "error must be nil just before the end of the buffer, even when it has been past the buffer")
}

func TestShifterSmall(t *testing.T) {
	s := `abcdefghi`
	z := NewShifterSize(test.NewPlainReader(bytes.NewBufferString(s)), 4)
	test.That(t, z.Peek(8) == 'i', "first character must be 'i' at position 8")
}

func TestShifterRunes(t *testing.T) {
	z := NewShifter(bytes.NewBufferString("aæ†\U00100000"))
	r, n := z.PeekRune(0)
	test.That(t, n == 1, "first character must be length 1")
	test.That(t, r == 'a', "first character must be rune 'a'")
	r, n = z.PeekRune(1)
	test.That(t, n == 2, "second character must be length 2")
	test.That(t, r == 'æ', "second character must be rune 'æ'")
	r, n = z.PeekRune(3)
	test.That(t, n == 3, "fourth character must be length 3")
	test.That(t, r == '†', "fourth character must be rune '†'")
	r, n = z.PeekRune(6)
	test.That(t, n == 4, "seventh character must be length 4")
	test.That(t, r == '\U00100000', "seventh character must be rune '\U00100000'")
}

func TestShifterZeroLen(t *testing.T) {
	var z = NewShifter(test.NewPlainReader(bytes.NewBufferString("")))
	test.That(t, z.Peek(0) == 0, "first character must yield error")
}

func TestShifterEmptyReader(t *testing.T) {
	var z = NewShifter(test.NewEmptyReader())
	test.That(t, z.Peek(0) == 0, "first character must yield error")
	test.That(t, z.IsEOF(), "empty reader must return EOF")
}

////////////////////////////////////////////////////////////////

func ExampleNewShifter() {
	b := bytes.NewBufferString("Lorem ipsum")
	z := NewShifter(b)
	for {
		c := z.Peek(0)
		if c == ' ' {
			break
		}
		z.Move(1)
	}
	fmt.Println(string(z.Shift()))
	// Output: Lorem
}

func ExampleShifter_PeekRune() {
	b := bytes.NewBufferString("† dagger") // † has a byte length of 3
	z := NewShifter(b)

	c, n := z.PeekRune(0)
	fmt.Println(string(c), n)
	// Output: † 3
}

func ExampleShifter_IsEOF() {
	b := bytes.NewBufferString("Lorem ipsum") // bytes.Buffer provides a Bytes function, NewShifter uses that and r.IsEOF() always returns true
	z := NewShifter(b)
	z.Move(5)

	lorem := z.Shift()
	if !z.IsEOF() { // required when io.Reader doesn't provide a Bytes function
		buf := make([]byte, len(lorem))
		copy(buf, lorem)
		lorem = buf
	}

	z.Peek(0) // might reallocate the internal buffer
	fmt.Println(string(lorem))
	// Output: Lorem
}

////////////////////////////////////////////////////////////////

func BenchmarkPeek(b *testing.B) {
	z := NewShifter(bytes.NewBufferString("Lorem ipsum"))
	for i := 0; i < b.N; i++ {
		j := i % 11
		z.Peek(j)
	}
}

var _c = 0
var _haystack = []byte("abcdefghijklmnopqrstuvwxyz")

func BenchmarkBytesEqual(b *testing.B) {
	for i := 0; i < b.N; i++ {
		j := i % (len(_haystack) - 3)
		if bytes.Equal([]byte("wxyz"), _haystack[j:j+4]) {
			_c++
		}
	}
}

func BenchmarkBytesEqual2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		j := i % (len(_haystack) - 3)
		if bytes.Equal([]byte{'w', 'x', 'y', 'z'}, _haystack[j:j+4]) {
			_c++
		}
	}
}

func BenchmarkBytesEqual3(b *testing.B) {
	match := []byte{'w', 'x', 'y', 'z'}
	for i := 0; i < b.N; i++ {
		j := i % (len(_haystack) - 3)
		if bytes.Equal(match, _haystack[j:j+4]) {
			_c++
		}
	}
}

func BenchmarkBytesEqual4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		j := i % (len(_haystack) - 3)
		if bytesEqual(_haystack[j:j+4], 'w', 'x', 'y', 'z') {
			_c++
		}
	}
}

func bytesEqual(stack []byte, match ...byte) bool {
	return bytes.Equal(stack, match)
}

func BenchmarkCharsEqual(b *testing.B) {
	for i := 0; i < b.N; i++ {
		j := i % (len(_haystack) - 3)
		if _haystack[j] == 'w' && _haystack[j+1] == 'x' && _haystack[j+2] == 'y' && _haystack[j+3] == 'z' {
			_c++
		}
	}
}

func BenchmarkCharsLoopEqual(b *testing.B) {
	match := []byte("wxyz")
	for i := 0; i < b.N; i++ {
		j := i % (len(_haystack) - 3)
		equal := true
		for k := 0; k < 4; k++ {
			if _haystack[j+k] != match[k] {
				equal = false
				break
			}
		}
		if equal {
			_c++
		}
	}
}

func BenchmarkCharsFuncEqual(b *testing.B) {
	match := []byte("wxyz")
	for i := 0; i < b.N; i++ {
		j := i % (len(_haystack) - 3)
		if at(match, _haystack[j:]) {
			_c++
		}
	}
}

func at(match []byte, stack []byte) bool {
	if len(stack) < len(match) {
		return false
	}
	for i, c := range match {
		if stack[i] != c {
			return false
		}
	}
	return true
}

package test

import "io"

// PlainReader implements an io.Reader and wraps over an existing io.Reader to hide other functions it implements.
type PlainReader struct {
	r io.Reader
}

// NewPlainReader returns a new PlainReader.
func NewPlainReader(r io.Reader) *PlainReader {
	return &PlainReader{r}
}

// Read implements the io.Reader interface.
func (r *PlainReader) Read(p []byte) (int, error) {
	return r.r.Read(p)
}

////////////////////////////////////////////////////////////////

// ErrorReader implements an io.Reader that will do N successive reads before it returns ErrPlain.
type ErrorReader struct {
	n int
}

// NewErrorReader returns a new ErrorReader.
func NewErrorReader(n int) *ErrorReader {
	return &ErrorReader{n}
}

// Read implements the io.Reader interface.
func (r *ErrorReader) Read(b []byte) (n int, err error) {
	if len(b) == 0 {
		return 0, nil
	}
	if r.n == 0 {
		return 0, ErrPlain
	}
	r.n--
	b[0] = '.'
	return 1, nil
}

////////////////////////////////////////////////////////////////

// InfiniteReader implements an io.Reader that will always read-in one character.
type InfiniteReader struct {
	n int
}

// NewInfiniteReader returns a new InfiniteReader.
func NewInfiniteReader(n int) *InfiniteReader {
	return &InfiniteReader{n}
}

// Read implements the io.Reader interface.
func (r *InfiniteReader) Read(b []byte) (n int, err error) {
	if len(b) == 0 {
		return 0, nil
	}
	b[0] = '.'
	return 1, nil
}

////////////////////////////////////////////////////////////////

// EmptyReader implements an io.Reader that will always return 0, nil.
type EmptyReader struct {
}

// NewEmptyReader returns a new EmptyReader.
func NewEmptyReader() *EmptyReader {
	return &EmptyReader{}
}

// Read implements the io.Reader interface.
func (r *EmptyReader) Read(b []byte) (n int, err error) {
	return 0, nil
}

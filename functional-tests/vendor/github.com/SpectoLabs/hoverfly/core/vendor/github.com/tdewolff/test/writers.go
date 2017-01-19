package test

// ErrorWriter implements an io.Writer that will do N successive writes before it returns ErrPlain.
type ErrorWriter struct {
	n int
}

// NewErrorWriter returns a new ErrorWriter.
func NewErrorWriter(n int) *ErrorWriter {
	return &ErrorWriter{n}
}

// Write implements the io.Writer interface.
func (w *ErrorWriter) Write(b []byte) (n int, err error) {
	if w.n == 0 {
		return 0, ErrPlain
	}
	w.n--
	return len(b), nil
}

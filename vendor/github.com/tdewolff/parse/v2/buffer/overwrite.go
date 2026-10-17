package buffer

import "io"

// OverwriteReaderWriter makes sure that any write always trails reads, so that we can safely overwrite the buffer.
type OverwriteReaderWriter struct {
	buf  []byte
	i, o int
}

func NewOverwriteReaderWriter(b []byte) *OverwriteReaderWriter {
	return &OverwriteReaderWriter{
		buf: b,
	}
}

func (r *OverwriteReaderWriter) Read(b []byte) (int, error) {
	if len(r.buf) <= r.i {
		return 0, io.EOF
	}
	n := copy(b, r.buf[r.i:])
	r.i += n
	return n, nil
}

func (r *OverwriteReaderWriter) Write(b []byte) (int, error) {
	if r.i < r.o+len(b) {
		panic("write further than read")
	} else if len(r.buf) < r.o+len(b) {
		panic("write overflows buffer")
	}
	n := copy(r.buf[r.o:], b)
	r.o += n
	return n, nil
}

func (r *OverwriteReaderWriter) WrittenLen() int {
	return r.o
}

func (r *OverwriteReaderWriter) WrittenBytes() []byte {
	return r.buf[:r.o]
}

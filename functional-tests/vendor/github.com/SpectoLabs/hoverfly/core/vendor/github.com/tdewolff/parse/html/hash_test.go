package html // import "github.com/tdewolff/parse/html"

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashTable(t *testing.T) {
	assert.Equal(t, Address, ToHash([]byte("address")), "'address' must resolve to Address")
	assert.Equal(t, "address", Address.String(), "Address must resolve to 'address'")
	assert.Equal(t, "accept-charset", Accept_Charset.String(), "Accept_Charset must resolve to 'accept-charset'")
	assert.Equal(t, Hash(0), ToHash([]byte("")), "empty string must resolve to zero")
	assert.Equal(t, "", Hash(0xffffff).String(), "Hash(0xffffff) must resolve to empty string")
	assert.Equal(t, Hash(0), ToHash([]byte("iter")), "'iter' must resolve to zero")
	assert.Equal(t, Hash(0), ToHash([]byte("test")), "'test' must resolve to zero")
}

////////////////////////////////////////////////////////////////

var result int

// naive scenario
func BenchmarkCompareBytes(b *testing.B) {
	var r int
	val := []byte("span")
	for n := 0; n < b.N; n++ {
		if bytes.Equal(val, []byte("span")) {
			r++
		}
	}
	result = r
}

// using-atoms scenario
func BenchmarkFindAndCompareAtom(b *testing.B) {
	var r int
	val := []byte("span")
	for n := 0; n < b.N; n++ {
		if ToHash(val) == Span {
			r++
		}
	}
	result = r
}

// using-atoms worst-case scenario
func BenchmarkFindAtomCompareBytes(b *testing.B) {
	var r int
	val := []byte("zzzz")
	for n := 0; n < b.N; n++ {
		if h := ToHash(val); h == 0 && bytes.Equal(val, []byte("zzzz")) {
			r++
		}
	}
	result = r
}

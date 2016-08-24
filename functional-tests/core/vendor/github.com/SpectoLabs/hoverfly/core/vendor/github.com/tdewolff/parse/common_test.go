package parse // import "github.com/tdewolff/parse"

import (
	"encoding/base64"
	"mime"
	"testing"

	"github.com/tdewolff/test"
)

func TestParseNumber(t *testing.T) {
	var numberTests = []struct {
		number   string
		expected int
	}{
		{"5", 1},
		{"0.51", 4},
		{"0.5e-99", 7},
		{"0.5e-", 3},
		{"+50.0", 5},
		{".0", 2},
		{"0.", 1},
		{"", 0},
		{"+", 0},
		{".", 0},
		{"a", 0},
	}
	for _, tt := range numberTests {
		n := Number([]byte(tt.number))
		test.That(t, n == tt.expected)
	}
}

func TestParseDimension(t *testing.T) {
	var dimensionTests = []struct {
		dimension    string
		expectedNum  int
		expectedUnit int
	}{
		{"5px", 1, 2},
		{"5px ", 1, 2},
		{"5%", 1, 1},
		{"5em", 1, 2},
		{"px", 0, 0},
		{"1", 1, 0},
		{"1~", 1, 0},
	}
	for _, tt := range dimensionTests {
		num, unit := Dimension([]byte(tt.dimension))
		test.That(t, num == tt.expectedNum, "number")
		test.That(t, unit == tt.expectedUnit, "unit")
	}
}

func TestMediatype(t *testing.T) {
	var mediatypeTests = []struct {
		mediatype        string
		expectedMimetype string
		expectedParams   map[string]string
	}{
		{"text/plain", "text/plain", nil},
		{"text/plain;charset=US-ASCII", "text/plain", map[string]string{"charset": "US-ASCII"}},
		{" text/plain  ; charset = US-ASCII ", "text/plain", map[string]string{"charset": "US-ASCII"}},
		{" text/plain  a", "text/plain", nil},
		{"text/plain;base64", "text/plain", map[string]string{"base64": ""}},
		{"text/plain;inline=;base64", "text/plain", map[string]string{"inline": "", "base64": ""}},
	}
	for _, tt := range mediatypeTests {
		mimetype, _ := Mediatype([]byte(tt.mediatype))
		test.Bytes(t, mimetype, []byte(tt.expectedMimetype), "mimetype")
		//test.That(t, params == tt.expectedParams, "parameters") // TODO
	}
}

func TestParseDataURI(t *testing.T) {
	var dataURITests = []struct {
		dataURI          string
		expectedMimetype string
		expectedData     string
		expectedErr      error
	}{
		{"www.domain.com", "", "", ErrBadDataURI},
		{"data:,", "text/plain", "", nil},
		{"data:text/xml,", "text/xml", "", nil},
		{"data:,text", "text/plain", "text", nil},
		{"data:;base64,dGV4dA==", "text/plain", "text", nil},
		{"data:image/svg+xml,", "image/svg+xml", "", nil},
		{"data:;base64,()", "", "", base64.CorruptInputError(0)},
	}
	for _, tt := range dataURITests {
		mimetype, data, err := DataURI([]byte(tt.dataURI))
		test.Bytes(t, mimetype, []byte(tt.expectedMimetype), "mimetype")
		test.Bytes(t, data, []byte(tt.expectedData), "data")
		test.Error(t, err, tt.expectedErr)
	}
}

func TestParseQuoteEntity(t *testing.T) {
	var quoteEntityTests = []struct {
		quoteEntity   string
		expectedQuote byte
		expectedN     int
	}{
		{"&#34;", '"', 5},
		{"&#039;", '\'', 6},
		{"&#x0022;", '"', 8},
		{"&#x27;", '\'', 6},
		{"&quot;", '"', 6},
		{"&apos;", '\'', 6},
		{"&gt;", 0x00, 0},
		{"&amp;", 0x00, 0},
	}
	for _, tt := range quoteEntityTests {
		quote, n := QuoteEntity([]byte(tt.quoteEntity))
		test.That(t, quote == tt.expectedQuote, "quote", quote, "must equal", tt.expectedQuote)
		test.That(t, n == tt.expectedN, "quote length", n, "must equal", tt.expectedN)
	}
}

////////////////////////////////////////////////////////////////

func BenchmarkParseMediatypeStd(b *testing.B) {
	mediatype := "text/plain"
	for i := 0; i < b.N; i++ {
		mime.ParseMediaType(mediatype)
	}
}

func BenchmarkParseMediatypeParamStd(b *testing.B) {
	mediatype := "text/plain;inline=1"
	for i := 0; i < b.N; i++ {
		mime.ParseMediaType(mediatype)
	}
}

func BenchmarkParseMediatypeParamsStd(b *testing.B) {
	mediatype := "text/plain;charset=US-ASCII;language=US-EN;compression=gzip;base64"
	for i := 0; i < b.N; i++ {
		mime.ParseMediaType(mediatype)
	}
}

func BenchmarkParseMediatypeParse(b *testing.B) {
	mediatype := []byte("text/plain")
	for i := 0; i < b.N; i++ {
		Mediatype(mediatype)
	}
}

func BenchmarkParseMediatypeParamParse(b *testing.B) {
	mediatype := []byte("text/plain;inline=1")
	for i := 0; i < b.N; i++ {
		Mediatype(mediatype)
	}
}

func BenchmarkParseMediatypeParamsParse(b *testing.B) {
	mediatype := []byte("text/plain;charset=US-ASCII;language=US-EN;compression=gzip;base64")
	for i := 0; i < b.N; i++ {
		Mediatype(mediatype)
	}
}

package css // import "github.com/tdewolff/parse/css"

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashTable(t *testing.T) {
	assert.Equal(t, Font, ToHash([]byte("font")), "'font' must resolve to hash.Font")
	assert.Equal(t, "font", Font.String(), "hash.Font must resolve to 'font'")
	assert.Equal(t, "margin-left", Margin_Left.String(), "hash.Margin_Left must resolve to 'margin-left'")
	assert.Equal(t, Hash(0), ToHash([]byte("")), "empty string must resolve to zero")
	assert.Equal(t, "", Hash(0xffffff).String(), "Hash(0xffffff) must resolve to empty string")
	assert.Equal(t, Hash(0), ToHash([]byte("fonts")), "'fonts' must resolve to zero")
}

package svg // import "github.com/tdewolff/parse/svg"

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashTable(t *testing.T) {
	assert.Equal(t, Svg, ToHash([]byte("svg")), "'svg' must resolve to hash.Svg")
	assert.Equal(t, Width, ToHash([]byte("width")), "'width' must resolve to hash.Width")
	assert.Equal(t, "svg", Svg.String(), "hash.Svg must resolve to 'svg'")
	assert.Equal(t, Hash(0), ToHash([]byte("")), "empty string must resolve to zero")
	assert.Equal(t, "", Hash(0xffffff).String(), "Hash(0xffffff) must resolve to empty string")
	assert.Equal(t, Hash(0), ToHash([]byte("svgs")), "'svgs' must resolve to zero")
	assert.Equal(t, Hash(0), ToHash([]byte("uopi")), "'uopi' must resolve to zero")
}

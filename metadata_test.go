package hoverfly

import (
	"testing"
)

func TestSetMetadataKey(t *testing.T) {

	metaBucket := GetRandomName(10)
	md := NewBoltDBMetadata(TestDB, metaBucket)

	md.Set("foo", "bar")

	expect(t, md.Get("foo"), "bar")
}

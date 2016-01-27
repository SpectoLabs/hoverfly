package main

import (
	"testing"
)

func TestVirtualizeInc(t *testing.T) {
	counter := NewModeCounter()

	counter.Count(VirtualizeMode)

	expect(t, int(counter.counterVirtualize.Count()), 1)
}

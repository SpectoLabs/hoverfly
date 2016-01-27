package main

import (
	"testing"
)

func TestVirtualizeInc(t *testing.T) {
	counter := NewModeCounter()

	counter.Count(VirtualizeMode)

	expect(t, int(counter.counterVirtualize.Count()), 1)
}

func TestCaptureInc(t *testing.T) {
	counter := NewModeCounter()

	counter.Count(CaptureMode)

	expect(t, int(counter.counterCapture.Count()), 1)
}

func TestModifyInc(t *testing.T) {
	counter := NewModeCounter()

	counter.Count(ModifyMode)

	expect(t, int(counter.counterModify.Count()), 1)
}

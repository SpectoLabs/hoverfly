package metrics

import (
	"testing"
)

func TestSimulateInc(t *testing.T) {
	counter := NewModeCounter([]string{"name"})

	counter.Count("name")

	count := counter.Counters["name"].Count()

	if count != 1 {
		t.Fatalf("Expected counter to have size %v but was %v", 1, counter)
	}
}

func TestFlush(t *testing.T) {
	counter := NewModeCounter([]string{"name"})

	counter.Counters["name"].Inc(1)

	fl := counter.Flush()

	count := fl.Counters["name"]

	if count != 1 {
		t.Fatalf("Expected counter to have size %v but was %v", 1, counter)
	}
}

package metrics

import (
	. "github.com/onsi/gomega"
	"testing"
)

func TestSimulateInc(t *testing.T) {
	RegisterTestingT(t)
	counter := NewModeCounter([]string{"name"})

	counter.Count("name")

	count := counter.Counters["name"].Count()

	Expect(count).To(Equal(int64(1)))
}

func TestFlush(t *testing.T) {
	RegisterTestingT(t)
	counter := NewModeCounter([]string{"name"})

	counter.Counters["name"].Inc(1)

	fl := counter.Flush()

	count := fl.Counters["name"]

	Expect(count).To(Equal(int64(1)))
}

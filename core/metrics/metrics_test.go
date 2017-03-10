package metrics_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/metrics"
	. "github.com/onsi/gomega"
)

func TestSimulateInc(t *testing.T) {
	RegisterTestingT(t)
	counter := metrics.NewModeCounter([]string{"name"})

	counter.Count("name")

	count := counter.Counters["name"].Count()

	Expect(count).To(Equal(int64(1)))
}

func TestFlush(t *testing.T) {
	RegisterTestingT(t)
	counter := metrics.NewModeCounter([]string{"name"})

	counter.Counters["name"].Inc(1)

	fl := counter.Flush()

	count := fl.Counters["name"]

	Expect(count).To(Equal(int64(1)))
}

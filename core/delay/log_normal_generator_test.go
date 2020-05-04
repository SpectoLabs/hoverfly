package delay

import (
	. "github.com/onsi/gomega"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/stat"
	"sort"
	"testing"
)

const tolerance = 10

func TestLogNormalGenerator_GenerateDelay(t *testing.T) {
	RegisterTestingT(t)


	opts := LogNormalDelayOptions{Median: 500, Mean: 1000, Max: 20000, Min: 100}
	gen := NewLogNormalGenerator(opts)

	const n = 1e5
	sample := make([]float64, n)
	for i := range sample {
		sample[i] = float64(gen.GenerateDelay())
	}
	sort.Float64s(sample)

	actualMean := stat.Mean(sample, nil)
	Expect(opts.Mean).To(BeNumerically("~", actualMean, tolerance), "mean diff must be less than tolerance")

	actualMedian := stat.Quantile(0.5, stat.Empirical, sample, nil)
	Expect(opts.Median).To(BeNumerically("~", actualMedian, tolerance), "median diff must be less than tolerance")

	Expect(opts.Max).To(BeNumerically(">=", floats.Max(sample)), "max generated value must be less or equal than `max`")
	Expect(opts.Min).To(BeNumerically("<=", floats.Min(sample)), "min generated value must be less or equal than `min`")

}

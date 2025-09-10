package delay

import (
	"sort"
	"testing"

	. "github.com/onsi/gomega"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/stat"
)

const tolerance = 50

func TestLogNormalGenerator_GenerateDelay(t *testing.T) {
	RegisterTestingT(t)

	min, max, mean, median := 100, 20000, 1000, 500

	gen := NewLogNormalGenerator(min, max, mean, median)

	const n = 1e5
	sample := make([]float64, n)
	for i := range sample {
		sample[i] = float64(gen.GenerateDelay())
	}
	sort.Float64s(sample)

	actualMean := stat.Mean(sample, nil)
	Expect(mean).To(BeNumerically("~", actualMean, tolerance), "mean diff must be less than tolerance")

	actualMedian := stat.Quantile(0.5, stat.Empirical, sample, nil)
	Expect(median).To(BeNumerically("~", actualMedian, tolerance), "median diff must be less than tolerance")

	Expect(max).To(BeNumerically(">=", floats.Max(sample)), "max generated value must be less or equal than `max`")
	Expect(min).To(BeNumerically("<=", floats.Min(sample)), "min generated value must be less or equal than `min`")

}

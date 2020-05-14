package delay

import (
	"errors"
	"gonum.org/v1/gonum/stat/distuv"
	"math"
)

func ValidateLogNormalDelayOptions(min int, max int, mean int, median int) error {
	if max < 0 || min < 0 {
		return errors.New("Config error - delay min and max can't be less than 0")
	}
	if mean <= 0 || median <= 0 {
		return errors.New("Config error - delay mean and median params can't be less or equals 0")
	}

	if max != 0 {
		if max < min {
			return errors.New("Config error - min delay must be less than max one")
		}
		if mean > max {
			return errors.New("Config error - mean delay can't be greather than max one")
		}
		if median > max {
			return errors.New("Config error - median delay can't be and greather than max one")
		}
	}

	if min != 0 {
		if mean < min {
			return errors.New("Config error - mean delay can't be less than min one")
		}
		if median < min {
			return errors.New("Config error - median delay can't be less than min one")
		}
	}

	if median > mean {
		return errors.New("Config error - mean delay can't be less than median one")
	}

	return nil
}

type LogNormalGenerator struct {
	Min  int
	Max  int
	dist *distuv.LogNormal
}

func NewLogNormalGenerator(min int, max int, mean int, median int) *LogNormalGenerator {
	mu := math.Log(float64(median))
	sigma := math.Sqrt(2 * (math.Log(float64(mean)) - mu))
	dist := &distuv.LogNormal{
		Mu:    mu,
		Sigma: sigma,
	}
	return &LogNormalGenerator{
		Min:  min,
		Max:  max,
		dist: dist,
	}
}

func (g *LogNormalGenerator) GenerateDelay() int {
	delay := int(g.dist.Rand())
	if g.Max != 0 && delay > g.Max {
		return g.Max
	}
	if g.Min != 0 && delay < g.Min {
		return g.Min
	}
	return delay
}

package delay

import (
	"errors"
	"gonum.org/v1/gonum/stat/distuv"
	"math"
)

type LogNormalDelayOptions struct {
	Min    int `json:"min"`
	Max    int `json:"max"`
	Mean   int `json:"mean"`
	Median int `json:"median"`
}

func ValidateLogNormalDelayOptions(opts LogNormalDelayOptions) error {
	if opts.Max < 0 || opts.Min < 0 {
		return errors.New("Config error - delay min and max can't be less than 0")
	}
	if opts.Mean <= 0 || opts.Median <= 0 {
		return errors.New("Config error - delay mean and median params can't be less or equals 0")
	}

	if opts.Max != 0 {
		if opts.Max < opts.Min {
			return errors.New("Config error - min delay must be less than max one")
		}
		if opts.Mean > opts.Max {
			return errors.New("Config error - mean delay can't be greather than max one")
		}
		if opts.Median > opts.Max {
			return errors.New("Config error - median delay can't be and greather than max one")
		}
	}

	if opts.Min != 0 {
		if opts.Mean < opts.Min {
			return errors.New("Config error - mean delay can't be less than min one")
		}
		if opts.Median < opts.Min {
			return errors.New("Config error - median delay can't be less than min one")
		}
	}

	if opts.Median > opts.Mean {
		return errors.New("Config error - mean delay can't be less than median one")
	}

	return nil
}

type LogNormalGenerator struct {
	Min  int
	Max  int
	dist *distuv.LogNormal
}

func NewLogNormalGenerator(opts LogNormalDelayOptions) *LogNormalGenerator {
	mu := math.Log(float64(opts.Median))
	sigma := math.Sqrt(2 * (math.Log(float64(opts.Mean)) - mu))
	dist := &distuv.LogNormal{
		Mu:    mu,
		Sigma: sigma,
	}
	return &LogNormalGenerator{
		Min:  opts.Min,
		Max:  opts.Max,
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

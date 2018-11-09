package delay

import (
	"gonum.org/v1/gonum/stat/distuv"
	"math"
)

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

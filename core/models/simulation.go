package models

import (
	"reflect"
)

type Simulation struct {
	matchingPairs  []RequestMatcherResponsePair
	ResponseDelays ResponseDelays
}

func NewSimulation() *Simulation {

	return &Simulation{
		matchingPairs:  []RequestMatcherResponsePair{},
		ResponseDelays: &ResponseDelayList{},
	}
}

func (this *Simulation) AddRequestMatcherResponsePair(pair *RequestMatcherResponsePair) {
	var duplicate bool
	for _, savedPair := range this.matchingPairs {
		duplicate = reflect.DeepEqual(pair.RequestMatcher, savedPair.RequestMatcher)
		if duplicate {
			break
		}
	}
	if !duplicate {
		this.matchingPairs = append(this.matchingPairs, *pair)
	}
}

func (this *Simulation) GetMatchingPairs() []RequestMatcherResponsePair {
	return this.matchingPairs
}

func (this *Simulation) DeleteMatchingPairs() {
	var pairs []RequestMatcherResponsePair
	this.matchingPairs = pairs
}

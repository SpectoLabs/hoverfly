package models

import (
	"reflect"
)

type Simulation struct {
	MatchingPairs  []RequestMatcherResponsePair
	ResponseDelays ResponseDelays
}

func NewSimulation() *Simulation {

	return &Simulation{
		MatchingPairs:  []RequestMatcherResponsePair{},
		ResponseDelays: &ResponseDelayList{},
	}
}

func (this *Simulation) AddRequestMatcherResponsePair(pair *RequestMatcherResponsePair) {
	var duplicate bool
	for _, savedPair := range this.MatchingPairs {
		duplicate = reflect.DeepEqual(pair.RequestMatcher, savedPair.RequestMatcher)
		if duplicate {
			break
		}
	}
	if !duplicate {
		this.MatchingPairs = append(this.MatchingPairs, *pair)
	}
}

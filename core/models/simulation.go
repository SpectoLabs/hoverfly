package models

import (
	"fmt"
	"reflect"
	"strconv"
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

func (this *Simulation) AddPair(pair *RequestMatcherResponsePair) {
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

func (this *Simulation) AddPairInSequence(pair *RequestMatcherResponsePair, state map[string]string) {
	var duplicate bool

	updates := map[int]RequestMatcherResponsePair{}

	var counter int

	for i, savedPair := range this.matchingPairs {
		fmt.Println("loop")

		pairNoState := pair.RequestMatcher
		pairNoState.RequiresState = nil

		savedPairNoState := savedPair.RequestMatcher
		savedPairNoState.RequiresState = nil

		duplicate = reflect.DeepEqual(pairNoState, savedPairNoState)
		if duplicate {
			fmt.Println("dup")
			counter = counter + 1

			if savedPair.RequestMatcher.RequiresState == nil {
				savedPair.RequestMatcher.RequiresState = map[string]string{}
			}

			if pair.RequestMatcher.RequiresState == nil {
				pair.RequestMatcher.RequiresState = map[string]string{}
			}

			sequenceState := savedPair.RequestMatcher.RequiresState["sequence"]
			if sequenceState == "" {
				sequenceState = "1"
				state["sequence"] = "1"
				savedPair.RequestMatcher.RequiresState["sequence"] = sequenceState
				updates[i] = savedPair
			}

		}
	}

	for i, updatedPair := range updates {
		this.matchingPairs[i] = updatedPair
	}

	fmt.Println(counter)
	if counter != 0 {
		pair.RequestMatcher.RequiresState["sequence"] = strconv.Itoa(counter + 1)
	}

	this.matchingPairs = append(this.matchingPairs, *pair)
}

func (this *Simulation) GetMatchingPairs() []RequestMatcherResponsePair {
	return this.matchingPairs
}

func (this *Simulation) DeleteMatchingPairs() {
	var pairs []RequestMatcherResponsePair
	this.matchingPairs = pairs
}

package models

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/SpectoLabs/hoverfly/core/state"
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

func (this *Simulation) AddPairInSequence(pair *RequestMatcherResponsePair, state *state.State) {
	var duplicate bool

	updates := map[int]RequestMatcherResponsePair{}

	var counter int
	sequenceKey := "sequence:0"

	for i, savedPair := range this.matchingPairs {

		pairNoState := pair.RequestMatcher
		pairNoState.RequiresState = nil

		savedPairNoState := savedPair.RequestMatcher
		savedPairNoState.RequiresState = nil

		duplicate = reflect.DeepEqual(pairNoState, savedPairNoState)
		if duplicate {
			counter = counter + 1

			if savedPair.RequestMatcher.RequiresState == nil {
				savedPair.RequestMatcher.RequiresState = map[string]string{}
			}

			if savedPair.Response.TransitionsState == nil {
				savedPair.Response.TransitionsState = map[string]string{}
			}

			if pair.RequestMatcher.RequiresState == nil {
				pair.RequestMatcher.RequiresState = map[string]string{}
			}
			sequenceKey = getNewSequenceKey(state)
			for key, _ := range savedPair.RequestMatcher.RequiresState {
				if strings.Contains(key, "sequence:") {
					sequenceKey = key
					break
				}
			}

			sequenceState := savedPair.RequestMatcher.RequiresState[sequenceKey]
			nextSequenceState := ""
			if sequenceState == "" {
				sequenceState = "1"
				nextSequenceState = "2"
				state.State[sequenceKey] = "1"

			} else {
				currentSequenceState, _ := strconv.Atoi(sequenceState)
				nextSequenceState = strconv.Itoa(currentSequenceState + 1)
			}
			savedPair.RequestMatcher.RequiresState[sequenceKey] = sequenceState
			savedPair.Response.TransitionsState[sequenceKey] = nextSequenceState
			updates[i] = savedPair
		}
	}

	for i, updatedPair := range updates {
		this.matchingPairs[i] = updatedPair
	}

	if counter != 0 {
		pair.RequestMatcher.RequiresState[sequenceKey] = strconv.Itoa(counter + 1)
	}

	this.matchingPairs = append(this.matchingPairs, *pair)
}

func getNewSequenceKey(state *state.State) string {
	returnKey := ""
	i := 0
	for returnKey == "" {
		tempKey := fmt.Sprintf("sequence:%v", i)
		if state.GetState(tempKey) == "" {
			returnKey = tempKey
		} else {
			i = i + 1
		}
	}
	return returnKey
}

func (this *Simulation) GetMatchingPairs() []RequestMatcherResponsePair {
	return this.matchingPairs
}

func (this *Simulation) DeleteMatchingPairs() {
	var pairs []RequestMatcherResponsePair
	this.matchingPairs = pairs
}
